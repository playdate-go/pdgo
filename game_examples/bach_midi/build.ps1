# Build script for Windows - all complexity is handled by pdgoc

# Copy assets from SDK if not already present
if (-not (Test-Path 'Source\bach.mid')) {
    Write-Host 'Copying bach.mid from Playdate SDK...'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\bach.mid\Source\bach.mid') 'Source\'
}

if (-not (Test-Path 'Source\piano.wav')) {
    Write-Host 'Copying piano.wav from Playdate SDK...'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\bach.mid\Source\piano.wav') 'Source\'
}

pdgoc -sim -device -name="BachMIDI" -author="PdGo" -desc="Bach MIDI Player" -bundle-id=com.pdgo.bachmidi -version=1.0 -build-number=1
