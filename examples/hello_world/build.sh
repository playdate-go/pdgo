#!/bin/bash

export CGO_CFLAGS="-I$PLAYDATE_SDK_PATH/C_API -DTARGET_EXTENSION=1"

pdgoc -device -sim \
  -name=HelloWorld \
  -author=PdGo \
  -desc="HelloWorld Game" \
  -bundle-id=com.pdgo.hello_world \
  -version=1.0 \
  -build-number=1
