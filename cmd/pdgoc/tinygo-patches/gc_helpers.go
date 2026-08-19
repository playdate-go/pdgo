//go:build gc.playdate

// Pure helpers for the Playdate GC runtime.
//
// KEEP IN SYNC with cmd/pdgoc/gcpure/helpers.go.
// This file is package runtime (untestable directly); the gcpure copy is
// the testable version. Bodies are identical; only names and package differ.
// (HeaderSize is gcpure-only; the runtime uses unsafe.Sizeof(gcHeader{}) directly.)
// (gcpure exports Align; the runtime copy is gcAlign to avoid colliding with
// the stock TinyGo runtime's align() in arch_cortexm.go.)

package runtime

import "unsafe"

const numSizeClasses = 8

const (
	colorWhite uint8 = 0
	colorGray  uint8 = 1
	colorBlack uint8 = 2
)
const colorMask uint8 = 3

const (
	hasFinalFlag   uint8 = 1 << 2
	inFreeListFlag uint8 = 1 << 3
	pinFlag        uint8 = 1 << 4
)

func sizeClassOf(size uintptr) uint8 {
	thresholds := [numSizeClasses - 1]uintptr{16, 32, 64, 128, 256, 512, 1024}
	for i, t := range thresholds {
		if size <= t {
			return uint8(i)
		}
	}
	return numSizeClasses - 1
}

func bucketSize(sc uint8, size uintptr) uintptr {
	if sc == numSizeClasses-1 {
		return size
	}
	buckets := [numSizeClasses - 1]uintptr{48, 64, 96, 160, 288, 544, 1056}
	return buckets[sc]
}

// gcAlign is named with a gc prefix to avoid colliding with the stock
// TinyGo runtime's align() in arch_cortexm.go (pointer alignment, 8-byte).
func gcAlign(size uintptr) uintptr {
	if size == 0 {
		return 0
	}
	a := unsafe.Sizeof(uintptr(0))
	return (size + a - 1) &^ (a - 1)
}

func colorOf(c uint8) uint8         { return c & colorMask }
func setColor(c, color uint8) uint8 { return (c &^ colorMask) | (color & colorMask) }

func setHasFinal(c uint8, v bool) uint8 {
	if v {
		return c | hasFinalFlag
	}
	return c &^ hasFinalFlag
}
func hasFinal(c uint8) bool { return c&hasFinalFlag != 0 }

func setInFreeList(c uint8, v bool) uint8 {
	if v {
		return c | inFreeListFlag
	}
	return c &^ inFreeListFlag
}
func inFreeList(c uint8) bool { return c&inFreeListFlag != 0 }

// Pinned blocks could not be covered by the object bitmap (SDK refused to
// grow it); they are never marked, swept, or finalized — leaked, not freed.
func setPinned(c uint8, v bool) uint8 {
	if v {
		return c | pinFlag
	}
	return c &^ pinFlag
}
func isPinned(c uint8) bool { return c&pinFlag != 0 }
