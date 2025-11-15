<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

# Contributing

Thanks for helping improve DiliVet!

## Quick dev loop

Run tests and build the CLIs locally:

```bash
go test -race ./...
go build ./cmd/dilivet && ./dilivet -version
go build ./cmd/mldsa-vet && ./mldsa-vet -version
```

## Making changes

- Prefer small, test-first PRs. Add a unit test in `*_test.go` that demonstrates the desired behavior before implementing.
- Keep CLI entrypoints (`cmd/dilivet`, `cmd/mldsa-vet`) thin — implement features inside packages under `code/` and call them from the CLIs.
- When working with Known-Answer Tests (KATs), store raw vectors under `code/clean/testdata/kats/ml-dsa/` and extend the loaders in `code/clean/kats` (plus their table-driven tests) instead of inventing new formats.

## Linting and tidy

- Run `go fmt ./...` and `go vet ./...` before opening a PR.
- Run `go mod tidy` to keep `go.mod` consistent.

## Releases

- Do not change binary naming or artifact layout. Release automation expects `dilivet-$GOOS-$GOARCH.zip` and `mldsa-vet-$GOOS-$GOARCH.zip`.
