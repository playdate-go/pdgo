#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 39 Make Slice" \
  -author="PdGo" \
  -desc="Tour Of Go 39 Creating Slice with Make" \
  -bundle-id=com.pdgo.tourofgo-39 \
  -version=1.0 \
  -build-number=1
