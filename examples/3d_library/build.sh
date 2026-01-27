#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="3DLibrary" \
  -author="PdGo" \
  -desc="3D Library Demo" \
  -bundle-id=com.pdgo.3dlibrary \
  -version=1.0 \
  -build-number=1
