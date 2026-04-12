#!/bin/bash

# Copy snowflake images from Playdate SDK
if [ ! -d "Source/images" ]; then
    mkdir -p Source/images
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/images/snowflake1.png" Source/images/
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/images/snowflake2.png" Source/images/
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/images/snowflake3.png" Source/images/
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/images/snowflake4.png" Source/images/
fi
    
# Copy font
if [ ! -d "Source/font" ]; then
    mkdir -p Source/font
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/font/namco-1x-table-9-9.png" Source/font/
    cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Particles/Source/font/namco-1x.fnt" Source/font/
fi
    
pdgoc -sim -device \
  -name="Particles" \
  -author="PdGo" \
  -desc="Particle System Demo" \
  -bundle-id=com.pdgo.particles \
  -version=1.0 \
  -build-number=1
