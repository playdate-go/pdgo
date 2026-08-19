//go:build gc.playdate && !tinygo.wasm && !scheduler.threads && !scheduler.cores

package runtime

// pdStackTop and pdStackTopCaptured are declared in runtime_playdate.go
// and captured at runtime_init (the earliest Go entry point).

func gcMarkReachable() {
	if !pdStackTopCaptured {
		return // nothing we can do
	}
	markStack()
	findGlobals(markRoots)
}

func markStack() {
	scanCurrentStack()
}

// scanCurrentStack is implemented in TinyGo's stock src/runtime/asm_arm.S:
// it pushes callee-saved registers onto the stack, then calls
// tinygo_scanstack. The body-less declaration plus //go:export binds this
// Go function to the assembly symbol — the stock TinyGo pattern (see
// gc_stack_raw.go). With its own (empty) Go body instead, the asm was never
// called and stack scanning silently became a no-op.
//
// TinyGo only assembles extra-files like asm_arm.S when it links an
// executable itself; for `-o game.o` output it emits just the Go module.
// The device build scripts therefore assemble asm_arm.S and link it
// alongside game.o.
//
//go:export tinygo_scanCurrentStack
func scanCurrentStack()

// scanstack is called from assembly after pushing callee-saved registers.
//
//go:export tinygo_scanstack
func scanstack(sp uintptr) {
	if pdStackTop > sp {
		markRoots(sp, pdStackTop)
	}
}

func gcResumeWorld() {} // single-threaded
