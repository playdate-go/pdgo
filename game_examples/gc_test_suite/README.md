# GC Test Suite

Tests the pdgo garbage collector on the Playdate device using native Go memory constructs.
Output is via `LogToConsole` (visible in simulator console or device serial).

## Why This Exists

The device build runs a custom conservative mark-sweep GC (TinyGo patch), not the stock Go GC —
and the simulator cannot validate it, because simulator builds use the standard host Go runtime.
Every change to the collector therefore needs an on-hardware test suite.

## What It Proves

| Group | Tests | Proves |
|-------|-------|--------|
| Correctness across constructs | Slices, Maps, Structs, PointerChains, Trees, Interfaces, Channels, Closures, Strings, NestedSlices, AppendGrowth | every Go allocation construct works under the custom GC |
| Live data is never corrupted | RetainedMemory, StressTest, LargeLiveSet | reachable objects survive GC cycles byte-for-byte (500 KB / 4000 objects) — the class of bug a wrong conservative marker produces |
| C resources are cleaned up | FinalizerChurn | all 500 registered finalizers run, so SDK memory (bitmaps, sounds, files) is not leaked |
| Pauses fit a frame | PauseBudget | worst GC pause < 3 ms, inside the 50 ms frame budget |
| Memory is actually reclaimed | FreeListReuse | SDK allocation count stays flat under churn — the free list recycles instead of growing the heap without bound |

## What It Tests

All tests use **pure Go constructs** (no C API calls in test code):

| Test | Description |
|------|-------------|
| Slices | `make([]T, n)` allocations |
| Maps | `make(map[K]V)` allocations |
| Structs(new) | `new(T)` allocations |
| Structs(literal) | `&Struct{}` allocations |
| PointerChains | Linked list structures |
| Trees | Binary tree structures |
| Interfaces | Interface boxing |
| Channels | Buffered/unbuffered channels |
| Closures | Closure captures |
| Strings | String concatenation/conversions |
| NestedSlices | 2D/3D slices |
| AppendGrowth | Slice growth via append |
| RetainedMemory | Verifies GC doesn't corrupt live data |
| StressTest | Rapid allocation/deallocation |
| LargeLiveSet | 500KB across 4000 objects survives GC uncorrupted |
| FinalizerChurn | 500 finalizers all run after alloc pressure |
| PauseBudget | Worst GC pause < 3ms |
| FreeListReuse | SDK alloc count stays flat under churn (free-list reuse) |

## Building

```bash
cd examples/gc_test_suite
./build.sh
```

Or manually:
```bash
pdgoc -sim -device \
  -name="GCTestSuite" \
  -author="PdGo" \
  -desc="GC Test Suite" \
  -bundle-id=com.pdgo.gctestsuite \
  -version=1.0 \
  -build-number=1
```

## Running

### Simulator
```bash
open GCTestSuite_sim.pdx
```

View output in the Playdate Simulator console window.

### Device
1. Copy `GCTestSuite.pdx` to your Playdate
2. Run the game
3. Connect via serial to see log output

## Expected Output

```
========================================
GC Test Suite
========================================

Initial HeapAlloc: 45000
Initial NumGC: 0

----------------------------------------
Running test 1

Test: Slices - allocating slices of various types
Test: Slices - PASS
HeapAlloc: 52000
NumGC: 2

... (more tests)

========================================
FINAL RESULTS
========================================
Tests Passed: 6
Tests Failed: 0
RetainedMemory: OK

--- Final Memory Stats ---
HeapAlloc: 48000
HeapSys: 1048576
NumGC: 25

========================================
ALL TESTS PASSED!
GC is working correctly!
========================================
```

## Success Criteria

- All 6 active tests pass
- `RetainedMemory: OK` (1KB of data survives GC intact)
- `NumGC` > 0 (GC has run)
- Memory doesn't grow unboundedly
