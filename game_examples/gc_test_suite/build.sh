#!/bin/bash
# GC Test Suite - Build Script
# Tests the device GC with pure Go code (no direct C calls in test code)

pdgoc -sim -device \
  -name="GCTestSuite" \
  -author="PdGo" \
  -desc="GC Test Suite - Tests memory management with native Go constructs" \
  -bundle-id=com.pdgo.gctestsuite \
  -version=1.0 \
  -build-number=1
