#!/usr/bin/env bash

# DiliVet – ML-DSA diagnostics and vetting toolkit
# Web check script: runs backend and frontend tests and builds

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "[web] Checking DiliVet Web components..."
echo ""

# Check backend
echo "[web] Step 1/3: Backend tests (web/server)"
cd "$REPO_ROOT/web/server"
if ! go test -v ./...; then
    echo "[web] ERROR: Backend tests failed"
    exit 1
fi
echo ""

# Check frontend dependencies
echo "[web] Step 2/3: Frontend dependencies and tests (web/ui)"
cd "$REPO_ROOT/web/ui"
if [ ! -d "node_modules" ]; then
    echo "[web] Installing frontend dependencies..."
    if command -v pnpm >/dev/null 2>&1; then
        pnpm install
    elif command -v npm >/dev/null 2>&1; then
        npm install
    else
        echo "[web] ERROR: Neither pnpm nor npm found"
        exit 1
    fi
fi

if ! npm test -- --run; then
    echo "[web] ERROR: Frontend tests failed"
    exit 1
fi
echo ""

# Build frontend
echo "[web] Step 3/3: Frontend build (web/ui)"
if ! npm run build; then
    echo "[web] ERROR: Frontend build failed"
    exit 1
fi
echo ""

echo ""
echo "[web] Summary:"
echo "  - Backend tests: PASSED"
echo "  - Frontend tests: PASSED"
echo "  - Frontend build: PASSED"
echo ""
echo "[web] All web checks passed successfully!"
echo ""
echo "[web] Next steps:"
echo "  - Run E2E tests: ./scripts/test-web-e2e.sh"
echo "  - Run load tests: ./scripts/test-web-load.sh (requires running server)"

