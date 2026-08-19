package main

import (
	"os"
	"testing"
)

func TestParseCrashesDeviceLog(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"testdata fixture", "testdata/crashlog_device.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			crashes := ParseCrashes(string(data))
			if len(crashes) != 2 {
				t.Fatalf("got %d crashes, want 2", len(crashes))
			}

			first := crashes[0]
			if first.When != "2026/08/18 08:56:48" {
				t.Errorf("When = %q, want 2026/08/18 08:56:48", first.When)
			}
			if first.Build != "415038e2-3.0.5-release.202175-gitlab-runner" {
				t.Errorf("Build = %q", first.Build)
			}
			wantRegs := map[string]uint32{
				"r0": 0xe000ed04, "r1": 0x9000361c, "r2": 0x00000001, "r3": 0x20009b8c,
				"r12": 0x0016142b, "lr": 0x90001f73, "pc": 0x90001ae8, "psr": 0x210f0000,
				"cfsr": 0x00008200, "hfsr": 0x40000000, "mmfar": 0xe000ed04, "bfar": 0xe000ed04,
				"rcccsr": 0,
			}
			for name, want := range wantRegs {
				if got := first.Regs[name]; got != want {
					t.Errorf("first.Regs[%s] = %#x, want %#x", name, got, want)
				}
			}
			if first.HeapAlloc != 146144 {
				t.Errorf("HeapAlloc = %d, want 146144", first.HeapAlloc)
			}

			second := crashes[1]
			if second.When != "2026/08/18 17:49:51" {
				t.Errorf("When = %q, want 2026/08/18 17:49:51", second.When)
			}
			// Compact one-line dump: unanchored fallback must recover all registers.
			if second.Regs["pc"] != 0x90005600 || second.Regs["lr"] != 0x90001afd || second.Regs["cfsr"] != 0x00020000 {
				t.Errorf("second crash regs incomplete: %v", second.Regs)
			}
			if len(second.Regs) < 12 {
				t.Errorf("second crash parsed only %d regs, want >= 12", len(second.Regs))
			}
			if second.HeapAlloc != 237472 {
				t.Errorf("HeapAlloc = %d, want 237472", second.HeapAlloc)
			}
		})
	}
}

func TestParseCrashesTable(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		check     func(t *testing.T, c *Crash)
	}{
		{
			name:      "no crash content",
			input:     "hello world\nnothing to see\n",
			wantCount: 0,
		},
		{
			name:      "partial paste without header",
			input:     " pc:90005600 lr:90001afd cfsr:00020000 \n",
			wantCount: 1,
			check: func(t *testing.T, c *Crash) {
				if c.When != "" {
					t.Errorf("When = %q, want empty", c.When)
				}
				if c.Regs["pc"] != 0x90005600 {
					t.Errorf("pc = %#x", c.Regs["pc"])
				}
			},
		},
		{
			name:      "registers on the header line (one-line paste)",
			input:     "crash at 2026/08/18 17:49:51 pc:90005600 lr:90001afd r1:90005600 r2:90005600 psr:600f0000 cfsr:00020000\n",
			wantCount: 1,
			check: func(t *testing.T, c *Crash) {
				if c.When != "2026/08/18 17:49:51" {
					t.Errorf("When = %q, want bare timestamp", c.When)
				}
				if c.Regs["pc"] != 0x90005600 || c.Regs["lr"] != 0x90001afd || c.Regs["cfsr"] != 0x00020000 {
					t.Errorf("regs = %v", c.Regs)
				}
			},
		},
		{
			name: "register-looking text on non-register lines is ignored",
			input: "--- crash at 2026/08/18 08:56:48---\n" +
				"Lua totalbytes=0 GCdebt=0 GCestimate=0 stacksize=0\n",
			wantCount: 0,
		},
		{
			name: "too few registers is not a crash",
			input: "--- crash at 2026/08/18 08:56:48---\n" +
				"   r0:e000ed04\n",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crashes := ParseCrashes(tt.input)
			if len(crashes) != tt.wantCount {
				t.Fatalf("got %d crashes, want %d", len(crashes), tt.wantCount)
			}
			if tt.check != nil {
				tt.check(t, crashes[0])
			}
		})
	}
}
