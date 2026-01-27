#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="SpriteGame" \
  -author="PdGo" \
  -desc="Sprite Game Demo" \
  -bundle-id=com.pdgo.spritegame \
  -version=1.0 \
  -build-number=1
