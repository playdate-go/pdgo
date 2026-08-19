#!/bin/bash
set -e

# rawBuildScript is the embedded build script for Playdate device builds.
# It expects the following environment variables:
#   - GAME_NAME: Project name
#   - GAME_DIR: Project root directory
#   - GO_SRC_DIR: Directory containing Go source files (go.mod)
#   - BUILD_DIR: Directory for build artifacts

# TinyGo with Playdate support
TINYGO="${TINYGO_PLAYDATE:-$HOME/tinygo-playdate/bin/tinygo}"
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

# Step 1b: Assemble TinyGo runtime assembly (asm_arm.S).
# TinyGo only assembles extra-files (asm_arm.S) when it links an executable
# itself; with `-o game.o` it emits just the Go module and returns early.
# tinygo_scanCurrentStack (GC stack scanning) and tinygo_longjmp (recover)
# live in asm_arm.S, so assemble it here and link it in Step 5.
# The installed asm_arm.S is the pdgo replacement (installed by install.sh):
# it restores r4-r11 and keeps SP 8-byte aligned, and is 16-bit-Thumb-only
# so GNU as can assemble it without workarounds.
arm-none-eabi-gcc $MCFLAGS -c "$TINYGO_DIR/src/runtime/asm_arm.S" -o "$BUILD_DIR/asm_arm.o"

# Step 2: Configure TinyGo target using installed patches as base
echo "Step 2: Configuring build..."

# Copy linker script from TinyGo installation (installed by install.sh from tinygo-patches/)
[ -f "$TINYGO_DIR/targets/playdate.ld" ] || { echo "Error: playdate.ld not found in $TINYGO_DIR/targets/ — run install.sh first"; exit 1; }
cp "$TINYGO_DIR/targets/playdate.ld" "$BUILD_DIR/playdate.ld"

# Extend the installed playdate.json with build-specific linkerscript and ldflags.
# A pristine copy is kept so repeated builds rewrite (not append) these keys —
# appending left duplicate JSON keys that accumulated on every build.
[ -f "$TINYGO_DIR/targets/playdate.json" ] || { echo "Error: playdate.json not found in $TINYGO_DIR/targets/ — run install.sh first"; exit 1; }
if [ ! -f "$TINYGO_DIR/targets/playdate.json.pristine" ]; then
    cp "$TINYGO_DIR/targets/playdate.json" "$TINYGO_DIR/targets/playdate.json.pristine"
fi
BASE_JSON=$(cat "$TINYGO_DIR/targets/playdate.json.pristine")
BASE_JSON="${BASE_JSON%?}"  # strip trailing }
cat > "$TINYGO_DIR/targets/playdate.json" << EOF
${BASE_JSON},
    "linkerscript": "$BUILD_DIR/playdate.ld",
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
# Uses the pdgo linker script (copied to $BUILD_DIR in Step 2), not the SDK's
# link_map.ld: the GC's root scanner needs _globals_start/_globals_end/_stack_top,
# which the SDK script does not define.
echo "Step 5: Linking..."
arm-none-eabi-gcc $MCFLAGS -T"$BUILD_DIR/playdate.ld" \
    -Wl,--gc-sections \
    -Wl,--emit-relocs \
    -nostartfiles \
    "$BUILD_DIR/setup.o" \
    "$BUILD_DIR/pd_runtime.o" \
    "$BUILD_DIR/game.o" \
    "$BUILD_DIR/asm_arm.o" \
    -lm \
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