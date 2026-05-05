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
$targetsDir = Join-Path $TinyGoDir "targets"
$ldSource = Join-Path $targetsDir "playdate.ld"
$jsonSource = Join-Path $targetsDir "playdate.json"
$ldDest = Join-Path $BuildDirAbs "playdate.ld"

# Copy linker script from TinyGo installation (installed by install.sh/ps1)
if (-not (Test-Path $ldSource)) {
    Write-Error "Error: playdate.ld not found in $targetsDir\ — run install script first"
    exit 1
}
Copy-Item -Path $ldSource -Destination $ldDest -Force

# Extend the installed playdate.json with build-specific linkerscript and ldflags
if (-not (Test-Path $jsonSource)) {
    Write-Error "Error: playdate.json not found in $targetsDir\ — run install script first"
    exit 1
}

# Replace backslashes with forward slashes to prevent JSON escaping errors on Windows
$BuildDirJson = $BuildDirAbs -replace '\\', '/'

# Safely parse and update JSON properties (Idempotent equivalent to the bash cat << EOF append)
$targetJson = Get-Content -Path $jsonSource -Raw | ConvertFrom-Json
$targetJson | Add-Member -NotePropertyName "linkerscript" -NotePropertyValue "$BuildDirJson/playdate.ld" -Force
$targetJson | Add-Member -NotePropertyName "ldflags" -NotePropertyValue @("-L$BuildDirJson", "-lpd") -Force

$targetJson | ConvertTo-Json -Depth 10 | Set-Content -Path $jsonSource

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