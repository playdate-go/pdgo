//go:build gc.playdate

// Optimized Hybrid Conservative GC for Playdate
// - Uses offset-encoded side bitmap for O(1) object lookups during marking
// - Minimizes STW pauses by reducing markObject complexity to O(1)
// - Supports finalizers for automatic cleanup of C resources (Bitmap, Sound, etc.)
// - Uses Playdate SDK's realloc for memory allocation

package runtime

import (
	"unsafe"
)

const needsStaticHeap = false

// GC header for each allocation
type gcHeader struct {
	size      uintptr
	color     uint8 // white|gray|black + hasFinal + inFreeList + pinned flags
	sizeClass uint8
	_pad      [2]byte
	prev      *gcHeader // previous in alloc list (for O(1) unlink)
	next      *gcHeader // next in alloc list
	nextFree  *gcHeader // free-list link (only meaningful while the block sits on a free list)
	userStart uintptr
}

// GC state machine
type gcStateT uint8

const (
	gcStateIdle gcStateT = iota
	gcStateMarking
	gcStateFinalizing
)

// GC state variables
var (
	gcHeaderSize  uintptr                   // computed at init
	gcAllocList   *gcHeader                 // head of global allocation list
	gcAllocCount  uint32                    // number of active allocations
	gcTotalAlloc  uint64                    // total bytes ever allocated
	gcMallocs     uint64                    // total allocation calls
	gcFrees       uint64                    // total bytes freed
	gcNumGC       uint32                    // number of GC cycles
	gcLastGCSize  uint64                    // heap size after last GC
	gcLastPauseNs int64                     // duration of last GC cycle in nanoseconds
	gcEnabled     bool                      // GC enabled after initial threshold
	gcStateVal    gcStateT                  // prevent recursive GC
	freeLists     [numSizeClasses]*gcHeader // per-size-class free lists

	// gcAllocsSinceGC bounds GC frequency by work done since the last cycle,
	// not by live count: a game with more than gcMaxAllocs LIVE objects would
	// otherwise trigger a full GC on every single allocation.
	gcAllocsSinceGC uint32

	// TEMPORARY device diagnostics for the corruption hunt: observed SDK
	// heap address range and blocks pinned because the bitmap could not be
	// grown to cover them.
	gcPinnedCount uint32
	gcMinAddr     uintptr // lowest userStart ever returned by the SDK
	gcMaxAddr     uintptr // highest userStart+size ever returned
	gcPinnedWarn  bool    // one-shot console warning
)

// GC tuning parameters - adjusted for smoother operation
const (
	gcTriggerRatio   = 3     // trigger GC when heap grows 3x
	gcMinTriggerSize = 65536 // minimum heap size before GC (64KB)
	gcMaxAllocs      = 4096  // max allocations since last GC before forcing one
)

// freeListPush adds a header to the appropriate size-classed free list
func freeListPush(h *gcHeader) {
	sc := h.sizeClass
	h.color = setInFreeList(setColor(h.color, colorWhite), true)
	h.nextFree = freeLists[sc]
	freeLists[sc] = h
}

// freeListPop removes and returns a header from the given size-class free list
func freeListPop(sc uint8) *gcHeader {
	h := freeLists[sc]
	if h == nil {
		return nil
	}
	freeLists[sc] = h.nextFree
	h.color = colorWhite // reset all bits (color + flags)
	return h
}

// allocListInsert adds a header to the front of the global alloc list. O(1).
func allocListInsert(h *gcHeader) {
	h.prev = nil
	h.next = gcAllocList
	if gcAllocList != nil {
		gcAllocList.prev = h
	}
	gcAllocList = h
	gcAllocCount++
}

// allocListUnlink removes a header from the alloc list. O(1). No-op if the
// header is not in the list.
func allocListUnlink(h *gcHeader) {
	if h.prev != nil {
		h.prev.next = h.next
	} else if gcAllocList == h {
		gcAllocList = h.next
	} else {
		return // not in list
	}
	if h.next != nil {
		h.next.prev = h.prev
	}
	h.prev = nil
	h.next = nil
	if gcAllocCount > 0 {
		gcAllocCount--
	}
}

// flushFreeLists drains all free-lists back to the SDK allocator
func flushFreeLists() {
	for sc := uint8(0); sc < numSizeClasses; sc++ {
		for h := freeLists[sc]; h != nil; {
			next := h.nextFree
			// Capture size before the release: the SDK allocator may
			// scribble metadata over the freed header.
			size := h.size
			objectMapClear(h.userStart, size)
			_cgo_pd_realloc(unsafe.Pointer(h), 0)
			gcFrees += uint64(gcHeaderSize + size)
			h = next
		}
		freeLists[sc] = nil
	}
}

// maybeTriggerGC triggers a GC if conditions are met
func maybeTriggerGC() {
	if gcStateVal != gcStateIdle {
		return
	}
	if !gcEnabled {
		return
	}
	if gcAllocsSinceGC > gcMaxAllocs {
		GC()
		return
	}
	currentHeap := gcTotalAlloc - gcFrees
	if currentHeap > gcMinTriggerSize && currentHeap > gcLastGCSize*gcTriggerRatio {
		GC()
	}
}

