package gcpure

import "testing"

func TestObjectMapMarkClearStart(t *testing.T) {
	var m ObjectMap
	m.Init(0x10000000, 0x10000000+0x10000) // 64KB window

	// Nothing is in an object initially.
	if start := m.ObjectStart(0x10000010); start != 0 {
		t.Errorf("expected 0, got %#x", start)
	}

	// Mark an object: start at 0x10000010, size 32 bytes (8 four-byte slots).
	m.Mark(0x10000010, 32)

	// The start should resolve.
	if start := m.ObjectStart(0x10000010); start != 0x10000010 {
		t.Errorf("start: expected %#x, got %#x", 0x10000010, start)
	}

	// An interior pointer should also resolve to the start.
	if start := m.ObjectStart(0x10000020); start != 0x10000010 {
		t.Errorf("interior: expected %#x, got %#x", 0x10000010, start)
	}

	// Just past the end should not resolve.
	if start := m.ObjectStart(0x10000030); start != 0 {
		t.Errorf("past end: expected 0, got %#x", start)
	}

	// Clear.
	m.Clear(0x10000010, 32)
	if start := m.ObjectStart(0x10000010); start != 0 {
		t.Errorf("after clear: expected 0, got %#x", start)
	}
}

func TestObjectMapOutOfRange(t *testing.T) {
	var m ObjectMap
	m.Init(0x1000, 0x2000)
	// Below base.
	if start := m.ObjectStart(0x100); start != 0 {
		t.Errorf("below base: expected 0, got %#x", start)
	}
	// Above end.
	if start := m.ObjectStart(0x20000); start != 0 {
		t.Errorf("above end: expected 0, got %#x", start)
	}
}

func TestObjectMapClearIdempotent(t *testing.T) {
	var m ObjectMap
	m.Init(0x1000, 0x2000)
	m.Clear(0x1000, 16) // never set — must not panic
	m.Clear(0x1000, 16) // again
}

func TestObjectMapLargeObject(t *testing.T) {
	var m ObjectMap
	m.Init(0x1000, 0x10000) // 60KB window

	// Object larger than 254*4 = 1016 bytes — tests the 255 cap.
	start := uintptr(0x1000)
	size := uintptr(4096) // 1024 four-byte slots
	m.Mark(start, size)

	// The start itself.
	if s := m.ObjectStart(start); s != start {
		t.Errorf("start: expected %#x, got %#x", start, s)
	}

	// Interior pointer at offset 512 (within O(1) range, offset < 255).
	mid1 := start + 512
	if s := m.ObjectStart(mid1); s != start {
		t.Errorf("mid1: expected %#x, got %#x", start, s)
	}

	// Interior pointer at offset 2048 (in the 255-capped region).
	mid2 := start + 2048
	if s := m.ObjectStart(mid2); s != start {
		t.Errorf("mid2 (capped): expected %#x, got %#x", start, s)
	}

	// Very end of object (last 4-byte slot).
	end := start + size - 4
	if s := m.ObjectStart(end); s != start {
		t.Errorf("end: expected %#x, got %#x", start, s)
	}

	// Just past the end.
	if s := m.ObjectStart(start + size); s != 0 {
		t.Errorf("past end: expected 0, got %#x", s)
	}
}

func TestObjectMapAdjacentObjects(t *testing.T) {
	var m ObjectMap
	m.Init(0x1000, 0x2000)

	// Two adjacent objects.
	m.Mark(0x1000, 16)
	m.Mark(0x1010, 16)

	// Interior of first.
	if s := m.ObjectStart(0x1008); s != 0x1000 {
		t.Errorf("obj1 interior: expected %#x, got %#x", 0x1000, s)
	}
	// Interior of second.
	if s := m.ObjectStart(0x1018); s != 0x1010 {
		t.Errorf("obj2 interior: expected %#x, got %#x", 0x1010, s)
	}
	// Boundary — last slot of first.
	if s := m.ObjectStart(0x100C); s != 0x1000 {
		t.Errorf("obj1 last slot: expected %#x, got %#x", 0x1000, s)
	}
}
