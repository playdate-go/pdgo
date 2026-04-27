#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 41 Append Slice" \
  -author="PdGo" \
  -desc="Tour Of Go 41 Appending to Slice" \
  -bundle-id=com.pdgo.tourofgo-41 \
  -version=1.0 \
  -build-number=1
