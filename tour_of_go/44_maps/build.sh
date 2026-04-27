#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 44 Maps" \
  -author="PdGo" \
  -desc="Tour Of Go 44 Maps" \
  -bundle-id=com.pdgo.tourofgo-44 \
  -version=1.0 \
  -build-number=1
