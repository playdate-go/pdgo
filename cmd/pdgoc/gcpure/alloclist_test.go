package gcpure

import "testing"

func TestAllocListInsertUnlink(t *testing.T) {
	var l AllocList
	a := &AllocNode{ID: 1}
	b := &AllocNode{ID: 2}
	c := &AllocNode{ID: 3}

	l.Insert(a)
	l.Insert(b)
	l.Insert(c)
	if l.Len != 3 {
		t.Errorf("Len = %d, want 3", l.Len)
	}

	// Unlink middle — must be O(1) (no walk).
	l.Unlink(b)
	if l.Len != 2 {
		t.Errorf("Len after unlink = %d, want 2", l.Len)
	}

	// Unlink head.
	l.Unlink(c)
	if l.Len != 1 {
		t.Errorf("Len after head unlink = %d, want 1", l.Len)
	}

	// Unlink last.
	l.Unlink(a)
	if l.Len != 0 {
		t.Errorf("Len after last unlink = %d, want 0", l.Len)
	}
}

func TestAllocListUnlinkNotPresentNoOp(t *testing.T) {
	var l AllocList
	a := &AllocNode{ID: 1}
	l.Insert(a)
	orphan := &AllocNode{ID: 99}
	l.Unlink(orphan) // not in list — must not corrupt
	if l.Len != 1 {
		t.Errorf("Len = %d, want 1", l.Len)
	}
}
