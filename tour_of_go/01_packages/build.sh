#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 01 Packages" \
  -author="PdGo" \
  -desc="Tour Of Go 01 Packages" \
  -bundle-id=com.pdgo.tourofgo-01 \
  -version=1.0 \
  -build-number=1