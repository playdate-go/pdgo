<#
.SYNOPSIS
    pdgo Full Installer for Windows
.DESCRIPTION
    Installs everything needed to build Go games for Playdate:
      - pdgoc (build tool)
      - TinyGo with Playdate support (for device builds)
.EXAMPLE
    .\install.ps1
#>

$ErrorActionPreference = 'Stop'

# Options via environment variables
$skipTinyGo = $env:SKIP_TINYGO -eq "1"
$autoInstallDeps = $env:AUTO_INSTALL_DEPS -eq "1"
$jobs = if ($env:JOBS) { $env:JOBS } else { (Get-CimInstance Win32_ComputerSystem).NumberOfLogicalProcessors }

$tinygoDir = Join-Path $HOME "tinygo-playdate"
$tinygoVersion = "0.40.1"

Write-Host "╔══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║              PdGo - Full Installer (Windows)             ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# ============================================================================
# Step 1: Check dependencies
# ============================================================================
Write-Host "[1/4] Checking dependencies..." -ForegroundColor Yellow

$dependencies = @("go", "git")
$missing = @()

# Make sure Scoop shims are in the current session's PATH just in case 
# they were installed recently but the terminal wasn't restarted.
$scoopShims = Join-Path $env:USERPROFILE "scoop\shims"
if ((Test-Path $scoopShims) -and ($env:PATH -notmatch [regex]::Escape($scoopShims))) {
    $env:PATH = "$scoopShims;$env:PATH"
}

foreach ($dep in $dependencies) {
    if (-not (Get-Command $dep -ErrorAction SilentlyContinue)) {
        $missing += $dep
    }
}

if ($missing.Count -gt 0) {
    Write-Host "Missing dependencies: $($missing -join ', ')" -ForegroundColor Yellow
    
    $installNow = $autoInstallDeps
    if (-not $installNow) {
        Write-Host "We can automatically install these for you using Scoop (https://scoop.sh)."
        $response = Read-Host "Would you like to install missing dependencies via Scoop now? [Y/n]"
        if ($response -notmatch "^[Nn]") {
            $installNow = $true
        }
    }

    if ($installNow) {
        # Check if scoop itself is installed
        if (-not (Get-Command "scoop" -ErrorAction SilentlyContinue)) {
            Write-Host "`nScoop is not installed. Installing Scoop..." -ForegroundColor Cyan
            Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
            Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
            
            # Ensure scoop shims are in the current path for the rest of the script
            $env:PATH = "$scoopShims;$env:PATH"
        }
        
        Write-Host "`nInstalling dependencies via Scoop..." -ForegroundColor Cyan
        # Install missing dependencies (Scoop handles arrays nicely, but we'll do it cleanly)
        $missingArgs = $missing -join " "
        Invoke-Expression "scoop install $missingArgs"

        # Refresh environment variables to pick up any new paths set by Scoop
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        if ($env:PATH -notmatch [regex]::Escape($scoopShims)) {
            $env:PATH = "$scoopShims;$env:PATH"
        }
        
        # Verify installation succeeded
        $stillMissing = @()
        foreach ($dep in $missing) {
            if (-not (Get-Command $dep -ErrorAction SilentlyContinue)) {
                $stillMissing += $dep
            }
        }
        
        if ($stillMissing.Count -gt 0) {
            Write-Host "`nERROR: Failed to install the following dependencies: $($stillMissing -join ', ')" -ForegroundColor Red
            Write-Host "Please install them manually and try again."
            exit 1
        }
        Write-Host "Dependencies successfully installed!`n" -ForegroundColor Green
    } else {
        Write-Host "`nPlease install them manually and try again." -ForegroundColor Red
        Write-Host "You can use Scoop: scoop install $($missing -join ' ')"
        exit 1
    }
}

$goVersion = (go version).Split(' ')[2].Replace('go', '')
Write-Host "  Go: " -NoNewline; Write-Host $goVersion -ForegroundColor Green

# Check Playdate SDK
$sdkPath = $env:PLAYDATE_SDK_PATH
if ([string]::IsNullOrWhiteSpace($sdkPath)) {
    $sdkPath = Join-Path $HOME "Documents\PlaydateSDK"
    Write-Host "  PLAYDATE_SDK_PATH not set. Inferring default path..." -ForegroundColor DarkGray
}

if (-not (Test-Path $sdkPath)) {
    Write-Host "ERROR: Playdate SDK not found at: $sdkPath" -ForegroundColor Red
    Write-Host "Install it to the default location or set PLAYDATE_SDK_PATH."
        exit 1
}

$sdkVersionPath = Join-Path $sdkPath "VERSION.txt"
$sdkVersion = if (Test-Path $sdkVersionPath) { (Get-Content $sdkVersionPath).Trim() } else { "unknown" }
Write-Host "  Playdate SDK: " -NoNewline; Write-Host $sdkVersion -ForegroundColor Green

Write-Host "All dependencies OK" -ForegroundColor Green

# ============================================================================
# Step 2: Install pdgoc
# ============================================================================
Write-Host "`n[2/4] Installing pdgoc..." -ForegroundColor Yellow

go install github.com/playdate-go/pdgo/cmd/pdgoc@latest

$goPath = go env GOPATH
$goBin = if ($env:GOBIN) { $env:GOBIN } else { Join-Path $goPath "bin" }

