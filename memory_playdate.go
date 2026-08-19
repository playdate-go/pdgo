//go:build gc.playdate

package pdgo

import (
	"runtime"
)

// readPlaydateMemStats returns GC statistics from the runtime.
func readPlaydateMemStats() (heapAlloc uint64, numGC uint32, liveObjects uint32, lastPauseNs int64, totalAlloc uint64, totalFrees uint64) {
	return runtime.ReadPlaydateMemStats()
}

// resetGCStats zeroes the GC counters in the runtime.
func resetGCStats() {
	runtime.ResetGCStats()
}
