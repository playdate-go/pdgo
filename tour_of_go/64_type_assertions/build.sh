#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 64 Type Assertions" \
  -author="PdGo" \
  -desc="Tour Of Go 64 Type Assertions" \
  -bundle-id=com.pdgo.tourofgo-64 \
  -version=1.0 \
  -build-number=1
