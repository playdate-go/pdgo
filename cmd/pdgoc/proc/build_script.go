package proc

// rawBuildScript is the embedded build script for Playdate device builds.
// It expects the following environment variables:
//   - GAME_NAME: Project name
//   - GAME_DIR: Project root directory
//   - GO_SRC_DIR: Directory containing Go source files (go.mod)
//   - BUILD_DIR: Directory for build artifacts
//   - TINYGO_DIR: TinyGo installation directory
const rawBuildScript = `#!/bin/bash
set -e

TINYGO="$HOME/tinygo-playdate/build/tinygo"
TINYGO_DIR="$HOME/tinygo-playdate"

SDK="${PLAYDATE_SDK_PATH:-$(grep '^\s*SDKRoot' ~/.Playdate/config 2>/dev/null | head -n 1 | cut -c9-)}"
[ -z "$SDK" ] && { echo "Error: Playdate SDK not found"; exit 1; }
[ ! -x "$TINYGO" ] && { echo "Error: TinyGo not found at $TINYGO"; exit 1; }
! command -v arm-none-eabi-gcc &>/dev/null && { echo "Error: arm-none-eabi-gcc not found"; exit 1; }

# These are passed from pdgoc
GAME_DIR="${GAME_DIR:-.}"
GO_SRC_DIR="${GO_SRC_DIR:-$GAME_DIR/Source}"
BUILD_DIR="${BUILD_DIR:-$GAME_DIR/build}"

GAME_NAME=${GAME_NAME}

echo "🎮 Building $GAME_NAME for Playdate Device"
echo "📂 Go source: $GO_SRC_DIR"
echo "📂 Build dir: $BUILD_DIR"

cd "$GO_SRC_DIR"

# Step 1: Compile C runtime (pd_runtime.c should already exist in BUILD_DIR)
echo "📦 Step 1: Creating runtime library..."
if [ ! -f "$BUILD_DIR/pd_runtime.c" ]; then
    echo "Error: pd_runtime.c not found in $BUILD_DIR"
    exit 1
fi

MCFLAGS="-mthumb -mcpu=cortex-m7 -mfloat-abi=hard -mfpu=fpv5-sp-d16"
arm-none-eabi-gcc $MCFLAGS -O2 -DTARGET_PLAYDATE=1 -c "$BUILD_DIR/pd_runtime.c" -o "$BUILD_DIR/pd_runtime.o"
arm-none-eabi-ar rcs "$BUILD_DIR/libpd.a" "$BUILD_DIR/pd_runtime.o"

# Step 2: Create linker script and TinyGo target
echo "⚙️  Step 2: Configuring build..."
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
    "default-stack-size": 61800,
    "linkerscript": "$BUILD_DIR/playdate.ld",
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"],
    "ldflags": ["-L$BUILD_DIR", "-lpd"]
}
EOF

# Step 3: Build with TinyGo
echo "🔨 Step 3: Compiling with TinyGo..."
$TINYGO build -target=playdate -o "$BUILD_DIR/pdex.elf" .

cp "$BUILD_DIR/pdex.elf" "$GO_SRC_DIR/"

echo "📊 ELF info:"
arm-none-eabi-size "$BUILD_DIR/pdex.elf"

# Step 4: Create PDX
echo "📦 Step 4: Creating .pdx bundle..."
"$SDK/bin/pdc" "-k" "$GO_SRC_DIR" "$GAME_DIR/${GAME_NAME}_device.pdx"

# Clean up: remove pdex.elf from Source directory and build directory
rm -f "$GO_SRC_DIR/pdex.elf"
rm -rf "$BUILD_DIR"

echo ""
echo "✅ Build complete: $GAME_DIR/${GAME_NAME}_device.pdx"
`
