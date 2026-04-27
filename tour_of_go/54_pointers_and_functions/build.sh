#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 54 Pointers Funcs" \
  -author="PdGo" \
  -desc="Tour Of Go 54 Pointers and Functions" \
  -bundle-id=com.pdgo.tourofgo-54 \
  -version=1.0 \
  -build-number=1
