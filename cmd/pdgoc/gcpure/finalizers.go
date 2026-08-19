package gcpure

// FinalizerFunc is the function signature for object finalizers.
type FinalizerFunc func(p uintptr)

// FinalizerTable maps object pointers to finalizer functions.
// Backed by a Go map; grows unboundedly.
//
// The runtime version lives in tinygo-patches/gc_finalizer_playdate.go and
// uses unsafe.Pointer keys. The host version uses uintptr for testability.
type FinalizerTable struct {
	m map[uintptr]FinalizerFunc
}

func (t *FinalizerTable) Add(p uintptr, fn FinalizerFunc) {
	if t.m == nil {
		t.m = make(map[uintptr]FinalizerFunc)
	}
	t.m[p] = fn
}

// Get returns the finalizer and REMOVES it from the table.
// This is the sweep-time accessor: caller is expected to invoke the fn.
// For explicit-free, use Remove (which doesn't return the fn).
func (t *FinalizerTable) Get(p uintptr) FinalizerFunc {
	if t.m == nil {
		return nil
	}
	fn, ok := t.m[p]
	if !ok {
		return nil
	}
	delete(t.m, p)
	return fn
}

// Remove drops the finalizer without returning it. Used by explicit free.
func (t *FinalizerTable) Remove(p uintptr) {
	if t.m == nil {
		return
	}
	delete(t.m, p)
}

func (t *FinalizerTable) Len() int {
	return len(t.m)
}
