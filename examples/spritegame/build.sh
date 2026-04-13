#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

# Copy images from SDK's Sprite Game example if not already present
if [ ! -d "Source/images" ]; then
  echo "Copying images from Playdate SDK..."
  cp -r "$PLAYDATE_SDK_PATH/C_API/Examples/Sprite Game/Source/images" Source/
fi

pdgoc -sim -device \
  -name="SpriteGame" \
  -author="PdGo" \
  -desc="Sprite Game Demo" \
  -bundle-id=com.pdgo.spritegame \
  -version=1.0 \
  -build-number=1
