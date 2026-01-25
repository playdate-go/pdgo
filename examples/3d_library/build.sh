#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

pdgoc -device -sim \
  -name="3D Library" \
  -author=PdGo \
  -desc="3D rendering demo with rotating icosahedron" \
  -bundle-id=com.pdgo.3ddemo \
  -version=1.0 \
  -build-number=1