//go:noinline
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
	if gcHeaderSize == 0 {
		gcHeaderSize = unsafe.Sizeof(gcHeader{})
	}
	size = gcAlign(size)
	sc := sizeClassOf(size)
	// Round the user size up so the total block (header + user) equals the
	// class bucket. Every block in a size class then has identical capacity,
	// so free-list reuse can never hand out a block smaller than the request.
	// Without this, popping an undersized block for a larger same-class
	// allocation overran the next heap block's gcHeader — silent heap
	// corruption that made the GC free and memzero live objects.
	if bucket := bucketSize(sc, size); bucket > size+gcHeaderSize {
		size = bucket - gcHeaderSize
	}

	// Fast path: reuse from free-list. Class 7 never enters free lists (see
	// sweep): exact-size blocks have heterogeneous capacities, so reuse for
	// a larger same-class request could overrun the block.
	if sc < numSizeClasses-1 {
		if h := freeListPop(sc); h != nil {
			// TEMPORARY device diagnostics: catch a forged/stale header at
			// the moment of reuse rather than at the crash it would cause.
			if h.size < size {
				print("GCBUG: bad freelist pop hdr=")
				print(int(uint32(uintptr(unsafe.Pointer(h)))))
				print(" size=")
				print(int(uint32(h.size)))
				print(" want=")
				println(int(uint32(size)))
				runtimePanic("gc heap corrupt (freelist)")
			}
			h.size = size
			h.sizeClass = sc
			h.prev = nil // clear stale prev from free-list
			allocListInsert(h)
			objectMapMark(h.userStart, size) // re-mark bitmap (was cleared at push)
			userData := unsafe.Pointer(h.userStart)
			memzero(userData, size)
			gcTotalAlloc += uint64(gcHeaderSize + size)
			gcMallocs++
			gcAllocsSinceGC++
			if !gcEnabled && gcTotalAlloc > gcMinTriggerSize {
				gcEnabled = true
			}
			maybeTriggerGC()
			return userData
		}
	}

	// Slow path: fresh allocation from SDK
	totalSize := gcHeaderSize + size
	gcTotalAlloc += uint64(totalSize)
	gcMallocs++
	gcAllocsSinceGC++

	ptr := _cgo_pd_realloc(nil, totalSize)
	if ptr == nil {
		flushFreeLists()
		ptr = _cgo_pd_realloc(nil, totalSize)
		if ptr == nil {
			runtimePanic("out of memory")
		}
	}

	header := (*gcHeader)(ptr)
	header.size = size
	header.color = colorWhite
	header.sizeClass = sc
	header.nextFree = nil
	header.userStart = uintptr(ptr) + gcHeaderSize
	// Insert BEFORE covering: if the bitmap must grow, its rebuild walks the
	// alloc list, and the new block has to be in it to get marked.
	allocListInsert(header)

	// Track the observed SDK heap range. If the bitmap cannot be grown to
	// cover this block, pin it — an uncoverable block can never be marked,
	// so sweeping it would reclaim a live object.
	if gcMinAddr == 0 || header.userStart < gcMinAddr {
		gcMinAddr = header.userStart
	}
	if header.userStart+size > gcMaxAddr {
		gcMaxAddr = header.userStart + size
	}
	if objectMapCover(header.userStart, size) {
		objectMapMark(header.userStart, size)
	} else {
		header.color = setPinned(header.color, true)
		gcPinnedCount++
		if !gcPinnedWarn {
			gcPinnedWarn = true
			print("GCBUG: cannot grow object bitmap, pinning addr=")
			print(int(uint32(header.userStart)))
			print(" size=")
			println(int(uint32(size)))
		}
	}

	userData := unsafe.Pointer(header.userStart)
	memzero(userData, size)

	if !gcEnabled && gcTotalAlloc > gcMinTriggerSize {
		gcEnabled = true
	}
	maybeTriggerGC()
	return userData
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	if gcHeaderSize == 0 {
		gcHeaderSize = unsafe.Sizeof(gcHeader{})
	}
	if ptr == nil {
		return alloc(size, nil)
	}
	if size == 0 {
		free(ptr)
		return nil
	}
	size = gcAlign(size)
	header := (*gcHeader)(unsafe.Pointer(uintptr(ptr) - gcHeaderSize))
	if size <= header.size {
		return ptr
	}
	newPtr := alloc(size, nil)
	if newPtr == nil {
		return nil
	}
	memcpy(newPtr, ptr, header.size)
	free(ptr)
	return newPtr
}

