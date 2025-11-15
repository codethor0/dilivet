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
echo "  📸 DiliVet Web UI Screenshot Capture"
echo "═══════════════════════════════════════════════════════"
echo ""

# Cleanup function
cleanup() {
  echo ""
  echo "[screenshot] Cleaning up..."
  # Stop any background server processes
  pkill -f 'go run ./web/server' || true
  lsof -ti:8080 | xargs kill -9 2>/dev/null || true
  # Stop Docker container if running
  docker rm -f dilivet-web-container 2>/dev/null || true
  echo "[screenshot] Cleanup complete"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Step 1: Start lab profile server
echo "[screenshot] Step 1/3: Starting lab profile server..."
# Stop any existing server first
pkill -f 'go run ./web/server' || true
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
docker rm -f dilivet-web-container 2>/dev/null || true
sleep 1

# Build frontend if needed
if [ ! -d "./web/ui/dist" ]; then
  echo "[screenshot] Building frontend for static serving..."
  (cd web/ui && npm install --silent && npm run build --silent)
fi

# Start server using go run (more reliable than Docker for this use case)
echo "[screenshot] Starting server with lab profile..."
export MAX_BODY_SIZE=10485760
export REQUEST_TIMEOUT=30s
export ALLOWED_ORIGINS="http://localhost:8080,http://localhost:3000"
export REQUIRE_AUTH=true
export AUTH_TOKEN=$(openssl rand -hex 32)
echo "$AUTH_TOKEN" > /tmp/dilivet-auth-token.txt

PORT=8080 go run ./web/server >/tmp/dilivet-screenshot-server.log 2>&1 &
SERVER_PID=$!
echo "[screenshot] Server starting (PID: $SERVER_PID)"
echo ""

# Step 2: Wait for health endpoint
echo "[screenshot] Step 2/3: Waiting for server to be ready..."
MAX_WAIT=60
WAIT_COUNT=0
HEALTH_URL="http://localhost:8080/api/health"

# Load AUTH_TOKEN
AUTH_TOKEN="${AUTH_TOKEN:-$(cat /tmp/dilivet-auth-token.txt 2>/dev/null || echo '')}"

while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
  # Try with auth token
  if [ -n "$AUTH_TOKEN" ]; then
    if curl -s -f -H "Authorization: Bearer $AUTH_TOKEN" "$HEALTH_URL" >/dev/null 2>&1; then
      echo "[screenshot] ✅ Server is ready"
      break
    fi
  else
    # Try without auth (fallback)
    if curl -s -f "$HEALTH_URL" >/dev/null 2>&1; then
      echo "[screenshot] ✅ Server is ready"
      break
    fi
  fi
  WAIT_COUNT=$((WAIT_COUNT + 1))
  if [ $((WAIT_COUNT % 5)) -eq 0 ]; then
    echo "[screenshot] Still waiting... (${WAIT_COUNT}s/${MAX_WAIT}s)"
  fi
  sleep 1
done

if [ $WAIT_COUNT -ge $MAX_WAIT ]; then
  echo "[screenshot] ERROR: Server did not become ready within ${MAX_WAIT}s"
  echo "[screenshot] Server logs:"
  tail -20 /tmp/dilivet-screenshot-server.log || true
  exit 1
fi
echo ""

# Step 3: Run Playwright screenshot test
echo "[screenshot] Step 3/3: Capturing screenshot with Playwright..."
cd "$REPO_ROOT/tests/e2e"

# Ensure assets directory exists
mkdir -p "$REPO_ROOT/docs/assets"

# Run screenshot test
if BASE_URL="http://localhost:8080" npx playwright test tests/screenshot.spec.ts --project=chromium; then
  echo "[screenshot] ✅ Screenshot test passed"
else
  echo "[screenshot] ERROR: Screenshot test failed"
  exit 1
fi

# Verify screenshot exists
SCREENSHOT_PATH="$REPO_ROOT/docs/assets/dilivet-web-ui.png"
if [ -f "$SCREENSHOT_PATH" ]; then
  FILE_SIZE=$(stat -f%z "$SCREENSHOT_PATH" 2>/dev/null || stat -c%s "$SCREENSHOT_PATH" 2>/dev/null || echo "0")
  echo "[screenshot] ✅ Screenshot saved: $SCREENSHOT_PATH ($(numfmt --to=iec-i --suffix=B $FILE_SIZE 2>/dev/null || echo "${FILE_SIZE} bytes"))"
else
  echo "[screenshot] ERROR: Screenshot file not found at $SCREENSHOT_PATH"
  exit 1
fi

cd "$REPO_ROOT"

echo ""
echo "═══════════════════════════════════════════════════════"
echo "  ✅ Screenshot capture complete"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "Screenshot location: docs/assets/dilivet-web-ui.png"
echo ""

# Cleanup on success
cleanup
trap - EXIT

