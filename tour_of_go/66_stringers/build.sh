#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 66 Stringers" \
  -author="PdGo" \
  -desc="Tour Of Go 66 Stringers" \
  -bundle-id=com.pdgo.tourofgo-66 \
  -version=1.0 \
  -build-number=1
