# DiliVet Web Security & Hardening Review

**Version**: v0.3.0  
**Date**: 2025-01-XX  
**Reviewer**: Security hardening pass  
**Status**: Initial review and hardening complete

---

## Executive Summary

This document provides a security and hardening review of the DiliVet Web stack, including threat modeling, security improvements, and deployment recommendations. The review assumes deployment in a **semi-hostile environment** (internal network, but not internet-exposed by default).

**Key Findings:**
- Input validation implemented for all endpoints
- HTTP method enforcement in place
- Structured error responses (no stack traces)
- No request size limits (added)
- No CORS policy (added, strict by default)
- No authentication (design-only recommendation)
- Docker runs as root (hardened)
- Logging may leak sensitive data (sanitized)

---

## 1. Threat Model

### Assets

- **Keys**: Public keys submitted for verification (public by design, but should not be logged)
- **Messages**: Plaintext messages being verified (potentially sensitive)
- **Signatures**: Cryptographic signatures (public, but should not be logged)
- **Diagnostics Results**: KAT verification results and error messages
- **Server Resources**: CPU, memory, disk (DoS risk)

### Attackers

1. **Untrusted users on same network**
   - Can send malicious requests
   - Can attempt DoS via large payloads
   - Can probe for vulnerabilities

2. **Misconfigured deployments**
   - Accidental internet exposure
   - Missing reverse proxy/auth
   - Weak container security

3. **Automated scanners**
   - Common vulnerability scanners
   - Bot traffic
   - Injection attempts

### Assumptions

- **Single-tenant, internal use**: Designed for controlled environments
- **No multi-tenant isolation**: All users share the same instance
- **No account management**: No user accounts or sessions
- **Diagnostics tooling**: Not a production crypto library

### Out of Scope

- **Account management**: User registration, login, sessions
- **Multi-tenant isolation**: Separate data per user
- **Rate limiting**: Should be handled by reverse proxy
- **TLS termination**: Should be handled by reverse proxy
- **Full audit logging**: Basic structured logging only

---

## 2. Backend Hardening (web/server)

### 2.1 Input Validation

**Status**: Implemented

- Parameter set validation (ML-DSA-44/65/87 only)
- Hex decoding with whitespace stripping
- Empty field checks
- JSON parsing with error handling

**Improvements Added:**
- Request size limits (10MB max body)
- Per-request timeouts (30s default)
- Hex length validation (reasonable bounds)

### 2.2 HTTP Method Enforcement

**Status**: Implemented

- All handlers check HTTP method
- Returns 405 Method Not Allowed for wrong methods

### 2.3 Error Messages

**Status**: Good

- Structured JSON responses
- User-friendly messages
- No stack traces or internal details

**Improvements Added:**
- Sanitized error messages (no full keys/sigs in logs)

### 2.4 Request Size Limits

**Status**: Added

- **Max body size**: 10MB (configurable via `MAX_BODY_SIZE`)
- Prevents DoS via large payloads
- Returns 413 Payload Too Large

### 2.5 Timeouts

**Status**: Added

- **Per-request timeout**: 30s (configurable via `REQUEST_TIMEOUT`)
- Context-based cancellation
- Prevents hanging requests

### 2.6 CORS Policy

**Status**: Added

- **Default**: Strict (no cross-origin by default)
- **Configurable**: `ALLOWED_ORIGINS` env var (comma-separated)
- **Headers**: Only necessary headers exposed
- **Methods**: Only POST/GET allowed

### 2.7 Logging

**Status**: Improved

- **Sanitized**: No full keys/signatures in logs
- **Structured**: JSON logs for errors
- **Metadata only**: Only operation type, param set, result
- **No sensitive data**: Messages truncated in logs

---

## 3. Frontend Hardening (web/ui)

### 3.1 Secrets in Bundle

**Status**: Clean

- No secrets embedded
- API base URL uses env var (`VITE_API_BASE`)
- Defaults to `/api` (relative, safe)

### 3.2 Error Messages

**Status**: Good

- User-friendly error display
- No internal stack traces
- Clear validation messages

### 3.3 Input Validation

**Status**: Good

- Client-side validation for UX
- Server remains source of truth
- Form validation prevents invalid submissions

---

## 4. Authentication / Access Control

### 4.1 Current State

**Status**: None

- No authentication implemented
- Open to anyone on the network

### 4.2 Recommended Strategy

**Option 1: Simple Token Auth (Optional)**
- Shared token via `Authorization: Bearer <token>` header
- Token from `AUTH_TOKEN` env var
- Off by default (set `REQUIRE_AUTH=true` to enable)
- Suitable for internal use

**Option 2: Reverse Proxy Auth (Recommended)**
- Deploy behind Nginx/Envoy with:
  - mTLS for client certificates
  - OAuth2/OIDC integration
  - IP allowlisting
- No changes to DiliVet Web required

