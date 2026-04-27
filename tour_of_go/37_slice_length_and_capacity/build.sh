#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 37 Slice Len Cap" \
  -author="PdGo" \
  -desc="Tour Of Go 37 Slice Length and Capacity" \
  -bundle-id=com.pdgo.tourofgo-37 \
  -version=1.0 \
  -build-number=1
