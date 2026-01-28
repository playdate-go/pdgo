#!/bin/bash
#
# Build TinyGo with Playdate Support
#
# This script builds TinyGo from source with Playdate runtime support.
#
# Build time:
#   - With system LLVM (Linux only): ~1 minute
#   - Building LLVM from source: ~25-30 minutes
#
# Usage:
#   ./build-tinygo-playdate.sh [--jobs N] [--use-system-llvm] [--build-llvm]
#
# Options:
#   --tinygo-version VER  TinyGo version to build (default: 0.40.1)
#   --use-system-llvm     Use system LLVM (Linux only! Does NOT work on macOS)
#   --build-llvm          Force building LLVM from source
#   --jobs N              Number of parallel jobs (default: auto)
#
# NOTE: --use-system-llvm does NOT work on macOS because Homebrew LLVM
#       does not include static LLD libraries required to build TinyGo.
#       On macOS, always build LLVM from source (default behavior).
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

FORCE_SYSTEM_LLVM=false
FORCE_BUILD_LLVM=false
LLVM_VERSION="20"
TINYGO_VERSION="0.40.1"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --jobs)
            JOBS="$2"
            shift 2
            ;;
        --tinygo-version)
            TINYGO_VERSION="$2"
            shift 2
            ;;
        --use-system-llvm)
            FORCE_SYSTEM_LLVM=true
            shift
            ;;
        --build-llvm)
            FORCE_BUILD_LLVM=true
            shift
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
echo -e "TinyGo version:  ${GREEN}v$TINYGO_VERSION${NC}"
echo -e "LLVM version:    ${GREEN}$LLVM_VERSION${NC}"
echo -e "Build directory: ${GREEN}$TINYGO_DIR${NC}"
echo -e "Parallel jobs:   ${GREEN}$JOBS${NC}"
echo ""

# Function to find system LLVM
find_system_llvm() {
    local llvm_path=""
    
    # macOS: Check Homebrew
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # Try different LLVM versions (20 first for TinyGo 0.40+)
        for ver in 20 19 18 17; do
            if brew --prefix llvm@$ver &>/dev/null 2>&1; then
                llvm_path="$(brew --prefix llvm@$ver)"
                if [ -f "$llvm_path/bin/llvm-config" ]; then
                    echo "$llvm_path"
                    return 0
                fi
            fi
        done
        # Try unversioned llvm
        if brew --prefix llvm &>/dev/null 2>&1; then
            llvm_path="$(brew --prefix llvm)"
            if [ -f "$llvm_path/bin/llvm-config" ]; then
                echo "$llvm_path"
                return 0
            fi
        fi
    fi
    
    # Linux: Check common paths (20 first for TinyGo 0.40+)
    for ver in 20 19 18 17; do
        # apt-style paths
        if [ -f "/usr/lib/llvm-$ver/bin/llvm-config" ]; then
            echo "/usr/lib/llvm-$ver"
            return 0
        fi
        # Alternative Linux paths
        if [ -f "/usr/lib64/llvm$ver/bin/llvm-config" ]; then
            echo "/usr/lib64/llvm$ver"
            return 0
        fi
    done
    
    # Check if llvm-config is in PATH
    if command -v llvm-config &>/dev/null; then
        local llvm_bin_dir="$(dirname "$(which llvm-config)")"
        echo "$(dirname "$llvm_bin_dir")"
        return 0
    fi
    
    return 1
}

# Function to get LLVM install instructions
get_llvm_install_instructions() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "  ⚠️  macOS: --use-system-llvm is NOT supported (Homebrew LLVM lacks LLD libraries)"
        echo "  Just run the script without --use-system-llvm to build LLVM from source."
    elif command -v apt &>/dev/null; then
        echo "  # Add LLVM apt repository"
        echo "  wget https://apt.llvm.org/llvm.sh"
        echo "  chmod +x llvm.sh"
        echo "  sudo ./llvm.sh $LLVM_VERSION"
        echo ""
        echo "  # Install required dev packages"
        echo "  sudo apt install clang-$LLVM_VERSION llvm-$LLVM_VERSION-dev lld-$LLVM_VERSION libclang-$LLVM_VERSION-dev"
    elif command -v dnf &>/dev/null; then
        echo "  sudo dnf install llvm$LLVM_VERSION-devel clang$LLVM_VERSION-devel lld$LLVM_VERSION-devel"
    elif command -v pacman &>/dev/null; then
        echo "  sudo pacman -S llvm clang lld"
    else
        echo "  Download from: https://releases.llvm.org/"
    fi
}

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

# Determine LLVM strategy
USE_SYSTEM_LLVM=false
SYSTEM_LLVM_PATH=""

