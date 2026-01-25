#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

pdgoc -device -sim \
  -name=BouncingSquare \
  -author=PdGo \
  -desc="BouncingSquare Game" \
  -bundle-id=com.pdgo.bouncingsquare \
  -version=1.0 \
  -build-number=1
