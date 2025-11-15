#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "[load] Starting DiliVet Web load tests..."
echo "[load] Target: $BASE_URL"
echo ""

# Check if k6 is installed
if ! command -v k6 >/dev/null 2>&1; then
    echo "[load] ERROR: k6 is not installed"
    echo "[load] Install from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

# Check if server is reachable
if ! curl -f -s "$BASE_URL/api/health" >/dev/null 2>&1; then
    echo "[load] ERROR: Server is not reachable at $BASE_URL"
    echo "[load] Please start DiliVet Web first (docker compose -f docker-compose.e2e.yml up or go run ./web/server)"
    exit 1
fi

echo "[load] Server is reachable. Starting load tests..."
echo ""

# Run health endpoint load test
echo "[load] 1/3: Health endpoint load test"
cd "$REPO_ROOT/tests/load"
k6 run --env BASE_URL="$BASE_URL" health_load.js
echo ""

# Run verify endpoint load test
echo "[load] 2/3: Verify endpoint load test"
k6 run --env BASE_URL="$BASE_URL" verify_load.js
echo ""

# Run KAT endpoint load test (optional, takes longer)
if [ "${SKIP_KAT_LOAD:-}" != "1" ]; then
    echo "[load] 3/3: KAT endpoint load test (this may take a while)..."
    k6 run --env BASE_URL="$BASE_URL" kat_load.js
else
    echo "[load] 3/3: KAT endpoint load test skipped (set SKIP_KAT_LOAD=1 to skip)"
fi

echo ""
echo "[load] All load tests completed!"

