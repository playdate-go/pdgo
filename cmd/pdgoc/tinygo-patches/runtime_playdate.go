//go:build playdate

package runtime

import "unsafe"

//go:extern _cgo_pd_realloc
func _cgo_pd_realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer

//go:extern _cgo_pd_getCurrentTimeMS
func _cgo_pd_getCurrentTimeMS() uint32

//go:extern _cgo_pd_logToConsole
func _cgo_pd_logToConsole(msg *byte)

// runtime_init is called from C to initialize the Go runtime
//export runtime_init
func runtime_init() {
    initAll()
}

func ticks() timeUnit {
    return timeUnit(_cgo_pd_getCurrentTimeMS()) * 1000000
}

func sleepTicks(d timeUnit) {}

func nanosecondsToTicks(ns int64) timeUnit { return timeUnit(ns) }

func ticksToNanoseconds(t timeUnit) int64 { return int64(t) }

var printBuf [256]byte
var printBufIdx int

func putchar(c byte) {
    if c == '\n' || printBufIdx >= len(printBuf)-1 {
        printBuf[printBufIdx] = 0
        _cgo_pd_logToConsole(&printBuf[0])
        printBufIdx = 0
    } else {
        printBuf[printBufIdx] = c
        printBufIdx++
    }
}
