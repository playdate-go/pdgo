//go:build gc.playdate

package runtime

import (
	"unsafe"
)

// FinalizerFunc is the function signature for object finalizers.
type FinalizerFunc = func(unsafe.Pointer)

// finalizers maps mangled object keys to their finalizer functions.
// Replaces the old [128]fixed array — now unbounded.
//
// Keys are stored mangled (finalizerKey) rather than as raw pointers: the GC
// conservatively scans every global (including this map's buckets) for values
// that look like heap pointers. Raw pointer keys would therefore mark every
// registered object on every cycle, pinning them forever and leaking until OOM.
// XOR-mapping is a bijection, so lookups still work; mangled SRAM heap
// addresses (0x2000_0000–0x2003_FFFF) land at 0x7A5A_xxxx, outside any region
// the scanner treats as the heap.
var finalizers map[uintptr]func(unsafe.Pointer)

// TEMPORARY device diagnostics for the corruption hunt: registration and
// sweep-time lookup counters, printed in the per-cycle GC log line.
var (
	gcFinalAdds uint32
	gcFinalHits uint32 // sweep-time lookups that found a finalizer
	gcFinalMiss uint32 // sweep-time lookups on hasFinal blocks that found none
)

// finalizerKey bijectively maps an object pointer to a non-pointer-looking
// map key (and back).
func finalizerKey(p unsafe.Pointer) uintptr {
	return uintptr(p) ^ 0x5A5A0000
}

func finalizersAdd(obj unsafe.Pointer, fn func(unsafe.Pointer)) {
	if finalizers == nil {
		finalizers = make(map[uintptr]func(unsafe.Pointer), 64)
	}
	gcFinalAdds++
	finalizers[finalizerKey(obj)] = fn
}

// finalizersGet returns the finalizer AND removes it from the table.
// Use for sweep-time invocation. For explicit free, use finalizersDrop.
func finalizersGet(obj unsafe.Pointer) func(unsafe.Pointer) {
	if finalizers == nil {
		return nil
	}
	fn, ok := finalizers[finalizerKey(obj)]
	if !ok {
		gcFinalMiss++
		return nil
	}
	gcFinalHits++
	delete(finalizers, finalizerKey(obj))
	return fn
}

// finalizersDrop removes a finalizer without returning or calling it.
func finalizersDrop(obj unsafe.Pointer) {
	if finalizers == nil {
		return
	}
	delete(finalizers, finalizerKey(obj))
}

// SetFinalizer sets a finalizer for an object. When the object is garbage
// collected, the finalizer will be called with the object's pointer before
// the memory is recycled.
//
// Like standard Go, the finalizer may be typed: func(*T). (The simulator
// runtime *requires* typed finalizers, so pdgo call sites use them.) A single
// pointer argument shares one machine-level calling convention with
// unsafe.Pointer, so typed finalizers are reinterpreted internally.
func SetFinalizer(obj interface{}, finalizer interface{}) {
	objPtr := extractInterfacePointer(obj)
	if objPtr == nil {
		panic("runtime.SetFinalizer: obj is not a pointer")
	}
	fn, ok := finalizer.(func(unsafe.Pointer))
	if !ok {
		fn = reinterpretFinalizer(finalizer)
	}
	header := (*gcHeader)(unsafe.Pointer(uintptr(objPtr) - gcHeaderSize))
	header.color = setHasFinal(header.color, true)
	finalizersAdd(objPtr, fn)
}

// reinterpretFinalizer converts a typed finalizer value (func(*T)) to the
// internal func(unsafe.Pointer) representation.
//
// TinyGo func type descriptors do not encode signatures (reflect's NumIn is
// unimplemented), so we can only check that the value is a function, not the
// shape of its argument. Passing any signature other than a single
// pointer-sized argument is undefined behavior.
func reinterpretFinalizer(finalizer interface{}) func(unsafe.Pointer) {
	type ifaceHeader struct {
		typ  unsafe.Pointer
		data unsafe.Pointer
	}
	h := (*ifaceHeader)(unsafe.Pointer(&finalizer))
	if h.typ == nil {
		panic("runtime.SetFinalizer: finalizer must be func(unsafe.Pointer) or func(*T)")
	}

	// TinyGo typecode layout (see internal/reflectlite/type.go): the low 2
	// bits of a typecode are a pointer-type tag; otherwise the typecode
	// points at a descriptor whose first byte packs the kind in the low 5
	// bits plus flags.
	const (
		kindMask  = 31
		flagNamed = 32
		kindFunc  = 24 // reflectlite Kind value for Func
	)
	tc := uintptr(h.typ)
	if tc&0b11 != 0 {
		panic("runtime.SetFinalizer: finalizer must be func(unsafe.Pointer) or func(*T)")
	}

	// Named types store their underlying type in elem; unwrap to reach the
	// actual kind. Mirrors reflectlite's elemType.
	type namedType struct {
		meta      uint8
		numMethod uint16
		ptrTo     unsafe.Pointer
		elem      unsafe.Pointer
	}
	for i := 0; i < 8; i++ {
		meta := *(*uint8)(unsafe.Pointer(tc))
		if meta&flagNamed == 0 {
			if meta&kindMask != kindFunc || h.data == nil {
				panic("runtime.SetFinalizer: finalizer must be func(unsafe.Pointer) or func(*T)")
			}
			// TinyGo represents a func value as a two-word {context, code}
			// pair and boxes it in an interface as a POINTER to that pair
			// (a static $pack constant for non-capturing closures). The Go
			// func value itself IS the pair, so load it through the
			// interface's data word. Loading AT &h.data (the stack slot)
			// captured the pair pointer as the code half instead: the
			// finalizer call then blx'd into rodata and hardfaulted with
			// an undefined instruction (device: cfsr UNDEFINSTR).
			return *(*func(unsafe.Pointer))(h.data)
		}
		tc = uintptr((*namedType)(unsafe.Pointer(tc)).elem)
	}
	panic("runtime.SetFinalizer: finalizer must be func(unsafe.Pointer) or func(*T)")
}

//go:noinline
func extractInterfacePointer(i interface{}) unsafe.Pointer {
	type iface struct {
		typ  uintptr
		data unsafe.Pointer
	}
	return (*iface)(unsafe.Pointer(&i)).data
}

// runFinalizerSafe invokes fn with panic recovery. A panicking finalizer
// must not halt the device — log and continue.
func runFinalizerSafe(fn func(unsafe.Pointer), p unsafe.Pointer) {
	defer func() {
		if r := recover(); r != nil {
			// On bare-metal Playdate, write to serial via putchar.
			const msg = "GC: finalizer panicked\n"
			for i := 0; i < len(msg); i++ {
				putchar(msg[i])
			}
		}
	}()
	fn(p)
}
