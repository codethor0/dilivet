# DiliVet Invention Disclosure

**Date:** 2025-11-14  
**Inventor:** Thor Thor  
**Contact:** codethor@gmail.com  
**Project:** DiliVet - ML-DSA Diagnostics and Vetting Toolkit  
**Version:** 0.2.3

---

## Executive Summary

DiliVet is a novel diagnostics and vetting toolkit specifically designed for ML-DSA (Dilithium-like) post-quantum signature implementations. Unlike traditional cryptographic libraries that provide signing and verification functionality, DiliVet focuses on **detecting integration bugs, encoding errors, and implementation flaws** in ML-DSA deployments through reproducible testing, adversarial edge case generation, and interoperability validation.

**Key Innovation:** A specialized diagnostics framework that treats ML-DSA signatures as testable artifacts rather than black-box cryptographic operations, enabling implementers to catch subtle bugs before production deployment.

---

## Problem Statement

### Current State of ML-DSA Implementation Testing

1. **Limited Diagnostic Tools**: Most ML-DSA implementations rely on basic unit tests and official NIST test vectors. There is no specialized toolkit for detecting integration bugs, encoding errors, or edge case failures.

2. **Non-Reproducible Testing**: Traditional cryptographic testing introduces randomness, making it difficult to reproduce failures and debug issues across different environments.

3. **Interoperability Gaps**: Implementations in different languages (Rust, Python, Go, C++) cannot easily be tested against each other without custom integration code.

4. **Supply Chain Blind Spots**: Release artifacts lack verifiable provenance, making it difficult to trust binaries and verify their integrity.

5. **Edge Case Coverage**: Standard test vectors miss subtle bugs in:
   - Polynomial encoding/decoding (bit-packing errors)
   - Montgomery reduction edge cases
   - Rejection sampling boundary conditions
   - Matrix-vector multiplication precision issues

### Impact of These Problems

- **Production Bugs**: Subtle encoding errors lead to signature verification failures in production
- **Security Vulnerabilities**: Implementation flaws can create attack vectors
- **Interoperability Failures**: Different implementations cannot reliably verify each other's signatures
- **Debugging Difficulty**: Non-reproducible failures are difficult to diagnose and fix

---

## Novel Solution: DiliVet Architecture

### Core Innovation 1: Deterministic Testing Framework

**Problem Solved:** Non-reproducible test failures make debugging impossible.

**Solution:** DiliVet implements a deterministic hashing and sampling system that eliminates randomness from the testing pipeline.

**Technical Details:**
- `HashDeterministic()` function prefixes each input component with its length before hashing
- Prevents collisions between `pk || msg` and `msg || pk` variations
- Enables exact reproduction of test failures across different machines
- CLI uses same deterministic primitive, making demo workflows predictable

**Novel Aspects:**
- Length-prefixed deterministic hashing specifically designed for ML-DSA test vector generation
- Eliminates need for seed management or external randomness sources
- Makes test failures immediately reproducible without additional context

**Files:** `code/hash/hash.go`, `code/kat/kat.go`

---

### Core Innovation 2: Adversarial Edge Vector Generation

**Problem Solved:** Standard test vectors miss subtle encoding and boundary condition bugs.

**Solution:** DiliVet generates synthetic edge cases specifically designed to stress ML-DSA implementation components.

**Technical Details:**
- `EdgeMsgs` collection includes:
  - Empty messages (tests zero-length handling)
  - Long runs of zeros (tests polynomial reduction edge cases)
  - Alternating bit patterns (tests sign-extension bugs)
  - High-bit toggles (tests truncation errors)
- Edge vectors complement official FIPS 204 ACVP vectors
- Extensible format allows implementers to inject custom adversarial patterns

**Novel Aspects:**
- Systematic generation of edge cases targeting ML-DSA-specific failure modes
- Focus on polynomial arithmetic, bit-packing, and rejection sampling edge conditions
- Integration with KAT framework allows mixing official and synthetic vectors

**Files:** `code/kat/edgecases.go`, `code/kat/kat.go`

---

### Core Innovation 3: Interoperability Testing Framework

**Problem Solved:** Implementations in different languages cannot easily test against each other.

**Solution:** DiliVet provides a universal adapter that exercises external signers and verifiers via standard I/O.

**Technical Details:**
- `dilivet exec` subcommand routes messages, keys, and signatures through external binaries
- Hex-encoded input/output format works across languages
- Captures stderr and exit codes for precise failure reporting
- Timeout handling prevents hanging processes
- Works with any binary that accepts stdin/stdout (Rust, Python, C++, etc.)

