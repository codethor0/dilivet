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
- If you add or change Known-Answer Tests (KATs), put vectors under `code/clean/testdata` and extend the table-driven tests rather than editing formats in-place.

## Linting and tidy

- Run `go fmt ./...` and `go vet ./...` before opening a PR.
- Run `go mod tidy` to keep `go.mod` consistent.

## Releases

- Do not change binary naming or artifact layout. Release automation expects `dilivet-$GOOS-$GOARCH.zip` and `mldsa-vet-$GOOS-$GOARCH.zip`.
