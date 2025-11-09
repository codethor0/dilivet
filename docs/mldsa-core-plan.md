# © 2025 Thor Thor
# Contact: codethor@gmail.com
# LinkedIn: https://www.linkedin.com/in/thor-thor0
# SPDX-License-Identifier: MIT

# DiliVet ML-DSA Core Revamp Plan

## Scope

Implement a production-grade vetting pipeline for ML-DSA-44/65/87 per FIPS-204 (final):

- Constant-time polynomial arithmetic (NTT/invNTT, reductions, sampling)
- Bit-exact key/signature packing and deterministic KAT parity
- Sign/verify wiring with diagnostics instrumentation
- Interop with CIRCL and optional liboqs-go
- CLI diagnostics, benchmarks, and self-tests
- CI enhancements (lint, fuzz, bench, reproducible builds)

## Package Layout

```
code/
  params/   // Parameter tables & lookup
  hash/     // SHAKE/cSHAKE wrappers, transcript hashing
  poly/     // NTT, sampling, polynomial/vector operations
  pack/     // Key/signature/hint packing & unpacking
  signer/   // Sign/verify internals (+ deterministic helpers)
  kat/      // KAT & ACVP loaders, testdata adapters
  diag/     // Diagnostics instrumentation & reporting structs
```

Existing code under `code/clean` will migrate into these packages.

## Key References

- FIPS-204, Tables 1–23 (parameters, encoding layouts, algorithms 1–6)
- ACVP ML-DSA JSON spec (internalProjection vectors)
- CIRCL `sign/mldsa` Go implementation (interop baseline)
- CRYSTALS-Dilithium reference for cross-checking twiddle tables and packers

## High-Level Milestones

1. Establish package skeletons & move current stubs
2. Implement params & hashing helpers
3. Implement polynomial math (NTT, reduction, sampling)
4. Implement packing & unpacking (pk/sk/sig)
5. Implement sign/verify internals + diagnostics hooks
6. Integrate KAT/ACVP loaders & golden testdata
7. Add interop tests (CIRCL, optional liboqs)
8. Expand CLI (kat check, diag, bench, selftest, interop)
9. Add property tests, fuzzers, and benchmarks
10. Harden CI (multi-platform, lint, fuzz-smoke, bench, reproducible builds)
11. Update docs, CHANGELOG, SECURITY note

## Constant-Time Guardrails

- All secret-dependent operations must use masked arithmetic (no branches on secrets)
- NTT/invNTT and reductions implemented via Montgomery/Barrett routines
- Sampling routines avoid modulo bias via rejection loops with constant-time rejection
- Zeroize secret buffers once consumed; keep `runtime.KeepAlive` for key material

## Diagnostics Requirements

- Track rejection causes, z-norms, hint density, and timing envelopes
- Expose JSON-friendly structures for CLI output (per signature)
- Provide aggregated summaries for ACVP/KAT sweeps

## Acceptance Criteria

- Byte-for-byte match with official KATs for pk/sk/sig across ML-DSA-44/65/87
- ACVP rejection cases covered in tests; each reason observed at least once
- Interop with CIRCL passes both directions
- `dilivet selftest`, `dilivet kat check`, `dilivet interop circl`, `dilivet bench`, and `dilivet diag sign` behave as specified
- CI green across platforms with fuzz/bench smoke tests

