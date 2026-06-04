$ErrorActionPreference = 'Stop'

Write-Host 'Building all pdgo examples'
Write-Host '=============================='
Write-Host ''

$examples = @(
    'json',
    'json_lowlevel',
    '3d_library',
    'bach_midi',
    'spritegame',
    'life',
    'bouncing_square',
    'go_logo',
    'hello_world',
    'tilemap',
    'sprite_collisions',
    'particles',
    'exposure'
)

$failed = @()
$succeeded = @()

foreach ($example in $examples) {
    Write-Host "Building $example..."
    Write-Host '------------------------'

    $buildScript = Join-Path $PSScriptRoot "$example\build.ps1"
    if ((Test-Path (Join-Path $PSScriptRoot $example)) -and (Test-Path $buildScript)) {
        Set-Location (Join-Path $PSScriptRoot $example)
        try {
            & $buildScript
            if ($LASTEXITCODE -ne 0) {
                throw "pdgoc exited with code $LASTEXITCODE"
            }
            $succeeded += $example
            Write-Host "$example built successfully" -ForegroundColor Green
        } catch {
            $failed += $example
            Write-Host "$example failed to build" -ForegroundColor Red
            Write-Host $_.Exception.Message -ForegroundColor Red
        }
    } else {
        Write-Host "Skipping $example (no build.ps1 found)" -ForegroundColor Yellow
    }
    Write-Host ''
}

Write-Host '=============================='
Write-Host 'Build Summary'
Write-Host '=============================='
Write-Host "Succeeded: $($succeeded.Count)" -ForegroundColor Green
foreach ($ex in $succeeded) {
    Write-Host "   - $ex" -ForegroundColor Green
}

if ($failed.Count -gt 0) {
    Write-Host "Failed: $($failed.Count)" -ForegroundColor Red
    foreach ($ex in $failed) {
        Write-Host "   - $ex" -ForegroundColor Red
    }
    exit 1
}

Write-Host ''
Write-Host 'All examples built successfully!' -ForegroundColor Green
