#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 20 Forever" \
  -author="PdGo" \
  -desc="Tour Of Go 20 Forever" \
  -bundle-id=com.pdgo.tourofgo-20 \
  -version=1.0 \
  -build-number=1