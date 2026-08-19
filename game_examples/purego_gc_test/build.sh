#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="PureGoGcTest" \
  -author="PdGo" \
  -desc="GC Test" \
  -bundle-id=com.pdgo.puregogctest \
  -version=1.0 \
  -build-number=1
