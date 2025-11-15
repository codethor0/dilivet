# LinkedIn Post - DiliVet v0.3.0

DiliVet v0.3.0 is live. I've added a browser-based web interface for ML-DSA (Dilithium-like) signature diagnostics, making the toolkit accessible through any modern web browser.

**What's New:**
• Web UI: React + TypeScript interface for signature verification and KAT testing
• Backend API: Go HTTP server wrapping the existing DiliVet core
• Comprehensive testing: 25+ backend tests, 15+ frontend tests, 12 E2E tests, and load testing with k6
• Docker support: Production-ready containerization

**Key Features:**
The web UI provides three main pages:
- Dashboard with server health monitoring
- Interactive signature verification form
- KAT (known-answer test) verification with detailed results

All backed by the same proven DiliVet CLI core that's been battle-tested for ML-DSA diagnostics.

**Testing Infrastructure:**
I've built a complete testing stack:
- Unit and integration tests for backend and frontend
- End-to-end tests with Playwright (cross-browser: Chromium, Firefox, WebKit)
- Load testing with k6 for stress testing all endpoints
- All tests automated and passing

The web UI is designed for diagnostics and controlled environments. It's not a production crypto library, but a tool for implementers to catch integration bugs before shipping.

If you work on post-quantum signatures, secure implementations, or supply chain security, I'd appreciate feedback on the design, test coverage, and any edge cases that should be covered next.

Repo: https://github.com/codethor0/dilivet
Release: https://github.com/codethor0/dilivet/releases/tag/v0.3.0

#PostQuantum #Cryptography #MLDSA #WebDevelopment #Testing #OpenSource

