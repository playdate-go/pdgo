// pdgo Memory subsystem — GC tuning and stats.

package pdgo

/*
#include <stdint.h>

// Realloc debug functions (from pd_cgo.c)
void pd_sys_setReallocDebug(int enabled);
int pd_sys_getReallocStats(int* count, unsigned long* total_bytes, int* free_count);
void pd_sys_resetReallocStats(void);
*/
import "C"

import (
	"runtime"
)

// MemStats describes GC state.
type MemStats struct {
	HeapAlloc   uint64 // live bytes (total alloc - frees)
	HeapInUse   uint64 // same as HeapAlloc (no idle pool)
	NumGC       uint32
	LiveObjects uint32
	LastPauseNs uint64 // ms granularity on device; 0 on host builds
	TotalAlloc  uint64
	TotalFrees  uint64
}

// Memory exposes the GC subsystem.
type Memory struct{}

// Stats returns current memory statistics.
func (m *Memory) Stats() MemStats {
	heapAlloc, numGC, liveObjects, lastPauseNs, totalAlloc, totalFrees := readPlaydateMemStats()
	return MemStats{
		HeapAlloc:   heapAlloc,
		HeapInUse:   heapAlloc,
		NumGC:       numGC,
		LiveObjects: liveObjects,
		LastPauseNs: uint64(lastPauseNs),
		TotalAlloc:  totalAlloc,
		TotalFrees:  totalFrees,
	}
}

// RunGC forces a GC cycle. Returns the pause duration in nanoseconds.
// On host builds (no gc.playdate tag) the return is always 0.
func (m *Memory) RunGC() int64 {
	runtime.GC()
	_, _, _, lastPauseNs, _, _ := readPlaydateMemStats()
	return lastPauseNs
}

// SetMaxPause is currently unimplemented on Playdate; reserved for future use.
func (m *Memory) SetMaxPause(micros int) {
	_ = micros
}

// ResetStats zeroes the GC and realloc counters.
func (m *Memory) ResetStats() {
	resetGCStats()
	C.pd_sys_resetReallocStats()
}

// SetReallocDebug enables or disables per-realloc console logging.
func (m *Memory) SetReallocDebug(enabled bool) {
	var v C.int
	if enabled {
		v = 1
	}
	C.pd_sys_setReallocDebug(v)
}

// GetReallocStats returns current SDK realloc statistics and whether debug
// logging is enabled.
func (m *Memory) GetReallocStats() (ReallocStats, bool) {
	var count C.int
	var total C.ulong
	var freeN C.int
	debug := C.pd_sys_getReallocStats(&count, &total, &freeN)
	return ReallocStats{
		Count:      int(count),
		TotalBytes: uint64(total),
		FreeCount:  int(freeN),
	}, debug != 0
}
