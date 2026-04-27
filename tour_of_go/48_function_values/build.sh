#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="Tour Of Go 48 Function Values" \
  -author="PdGo" \
  -desc="Tour Of Go 48 Function Values" \
  -bundle-id=com.pdgo.tourofgo-48 \
  -version=1.0 \
  -build-number=1
