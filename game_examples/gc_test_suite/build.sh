#!/bin/bash
# GC Test Pure Go - Build Script
# Tests Go's garbage collector with pure Go code (no direct C calls)

pdgoc -sim -device \
  -name="GCTestPureGo" \
  -author="PdGo" \
  -desc="GC Test Pure Go - Tests memory management with native Go constructs" \
  -bundle-id=com.pdgo.gctestpurego \
  -version=1.0 \
  -build-number=1
