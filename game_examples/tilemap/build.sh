#!/bin/bash

# Copy font bitmap table asset from Playdate SDK
cp "${PLAYDATE_SDK_PATH}/C_API/Examples/Tilemap/Source/font-table-20-20.png" Source/

pdgoc -sim -device \
  -name="Tilemap" \
  -author="PdGo" \
  -desc="Tilemap Demo" \
  -bundle-id=com.pdgo.tilemap \
  -version=1.0 \
  -build-number=1
