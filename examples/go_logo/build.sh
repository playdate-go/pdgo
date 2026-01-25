#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

pdgoc -device -sim \
  -name=GoLogo \
  -author=PdGo \
  -desc="GoLogo Game" \
  -bundle-id=com.pdgo.go_logo \
  -version=1.0 \
  -build-number=1
