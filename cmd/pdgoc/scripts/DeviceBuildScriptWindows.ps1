<#
.SYNOPSIS
Build script for Playdate device builds using TinyGo.

.DESCRIPTION
Compiles C runtime, builds Go code, links artifacts, and creates a .pdx bundle.
#>

param (
    [string]$GameName = "$env:GAME_NAME",
    [string]$GameDir = "$env:GAME_DIR",
    [string]$GoSrcDir = "$env:GO_SRC_DIR",
    [string]$BuildDir = "$env:BUILD_DIR",
    [string]$TinyGoPath = "$env:TINYGO_PLAYDATE",
    [string]$TinyGoDir = "$env:TINYGO_PLAYDATE_DIR",
    [string]$SdkPath = "$env:PLAYDATE_SDK_PATH"
)

# Ensure scoop\shim is the first in path
$scoopShims = Join-Path $env:USERPROFILE "scoop\shims"
$env:PATH = "$scoopShims;$env:PATH"

# Equivalent to set -e
$ErrorActionPreference = "Stop"

# Validate requirements
if (-not (Test-Path $SdkPath)) { Write-Error "Error: Playdate SDK not found at $SdkPath"; exit 1 }
if (-not (Get-Command $TinyGoPath -ErrorAction Ignore)) { Write-Error "Error: TinyGo not found at $TinyGoPath"; exit 1 }
if (-not (Get-Command "arm-none-eabi-gcc" -ErrorAction Ignore)) { Write-Error "Error: arm-none-eabi-gcc not found in PATH"; exit 1 }

# Ensure directories exist and get absolute paths to avoid context issues
if (-not (Test-Path $GameDir)) { New-Item -ItemType Directory -Force -Path $GameDir | Out-Null }
if (-not (Test-Path $GoSrcDir)) { New-Item -ItemType Directory -Force -Path $GoSrcDir | Out-Null }
if (-not (Test-Path $BuildDir)) { New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null }

$GameDirAbs = (Resolve-Path $GameDir).Path
$GoSrcDirAbs = (Resolve-Path $GoSrcDir).Path
$BuildDirAbs = (Resolve-Path $BuildDir).Path

Write-Host "=== Building $GameName for Playdate ==="
Write-Host "Go source: $GoSrcDirAbs"
Write-Host "Build dir: $BuildDirAbs"

# ARM compiler flags
$McFlags = @("-mthumb", "-mcpu=cortex-m7", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16")
$CpFlags = @("-ffunction-sections", "-fdata-sections", "-mword-relocations", "-fno-common")

# cd "$GO_SRC_DIR"
Set-Location $GoSrcDirAbs

# Step 1: Compile C runtime (pd_runtime.c)
Write-Host "Step 1: Compiling C runtime..."
& "arm-none-eabi-gcc" @McFlags @CpFlags -O2 -DTARGET_PLAYDATE=1 -DTARGET_EXTENSION=1 "-I$SdkPath/C_API" -c "$BuildDirAbs/pd_runtime.c" -o "$BuildDirAbs/pd_runtime.o"

& "arm-none-eabi-ar" rcs "$BuildDirAbs/libpd.a" "$BuildDirAbs/pd_runtime.o"

# Step 2: Create linker script and TinyGo target
Write-Host "Step 2: Configuring build..."
$ldScript = @"
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
"@
Set-Content -Path "$BuildDirAbs/playdate.ld" -Value $ldScript

$targetsDir = Join-Path $TinyGoDir "targets"
if (-not (Test-Path $targetsDir)) { New-Item -ItemType Directory -Force -Path $targetsDir | Out-Null }

# Replace backslashes with forward slashes to prevent JSON escaping errors on Windows
$BuildDirJson = $BuildDirAbs -replace '\\', '/'

$targetJson = @"
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
    "linkerscript": "$BuildDirJson/playdate.ld",
    "cflags": ["-DTARGET_PLAYDATE=1", "-mfloat-abi=hard", "-mfpu=fpv5-sp-d16"],
    "ldflags": ["-L$BuildDirJson", "-lpd"]
}
"@
Set-Content -Path "$targetsDir/playdate.json" -Value $targetJson

# Step 3: Build with TinyGo
Write-Host "Step 3: Compiling Go code with TinyGo..."
scoop reset go@1.25.8
& $TinyGoPath build -target=playdate -o "$BuildDirAbs/game.o" .

# Step 4: Compile setup.c from SDK
Write-Host "Step 4: Compiling SDK setup..."
& "arm-none-eabi-gcc" @McFlags @CpFlags -O2 -DTARGET_PLAYDATE=1 -DTARGET_EXTENSION=1 "-I$SdkPath/C_API" -c "$SdkPath/C_API/buildsupport/setup.c" -o "$BuildDirAbs/setup.o"

# Step 5: Link everything
Write-Host "Step 5: Linking..."
& "arm-none-eabi-gcc" @McFlags "-T$SdkPath/C_API/buildsupport/link_map.ld" "-Wl,--gc-sections" "-Wl,--emit-relocs" "-nostartfiles" "$BuildDirAbs/setup.o" "$BuildDirAbs/pd_runtime.o" "$BuildDirAbs/game.o" -o "$BuildDirAbs/pdex.elf"

Copy-Item "$BuildDirAbs/pdex.elf" -Destination "$GoSrcDirAbs/"

Write-Host "ELF size:"
& "arm-none-eabi-size" "$BuildDirAbs/pdex.elf"

# Step 6: Create PDX bundle
Write-Host "Step 6: Creating .pdx bundle..."
$PdcPath = Join-Path $SdkPath "bin/pdc"
& $PdcPath -k "$GoSrcDirAbs" "$GameDirAbs/${GameName}.pdx"

# Cleanup
Remove-Item -Path "$GoSrcDirAbs/pdex.elf" -Force -ErrorAction Ignore

Write-Host "`n=== Build complete: $GameDirAbs/${GameName}.pdx ==="