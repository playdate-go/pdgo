#!/bin/bash
# GC Bench - Per-frame GC pause benchmark

pdgoc -sim -device \
  -name="GCBench" \
  -author="PdGo" \
  -desc="GC Pause Benchmark" \
  -bundle-id=com.pdgo.gcbench \
  -version=1.0 \
  -build-number=1
