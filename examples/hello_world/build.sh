#!/bin/bash

# Set Playdate SDK path for CGO (adjust if your SDK is in a different location)
export CGO_CFLAGS="-I$HOME/Developer/PlaydateSDK/C_API -DTARGET_EXTENSION=1"

pdgoc -device -sim \
  -name=HelloWorld \
  -author=PdGo \
  -desc="HelloWorld Game" \
  -bundle-id=com.pdgo.hello_world \
  -version=1.0 \
  -build-number=1
