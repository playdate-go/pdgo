package gcpure

// ObjectMap is a flat byte bitmap for O(1) "which object owns this address?"
// lookups during conservative marking.
//
// One byte per 4-byte heap slot (matching 32-bit pointer alignment). The byte
// value encodes the distance from this slot to the object's userStart:
//   - 0: not part of any object
//   - 1: this slot IS the object start
//   - 2..254: start is (value-1) slots back
//   - 255: start is at least 254 slots back (capped — scan back for large objects)
//
// This encoding handles interior pointers: any pointer within an object
// resolves to the object's start in O(1) for objects up to ~1 KB, and in
// O(slots/254) for larger objects.
//
// The bitmap covers a fixed heap window [base, end). Addresses outside the
// window return 0 from ObjectStart and are ignored by Mark/Clear.
type ObjectMap struct {
	bits []byte
	base uintptr
	end  uintptr
}

// Init allocates the bitmap for the heap window [base, end).
func (m *ObjectMap) Init(base, end uintptr) {
	m.base = base
	m.end = end
	if end < base {
		end = base
	}
	slots := (end - base) >> 2 // 4-byte slots
	m.bits = make([]byte, slots)
}

// Mark tags all 4-byte slots covering [start, start+size) as belonging to an
// object whose userStart is start.
func (m *ObjectMap) Mark(start, size uintptr) {
	if size == 0 || start < m.base || start >= m.end {
		return
	}
	startIdx := (start - m.base) >> 2
	numSlots := (size + 3) >> 2 // ceil(size / 4)
	for i := uintptr(0); i < numSlots; i++ {
		idx := startIdx + i
		if idx >= uintptr(len(m.bits)) {
			break
		}
		offset := i + 1
		if offset > 255 {
			offset = 255
		}
		m.bits[idx] = byte(offset)
	}
}

// Clear zeroes all 4-byte slots covering [start, start+size).
func (m *ObjectMap) Clear(start, size uintptr) {
	if size == 0 || start < m.base || start >= m.end {
		return
	}
	startIdx := (start - m.base) >> 2
	numSlots := (size + 3) >> 2
	for i := uintptr(0); i < numSlots; i++ {
		idx := startIdx + i
		if idx >= uintptr(len(m.bits)) {
			break
		}
		m.bits[idx] = 0
	}
}

// ObjectStart returns the userStart of the object containing addr, or 0 if
// addr is not within any object. Handles interior pointers.
func (m *ObjectMap) ObjectStart(addr uintptr) uintptr {
	if addr < m.base || addr >= m.end {
		return 0
	}
	idx := (addr - m.base) >> 2
	for {
		b := m.bits[idx]
		if b == 0 {
			return 0 // not in any object
		}
		if b == 1 {
			return m.base + idx*4 // found the start
		}
		if b < 255 {
			// Exact offset: start is (b-1) slots back.
			startIdx := idx - uintptr(b) + 1
			return m.base + startIdx*4
		}
		// b == 255: capped. Step back 254 slots and re-check.
		if idx < 254 {
			return 0 // shouldn't happen for valid encoding
		}
		idx -= 254
	}
}
