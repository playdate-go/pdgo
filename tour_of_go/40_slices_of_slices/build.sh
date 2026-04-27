#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 40 Slices of Slices" \
  -author="PdGo" \
  -desc="Tour Of Go 40 Slices of Slices" \
  -bundle-id=com.pdgo.tourofgo-40 \
  -version=1.0 \
  -build-number=1