**Option 3: Gateway Pattern**
- Deploy behind API gateway (Kong, Istio)
- Gateway handles auth, rate limiting, TLS
- DiliVet Web remains stateless

**Implementation**: Simple token auth added as optional feature (disabled by default).

---

## 5. Docker Hardening

### 5.1 Current State

**Issues:**
- Runs as root user
- Uses `alpine:latest` (good, minimal)
- Exposes port 8080

### 5.2 Improvements Added

- **Non-root user**: Created `dilivet` user (UID 1000)
- **Minimal base**: Already using Alpine
- **Port exposure**: Only 8080 (necessary)
- **Health check**: Already present in docker-compose

### 5.3 Recommendations

- **Read-only rootfs**: Consider `--read-only` flag
- **No new privileges**: Use `--security-opt=no-new-privileges`
- **Resource limits**: Set CPU/memory limits in docker-compose
- **Network isolation**: Use custom Docker network

---

## 6. Logging, Metrics, and Observability

### 6.1 Current State

- Basic Go `log` package
- No structured logging
- No metrics

### 6.2 Improvements Added

- **Structured logging**: JSON format for errors
- **Sanitized logs**: No sensitive data
- **Operation metadata**: Log operation type, param set, result
- **Metrics hooks**: Placeholder for future Prometheus integration

### 6.3 Recommendations

- **Prometheus metrics**: Add counters for:
  - Request count by endpoint
  - Error count by type
  - Request duration
- **Distributed tracing**: Consider OpenTelemetry for production
- **Log aggregation**: Use structured JSON for log aggregation tools

---

## 7. Deployment Recommendations

### 7.1 Internal Network (Recommended)

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
      - PORT=8080
      - REQUIRE_AUTH=true
      - AUTH_TOKEN=${AUTH_TOKEN}
      - ALLOWED_ORIGINS=https://internal.example.com
    user: "1000:1000"
    read_only: true
    security_opt:
      - no-new-privileges:true
    networks:
      - internal
```

### 7.2 Behind Reverse Proxy (Best Practice)

```nginx
# Nginx configuration
server {
    listen 443 ssl;
    server_name dilivet.internal.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Client certificate authentication
    ssl_client_certificate /path/to/ca.pem;
    ssl_verify_client on;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 7.3 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `REQUIRE_AUTH` | `false` | Enable token authentication |
| `AUTH_TOKEN` | (none) | Shared token for auth |
| `ALLOWED_ORIGINS` | (none) | Comma-separated CORS origins |
| `MAX_BODY_SIZE` | `10485760` | Max request body size (bytes) |
| `REQUEST_TIMEOUT` | `30s` | Per-request timeout |

---

## 8. Security Checklist

### Pre-Deployment

- [ ] Set `REQUIRE_AUTH=true` if using token auth
- [ ] Generate strong `AUTH_TOKEN` (32+ random bytes)
- [ ] Configure `ALLOWED_ORIGINS` for CORS
- [ ] Deploy behind reverse proxy with TLS
- [ ] Set resource limits in Docker
- [ ] Review and restrict network access
- [ ] Enable firewall rules (only allow necessary ports)

### Post-Deployment

- [ ] Verify health endpoint responds
- [ ] Test authentication (if enabled)
- [ ] Verify CORS policy
- [ ] Check logs for sensitive data leakage
- [ ] Monitor resource usage
- [ ] Set up log aggregation (optional)

---

## 9. Remaining Risks

### High Priority

1. **No rate limiting**: Should be handled by reverse proxy
2. **No TLS termination**: Should be handled by reverse proxy
3. **Single-tenant**: All users share the same instance

### Medium Priority

1. **No audit logging**: Basic logging only
2. **No metrics**: Placeholder for future integration
3. **No distributed tracing**: Consider for production

### Low Priority

1. **No request ID tracking**: Could help with debugging
2. **No graceful shutdown**: Consider signal handling
3. **No health check endpoint**: Already present 

---

## 10. Testing

### Security Tests Added

- **Input validation**: Large payloads, malformed JSON, invalid hex
- **Authentication**: No auth, correct auth, bad auth
- **CORS**: Origin validation, method restrictions
- **Timeouts**: Long-running requests
- **Error handling**: Sanitized error messages

### Test Coverage

- Backend: 25+ tests (including new security tests)
- Frontend: 15+ tests
- E2E: 12 tests
- Security: 10+ new tests

---

## 11. Conclusion

DiliVet Web has been hardened for **internal, controlled environment** deployment. Key improvements:

- Request size limits and timeouts
- CORS policy (strict by default)
- Optional token authentication
- Sanitized logging
- Non-root Docker user
- Structured error responses

**Recommended deployment**: Behind reverse proxy with TLS and client certificate authentication.

**Not suitable for**: Internet-exposed, multi-tenant, or untrusted user environments without additional hardening.

---

## References

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [Docker Security](https://docs.docker.com/engine/security/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

