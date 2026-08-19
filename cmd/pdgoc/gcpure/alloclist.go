package gcpure

// AllocList is a doubly-linked list of allocations. The host-test version
// uses AllocNode; the runtime version uses gcHeader.prev/next directly
// with the same logic.
type AllocNode struct {
	ID         int // for test identification
	prev, next *AllocNode
}

type AllocList struct {
	Head *AllocNode
	Len  int
}

func (l *AllocList) Insert(n *AllocNode) {
	n.prev = nil
	n.next = l.Head
	if l.Head != nil {
		l.Head.prev = n
	}
	l.Head = n
	l.Len++
}

func (l *AllocList) Unlink(n *AllocNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else if l.Head == n {
		l.Head = n.next
	} else {
		// Not in list (Head != n && n.prev == nil). No-op.
		return
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	n.prev = nil
	n.next = nil
	l.Len--
}
