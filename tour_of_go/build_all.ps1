$ErrorActionPreference = 'Stop'

Write-Host 'Building all Tour of Go examples'
Write-Host '=============================='
Write-Host ''

$failed = @()
$succeeded = @()

$dirs = Get-ChildItem -Path $PSScriptRoot -Directory | Where-Object {
    $_.Name -match '^\d' -or $_.Name -match '^generics_'
}

foreach ($dir in $dirs) {
    $buildScript = Join-Path $dir.FullName 'build.ps1'
    if (Test-Path $buildScript) {
        $name = $dir.Name
        Write-Host "Building $name..."
        Write-Host '------------------------'

        Set-Location $dir.FullName
        try {
            & $buildScript
            if ($LASTEXITCODE -ne 0) {
                throw "pdgoc exited with code $LASTEXITCODE"
            }
            $succeeded += $name
            Write-Host "$name built successfully" -ForegroundColor Green
        } catch {
            $failed += $name
            Write-Host "$name failed to build" -ForegroundColor Red
            Write-Host $_.Exception.Message -ForegroundColor Red
        }
        Write-Host ''
    }
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
