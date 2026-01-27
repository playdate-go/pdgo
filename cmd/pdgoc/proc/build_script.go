package proc

// rawBuildScript is the embedded build script for Playdate device builds.
// It expects the following environment variables:
//   - GAME_NAME: Project name
//   - GAME_DIR: Project root directory
//   - GO_SRC_DIR: Directory containing Go source files (go.mod)
//   - BUILD_DIR: Directory for build artifacts
const rawBuildScript = `#!/bin/bash
set -e

# TinyGo with Playdate support
TINYGO="${TINYGO_PLAYDATE:-$HOME/tinygo-playdate/build/tinygo}"
TINYGO_DIR="${TINYGO_PLAYDATE_DIR:-$HOME/tinygo-playdate}"

# Playdate SDK path
SDK="${PLAYDATE_SDK_PATH:-$(grep '^\s*SDKRoot' ~/.Playdate/config 2>/dev/null | head -n 1 | cut -c9-)}"

# Validate requirements
[ -z "$SDK" ] && { echo "Error: Playdate SDK not found"; exit 1; }
[ ! -x "$TINYGO" ] && { echo "Error: TinyGo not found at $TINYGO"; exit 1; }
! command -v arm-none-eabi-gcc &>/dev/null && { echo "Error: arm-none-eabi-gcc not found"; exit 1; }

# These are passed from pdgoc
GAME_DIR="${GAME_DIR:-.}"
GO_SRC_DIR="${GO_SRC_DIR:-$GAME_DIR/Source}"
BUILD_DIR="${BUILD_DIR:-$GAME_DIR/build}"
GAME_NAME="${GAME_NAME:-game}"

echo "=== Building $GAME_NAME for Playdate ==="
echo "Go source: $GO_SRC_DIR"
echo "Build dir: $BUILD_DIR"

# ARM compiler flags
MCFLAGS="-mthumb -mcpu=cortex-m7 -mfloat-abi=hard -mfpu=fpv5-sp-d16"
CPFLAGS="-ffunction-sections -fdata-sections -mword-relocations -fno-common"
LDFLAGS="--gc-sections --emit-relocs"

cd "$GO_SRC_DIR"

# Step 1: Compile C runtime (pd_runtime.c)
echo "Step 1: Compiling C runtime..."
arm-none-eabi-gcc $MCFLAGS $CPFLAGS -O2 -DTARGET_PLAYDATE=1 -DTARGET_EXTENSION=1 \
    -I"$SDK/C_API" \
    -c "$BUILD_DIR/pd_runtime.c" -o "$BUILD_DIR/pd_runtime.o"

# Create static library from pd_runtime.o
arm-none-eabi-ar rcs "$BUILD_DIR/libpd.a" "$BUILD_DIR/pd_runtime.o"

# Step 2: Create linker script and TinyGo target
echo "Step 2: Configuring build..."
cat > "$BUILD_DIR/playdate.ld" << 'LDSCRIPT'
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
LDSCRIPT

mkdir -p "$TINYGO_DIR/targets"
cat > "$TINYGO_DIR/targets/playdate.json" << EOF
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
    "linkerscript": "$BUILD_DIR/playdate.ld",
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"],
    "ldflags": ["-L$BUILD_DIR", "-lpd"]
}
EOF

# Step 3: Build with TinyGo
echo "Step 3: Compiling Go code with TinyGo..."
$TINYGO build -target=playdate -o "$BUILD_DIR/game.o" .

# Step 4: Compile setup.c from SDK
echo "Step 4: Compiling SDK setup..."
arm-none-eabi-gcc $MCFLAGS $CPFLAGS -O2 -DTARGET_PLAYDATE=1 -DTARGET_EXTENSION=1 \
    -I"$SDK/C_API" \
    -c "$SDK/C_API/buildsupport/setup.c" -o "$BUILD_DIR/setup.o"

# Step 5: Link everything
echo "Step 5: Linking..."
arm-none-eabi-gcc $MCFLAGS -T"$SDK/C_API/buildsupport/link_map.ld" \
    -Wl,--gc-sections \
    -Wl,--emit-relocs \
    -nostartfiles \
    "$BUILD_DIR/setup.o" \
    "$BUILD_DIR/pd_runtime.o" \
    "$BUILD_DIR/game.o" \
    -o "$BUILD_DIR/pdex.elf"

# Copy ELF to Source for pdc
cp "$BUILD_DIR/pdex.elf" "$GO_SRC_DIR/"

echo "ELF size:"
arm-none-eabi-size "$BUILD_DIR/pdex.elf"

# Step 6: Create PDX bundle
echo "Step 6: Creating .pdx bundle..."
"$SDK/bin/pdc" -k "$GO_SRC_DIR" "$GAME_DIR/${GAME_NAME}.pdx"

# Cleanup
rm -f "$GO_SRC_DIR/pdex.elf"

echo ""
echo "=== Build complete: $GAME_DIR/${GAME_NAME}.pdx ==="
`
