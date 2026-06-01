<#
.SYNOPSIS
    pdgo Full Installer for Windows
.DESCRIPTION
    Installs everything needed to build Go games for Playdate:
      - pdgoc (build tool)
      - TinyGo with Playdate support (for device builds)
.EXAMPLE
    .\install.ps1
.EXAMPLE
    .\install.ps1 -CleanInstall
.PARAMETER CleanInstall
    Remove all previously installed PdGo components before installing fresh:
      - ~/tinygo-playdate/ directory
      - pdgoc.exe from GOBIN
#>

param(
    [switch]$CleanInstall
)

$ErrorActionPreference = 'Stop'

$isLocalBuild = (Test-Path -Path .\cmd\pdgoc) -and (Test-Path -Path .\go.mod)
$localRoot = (Get-Location).Path

# Options via environment variables
$skipTinyGo = $env:SKIP_TINYGO -eq '1'
$autoInstallDeps = $env:AUTO_INSTALL_DEPS -eq '1'

$tinygoDir = Join-Path $HOME 'tinygo-playdate'

# ============================================================================
# Clean install: remove previous PdGo components
# ============================================================================
if ($CleanInstall) {
    Write-Host 'Cleaning previous PdGo installation...' -ForegroundColor Yellow

    $removed = @()

    if (Test-Path $tinygoDir) {
        Write-Host "  Removing TinyGo: $tinygoDir"
        Remove-Item -Recurse -Force $tinygoDir
        $removed += "TinyGo ($tinygoDir)"
    }

    $goPath = go env GOPATH 2> $null
    $goBin = if ($env:GOBIN) { $env:GOBIN } else { Join-Path $goPath 'bin' }
    $pdgocExe = Join-Path $goBin 'pdgoc.exe'
    if (Test-Path $pdgocExe) {
        Write-Host "  Removing pdgoc: $pdgocExe"
        Remove-Item -Force $pdgocExe
        $removed += "pdgoc ($pdgocExe)"
    }

    # Remove dependencies installed by the script via Scoop
    if (Get-Command 'scoop' -ErrorAction SilentlyContinue) {
        $scoopDeps = @('mingw', 'gcc-arm-none-eabi', 'go')
        $installedList = scoop list 2> $null
        foreach ($dep in $scoopDeps) {
            $isInstalled = $false
            foreach ($line in $installedList) {
                if ($line -match "^\s*$dep\b") {
                    $isInstalled = $true
                    break
                }
            }
            if ($isInstalled) {
                Write-Host "  Removing Scoop package: $dep"
                $uninstallOutput = scoop uninstall $dep 2>&1
                if ($LASTEXITCODE -eq 0) {
                    $removed += "Scoop package ($dep)"
                } else {
                    Write-Host "    Warning: failed to uninstall $dep" -ForegroundColor Yellow
                }
            }
        }
    }

    if ($removed.Count -gt 0) {
        Write-Host '  Cleaned:' -ForegroundColor Green
        foreach ($item in $removed) {
            Write-Host "    - $item" -ForegroundColor Green
        }
    } else {
        Write-Host '  Nothing to clean' -ForegroundColor DarkGray
    }
    Write-Host ''
}
$tinygoVersion = '0.40.1'
$requiredGoVersion = '1.25.8'

# Version info for building pdgoc and fetching patches
$pdgoVersion = 'latest'
$pdgoCommit = 'unknown'
$pdgoDate = (Get-Date).ToUniversalTime().ToString('yyyy-MM-dd HH:mm:ss UTC')

if ($isLocalBuild) {
    try {
        $pdgoVersion = git describe --tags --always --dirty
        $pdgoCommit = git rev-parse --short HEAD
    } catch {
        Write-Host 'Warning: Failed to get local version info from git. Using defaults.' -ForegroundColor Yellow
    }
} else {
    try {
        $githubApi = 'https://api.github.com/repos/playdate-go/pdgo'
        $headers = @{'User-Agent' = 'pdgo-installer'}

        $releaseData = Invoke-RestMethod -Uri "$githubApi/releases/latest" -Headers $headers -ErrorAction SilentlyContinue
        if ($releaseData -and $releaseData.tag_name) {
            $pdgoVersion = $releaseData.tag_name
        }

        $commitData = Invoke-RestMethod -Uri "$githubApi/commits/main" -Headers $headers -ErrorAction SilentlyContinue
        if ($commitData -and $commitData.sha) {
            $pdgoCommit = $commitData.sha.Substring(0, 7)
        }
    } catch {
        Write-Host 'Warning: Failed to fetch version info from GitHub API. Using defaults.' -ForegroundColor Yellow
    }
}

