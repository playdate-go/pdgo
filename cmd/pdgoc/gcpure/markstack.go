package gcpure

// MarkStack is a growable LIFO stack of uintptr used as the gray-work
// queue during the mark phase. Unlike a fixed-size buffer it never
// silently drops entries on overflow.
type MarkStack struct {
	data []uintptr
}

func (s *MarkStack) Push(p uintptr) {
	s.data = append(s.data, p)
}

func (s *MarkStack) Pop() uintptr {
	// Note: runtime version uses no panic — caller checks Len first.
	// Host version panics to make bugs visible in tests.
	if len(s.data) == 0 {
		panic("MarkStack: Pop on empty")
	}
	n := len(s.data) - 1
	v := s.data[n]
	s.data = s.data[:n]
	return v
}

func (s *MarkStack) Len() int { return len(s.data) }
