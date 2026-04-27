#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 33 Slices" \
  -author="PdGo" \
  -desc="Tour Of Go 33 Slices" \
  -bundle-id=com.pdgo.tourofgo-33 \
  -version=1.0 \
  -build-number=1
