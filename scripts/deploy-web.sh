#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

PROFILE="${1:-local}"

case "$PROFILE" in
  local)
    echo "=== Deploying DiliVet Web: Local Development Profile ==="
    export MAX_BODY_SIZE=10485760
    export REQUEST_TIMEOUT=30s
    export ALLOWED_ORIGINS="*"
    export REQUIRE_AUTH=false
    ;;
  lab)
    echo "=== Deploying DiliVet Web: Internal Lab Profile ==="
    export MAX_BODY_SIZE=10485760
    export REQUEST_TIMEOUT=30s
    export ALLOWED_ORIGINS="${ALLOWED_ORIGINS:-https://dilivet.internal.example.com}"
    export REQUIRE_AUTH=true
    if [ -z "${AUTH_TOKEN:-}" ]; then
      echo "⚠️  AUTH_TOKEN not set. Generating one..."
      export AUTH_TOKEN=$(openssl rand -hex 32)
      echo "✅ Generated AUTH_TOKEN (save this): $AUTH_TOKEN"
    fi
    ;;
  hardened)
    echo "=== Deploying DiliVet Web: Hardened Internal Profile ==="
    export MAX_BODY_SIZE=5242880
    export REQUEST_TIMEOUT=20s
    export ALLOWED_ORIGINS="${ALLOWED_ORIGINS:-https://dilivet.internal.example.com}"
    export REQUIRE_AUTH=true
    if [ -z "${AUTH_TOKEN:-}" ]; then
      echo "⚠️  AUTH_TOKEN not set. Generating one..."
      export AUTH_TOKEN=$(openssl rand -hex 32)
      echo "✅ Generated AUTH_TOKEN (save this): $AUTH_TOKEN"
    fi
    ;;
  *)
    echo "Usage: $0 [local|lab|hardened]"
    exit 1
    ;;
esac

echo ""
echo "Configuration:"
echo "  MAX_BODY_SIZE: $MAX_BODY_SIZE"
echo "  REQUEST_TIMEOUT: $REQUEST_TIMEOUT"
echo "  ALLOWED_ORIGINS: $ALLOWED_ORIGINS"
echo "  REQUIRE_AUTH: $REQUIRE_AUTH"
if [ "$REQUIRE_AUTH" = "true" ]; then
  echo "  AUTH_TOKEN: ${AUTH_TOKEN:0:8}... (hidden)"
fi
echo ""

# Check if Docker is available
if command -v docker >/dev/null 2>&1; then
  echo "🐳 Building Docker image..."
  docker build -f Dockerfile.web -t dilivet-web:v0.3.0 .
  
  echo ""
  echo "🚀 Starting container..."
  docker run --rm -p 8080:8080 \
    -e MAX_BODY_SIZE="$MAX_BODY_SIZE" \
    -e REQUEST_TIMEOUT="$REQUEST_TIMEOUT" \
    -e ALLOWED_ORIGINS="$ALLOWED_ORIGINS" \
    -e REQUIRE_AUTH="$REQUIRE_AUTH" \
    ${AUTH_TOKEN:+-e AUTH_TOKEN="$AUTH_TOKEN"} \
    dilivet-web:v0.3.0
else
  echo "📦 Running with Go (Docker not available)..."
  echo ""
  echo "Make sure frontend is built:"
  echo "  cd web/ui && npm install && npm run build"
  echo ""
  echo "Starting server..."
  go run ./web/server
fi

