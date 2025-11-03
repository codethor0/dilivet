# DiliVet

[![Go Test](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/go-test.yml)
[![Lint](https://github.com/codethor0/dilivet/actions/workflows/lint.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/lint.yml)
[![Release](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml/badge.svg)](https://github.com/codethor0/dilivet/actions/workflows/release-dilivet.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

A systematic framework for ML-DSA (post-quantum signature) vetting and diagnostics.

## Overview

DiliVet provides tools to analyze and validate ML-DSA (Modular Lattice Digital Signature Algorithm) implementations. It includes:

- Signature verification testing
- Known-answer test (KAT) vectors
- Diagnostic utilities

## Quick Start

```bash
# Install latest release
go install github.com/codethor0/dilivet/cmd/dilivet@latest

# Verify installation
dilivet -version

## Alternative CLI

An alias CLI is also available:

```bash
go install github.com/codethor0/dilivet/cmd/mldsa-vet@latest
mldsa-vet -version  # should match dilivet version

## Verify
```bash
dilivet -version
# After downloading release zips + SHA256SUMS.txt:
# GOOS="$(uname -s | tr A-Z a-z)"; GOARCH="$(uname -m | sed 's/aarch64/arm64/;s/x86_64/amd64/')"
# grep -E "(dilivet|mldsa-vet)-${GOOS}-${GOARCH}\.zip" SHA256SUMS.txt | shasum -a 256 -c
```
