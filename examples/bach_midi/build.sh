#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

# Copy assets from SDK if not already present
if [ ! -f "Source/bach.mid" ]; then
  echo "Copying bach.mid from Playdate SDK..."
  cp "$PLAYDATE_SDK_PATH/C_API/Examples/bach.mid/Source/bach.mid" Source/
fi

if [ ! -f "Source/piano.wav" ]; then
  echo "Copying piano.wav from Playdate SDK..."
  cp "$PLAYDATE_SDK_PATH/C_API/Examples/bach.mid/Source/piano.wav" Source/
fi


pdgoc -sim -device \
  -name="BachMIDI" \
  -author="PdGo" \
  -desc="Bach MIDI Player" \
  -bundle-id=com.pdgo.bachmidi \
  -version=1.0 \
  -build-number=1
