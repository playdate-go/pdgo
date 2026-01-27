#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="HelloWorld" \
  -author="PdGo" \
  -desc="Hello World Demo" \
  -bundle-id=com.pdgo.helloworld \
  -version=1.0 \
  -build-number=1
