#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 68 Readers" \
  -author="PdGo" \
  -desc="Tour Of Go 68 Readers" \
  -bundle-id=com.pdgo.tourofgo-68 \
  -version=1.0 \
  -build-number=1
