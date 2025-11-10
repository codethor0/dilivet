dilivet -version
dilivet -version
# DiliVet

[![Go Report Card](https://goreportcard.com/badge/github.com/codethor0/dilivet)](https://goreportcard.com/report/github.com/codethor0/dilivet)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/codethor0/dilivet)](https://pkg.go.dev/github.com/codethor0/dilivet)
[![Go Test](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml)
[![Lint](https://github.com/codethor0/dilivet/actions/workflows/lint.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/lint.yml)
[![Release](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Tip](https://img.shields.io/badge/Tip-❤-brightgreen)](https://buy.stripe.com/00w6oA7kM4wc4co5RB3Nm01)
[![Monthly](https://img.shields.io/badge/Monthly-♻️-blue)](https://buy.stripe.com/7sY3cobB2bYEdMYa7R3Nm00)

A small toolkit for ML-DSA (Dilithium-like) signature diagnostics and vetting. DiliVet provides test harnesses, known-answer vectors, and simple CLI tools to validate implementations and help developers catch common implementation errors.

## Quick start

Install the primary CLI and verify the binary:

```bash
go install github.com/codethor0/dilivet/cmd/dilivet@latest
dilivet -version
```

Alternative/alias CLI:

```bash
go install github.com/codethor0/dilivet/cmd/mldsa-vet@latest
mldsa-vet -version
```

Verify downloaded release artifacts (when using release zips):

```bash
# Example: check matching SHA256 in SHA256SUMS.txt
# GOOS="$(uname -s | tr 'A-Z' 'a-z')"
# GOARCH="$(uname -m | sed 's/aarch64/arm64/;s/x86_64/amd64/')"
# grep -E "(dilivet|mldsa-vet)-${GOOS}-${GOARCH}\.zip" SHA256SUMS.txt | shasum -a 256 -c
```

## Where to look

- `cmd/` — CLI entrypoints (`dilivet`, `mldsa-vet`)
- `code/` — core packages and tests (KATs under `code/*/testdata`)
- `.github/workflows` — CI (tests, lint, release)
- `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md` — project metadata

## Contributing

See `CONTRIBUTING.md` for the developer quick loop, testing, and release notes.

## License

This project is licensed under the MIT License — see `LICENSE`.

[![Tip](https://img.shields.io/badge/Tip-❤-brightgreen)](https://buy.stripe.com/00w6oA7kM4wc4co5RB3Nm01)
[![Monthly](https://img.shields.io/badge/Monthly-♻️-blue)](https://buy.stripe.com/7sY3cobB2bYEdMYa7R3Nm00)
