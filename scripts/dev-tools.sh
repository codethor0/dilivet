#!/usr/bin/env bash
set -euo pipefail
echo "Installing local dev tools (if missing)…"
if command -v brew >/dev/null 2>&1; then
  brew install golangci-lint staticcheck || true
else
  echo "Homebrew not found; skipping brew installs."
fi
# Always install/refresh govulncheck from x/vuln
go install golang.org/x/vuln/cmd/govulncheck@latest
echo "✅ Tools ready. Make sure GOPATH/bin is on your PATH: $(go env GOPATH)/bin"
