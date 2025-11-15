# DiliVet v0.3.0 Release Notes

## Overview

DiliVet v0.3.0 introduces a **web-based user interface** for ML-DSA diagnostics, making the toolkit accessible through a browser. This release includes comprehensive testing infrastructure, Docker support, and production-ready deployment configurations.

---

## Major Features

### Web UI

- **Browser-based interface**: Access DiliVet diagnostics from any modern web browser
- **Three main pages**:
  - **Dashboard**: Server health and version information
  - **Verify Signature**: Interactive form for ML-DSA signature verification
  - **KAT Verification**: Run known-answer test verification with detailed results
- **Modern stack**: React 18 + TypeScript, built with Vite
- **Responsive design**: Works on desktop and mobile browsers

### Backend API

- **Go HTTP server**: Lightweight, efficient backend wrapping existing DiliVet library
- **Three endpoints**:
  - `GET /api/health` - Server health and version
  - `POST /api/verify` - Signature verification
  - `POST /api/kat-verify` - KAT verification
- **SPA routing support**: Properly handles React Router client-side routes
- **Error handling**: Structured JSON responses with clear error messages

### Testing Infrastructure

- **Backend tests**: 25+ tests covering edge cases, concurrency, large inputs, boundary conditions
- **Frontend tests**: 15+ component and API client tests
- **E2E tests**: 12 Playwright tests across Chromium, Firefox, and WebKit
- **Load tests**: k6 scripts for stress testing all endpoints
- **All tests passing**: Comprehensive coverage ensures reliability

### Docker Support

- **Multi-stage Dockerfile**: Optimized builds for production
- **docker-compose.e2e.yml**: Easy setup for E2E testing
- **Health checks**: Built-in container health monitoring

---

## Technical Details

### Architecture

- **Backend** (`web/server/`): Go HTTP server using standard library
- **Frontend** (`web/ui/`): React 18 + TypeScript SPA
- **Separation**: Web components live in `web/` directory, core CLI/library unchanged

### Testing

- **Unit/Integration**: Backend and frontend unit tests
- **E2E**: Playwright with Docker automation
- **Load**: k6-based stress testing
- **Scripts**: `check-web.sh`, `test-web-e2e.sh`, `test-web-load.sh`

### Documentation

- `docs/WEB_UI.md` - Complete user guide and API documentation
- `docs/WEB_TEST_REPORT.md` - Comprehensive test documentation
- `docs/WEB_TESTING_SUMMARY.md` - Quick reference guide

---

## Usage

### Quick Start

```bash
# Start backend server
go run ./web/server

# In another terminal, start frontend dev server
cd web/ui && npm install && npm run dev
```

### Production Build

```bash
# Build frontend
cd web/ui && npm install && npm run build

# Start server (serves built frontend)
go run ./web/server
```

### Docker

```bash
# Build and run
docker compose -f docker-compose.e2e.yml up --build
```

### Testing

```bash
# Quick check (unit + integration)
./scripts/check-web.sh

# E2E tests
./scripts/test-web-e2e.sh

# Load tests (requires running server)
./scripts/test-web-load.sh
```

---

## Breaking Changes

None. The web UI is additive and does not change existing CLI or library behavior.

---

## Security Notes

- **Diagnostics tooling**: Intended for controlled environments
- **No authentication**: Not hardened for untrusted multi-tenant use
- **Input validation**: All inputs validated on backend
- **Error messages**: Structured, user-friendly, no internal details leaked

See `docs/WEB_UI.md` for deployment recommendations.

---

## What's Next

- Security hardening review (in progress)
- Optional authentication layer
- Enhanced error messages and user guidance
- Additional diagnostic endpoints

---

## Full Changelog

### Added
- Web UI with React + TypeScript frontend
- Go HTTP backend server
- Comprehensive test suite (backend, frontend, E2E, load)
- Docker support (Dockerfile.web, docker-compose.e2e.yml)
- Testing scripts (check-web.sh, test-web-e2e.sh, test-web-load.sh)
- Documentation (WEB_UI.md, WEB_TEST_REPORT.md, WEB_TESTING_SUMMARY.md)

### Fixed
- SPA routing: Server now properly serves index.html for all non-API routes
- Test selectors: Made more specific to avoid false positives
- TypeScript errors: Fixed test file compatibility issues

### Changed
- CI: Added web component testing to standard CI workflow

---

## Contributors

- Thor Thor (codethor@gmail.com)

---

## Links

- **Repository**: https://github.com/codethor0/dilivet
- **Release**: https://github.com/codethor0/dilivet/releases/tag/v0.3.0
- **Documentation**: See `docs/WEB_UI.md`

---

**DiliVet v0.3.0** - Browser-based ML-DSA diagnostics, backed by comprehensive testing and the proven CLI core.