func free(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	header := (*gcHeader)(unsafe.Pointer(uintptr(ptr) - gcHeaderSize))
	if hasFinal(header.color) {
		finalizersDrop(ptr)
	}
	// Capture size before the release paths: the SDK allocator may scribble
	// metadata over the freed header.
	size := header.size
	allocListUnlink(header)
	objectMapClear(header.userStart, size) // no-op if uncovered
	if header.sizeClass == numSizeClasses-1 || isPinned(header.color) {
		// Class 7 blocks are exact-size with heterogeneous capacities, and
		// pinned blocks are uncoverable: neither may be recycled through
		// the free list. The caller declared this object dead, so returning
		// it straight to the SDK is safe.
		_cgo_pd_realloc(unsafe.Pointer(header), 0)
	} else {
		freeListPush(header)
	}
	gcFrees += uint64(gcHeaderSize + size)
}

// GC performs a garbage collection cycle
func GC() {
	if gcStateVal != gcStateIdle {
		return
	}
	gcStateVal = gcStateMarking
	defer func() { gcStateVal = gcStateIdle }()

	gcNumGC++
	gcDebugMarked = 0
	gcDebugSwept = 0
	gcAllocsSinceGC = 0
	// Fresh slice each cycle. The previous backing array may have been swept
	// and REUSED by a game allocation: arrays grown mid-processWorkQueue are
	// white at sweep time (the globals scan ran before they existed), so
	// sweep frees them while this global still points at the recycled block.
	// Keeping [:0] would append mark pairs straight through live game data.
	gcMarkStack = nil

	start := ticks() // ms granularity; pauses < 1ms report as 0
	gcMarkReachable()
	processWorkQueue()
	sweep()

	// TEMPORARY device diagnostics for the corruption hunt. One console
	// line per cycle: if marked << live, root scanning is broken; pin>0
	// means the SDK refused to grow the bitmap (those blocks leak);
	// lo/hi show the observed SDK heap address range.
	print("GC ")
	print(int(gcNumGC))
	print(": marked=")
	print(int(gcDebugMarked))
	print(" swept=")
	print(int(gcDebugSwept))
	print(" live=")
	print(int(gcAllocCount))
	print(" heap=")
	print(int(uint32(gcTotalAlloc - gcFrees)))
	print(" pin=")
	print(int(gcPinnedCount))
	print(" lo=")
	print(int(uint32(gcMinAddr)))
	print(" hi=")
	print(int(uint32(gcMaxAddr)))
	println()
	gcVerifyHeap()

	gcLastPauseNs = int64(ticks() - start)

	gcLastGCSize = gcTotalAlloc - gcFrees
}

// ReadPlaydateMemStats returns GC statistics for the pdgo.Memory API.
// Exported so pdgo can call via runtime.ReadPlaydateMemStats().
func ReadPlaydateMemStats() (heapAlloc uint64, numGC uint32, liveObjects uint32, lastPauseNs int64, totalAlloc uint64, totalFrees uint64) {
	heapAlloc = gcTotalAlloc - gcFrees
	numGC = gcNumGC
	liveObjects = gcAllocCount
	lastPauseNs = gcLastPauseNs
	totalAlloc = gcTotalAlloc
	totalFrees = gcFrees
	return
}

// ResetGCStats zeroes the GC counters for benchmarking. Exported for
// pdgo.Memory.ResetStats(). Preserves the live-heap invariant: gcTotalAlloc
// is rebased to the current live heap so that heapAlloc stays consistent with
// gcAllocCount. gcLastGCSize is also rebased so the trigger heuristic
// recalibrates from the current heap size.
func ResetGCStats() {
	liveHeap := gcTotalAlloc - gcFrees
	gcNumGC = 0
	gcTotalAlloc = liveHeap
	gcFrees = 0
	gcMallocs = 0
	gcLastPauseNs = 0
	gcLastGCSize = liveHeap
}

func ReadMemStats(m *MemStats) {
	heapInuse := gcTotalAlloc - gcFrees
	m.HeapIdle = 0
	m.HeapInuse = heapInuse
	m.HeapReleased = 0
	m.HeapSys = m.HeapInuse + m.HeapIdle
	m.GCSys = 0
	m.TotalAlloc = gcTotalAlloc
	m.Mallocs = gcMallocs
	m.Frees = gcFrees
	m.Sys = gcTotalAlloc
	m.HeapAlloc = heapInuse
	m.Alloc = m.HeapAlloc
}

func initHeap() {
	gcHeaderSize = unsafe.Sizeof(gcHeader{})
	gcAllocList = nil
	gcAllocCount = 0
	gcTotalAlloc = 0
	gcMallocs = 0
	gcFrees = 0
	gcNumGC = 0
	gcLastGCSize = 0
	gcLastPauseNs = 0
	gcEnabled = false
	gcAllocsSinceGC = 0
	gcMarkStack = gcMarkStack[:0]
	finalizers = nil
	gcPinnedCount = 0
	gcMinAddr = 0
	gcMaxAddr = 0
	gcPinnedWarn = false
	objectMap = nil
	objectMapBase = 0
	for i := range freeLists {
		freeLists[i] = nil
	}
}

func setHeapEnd(newHeapEnd uintptr) {}
