# DiliVet Web Status

**Date:** 2025-11-15  
**Version:** v0.2.4  
**Status:** Implemented and tested

## Overview

DiliVet Web provides a browser-based interface for running ML-DSA diagnostics. It consists of a Go HTTP backend that wraps the existing DiliVet library and a React + TypeScript frontend.

## Architecture

### Backend (`web/server/`)

- **Language:** Go (standard library only)
- **Port:** 8080 (configurable via `PORT` env var)
- **Endpoints:**
  - `GET /api/health` - Server health and version
  - `POST /api/verify` - Signature verification
  - `POST /api/kat-verify` - KAT verification
- **Static Files:** Serves built frontend from `web/ui/dist` (if present)

### Frontend (`web/ui/`)

- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Testing:** Vitest + React Testing Library
- **Pages:**
  - Dashboard - Overview and server health
  - Verify Signature - Form for signature verification
  - KAT Verification - Run and view KAT results

## Test Coverage

### Backend

- Unit tests for all HTTP handlers
- Tests for error cases (invalid input, malformed hex, missing fields)
- Integration-style tests using `httptest`

**Status:** All tests passing

### Frontend

- Component tests for Verify form
- API client tests (mocked)
- Basic rendering tests

**Status:** Tests implemented

## CI Integration

The web components are integrated into `.github/workflows/ci.yml`:

- Backend tests run on Linux (`ubuntu-latest`)
- Frontend dependencies installed via npm
- Frontend tests run
- Frontend build verified

**Status:** Integrated into CI

## Known Limitations

1. **No Authentication:** The web UI does not include authentication or authorization. Suitable for controlled environments only.

2. **Not Multi-Tenant Ready:** Not hardened for untrusted users or multi-tenant deployments.

3. **Error Messages:** Error messages are displayed to users. Ensure no sensitive information is leaked.

4. **Performance:** Large inputs may impact performance. Consider rate limiting for production use.

5. **Static Serving:** In production, the Go server serves static files. For high-traffic scenarios, consider a dedicated web server.

## Endpoints Provided

### Health Check

- **Endpoint:** `GET /api/health`
- **Purpose:** Server health and version information
- **Response:** `{ "status": "ok", "version": "..." }`

### Signature Verification

- **Endpoint:** `POST /api/verify`
- **Purpose:** Verify ML-DSA signatures
- **Request:** Parameter set, public key (hex), signature (hex), message (hex or UTF-8)
- **Response:** Verification result (valid/invalid) or error

### KAT Verification

- **Endpoint:** `POST /api/kat-verify`
- **Purpose:** Run known-answer test verification
- **Request:** Optional vectors path (defaults to bundled vectors)
- **Response:** Summary with total vectors, passed/failed counts, and detailed results

## Next Steps

Potential enhancements for future versions:

- Additional diagnostic endpoints
- Real-time progress updates for long-running operations
- Export results as JSON/CSV
- Support for custom KAT vector uploads
- Enhanced error messages with suggestions
- Dark mode theme
- Authentication/authorization (if needed for production use)

## Verification

To verify the web implementation:

```bash
# Backend tests
go test ./web/server/...

# Frontend tests
cd web/ui && npm test

# Full check
./scripts/check-web.sh
```

All checks should pass.

