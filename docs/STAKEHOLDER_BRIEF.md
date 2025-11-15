# DiliVet – ML-DSA Diagnostics and Vetting Toolkit

## Patent Summary and Stakeholder Brief

**Date:** 2025-11-14  
**Version:** 0.2.3  
**Status:** Patent Documentation Complete, Ready for Review

---

## 1. Overview

DiliVet is a diagnostics and vetting toolkit for ML-DSA (Dilithium-like) post-quantum signature implementations. It is designed to help implementers and security engineers catch subtle integration bugs before deployment by providing a focused CLI harness for tests, adversarial vectors, and reproducible verification.

The current codebase is open source and versioned (v0.2.3). The attached patent documentation package formalizes the underlying technical innovations and frames how this work could support both internal security engineering and external ecosystem hardening.

**Repository:** https://github.com/codethor0/dilivet  
**License:** MIT (Open Source)

---

## 2. Problem

ML-DSA / Dilithium-style signatures are moving from research into real deployments, but testing and diagnostics lag behind:

- **Limited Diagnostic Tools**: Most existing tools focus on conformance vectors and throughput benchmarks, not on *debuggable failures*.

- **Hidden Bugs**: Bugs often hide in packing/unpacking, polynomial encoders, random sampling, and glue code around external signers/verifiers.

- **No Unified Harness**: There is no common, reusable harness that combines known-answer tests, adversarial inputs, deterministic sampling, interop, and supply-chain level verification in one place.

This creates risk: implementations can "pass basic tests" yet still be fragile in edge cases, and downstream consumers may not have confidence in binaries they ship or integrate.

---

## 3. Solution (DiliVet)

DiliVet provides a CLI toolkit and harness that treats ML-DSA signatures as something you can interrogate, not just call:

- **Reproducible known-answer testing** (KAT/ACVP-style vectors)
- **Adversarial edge-case generation** targeting encoders, decoders, and verifiers
- **Deterministic sampling** to make failures repeatable and debuggable
- **Interop hooks** for driving external signers and verifiers under controlled scenarios
- **Reproducible release artifacts** with signing and provenance to improve supply-chain transparency

The patent documentation identifies and formalizes the key architectural and algorithmic ideas behind this design.

---

## 4. Core Innovations (High-Level)

The five core innovations captured in the patent package can be summarized as:

### 1) Unified Diagnostics Harness for ML-DSA

A structured harness that treats multiple ML-DSA implementations and binaries as pluggable backends under one interface, with consistent result capture and reporting.

**Key Differentiator:** Extensible framework accepting both official NIST vectors and custom edge cases in unified format.

### 2) Adversarial Vector Engine for Signature Pipelines

A generator for edge-case inputs that specifically targets weak spots such as packing, length limits, truncated encodings, and boundary conditions in message and key handling.

**Key Differentiator:** ML-DSA-specific edge cases (empty messages, zero runs, alternating patterns, high-bit toggles) not found in generic fuzzing tools.

### 3) Deterministic Sampling and Replay Model

A design that couples random sampling with deterministic seeds and logging, so any failure can be faithfully replayed and debugged on demand across machines and builds.

**Key Differentiator:** Length-prefixed deterministic hashing eliminates randomness while preventing collision attacks in test generation.

### 4) Cross-Implementation Interop and "Round-Trip" Testing

A mechanism where one implementation's signer and another implementation's verifier can be exercised in round-trip fashion under the same harness, with structured reporting of mismatches.

**Key Differentiator:** Language-agnostic testing via standard I/O, requiring no code modifications in target implementations.

### 5) Supply-Chain Aware Release Verification

A workflow and metadata model for producing, signing, and verifying release artifacts (binaries, checksums, provenance) alongside the cryptographic diagnostics, so that downstream users can trust both the bits and the tests.

**Key Differentiator:** Integrated cosign signing, SLSA provenance, and SBOM generation in diagnostics toolkit release process.

---

## 5. Business and Strategic Value

### For Internal Teams

DiliVet-style tooling can be used to test internal PQC implementations, vendor libraries, and third-party code before adoption. It can help reduce cryptographic integration risk and provide better assurance to customers.

**Use Cases:**
- Pre-production validation of ML-DSA implementations
- Vendor library evaluation and security assessment
- Interoperability testing across different implementations
- Supply chain security validation

### For Customers and Ecosystem

A hardened and productized version could be offered as part of PQC readiness assessments, secure SDLC tooling, or supply-chain validation offerings.

**Potential Offerings:**
- PQC implementation security audits
- Interoperability validation services
- Supply chain security consulting
- Custom edge case generation for specific implementations

### Differentiation

While many vendors will implement ML-DSA, fewer will have opinionated, reusable diagnostic and supply-chain tooling for those implementations. This creates an opportunity to lead on "how to safely adopt PQC," not just on performance numbers.

