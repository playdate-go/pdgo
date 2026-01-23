#!/bin/bash

# Set Playdate SDK path for CGO (adjust if your SDK is in a different location)
export CGO_CFLAGS="-I$HOME/Developer/PlaydateSDK/C_API -DTARGET_EXTENSION=1"

pdgoc -device -sim \
  -name=Life \
  -author=PdGo \
  -desc="Conway's Game of Life" \
  -bundle-id=com.pdgo.life \
  -version=1.0 \
  -build-number=1