**Novel Aspects:**
- Language-agnostic testing framework for cryptographic implementations
- No code changes required in target implementations
- Enables cross-language signature verification testing
- Provides reproducible command-line interface for interop validation

**Files:** `code/adapter/execsign/exec.go`, `code/cli/app.go`

---

### Core Innovation 4: Extensible KAT Harness

**Problem Solved:** Test vectors are locked into specific formats, making it hard to add custom cases.

**Solution:** DiliVet provides a flexible KAT framework that accepts both official NIST vectors and custom edge cases.

**Technical Details:**
- Supports `.req` format for hand-written edge vectors
- Loads official FIPS 204 ACVP JSON fixtures
- Merges multiple vector sources into unified test suite
- CLI commands: `kat`, `kat-verify` for running test suites
- Diagnostic reporting with pass/fail categorization

**Novel Aspects:**
- Unified framework for official and custom test vectors
- Extensible without rewriting test infrastructure
- Provides diagnostic output for debugging failures
- Enables community contribution of edge cases

**Files:** `code/clean/kats/loader.go`, `code/cli/kat.go`

---

### Core Innovation 5: Supply Chain Transparency Integration

**Problem Solved:** Release artifacts lack verifiable provenance and integrity guarantees.

**Solution:** DiliVet release workflow integrates cosign signing and SLSA provenance generation.

**Technical Details:**
- Keyless cosign signing of `SHA256SUMS.txt` using OIDC
- SLSA3 provenance generation linked to Git tags
- SPDX SBOM generation for dependency transparency
- Reproducible builds with `-trimpath` and `CGO_ENABLED=0`
- Cross-platform binaries (Linux, macOS, Windows) with consistent verification

**Novel Aspects:**
- Integration of supply chain security into diagnostics toolkit release process
- Demonstrates best practices for cryptographic tool distribution
- Enables downstream consumers to verify artifact integrity

**Files:** `.github/workflows/release.yml`

---

## Technical Implementation Highlights

### Deterministic Hashing Algorithm

```go
// HashDeterministic prefixes each component with length to prevent collisions
func HashDeterministic(parts ...[]byte) []byte {
    h := sha256.New()
    for _, part := range parts {
        var lenBuf [8]byte
        binary.BigEndian.PutUint64(lenBuf[:], uint64(len(part)))
        h.Write(lenBuf[:])
        h.Write(part)
    }
    return h.Sum(nil)
}
```

**Why This Matters:** Standard cryptographic hashing doesn't prevent collisions between different concatenation orders. Length-prefixing ensures `Hash(pk, msg)` ≠ `Hash(msg, pk)` even if the bytes are identical.

### Edge Case Generation Strategy

DiliVet generates edge cases targeting specific ML-DSA failure modes:

1. **Polynomial Encoding Bugs**: Empty messages, all-zero patterns
2. **Bit-Packing Errors**: Alternating patterns, high-bit toggles
3. **Rejection Sampling**: Boundary conditions around rejection thresholds
4. **Matrix Operations**: Edge cases in NTT (Number Theoretic Transform) computations

### Interoperability Adapter Design

The exec adapter uses a simple protocol:
- Input: Hex-encoded message, public key, private key (for signing)
- Output: Hex-encoded signature or verification result
- Error handling: Exit codes, stderr capture, timeout protection

This design allows testing any implementation without modification.

---

## Novelty and Non-Obviousness

### Why This Is Novel

1. **Specialized for ML-DSA**: Unlike generic cryptographic testing tools, DiliVet is specifically designed for ML-DSA's unique characteristics (polynomial arithmetic, rejection sampling, bit-packing).

2. **Deterministic by Default**: Most cryptographic testing introduces randomness. DiliVet eliminates randomness to enable reproducible debugging.

3. **Interoperability Focus**: No existing tool provides language-agnostic testing for ML-DSA implementations across different programming languages.

4. **Edge Case Generation**: Systematic generation of adversarial inputs targeting ML-DSA-specific failure modes is not found in existing tools.

5. **Toolkit vs Library**: DiliVet is a diagnostics toolkit, not a cryptographic library. This separation of concerns is novel in the post-quantum cryptography space.

### Non-Obvious Combinations

- Combining deterministic hashing with ML-DSA test vector generation
- Integrating interoperability testing with supply chain security
- Merging official NIST vectors with synthetic edge cases in unified framework
- Using CLI-based architecture for language-agnostic cryptographic testing

---

