#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 50 Methods" \
  -author="PdGo" \
  -desc="Tour Of Go 50 Methods" \
  -bundle-id=com.pdgo.tourofgo-50 \
  -version=1.0 \
  -build-number=1
