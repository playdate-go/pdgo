#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 67 Errors" \
  -author="PdGo" \
  -desc="Tour Of Go 67 Errors" \
  -bundle-id=com.pdgo.tourofgo-67 \
  -version=1.0 \
  -build-number=1
