#!/bin/bash
# Build script for Networking Demo
# Note: Network is only available on Playdate Simulator, not on device

pdgoc -sim \
  -name="Networking" \
  -author="PdGo" \
  -desc="Networking Demo" \
  -bundle-id=com.pdgo.networking \
  -version=0.1 \
  -build-number=1
