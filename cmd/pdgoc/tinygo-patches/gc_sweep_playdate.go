//go:build gc.playdate

package runtime

import "unsafe"

// Debug counter
var gcDebugSwept uint32

// sweep frees all unmarked allocations, calling finalizers if present
func sweep() {
	current := gcAllocList

	for current != nil {
		next := current.next

		// Pinned blocks are not covered by the object bitmap and can never
		// be marked, so they are always white here; recycling them would
		// reclaim LIVE objects. Pinning leaks them instead of corrupting
		// the heap. Their finalizers never run; explicit free() still
		// releases them to the SDK.
		if isPinned(current.color) {
			current = next
			continue
		}

		if colorOf(current.color) == colorWhite && !isFresh(current.color) {
			objPtr := unsafe.Pointer(current.userStart)
			// Capture size now: the class-7 branch below releases the block
			// to the SDK, whose allocator may scribble metadata over the
			// header before gcFrees reads it back.
			size := current.size

			// Call finalizer if present BEFORE recycling
			if hasFinal(current.color) {
				if fn := finalizersGet(objPtr); fn != nil {
					runFinalizerSafe(fn, objPtr)
				}
			}

			// Unlink from alloc list (O(1) now via doubly-linked list)
			allocListUnlink(current)

			if current.sizeClass == numSizeClasses-1 {
				// Class 7 blocks are exact-size with heterogeneous
				// capacities; free-list reuse for a larger same-class
				// request would hand out an undersized block and overrun
				// its neighbor. Release straight back to the SDK instead.
				// Clear before the release: the memory leaves our control
				// and the SDK may scribble over the header.
				objectMapClear(current.userStart, size)
				_cgo_pd_realloc(unsafe.Pointer(current), 0)
			} else {
				// Push to free-list (recycle — no SDK call). No bitmap
				// clear: reuse re-marks the FULL bucket span (alloc rounds
				// the size up to the bucket before marking), and markObject
				// ignores inFreeList blocks, so stale entries are inert.
				// The sweep-time clear was the dominant per-block cost
				// (device: 10-15µs per 264-byte block).
				current.color = colorWhite // reset all flags
				freeListPush(current)
			}
			gcFrees += uint64(gcHeaderSize + size)
			gcDebugSwept++
		} else {
			// Kept (marked, or fresh from this cycle's allocations): reset
			// the color for the next cycle and drop the fresh protection —
			// the block has survived the allocation cycle and must now
			// stand on its own reachability.
			current.color = setFresh(setColor(current.color, colorWhite), false)
		}

		current = next
	}

	// Note: gcAllocCount is already decremented by allocListUnlink.
	// gcDebugSwept tracks how many were swept for debug reporting.
}

// gcVerifyHeap checks post-sweep invariants. TEMPORARY device diagnostics
// for the corruption hunt — every check that fires names a subsystem:
//   - a block in both the alloc list and a free list  → double push/pop
//   - size not equal to the class bucket capacity     → forged/overrun header
//   - list walk count != gcAllocCount                 → list corruption
func gcVerifyHeap() {
	seen := 0
	for h := gcAllocList; h != nil; h = h.next {
		seen++
		if inFreeList(h.color) {
			print("GCBUG: block in both lists: ")
			println(int(uint32(uintptr(unsafe.Pointer(h)))))
		}
		want := h.size
		if h.sizeClass < numSizeClasses-1 {
			want = bucketSize(h.sizeClass, h.size) - gcHeaderSize
		}
		if h.size != want || h.userStart != uintptr(unsafe.Pointer(h))+gcHeaderSize {
			print("GCBUG: bad header at ")
			println(int(uint32(uintptr(unsafe.Pointer(h)))))
		}
	}
	if seen != int(gcAllocCount) {
		print("GCBUG: list count ")
		print(seen)
		print(" != gcAllocCount ")
		println(int(gcAllocCount))
	}
}