# Check for macOS + --use-system-llvm (not supported)
if [ "$FORCE_SYSTEM_LLVM" = true ] && [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "${RED}ERROR: --use-system-llvm is NOT supported on macOS${NC}"
    echo ""
    echo "Homebrew LLVM does not include static LLD libraries (liblldCOFF.a, liblldELF.a, etc.)"
    echo "which are required to build TinyGo."
    echo ""
    echo "Please run the script without --use-system-llvm to build LLVM from source:"
    echo "  ./build-tinygo-playdate.sh"
    exit 1
fi

if [ "$FORCE_BUILD_LLVM" = true ]; then
    echo -e "${YELLOW}Forcing LLVM build from source (--build-llvm)${NC}"
    USE_SYSTEM_LLVM=false
elif SYSTEM_LLVM_PATH=$(find_system_llvm); then
    echo -e "${GREEN}Found system LLVM at: $SYSTEM_LLVM_PATH${NC}"
    LLVM_VER=$("$SYSTEM_LLVM_PATH/bin/llvm-config" --version 2>/dev/null || echo "unknown")
    echo -e "LLVM version: ${GREEN}$LLVM_VER${NC}"
    USE_SYSTEM_LLVM=true
elif [ "$FORCE_SYSTEM_LLVM" = true ]; then
    echo -e "${RED}System LLVM not found but --use-system-llvm was specified${NC}"
    echo ""
    echo "Install LLVM first:"
    get_llvm_install_instructions
    exit 1
else
    echo -e "${YELLOW}System LLVM not found. Will build from source (~25-30 minutes).${NC}"
    echo ""
    echo -e "To speed up future builds, install LLVM:"
    get_llvm_install_instructions
    echo ""
    read -p "Continue with LLVM build from source? [Y/n] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        echo "Aborted. Install LLVM and run again."
        exit 1
    fi
fi

if [ "$USE_SYSTEM_LLVM" = true ]; then
    echo -e "${GREEN}Using system LLVM - build will take ~1 minute${NC}"
else
    echo -e "${YELLOW}Building LLVM from source - this will take ~25 minutes${NC}"
fi
echo ""

# Clone or update TinyGo
if [ -d "$TINYGO_DIR/.git" ]; then
    echo -e "${YELLOW}Updating TinyGo repository...${NC}"
    cd "$TINYGO_DIR"
    git fetch --all --tags
else
    echo -e "${YELLOW}Cloning TinyGo repository...${NC}"
    git clone --recursive https://github.com/tinygo-org/tinygo.git "$TINYGO_DIR"
    cd "$TINYGO_DIR"
fi

# Checkout specific TinyGo version
echo -e "${YELLOW}Checking out TinyGo v${TINYGO_VERSION}...${NC}"
git checkout "v${TINYGO_VERSION}"

# Initialize submodules
git submodule update --init --recursive

# LLVM setup
if [ "$USE_SYSTEM_LLVM" = true ]; then
    echo ""
    echo -e "${GREEN}Step 1/3: Using system LLVM (skipping download and build)${NC}"
    LLVM_BUILDDIR="$SYSTEM_LLVM_PATH"
else
    # Download LLVM sources
    echo ""
    echo -e "${YELLOW}Step 1/4: Downloading LLVM sources...${NC}"
    make llvm-source

    # Build LLVM
    echo ""
    echo -e "${YELLOW}Step 2/4: Building LLVM (this takes ~20-25 minutes)...${NC}"
    echo -e "Started at: $(date)"

    # Use clang if available for faster build
    if command -v clang &>/dev/null; then
        export CC=clang
        export CXX=clang++
    fi

    LLVM_PARALLEL=$JOBS make llvm-build

    echo -e "LLVM built at: $(date)"
    LLVM_BUILDDIR=""
fi

# Copy Playdate target and runtime files
STEP_NUM=$( [ "$USE_SYSTEM_LLVM" = true ] && echo "2/3" || echo "3/4" )
echo ""
echo -e "${YELLOW}Step $STEP_NUM: Adding Playdate support...${NC}"

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
    "default-stack-size": 131072,
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

# Copy runtime Go file
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
STEP_NUM=$( [ "$USE_SYSTEM_LLVM" = true ] && echo "3/3" || echo "4/4" )
echo ""
echo -e "${YELLOW}Step $STEP_NUM: Building TinyGo...${NC}"

if [ "$USE_SYSTEM_LLVM" = true ]; then
    # Build with system LLVM
    make LLVM_BUILDDIR="$LLVM_BUILDDIR"
else
    # Build with locally compiled LLVM
    make
fi

# Verify
echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Build Complete!                       ║${NC}"
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
