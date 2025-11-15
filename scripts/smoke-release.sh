#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

# Smoke release script: runs all critical checks before release
# Exits with non-zero status if any step fails

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "═══════════════════════════════════════════════════════"
echo "  DiliVet Release Smoke Test"
echo "═══════════════════════════════════════════════════════"
echo ""

CORE_STATUS="FAIL"
WEB_STATUS="FAIL"
AUTH_STATUS="SKIP"

# Core gate
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Core gate: ./scripts/check-all.sh"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if ./scripts/check-all.sh >/tmp/dilivet-smoke-core.log 2>&1; then
    CORE_STATUS="PASS"
    echo "✅ Core: PASS"
else
    echo "❌ Core: FAIL"
    echo "Last 20 lines of output:"
    tail -20 /tmp/dilivet-smoke-core.log
    exit 1
fi

echo ""

# Web gate
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Web gate: ./scripts/check-web.sh"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f "./scripts/check-web.sh" ]; then
    if ./scripts/check-web.sh >/tmp/dilivet-smoke-web.log 2>&1; then
        WEB_STATUS="PASS"
        echo "✅ Web: PASS"
    else
        echo "❌ Web: FAIL"
        echo "Last 20 lines of output:"
        tail -20 /tmp/dilivet-smoke-web.log
        exit 1
    fi
else
    echo "⚠️  check-web.sh not found, skipping web gate"
    WEB_STATUS="SKIP"
fi

echo ""

# Auth gate (if script exists)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Auth gate: ./scripts/test-web-auth.sh"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f "./scripts/test-web-auth.sh" ]; then
    if ./scripts/test-web-auth.sh >/tmp/dilivet-smoke-auth.log 2>&1; then
        AUTH_STATUS="PASS"
        echo "✅ Auth: PASS"
    else
        echo "❌ Auth: FAIL"
        echo "Last 20 lines of output:"
        tail -20 /tmp/dilivet-smoke-auth.log
        exit 1
    fi
else
    echo "⚠️  test-web-auth.sh not found, skipping auth gate"
    AUTH_STATUS="SKIP"
fi

echo ""
echo "═══════════════════════════════════════════════════════"
echo "  Summary"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "  Core: $CORE_STATUS"
echo "  Web:  $WEB_STATUS"
echo "  Auth: $AUTH_STATUS"
echo ""

if [ "$CORE_STATUS" = "PASS" ] && [ "$WEB_STATUS" != "FAIL" ]; then
    echo "✅ All critical checks passed"
    exit 0
else
    echo "❌ One or more checks failed"
    exit 1
fi