if (Test-Path (Join-Path $goBin "pdgoc.exe")) {
    Write-Host "  pdgoc installed at: " -NoNewline; Write-Host "$goBin\pdgoc.exe" -ForegroundColor Green
} else {
    Write-Host "Failed to install pdgoc" -ForegroundColor Red
    exit 1
}

# ============================================================================
# Step 3: Install TinyGo with Playdate support (Pre-Compiled)
# ============================================================================
if ($skipTinyGo) {
    Write-Host "`n[3/4] Skipping TinyGo install (SKIP_TINYGO=1)" -ForegroundColor Yellow
} else {
    Write-Host "`n[3/4] Installing TinyGo with Playdate support..." -ForegroundColor Yellow
    Write-Host "  TinyGo version: " -NoNewline; Write-Host "v$tinygoVersion (Pre-compiled)" -ForegroundColor Green
    Write-Host "  Install path: " -NoNewline; Write-Host $tinygoDir -ForegroundColor Green
    Write-Host ""

    $tempZip = Join-Path $env:TEMP "tinygo.zip"
    $tinygoUrl = "https://github.com/tinygo-org/tinygo/releases/download/v$tinygoVersion/tinygo$tinygoVersion.windows-amd64.zip"

    Write-Host "  Downloading official Windows release..."
    Invoke-WebRequest -Uri $tinygoUrl -OutFile $tempZip

    if (Test-Path $tinygoDir) {
        Write-Host "  Removing old installation..."
        Remove-Item -Recurse -Force $tinygoDir
    }

    Write-Host "  Extracting to target directory..."
    Expand-Archive -Path $tempZip -DestinationPath $HOME -Force
    Rename-Item -Path (Join-Path $HOME "tinygo") -NewName (Split-Path $tinygoDir -Leaf)

    Remove-Item $tempZip

    Write-Host "  Injecting Playdate support files..."

    $targetsDir = Join-Path $tinygoDir "targets"
    $runtimeDir = Join-Path (Join-Path $tinygoDir "src") "runtime"

    # targets/playdate.json
    @'
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
'@ | Set-Content (Join-Path $targetsDir "playdate.json") -Encoding UTF8

    # targets/playdate.ld
    @'
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
'@ | Set-Content (Join-Path $targetsDir "playdate.ld") -Encoding UTF8

    # src/runtime/runtime_playdate.go
    @'
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
'@ | Set-Content (Join-Path $runtimeDir "runtime_playdate.go") -Encoding UTF8

    # src/runtime/gc_playdate.go
    @'
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
'@ | Set-Content (Join-Path $runtimeDir "gc_playdate.go") -Encoding UTF8
    Write-Host "  TinyGo with Playdate support ready!" -ForegroundColor Green
}

# ============================================================================
# Step 4: Setup PATH
# ============================================================================
Write-Host "`n[4/4] Configuring PATH..." -ForegroundColor Yellow

$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$newPaths = @()
$needsUpdate = $false

if ($userPath -notmatch [regex]::Escape($goBin)) {
    $newPaths += $goBin
    $needsUpdate = $true
}

$tinygoBinDir = Join-Path $tinygoDir "bin"
if (-not $skipTinyGo -and $userPath -notmatch [regex]::Escape($tinygoBinDir)) {
    $newPaths += $tinygoBinDir
    $needsUpdate = $true
}

if ($needsUpdate) {
    Write-Host "`nThe following paths need to be added to your User PATH:" -ForegroundColor Yellow
    foreach ($p in $newPaths) {
        Write-Host "  $p"
    }
    Write-Host ""
    
    $response = Read-Host "Add them automatically to your User Environment? [Y/n]"
    if ($response -notmatch "^[Nn]") {
        $updatedPath = $userPath
        if (-not $updatedPath.EndsWith(";")) { $updatedPath += ";" }
        $updatedPath += ($newPaths -join ";")
        
        [Environment]::SetEnvironmentVariable("PATH", $updatedPath, "User")
        $env:PATH = $updatedPath
        
        Write-Host "Added to User PATH. (You may need to restart your terminal for all changes to take effect)" -ForegroundColor Green
    }
} else {
    Write-Host "PATH already configured" -ForegroundColor Green
}

# ============================================================================
# Done!
# ============================================================================
Write-Host "`n╔══════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║              Installation Complete!                      ║" -ForegroundColor Green
Write-Host "╚══════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host "`nInstalled:"
Write-Host "  pdgoc:  " -NoNewline; Write-Host "$goBin\pdgoc.exe" -ForegroundColor Green
if (-not $skipTinyGo) {
    Write-Host "  TinyGo: " -NoNewline; Write-Host "$tinygoBinDir\tinygo.exe" -ForegroundColor Green
}

Write-Host "`nUsage:`n"
Write-Host "  # Simulator build (run in simulator)"
Write-Host "  pdgoc -sim -run -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
Write-Host ""

if (-not $skipTinyGo) {
    Write-Host "  # Device build (for real Playdate)"
    Write-Host "  pdgoc -device -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
    Write-Host ""
    Write-Host "  # Device build + deploy to connected Playdate"
    Write-Host "  pdgoc -device -deploy -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
    Write-Host ""
}

Write-Host "Documentation: https://github.com/playdate-go/pdgo`n"