//go:build gc.playdate

package runtime

import "unsafe"

// The Playdate OS serves game-heap allocations from the QSPI-mapped game
// flash region (observed on device: userStart values from 0x9001xxxx and
// upward), NOT from the 256KB SRAM window a fixed bitmap assumed. The object
// bitmap is therefore itself allocated from the SDK heap and grown/rebased on
// demand: when an allocation falls outside the covered range, a bigger
// bitmap is allocated, zeroed, and repopulated exactly from the allocation
// list (the source of truth). If the SDK cannot provide a bigger bitmap, the
// allocation is flagged pinned and never reclaimed (leak, not corruption).
var (
	objectMap     []byte  // 1 byte per 4-byte heap slot
	objectMapBase uintptr // heap address covered by objectMap[0]
)

// Extend the covered range by at least this much to amortize rebuilds.
const objectMapMinGrow = 128 << 10

// objectMapCovers reports whether the whole block [addr, addr+size) lies
// inside the tracked window.
func objectMapCovers(addr, size uintptr) bool {
	if objectMap == nil || size == 0 {
		return false
	}
	end := objectMapBase + uintptr(len(objectMap))*4
	return addr >= objectMapBase && addr < end && addr+size <= end
}

// objectMapCover ensures the bitmap covers [addr, addr+size), growing and
// rebuilding it from the allocation list when needed. Returns false when the
// SDK cannot provide a bigger bitmap; the caller then pins the allocation.
// Only call from alloc paths (allocates via _cgo_pd_realloc directly, never
// through alloc(), so it cannot re-enter the GC).
func objectMapCover(addr, size uintptr) bool {
	if objectMapCovers(addr, size) {
		return true
	}
	base := addr
	end := addr + size + objectMapMinGrow
	if objectMap != nil {
		oldEnd := objectMapBase + uintptr(len(objectMap))*4
		if objectMapBase < base {
			base = objectMapBase
		}
		if oldEnd > end {
			end = oldEnd
		}
	}
	base &^= 3
	n := (end - base + 3) >> 2
	buf := _cgo_pd_realloc(nil, n)
	if buf == nil {
		return false
	}
	memzero(buf, n)
	old := objectMap
	objectMap = unsafe.Slice((*byte)(buf), int(n))
	objectMapBase = base
	// Rebuild exactly from the alloc list. Free-list blocks are absent from
	// the list, so their slots correctly stay 0.
	for h := gcAllocList; h != nil; h = h.next {
		if !isPinned(h.color) && objectMapCovers(h.userStart, h.size) {
			objectMapMark(h.userStart, h.size)
		}
	}
	if old != nil {
		_cgo_pd_realloc(unsafe.Pointer(&old[0]), 0)
	}
	return objectMapCovers(addr, size)
}

// objectMapMark tags all 4-byte slots covering [start, start+size).
func objectMapMark(start, size uintptr) {
	if !objectMapCovers(start, size) {
		return
	}
	startIdx := (start - objectMapBase) >> 2
	numSlots := (size + 3) >> 2
	for i := uintptr(0); i < numSlots; i++ {
		offset := i + 1
		if offset > 255 {
			offset = 255
		}
		objectMap[startIdx+i] = byte(offset)
	}
}

// objectMapClear zeroes all 4-byte slots covering [start, start+size).
func objectMapClear(start, size uintptr) {
	if !objectMapCovers(start, size) {
		return
	}
	startIdx := (start - objectMapBase) >> 2
	numSlots := (size + 3) >> 2
	for i := uintptr(0); i < numSlots; i++ {
		objectMap[startIdx+i] = 0
	}
}

// objectMapStart returns the userStart of the object containing addr, or 0.
func objectMapStart(addr uintptr) uintptr {
	if !objectMapCovers(addr, 1) {
		return 0
	}
	idx := (addr - objectMapBase) >> 2
	for {
		b := objectMap[idx]
		if b == 0 {
			return 0
		}
		if b == 1 {
			return objectMapBase + idx*4
		}
		if b < 255 {
			return objectMapBase + (idx-uintptr(b)+1)*4
		}
		if idx < 254 {
			return 0
		}
		idx -= 254
	}
}

// headerOf returns the gcHeader for any pointer into an object (including
// interior pointers), or nil if the pointer is not within a known object.
func headerOf(userPtr uintptr) *gcHeader {
	start := objectMapStart(userPtr)
	if start == 0 {
		return nil
	}
	return (*gcHeader)(unsafe.Pointer(start - gcHeaderSize))
}
