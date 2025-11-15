# DiliVet Web UI

DiliVet Web provides a browser-based interface for running ML-DSA diagnostics. It includes signature verification and known-answer test (KAT) verification capabilities.

## Architecture

The web version consists of:

- **Backend** (`web/server/`): Go HTTP server that wraps the existing DiliVet library
- **Frontend** (`web/ui/`): React + TypeScript single-page application (SPA)

The backend exposes JSON APIs that the frontend consumes. The server can also serve the built frontend static files.

## Quick Smoke Test

For a quick verification that the Web UI is working, run:

```bash
./scripts/smoke-web.sh
```

This script:
- Runs lightweight web tests (backend + frontend)
- Starts the lab profile server
- Waits for the health endpoint to be ready
- Prints a summary with the server URL

The server will remain running after the smoke test completes, so you can access it at `http://localhost:8080`.

## Running Locally

### Option 1: Development Mode (Recommended for Development)

**Terminal 1 - Backend:**
```bash
go run ./web/server
```

The server starts on `http://localhost:8080` by default. Set the `PORT` environment variable to use a different port.

**Terminal 2 - Frontend:**
```bash
cd web/ui
npm install  # or: pnpm install
npm run dev
```

The frontend dev server starts on `http://localhost:3000` and proxies API requests to the backend.

### Option 2: Production Build

**Build the frontend:**
```bash
cd web/ui
npm install
npm run build
```

**Start the backend (serves static files):**
```bash
go run ./web/server
```

The server will serve the built frontend from `web/ui/dist` and handle API requests.

## API Endpoints

### `GET /api/health`

Returns server health and version information.

**Response:**
```json
{
  "status": "ok",
  "version": "0.2.4"
}
```

### `POST /api/verify`

Verifies an ML-DSA signature.

**Request:**
```json
{
  "paramSet": "ML-DSA-44",
  "publicKeyHex": "deadbeef...",
  "signatureHex": "cafebabe...",
  "messageHex": "616263...",
  "message": "optional UTF-8 text (if messageHex not provided)"
}
```

**Response (success):**
```json
{
  "ok": true,
  "result": "valid"
}
```

**Response (invalid signature):**
```json
{
  "ok": true,
  "result": "invalid"
}
```

**Response (error):**
```json
{
  "ok": false,
  "error": "Invalid hex in public key: ..."
}
```

### `POST /api/kat-verify`

Runs KAT verification against ACVP sigVer vectors.

**Request (optional):**
```json
{
  "vectorsPath": "code/clean/testdata/kats/ml-dsa/ML-DSA-sigVer-FIPS204-internalProjection.json"
}
```

If `vectorsPath` is omitted, the default path is used.

**Response:**
```json
{
  "ok": true,
  "totalVectors": 1234,
  "passed": 1234,
  "failed": 0,
  "decodeFailures": 0,
  "details": [
    {
      "caseId": 1,
      "passed": true,
      "parameterSet": "ML-DSA-44",
      "reason": null
    }
  ]
}
```

## Testing

### Backend Tests

```bash
go test ./web/server/...
```

### Frontend Tests

```bash
cd web/ui
npm test
```

### Full Web Check

Run the comprehensive check script:

```bash
./scripts/check-web.sh
```

This script:
1. Runs backend tests
2. Installs frontend dependencies (if needed)
3. Runs frontend tests
4. Builds the frontend

## Development

### Backend

The backend is a standard Go HTTP server using only the standard library. It reuses the existing DiliVet library functions from `code/clean/mldsa` and `code/clean/kats`.

Key files:
- `web/server/main.go` - HTTP handlers and server setup
- `web/server/server_test.go` - Handler tests

### Frontend

The frontend uses:
- **React 18** with TypeScript
- **Vite** for build tooling
- **React Router** for navigation
- **Vitest** for testing

Key directories:
- `web/ui/src/pages/` - Page components (Dashboard, Verify, KATVerify)
- `web/ui/src/api/` - API client functions
- `web/ui/src/components/` - Reusable components (currently minimal)

## Configuration

### Backend

- `PORT` environment variable: Server port (default: `8080`)

### Frontend

- `VITE_API_BASE` environment variable: API base URL (default: `/api`)

For development, the Vite dev server proxies `/api` requests to `http://localhost:8080` (configured in `vite.config.ts`).

## Limitations and Security Notes

1. **Not for Production Multi-Tenant Use**: The web UI is designed for diagnostics and controlled environments. It is not hardened for untrusted users or multi-tenant deployments.

2. **No Authentication**: The web UI does not include authentication or authorization. If exposing publicly, add appropriate security measures.

3. **Error Handling**: The backend returns structured errors, but the frontend displays them to users. Ensure sensitive information is not leaked in error messages.

4. **Input Validation**: All inputs are validated, but large inputs may impact performance. Consider adding rate limiting for production use.

5. **Static File Serving**: In production mode, the Go server serves static files from `web/ui/dist`. For high-traffic scenarios, consider using a dedicated web server (nginx, Caddy) for static files.

## CI Integration

The web components are tested as part of the main CI workflow (`.github/workflows/ci.yml`):

- Backend tests run on Linux
- Frontend dependencies are installed
- Frontend tests run
- Frontend is built

All web checks are integrated into the existing CI pipeline without breaking existing checks.

## Security & Deployment

### Threat Model

DiliVet Web is designed for **internal, controlled environment** deployment. It is **not hardened for internet-exposed, multi-tenant, or untrusted user environments** without additional security layers.

**Assets:**
- Public keys, messages, signatures (diagnostic data)
- Server resources (CPU, memory)

**Assumptions:**
- Single-tenant use (all users share the same instance)
- Internal network deployment
- No account management or user isolation

