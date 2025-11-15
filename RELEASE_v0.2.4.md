# DiliVet v0.2.4 — FIPS 204 hint fix, property tests, and stress tooling

## Overview

DiliVet v0.2.4 makes the ML-DSA diagnostics toolkit production-ready.

This release fixes the last known correctness gap in the ML-DSA hint logic, adds property-based and fuzz tests around the core algorithms, and introduces CLI stress tooling for large inputs and weird encodings.

---

## Changes

### 1. Correctness: FIPS 204 hint application

- Implemented full FIPS 204 Algorithm 3 for hint application in `useHint()`.
- Hint vectors are now applied exactly as specified:
  - For each coefficient index `i` in the hint vector `h`, the corresponding `r0`/`r1` coefficient is adjusted according to the FIPS 204 rules.
- Added:
  - 3 property-based tests covering:
    - `decompose` / `makeHint` / `useHint` round-trips
    - Boundary conditions for `r0` near ±γ₂
  - 1 fuzz test (5M+ executions) to explore hint edge cases and ensure no panics.

Result: hint logic is now spec-accurate and well covered by tests.

### 2. Matrix A expansion benchmarks

Benchmarks for matrix A expansion:

- ML-DSA-44: ~43 µs/op, 33 KB
- ML-DSA-65: ~82 µs/op, 62 KB
- ML-DSA-87: ~154 µs/op, 115 KB

Analysis:

- Fast enough for single verifications and typical diagnostics use.
- Caching was evaluated and deferred; it can be revisited if profiling ever shows this as a bottleneck in batch scenarios.

### 3. CLI stress and robustness testing

New scripts:

- `scripts/stress-cli.sh`
  - Large message inputs (10 MB, 100 MB)  
  - Weird hex formats (CRLF, BOM, whitespace, mixed case)  
  - Error-path stress loop (50 iterations)

- `scripts/stress-soak.sh`
  - Repeated `kat-verify` runs as a soak test

Results:

- Large inputs handled correctly.
- Weird hex formats parsed or rejected cleanly.
- All error-path stress runs are clean:
  - No panics
  - No hangs
  - Clear error messages
- KAT soak runs complete successfully.

### 4. Documentation and audit

- `docs/AUDIT_DILIVET.md` updated with:
  - Hint fix details and coverage
  - Matrix A benchmark results
  - Stress test design and outcomes
  - TODOs marked complete or explicitly evaluated

---

## Status

- All tests passing: `./scripts/check-all.sh`
  - `go vet`, optional `golangci-lint`
  - `go test -race ./...`
  - Fuzz smoke tests
  - Cross-build matrix
- No critical bugs known.
- Robust under error-path and stress testing.

DiliVet v0.2.4 is production-ready as a diagnostics and vetting toolkit for ML-DSA (Dilithium-like) signature implementations.

