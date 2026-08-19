//go:build !gc.playdate

package pdgo

import (
	"runtime"
)

// readPlaydateMemStats returns GC statistics using standard runtime.ReadMemStats.
// This is a fallback for simulator/host builds where gc_playdate.go is not included.
func readPlaydateMemStats() (heapAlloc uint64, numGC uint32, liveObjects uint32, lastPauseNs int64, totalAlloc uint64, totalFrees uint64) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats.HeapAlloc, stats.NumGC, 0, 0, stats.TotalAlloc, stats.Frees
}

// resetGCStats is a no-op for host builds (runtime doesn't support resetting stats).
func resetGCStats() {
	// Standard runtime doesn't support resetting stats
}
