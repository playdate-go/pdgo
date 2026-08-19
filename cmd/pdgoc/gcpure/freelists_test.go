package gcpure

import "testing"

func TestFreeListPushPopByClass(t *testing.T) {
	var fl FreeLists
	fl.Push(2, 0x1000)
	fl.Push(2, 0x2000)
	fl.Push(0, 0x3000)

	got := fl.Pop(2)
	if got != 0x2000 { // LIFO
		t.Errorf("Pop(2) = %#x, want 0x2000 (LIFO)", got)
	}
	got = fl.Pop(2)
	if got != 0x1000 {
		t.Errorf("Pop(2) = %#x, want 0x1000", got)
	}
	got = fl.Pop(2)
	if got != 0 { // empty
		t.Errorf("Pop(2) on empty = %#x, want 0", got)
	}
	got = fl.Pop(0)
	if got != 0x3000 {
		t.Errorf("Pop(0) = %#x, want 0x3000", got)
	}
}

func TestFreeListPopEmptyReturnsZero(t *testing.T) {
	var fl FreeLists
	for sc := uint8(0); sc < NumSizeClasses; sc++ {
		if v := fl.Pop(sc); v != 0 {
			t.Errorf("Pop(%d) on empty = %#x, want 0", sc, v)
		}
	}
}

func TestFreeListLen(t *testing.T) {
	var fl FreeLists
	if fl.Len(3) != 0 {
		t.Error("empty Len != 0")
	}
	fl.Push(3, 1)
	fl.Push(3, 2)
	fl.Push(3, 3)
	if fl.Len(3) != 3 {
		t.Errorf("Len(3) = %d, want 3", fl.Len(3))
	}
}
