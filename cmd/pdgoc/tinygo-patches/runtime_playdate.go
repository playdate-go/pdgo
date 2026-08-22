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

// pdStackTop is the highest stack address the conservative root scan may
// read. The game runs on the OS-provided stack, which may live in the game
// RAM region (0x9xxxxxxx), not SRAM, and its top is unknown — so it is
// tracked as the maximum SP ever observed at a C→Go entry point.
// runtime_init establishes it (kEventInit dispatch); runtime_note_stack_top
// raises it at every later entry (per-frame update wrapper, event handler),
// because the SDK dispatches those callbacks at a SHALLOWER depth than the
// kEventInit dispatch — with only the init capture, scanstack's
// pdStackTop > sp range check fails and the stack scan silently contributes
// zero roots, so the GC sweeps objects reachable only from the current call
// chain (stack-rooted locals) while they are still in use.
// Uses pd- prefix to avoid collision with baremetal.go's stackTop.
var (
	pdStackTop         uintptr
	pdStackTopCaptured bool
	pdStackTopInit     uintptr // TEMPORARY diagnostic: original runtime_init capture
)

// runtime_init is called from C to initialize the Go runtime
//
//export runtime_init
func runtime_init() {
	if !pdStackTopCaptured {
		pdStackTop = getCurrentStackPointer()
		pdStackTopInit = pdStackTop
		pdStackTopCaptured = true
	}
	initAll()
}

// runtime_note_stack_top is called from C at every entry into Go code
// (update_callback_wrapper per frame, eventHandler per event). Raises
// pdStackTop to the highest SP seen at an entry so a GC running inside a
// callback dispatched shallower than the kEventInit capture still scans the
// frames between the collector and that callback.
//
//export runtime_note_stack_top
func noteStackTop() {
	if sp := getCurrentStackPointer(); sp > pdStackTop {
		pdStackTop = sp
	}
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
