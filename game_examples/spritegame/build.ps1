# Build script for Windows - all complexity is handled by pdgoc

# Copy images from SDK Sprite Game example if not already present
if (-not (Test-Path 'Source\images')) {
    Write-Host 'Copying images from Playdate SDK...'
    Copy-Item -Recurse (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Sprite Game\Source\images') 'Source\'
}

pdgoc -sim -device -name="SpriteGame" -author="PdGo" -desc="Sprite Game Demo" -bundle-id=com.pdgo.spritegame -version=1.0 -build-number=1
