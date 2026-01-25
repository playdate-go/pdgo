#!/bin/bash

cd Source

go get github.com/playdate-go/pdgo@latest

cd ..

pdgoc -device -sim \
  -name=HelloWorld \
  -author=PdGo \
  -desc="HelloWorld Game" \
  -bundle-id=com.pdgo.hello_world \
  -version=1.0 \
  -build-number=1
