#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

# Copy assets from SDK if not already present
if [ ! -f "Source/bach.mid" ]; then
  echo "Copying bach.mid from Playdate SDK..."
  cp "$PLAYDATE_SDK_PATH/C_API/Examples/bach.mid/Source/bach.mid" Source/
fi

if [ ! -f "Source/piano.wav" ]; then
  echo "Copying piano.wav from Playdate SDK..."
  cp "$PLAYDATE_SDK_PATH/C_API/Examples/bach.mid/Source/piano.wav" Source/
fi

pdgoc -sim \
  -name=BachMidi \
  -author=PdGo \
  -desc="Bach MIDI player with note visualization" \
  -bundle-id=com.pdgo.bachmidi \
  -version=1.0 \
  -build-number=1
