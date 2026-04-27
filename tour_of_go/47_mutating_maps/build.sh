#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 47 Mutating Maps" \
  -author="PdGo" \
  -desc="Tour Of Go 47 Mutating Maps" \
  -bundle-id=com.pdgo.tourofgo-47 \
  -version=1.0 \
  -build-number=1
