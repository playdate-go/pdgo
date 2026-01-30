#!/bin/bash
#
# pdgo Full Installer
#
# Installs everything needed to build Go games for Playdate:
#   - pdgoc (build tool)
#   - TinyGo with Playdate support (for device builds)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/playdate-go/pdgo/main/install.sh | bash
#
# Options (via environment variables):
#   SKIP_TINYGO=1  - Skip TinyGo build (simulator-only)
#   JOBS=N         - Number of parallel jobs for LLVM build
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

TINYGO_DIR="$HOME/tinygo-playdate"
TINYGO_VERSION="0.40.1"
LLVM_VERSION="20"
JOBS=${JOBS:-$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)}

echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              PdGo - Full Installer                       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check OS
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo -e "${RED}ERROR: Windows is not supported${NC}"
    echo "Supported: macOS, Linux"
    exit 1
fi

# ============================================================================
# Step 1: Check dependencies
# ============================================================================
echo -e "${YELLOW}[1/4] Checking dependencies...${NC}"

missing=""
command -v go &>/dev/null || missing="$missing go"
command -v git &>/dev/null || missing="$missing git"
command -v cmake &>/dev/null || missing="$missing cmake"
command -v ninja &>/dev/null || missing="$missing ninja"

if [ -n "$missing" ]; then
    echo -e "${RED}Missing dependencies:$missing${NC}"
    echo ""
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Install with: brew install$missing"
    elif command -v apt &>/dev/null; then
        echo "Install with: sudo apt install$missing"
    elif command -v dnf &>/dev/null; then
        echo "Install with: sudo dnf install$missing"
    elif command -v pacman &>/dev/null; then
        echo "Install with: sudo pacman -S$missing"
    fi
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "  Go: ${GREEN}$GO_VERSION${NC}"

# Check Playdate SDK
if [ -z "$PLAYDATE_SDK_PATH" ]; then
    echo -e "${RED}ERROR: PLAYDATE_SDK_PATH is not set${NC}"
    echo ""
    echo "1. Download Playdate SDK from https://play.date/dev/"
    echo "2. Add to your shell profile (~/.zshrc or ~/.bashrc):"
    echo "   export PLAYDATE_SDK_PATH=\"/path/to/PlaydateSDK\""
    exit 1
fi

if [ ! -d "$PLAYDATE_SDK_PATH" ]; then
    echo -e "${RED}ERROR: PLAYDATE_SDK_PATH does not exist: $PLAYDATE_SDK_PATH${NC}"
    exit 1
fi

SDK_VERSION=$(cat "$PLAYDATE_SDK_PATH/VERSION.txt" 2>/dev/null | tr -d '\n' || echo "unknown")
echo -e "  Playdate SDK: ${GREEN}$SDK_VERSION${NC}"
echo -e "  cmake: ${GREEN}$(cmake --version | head -1 | awk '{print $3}')${NC}"
echo -e "  ninja: ${GREEN}$(ninja --version)${NC}"

echo -e "${GREEN}All dependencies OK${NC}"

# ============================================================================
# Step 2: Install pdgoc
# ============================================================================
echo ""
echo -e "${YELLOW}[2/4] Installing pdgoc...${NC}"

go install github.com/playdate-go/pdgo/cmd/pdgoc@latest

GOBIN="${GOBIN:-$(go env GOPATH)/bin}"
if [ -f "$GOBIN/pdgoc" ]; then
    echo -e "  pdgoc installed at: ${GREEN}$GOBIN/pdgoc${NC}"
else
    echo -e "${RED}Failed to install pdgoc${NC}"
    exit 1
fi

# Check if GOBIN is in PATH
if ! command -v pdgoc &>/dev/null; then
    echo -e "${YELLOW}  Note: Add GOBIN to PATH: export PATH=\"\$PATH:$GOBIN\"${NC}"
fi

# ============================================================================
# Step 3: Build TinyGo with Playdate support
# ============================================================================
if [ "$SKIP_TINYGO" = "1" ]; then
    echo ""
    echo -e "${YELLOW}[3/4] Skipping TinyGo build (SKIP_TINYGO=1)${NC}"
    echo "  You can build it later with:"
    echo "  git clone https://github.com/playdate-go/pdgo.git && cd pdgo"
    echo "  ./cmd/pdgoc/scripts/build-tinygo-playdate.sh"
else
    echo ""
    echo -e "${YELLOW}[3/4] Building TinyGo with Playdate support...${NC}"
    echo -e "  TinyGo version: ${GREEN}v$TINYGO_VERSION${NC}"
    echo -e "  Install path: ${GREEN}$TINYGO_DIR${NC}"
    echo -e "  Parallel jobs: ${GREEN}$JOBS${NC}"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo -e "  LLVM: ${GREEN}Building from source (~25 minutes)${NC}"
    else
        echo -e "  LLVM: ${GREEN}Will use system LLVM if available${NC}"
    fi
    echo ""

    # Clone or update TinyGo
    if [ -d "$TINYGO_DIR/.git" ]; then
        echo "  Updating TinyGo repository..."
        cd "$TINYGO_DIR"
        git fetch --all --tags
    else
        echo "  Cloning TinyGo repository..."
        git clone --recursive https://github.com/tinygo-org/tinygo.git "$TINYGO_DIR"
        cd "$TINYGO_DIR"
    fi

    # Checkout specific version
    echo "  Checking out v$TINYGO_VERSION..."
    git checkout "v$TINYGO_VERSION"
    git submodule update --init --recursive

    # Clean previous build to ensure Playdate files are included
    if [ -d "$TINYGO_DIR/build" ]; then
        echo "  Cleaning previous build..."
        make clean 2>/dev/null || rm -rf "$TINYGO_DIR/build"
    fi

    # Add Playdate support files BEFORE building (required for gc.playdate)
    echo "  Adding Playdate support files..."

    # playdate.json
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

    # playdate.ld
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

    # runtime_playdate.go
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

