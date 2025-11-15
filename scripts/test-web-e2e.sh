#!/usr/bin/env bash

# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "[e2e] Starting DiliVet Web E2E tests..."
echo ""

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "[e2e] ERROR: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Playwright is installed
if [ ! -d "tests/e2e/node_modules" ]; then
    echo "[e2e] Installing Playwright dependencies..."
    cd tests/e2e
    if command -v pnpm >/dev/null 2>&1; then
        pnpm install
    elif command -v npm >/dev/null 2>&1; then
        npm install
    else
        echo "[e2e] ERROR: Neither pnpm nor npm found"
        exit 1
    fi
    cd "$REPO_ROOT"
fi

# Start Docker stack
echo "[e2e] Building and starting Docker stack..."
docker compose -f docker-compose.e2e.yml up -d --build

# Wait for health endpoint
echo "[e2e] Waiting for server to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f -s http://localhost:8080/api/health >/dev/null 2>&1; then
        echo "[e2e] Server is ready!"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "[e2e] Waiting for server... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "[e2e] ERROR: Server did not become ready in time"
    docker compose -f docker-compose.e2e.yml down
    exit 1
fi

# Run Playwright tests
echo ""
echo "[e2e] Running Playwright tests..."
cd tests/e2e

if command -v pnpm >/dev/null 2>&1; then
    pnpm exec playwright test
    TEST_EXIT_CODE=$?
elif command -v npm >/dev/null 2>&1; then
    npm exec playwright test
    TEST_EXIT_CODE=$?
else
    echo "[e2e] ERROR: Neither pnpm nor npm found"
    cd "$REPO_ROOT"
    docker compose -f docker-compose.e2e.yml down
    exit 1
fi

cd "$REPO_ROOT"

# Cleanup
echo ""
echo "[e2e] Stopping Docker stack..."
docker compose -f docker-compose.e2e.yml down

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "[e2e] All E2E tests passed!"
    exit 0
else
    echo ""
    echo "[e2e] E2E tests failed with exit code $TEST_EXIT_CODE"
    exit $TEST_EXIT_CODE
fi

