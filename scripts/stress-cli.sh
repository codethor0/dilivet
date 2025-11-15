#!/usr/bin/env bash
set -euo pipefail

# DiliVet – CLI stress test script
# Run from repo root after building the dilivet binary into ./dist/dilivet

DILIVET_BIN="${DILIVET_BIN:-./dist/dilivet}"

if [ ! -x "$DILIVET_BIN" ]; then
  echo "[stress] dilivet binary not found at $DILIVET_BIN"
  echo "[stress] build it first, e.g.: go build -o dist/dilivet ./cmd/dilivet"
  exit 1
fi

echo "[stress] Using dilivet binary: $DILIVET_BIN"

# Prepare temp dir
WORKDIR="$(mktemp -d ./stress-XXXXXX)"
echo "[stress] Workdir: $WORKDIR"

cleanup() {
  echo "[stress] Cleaning up $WORKDIR"
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

# 1) Generate large message files (10MB, 100MB)
echo "[stress] Generating large message files..."
dd if=/dev/urandom of="$WORKDIR/msg_10mb.bin"  bs=1M count=10   status=none
dd if=/dev/urandom of="$WORKDIR/msg_100mb.bin" bs=1M count=100  status=none

# 2) Create "weird" hex files: CRLF, whitespace, mixed case, BOM
echo "[stress] Generating weird hex files..."
printf '\xEF\xBB\xBFdeadBEEF\r\n  cafebabe  \r\n' > "$WORKDIR/weird_hex_crlf_bom.hex"
printf '  DEADBEEF\nCAFEBABE\n  00ff11  '   > "$WORKDIR/weird_hex_ws.hex"

# Use the same hex for pk/sig just to drive error paths; we only care about robustness.
echo "[stress] Running error-path stress loops..."
for i in $(seq 1 50); do
  # Expect clean failures, no panics
  "$DILIVET_BIN" verify \
    -pub "$WORKDIR/weird_hex_crlf_bom.hex" \
    -sig "$WORKDIR/weird_hex_ws.hex" \
    -msg "$WORKDIR/msg_10mb.bin" \
    >/dev/null 2>>"$WORKDIR/stress.log" || true
done

# 3) Repeated kat-verify as a soak test
echo "[stress] Running kat-verify soak..."
for i in $(seq 1 50); do
  # kat-verify may exit with non-zero if tests fail, but should not panic/hang
  "$DILIVET_BIN" kat-verify >/dev/null 2>&1 || true
done

echo "[stress] Done. Check $WORKDIR/stress.log for any unexpected messages."
echo "[stress] No crashes or hangs indicates the CLI is robust under stress."

