package main

import (
	"debug/elf"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseAddr2line(t *testing.T) {
	out := "" +
		"0x1afd\n" +
		"runtime.runFinalizerSafe\n" +
		"/Users/rom/pdgo/cmd/pdgoc/tinygo-patches/gc_finalizer_playdate.go:157\n" +
		"0x5600\n" +
		"??\n" +
		"??:0\n" +
		"0x3000\n" +
		"runtime.gcMarkReachable\n" +
		"/Users/rom/pdgo/gc_mark.go:40\n" +
		"runtime.markRoots\n" +
		"/Users/rom/pdgo/gc_mark.go:12\n"

	res := parseAddr2line(out)
	if len(res) != 2 {
		t.Fatalf("got %d resolved addresses, want 2 (?? address must drop out): %v", len(res), res)
	}
	frames := res[0x1afd]
	if len(frames) != 1 || frames[0].Name != "runtime.runFinalizerSafe" {
		t.Fatalf("0x1afd frames = %+v", frames)
	}
	if frames[0].Line != 157 || frames[0].File != "/Users/rom/pdgo/cmd/pdgoc/tinygo-patches/gc_finalizer_playdate.go" {
		t.Fatalf("0x1afd frame = %+v", frames[0])
	}

	chain := res[0x3000]
	if len(chain) != 2 {
		t.Fatalf("inline chain frames = %+v, want 2", chain)
	}
	if chain[0].Name != "runtime.gcMarkReachable" || chain[1].Name != "runtime.markRoots" {
		t.Fatalf("inline chain order = %+v", chain)
	}
}

func TestParseLocation(t *testing.T) {
	tests := []struct {
		loc      string
		wantFile string
		wantLine int
	}{
		{"gc.go:157", "gc.go", 157},
		{"gc.go:157:9", "gc.go", 157},
		{"/abs/path/a.go:12", "/abs/path/a.go", 12},
		{"??:0", "??", 0},
		{"noname", "noname", 0},
	}
	for _, tt := range tests {
		file, line := parseLocation(tt.loc)
		if file != tt.wantFile || line != tt.wantLine {
			t.Errorf("parseLocation(%q) = (%q, %d), want (%q, %d)", tt.loc, file, line, tt.wantFile, tt.wantLine)
		}
	}
}

func TestNearestSymbol(t *testing.T) {
	img := &ElfImage{syms: []elf.Symbol{
		{Value: 0x1000, Name: "runtime.runFinalizerSafe"},
		{Value: 0x2000, Name: "runtime.GC"},
	}}
	tests := []struct {
		off  uint64
		want string
		ok   bool
	}{
		{0x1000, "runtime.runFinalizerSafe", true},
		{0x1014, "runtime.runFinalizerSafe+0x14", true},
		{0x2000, "runtime.GC", true},
		{0x900, "", false},             // below the first symbol
		{0x2000 + 0x200000, "", false}, // suspiciously far past any symbol
	}
	for _, tt := range tests {
		got, ok := img.NearestSymbol(tt.off)
		if got != tt.want || ok != tt.ok {
			t.Errorf("NearestSymbol(%#x) = (%q, %v), want (%q, %v)", tt.off, got, ok, tt.want, tt.ok)
		}
	}
}

func TestOffsetAndSection(t *testing.T) {
	img := &ElfImage{
		LoadEnd: 0x40000,
		sections: []sectionRange{
			{Name: ".text", Start: 0x0, End: 0x2000, Exec: true},
			{Name: ".rodata", Start: 0x3000, End: 0x3400},
		},
	}
	if off, ok := img.Offset(0x90005600); !ok || off != 0x5600 {
		t.Errorf("Offset(0x90005600) = (%#x, %v)", off, ok)
	}
	if _, ok := img.Offset(0x20000000); ok {
		t.Error("SRAM address classified as ELF window")
	}
	if _, ok := img.Offset(0x90040000); ok {
		t.Error("address past LoadEnd classified as ELF window")
	}
	if s, exec := img.SectionAt(0x3200); s != ".rodata" || exec {
		t.Errorf("SectionAt(0x3200) = (%q, %v), want (.rodata, false)", s, exec)
	}
	if s, exec := img.SectionAt(0x1000); s != ".text" || !exec {
		t.Errorf("SectionAt(0x1000) = (%q, %v), want (.text, true)", s, exec)
	}
	if s, _ := img.SectionAt(0x2800); s != "" {
		t.Errorf("SectionAt(0x2800) = %q, want empty (gap between sections)", s)
	}
	// A section starting at vaddr 0 must still match offset 0.
	if s, _ := img.SectionAt(0); s != ".text" {
		t.Errorf("SectionAt(0) = %q, want .text", s)
	}
}

// stubSymbolizer returns canned frames regardless of input.
type stubSymbolizer struct {
	frames   map[uint64][]Frame
	requests []uint64
}

func (s *stubSymbolizer) Resolve(offsets []uint64) map[uint64][]Frame {
	s.requests = append(s.requests, offsets...)
	return s.frames
}

// An exact non-function symbol (a data label such as a TinyGo $pack
// constant) must override addr2line's answer, which can be DWARF noise
// for data addresses trailing a function.
func TestExactDataSymbolOverridesAddr2line(t *testing.T) {
	img := &ElfImage{LoadEnd: 0x40000, syms: []elf.Symbol{
		{Value: 0x1000, Name: "__aeabi_memcpy8", Info: byte(elf.STT_FUNC)},
		{Value: 0x5600, Name: "pdgo$pack.55", Info: byte(elf.STT_OBJECT)},
	}}
	sym := &stubSymbolizer{frames: map[uint64][]Frame{
		0x5600: {{Name: "__aeabi_memcpy8", File: "memcpy.c", Line: 9}},
		0x1000: {{Name: "__aeabi_memcpy8", File: "memcpy.c", Line: 3}},
	}}
	res := resolveAll(img, sym, []uint64{0x5600, 0x1000})
	if f := res[0x5600]; len(f) != 1 || f[0].Name != "pdgo$pack.55" {
		t.Errorf("0x5600 = %+v, want the exact data symbol pdgo$pack.55", f)
	}
	if f := res[0x1000]; len(f) != 1 || f[0].Name != "__aeabi_memcpy8" || f[0].Line != 3 {
		t.Errorf("0x1000 = %+v, want addr2line answer kept for real functions", f)
	}
}

func TestFuncAt(t *testing.T) {
	img := &ElfImage{syms: []elf.Symbol{
		{Value: 0x1000, Name: "a", Info: byte(elf.STT_FUNC), Size: 0x100},
		{Value: 0x2000, Name: "pack", Info: byte(elf.STT_OBJECT), Size: 0x10},
		{Value: 0x2100, Name: "b", Info: byte(elf.STT_FUNC), Size: 0x10},
	}}
	if s, ok := img.FuncAt(0x1080); !ok || s.Name != "a" {
		t.Errorf("FuncAt(0x1080) = (%q, %v), want a", s.Name, ok)
	}
	if s, ok := img.FuncAt(0x1100); ok {
		t.Errorf("FuncAt(0x1100) = %q, want none (past a's end)", s.Name)
	}
	if s, ok := img.FuncAt(0x2108); !ok || s.Name != "b" {
		t.Errorf("FuncAt(0x2108) = (%q, %v), want b", s.Name, ok)
	}
	// A data symbol between must not turn the range check into a match.
	if s, ok := img.FuncAt(0x2008); ok {
		t.Errorf("FuncAt(0x2008) = %q, want none (inside a data symbol)", s.Name)
	}
}

func TestResolveAllFallsBackToSymtab(t *testing.T) {
	img := &ElfImage{LoadEnd: 0x40000, syms: []elf.Symbol{
		{Value: 0x3000, Name: "f", Info: byte(elf.STT_FUNC), Size: 0x100},
		{Value: 0x5600, Name: "pdgo$pack.55"},
	}}
	// nil Symbolizer exercises the fallback path.
	res := resolveAll(img, nil, []uint64{0x5600, 0x5700, 0x3080, 0x5600 + (1 << 20) + 5})
	if f := res[0x5600]; len(f) != 1 || f[0].Name != "pdgo$pack.55" || f[0].Approx {
		t.Errorf("0x5600 = %+v, want solid exact data symbol", f)
	}
	if f := res[0x5700]; len(f) != 1 || f[0].Name != "pdgo$pack.55+0x100" || !f[0].Approx {
		t.Errorf("0x5700 = %+v, want approximate nearest-symbol guess", f)
	}
	if f := res[0x3080]; len(f) != 1 || f[0].Name != "f" || f[0].Approx {
		t.Errorf("0x3080 = %+v, want solid frame inside function f", f)
	}
	if _, ok := res[0x5600+(1<<20)+5]; ok {
		t.Error("offset with no nearby symbol should be absent")
	}
}

// TestAddr2lineExec runs the exec path against a fake addr2line script,
// so CI does not need the ARM toolchain. Skipped on Windows (no shell
// scripts).
func TestAddr2lineExec(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake addr2line is a shell script")
	}
	dir := t.TempDir()
	fake := filepath.Join(dir, "fake-addr2line")
	script := "#!/bin/sh\n" +
		"cat >/dev/null\n" +
		"printf '0x10\\nfoo\\nbar.go:3\\n0x20\\n??\\n??:0\\n'\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	a := &Addr2line{Bin: fake, ELF: "whatever.elf"}
	res := a.Resolve([]uint64{0x10, 0x20})
	if res == nil {
		t.Fatal("Resolve returned nil")
	}
	if len(res) != 1 {
		t.Fatalf("resolved %d addresses, want 1 (?? drops out): %v", len(res), res)
	}
	if f := res[0x10]; len(f) != 1 || f[0].Name != "foo" || f[0].Line != 3 {
		t.Errorf("0x10 = %+v", f)
	}
}