**Out of Scope:**
- Multi-tenant isolation
- Rate limiting (should be handled by reverse proxy)
- TLS termination (should be handled by reverse proxy)
- Full audit logging

### Security Features

**Backend Hardening:**
- Request size limits (10MB default, configurable)
- Per-request timeouts (30s default, configurable)
- CORS policy (strict by default, configurable)
- Optional token authentication (disabled by default)
- Input validation (parameter sets, hex encoding, lengths)
- Sanitized logging (no sensitive data in logs)
- HTTP method enforcement

**Docker Hardening:**
- Non-root user (UID 1000)
- Minimal base image (Alpine)
- Health checks

### Configuration

Environment variables for security:

| Variable | Default | Description |
|----------|---------|-------------|
| `REQUIRE_AUTH` | `false` | Enable token authentication |
| `AUTH_TOKEN` | (none) | Shared token for authentication |
| `ALLOWED_ORIGINS` | (none) | Comma-separated CORS origins |
| `MAX_BODY_SIZE` | `10485760` | Max request body size (bytes) |
| `REQUEST_TIMEOUT` | `30s` | Per-request timeout |

### Recommended Deployment

**Option 1: Behind Reverse Proxy (Best Practice)**

Deploy behind Nginx/Envoy with:
- TLS termination
- Client certificate authentication (mTLS)
- Rate limiting
- IP allowlisting

Example Nginx configuration:
```nginx
server {
    listen 443 ssl;
    server_name dilivet.internal.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    ssl_client_certificate /path/to/ca.pem;
    ssl_verify_client on;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Option 2: Internal Network with Token Auth**

```bash
export REQUIRE_AUTH=true
export AUTH_TOKEN=$(openssl rand -hex 32)
export ALLOWED_ORIGINS=https://internal.example.com
go run ./web/server
```

**Option 3: Docker with Security Settings**

```yaml
# docker-compose.prod.yml
services:
  dilivet-web:
    build:
      context: .
      dockerfile: Dockerfile.web
    ports:
      - "127.0.0.1:8080:8080"  # Bind to localhost only
    environment:
      - REQUIRE_AUTH=true
      - AUTH_TOKEN=${AUTH_TOKEN}
      - ALLOWED_ORIGINS=https://internal.example.com
    user: "1000:1000"
    read_only: true
    security_opt:
      - no-new-privileges:true
```

### Security Checklist

**Pre-Deployment:**
- [ ] Set `REQUIRE_AUTH=true` if using token auth
- [ ] Generate strong `AUTH_TOKEN` (32+ random bytes)
- [ ] Configure `ALLOWED_ORIGINS` for CORS
- [ ] Deploy behind reverse proxy with TLS
- [ ] Set resource limits in Docker
- [ ] Review and restrict network access

**Post-Deployment:**
- [ ] Verify health endpoint responds
- [ ] Test authentication (if enabled)
- [ ] Verify CORS policy
- [ ] Check logs for sensitive data leakage
- [ ] Monitor resource usage

For detailed security information, see `docs/WEB_SECURITY_REVIEW.md`.

## Testing and QA

DiliVet Web includes a comprehensive test suite covering unit tests, integration tests, end-to-end (E2E) tests, and load testing.

### Quick Check

Run all unit and integration tests:

```bash
./scripts/check-web.sh
```

This script:
1. Runs backend tests (`go test ./web/server/...`)
2. Runs frontend tests (`npm test` in `web/ui`)
3. Builds the frontend (`npm run build` in `web/ui`)

### End-to-End Tests

Run E2E tests with Docker:

```bash
./scripts/test-web-e2e.sh
```

**Prerequisites:**
- Docker installed and running
- Node.js 18+ (for Playwright)

This script:
1. Builds and starts the Docker stack
2. Waits for server health
3. Runs Playwright tests across Chromium, Firefox, and WebKit
4. Cleans up the Docker stack

### Load Testing

Run load tests (requires a running server):

```bash
# Start server (choose one):
docker compose -f docker-compose.e2e.yml up -d
# OR
go run ./web/server

# Run load tests
./scripts/test-web-load.sh
```

**Prerequisites:**
- k6 installed (https://k6.io/docs/getting-started/installation/)
- Server running on `http://localhost:8080`

Load tests cover:
- Health endpoint (50-100 concurrent users)
- Verify endpoint (10-50 concurrent users)
- KAT endpoint (2-5 concurrent users, slower)

**Note:** Load tests should be run in controlled environments, not on shared infrastructure.

### Test Documentation

For detailed test information, see:
- `docs/WEB_TEST_REPORT.md` - Comprehensive test report
- `docs/WEB_STATUS.md` - Implementation status

## Troubleshooting

### Backend won't start

- Check if port 8080 is already in use: `lsof -i :8080`
- Use a different port: `PORT=3001 go run ./web/server`

### Frontend can't connect to backend

- Ensure the backend is running
- Check the proxy configuration in `vite.config.ts`
- Verify `VITE_API_BASE` is set correctly if using a custom API URL

### Build fails

- Ensure Node.js 18+ is installed: `node --version`
- Clear node_modules and reinstall: `rm -rf node_modules && npm install`
- Check for TypeScript errors: `npm run build` (will show errors)

### Tests fail

- Backend: Ensure Go tests pass: `go test ./web/server/...`
- Frontend: Check test output: `cd web/ui && npm test`

## Future Enhancements

Potential improvements for future versions:

- Additional diagnostic endpoints
- Real-time progress updates for long-running KAT verification
- Export results as JSON/CSV
- Support for custom KAT vector uploads
- Enhanced error messages with suggestions
- Dark mode theme

