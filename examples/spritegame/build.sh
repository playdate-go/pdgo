#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

# Copy images from SDK's Sprite Game example if not already present
if [ ! -d "Source/images" ]; then
  echo "Copying images from Playdate SDK..."
  cp -r "$PLAYDATE_SDK_PATH/C_API/Examples/Sprite Game/Source/images" Source/
fi

pdgoc -device -sim \
  -name=SpriteGame \
  -author=PdGo \
  -desc="Shoot-em-up game" \
  -bundle-id=com.pdgo.spritegame \
  -version=1.0 \
  -build-number=1
