#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="BouncingSquare" \
  -author="PdGo" \
  -desc="Bouncing Square Demo" \
  -bundle-id=com.pdgo.bouncingsquare \
  -version=1.0 \
  -build-number=1
