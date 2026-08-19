//go:build cortexm

package interrupt

import (
	"device/arm"
)

// Enable enables this interrupt. Right after calling this function, the
// interrupt may be invoked if it was already pending.
func (irq Interrupt) Enable() {
	// Clear the ARM pending bit, an asserting device may still
	// trigger the interrupt once enabled.
	arm.ClearPendingIRQ(uint32(irq.num))
	arm.EnableIRQ(uint32(irq.num))
}

// Disable disables this interrupt.
func (irq Interrupt) Disable() {
	arm.DisableIRQ(uint32(irq.num))
}

// SetPriority sets the interrupt priority for this interrupt. A lower number
// means a higher priority. Additionally, most hardware doesn't implement all
// priority bits (only the uppoer bits).
// Examples: 0xff (lowest priority), 0xc0 (low), 0x00 (highest possible).
func (irq Interrupt) SetPriority(priority uint8) {
	arm.SetPriority(uint32(irq.num), uint32(priority))
}

// State represents the previous global interrupt state.
type State uintptr

// Disable disables all interrupts and returns the previous interrupt state. It
// can be used in a critical section like this:
//
//	state := interrupt.Disable()
//	// critical section
//	interrupt.Restore(state)
//
// Critical sections can be nested. Make sure to call Restore in the same order
// as you called Disable (this happens naturally with the pattern above).
func Disable() (state State) {
	return State(arm.DisableInterrupts())
}

// Restore restores interrupts to what they were before. Give the previous state
// returned by Disable as a parameter. If interrupts were disabled before
// calling Disable, this will not re-enable them, allowing for nested critical
// sections.
func Restore(state State) {
	arm.EnableInterrupts(uintptr(state))
}

// In returns whether the system is currently in an interrupt.
//
// PLAYDATE PATCH: the stock implementation reads SCB.ICSR (0xE000ED04) to
// check the VECTACTIVE field. The Playdate OS runs game code unprivileged,
// and any System Control Space access from unprivileged code raises a
// BusFault that escalates to HardFault — the system crash screen. This
// fires inside panicOrGoexit() *before* the panic message is printed, so
// every Go panic surfaces as an immediate crash with no diagnostics
// (pc inside panicOrGoexit, bfar=0xE000ED04).
//
// Playdate callbacks (eventHandler, update, etc.) always run in thread mode
// (IPSR=0), and this runtime never installs interrupt handlers, so "in
// interrupt" is always false here.
func In() bool {
	return false
}
