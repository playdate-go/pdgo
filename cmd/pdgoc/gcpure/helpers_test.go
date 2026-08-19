package gcpure

import "testing"
import "unsafe"

func TestSizeClassOf(t *testing.T) {
	cases := []struct {
		size uintptr
		want uint8
	}{
		{0, 0}, {1, 0}, {16, 0},
		{17, 1}, {32, 1},
		{33, 2}, {64, 2},
		{65, 3}, {128, 3},
		{129, 4}, {256, 4},
		{257, 5}, {512, 5},
		{513, 6}, {1024, 6},
		{1025, 7}, {1 << 20, 7},
	}
	for _, c := range cases {
		got := SizeClassOf(c.size)
		if got != c.want {
			t.Errorf("SizeClassOf(%d) = %d, want %d", c.size, got, c.want)
		}
	}
}

func TestBucketSize(t *testing.T) {
	cases := []struct {
		sc   uint8
		size uintptr
		want uintptr
	}{
		{0, 1, 48}, {0, 16, 48},
		{1, 17, 64}, {1, 32, 64},
		{2, 33, 96}, {2, 64, 96},
		{3, 65, 160}, {3, 128, 160},
		{4, 129, 288}, {4, 256, 288},
		{5, 257, 544}, {5, 512, 544},
		{6, 513, 1056}, {6, 1024, 1056},
		{7, 1025, 1025}, {7, 1 << 20, 1 << 20}, // exact
	}
	for _, c := range cases {
		got := BucketSize(c.sc, c.size)
		if got != c.want {
			t.Errorf("BucketSize(%d, %d) = %d, want %d", c.sc, c.size, got, c.want)
		}
	}
}

func TestAlign(t *testing.T) {
	if Align(0) != 0 {
		t.Errorf("Align(0) = %d, want 0", Align(0))
	}
	// Align is at least 8-byte (must hold a pointer on 32-bit ARM)
	got := Align(1)
	if got != 8 {
		t.Errorf("Align(1) = %d, want 8", got)
	}
	if Align(8) != 8 {
		t.Errorf("Align(8) = %d, want 8", Align(8))
	}
	if Align(9) != 16 {
		t.Errorf("Align(9) = %d, want 16", Align(9))
	}
}

func TestColorFlagOps(t *testing.T) {
	var c uint8 = ColorWhite
	if ColorOf(c) != ColorWhite {
		t.Errorf("ColorOf(white) wrong")
	}
	c = SetColor(c, ColorGray)
	if ColorOf(c) != ColorGray {
		t.Errorf("SetColor gray failed")
	}
	c = SetColor(c, ColorBlack)
	if ColorOf(c) != ColorBlack {
		t.Errorf("SetColor black failed")
	}
	c = ColorWhite
	c = SetHasFinal(c, true)
	if !HasFinal(c) {
		t.Errorf("SetHasFinal true failed")
	}
	c = SetHasFinal(c, false)
	if HasFinal(c) {
		t.Errorf("SetHasFinal false failed")
	}
	c = SetInFreeList(c, true)
	if !InFreeList(c) {
		t.Errorf("SetInFreeList true failed")
	}
	// Color must survive flag changes
	c = SetColor(ColorWhite, ColorBlack)
	c = SetHasFinal(c, true)
	c = SetInFreeList(c, true)
	if ColorOf(c) != ColorBlack {
		t.Errorf("color lost when setting flags")
	}
}

func TestHeaderSize(t *testing.T) {
	// Header must be pointer-aligned (4 bytes on the 32-bit Playdate target).
	if HeaderSize()%4 != 0 {
		t.Errorf("HeaderSize() = %d, must be 4-aligned", HeaderSize())
	}
	// 20 bytes expected (size + color + sizeClass + pad + prev + next + nextFree).
	// fakeHeader uses uint32 fields so the size is the same on 32-bit and 64-bit
	// hosts (the real gcHeader on target uses uintptr/pointers which are 4 bytes).
	if HeaderSize() != 20 {
		t.Errorf("HeaderSize() = %d, want 20", HeaderSize())
	}
	// Suppress unused-import warning if no other use of unsafe
	_ = unsafe.Sizeof(uintptr(0))
}
