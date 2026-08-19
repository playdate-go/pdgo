// Symbolization: batch addr2line over device addresses translated to ELF
// offsets, with a pure-Go symtab fallback (debug/elf) for when addr2line is
// missing or returns ?? for an address.

package main

import (
	"bytes"
	"debug/elf"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// playdateLoadBase is where the device maps the game ELF, which links at 0.
const playdateLoadBase = 0x90000000

// flashWindowEnd bounds the QSPI XIP range the game's flash is mapped
// into. Device addresses in [playdateLoadBase, flashWindowEnd) that fall
// past this ELF's load extent are still game-flash references (the SDK
// heap arena sits right after the loaded image).
const flashWindowEnd = 0xA0000000

// Frame is one symbolized frame; an address may resolve to a chain of
// inlined frames (outermost first, as addr2line -i prints them). Approx
// marks a nearest-symtab guess, not a real resolution.
type Frame struct {
	Name   string
	File   string
	Line   int
	Approx bool
}

// Symbolizer resolves ELF offsets to frame chains. Addresses that resolve
// to nothing are absent from the returned map.
type Symbolizer interface {
	Resolve(offsets []uint64) map[uint64][]Frame
}

// sectionRange is an ELF section mapped by vaddr (ELF is linked at 0).
type sectionRange struct {
	Name       string
	Start, End uint64
	Exec       bool // SHF_EXECINSTR
}

// ElfImage is the parsed game ELF: load window for classifying device
// addresses, section table for naming where pc landed, and the symbol
// table for the no-addr2line fallback.
type ElfImage struct {
	Path     string
	LoadEnd  uint64 // end of the highest PT_LOAD (max vaddr+memsz)
	sections []sectionRange
	syms     []elf.Symbol // sorted by Value
}

// OpenElf parses the ELF file at path and precomputes lookup tables.
func OpenElf(path string) (*ElfImage, error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img := &ElfImage{Path: path}
	for _, seg := range f.Progs {
		if seg.Type != elf.PT_LOAD {
			continue
		}
		if end := seg.Vaddr + seg.Memsz; end > img.LoadEnd {
			img.LoadEnd = end
		}
	}
	for _, s := range f.Sections {
		// Only ALLOC sections occupy the loaded image; debug/symtab
		// sections would otherwise shadow real ones. Note the game ELF
		// links at 0, so Addr==0 is valid (the playdate.ld script even
		// folds code and constants into a single .text at 0).
		if s.Flags&elf.SHF_ALLOC == 0 || s.Size == 0 {
			continue
		}
		img.sections = append(img.sections, sectionRange{
			Name:  s.Name,
			Start: s.Addr,
			End:   s.Addr + s.Size,
			Exec:  s.Flags&elf.SHF_EXECINSTR != 0,
		})
	}
	if syms, err := f.Symbols(); err == nil {
		defined := make([]elf.Symbol, 0, len(syms))
		for _, s := range syms {
			if s.Value != 0 && s.Section != elf.SHN_UNDEF {
				defined = append(defined, s)
			}
		}
		sort.Slice(defined, func(i, j int) bool { return defined[i].Value < defined[j].Value })
		img.syms = defined
	}
	return img, nil
}

// Offset maps a device address to an ELF offset. The second result is
// false when the address is outside the mapped game image.
func (img *ElfImage) Offset(addr uint32) (uint64, bool) {
	if uint64(addr) < playdateLoadBase {
		return 0, false
	}
	off := uint64(addr) - playdateLoadBase
	if off >= img.LoadEnd {
		return 0, false
	}
	return off, true
}

// SectionAt names the section containing the ELF offset (".text",
// ".rodata", ...) and whether that section is executable.
func (img *ElfImage) SectionAt(off uint64) (string, bool) {
	for _, s := range img.sections {
		if off >= s.Start && off < s.End {
			return s.Name, s.Exec
		}
	}
	return "", false
}

// ExactSymbol returns the symbol located exactly at the ELF offset, if any.
// An exact non-function symbol (e.g. a TinyGo $pack constant) is stronger
// evidence than addr2line, whose DWARF line tables can attribute trailing
// data addresses to the preceding function.
func (img *ElfImage) ExactSymbol(off uint64) (elf.Symbol, bool) {
	for _, s := range img.syms {
		if s.Value == off {
			return s, true
		}
		if s.Value > off {
			break
		}
	}
	return elf.Symbol{}, false
}

// NearestSymbol returns "symbol+0xoff" for the closest symbol at or below
// the ELF offset — the nm-style fallback when addr2line has no answer.
func (img *ElfImage) NearestSymbol(off uint64) (string, bool) {
	i := sort.Search(len(img.syms), func(i int) bool { return img.syms[i].Value > off }) - 1
	if i < 0 {
		return "", false
	}
	s := img.syms[i]
	if off-s.Value > 1<<20 {
		return "", false // suspiciously far: more likely noise than a match
	}
	if off == s.Value {
		return s.Name, true
	}
	return fmt.Sprintf("%s+0x%x", s.Name, off-s.Value), true
}

// FuncAt returns the STT_FUNC symbol whose [Value, Value+Size) range
// contains the ELF offset — the address is genuinely inside that
// function's code, not merely past its end in a gap or a data section.
func (img *ElfImage) FuncAt(off uint64) (elf.Symbol, bool) {
	i := sort.Search(len(img.syms), func(i int) bool { return img.syms[i].Value > off }) - 1
	for ; i >= 0; i-- {
		s := img.syms[i]
		if elf.ST_TYPE(s.Info) != elf.STT_FUNC || s.Size == 0 {
			continue
		}
		// Symbols do not overlap, so the first sized function at or
		// below off decides: inside it, or in a gap.
		if off < s.Value+s.Size {
			return s, true
		}
		return elf.Symbol{}, false
	}
	return elf.Symbol{}, false
}

// Addr2line is a Symbolizer backed by one arm-none-eabi-addr2line
// invocation (batched: all offsets on stdin).
type Addr2line struct {
	Bin string // addr2line executable
	ELF string
}

// NewAddr2line returns an Addr2line using the toolchain binary from PATH,
// or nil (with the reason) when it is not installed.
func NewAddr2line(elfPath string) (*Addr2line, string) {
	bin, err := exec.LookPath("arm-none-eabi-addr2line")
	if err != nil {
		return nil, "arm-none-eabi-addr2line not found on PATH; falling back to ELF symbol names"
	}
	return &Addr2line{Bin: bin, ELF: elfPath}, ""
}

var addrLineRe = regexp.MustCompile(`^0x[0-9a-fA-F]+$`)

func (a *Addr2line) Resolve(offsets []uint64) map[uint64][]Frame {
	if len(offsets) == 0 {
		return map[uint64][]Frame{}
	}
	var in strings.Builder
	for _, off := range offsets {
		fmt.Fprintf(&in, "0x%x\n", off)
	}

	cmd := exec.Command(a.Bin, "-e", a.ELF, "-a", "-f", "-i", "-C")
	cmd.Stdin = strings.NewReader(in.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil // caller falls back to the symtab
	}
	return parseAddr2line(out.String())
}

// parseAddr2line parses "-a -f -i -C" output: an address line, then
// alternating function/file:line pairs until the next address line.
func parseAddr2line(out string) map[uint64][]Frame {
	res := map[uint64][]Frame{}
	lines := strings.Split(out, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !addrLineRe.MatchString(line) {
			continue
		}
		addr, err := strconv.ParseUint(line, 0, 64)
		if err != nil {
			continue
		}
		var frames []Frame
		for i+1 < len(lines) {
			name := strings.TrimSpace(lines[i+1])
			if addrLineRe.MatchString(name) || name == "" {
				break // next address block (trailing empty line)
			}
			if i+2 >= len(lines) {
				break
			}
			loc := strings.TrimSpace(lines[i+2])
			file, ln := parseLocation(loc)
			if name != "??" {
				frames = append(frames, Frame{Name: name, File: file, Line: ln})
			}
			i += 2
		}
		if len(frames) > 0 {
			res[addr] = frames
		}
	}
	return res
}

func parseLocation(loc string) (string, int) {
	// "file.go:157", "file.go:157:9" or "??:0"; split from the right so
	// paths never get chopped.
	if i := strings.LastIndex(loc, ":"); i >= 0 {
		if j := strings.LastIndex(loc[:i], ":"); j >= 0 {
			loc = loc[:i]
			i = j
		}
		ln, err := strconv.Atoi(loc[i+1:])
		if err == nil {
			return loc[:i], ln
		}
	}
	return loc, 0
}

// resolveAll merges addr2line results with the symtab fallback: every
// offset must come back with something, or rendering gets holes. An exact
// non-function symbol overrides addr2line — a data label at the address
// (a TinyGo $pack constant) beats a DWARF attribution to the preceding
// function. Addresses inside a function's symtab range resolve solidly;
// the bare nearest-symbol match is marked Approx so consumers can tell a
// guess from a resolution.
func resolveAll(img *ElfImage, sym Symbolizer, offsets []uint64) map[uint64][]Frame {
	res := map[uint64][]Frame{}
	if sym != nil {
		res = sym.Resolve(offsets)
	}
	for _, off := range offsets {
		if s, ok := img.ExactSymbol(off); ok && elf.ST_TYPE(s.Info) != elf.STT_FUNC {
			res[off] = []Frame{{Name: s.Name}}
			continue
		}
		if _, ok := res[off]; ok {
			continue
		}
		if s, ok := img.FuncAt(off); ok {
			res[off] = []Frame{{Name: s.Name}}
			continue
		}
		if name, ok := img.NearestSymbol(off); ok {
			res[off] = []Frame{{Name: name, Approx: true}}
		}
	}
	return res
}

// fileMagic reports "elf", "pdx" or "" for the file at path — enough to
// steer discovery without parsing headers fully.
func fileMagic(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 16)
	n, _ := f.Read(buf)
	if n >= 4 && bytes.Equal(buf[:4], []byte{0x7f, 'E', 'L', 'F'}) {
		return "elf"
	}
	if n >= 12 && string(buf[:12]) == "Playdate PDX" {
		return "pdx"
	}
	return ""
}
