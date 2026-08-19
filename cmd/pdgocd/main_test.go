package main

import (
	"debug/elf"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A .pdx bundle carries only the pdc-encrypted pdex.bin; the symbolizable
// ELF sits in the game dir beside the bundle (build/pdex.elf from manual
// DeviceBuildScriptUnix.sh runs). Passing the bundle must find it there.
func TestFindElfFromPdxBundle(t *testing.T) {
	gameDir := filepath.Join(t.TempDir(), "game")
	bundle := filepath.Join(gameDir, "MyGame.pdx")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(gameDir, "build"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "pdex.bin"), []byte("Playdate PDX\x00\x00\x00\x00junk"), 0o644); err != nil {
		t.Fatal(err)
	}
	elfPath := filepath.Join(gameDir, "build", "pdex.elf")
	if err := os.WriteFile(elfPath, []byte("\x7fELF-fake"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := findElf(bundle, "")
	if err != nil || got != elfPath {
		t.Fatalf("findElf(bundle) = (%q, %v), want (%q, nil)", got, err, elfPath)
	}

	// Without the sibling ELF the encrypted-bin diagnostic must win
	// over a silent cwd walk-up.
	if err := os.RemoveAll(filepath.Join(gameDir, "build")); err != nil {
		t.Fatal(err)
	}
	if _, err := findElf(bundle, ""); err == nil || !strings.Contains(err.Error(), "pdc-encrypted") {
		t.Fatalf("want encrypted-bin error, got %v", err)
	}
}

func TestLogFlag(t *testing.T) {
	var out, errw strings.Builder

	// Raw -log text reaches the parser (fails there for lack of a
	// crash record — proving the input came from the flag).
	errw.Reset()
	if code := run([]string{"-log", "no crash in here"}, strings.NewReader(""), &out, &errw); code != 2 || !strings.Contains(errw.String(), "no crash records found") {
		t.Errorf("code=%d stderr=%q", code, errw.String())
	}

	// -log together with a log file is rejected.
	errw.Reset()
	if code := run([]string{"-log", "x", "somelog.txt"}, strings.NewReader(""), &out, &errw); code != 2 || !strings.Contains(errw.String(), "mutually exclusive") {
		t.Errorf("code=%d stderr=%q", code, errw.String())
	}
}

func TestSymbolizationWarning(t *testing.T) {
	img := &ElfImage{LoadEnd: 0x40000, syms: []elf.Symbol{
		{Value: 0x1afd, Name: "runtime.runFinalizerSafe", Info: byte(elf.STT_FUNC)},
	}}
	crash := &Crash{Regs: map[string]uint32{
		"pc": 0x90001afd, "lr": 0x90002000,
	}}

	tests := []struct {
		name   string
		frames map[uint64][]Frame
		want   string
	}{
		{
			name:   "pc and lr resolve - no warning",
			frames: map[uint64][]Frame{0x1afd: {{Name: "runtime.runFinalizerSafe"}}, 0x2000: {{Name: "runtime.GC"}}},
			want:   "",
		},
		{
			name:   "lr resolves - no warning",
			frames: map[uint64][]Frame{0x2000: {{Name: "runtime.GC"}}},
			want:   "",
		},
		{
			name: "nearest-symbol guesses only - different ELF",
			frames: map[uint64][]Frame{
				0x1afd: {{Name: "__bss_start__+0xdc", Approx: true}},
				0x2000: {{Name: "__heap_start__", Approx: true}},
			},
			want: "cannot resolve the crash addresses to Go function names",
		},
		{
			name:   "one solid anchor outweighs a guessed pc",
			frames: map[uint64][]Frame{0x1afd: {{Name: "__bss_start__+0xdc", Approx: true}}, 0x2000: {{Name: "runtime.GC"}}},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := symbolizationWarning(crash, img, tt.frames)
			if tt.want == "" && got != "" {
				t.Errorf("unexpected warning: %q", got)
			}
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("warning = %q, want it to contain %q", got, tt.want)
			}
		})
	}

	// pc in the flash window but past the ELF's loaded image: the ELF
	// cannot be the crashing build.
	pastEnd := &Crash{Regs: map[string]uint32{"pc": 0x90041000, "lr": 0x90001afd}}
	if w := symbolizationWarning(pastEnd, img, nil); !strings.Contains(w, "past the end of this ELF") {
		t.Errorf("want past-image warning, got %q", w)
	}

	// pc/lr outside the flash window (e.g. in system space): nothing to
	// judge symbolization by, so no warning.
	outside := &Crash{Regs: map[string]uint32{"pc": 0xe000ed04, "lr": 0x20000000}}
	if w := symbolizationWarning(outside, img, nil); w != "" {
		t.Errorf("unexpected warning for out-of-window pc/lr: %q", w)
	}
}
