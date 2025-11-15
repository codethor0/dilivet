#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "═══════════════════════════════════════════════════════"
echo "  🔍 DiliVet Release Smoke Test"
echo "═══════════════════════════════════════════════════════"
echo ""

CORE_STATUS="FAIL"
WEB_STATUS="FAIL"
SCREENSHOT_STATUS="FAIL"

# Core checks
echo "[smoke] Running core checks..."
if ./scripts/check-all.sh >/tmp/dilivet-smoke-core.log 2>&1; then
  CORE_STATUS="PASS"
  echo "[smoke] ✅ Core checks: PASS"
else
  echo "[smoke] ❌ Core checks: FAIL"
  echo "[smoke] See /tmp/dilivet-smoke-core.log for details"
fi
echo ""

# Web checks
echo "[smoke] Running web checks..."
if ./scripts/check-web.sh >/tmp/dilivet-smoke-web.log 2>&1; then
  WEB_STATUS="PASS"
  echo "[smoke] ✅ Web checks: PASS"
else
  echo "[smoke] ❌ Web checks: FAIL"
  echo "[smoke] See /tmp/dilivet-smoke-web.log for details"
fi
echo ""

# Screenshot capture
echo "[smoke] Capturing Web UI screenshot..."
if ./scripts/capture-web-screenshot.sh >/tmp/dilivet-smoke-screenshot.log 2>&1; then
  SCREENSHOT_PATH="$REPO_ROOT/docs/assets/dilivet-web-ui.png"
  if [ -f "$SCREENSHOT_PATH" ]; then
    # Verify PNG magic bytes
    PNG_MAGIC=$(head -c 4 "$SCREENSHOT_PATH" | hexdump -C | head -1)
    if echo "$PNG_MAGIC" | grep -qE "89 50 4e 47|89504e47"; then
      SCREENSHOT_STATUS="PASS (magic bytes OK)"
      echo "[smoke] ✅ Screenshot: PASS (magic bytes OK)"
    else
      SCREENSHOT_STATUS="FAIL (NOT OK - not a PNG)"
      echo "[smoke] ❌ Screenshot: FAIL (NOT OK - not a PNG)"
      echo "[smoke] First 4 bytes: $(head -c 4 "$SCREENSHOT_PATH" | hexdump -C | head -1)"
    fi
  else
    SCREENSHOT_STATUS="FAIL (file not found)"
    echo "[smoke] ❌ Screenshot: FAIL (file not found)"
  fi
else
  echo "[smoke] ❌ Screenshot capture: FAIL"
  echo "[smoke] See /tmp/dilivet-smoke-screenshot.log for details"
fi
echo ""

# Summary
echo "═══════════════════════════════════════════════════════"
echo "  📊 Smoke Test Summary"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "core:      $CORE_STATUS"
echo "web:       $WEB_STATUS"
echo "screenshot: $SCREENSHOT_STATUS"
echo ""

# Exit with error if any check failed
if [ "$CORE_STATUS" != "PASS" ] || [ "$WEB_STATUS" != "PASS" ] || [ "$SCREENSHOT_STATUS" != "PASS (magic bytes OK)" ]; then
  echo "❌ Smoke test failed - see logs above"
  exit 1
fi

echo "✅ All smoke tests passed!"
exit 0

