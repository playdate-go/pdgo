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

$dependencies = @("go", "git", "mingw", "gcc-arm-none-eabi")
$testCommmands = @("go", "git", "mingw32-make", "arm-none-eabi-gcc")
$missing = @()

# Make sure Scoop shims are in the current session's PATH just in case 
# they were installed recently but the terminal wasn't restarted.
$scoopShims = Join-Path $env:USERPROFILE "scoop\shims"
if ((Test-Path $scoopShims) -and ($env:PATH -notmatch [regex]::Escape($scoopShims))) {
    $env:PATH = "$scoopShims;$env:PATH"
}

$index = 0
foreach ($cmd in $testCommmands) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        $missing += $dependencies[$index]
    }
    $index++
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
        $missingArgs = $missing -join " "
        Invoke-Expression "scoop bucket add extras"
        Invoke-Expression "scoop install $missingArgs"

        # Refresh environment variables to pick up any new paths set by Scoop
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        if ($env:PATH -notmatch [regex]::Escape($scoopShims)) {
            $env:PATH = "$scoopShims;$env:PATH"
        }
        
        # Verify installation succeeded
        $stillMissing = @()
        $installedList = (scoop list)
        foreach ($dep in $missing) {
            # The regex "^\b$packageName\b" ensures we match the exact package name at the start of a line
            # -Quiet makes Select-String return a simple $true or $false boolean
            $isInstalled = [bool]($installedList | Select-String -Pattern "^\b$dep\b" -Quiet)

            if (-not $isInstalled) {
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

    $pdgoReleaseData = Invoke-RestMethod -Uri "https://api.github.com/repos/playdate-go/pdgo/releases/latest"
    $pdgoReleaseTag = $pdgoReleaseData.tag_name

    $targetsDir = Join-Path $tinygoDir "targets"
    $runtimeDir = Join-Path (Join-Path $tinygoDir "src") "runtime"

    # targets/playdate.json
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/playdate-go/pdgo/refs/tags/$pdgoReleaseTag/cmd/pdgoc/tinygo-patches/playdate.json" -OutFile (Join-Path $targetsDir "playdate.json")
    # targets/playdate.ld
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/playdate-go/pdgo/refs/tags/$pdgoReleaseTag/cmd/pdgoc/tinygo-patches/playdate.ld" -OutFile (Join-Path $targetsDir "playdate.ld")
    # src/runtime/runtime_playdate.go
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/playdate-go/pdgo/refs/tags/$pdgoReleaseTag/cmd/pdgoc/tinygo-patches/runtime_playdate.go" -OutFile (Join-Path $runtimeDir "runtime_playdate.go")
    # src/runtime/gc_playdate.go
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/playdate-go/pdgo/refs/tags/$pdgoReleaseTag/cmd/pdgoc/tinygo-patches/gc_playdate.go" -OutFile (Join-Path $runtimeDir "gc_playdate.go")

    Write-Host "  Installing Go 1.25.8 as required by TinyGo:"
    Invoke-Expression "scoop install go@1.25.8"

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
Write-Host "Logout and login is recommended after installation`n"