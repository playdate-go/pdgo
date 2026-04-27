#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 45 Map Literals" \
  -author="PdGo" \
  -desc="Tour Of Go 45 Map Literals" \
  -bundle-id=com.pdgo.tourofgo-45 \
  -version=1.0 \
  -build-number=1
