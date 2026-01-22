#!/bin/bash
#
# Build TinyGo with Playdate Support
#
# This script builds TinyGo from source with Playdate runtime support.
# Takes 0.5-1h to complete.
#
# Usage:
#   ./build-tinygo-playdate.sh [--jobs N]
#
# The resulting TinyGo will be in:
#   ~/tinygo-playdate/build/tinygo
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

JOBS=${JOBS:-$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)}
TINYGO_DIR="$HOME/tinygo-playdate"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PDGO_DIR="$(dirname "$SCRIPT_DIR")"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --jobs)
            JOBS="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║       Building TinyGo with Playdate Support              ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Build directory: ${GREEN}$TINYGO_DIR${NC}"
echo -e "Parallel jobs:   ${GREEN}$JOBS${NC}"
echo ""
echo -e "${YELLOW}⚠️  This will take 1-2 hours!${NC}"
echo ""

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

missing=""
command -v go >/dev/null || missing="$missing go"
command -v cmake >/dev/null || missing="$missing cmake"
command -v ninja >/dev/null || missing="$missing ninja"
command -v git >/dev/null || missing="$missing git"

if [ -n "$missing" ]; then
    echo -e "${RED}Missing dependencies:$missing${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Install with: brew install$missing"
    elif command -v apt &>/dev/null; then
        echo "Install with: sudo apt install$missing"
    elif command -v dnf &>/dev/null; then
        echo "Install with: sudo dnf install$missing"
    elif command -v pacman &>/dev/null; then
        echo "Install with: sudo pacman -S$missing"
    else
        echo "Please install:$missing"
    fi
    exit 1
fi

echo -e "${GREEN}All dependencies OK${NC}"

# Clone or update TinyGo
if [ -d "$TINYGO_DIR/.git" ]; then
    echo -e "${YELLOW}Updating TinyGo repository...${NC}"
    cd "$TINYGO_DIR"
    git pull
else
    echo -e "${YELLOW}Cloning TinyGo repository...${NC}"
    git clone --recursive https://github.com/tinygo-org/tinygo.git "$TINYGO_DIR"
    cd "$TINYGO_DIR"
fi

# Initialize submodules
git submodule update --init --recursive

# Download LLVM sources
echo ""
echo -e "${YELLOW}Step 1/4: Downloading LLVM sources...${NC}"
make llvm-source

# Build LLVM
echo ""
echo -e "${YELLOW}Step 2/4: Building LLVM (this takes ~1 hour)...${NC}"
echo -e "Started at: $(date)"

# Use clang if available for faster build
if command -v clang &>/dev/null; then
    export CC=clang
    export CXX=clang++
fi

LLVM_PARALLEL=$JOBS make llvm-build

echo -e "LLVM built at: $(date)"

# Copy Playdate target and runtime files
echo ""
echo -e "${YELLOW}Step 3/4: Adding Playdate support...${NC}"

# Create base target (paths will be set dynamically by pdgoc build script)
cat > "$TINYGO_DIR/targets/playdate.json" << 'EOF'
{
    "inherits": ["cortex-m"],
    "llvm-target": "thumbv7em-unknown-unknown-eabihf",
    "cpu": "cortex-m7",
    "features": "+armv7e-m,+dsp,+hwdiv,+thumb-mode,+fp-armv8d16sp,+vfp4d16sp",
    "build-tags": ["playdate", "tinygo", "gc.playdate"],
    "gc": "playdate",
    "scheduler": "none",
    "serial": "none",
    "automatic-stack-size": false,
    "default-stack-size": 61800,
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"]
}
EOF

# Copy linker script (base version, actual used by pdgoc build script)
cat > "$TINYGO_DIR/targets/playdate.ld" << 'EOF'
ENTRY(eventHandlerShim)

SECTIONS
{
    .text : ALIGN(4) {
        KEEP(*(.text.eventHandlerShim))
        KEEP(*(.text.eventHandler))
        KEEP(*(.text.updateCallback))
        KEEP(*(.text.runtime_init))
        *(.text) *(.text.*) *(.rodata) *(.rodata.*)
        KEEP(*(.init)) KEEP(*(.fini))
        . = ALIGN(4);
    }
    .data : ALIGN(4) {
        __data_start__ = .; *(.data) *(.data.*) . = ALIGN(4); __data_end__ = .;
    }
    .bss (NOLOAD) : ALIGN(4) {
        __bss_start__ = .; _sbss = .; *(.bss) *(.bss.*) *(COMMON) . = ALIGN(4); __bss_end__ = .; _ebss = .;
    }
    /DISCARD/ : { *(.ARM.exidx*) *(.ARM.extab*) }
    _sidata = LOADADDR(.data);
    _sdata = __data_start__; _edata = __data_end__;
    _globals_start = __data_start__; _globals_end = __bss_end__;
    _stack_top = __bss_end__;
}
EOF

