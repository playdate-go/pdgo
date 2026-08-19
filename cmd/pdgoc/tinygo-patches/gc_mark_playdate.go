//go:build gc.playdate

package runtime

import "unsafe"

// Growable mark stack for iterative marking (never silently drops entries)
var gcMarkStack []uintptr

// Debug counter
var gcDebugMarked uint32

// markRoots scans a memory range conservatively for heap pointers
func markRoots(start, end uintptr) {
	if start >= end {
		return
	}

	ptrSize := unsafe.Sizeof(uintptr(0))

	if end-start > 1024*1024 {
		end = start + 1024*1024
	}

	start = (start + ptrSize - 1) &^ (ptrSize - 1)

	for addr := start; addr+ptrSize <= end; addr += ptrSize {
		ptr := *(*uintptr)(unsafe.Pointer(addr))
		if ptr != 0 {
			markObject(ptr)
		}
	}
}

// markObject marks a heap object as reachable.
// Uses the offset-encoded side bitmap for O(1) headerOf lookup.
func markObject(ptr uintptr) {
	h := headerOf(ptr)
	if h == nil || isPinned(h.color) {
		return
	}
	if colorOf(h.color) != colorWhite {
		return
	}
	h.color = setColor(h.color, colorGray)
	gcDebugMarked++
	gcMarkStack = append(gcMarkStack, h.userStart, h.size)
}

func processWorkQueue() {
	ptrSize := unsafe.Sizeof(uintptr(0))
	for len(gcMarkStack) >= 2 {
		size := gcMarkStack[len(gcMarkStack)-1]
		start := gcMarkStack[len(gcMarkStack)-2]
		// Zero the popped slots before truncating: the slice's backing array
		// is itself a heap object the GC scans conservatively, so a stale
		// (start,size) tail would keep pinning dead objects forever.
		n := len(gcMarkStack)
		gcMarkStack[n-1] = 0
		gcMarkStack[n-2] = 0
		gcMarkStack = gcMarkStack[:n-2]
		end := start + size
		start = (start + ptrSize - 1) &^ (ptrSize - 1)
		for addr := start; addr+ptrSize <= end; addr += ptrSize {
			ptr := *(*uintptr)(unsafe.Pointer(addr))
			if ptr != 0 {
				markObject(ptr)
			}
		}
	}
}
