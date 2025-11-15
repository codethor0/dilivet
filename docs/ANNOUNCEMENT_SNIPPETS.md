# DiliVet Announcement Snippets

Version: v0.2.3  
Repo: https://github.com/codethor0/dilivet  
Releases: https://github.com/codethor0/dilivet/releases

---

## 1) Twitter / X (short)

DiliVet v0.2.3 is live: a CLI toolkit for debugging ML-DSA (Dilithium-like) signatures with known-answer tests, adversarial edge cases, and deterministic sampling. Binaries for Linux, macOS, and Windows.

Repo: https://github.com/codethor0/dilivet

---

## 2) LinkedIn post (longer)

DiliVet v0.2.3 is live.

DiliVet is a diagnostics and vetting toolkit for ML-DSA (Dilithium-like) signature implementations. The goal is simple: make it easier to catch subtle integration bugs before they ship.

The CLI focuses on:

- Reproducible known-answer testing and ACVP style vectors
- Adversarial edge inputs to stress encoders, decoders, and verifiers
- Deterministic sampling so test runs are stable and debuggable
- Interop helpers that let you exercise external signers and verifiers
- Reproducible release artifacts with cosign and SLSA provenance

The v0.2.x line adds:

- New branding assets and README banner and logo
- ML-DSA polynomial and vector implementations with a deterministic signer stub
- Packing utilities, CLI fuzz harnesses, and cross platform build scripts
- Consolidated CI and fuzz workflows plus OpenSSF Scorecard automation

If you work on post quantum signatures, secure implementations, or supply chain hardening, I would appreciate feedback on the design, test surface, and edge cases that should be covered next.

Repo: https://github.com/codethor0/dilivet  
Latest release: https://github.com/codethor0/dilivet/releases/tag/v0.2.3

---

## 3) Technical blog opening

Modern ML-DSA (Dilithium like) implementations are easy to benchmark and hard to debug. Most of the interesting bugs hide in packing routines, polynomial encoders, and the glue code that connects a signer to real world key and message formats.

DiliVet is a CLI toolkit for treating ML-DSA signatures as something you can interrogate rather than just call. It ships with known answer tests, adversarial edge vectors, and deterministic sampling so implementers can replay the same failure twice and actually understand what went wrong. The same tool also exposes interop hooks for external signers and verifiers, plus reproducible release artifacts with cosign and SLSA provenance for downstream consumers.

---

## 4) Academic / research context

This repository provides a reusable diagnostics harness for ML-DSA style signatures. It is intended for researchers and implementers who need to:

- Validate an implementation against FIPS 204 style ACVP vectors
- Explore adversarial edge cases in encoders, decoders, and verifiers
- Exercise external signing and verification binaries in a controlled loop
- Attach provenance and reproducible checksums to release artifacts

DiliVet is research and diagnostics tooling, not a hardened cryptographic library. It is designed to support experiments on correctness, robustness, and supply chain transparency for ML-DSA deployments.

---

## 5) Reddit / Hacker News blurb

I have been working on DiliVet, a CLI toolkit for debugging ML-DSA (Dilithium like) signature implementations.

The idea is to give implementers a focused harness for:

- Running known answer tests and ACVP style vectors
- Exploring adversarial edge cases around encoders and verifiers
- Wiring in external signers and verifiers for interop testing
- Verifying release artifacts with cosign and SLSA provenance

This is research and diagnostics tooling, not a production crypto library. I would be very interested in feedback from people working on PQC implementations, code review, and supply chain security about what is missing or what should be added next.

Repo: https://github.com/codethor0/dilivet
Latest release: https://github.com/codethor0/dilivet/releases/tag/v0.2.3

