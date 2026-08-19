//go:build gc.leaking

// Leaking GC — no-op. Allocations are never reclaimed.
// This is the fallback if the conservative GC (gc.playdate) shows
// problems in production. Switch via playdate.json's "gc" field:
// change "gc": "playdate" to "gc": "leaking" and update build-tags
// from "gc.playdate" to "gc.leaking".

package runtime

import "unsafe"

const needsStaticHeap = false

var (
	gcTotalAlloc uint64
	gcMallocs    uint64
	gcFrees      uint64
)

// leakAlign aligns size to pointer-width (matches gc_helpers.align).
func leakAlign(size uintptr) uintptr {
	if size == 0 {
		return 0
	}
	a := unsafe.Sizeof(uintptr(0))
	return (size + a - 1) &^ (a - 1)
}

//go:noinline
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
	size = leakAlign(size)
	gcTotalAlloc += uint64(size)
	gcMallocs++
	ptr := _cgo_pd_realloc(nil, size)
	if ptr == nil {
		runtimePanic("out of memory")
	}
	memzero(ptr, size)
	return ptr
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	size = leakAlign(size)
	newPtr := _cgo_pd_realloc(ptr, size)
	if newPtr == nil && size > 0 {
		runtimePanic("out of memory")
	}
	return newPtr
}

func free(ptr unsafe.Pointer) {
	if ptr != nil {
		_cgo_pd_realloc(ptr, 0)
		gcFrees++
	}
}

func GC()                             {}
func SetFinalizer(obj, finalizer any) {}

func ReadMemStats(m *MemStats) {
	m.HeapAlloc = gcTotalAlloc
	m.TotalAlloc = gcTotalAlloc
	m.Mallocs = gcMallocs
	m.Frees = gcFrees
}

func initHeap()          {}
func setHeapEnd(uintptr) {}
