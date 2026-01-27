#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="GoLogo" \
  -author="PdGo" \
  -desc="Go Logo Demo" \
  -bundle-id=com.pdgo.gologo \
  -version=1.0 \
  -build-number=1
