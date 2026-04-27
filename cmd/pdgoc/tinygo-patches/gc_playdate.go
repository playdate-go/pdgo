//go:build gc.playdate

package runtime

import (
    "unsafe"
)

const needsStaticHeap = false

var gcTotalAlloc uint64
var gcMallocs uint64
var gcFrees uint64

//go:noinline
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
    size = align(size)
    gcTotalAlloc += uint64(size)
    gcMallocs++

    ptr := _cgo_pd_realloc(nil, size)
    if ptr == nil {
        runtimePanic("out of memory")
    }

    memzero(ptr, size)
    return ptr
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
    size = align(size)
    newPtr := _cgo_pd_realloc(ptr, size)
    if newPtr == nil && size > 0 {
        runtimePanic("out of memory")
    }
    return newPtr
}

func free(ptr unsafe.Pointer) {
    if ptr != nil {
        _cgo_pd_realloc(ptr, 0)
        gcFrees++
    }
}

func markRoots(start, end uintptr) {}

func ReadMemStats(m *MemStats) {
    m.HeapIdle = 0
    m.HeapInuse = gcTotalAlloc
    m.HeapReleased = 0
    m.HeapSys = m.HeapInuse + m.HeapIdle
    m.GCSys = 0
    m.TotalAlloc = gcTotalAlloc
    m.Mallocs = gcMallocs
    m.Frees = gcFrees
    m.Sys = gcTotalAlloc
    m.HeapAlloc = gcTotalAlloc
    m.Alloc = m.HeapAlloc
}

func GC() {}

func SetFinalizer(obj interface{}, finalizer interface{}) {}

func initHeap() {}

func setHeapEnd(newHeapEnd uintptr) {}
