//go:build playdate

package runtime

import "unsafe"

// External Playdate functions - implemented in pd_cgo.c
// Note: These use the runtime. prefix because they're exported from C
//
//go:extern runtime._cgo_pd_realloc
func _cgo_pd_realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer

//go:extern runtime._cgo_pd_getCurrentTimeMS
func _cgo_pd_getCurrentTimeMS() uint32

//go:extern runtime._cgo_pd_logToConsole
func _cgo_pd_logToConsole(msg *byte)

// pdStackTop is captured at runtime_init (earliest Go entry point).
// Uses pd- prefix to avoid collision with baremetal.go's stackTop.
var (
	pdStackTop         uintptr
	pdStackTopCaptured bool
)

// runtime_init is called from C to initialize the Go runtime
//
//export runtime_init
func runtime_init() {
	if !pdStackTopCaptured {
		// Capture the actual SP at the shallowest Go entry point. The game
		// runs on the OS-provided stack, which may live in the game RAM
		// region (0x9xxxxxxx), not SRAM. Never clamp or substitute a
		// hardcoded address: a wrong pdStackTop below SP makes
		// scanstack's range check fail and silently disables stack
		// scanning — the GC then sweeps objects reachable only from the
		// stack, causing use-after-free crashes shortly after the first
		// collection.
		pdStackTop = getCurrentStackPointer()
		pdStackTopCaptured = true
	}
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
