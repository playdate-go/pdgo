package gcpure

import "testing"

func TestFinalizerTableAddRemove(t *testing.T) {
	var ft FinalizerTable
	called := false
	ft.Add(0x1000, func(p uintptr) { called = true })
	if ft.Len() != 1 {
		t.Errorf("Len = %d, want 1", ft.Len())
	}
	fn := ft.Get(0x1000)
	if fn == nil {
		t.Fatal("Get returned nil")
	}
	fn(0x1000)
	if !called {
		t.Error("finalizer not called")
	}
	ft.Remove(0x1000)
	if ft.Len() != 0 {
		t.Errorf("Len after Remove = %d, got %d", ft.Len(), 0)
	}
	if ft.Get(0x1000) != nil {
		t.Error("Get after Remove should be nil")
	}
}

func TestFinalizerTableBeyond128Entries(t *testing.T) {
	var ft FinalizerTable
	for i := uintptr(0); i < 200; i++ {
		ft.Add(i, func(p uintptr) {})
	}
	if ft.Len() != 200 {
		t.Errorf("Len = %d, want 200 (no fixed cap)", ft.Len())
	}
}

func TestFinalizerTableGetUncalledRemoves(t *testing.T) {
	var ft FinalizerTable
	called := false
	ft.Add(0x42, func(p uintptr) { called = true })
	// Get + Remove is the explicit-free path — finalizer NOT called.
	_ = ft.Get(0x42)
	ft.Remove(0x42)
	if called {
		t.Error("Remove triggered finalizer; should not")
	}
}
