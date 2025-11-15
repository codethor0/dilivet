#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

DILIVET_BIN="${DILIVET_BIN:-./dist/dilivet}"

if [ ! -x "$DILIVET_BIN" ]; then
  echo "[soak] dilivet binary not found at $DILIVET_BIN"
  echo "[soak] build it first, e.g.: go build -o dist/dilivet ./cmd/dilivet"
  exit 1
fi

LOG="./stress-soak.log"

echo "[soak] Starting soak run; logging to $LOG"
echo "[soak] Using dilivet binary: $DILIVET_BIN"
: > "$LOG"

for i in $(seq 1 1000); do
  "$DILIVET_BIN" kat-verify >>"$LOG" 2>&1 || {
    echo "[soak] iteration $i failed; see $LOG"
    exit 1
  }
  if [ $((i % 100)) -eq 0 ]; then
    echo "[soak] Completed $i iterations..."
  fi
done

echo "[soak] Completed 1000 iterations with no failures."