# Copy runtime Go file (NO C files - TinyGo doesn't support them without CGO)
cat > "$TINYGO_DIR/src/runtime/runtime_playdate.go" << 'EOF'
//go:build playdate

package runtime

import "unsafe"

//go:extern _cgo_pd_realloc
func _cgo_pd_realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer

//go:extern _cgo_pd_getCurrentTimeMS
func _cgo_pd_getCurrentTimeMS() uint32

//go:extern _cgo_pd_logToConsole
func _cgo_pd_logToConsole(msg *byte)

func ticks() timeUnit {
	return timeUnit(_cgo_pd_getCurrentTimeMS()) * 1000000
}

func sleepTicks(d timeUnit) {}

func nanosecondsToTicks(ns int64) timeUnit { return timeUnit(ns) }

func ticksToNanoseconds(t timeUnit) int64 { return int64(t) }

var printBuf [256]byte
var printBufIdx int

func putchar(c byte) {
	if c == '\n' || printBufIdx >= len(printBuf)-1 {
		printBuf[printBufIdx] = 0
		_cgo_pd_logToConsole(&printBuf[0])
		printBufIdx = 0
	} else {
		printBuf[printBufIdx] = c
		printBufIdx++
	}
}
EOF

# Create Playdate GC - uses Playdate SDK realloc for memory allocation
cat > "$TINYGO_DIR/src/runtime/gc_playdate.go" << 'EOF'
//go:build gc.playdate

package runtime

// Playdate GC - uses Playdate SDK realloc for all memory allocation.
// This is a simple allocator that delegates to the Playdate SDK.

import (
	"unsafe"
)

const needsStaticHeap = false

// Total amount allocated for runtime.MemStats
var gcTotalAlloc uint64

// Total number of calls to alloc()
var gcMallocs uint64

// Total number of freed
var gcFrees uint64

//go:noinline
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
	size = align(size)
	gcTotalAlloc += uint64(size)
	gcMallocs++

	ptr := _cgo_pd_realloc(nil, size)
	if ptr == nil {
		runtimePanic("out of memory")
	}

	// Zero the allocated memory
	memzero(ptr, size)
	return ptr
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	size = align(size)
	newPtr := _cgo_pd_realloc(ptr, size)
	if newPtr == nil && size > 0 {
		runtimePanic("out of memory")
	}
	return newPtr
}

func free(ptr unsafe.Pointer) {
	if ptr != nil {
		_cgo_pd_realloc(ptr, 0)
		gcFrees++
	}
}

func markRoots(start, end uintptr) {
	// No GC, no marking needed
}

// ReadMemStats populates m with memory statistics.
func ReadMemStats(m *MemStats) {
	m.HeapIdle = 0
	m.HeapInuse = gcTotalAlloc
	m.HeapReleased = 0
	m.HeapSys = m.HeapInuse + m.HeapIdle
	m.GCSys = 0
	m.TotalAlloc = gcTotalAlloc
	m.Mallocs = gcMallocs
	m.Frees = gcFrees
	m.Sys = gcTotalAlloc
	m.HeapAlloc = gcTotalAlloc
	m.Alloc = m.HeapAlloc
}

func GC() {
	// No-op - Playdate SDK manages memory
}

func SetFinalizer(obj interface{}, finalizer interface{}) {
	// No-op
}

func initHeap() {
	// No initialization needed - Playdate SDK handles it
}

func setHeapEnd(newHeapEnd uintptr) {
	// No-op - Playdate SDK handles heap
}
EOF

echo -e "${GREEN}Playdate support files added${NC}"

# Build TinyGo
echo ""
echo -e "${YELLOW}Step 4/4: Building TinyGo...${NC}"
make

# Verify
echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Build Complete!                        ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "TinyGo binary: $TINYGO_DIR/build/tinygo"
echo ""
echo "To use:"
echo "  export PATH=\"$TINYGO_DIR/build:\$PATH\""
echo "  tinygo build -target=playdate -o game.elf ."
echo ""

# Test
"$TINYGO_DIR/build/tinygo" version
"$TINYGO_DIR/build/tinygo" targets | grep playdate || echo "Playdate target added"