Write-Host ''
Write-Host '  PdGo - Full Installer (Windows)' -ForegroundColor Cyan
Write-Host ''

# ============================================================================
# Step 1: Check dependencies
# ============================================================================
Write-Host '[1/4] Checking dependencies...' -ForegroundColor Yellow

$dependencies = @('go', 'git', 'mingw', 'gcc-arm-none-eabi')
$testCommands = @('go', 'git', 'gcc', 'arm-none-eabi-gcc')
$missing = @()

# Make sure Scoop shims are in the current session PATH just in case
# they were installed recently but the terminal was not restarted.
$scoopShims = Join-Path $env:USERPROFILE 'scoop\shims'
if ((Test-Path $scoopShims) -and ($env:PATH -notmatch [regex]::Escape($scoopShims))) {
    $env:PATH = $scoopShims + ';' + $env:PATH
}

$index = 0
foreach ($cmd in $testCommands) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        $missing += $dependencies[$index]
    }
    $index++
}

if ($missing.Count -gt 0) {
    Write-Host "Missing dependencies: $($missing -join ', ')" -ForegroundColor Yellow

    $installNow = $autoInstallDeps
    if (-not $installNow) {
        Write-Host 'We can automatically install these for you using Scoop (https://scoop.sh).'
        if ($env:CI -eq '1') {
            Write-Host 'CI detected - skipping interactive prompt, installing dependencies automatically' -ForegroundColor Yellow
            $installNow = $true
        } else {
            $response = Read-Host 'Would you like to install missing dependencies via Scoop now? [Y/n]'
            if ($response -notmatch '^[Nn]') {
                $installNow = $true
            }
        }
    }

    if ($installNow) {
        # Check if scoop itself is installed
        if (-not (Get-Command 'scoop' -ErrorAction SilentlyContinue)) {
            Write-Host ''
            Write-Host 'Scoop is not installed. Installing Scoop...' -ForegroundColor Cyan
            Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
            Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

            # Ensure scoop shims are in the current path for the rest of the script
            $env:PATH = $scoopShims + ';' + $env:PATH
        }

        Write-Host ''
        Write-Host 'Installing dependencies via Scoop...' -ForegroundColor Cyan
        $missingArgs = $missing -join ' '
        Invoke-Expression 'scoop bucket add extras' 2> $null
        Invoke-Expression "scoop install $missingArgs"

        # Refresh environment variables to pick up any new paths set by Scoop
        $env:PATH = [System.Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [System.Environment]::GetEnvironmentVariable('Path','User')
        if ($env:PATH -notmatch [regex]::Escape($scoopShims)) {
            $env:PATH = $scoopShims + ';' + $env:PATH
        }

        # Verify installation succeeded
        $stillMissing = @()
        $installedList = (scoop list)
        foreach ($dep in $missing) {
            $isInstalled = [bool]($installedList | Select-String -Pattern "^\b$dep\b" -Quiet)

            if (-not $isInstalled) {
                $stillMissing += $dep
            }
        }

        if ($stillMissing.Count -gt 0) {
            Write-Host ''
            Write-Host "ERROR: Failed to install the following dependencies: $($stillMissing -join ', ')" -ForegroundColor Red
            Write-Host 'Please install them manually and try again.'
            exit 1
        }
        Write-Host 'Dependencies successfully installed!' -ForegroundColor Green
        Write-Host ''
    } else {
        Write-Host ''
        Write-Host 'Please install them manually and try again.' -ForegroundColor Red
        Write-Host "You can use Scoop: scoop install $($missing -join ' ')"
        exit 1
    }
}

$goVersion = (go version).Split(' ')[2].Replace('go', '')
Write-Host '  Go: ' -NoNewline; Write-Host $goVersion -ForegroundColor Green

# Check Playdate SDK
$sdkPath = $env:PLAYDATE_SDK_PATH
if ([string]::IsNullOrWhiteSpace($sdkPath)) {
    $sdkPath = Join-Path $HOME 'Documents\PlaydateSDK'
    Write-Host '  PLAYDATE_SDK_PATH not set. Inferring default path...' -ForegroundColor DarkGray
}

if (-not (Test-Path $sdkPath)) {
    Write-Host "ERROR: Playdate SDK not found at: $sdkPath" -ForegroundColor Red
    Write-Host 'Install it to the default location or set PLAYDATE_SDK_PATH.'
    exit 1
}

$sdkVersionPath = Join-Path $sdkPath 'VERSION.txt'
$sdkVersion = if (Test-Path $sdkVersionPath) { (Get-Content $sdkVersionPath).Trim() } else { 'unknown' }
Write-Host '  Playdate SDK: ' -NoNewline; Write-Host $sdkVersion -ForegroundColor Green

Write-Host 'All dependencies OK' -ForegroundColor Green

# Pin Go to the version required by TinyGo
$goVersion = (go version).Split(' ')[2].Replace('go', '')
if ($goVersion -ne $requiredGoVersion) {
    Write-Host ''
    Write-Host "Go $goVersion detected, but TinyGo requires Go $requiredGoVersion." -ForegroundColor Yellow
    Write-Host "Installing Go $requiredGoVersion via Scoop..." -ForegroundColor Yellow
    if (Get-Command 'scoop' -ErrorAction SilentlyContinue) {
        Invoke-Expression "scoop install go@$requiredGoVersion"
        # Refresh PATH
        $env:PATH = [System.Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [System.Environment]::GetEnvironmentVariable('Path','User')
        $scoopShims = Join-Path $env:USERPROFILE 'scoop\shims'
        if ($env:PATH -notmatch [regex]::Escape($scoopShims)) {
            $env:PATH = $scoopShims + ';' + $env:PATH
        }
        $goVersionCheck = (go version).Split(' ')[2].Replace('go', '')
        if ($goVersionCheck -eq $requiredGoVersion) {
            Write-Host "  Go $requiredGoVersion installed successfully." -ForegroundColor Green
        } else {
            Write-Host "  WARNING: Go version is still $goVersionCheck. Build may fail." -ForegroundColor Yellow
        }
    } else {
        Write-Host "  Scoop not found. Please install Go $requiredGoVersion manually." -ForegroundColor Red
        exit 1
    }
}

# ============================================================================
# Step 2: Install pdgoc
# ============================================================================
Write-Host ''
Write-Host '[2/4] Installing pdgoc...' -ForegroundColor Yellow

$goPath = go env GOPATH
$goBin = if ($env:GOBIN) { $env:GOBIN } else { Join-Path $goPath 'bin' }
if (-not (Test-Path $goBin)) {
    New-Item -ItemType Directory -Path $goBin | Out-Null
}

if ($isLocalBuild) {
    Write-Host 'Installing pdgoc from local directory'
    $originalLocation = (Get-Location).Path
    Set-Location $localRoot\cmd\pdgoc

    Write-Host '  Downloading Go dependencies...'
    go mod tidy

    # Build with ldflags to inject version information
    go build "-ldflags=-X 'main.Version=$pdgoVersion' -X 'main.Commit=$pdgoCommit' -X 'main.Date=$pdgoDate'" -o (Join-Path $goBin 'pdgoc.exe') .

    Set-Location $originalLocation
} else {
    Write-Host '  Installing pdgoc from GitHub...'
    $tempDir = Join-Path $env:TEMP "pdgo-build-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir | Out-Null

    try {
        $zipUrl = ''
        if ($pdgoVersion -ne 'latest') {
            Write-Host "  Downloading release $pdgoVersion..."
            $zipUrl = "https://github.com/playdate-go/pdgo/archive/refs/tags/$pdgoVersion.zip"
        } else {
            Write-Host '  Downloading latest source from main...'
            $zipUrl = 'https://github.com/playdate-go/pdgo/archive/refs/heads/main.zip'
        }

        $tempZip = Join-Path $tempDir 'source.zip'
        Invoke-WebRequest -Uri $zipUrl -OutFile $tempZip

        Write-Host '  Extracting source...'
        Expand-Archive -Path $tempZip -DestinationPath $tempDir -Force

        $extractedDir = Get-ChildItem -Path $tempDir -Directory | Select-Object -First 1 | ForEach-Object { $_.FullName }
        $pdgocSourceDir = Join-Path $extractedDir 'cmd\pdgoc'

        $originalLocation = (Get-Location).Path
        Set-Location $pdgocSourceDir

        Write-Host '  Downloading Go dependencies...'
        go mod tidy

        Write-Host '  Building pdgoc...'
        Write-Host "    Version: $pdgoVersion"
        Write-Host "    Commit:  $pdgoCommit"
        Write-Host "    Date:    $pdgoDate"

        # Build with ldflags to inject version information
        go build "-ldflags=-X 'main.Version=$pdgoVersion' -X 'main.Commit=$pdgoCommit' -X 'main.Date=$pdgoDate'" -o (Join-Path $goBin 'pdgoc.exe') .
        if ($LASTEXITCODE -ne 0) {
            Write-Host 'ERROR: Failed to build pdgoc' -ForegroundColor Red
            Set-Location $originalLocation
            exit 1
        }

        Set-Location $originalLocation
        Write-Host '  pdgoc successfully built and installed!' -ForegroundColor Green
    }
    finally {
        if (Test-Path $tempDir) {
            Remove-Item -Recurse -Force $tempDir
        }
    }
}

if (Test-Path (Join-Path $goBin 'pdgoc.exe')) {
    Write-Host "  pdgoc installed at: $goBin\pdgoc.exe" -ForegroundColor Green
} else {
    Write-Host 'Failed to install pdgoc' -ForegroundColor Red
    exit 1
}

# ============================================================================
# Step 3: Install TinyGo with Playdate support (Pre-Compiled)
# ============================================================================
if ($skipTinyGo) {
    Write-Host ''
    Write-Host '[3/4] Skipping TinyGo install (SKIP_TINYGO=1)' -ForegroundColor Yellow
} else {
    $tinygoBinPath = Join-Path $tinygoDir 'bin\tinygo.exe'
    $tinygoTargetsJson = Join-Path $tinygoDir 'targets\playdate.json'
    $tinygoAlreadyInstalled = (Test-Path $tinygoBinPath) -and (Test-Path $tinygoTargetsJson)

    if ($tinygoAlreadyInstalled) {
        Write-Host ''
        Write-Host '[3/4] TinyGo already installed, skipping download.' -ForegroundColor Green
        Write-Host "  TinyGo: $tinygoBinPath" -ForegroundColor Green
    } else {
        Write-Host ''
        Write-Host '[3/4] Installing TinyGo with Playdate support...' -ForegroundColor Yellow
        Write-Host "  TinyGo version: v$tinygoVersion (Pre-compiled)" -ForegroundColor Green
        Write-Host "  Install path: $tinygoDir" -ForegroundColor Green
        Write-Host ''

        $tempZip = Join-Path $env:TEMP 'tinygo.zip'
        $tinygoUrl = "https://github.com/tinygo-org/tinygo/releases/download/v$tinygoVersion/tinygo$tinygoVersion.windows-amd64.zip"

        Write-Host '  Downloading official Windows release...'
        Invoke-WebRequest -Uri $tinygoUrl -OutFile $tempZip

        if (Test-Path $tinygoDir) {
            Write-Host '  Removing old installation...'
            Remove-Item -Recurse -Force $tinygoDir
        }

        Write-Host '  Extracting to target directory...'
        Expand-Archive -Path $tempZip -DestinationPath $HOME -Force
        Rename-Item -Path (Join-Path $HOME 'tinygo') -NewName (Split-Path $tinygoDir -Leaf)

        Remove-Item $tempZip

        Write-Host '  Injecting Playdate support files...'

        $targetsDir = Join-Path $tinygoDir 'targets'
        $runtimeDir = Join-Path (Join-Path $tinygoDir 'src') 'runtime'

        if ($isLocalBuild) {
            Write-Host 'Patching TinyGo with local files'
            Set-Location $localRoot
            Copy-Item -Path '.\cmd\pdgoc\tinygo-patches\playdate.json' -Destination (Join-Path $targetsDir 'playdate.json')
            Copy-Item -Path '.\cmd\pdgoc\tinygo-patches\playdate.ld' -Destination (Join-Path $targetsDir 'playdate.ld')
            Copy-Item -Path '.\cmd\pdgoc\tinygo-patches\runtime_playdate.go' -Destination (Join-Path $runtimeDir 'runtime_playdate.go')
            Copy-Item -Path '.\cmd\pdgoc\tinygo-patches\gc_playdate.go' -Destination (Join-Path $runtimeDir 'gc_playdate.go')
        } else {
            # Determine the branch/tag/ref to use for downloading patches
            $pdgoBranch = if ($pdgoVersion -ne 'latest') { "refs/tags/$pdgoVersion" } else { 'main' }
            $patchesBaseUrl = "https://raw.githubusercontent.com/playdate-go/pdgo/$pdgoBranch/cmd/pdgoc/tinygo-patches"

            Write-Host "  Downloading patches from $patchesBaseUrl..."
            Invoke-WebRequest -Uri "$patchesBaseUrl/playdate.json" -OutFile (Join-Path $targetsDir 'playdate.json')
            Invoke-WebRequest -Uri "$patchesBaseUrl/playdate.ld" -OutFile (Join-Path $targetsDir 'playdate.ld')
            Invoke-WebRequest -Uri "$patchesBaseUrl/runtime_playdate.go" -OutFile (Join-Path $runtimeDir 'runtime_playdate.go')
            Invoke-WebRequest -Uri "$patchesBaseUrl/gc_playdate.go" -OutFile (Join-Path $runtimeDir 'gc_playdate.go')
        }

        Write-Host '  TinyGo with Playdate support ready!' -ForegroundColor Green
    }
}

# ============================================================================
# Step 4: Setup PATH
# ============================================================================
Write-Host ''
Write-Host '[4/4] Configuring PATH...' -ForegroundColor Yellow

$userPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
$newPaths = @()
$needsUpdate = $false

# pdgoc (GOBIN)
if ($userPath -notmatch [regex]::Escape($goBin)) {
    $newPaths += $goBin
    $needsUpdate = $true
}

# TinyGo
$tinygoBinDir = Join-Path $tinygoDir 'bin'
if (-not $skipTinyGo -and $userPath -notmatch [regex]::Escape($tinygoBinDir)) {
    $newPaths += $tinygoBinDir
    $needsUpdate = $true
}

# arm-none-eabi-gcc (installed by Scoop to its own dir)
$armGccCmd = Get-Command 'arm-none-eabi-gcc' -ErrorAction SilentlyContinue
if ($armGccCmd) {
    $armGccDir = Split-Path $armGccCmd.Source
    if ($userPath -notmatch [regex]::Escape($armGccDir)) {
        $newPaths += $armGccDir
        $needsUpdate = $true
    }
}

if ($needsUpdate) {
    Write-Host ''
    Write-Host 'The following paths need to be added to your User PATH:' -ForegroundColor Yellow
    foreach ($p in $newPaths) {
        Write-Host "  $p"
    }
    Write-Host ''

    if ($env:CI -eq '1') {
        Write-Host 'CI detected - skipping interactive PATH setup' -ForegroundColor Yellow
    } else {
        $response = Read-Host 'Add them automatically to your User Environment? [Y/n]'
        if ($response -notmatch '^[Nn]') {
            $updatedPath = $userPath
            if (-not $updatedPath.EndsWith(';')) { $updatedPath += ';' }
            $updatedPath += ($newPaths -join ';')

            [Environment]::SetEnvironmentVariable('PATH', $updatedPath, 'User')
            $env:PATH = $updatedPath

            Write-Host 'Added to User PATH. (You may need to restart your terminal for all changes to take effect)' -ForegroundColor Green
        }
    }
} else {
    Write-Host 'PATH already configured' -ForegroundColor Green
}

# ============================================================================
# Done!
# ============================================================================
Write-Host ''
Write-Host '  Installation Complete!' -ForegroundColor Green
Write-Host ''
Write-Host 'Installed:'
Write-Host "  pdgoc:  $goBin\pdgoc.exe" -ForegroundColor Green
if (-not $skipTinyGo) {
    Write-Host "  TinyGo: $tinygoBinDir\tinygo.exe" -ForegroundColor Green
}

Write-Host ''
Write-Host 'Usage:'
Write-Host ''
Write-Host "  # Simulator build (run in simulator)"
Write-Host "  pdgoc -sim -run -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
Write-Host ''

if (-not $skipTinyGo) {
    Write-Host '  # Device build (for real Playdate)'
    Write-Host "  pdgoc -device -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
    Write-Host ''
    Write-Host '  # Device build + deploy to connected Playdate'
    Write-Host "  pdgoc -device -deploy -name MyGame -author Me -desc 'Game' -bundle-id com.me.game -version 1.0 -build-number 1"
    Write-Host ''
}

Write-Host 'Documentation: https://github.com/playdate-go/pdgo'
Write-Host ''
Write-Host 'Logout and login is recommended after installation'
Write-Host ''