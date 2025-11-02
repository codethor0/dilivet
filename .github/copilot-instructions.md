# Copilot Instructions for this Repo  (`(generated on 2025-11-02)`)

## What this project is
- Go module: `github.com/codethor0/dilivet`.
- Provides two CLIs (thin entrypoints):
  - `dilivet` — primary user-facing CLI under `cmd/dilivet`.
  - `mldsa-vet` — companion/alias CLI under `cmd/mldsa-vet`.
- Focus: ML-DSA (post-quantum signatures) diagnostics/vetting and developer utilities.
- Code layout:
  - `cmd/dilivet/` and `cmd/mldsa-vet/` are thin CLI entrypoints that set version via `-ldflags "-X main.version=${tag}"`.
  - `code/` (and sibling packages) hold core logic. Keep CLI layers minimal; put reusable logic in packages under `code/` and import from CLIs.

## Build & run (local)
- Build current OS/arch:
  - `go build ./cmd/dilivet`
  - `go build ./cmd/mldsa-vet`
- Install from source:
  - `go install github.com/codethor0/dilivet/cmd/dilivet@latest`
- Version flag (wired by CI `-ldflags`):
  - `dilivet -version` and `mldsa-vet -version` print the embedded version.
- For reproducible, small builds prefer `CGO_ENABLED=0` and strip flags: `-trimpath -ldflags "-s -w"`.

## Tests & KATs
- Run tests with race detector: `go test -race ./...` (CI uses the same: see `.github/workflows/go.yml`).
- Known-answer tests (KAT): the repo includes RSP files (e.g., `code/clean/testdata/mldsa_kat.rsp`).
  - Follow table-driven patterns in `*_test.go` for parsing and asserting KAT vectors.
  - Do not change vector formats in-place; add cases or new files and extend the parser.

## Linting & formatting
- CI includes linting (`.github/workflows/lint.yml`). Locally, run `golangci-lint run` if installed.
- Always run `go mod tidy` before committing; CI also runs it.

## Releases (GitHub Actions)
- Tagging `vX.Y.Z` triggers cross-compile and publishes zips + `SHA256SUMS.txt` (see `.github/workflows/release-dilivet.yml`).
- Binaries are packaged as `mldsa-vet-$GOOS-$GOARCH.zip` and `dilivet-$GOOS-$GOARCH.zip` — do not rename these artifacts.

## Conventions & patterns
- Two-entrypoint rule: keep feature logic in packages under `code/`; CLIs only parse flags and call package APIs.
- Deterministic builds: avoid embedding timestamps or nondeterministic data into version outputs.
- Error handling: packages return `error`; CLIs map errors to non-zero exit codes and concise messages.
- Flags: preserve existing flags and semantics (notably `-version`). Add new flags with clear short names.

## Typical workflows
- Quick dev loop:
  1) `go test -race ./...`
  2) `go build ./cmd/dilivet && ./dilivet -version`
- Add a feature:
  1) Implement in a package under `code/`.
  2) Add unit tests in `*_test.go` (table-driven if KATs involved).
  3) Wire CLI flag in `cmd/dilivet/main.go` (and `cmd/mldsa-vet/main.go` if needed).

## Files to look at first
- `go.mod` — module path and dependencies
- `README.md` — quickstart & context
- `cmd/dilivet/main.go`, `cmd/mldsa-vet/main.go` — CLI entrypoints (they already include a `-version` flag wired to `main.version`)
- `code/` — core logic and useful patterns
- `.github/workflows/` — CI (test/lint/release)
- `seed.sh` — bootstrap/maintenance helper used to generate example stubs

## Guardrails for AI edits
- Preserve `CGO_ENABLED=0`, `-trimpath`, and `-ldflags "-s -w -X main.version=..."` for builds used by release workflows.
- Avoid breaking changes to binary names, tag handling, or artifact layout.
- Prefer adding new package APIs rather than pushing logic into `main`.
- If touching KAT parsing, write tests first and keep backward compatibility.

## One-time scan (useful files found in repository)
- README.md
- seed.sh
- go.mod
- cmd/dilivet/main.go
- cmd/mldsa-vet/main.go
- code/clean/mldsa.go
- code/clean/mldsa_test.go
- code/clean/testdata/mldsa_kat.rsp
- .github/workflows/go.yml
- .github/workflows/lint.yml
- .github/workflows/release-dilivet.yml

---
If you want this expanded (examples: a small KAT template, a PR template, or stricter lint/run commands), tell me which area to expand and I will update this file.

_Last updated: 2025-11-02_
