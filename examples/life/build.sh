#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="GameOfLife" \
  -author="PdGo" \
  -desc="Conway's Game of Life" \
  -bundle-id=com.pdgo.life \
  -version=1.0 \
  -build-number=1
