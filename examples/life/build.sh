#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

pdgoc -device -sim \
  -name=Life \
  -author=PdGo \
  -desc="Conway's Game of Life" \
  -bundle-id=com.pdgo.life \
  -version=1.0 \
  -build-number=1
