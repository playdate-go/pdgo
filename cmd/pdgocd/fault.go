// ARMv7-M fault status decoding for Playdate crash logs (Cortex-M7).
// Bit positions per the ARMv7-M Architecture Reference Manual: CFSR packs
// MMFSR (byte 0), BFSR (byte 1) and UFSR (halfword 2).

package main

import "fmt"

type faultBit struct {
	mask uint32
	text string
}

var (
	mmfsrBits = []faultBit{
		{0x01, "instruction access violation (MMFSR.IACCVIOL)"},
		{0x02, "data access violation (MMFSR.DACCVIOL)"},
	}
	bfsrBits = []faultBit{
		{0x100, "instruction bus error (BFSR.IBUSERR)"},
		{0x200, "precise data bus error (BFSR.PRECISERR)"},
		{0x400, "imprecise data bus error (BFSR.IMPRECISERR)"},
		{0x800, "fault on exception-return unstacking (BFSR.UNSTKERR)"},
		{0x1000, "fault on exception-entry stacking (BFSR.STKERR)"},
	}
	ufsrBits = []faultBit{
		{0x10000, "undefined instruction executed (UFSR.UNDEFINSTR)"},
		{0x20000, "invalid state: branch to a non-Thumb (even) address (UFSR.INVSTATE)"},
		{0x40000, "invalid PC on exception return (UFSR.INVPC)"},
		{0x80000, "access to disabled coprocessor (UFSR.NOCP)"},
	}
)

// decodeFault renders the fault reason lines for one crash record. It is
// pure register arithmetic; ELF-dependent context (e.g. whether pc points
// into .text) is added by the caller.
func decodeFault(c *Crash) []string {
	var lines []string
	cfsr := c.Regs["cfsr"]
	hfsr := c.Regs["hfsr"]

	var reasons []string
	for _, group := range [][]faultBit{mmfsrBits, bfsrBits, ufsrBits} {
		for _, b := range group {
			if cfsr&b.mask != 0 {
				reasons = append(reasons, b.text)
			}
		}
	}
	if len(reasons) == 0 && cfsr != 0 {
		reasons = append(reasons, fmt.Sprintf("undecoded CFSR bits 0x%08x", cfsr))
	}
	if len(reasons) > 0 {
		kind := "fault"
		switch {
		case cfsr&0xff0000 != 0:
			kind = "usage fault"
		case cfsr&0xff00 != 0:
			kind = "bus fault"
		case cfsr&0xff != 0:
			kind = "memmanage fault"
		}
		lines = append(lines, kind+": "+reasons[0])
		for _, r := range reasons[1:] {
			lines = append(lines, "  "+r)
		}
	}

	if cfsr&0x8000 != 0 { // BFARVALID
		lines = append(lines, fmt.Sprintf("faulting address bfar=0x%08x%s", c.Regs["bfar"], scsNote(c.Regs["bfar"])))
	}
	if cfsr&0x80 != 0 { // MMARVALID
		lines = append(lines, fmt.Sprintf("faulting address mmfar=0x%08x%s", c.Regs["mmfar"], scsNote(c.Regs["mmfar"])))
	}

	if hfsr&(1<<30) != 0 {
		lines = append(lines, "escalated to HardFault (HFSR.FORCED)")
	}
	if hfsr&0x2 != 0 {
		lines = append(lines, "vector table read fault (HFSR.VECTTBL)")
	}
	return lines
}

// scsNote annotates fault addresses that land in the ARM System Control
// Space — always a sign that a call/jump went to a non-code address.
func scsNote(addr uint32) string {
	if addr >= 0xE0000000 {
		return " — ARM System Control Space: a call/jump went to a non-code address"
	}
	return ""
}

// decodePSR renders the merged xPSR as a one-liner. Approximate: the low
// bits of the fault-report PSR hold the pre-fault IPSR exception number.
func decodePSR(psr uint32) string {
	mode := "thread mode"
	switch exc := psr & 0x1ff; exc {
	case 0:
	case 3:
		mode = "exception 3 (HardFault)"
	case 4:
		mode = "exception 4 (MemManage)"
	case 5:
		mode = "exception 5 (BusFault)"
	case 6:
		mode = "exception 6 (UsageFault)"
	case 11:
		mode = "exception 11 (SVCall)"
	case 14:
		mode = "exception 14 (PendSV)"
	case 15:
		mode = "exception 15 (SysTick)"
	default:
		mode = fmt.Sprintf("exception %d", exc)
	}
	state := "Thumb"
	if psr&(1<<24) == 0 {
		// On Cortex-M the T bit is always set for executing Thumb code; a
		// cleared T bit in the faulted PSR means a blx/branch to an even
		// (non-Thumb) address switched the core to ARM state — the classic
		// UFSR.INVSTATE signature.
		state = "ARM state (T-bit clear — attempted non-Thumb execution)"
	}
	return state + ", " + mode
}
