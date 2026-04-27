#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 51 Methods Funcs" \
  -author="PdGo" \
  -desc="Tour Of Go 51 Methods are Functions" \
  -bundle-id=com.pdgo.tourofgo-51 \
  -version=1.0 \
  -build-number=1