**Competitive Advantages:**
- Specialized focus on ML-DSA (not generic cryptographic testing)
- Deterministic reproducibility (not found in prior art)
- Language-agnostic interoperability testing
- Integrated supply chain security

---

## 6. Status and Next Steps

### Current Status

**Codebase:**
- Open source CLI toolkit (v0.2.3)
- Fully functional with all 5 core innovations implemented
- Comprehensive test coverage and documentation
- Production-ready release workflow

**Patent Package:**
- Technical disclosure (`docs/INVENTION_DISCLOSURE.md`)
- Formal patent draft (`docs/PATENT_DRAFT.md`)
- Executive summary and quick reference (`docs/PATENT_SUMMARY.md`)
- ~1,000+ lines of combined documentation
- 7 formal patent claims (5 independent, 2 dependent)

### Recommended Next Steps

#### 1. Technical Review

- [ ] Confirm the 5 core innovations are described accurately and match the implementation
- [ ] Ensure no sensitive, non-public internal information is unintentionally included
- [ ] Verify all code references and file paths are correct
- [ ] Complete the review checklist in `PATENT_SUMMARY.md`

#### 2. Legal Review

- [ ] Have a patent attorney review the draft claims and narrow/strengthen as needed
- [ ] Conduct prior art search (recommended before filing)
- [ ] Decide whether to file a provisional application to establish an early priority date
- [ ] Evaluate international filing strategy (PCT application)

#### 3. Strategy

- [ ] Decide how this work should align with broader post-quantum and supply-chain security roadmaps
- [ ] Evaluate whether to maintain a fully open-source implementation, a hybrid model, or a productized offering on top
- [ ] Consider licensing strategy (defensive patents, cross-licensing opportunities)
- [ ] Assess market timing and competitive landscape

---

## 7. Documentation Package Contents

### Invention Disclosure (`docs/INVENTION_DISCLOSURE.md`)
- Technical disclosure for internal use
- Problem statement and motivation
- 5 core innovations with technical details
- Novelty analysis and prior art considerations
- Commercial applications and use cases

### Patent Draft (`docs/PATENT_DRAFT.md`)
- Formal patent application draft
- Abstract, background, detailed description
- 7 formal patent claims (5 independent, 2 dependent)
- Industrial applicability section
- Notes for patent attorney

### Executive Summary (`docs/PATENT_SUMMARY.md`)
- High-level overview of all 5 innovations
- Patent claims summary
- Prior art differentiation
- Next steps checklist
- Quick reference guide

### Supporting Documentation
- Project rationale (`docs/blog/0001-dilivet-rationale.md`)
- Announcement snippets (`docs/ANNOUNCEMENT_SNIPPETS.md`)
- Citation metadata (`CITATION.cff`)

---

## 8. Key Metrics

### Technical Metrics
- **Codebase Size:** ~5,000+ lines of Go code
- **Test Coverage:** Comprehensive with race detector
- **Documentation:** 1,000+ lines of patent documentation
- **Innovations Documented:** 5 core innovations
- **Patent Claims:** 7 claims (5 independent, 2 dependent)

### Market Metrics
- **Target Market:** Post-quantum cryptography implementers, security auditors, research institutions
- **Competitive Advantage:** Specialized ML-DSA diagnostics (not generic testing)
- **Differentiation:** Deterministic reproducibility, language-agnostic interop, integrated supply chain security

---

## 9. Risk Assessment

### Technical Risks
- **Low:** All innovations are implemented and tested
- **Low:** Open source codebase provides transparency
- **Medium:** Prior art search may reveal similar approaches (mitigation: narrow claims if needed)

### Legal Risks
- **Low:** Open source (MIT) distribution can continue regardless of patent status
- **Medium:** Patent application may be rejected if prior art is too similar (mitigation: conduct thorough prior art search)
- **Low:** Patent provides defensive protection even if not commercialized

### Business Risks
- **Low:** Open source model allows community adoption regardless of patent status
- **Medium:** Market timing for PQC adoption (mitigation: early filing establishes priority)
- **Low:** Competitive response (mitigation: first-mover advantage in specialized diagnostics)

---

## 10. Conclusion

DiliVet represents a novel approach to ML-DSA implementation testing and validation. The patent documentation formalizes the technical innovations and provides a foundation for both defensive IP protection and potential commercialization opportunities.

The open source nature of the codebase allows for community adoption and contribution while the patent documentation protects the underlying innovations. This dual approach provides flexibility for future strategic decisions regarding productization, licensing, or continued open source development.

**Next Action:** Complete technical review checklist, then engage patent attorney for legal review and filing strategy.

---

**Document Version:** 1.0  
**Last Updated:** 2025-11-14  
**Contact:** codethor@gmail.com

