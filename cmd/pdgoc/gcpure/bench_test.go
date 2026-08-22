package gcpure

// Benchmarks for the GC's core data structures.
//
// These run on the host with `go test -bench . ./gcpure` (from cmd/pdgoc)
// and measure the hot paths behind every GC pause:
//
//   - SizeClassOf          per-allocation size-class lookup
//   - FreeLists            alloc pop / sweep push
//   - AllocList            alloc insert / sweep unlink
//   - ObjectMap            bitmap mark/clear on alloc/free, ObjectStart on
//                          every candidate pointer during conservative marking
//   - MarkStack            gray-queue push/pop during mark
//   - FinalizerTable       SetFinalizer registration / sweep-time lookup
//
// The host versions of FreeLists and FinalizerTable use side maps where the
// runtime uses intrusive links and unsafe.Pointer map keys, so absolute
// numbers overstate device cost — treat these as regression tracking and
// complexity verification (O(1) claims), not device pause predictions.
// End-to-end pause times are measured on hardware by game_examples/gc_pause_benchmark.

import "testing"

// --- size classes ---------------------------------------------------------

func BenchmarkSizeClassOf(b *testing.B) {
	sizes := []uintptr{8, 24, 48, 100, 200, 400, 800, 1500, 4096}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		SizeClassOf(sizes[i%len(sizes)])
		i++
	}
}

// --- free lists (alloc pop / sweep push) ----------------------------------

func BenchmarkFreeListPushPop(b *testing.B) {
	var f FreeLists
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		f.Push(2, uintptr(0x1000+i*16))
		f.Pop(2)
		i++
	}
}

// Pop from a deep list: head access must not walk the list. The list is
// reseeded outside the timed region whenever it drains.
func BenchmarkFreeListPopDeep(b *testing.B) {
	var f FreeLists
	const depth = 1 << 16
	seed := func() {
		for j := range depth {
			f.Push(3, uintptr(0x1000+j*16))
		}
	}
	seed()
	b.ReportAllocs()
	for b.Loop() {
		if f.Len(3) == 0 {
			b.StopTimer()
			seed()
			b.StartTimer()
		}
		if f.Pop(3) == 0 {
			b.Fatal("unexpected empty list")
		}
	}
}

// --- alloc list (alloc insert / sweep unlink) ------------------------------

func BenchmarkAllocListInsertUnlink(b *testing.B) {
	var l AllocList
	n := &AllocNode{}
	b.ReportAllocs()
	for b.Loop() {
		l.Insert(n)
		l.Unlink(n)
	}
}

// Unlink from the middle of a populated list: must stay O(1).
func BenchmarkAllocListUnlinkMiddle(b *testing.B) {
	nodes := make([]AllocNode, 10000)
	var l AllocList
	for i := range nodes {
		l.Insert(&nodes[i])
	}
	mid := &nodes[len(nodes)/2]
	b.ReportAllocs()
	for b.Loop() {
		l.Unlink(mid)
		l.Insert(mid)
	}
}

// --- object bitmap (alloc/free coverage + mark lookup) ---------------------

func benchObjectMapMarkClear(b *testing.B, size uintptr) {
	var m ObjectMap
	m.Init(0x1000, 0x1000+1<<20) // 1 MB heap window
	b.SetBytes(int64(size))
	b.ReportAllocs()
	for b.Loop() {
		m.Mark(0x2000, size)
		m.Clear(0x2000, size)
	}
}

func BenchmarkObjectMapMarkClear16(b *testing.B)  { benchObjectMapMarkClear(b, 16) }
func BenchmarkObjectMapMarkClear64(b *testing.B)  { benchObjectMapMarkClear(b, 64) }
func BenchmarkObjectMapMarkClear256(b *testing.B) { benchObjectMapMarkClear(b, 256) }
func BenchmarkObjectMapMarkClear1K(b *testing.B)  { benchObjectMapMarkClear(b, 1024) }

// ObjectStart at head/middle/tail of a 256-byte object and into free space:
// the conservative marking inner loop calls this for every candidate pointer.
func BenchmarkObjectMapObjectStart(b *testing.B) {
	var m ObjectMap
	m.Init(0x1000, 0x1000+1<<20)
	m.Mark(0x2000, 256)
	addrs := []uintptr{0x2000, 0x2080, 0x20FC, 0x5000}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		m.ObjectStart(addrs[i%len(addrs)])
		i++
	}
}

// Large-object worst case: capped offset bytes (255) make ObjectStart walk
// back in 254-slot steps.
func BenchmarkObjectMapObjectStartLarge(b *testing.B) {
	var m ObjectMap
	m.Init(0x1000, 0x1000+1<<20)
	m.Mark(0x2000, 64<<10) // 64 KB object
	b.ReportAllocs()
	for b.Loop() {
		m.ObjectStart(0x2000 + 64<<10 - 4) // last slot
	}
}

// Conservative scan: resolve every 4-byte slot in a heap window as if each
// were a candidate pointer — the shape of gcMarkReachable's inner loop.
// Half the window is covered by 64-byte objects, half is free.
func BenchmarkConservativeScan(b *testing.B) {
	var m ObjectMap
	const window = 1 << 20 // 1 MB
	m.Init(0x1000, 0x1000+window)
	for a := uintptr(0x1000); a < 0x1000+window/2; a += 64 {
		m.Mark(a, 48) // 64-byte buckets, 48 bytes user data
	}
	b.SetBytes(window)
	b.ReportAllocs()
	var sink uintptr
	for b.Loop() {
		for a := uintptr(0x1000); a < 0x1000+window; a += 4 {
			sink += m.ObjectStart(a)
		}
	}
	if sink == 1 { // keep the loop from being optimized away
		b.Fatal("impossible")
	}
}

// --- mark stack (gray queue) -----------------------------------------------

func BenchmarkMarkStackPushPop(b *testing.B) {
	var s MarkStack
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		s.Push(uintptr(i))
		s.Pop()
		i++
	}
}

// --- finalizer table --------------------------------------------------------

func BenchmarkFinalizerAddGet(b *testing.B) {
	var t FinalizerTable
	fn := func(uintptr) {}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		p := uintptr(0x1000 + i%1000*16)
		t.Add(p, fn)
		t.Get(p) // sweep-time: returns and removes
		i++
	}
}

// --- composite: per-object sweep cost ---------------------------------------

// One dead object through the sweep path: free-list push + objectmap clear +
// alloc-list unlink. The three operations above measured together.
func BenchmarkSweepObject(b *testing.B) {
	var f FreeLists
	var m ObjectMap
	var l AllocList
	nodes := make([]AllocNode, 1)
	m.Init(0x1000, 0x1000+1<<20)
	const addr = uintptr(0x2000)
	b.ReportAllocs()
	for b.Loop() {
		m.Mark(addr, 48)
		l.Insert(&nodes[0])
		// sweep:
		f.Push(1, addr)
		m.Clear(addr, 48)
		l.Unlink(&nodes[0])
	}
}
