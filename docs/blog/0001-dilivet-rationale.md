<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

---
title: "Why DiliVet Exists"
date: 2025-11-09
---

## 1. A toolkit, not a crypto library

Post-quantum signatures are finally exiting the lab. Standards bodies (NIST,
IETF) and implementers (language runtimes, hardware vendors, wallet builders)
need a neutral place to sanity check what they ship. DiliVet exists to fill the
gap between “I copied the reference implementation” and “I have a production
constant-time stack with regression tests”. We do **not** provide a signing
stack; we keep the surface area narrow: load vectors, exercise primitives,
compare behaviours, and highlight mismatches.

The guiding principle is *harsh light, soft edges*. Harsh light means deterministic
hashing, explicit edge vectors, and simple report output so nothing stays hidden.
Soft edges means no YubiHSM or SGX requirement to run the toolkit — just `go`,
`dilivet`, and the vectors you care about.

## 2. Personal safety nets

Implementing ML-DSA involves juggling Montgomery reductions, rejection sampling,
NTTs, and bit packing. Every layer introduces corner cases: truncation bugs,
sign-extension, lenient decoders that accept more than the spec allows. DiliVet
ships with two net types:

1. **Known Answer Tests (KATs)** — canonical `.req` files describing message, key,
   and signature triples. The harness accepts both hand-written edge vectors and
   generated fixtures, merging them into a repeatable suite.
2. **Synthetic edge messages** — `EdgeMsgs` includes empty messages, long runs of
   zeros, alternating patterns, and high-bit toggles. They sound trivial, yet they
   catch many “forgot to reduce mod q” bugs.

Neither replace the official FIPS 204 vectors; they sit next to them, ensuring you
can inject your own adversarial patterns without rewriting supporting code.

## 3. Deterministic by default

When a build fails the toolkit should tell you _why_ in plain text. Deterministic
hashing (`HashDeterministic`) prefixes each part with its length, then feeds it to
SHA-256. This avoids collisions between `pk || msg` and `msg || pk` variations and
makes it trivial to reproduce on other machines. The CLI uses the same primitive,
meaning the `dilivet kat -mode sign` demo teaches the expected request/response
workflow without any randomness.

On the verification side we settle for a practical stance: a signature matches if
its deterministic hash equals the computed one. These stubs are blunt, but safe.
They allow downstream teams to plug in their real implementations via the exec
adapter or the Go interface, swapping our stubs for real ML-DSA implementations.

## 4. Interoperability first

Real-world deployments mix languages and ABI conventions. The `dilivet exec`
subcommand bridges the gap, routing messages, keys, and signatures through
external binaries over stdin/stdout. A Rust signer? A Python verifier? Feed them
hex-encoded inputs and let the adapter wrangle timeouts and exit codes. This makes
interoperability testing a one-liner:

```bash
dilivet exec \
  --sign ./target/release/my-signer \
  --verify ./venv/bin/my-verifier \
  --msg 616263 \
  --pk deadbeef... \
  --sk 0badf00d...
```

If either binary misbehaves, DiliVet captures stderr and folds it into a Go error.
Users get a precise failure reason with a reproducible command line.

## 5. Supply-chain guarantees

Shipping a vetting tool is pointless if the binary cannot be trusted. The release
workflow ticks three boxes:

- **Reproducible builds** — `go build -trimpath -buildvcs=false` with CGO disabled
  across Linux, macOS, and Windows. The produced archives contain only the binary.
- **Cryptographic signing** — we sign `SHA256SUMS.txt` using keyless cosign. Users
  can verify via OIDC without maintaining a keyring.
- **Provenance** — the SLSA3 generator emits provenance metadata linked to the tag
  build. Chain-of-custody is no longer a manual chore.

On top of this, the release workflow produces an SPDX SBOM using Anchore's action,
so consumers know the exact Go stdlib revision and dependency set.

## 6. CI that matches reality

The CI workflow runs on a matrix: Linux, macOS, Windows; Go 1.22 and 1.23. This
mirrors the shipping targets for the binaries and ensures we catch platform quirks
early. Each run performs:

1. `go vet`
2. `golangci-lint`
3. `go test -race` with coverage uploaded to Codecov

The fuzz workflow opens the door to structured chaos. Every PR gets a one-minute
fuzz blast; a nightly cron extends that to catch latent issues. Fuzzing outputs
end up as regression tests or trimmed corpus entries, keeping noise low while
maximising coverage.

## 7. Documentation as onboarding

People glance at the README to decide whether to invest time. The new quickstart
offers a three-command onboarding path: install, version check, run a sample KAT.
The release verification notes demonstrate cosign and SLSA provenance checks in a
copy/paste-friendly snippet. Interop examples show how to wire `dilivet exec`
against external signers, giving teams confidence they can test without refactors.

Under `docs/`, the fuzzing note explains why/when to fuzz, the Wycheproof plan
describes how to curate adversarial vectors, and this post outlines the project
philosophy. Each document is short, actionable, and focused on enabling external
contributors.

## 8. Contribution roadmap

Short term:

- Expand the KAT catalogue with real-world edge cases contributed by downstream
  implementers.
- Land additional CLI tooling for batch analysis (e.g., verifying entire
  directories of signatures).
- Wire the exec adapter into CI pipelines for popular libraries.

Medium term:

- Host nightly fuzz/artifact builds so downstream projects can mirror the latest
  checks without waiting for a release.
- Integrate Wycheproof-style JSON cases once the catalogue stabilises, ensuring
  consistent semantics across the req/rsp and JSON formats.

DiliVet deliberately avoids chasing every spec permutation. Instead it supplies a
boring, reliable harness that makes it easy to spot mistakes. If you ship ML-DSA,
your future self will thank you for running the suite early and often. Contributions
— whether they are vectors, adapters, or documentation — are welcome. The roadmap
is public; discussions and PRs keep us honest. Together we can make the post-quantum
transition less painful.

