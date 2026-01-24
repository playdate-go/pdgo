#!/bin/bash

export CGO_CFLAGS="-I$PLAYDATE_SDK_PATH/C_API -DTARGET_EXTENSION=1"

go get github.com/playdate-go/pdgo@latest

pdgoc -device -sim \
  -name=GoLogo \
  -author=PdGo \
  -desc="GoLogo Game" \
  -bundle-id=com.pdgo.go_logo \
  -version=1.0 \
  -build-number=1
