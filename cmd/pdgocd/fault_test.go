package main

import (
	"strings"
	"testing"
)

// The two cfsr values seen on real devices during the GC corruption hunt.
func TestDecodeFaultRealLogs(t *testing.T) {
	tests := []struct {
		name       string
		regs       map[string]uint32
		wantSubstr []string
	}{
		{
			name: "precise bus fault into system control space",
			regs: map[string]uint32{
				"cfsr": 0x00008200, "hfsr": 0x40000000,
				"mmfar": 0xe000ed04, "bfar": 0xe000ed04,
			},
			wantSubstr: []string{
				"bus fault",
				"precise data bus error",
				"bfar=0xe000ed04",
				"System Control Space",
				"HardFault",
			},
		},
		{
			name: "invstate from blx to even address",
			regs: map[string]uint32{
				"cfsr": 0x00020000, "hfsr": 0x40000000,
			},
			wantSubstr: []string{
				"usage fault",
				"invalid state",
				"non-Thumb",
				"HardFault",
			},
		},
		{
			name: "undefined instruction",
			regs: map[string]uint32{
				"cfsr": 0x00010000, "hfsr": 0x40000000,
			},
			wantSubstr: []string{
				"usage fault",
				"undefined instruction",
			},
		},
		{
			name: "clean state decodes to nothing",
			regs: map[string]uint32{
				"cfsr": 0, "hfsr": 0,
			},
			wantSubstr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crash{Regs: tt.regs}
			lines := decodeFault(c)
			joined := strings.Join(lines, "\n")
			for _, want := range tt.wantSubstr {
				if !strings.Contains(joined, want) {
					t.Errorf("decodeFault output missing %q:\n%s", want, joined)
				}
			}
			if tt.wantSubstr == nil && len(lines) != 0 {
				t.Errorf("want no lines, got:\n%s", joined)
			}
		})
	}
}

func TestDecodePSR(t *testing.T) {
	tests := []struct {
		psr  uint32
		want string
	}{
		{0x210f0000, "Thumb, thread mode"},
		{0x600f0000, "ARM state (T-bit clear — attempted non-Thumb execution), thread mode"},
		{0x00000003, "ARM state (T-bit clear — attempted non-Thumb execution), exception 3 (HardFault)"},
	}
	for _, tt := range tests {
		if got := decodePSR(tt.psr); got != tt.want {
			t.Errorf("decodePSR(%#x) = %q, want %q", tt.psr, got, tt.want)
		}
	}
}
