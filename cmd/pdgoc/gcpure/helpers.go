// Package gcpure contains pure-Go helpers for the Playdate GC.
//
// KEEP IN SYNC with cmd/pdgoc/tinygo-patches/gc_helpers.go.
// The runtime copy uses unexported names and lives in package runtime.
// This package exists so the helpers can be unit-tested on host.
package gcpure

import "unsafe"

// Number of size classes. Class 7 is "large/exact" (no bucketing).
const NumSizeClasses = 8

// Color values (tri-color marking). Stored in the low 2 bits of gcHeader.color.
const (
	ColorWhite uint8 = 0 // unreachable (default at alloc, dead at sweep)
	ColorGray  uint8 = 1 // on mark stack, not yet scanned
	ColorBlack uint8 = 2 // scanned, reachable
)
const colorMask uint8 = 3

// Flag bits in gcHeader.color (above the 2 color bits).
const (
	HasFinalFlag   uint8 = 1 << 2
	InFreeListFlag uint8 = 1 << 3
	PinFlag        uint8 = 1 << 4
)

// SizeClassOf returns the bucket index for a user-data size.
// Class 7 is "exact" — no bucketing, used for sizes above 1024.
func SizeClassOf(size uintptr) uint8 {
	thresholds := [NumSizeClasses - 1]uintptr{16, 32, 64, 128, 256, 512, 1024}
	for i, t := range thresholds {
		if size <= t {
			return uint8(i)
		}
	}
	return NumSizeClasses - 1
}

// BucketSize returns the total allocation size (header + user) for a given
// class and user size. For class 7, returns the exact size.
//
// Bucket sizes account for the 20-byte header and 16-byte alignment.
func BucketSize(sc uint8, size uintptr) uintptr {
	if sc == NumSizeClasses-1 {
		return size
	}
	buckets := [NumSizeClasses - 1]uintptr{48, 64, 96, 160, 288, 544, 1056}
	return buckets[sc]
}

// Align rounds size up to the GC's natural alignment (pointer-size).
// On the 64-bit host this is 8; on the 32-bit Playdate target it is 4.
// Sizes of 0 stay 0.
func Align(size uintptr) uintptr {
	if size == 0 {
		return 0
	}
	a := unsafe.Sizeof(uintptr(0))
	return (size + a - 1) &^ (a - 1)
}

// ColorOf extracts the color from a gcHeader.color byte.
func ColorOf(c uint8) uint8 {
	return c & colorMask
}

// SetColor returns a color byte with the color portion replaced.
func SetColor(c uint8, color uint8) uint8 {
	return (c &^ colorMask) | (color & colorMask)
}

// SetHasFinal sets or clears the has-finalizer flag bit, preserving color and other flags.
func SetHasFinal(c uint8, v bool) uint8 {
	if v {
		return c | HasFinalFlag
	}
	return c &^ HasFinalFlag
}

// HasFinal reports whether the has-finalizer flag bit is set.
func HasFinal(c uint8) bool { return c&HasFinalFlag != 0 }

// SetInFreeList sets or clears the in-free-list flag bit, preserving color and other flags.
func SetInFreeList(c uint8, v bool) uint8 {
	if v {
		return c | InFreeListFlag
	}
	return c &^ InFreeListFlag
}

// InFreeList reports whether the in-free-list flag bit is set.
func InFreeList(c uint8) bool { return c&InFreeListFlag != 0 }

// Pinned blocks could not be covered by the object bitmap; they are never
// marked, swept, or finalized — leaked, not freed.
func SetPinned(c uint8, v bool) uint8 {
	if v {
		return c | PinFlag
	}
	return c &^ PinFlag
}
func IsPinned(c uint8) bool { return c&PinFlag != 0 }

// HeaderSize returns the size of the END-STATE gcHeader (after Task 7 removes
// the transitional userStart field). The current transitional gcHeader in
// tinygo-patches/gc_playdate.go is larger (24 bytes on target) because it
// still carries userStart; fakeHeader here models the final 20-byte shape so
// that BucketSize and other consumers compute against the target layout.
//
// fakeHeader uses fixed-width uint32 fields (not uintptr/pointers) so the
// size is identical on 32-bit (Playdate target) and 64-bit (host) builds —
// the real gcHeader uses uintptr/*gcHeader which are 4 bytes on target.
func HeaderSize() uintptr {
	type fakeHeader struct {
		size      uint32
		color     uint8
		sizeClass uint8
		_pad      [2]byte
		prevAlloc uint32
		nextAlloc uint32
		nextFree  uint32
	}
	return unsafe.Sizeof(fakeHeader{})
}
