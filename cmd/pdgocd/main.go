// pdgocd symbolizes Playdate device crash logs against a game ELF.
//
// Usage:
//
//	pdgocd [-e <pdex.elf|game dir>] [-d] [-log <text>] [crashlog.txt]
//
// The crash log and the ELF source may also be passed together as two
// positional arguments, in either order. The log itself comes from a
// file, stdin, or the -log flag (raw text), but only one of those.
//
// The crash log may be piped on stdin instead of passed as a file. Device
// addresses in the 0x9000_0000 window are translated to ELF offsets (the
// game ELF links at 0 and the device maps it at 0x90000000) and resolved
// with arm-none-eabi-addr2line, falling back to raw ELF symbol names.
//
// The ELF source is an ELF file or a directory searched for
// build/pdex.elf and pdex.elf. A .pdx bundle is not a usable source: it
// contains only the pdc-encrypted pdex.bin. Build with
// 'pdgoc -device -keep' to keep the game's build/pdex.elf.
package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const usageText = `Usage: pdgocd [-e <elf>] [-d] [-log <text>] [crashlog.txt] [<pdex.elf|game dir>]

Symbolizes Playdate device crash logs against a game ELF.
The crash log comes from a file argument, stdin, or -log (raw text) —
but only one of those. A positional argument that is an ELF file or a
directory (searched for build/pdex.elf, pdex.elf) is used as the ELF
source; the crash log and the ELF source may be passed together in
either order.

  -e string   game ELF file, or a directory searched for
              build/pdex.elf and pdex.elf
  -log string raw crash-log text, instead of a log file or stdin
  -d          disassemble instructions around the faulting pc

Exit codes: 0 analyzed, 2 no crash found / bad input, 3 no usable ELF.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, out, errw io.Writer) int {
	fs := flag.NewFlagSet("pdgocd", flag.ContinueOnError)
	fs.SetOutput(errw)
	elfFlag := fs.String("e", "", "game ELF file, or a directory searched for build/pdex.elf and pdex.elf")
	logFlag := fs.String("log", "", "raw crash-log text, instead of a log file or stdin")
	disasm := fs.Bool("d", false, "disassemble instructions around the faulting pc")
	if err := fs.Parse(args); err != nil {
		fmt.Fprint(errw, usageText)
		return 2
	}
	if fs.NArg() > 2 {
		fmt.Fprint(errw, usageText)
		return 2
	}

	// Positional arguments are the crash log and/or the ELF source,
	// told apart by the ELF source being a directory or having ELF/PDX
	// magic. With two, either order is accepted.
	isElfSrc := func(p string) bool {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			return true
		}
		m := fileMagic(p)
		return m == "elf" || m == "pdx"
	}
	var logPath, elfHint string
	switch fs.NArg() {
	case 1:
		if isElfSrc(fs.Arg(0)) {
			elfHint = fs.Arg(0)
		} else {
			logPath = fs.Arg(0)
		}
	case 2:
		a, b := fs.Arg(0), fs.Arg(1)
		switch {
		case isElfSrc(a) && !isElfSrc(b):
			elfHint, logPath = a, b
		case isElfSrc(b) && !isElfSrc(a):
			elfHint, logPath = b, a
		default:
			fmt.Fprint(errw, usageText)
			return 2
		}
	}

	if *logFlag != "" && logPath != "" {
		fmt.Fprintln(errw, "pdgocd: -log and a crash log file argument are mutually exclusive")
		fmt.Fprint(errw, usageText)
		return 2
	}

	var data []byte
	var err error
	switch {
	case *logFlag != "":
		data = []byte(*logFlag)
	case logPath != "":
		data, err = os.ReadFile(logPath)
	default:
		data, err = io.ReadAll(stdin)
	}
	if err != nil {
		fmt.Fprintf(errw, "pdgocd: reading input: %v\n", err)
		return 2
	}

	crashes := ParseCrashes(string(data))
	if len(crashes) == 0 {
		fmt.Fprintln(errw, "pdgocd: no crash records found in input")
		return 2
	}

	elfPath, err := findElf(*elfFlag, elfHint)
	if err != nil {
		fmt.Fprintf(errw, "pdgocd: %v\n", err)
		return 3
	}
	img, err := OpenElf(elfPath)
	if err != nil {
		fmt.Fprintf(errw, "pdgocd: %v\n", err)
		return 3
	}

	var sym Symbolizer
	if a, warn := NewAddr2line(elfPath); a != nil {
		sym = a
	} else {
		fmt.Fprintln(errw, "pdgocd:", warn)
	}

	var elfMod time.Time
	if fi, err := os.Stat(elfPath); err == nil {
		elfMod = fi.ModTime()
	}

	// Collect every register value that lands in the ELF window and
	// resolve them in one addr2line batch.
	offsetSet := map[uint64]bool{}
	for _, c := range crashes {
		for _, name := range []string{"pc", "lr", "r0", "r1", "r2", "r3", "r12"} {
			if off, ok := img.Offset(c.Regs[name]); ok {
				offsetSet[off] = true
			}
		}
	}
	offsets := make([]uint64, 0, len(offsetSet))
	for off := range offsetSet {
		offsets = append(offsets, off)
	}
	frames := resolveAll(img, sym, offsets)

	for i, c := range crashes {
		if i > 0 {
			fmt.Fprintln(out)
			fmt.Fprintln(out, strings.Repeat("-", 60))
			fmt.Fprintln(out)
		}
		renderCrash(out, c, i+1, img, elfPath, elfMod, frames)
		if *disasm {
			renderDisasm(out, errw, img, elfPath, c)
		}
	}
	return 0
}

func renderCrash(out io.Writer, c *Crash, n int, img *ElfImage, elfPath string, elfMod time.Time, frames map[uint64][]Frame) {
	fmt.Fprintf(out, "Crash #%d - %s\n", n, c.When)
	if c.Build != "" {
		fmt.Fprintf(out, "  build: %s\n", c.Build)
	}
	fmt.Fprintf(out, "  ELF:   %s", elfPath)
	if !elfMod.IsZero() {
		fmt.Fprintf(out, " (modified %s)", elfMod.Format("2006-01-02 15:04"))
	}
	fmt.Fprintln(out)

	// Crash timestamps are device-local and mtimes host-local; this is a
	// heuristic, but an ELF newer than the crash usually means a rebuild
	// since, i.e. drifted symbols.
	if when, err := time.Parse("2006/01/02 15:04:05", c.When); err == nil && elfMod.After(when) {
		fmt.Fprintf(out, "  WARNING: ELF is newer than the crash - symbols may have drifted\n")
	}
	if w := symbolizationWarning(c, img, frames); w != "" {
		fmt.Fprintf(out, "  WARNING: %s\n", w)
	}
	fmt.Fprintln(out)

	for _, ln := range decodeFault(c) {
		fmt.Fprintf(out, "  %s\n", ln)
	}
	if psr, ok := c.Regs["psr"]; ok {
		fmt.Fprintf(out, "  psr %08x: %s\n", psr, decodePSR(psr))
	}
	// The full register listing — the crash's stack dump — wrapped in
	// slash markers so it stands out from the decode prose.
	fmt.Fprintln(out)
	fmt.Fprintln(out, strings.Repeat("/", 60))
	for _, name := range regNames {
		v, ok := c.Regs[name]
		if !ok {
			continue
		}
		switch name {
		case "psr", "cfsr", "hfsr", "mmfar", "bfar", "rcccsr":
			// Already covered by the fault decode above.
			fmt.Fprintf(out, "  %-6s %08x\n", name, v)
			continue
		}
		if off, ok := img.Offset(v); ok {
			fmt.Fprintf(out, "  %-6s %08x -> %05x", name, v, off)
			if s, exec := img.SectionAt(off); s != "" && !exec {
				fmt.Fprintf(out, "  [%s]", s)
			}
			if f := frames[off]; len(f) > 0 {
				fmt.Fprintf(out, "  %s", frameStr(f))
			}
			if name == "lr" {
				fmt.Fprint(out, "  (return address: caller)")
			}
			fmt.Fprintln(out)
		} else if v >= 0x20000000 && v < 0x20040000 {
			fmt.Fprintf(out, "  %-6s %08x  (SRAM: globals/stack)\n", name, v)
		} else if v >= playdateLoadBase && v < flashWindowEnd {
			fmt.Fprintf(out, "  %-6s %08x  (flash: past this ELF's image - SDK heap, or a different build)\n", name, v)
		} else if v >= 0xE0000000 {
			fmt.Fprintf(out, "  %-6s %08x  (ARM system space)\n", name, v)
		} else {
			fmt.Fprintf(out, "  %-6s %08x\n", name, v)
		}
	}
	fmt.Fprintln(out, strings.Repeat("/", 60))

	for _, hint := range crashHints(c, img) {
		fmt.Fprintf(out, "  ! %s\n", hint)
	}
	if c.HeapAlloc > 0 {
		fmt.Fprintf(out, "  heap allocated: %d bytes\n", c.HeapAlloc)
	}
}

// symbolizationWarning reports when the crash addresses cannot be trusted
// as Go function names. Two cases:
//
//   - pc/lr in the flash window but past this ELF's loaded image: the
//     crashing game mapped code this ELF does not contain, so it is from
//     a different game/build (or execution ran into the flash SDK heap).
//   - in-window pc/lr resolving only to nearest-symbol guesses (no DWARF
//     answer, no exact symbol, not inside any function): TinyGo symbol
//     tables are dense across the whole image, so a fallback-only match
//     usually means the ELF belongs to a different game or is stripped.
func symbolizationWarning(c *Crash, img *ElfImage, frames map[uint64][]Frame) string {
	var beyond []string
	candidates, solid := 0, 0
	for _, name := range []string{"pc", "lr"} {
		a := c.Regs[name]
		if a < playdateLoadBase || a >= flashWindowEnd {
			continue
		}
		off, ok := img.Offset(a)
		if !ok {
			beyond = append(beyond, fmt.Sprintf("%s=%08x", name, a))
			continue
		}
		candidates++
		if f := frames[off]; len(f) > 0 && !f[0].Approx {
			solid++
		}
	}
	if len(beyond) > 0 {
		return fmt.Sprintf("%s lies past the end of this ELF's loaded image (image ends at 0x%x) - this ELF is from a different game/build than the crashing one, or execution ran outside game code (e.g. into the flash-mapped SDK heap)", strings.Join(beyond, ", "), img.LoadEnd)
	}
	if candidates > 0 && solid == 0 {
		return "cannot resolve the crash addresses to Go function names - this ELF is probably from a different game/build than the crashing one, or its symbols are stripped"
	}
	return ""
}

// crashHints names patterns that keep recurring in real Playdate crashes:
// indirect calls through a register and execution landing in data.
func crashHints(c *Crash, img *ElfImage) []string {
	var hints []string
	pc := c.Regs["pc"]
	if off, ok := img.Offset(pc); ok {
		if s, ok := img.ExactSymbol(off); ok && elf.ST_TYPE(s.Info) != elf.STT_FUNC {
			hints = append(hints, fmt.Sprintf("pc sits exactly on data symbol %q, not a function - a non-code value (e.g. a TinyGo $pack interface constant) was called as code", s.Name))
		} else if sec, exec := img.SectionAt(off); sec != "" && !exec {
			hints = append(hints, fmt.Sprintf("pc is in %s, not executable code - execution jumped into data (e.g. blx through a non-function pointer)", sec))
		}
		for _, rn := range []string{"r1", "r2", "r3", "r12"} {
			if c.Regs[rn] == pc {
				hints = append(hints, fmt.Sprintf("pc == %s - indirect call (blx %s) through that register", rn, rn))
			}
		}
	}
	return hints
}

func renderDisasm(out, errw io.Writer, img *ElfImage, elfPath string, c *Crash) {
	off, ok := img.Offset(c.Regs["pc"])
	if !ok {
		return
	}
	bin, err := exec.LookPath("arm-none-eabi-objdump")
	if err != nil {
		fmt.Fprintln(errw, "pdgocd: arm-none-eabi-objdump not found on PATH; skipping disassembly")
		return
	}
	start := uint64(0)
	if off > 24 {
		start = off - 24
	}
	cmd := exec.Command(bin, "-d", "--start-address="+fmt.Sprint(start), "--stop-address="+fmt.Sprint(off+16), elfPath)
	dis, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(errw, "pdgocd: objdump failed: %v\n", err)
		return
	}
	fmt.Fprintln(out, "\n  disassembly around pc:")
	for _, ln := range strings.Split(strings.TrimRight(string(dis), "\n"), "\n") {
		fmt.Fprintf(out, "  %s\n", ln)
	}
}

func frameStr(fs []Frame) string {
	parts := make([]string, len(fs))
	for i, f := range fs {
		s := f.Name
		if f.Approx {
			s += " (nearest symbol)"
		} else if f.Line > 0 {
			s += fmt.Sprintf(" (%s:%d)", f.File, f.Line)
		}
		parts[i] = s
	}
	return strings.Join(parts, " [inlined] ")
}

// findElf locates a symbolizable ELF. explicit comes from -e, hint from the
// positional argument; both may be an ELF file or a directory searched for
// build/pdex.elf and pdex.elf. Without either, the search walks up from
// the working directory. A .pdx bundle is never a source: it contains only
// the pdc-encrypted pdex.bin.
func findElf(explicit, hint string) (string, error) {
	searchDir := func(dir string) (string, bool) {
		for _, cand := range []string{
			filepath.Join(dir, "build", "pdex.elf"), // kept by 'pdgoc -device -keep'
			filepath.Join(dir, "pdex.elf"),
		} {
			if fileMagic(cand) == "elf" {
				return cand, true
			}
		}
		return "", false
	}

	// An explicit -e path wins; the positional hint is the fallback root.
	root := explicit
	if root == "" {
		root = hint
	}
	if root != "" {
		if fi, err := os.Stat(root); err != nil {
			return "", fmt.Errorf("path not found: %s", root)
		} else if !fi.IsDir() {
			switch fileMagic(root) {
			case "elf":
				return root, nil
			case "pdx":
				return "", errEncrypted(root)
			default:
				return "", fmt.Errorf("%s is not an ELF file", root)
			}
		}
		if strings.HasSuffix(strings.TrimSuffix(root, string(filepath.Separator)), ".pdx") {
			return "", errBundle(root)
		}
		if found, ok := searchDir(root); ok {
			return found, nil
		}
		return "", fmt.Errorf("no pdex.elf found in %s (looked for build/pdex.elf, pdex.elf); pass -e <path-to-pdex.elf>", root)
	}

	cwd, _ := os.Getwd()
	for dir := cwd; ; dir = filepath.Dir(dir) {
		if found, ok := searchDir(dir); ok {
			return found, nil
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	return "", fmt.Errorf("no pdex.elf found (searched build/pdex.elf, pdex.elf upward from %s); pass -e <path>", cwd)
}

func errEncrypted(p string) error {
	return fmt.Errorf("%s is a pdc-encrypted bundle binary (\"Playdate PDX\" header) and cannot be symbolized; pass -e <path-to-pdex.elf>", p)
}

func errBundle(p string) error {
	return fmt.Errorf("%s is a .pdx bundle: it contains only the pdc-encrypted pdex.bin and cannot be symbolized; pass the game's build/pdex.elf (kept by 'pdgoc -device -keep')", p)
}
