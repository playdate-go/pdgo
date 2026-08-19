#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="ReallocDebug" \
  -author="PdGo" \
  -desc="Realloc Debug Example" \
  -bundle-id=com.pdgo.reallocdebug \
  -version=1.0 \
  -build-number=1