## Potential Patent Claims (Draft)

### Claim 1: Deterministic Testing Method

A method for deterministic testing of ML-DSA signature implementations, comprising:
- Generating test vectors using length-prefixed deterministic hashing
- Eliminating randomness from the test generation process
- Enabling exact reproduction of test failures across different computing environments

### Claim 2: Adversarial Edge Vector Generation

A system for generating adversarial test vectors for ML-DSA implementations, comprising:
- Systematic generation of edge cases targeting polynomial encoding, bit-packing, and rejection sampling
- Integration with official NIST test vectors in unified framework
- Extensible format for custom edge case injection

### Claim 3: Interoperability Testing Framework

A method for testing ML-DSA implementations across different programming languages, comprising:
- Universal adapter using standard I/O for binary communication
- Hex-encoded input/output format for language-agnostic testing
- Reproducible command-line interface for cross-language validation

### Claim 4: Extensible KAT Harness

A system for managing known-answer tests for ML-DSA implementations, comprising:
- Unified framework accepting both official NIST vectors and custom edge cases
- Extensible format enabling community contribution of test vectors
- Diagnostic reporting with pass/fail categorization

### Claim 5: Supply Chain Integration

A method for distributing cryptographic diagnostics tools with verifiable provenance, comprising:
- Integration of cosign signing and SLSA provenance generation
- Reproducible build process with cross-platform binary generation
- SBOM generation for dependency transparency

---

## Prior Art Considerations

### Known Prior Art

1. **NIST ACVP Test Vectors**: Official FIPS 204 test vectors (public domain)
2. **Wycheproof Test Suite**: Google's cryptographic test suite (different focus, different algorithms)
3. **Generic Fuzzing Tools**: AFL, libFuzzer (not specialized for ML-DSA)
4. **Cryptographic Libraries**: liboqs, pqclean (provide implementations, not diagnostics)

### Differentiation from Prior Art

- **Wycheproof**: Focuses on multiple algorithms, not specialized for ML-DSA edge cases
- **Generic Fuzzers**: Don't understand ML-DSA structure (polynomials, matrices, rejection sampling)
- **Test Vector Loaders**: Don't provide deterministic hashing or interoperability testing
- **Cryptographic Libraries**: Provide implementations, not diagnostic toolkits

### Novel Aspects Not Found in Prior Art

1. Deterministic hashing specifically for ML-DSA test reproducibility
2. Systematic edge case generation targeting ML-DSA failure modes
3. Language-agnostic interoperability testing via standard I/O
4. Unified framework for official and custom test vectors
5. Supply chain security integration in diagnostics toolkit

---

## Commercial and Research Applications

### Target Markets

1. **Cryptographic Implementers**: Companies building ML-DSA into products
2. **Security Auditors**: Third-party security assessment firms
3. **Research Institutions**: Academic researchers working on post-quantum cryptography
4. **Standards Bodies**: NIST, IETF for validation and compliance testing
5. **Open Source Projects**: Community-driven ML-DSA implementations

### Use Cases

1. **Pre-Production Validation**: Catch bugs before shipping
2. **Interoperability Testing**: Verify cross-language compatibility
3. **Security Audits**: Systematic edge case coverage
4. **Research Validation**: Reproducible experimental results
5. **Compliance Testing**: FIPS 204 validation support

---

## Implementation Status

**Current State:** v0.2.3 - Fully functional with:
-  Deterministic hashing implementation
-  Edge case generation
-  Interoperability adapter
-  KAT harness with official NIST vector support
-  Supply chain security integration
-  Full ML-DSA verification implementation (FIPS 204 Algorithm 3)

**Future Enhancements:**
- Expanded edge case catalog
- Additional CLI tooling for batch analysis
- CI integration for popular ML-DSA libraries
- Wycheproof-style JSON case integration

---

## Disclosure Notes

- **Confidentiality**: This disclosure contains proprietary information
- **Public Repository**: DiliVet is open source (MIT license) but this disclosure document is for patent purposes
- **Prior Art Search**: Recommended before filing
- **International Filing**: Consider PCT application for international protection

---

## Contact and Next Steps

**Inventor:** Thor Thor  
**Email:** codethor@gmail.com  
**Repository:** https://github.com/codethor0/dilivet

**Recommended Actions:**
1. Conduct prior art search
2. Review claims with patent attorney
3. Consider provisional patent application
4. Evaluate international filing strategy
5. Document additional novel aspects as they emerge

---

**Document Version:** 1.0  
**Last Updated:** 2025-11-14

