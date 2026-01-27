#!/bin/bash
# Minimal build script - all complexity is handled by pdgoc

pdgoc -sim -device \
  -name="BachMIDI" \
  -author="PdGo" \
  -desc="Bach MIDI Player" \
  -bundle-id=com.pdgo.bachmidi \
  -version=1.0 \
  -build-number=1
