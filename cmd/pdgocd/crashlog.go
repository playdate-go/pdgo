// pdgocd crash log parser.
//
// Playdate device crash logs are pasted from the device's console and come
// in two shapes: the multi-line register dump printed by the OS crash
// handler, and compact one-line re-wraps of the same data. Both are observed
// in the wild, so the parser anchors on the "crash at" header first and
// falls back to an unanchored register scan per block.

package main

import (
	"regexp"
	"strconv"
)

// Crash is a single crash record parsed from a device log.
type Crash struct {
	When      string // timestamp text from the "crash at ..." header line
	Build     string // OS build id from the "build:" line
	Regs      map[string]uint32
	HeapAlloc uint64
}

// regNames is the dump order the device uses, for stable rendering.
var regNames = []string{
	"r0", "r1", "r2", "r3", "r12", "lr", "pc", "psr",
	"cfsr", "hfsr", "mmfar", "bfar", "rcccsr",
}

var (
	crashHeaderRe = regexp.MustCompile(`(?im)^-{0,3}\s*crash at\s*(.+?)\s*-*\s*$`)
	tsRe          = regexp.MustCompile(`\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2}`)
	buildRe       = regexp.MustCompile(`(?im)^\s*build:\s*(.+?)\s*$`)
	heapRe        = regexp.MustCompile(`(?im)heap allocated:\s*(\d+)`)

	// lineRegRe matches the device's one-register-per-line dump. Anchoring
	// to line starts keeps Lua/heap lines from false-positive matches.
	lineRegRe = regexp.MustCompile(`(?im)^\s*(r\d{1,2}|lr|pc|psr|cfsr|hfsr|mmfar|bfar|rcccsr)\s*:\s*([0-9a-fA-F]{8})\b`)

	// inlineRegRe matches compact one-line dumps ("r0:900371d8 r1:...").
	inlineRegRe = regexp.MustCompile(`\b(r\d{1,2}|lr|pc|psr|cfsr|hfsr|mmfar|bfar|rcccsr)\s*:\s*([0-9a-fA-F]{8})\b`)
)

// ParseCrashes splits a device log into crash records. Input without a
// "crash at" header but with register lines is accepted as a single
// anonymous record (partial paste).
func ParseCrashes(data string) []*Crash {
	headers := crashHeaderRe.FindAllStringSubmatch(data, -1)

	var crashes []*Crash
	if len(headers) == 0 {
		if c := parseBlock(data, ""); c != nil {
			crashes = append(crashes, c)
		}
		return crashes
	}

	matches := crashHeaderRe.FindAllStringSubmatchIndex(data, -1)
	for i, m := range matches {
		when := headers[i][1]
		start := m[1] // end of the header line
		// Registers may share the header line (one-line pastes): resume
		// the block right after the timestamp instead of dropping the
		// rest of the line, and normalize When to the bare timestamp.
		headerLine := data[m[0]:m[1]]
		if ts := tsRe.FindStringIndex(headerLine); ts != nil {
			when = headerLine[ts[0]:ts[1]]
			start = m[0] + ts[1]
		}
		end := len(data)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		if c := parseBlock(data[start:end], when); c != nil {
			crashes = append(crashes, c)
		}
	}
	return crashes
}

func parseBlock(block, when string) *Crash {
	regs := map[string]uint32{}
	for _, m := range lineRegRe.FindAllStringSubmatch(block, -1) {
		if v, err := strconv.ParseUint(m[2], 16, 32); err == nil {
			regs[m[1]] = uint32(v)
		}
	}
	// Compact one-line dumps have no line structure to anchor on; rescan
	// unanchored when the line pass found too little.
	if len(regs) < 5 {
		for _, m := range inlineRegRe.FindAllStringSubmatch(block, -1) {
			if v, err := strconv.ParseUint(m[2], 16, 32); err == nil {
				regs[m[1]] = uint32(v)
			}
		}
	}
	if len(regs) < 3 {
		return nil
	}

	c := &Crash{When: when, Regs: regs}
	if m := buildRe.FindStringSubmatch(block); m != nil {
		c.Build = m[1]
	}
	if m := heapRe.FindStringSubmatch(block); m != nil {
		c.HeapAlloc, _ = strconv.ParseUint(m[1], 10, 64)
	}
	return c
}
