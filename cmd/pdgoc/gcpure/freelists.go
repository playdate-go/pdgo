package gcpure

// FreeLists is a set of per-size-class LIFO free lists. The host-testable
// version uses a side map for linking because it has no real gcHeader; the
// runtime version (gc_playdate.go) uses intrusive linking through
// gcHeader.nextFree.
type FreeLists struct {
	head [NumSizeClasses]uintptr
	len  [NumSizeClasses]int
	side map[uintptr]uintptr // links ptr -> next ptr (host-test crutch)
}

func (f *FreeLists) Push(sc uint8, ptr uintptr) {
	if sc >= NumSizeClasses {
		return
	}
	if f.side == nil {
		f.side = make(map[uintptr]uintptr)
	}
	f.side[ptr] = f.head[sc]
	f.head[sc] = ptr
	f.len[sc]++
}

func (f *FreeLists) Pop(sc uint8) uintptr {
	if sc >= NumSizeClasses || f.head[sc] == 0 {
		return 0
	}
	p := f.head[sc]
	f.head[sc] = f.side[p]
	delete(f.side, p)
	f.len[sc]--
	return p
}

func (f *FreeLists) Len(sc uint8) int {
	if sc >= NumSizeClasses {
		return 0
	}
	return f.len[sc]
}
