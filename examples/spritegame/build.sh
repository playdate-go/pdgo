#!/bin/bash

PLAYDATE_SDK="$PLAYDATE_SDK_PATH"

# Set CGO flags for simulator build
export CGO_CFLAGS="-I$PLAYDATE_SDK/C_API -DTARGET_EXTENSION=1"

go get github.com/playdate-go/pdgo@latest

# Copy images from SDK's Sprite Game example if not already present
if [ ! -d "Source/images" ]; then
  echo "Copying images from Playdate SDK..."
  cp -r "$PLAYDATE_SDK/C_API/Examples/Sprite Game/Source/images" Source/
fi

pdgoc -device -sim \
  -name=SpriteGame \
  -author=PdGo \
  -desc="Shoot-em-up game" \
  -bundle-id=com.pdgo.spritegame \
  -version=1.0 \
  -build-number=1
