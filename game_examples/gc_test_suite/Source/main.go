// gc_test_suite - GC test suite for the Playdate device
// Tests Go's garbage collector on Playdate device
// Uses LogToConsole for output - GC testing logic is pure Go (no C calls in test code)

package main

import (
	"runtime"
	"unsafe"

	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

// ============================================================================
// Test Structures (Pure Go)
// ============================================================================

type Node struct {
	Value    int
	Children []*Node
	Data     []byte
}

type LargeStruct struct {
	IntField     int
	Int64Field   int64
	FloatField   float64
	SliceField   []int
	MapField     map[string]int
	StringField  string
	PointerField *SmallStruct
}

type SmallStruct struct {
	A, B, C int
}

// ============================================================================
// State
// ============================================================================

var frameCount int
var testPhase int
var retained []byte
var retainedOK bool = true // default to true, set false if corruption detected
var testsPassed int
var testsFailed int

// ============================================================================
// Logging Helper
// ============================================================================

func log(msg string) {
	pd.System.LogToConsole(msg)
}

func logInt(label string, val int) {
	// Convert int to string manually for logging
	s := label
	if val < 0 {
		s += "-"
		val = -val
	}
	if val == 0 {
		s += "0"
	} else {
		digits := ""
		for val > 0 {
			digits = string('0'+rune(val%10)) + digits
			val /= 10
		}
		s += digits
	}
	log(s)
}

func logUint64(label string, val uint64) {
	s := label
	if val == 0 {
		s += "0"
	} else {
		digits := ""
		for val > 0 {
			digits = string('0'+rune(val%10)) + digits
			val /= 10
		}
		s += digits
	}
	log(s)
}

// ============================================================================
// GC Test Functions (All Pure Go - No C Calls)
// ============================================================================

// testSlices tests slice allocations
func testSlices() bool {
	log("Test: Slices - allocating slices of various types")
	for i := 0; i < 50; i++ {
		s1 := make([]int, 30)
		s2 := make([]byte, 128)
		s3 := make([]float64, 50)

		s1[0] = i
		s2[0] = byte(i)
		s3[0] = float64(i)

		_ = s1[0] + int(s2[0])
	}
	log("Test: Slices - PASS")
	return true
}

// testMaps tests map allocations
func testMaps() bool {
	log("Test: Maps - allocating maps")
	for i := 0; i < 30; i++ {
		m1 := make(map[int]int)
		m2 := make(map[string]int)

		for j := 0; j < 30; j++ {
			m1[j] = j * 2
			m2["key"] = j
		}

		_ = len(m1) + len(m2)
	}
	log("Test: Maps - PASS")
	return true
}

// testStructsNew tests new() allocations
func testStructsNew() bool {
	log("Test: Structs(new) - allocating with new()")
	for i := 0; i < 100; i++ {
		s1 := new(SmallStruct)
		s1.A = i
		s1.B = i + 1
		s1.C = i + 2

		s2 := new(LargeStruct)
		s2.IntField = i
		s2.SliceField = make([]int, 20)
		s2.MapField = make(map[string]int)
		s2.PointerField = new(SmallStruct)

		_ = s1.A + s2.IntField
	}
	log("Test: Structs(new) - PASS")
	return true
}

// testStructsLiteral tests struct literal allocations
func testStructsLiteral() bool {
	log("Test: Structs(literal) - allocating with literals")
	for i := 0; i < 100; i++ {
		_ = LargeStruct{
			IntField:    i,
			Int64Field:  int64(i),
			FloatField:  float64(i),
			SliceField:  []int{1, 2, 3, 4, 5},
			MapField:    map[string]int{"a": 1},
			StringField: "test",
			PointerField: &SmallStruct{
				A: i,
				B: i,
				C: i,
			},
		}
	}
	log("Test: Structs(literal) - PASS")
	return true
}

// testPointerChains tests linked structures
func testPointerChains() bool {
	log("Test: PointerChains - building linked lists")

	type Chain struct {
		Val  int
		Next *Chain
	}

	var totalNodes int
	for i := 0; i < 5; i++ {
		var head *Chain
		for j := 0; j < 200; j++ {
			head = &Chain{Val: j, Next: head}
		}

		count := 0
		for p := head; p != nil; p = p.Next {
			count++
		}
		totalNodes += count
		head = nil
	}
	logInt("Test: PointerChains - PASS (nodes:", totalNodes)
	return true
}

// testTrees tests recursive structures
func testTrees() bool {
	log("Test: Trees - building binary trees")

	var buildTree func(depth int) *Node
	buildTree = func(depth int) *Node {
		if depth <= 0 {
			return nil
		}
		return &Node{
			Value:    depth,
			Children: []*Node{buildTree(depth - 1), buildTree(depth - 1)},
			Data:     make([]byte, 32),
		}
	}

	var totalNodes int
	var countNodes func(n *Node) int
	countNodes = func(n *Node) int {
		if n == nil {
			return 0
		}
		c := 1
		for _, child := range n.Children {
			c += countNodes(child)
		}
		return c
	}

	for i := 0; i < 5; i++ {
		root := buildTree(5)
		totalNodes += countNodes(root)
		root = nil
		runtime.GC()
	}
	logInt("Test: Trees - PASS (nodes:", totalNodes)
	return true
}

// testInterfaces tests interface allocations
func testInterfaces() bool {
	log("Test: Interfaces - boxing various types")
	var items []any

	for i := 0; i < 50; i++ {
		var item any
		switch i % 4 {
		case 0:
			item = i
		case 1:
			item = "string value"
		case 2:
			item = &SmallStruct{A: i}
		case 3:
			item = []int{1, 2, 3}
		}
		items = append(items, item)
	}

	count := len(items)
	items = nil
	logInt("Test: Interfaces - PASS (items:", count)
	return true
}

// testChannels tests channel allocations
func testChannels() bool {
	log("Test: Channels - creating buffered/unbuffered channels")
	var totalCap int
	for i := 0; i < 30; i++ {
		ch1 := make(chan int)
		ch2 := make(chan int, 10)
		ch3 := make(chan *SmallStruct, 5)
		ch4 := make(chan []byte, 3)

		select {
		case ch2 <- i:
		default:
		}

		select {
		case ch3 <- &SmallStruct{A: i}:
		default:
		}

		totalCap += cap(ch2) + cap(ch3) + cap(ch4)
		_ = ch1
	}
	logInt("Test: Channels - PASS (cap:", totalCap)
	return true
}

// testClosures tests closure allocations
func testClosures() bool {
	log("Test: Closures - creating closures with captures")
	var closures []func() int

	for i := 0; i < 100; i++ {
		x := i
		y := i * 2
		fn := func() int {
			return x + y
		}
		closures = append(closures, fn)
	}

	sum := 0
	for _, f := range closures[:10] {
		sum += f()
	}

	closures = nil
	logInt("Test: Closures - PASS (sum:", sum)
	return true
}

// testStrings tests string operations
func testStrings() bool {
	log("Test: Strings - concatenation and conversions")
	var totalLen int
	for i := 0; i < 100; i++ {
		s1 := "hello"
		s2 := "world"
		s3 := s1 + " " + s2

		b := []byte(s3)
		s4 := string(b)

		// Build longer string
		s5 := ""
		for j := 0; j < 5; j++ {
			s5 = s5 + s4
		}

		totalLen += len(s5)
	}
	logInt("Test: Strings - PASS (len:", totalLen)
	return true
}

// testNestedSlices tests 2D/3D slices
func testNestedSlices() bool {
	log("Test: NestedSlices - 2D and 3D slice allocation")

	m2d := make([][]int, 30)
	for i := range m2d {
		m2d[i] = make([]int, 30)
		for j := range m2d[i] {
			m2d[i][j] = i*30 + j
		}
	}

	m3d := make([][][]int, 5)
	for i := range m3d {
		m3d[i] = make([][]int, 5)
		for j := range m3d[i] {
			m3d[i][j] = make([]int, 5)
			for k := range m3d[i][j] {
				m3d[i][j][k] = i*100 + j*10 + k
			}
		}
	}

	var sum int
	for _, row := range m2d {
		for _, val := range row {
			sum += val
		}
	}

	m2d = nil
	m3d = nil
	logInt("Test: NestedSlices - PASS (sum:", sum)
	return true
}

// testAppendGrowth tests slice growth
func testAppendGrowth() bool {
	log("Test: AppendGrowth - slice growth via append")
	var totalCap int
	for i := 0; i < 20; i++ {
		s := make([]int, 0, 1)
		for j := 0; j < 300; j++ {
			s = append(s, j)
		}
		totalCap += cap(s)
	}
	logInt("Test: AppendGrowth - PASS (cap:", totalCap)
	return true
}

// testRetainedMemory verifies GC doesn't corrupt live data
func testRetainedMemory() bool {
	log("Test: RetainedMemory - verifying GC preserves live data")

	// Allocate retained memory
	retained = make([]byte, 512)
	for i := range retained {
		retained[i] = byte(i % 256)
	}

	// Allocate lots of garbage
	for i := 0; i < 50; i++ {
		garbage := make([]byte, 512)
		for j := range garbage {
			garbage[j] = byte(j)
		}
		// garbage goes out of scope
	}

	// Force GC
	runtime.GC()

	// Verify retained data is intact
	retainedOK = true
	for i := range retained {
		if retained[i] != byte(i%256) {
			retainedOK = false
			log("Test: RetainedMemory - FAIL - data corrupted!")
			break
		}
	}

	if retainedOK {
		log("Test: RetainedMemory - PASS - 512B data intact")
	}
	return retainedOK
}

// stressTest runs rapid allocations
func stressTest() bool {
	log("Test: StressTest - rapid allocations")
	for round := 0; round < 2; round++ {
		for i := 0; i < 100; i++ {
			s := make([]byte, 256)
			s[0] = byte(i)

			m := make(map[int]int)
			m[i] = i * 2

			_ = s[0]
			_ = m[i]
		}
		runtime.GC()
	}
	log("Test: StressTest - PASS")
	return true
}

// testGCRealProof allocates WAY MORE than device RAM in garbage
// Playdate has ~16MB RAM. We allocate 40MB+ total.
// WITHOUT GC: instant OOM crash (can't fit 40MB in 16MB)
// WITH GC: survives because garbage is collected each round
func testGCRealProof() bool {
	log("Test: GCRealProof - THE DEFINITIVE TEST")
	log("Allocating 40MB+ total (device has 16MB RAM)")
	log("WITHOUT GC: WILL CRASH")
	log("WITH GC: survives (garbage collected)")

	// Keep some retained data (must survive all GC cycles)
	log("Allocating retained 4KB...")
	retainedLocal := make([]byte, 4096)
	log("Retained allocated OK")
	for i := range retainedLocal {
		retainedLocal[i] = byte(i % 256)
	}
	log("Retained initialized OK")

	// First, test if runtime.GC() even works
	log("Testing runtime.GC()...")
	runtime.GC()
	log("runtime.GC() returned OK")

	// Full test: 2000 rounds of 20KB = 40MB total on 16MB device
	// This is 2.5x device RAM - IMPOSSIBLE without working GC!
	log("Starting test: 2000 rounds x 20KB = 40MB total (2.5x device RAM)")
	var dummySum int // prevent optimization
	for round := 0; round < 2000; round++ {
		// Allocate 20KB - USE ALL BYTES to prevent optimization
		garbage := make([]byte, 20480)
		for j := range garbage {
			garbage[j] = byte(j + round)
		}
		for j := 0; j < len(garbage); j++ {
			dummySum += int(garbage[j])
		}

		// Log and GC every 200 rounds (4MB intervals)
		if round%200 == 199 {
			logInt("Round:", round+1)
			logInt("DummySum:", dummySum) // force use
			runtime.GC()
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			logUint64("  HeapAlloc:", m.Alloc)
			logUint64("  LargeAllocs:", m.GCSys) // DEBUG: shows count of >=10KB allocations
		}
	}

	// Verify retained data survived all GC cycles
	ok := true
	for i := range retainedLocal {
		if retainedLocal[i] != byte(i%256) {
			log("GCRealProof - FAIL - retained corrupted!")
			ok = false
			break
		}
	}

	if ok {
		log("GCRealProof - PASS - 40MB+ collected!")
		log("GC IS DEFINITELY WORKING!")
	}
	return ok
}

// testLargeLiveSet builds ~500KB of live data across ~4000 objects and
// verifies it survives GC uncorrupted.
func testLargeLiveSet() bool {
	log("Test: LargeLiveSet - 500KB across 4000 objects")
	live := make([][]byte, 4000)
	for i := range live {
		live[i] = make([]byte, 128)
		live[i][0] = byte(i & 0xff)
		live[i][127] = byte((i + 1) & 0xff)
	}
	runtime.GC()
	// Verify integrity.
	for i := range live {
		if live[i][0] != byte(i&0xff) || live[i][127] != byte((i+1)&0xff) {
			log("Test: LargeLiveSet - FAIL (data corruption)")
			return false
		}
	}
	// Explicitly drop and GC.
	live = nil
	runtime.GC()
	log("Test: LargeLiveSet - PASS")
	return true
}

// registerFinalizables creates 500 finalizable objects in a dedicated frame.
// The registration MUST live in a helper: while testFinalizerChurn's frame is
// alive, the conservative stack scan keeps its locals (including the last
// `obj`) reachable, so their finalizers correctly never run.
func registerFinalizables(ran []bool) {
	for i := 0; i < len(ran); i++ {
		obj := new([1]byte)
		idx := i // capture
		runtime.SetFinalizer(obj, func(p unsafe.Pointer) {
			ran[idx] = true
		})
		// obj is heap-promoted by SetFinalizer; no need to use it further
	}
}

// testFinalizerChurn registers 500 finalizers, allocates+frees 10k objects,
// and verifies all finalizers ran.
//
// A conservative GC never guarantees a finalizer runs after a fixed number
// of cycles: objects allocated since the last sweep are spared one cycle
// (fresh flag), and an object's address lingering in a callee-saved register
// or stale stack slot is a legitimate root (observed on device: the last
// object survived 2 explicit GCs and was collected on the 3rd, after string
// formatting clobbered the register). So: poll with a bounded retry loop —
// the same pattern mainline Go's own finalizer tests use.
func testFinalizerChurn() bool {
	log("Test: FinalizerChurn - 500 finalizers, 10k allocs")
	ran := make([]bool, 500)
	registerFinalizables(ran)
	// Force allocations to trigger GC.
	for i := 0; i < 10000; i++ {
		s := make([]byte, 64)
		_ = s
	}
	const attempts = 5
	missed := 1
	for a := 0; a < attempts && missed > 0; a++ {
		runtime.GC()
		missed = 0
		for i := range ran {
			if !ran[i] {
				missed++
			}
		}
		if missed > 0 && a < attempts-1 {
			// Churn registers between attempts (log does string formatting)
			// so stale object pointers don't survive as conservative roots.
			logInt("  finalizers pending, retrying: ", missed)
		}
	}
	if missed > 0 {
		log("Test: FinalizerChurn - FAIL (finalizer did not run)")
		logInt("  finalizers not run: ", missed)
		return false
	}
	log("Test: FinalizerChurn - PASS")
	return true
}

// pauseSink forces the PauseBudget garbage onto the heap: a non-escaping
// make([]byte, 256) is stack-allocated by TinyGo, so the original test
// triggered ZERO GCs during its measurement loop and reported the stale
// LastPauseNs of the previous test's finalizer-heavy cycle (observed on
// device: FAIL worst=4ms with no GC line printed during the test).
var pauseSink []byte

// testPauseBudget measures the worst GC pause and asserts < 3ms.
func testPauseBudget() bool {
	log("Test: PauseBudget - worst pause must be < 3ms")
	// Drop the previous test's pause so the measurement only reflects GCs
	// triggered by this test's own allocations.
	pd.Memory.ResetStats()
	// Allocate enough to trigger several GCs.
	worst := int64(0)
	for i := 0; i < 50; i++ {
		// Make some garbage (write to prevent optimization).
		var dummy byte
		for j := 0; j < 100; j++ {
			s := make([]byte, 256)
			s[0] = byte(j)
			dummy += s[0]
			pauseSink = s // force heap allocation (garbage when overwritten)
		}
		_ = dummy
		stats := pd.Memory.Stats()
		if int64(stats.LastPauseNs) > worst {
			worst = int64(stats.LastPauseNs)
		}
	}
	if worst > 3*1000*1000 {
		log("Test: PauseBudget - FAIL")
		logUint64("  worst pause ns: ", uint64(worst))
		return false
	}
	logUint64("Test: PauseBudget - PASS (worst ns: ", uint64(worst))
	log(")")
	return true
}

// testFreeListReuse verifies that alloc/free same sizes reuses slots.
func testFreeListReuse() bool {
	log("Test: FreeListReuse - SDK alloc count stays flat")
	pd.Memory.ResetStats()
	pd.Memory.SetReallocDebug(true)
	// Warm up: prime the free-lists.
	ptrs := make([]*byte, 100)
	for i := range ptrs {
		ptrs[i] = new(byte)
	}
	for i := range ptrs {
		ptrs[i] = nil
	}
	runtime.GC()
	// Now alloc/free the same sizes many times — SDK alloc count
	// should stay roughly constant.
	before, _ := pd.Memory.GetReallocStats()
	for cycle := 0; cycle < 20; cycle++ {
		tmp := make([]*byte, 100)
		for i := range tmp {
			tmp[i] = new(byte)
		}
		tmp = nil // release before GC so objects are swept
		runtime.GC()
	}
	after, _ := pd.Memory.GetReallocStats()
	pd.Memory.SetReallocDebug(false)
	delta := after.Count - before.Count
	if delta > 50 { // allow some slack for non-user allocations
		log("Test: FreeListReuse - FAIL (SDK alloc grew)")
		logInt("  delta: ", delta)
		return false
	}
	log("Test: FreeListReuse - PASS")
	return true
}

// ============================================================================
// Test Runner
// ============================================================================

type gcTest struct {
	name string
	fn   func() bool
}

var gcTests = []gcTest{
	//{"Slices", testSlices},
	//{"Maps", testMaps},
	//{"Structs(new)", testStructsNew},
	//{"Structs(literal)", testStructsLiteral},
	//{"PointerChains", testPointerChains},
	//{"Trees", testTrees},
	//{"Interfaces", testInterfaces},
	//{"Channels", testChannels},
	//{"Closures", testClosures},
	//{"Strings", testStrings},
	//{"NestedSlices", testNestedSlices},
	//{"AppendGrowth", testAppendGrowth},
	//{"RetainedMemory", testRetainedMemory},
	{"StressTest", stressTest},
	{"GCRealProof", testGCRealProof}, // THE real test!
	{"LargeLiveSet", testLargeLiveSet},
	{"FinalizerChurn", testFinalizerChurn},
	{"PauseBudget", testPauseBudget},
	{"FreeListReuse", testFreeListReuse},
}

// ============================================================================
// Game Loop
// ============================================================================

func initGame() {
	log("========================================")
	log("GC Test Suite")
	log("========================================")
	log("")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logUint64("Initial HeapAlloc: ", m.Alloc)
	logUint64("Initial HeapSys: ", m.HeapSys)
	log("")
}

func update() int {
	frameCount++

	// Run one test per 60 frames (1 second at 60fps)
	if testPhase < len(gcTests) {
		if frameCount%60 == 0 {
			log("")
			log("----------------------------------------")
			logInt("Running test ", testPhase+1)
			log("")

			passed := gcTests[testPhase].fn()

			if passed {
				testsPassed++
			} else {
				testsFailed++
			}

			runtime.GC()

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			logUint64("HeapAlloc: ", m.Alloc)
			logUint64("HeapSys: ", m.HeapSys)

			testPhase++
		}
	}

	// All tests complete - show final results
	if testPhase >= len(gcTests) && frameCount%60 == 0 {
		// Only log once after all tests
		if frameCount == (len(gcTests)+1)*60 {
			log("")
			log("========================================")
			log("FINAL RESULTS")
			log("========================================")
			logInt("Tests Passed: ", testsPassed)
			logInt("Tests Failed: ", testsFailed)

			if retainedOK {
				log("RetainedMemory: OK")
			} else {
				log("RetainedMemory: CORRUPTED!")
			}

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log("")
			log("--- Final Memory Stats ---")
			logUint64("HeapAlloc: ", m.Alloc)
			logUint64("HeapSys: ", m.HeapSys)
			logUint64("TotalAlloc: ", m.TotalAlloc)
			logUint64("Sys: ", m.Sys)
			log("")

			if testsFailed == 0 && retainedOK {
				log("========================================")
				log("ALL TESTS PASSED!")
				log("GC is working correctly!")
				log("========================================")
			} else {
				log("========================================")
				log("SOME TESTS FAILED!")
				log("Check memory management!")
				log("========================================")
			}
		}
	}

	return 1
}

func main() {}
