dilivet -version
dilivet -version
# DiliVet

[![Go Report Card](https://goreportcard.com/badge/github.com/codethor0/dilivet)](https://goreportcard.com/report/github.com/codethor0/dilivet)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/codethor0/dilivet)](https://pkg.go.dev/github.com/codethor0/dilivet)
[![Go Test](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml)
[![Lint](https://github.com/codethor0/dilivet/actions/workflows/lint.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/lint.yml)
[![Release](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/codethor0/dilivet/badge)](https://securityscorecards.dev/viewer/?uri=github.com/codethor0/dilivet)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

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
- `code/` — core packages and tests (official ML-DSA KAT loaders live in `code/clean/kats`)
- `code/clean/testdata/kats/ml-dsa/` — bundled FIPS 204 ACVP vectors for offline testing
- `.github/workflows` — CI (tests, lint, release)
- `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md` — project metadata

## Support

- **Questions?** Open a [Discussion](https://github.com/codethor0/dilivet/discussions)
- **Bug?** File an [Issue](https://github.com/codethor0/dilivet/issues/new?template=bug_report.yml)
- **Security?** See [SECURITY.md](./SECURITY.md)

## Contributing

See `CONTRIBUTING.md` for the developer quick loop, testing, and release notes.

## License

This project is licensed under the MIT License — see `LICENSE`.

---

### Test vectors

The repository ships with the official FIPS 204 specification (`code/clean/testdata/fips_204.pdf`) and ACVP “internalProjection” JSON fixtures in `code/clean/testdata/kats/ml-dsa/`. You can sanity-check the parser and data integrity with:

```bash
go test ./code/clean/kats
```

If NIST republishes updated vectors, drop the new JSON files in the same directory and extend the loader tests as needed.
