#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 49 Closures" \
  -author="PdGo" \
  -desc="Tour Of Go 49 Function Closures" \
  -bundle-id=com.pdgo.tourofgo-49 \
  -version=1.0 \
  -build-number=1