// runtime_init is called from C to initialize the Go runtime
//export runtime_init
func runtime_init() {
    initAll()
}

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

    # gc_playdate.go
    cat > "$TINYGO_DIR/src/runtime/gc_playdate.go" << 'EOF'
//go:build gc.playdate

package runtime

import (
    "unsafe"
)

const needsStaticHeap = false

var gcTotalAlloc uint64
var gcMallocs uint64
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

func markRoots(start, end uintptr) {}

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

func GC() {}

func SetFinalizer(obj interface{}, finalizer interface{}) {}

func initHeap() {}

func setHeapEnd(newHeapEnd uintptr) {}
EOF

    # Check for system LLVM (Linux only)
    USE_SYSTEM_LLVM=false
    SYSTEM_LLVM_PATH=""

    if [[ "$OSTYPE" != "darwin"* ]]; then
        # Try to find system LLVM on Linux
        for ver in 20 19 18 17; do
            if [ -f "/usr/lib/llvm-$ver/bin/llvm-config" ]; then
                SYSTEM_LLVM_PATH="/usr/lib/llvm-$ver"
                USE_SYSTEM_LLVM=true
                break
            fi
        done
    fi

    # Build TinyGo
    if [ "$USE_SYSTEM_LLVM" = true ]; then
        echo -e "  Found system LLVM at: ${GREEN}$SYSTEM_LLVM_PATH${NC}"
        echo "  Building TinyGo (~1 minute)..."
        make LLVM_BUILDDIR="$SYSTEM_LLVM_PATH"
    else
        echo "  Downloading LLVM sources..."
        make llvm-source

        echo "  Building LLVM (~20-25 minutes)..."
        echo "  Started at: $(date)"
        
        # Use clang if available
        if command -v clang &>/dev/null; then
            export CC=clang
            export CXX=clang++
        fi

        LLVM_PARALLEL=$JOBS make llvm-build
        echo "  LLVM completed at: $(date)"

        echo "  Building TinyGo..."
        make
    fi

    echo -e "  ${GREEN}TinyGo built successfully!${NC}"
fi

# ============================================================================
# Step 4: Setup PATH
# ============================================================================
echo ""
echo -e "${YELLOW}[4/4] Configuring PATH...${NC}"

SHELL_RC=""
if [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
fi

NEED_GOBIN=false
NEED_TINYGO=false

if ! echo "$PATH" | grep -q "$GOBIN"; then
    NEED_GOBIN=true
fi

if [ "$SKIP_TINYGO" != "1" ] && ! echo "$PATH" | grep -q "$TINYGO_DIR/build"; then
    NEED_TINYGO=true
fi

if [ "$NEED_GOBIN" = true ] || [ "$NEED_TINYGO" = true ]; then
    echo ""
    echo -e "${YELLOW}Add these lines to $SHELL_RC:${NC}"
    echo ""
    if [ "$NEED_GOBIN" = true ]; then
        echo "  export PATH=\"\$PATH:$GOBIN\""
    fi
    if [ "$NEED_TINYGO" = true ]; then
        echo "  export PATH=\"\$PATH:$TINYGO_DIR/build\""
    fi
    echo ""
    
    read -p "Add automatically? [Y/n] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        echo "" >> "$SHELL_RC"
        echo "# pdgo - Playdate Go toolkit" >> "$SHELL_RC"
        if [ "$NEED_GOBIN" = true ]; then
            echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
        fi
        if [ "$NEED_TINYGO" = true ]; then
            echo "export PATH=\"\$PATH:$TINYGO_DIR/build\"" >> "$SHELL_RC"
        fi
        echo -e "${GREEN}Added to $SHELL_RC${NC}"
        echo "Run: source $SHELL_RC"
    fi
else
    echo -e "${GREEN}PATH already configured${NC}"
fi

# ============================================================================
# Done!
# ============================================================================
echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Installation Complete!                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Installed:"
echo -e "  pdgoc:  ${GREEN}$GOBIN/pdgoc${NC}"
if [ "$SKIP_TINYGO" != "1" ]; then
echo -e "  TinyGo: ${GREEN}$TINYGO_DIR/build/tinygo${NC}"
fi
echo ""
echo "Usage:"
echo ""
echo "  # Simulator build (run in simulator)"
echo "  pdgoc -sim -run -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
if [ "$SKIP_TINYGO" != "1" ]; then
echo "  # Device build (for real Playdate)"
echo "  pdgoc -device -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
echo "  # Device build + deploy to connected Playdate"
echo "  pdgoc -device -deploy -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
fi
echo "Documentation: https://github.com/playdate-go/pdgo"
echo ""
