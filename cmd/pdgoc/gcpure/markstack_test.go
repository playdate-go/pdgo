package gcpure

import "testing"

func TestMarkStackPushPop(t *testing.T) {
	var s MarkStack
	if s.Len() != 0 {
		t.Fatalf("empty Len() = %d", s.Len())
	}
	s.Push(0x1000)
	s.Push(0x2000)
	if s.Len() != 2 {
		t.Fatalf("Len after 2 pushes = %d", s.Len())
	}
	if v := s.Pop(); v != 0x2000 {
		t.Errorf("Pop = %#x, want 0x2000", v)
	}
	if v := s.Pop(); v != 0x1000 {
		t.Errorf("Pop = %#x, want 0x1000", v)
	}
	if s.Len() != 0 {
		t.Errorf("Len after drain = %d", s.Len())
	}
}

func TestMarkStackGrowsBeyondInitial(t *testing.T) {
	var s MarkStack
	// Push beyond any plausible initial capacity.
	for i := uintptr(1); i <= 5000; i++ {
		s.Push(i)
	}
	if s.Len() != 5000 {
		t.Fatalf("Len = %d after 5000 pushes", s.Len())
	}
	for i := uintptr(5000); i >= 1; i-- {
		if v := s.Pop(); v != i {
			t.Fatalf("Pop at i=%d = %d", i, v)
		}
	}
}

func TestMarkStackPopEmptyPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Pop on empty did not panic")
		}
	}()
	var s MarkStack
	s.Pop()
}
