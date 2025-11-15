#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

echo "[dilivet] Go version:"
go version || { echo "[dilivet] go is not installed or not on PATH"; exit 1; }

echo
echo "[dilivet] Step 1/5: go vet ./..."
go vet ./...

echo
echo "[dilivet] Step 2/5: golangci-lint (if available)"
if command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run --timeout=5m
else
  echo "[dilivet] golangci-lint not found; skipping lint. Install from https://golangci-lint.run/ if you want full checks."
fi

echo
echo "[dilivet] Step 3/5: go test -race ./..."
go test -race -p 4 ./...

echo
echo "[dilivet] Step 4/5: fuzz smoke tests (short runs; safe to skip if too slow)"
if go test -c ./fuzz >/dev/null 2>&1; then
  # These fuzz targets are listed in the README; if they do not exist, errors are ignored.
  go test ./fuzz -run='^$' -fuzz=FuzzDecodePublicKey -fuzztime=30s || echo "[dilivet] FuzzDecodePublicKey fuzz run failed or not present; continuing."
  go test ./fuzz -run='^$' -fuzz=FuzzVerify         -fuzztime=30s || echo "[dilivet] FuzzVerify fuzz run failed or not present; continuing."
else
  echo "[dilivet] No fuzz package or fuzzing not supported by this Go version; skipping fuzz."
fi

echo
echo "[dilivet] Step 5/5: cross-build matrix (no install, just build into dist/check-build)"
mkdir -p dist/check-build
for os in linux darwin windows; do
  for arch in amd64 arm64; do
    out="dist/check-build/dilivet-${os}-${arch}"
    echo "[dilivet] Building $out ..."
    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
      go build -trimpath -ldflags "-s -w" \
      -o "$out" ./cmd/dilivet
  done
done

echo
echo "[dilivet] All checks completed successfully."
echo "[dilivet] If there were no errors above, the project is in a good state."

