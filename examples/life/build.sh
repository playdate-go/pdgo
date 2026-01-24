#!/bin/bash

export CGO_CFLAGS="-I$PLAYDATE_SDK_PATH/C_API -DTARGET_EXTENSION=1"

go get github.com/playdate-go/pdgo@latest

pdgoc -device -sim \
  -name=Life \
  -author=PdGo \
  -desc="Conway's Game of Life" \
  -bundle-id=com.pdgo.life \
  -version=1.0 \
  -build-number=1
