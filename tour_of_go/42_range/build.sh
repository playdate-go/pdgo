#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 42 Range" \
  -author="PdGo" \
  -desc="Tour Of Go 42 Range" \
  -bundle-id=com.pdgo.tourofgo-42 \
  -version=1.0 \
  -build-number=1
