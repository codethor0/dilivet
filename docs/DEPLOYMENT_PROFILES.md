# DiliVet Web Deployment Profiles

This document describes three standard deployment profiles for DiliVet Web, from local development to hardened internal deployments.

---

## Profile A: Local Development

**Goal**: Hacking, demos, and local testing with no authentication friction.

**Use Case**: 
- Local development
- Quick demos
- Testing and debugging

**Environment Variables**:
```bash
export MAX_BODY_SIZE=10485760         # 10MB
export REQUEST_TIMEOUT=30s
export ALLOWED_ORIGINS="*"            # Allow all origins (dev only)
export REQUIRE_AUTH=false             # No authentication
```

**Run**:
```bash
# Backend
go run ./web/server

# Frontend dev (optional, for hot reload)
cd web/ui
npm install
npm run dev
```

**Access**: http://localhost:8080

**Security Notes**:
- No authentication
- CORS allows all origins
- Not suitable for network exposure

---

## Profile B: Internal Lab / Single-Tenant (Recommended Default)

**Goal**: Internal PQC lab or trusted VPN-only network deployment.

**Use Case**:
- Internal research labs
- Trusted VPN networks
- Single-tenant deployments
- Controlled internal environments

**Environment Variables**:
```bash
export MAX_BODY_SIZE=10485760         # 10MB (adjust if needed)
export REQUEST_TIMEOUT=30s
export ALLOWED_ORIGINS="https://dilivet.internal.example.com"
export REQUIRE_AUTH=true
export AUTH_TOKEN="$(openssl rand -hex 32)"  # Generate strong token
```

**Deployment**:

**Option 1: Direct (with reverse proxy)**
```bash
# Generate token
export AUTH_TOKEN=$(openssl rand -hex 32)

# Run server
go run ./web/server
```

**Option 2: Docker**
```bash
# Build
docker build -f Dockerfile.web -t dilivet-web:v0.3.0 .

# Run
docker run --rm -p 8080:8080 \
  -e MAX_BODY_SIZE=10485760 \
  -e REQUEST_TIMEOUT=30s \
  -e ALLOWED_ORIGINS="https://dilivet.internal.example.com" \
  -e REQUIRE_AUTH=true \
  -e AUTH_TOKEN="$(openssl rand -hex 32)" \
  dilivet-web:v0.3.0
```

**Reverse Proxy Configuration** (Nginx example):
```nginx
server {
    listen 443 ssl;
    server_name dilivet.internal.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Optional: Client certificate authentication
    # ssl_client_certificate /path/to/ca.pem;
    # ssl_verify_client on;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Pass auth token (if using proxy-level auth)
        # proxy_set_header Authorization "Bearer ${AUTH_TOKEN}";
    }
}
```

**Security Notes**:
- Token authentication enabled
- CORS restricted to specific origin
- TLS termination at reverse proxy
- Network access should be restricted (VPN/firewall)

---

## Profile C: Hardened Internal

**Goal**: Production-ready internal deployment with additional security layers.

**Use Case**:
- Production internal deployments
- Environments with multiple users
- Compliance-sensitive deployments
- High-security internal networks

**Environment Variables**:
```bash
export MAX_BODY_SIZE=5242880          # 5MB (tighter limit)
export REQUEST_TIMEOUT=20s            # Shorter timeout
export ALLOWED_ORIGINS="https://dilivet.internal.example.com"
export REQUIRE_AUTH=true
export AUTH_TOKEN="$(openssl rand -hex 32)"
```

**Additional Hardening**:

1. **Firewall Rules**:
   - Only allow specific IP ranges to access the reverse proxy
   - Block direct access to port 8080 from external networks

2. **Reverse Proxy** (Nginx with mTLS):
```nginx
server {
    listen 443 ssl;
    server_name dilivet.internal.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # Client certificate authentication (mTLS)
    ssl_client_certificate /path/to/ca.pem;
    ssl_verify_client on;
    ssl_verify_depth 2;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=dilivet_limit:10m rate=10r/s;
    limit_req zone=dilivet_limit burst=20 nodelay;

    # IP allowlist (optional)
    # allow 10.0.0.0/8;
    # allow 192.168.0.0/16;
    # deny all;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

3. **Log Shipping**:
   - Ship logs to SIEM/log aggregation system
   - Ensure no sensitive data (keys/sigs/messages) in logs
   - Log only metadata: param set, success/failure, latency, approximate sizes

4. **Docker Deployment**:
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
      - MAX_BODY_SIZE=5242880
      - REQUEST_TIMEOUT=20s
      - ALLOWED_ORIGINS=https://dilivet.internal.example.com
      - REQUIRE_AUTH=true
      - AUTH_TOKEN=${AUTH_TOKEN}
    user: "1000:1000"
    read_only: true
    security_opt:
      - no-new-privileges:true
    networks:
      - internal
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

**Security Notes**:
- All Profile B features
- Tighter resource limits
- mTLS client certificate authentication
- Rate limiting at reverse proxy
- IP allowlisting
- Resource limits in Docker
- Log shipping to SIEM
- Read-only filesystem (optional)

---

## Choosing a Profile

| Profile | Authentication | CORS | Network | Use Case |
|---------|---------------|------|---------|----------|
| **A: Local Dev** | None | Open | Localhost | Development, demos |
| **B: Internal Lab** | Token | Restricted | VPN/Internal | Labs, trusted networks |
| **C: Hardened** | Token + mTLS | Restricted | Firewalled | Production, compliance |

---

## Migration Between Profiles

**From A → B**:
1. Set `REQUIRE_AUTH=true`
2. Generate `AUTH_TOKEN`
3. Set `ALLOWED_ORIGINS` to your domain
4. Deploy behind reverse proxy with TLS

**From B → C**:
1. Tighten `MAX_BODY_SIZE` and `REQUEST_TIMEOUT`
2. Add mTLS to reverse proxy
3. Configure firewall rules
4. Set up log shipping
5. Add resource limits

---

## Quick Reference

**Generate Auth Token**:
```bash
openssl rand -hex 32
```

**Test Authentication**:
```bash
# Without auth (should fail if REQUIRE_AUTH=true)
curl http://localhost:8080/api/health

# With auth
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/health
```

**Test CORS**:
```bash
curl -H "Origin: https://example.com" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     http://localhost:8080/api/verify
```

---

## See Also

- `docs/WEB_SECURITY_REVIEW.md` - Comprehensive security review
- `docs/WEB_UI.md` - User guide and API documentation
- `README.md` - Project overview

