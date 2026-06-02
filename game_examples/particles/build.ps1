# Build script for Windows - all complexity is handled by pdgoc

# Copy snowflake images from Playdate SDK
if (-not (Test-Path 'Source\images')) {
    New-Item -ItemType Directory -Path 'Source\images' -Force | Out-Null
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\images\snowflake1.png') 'Source\images\'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\images\snowflake2.png') 'Source\images\'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\images\snowflake3.png') 'Source\images\'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\images\snowflake4.png') 'Source\images\'
}

# Copy font
if (-not (Test-Path 'Source\font')) {
    New-Item -ItemType Directory -Path 'Source\font' -Force | Out-Null
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\font\namco-1x-table-9-9.png') 'Source\font\'
    Copy-Item (Join-Path $env:PLAYDATE_SDK_PATH 'C_API\Examples\Particles\Source\font\namco-1x.fnt') 'Source\font\'
}

pdgoc -sim -device -name "Particles" -author "PdGo" -desc "Particle System Demo" -bundle-id com.pdgo.particles -version 1.0 -build-number 1
