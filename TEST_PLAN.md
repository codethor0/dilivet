# Test Plan for DiliVet Repository

## Repository Overview

**Language**: Go 1.23.0  
**Project Type**: ML-DSA (FIPS 204) signature verification toolkit  
**Package Manager**: Go modules (`go.mod`)

## Detected Languages and Frameworks

- **Primary**: Go (golang.org/x/crypto, golang.org/x/sys)
- **Testing**: Go's built-in testing framework (`go test`)
- **Linting**: `golangci-lint`
- **Fuzzing**: Go's native fuzzing (`go test -fuzz`)

## Test Commands

### Primary Test Suite
```bash
go test -race -p 4 ./...
```
- Runs all tests with race detector enabled
- Uses 4 parallel test processes
- Matches CI configuration

### Individual Package Tests
```bash
go test ./code/clean
go test ./code/clean/kats
go test ./code/cli
go test ./code/pack
go test ./code/poly
# etc.
```

### Fuzz Tests
```bash
go test -fuzz=FuzzDecodePublicKey -fuzztime=1m ./fuzz
go test -fuzz=FuzzVerify -fuzztime=1m ./fuzz
```
- Fuzz testing for public key decoding and verification
- 1 minute timeout per fuzz target

## Lint and Format Commands

### Linting
```bash
golangci-lint run --timeout=5m
```
- Uses `golangci-lint` with 5 minute timeout
- Configured via `.golangci.yml` (if present) or defaults

### Formatting
```bash
go fmt ./...
```
- Standard Go formatter

### Static Analysis
```bash
go vet ./...
```
- Go's built-in static analysis tool

## Build Commands

### Standard Build
```bash
go build ./cmd/dilivet
go build ./cmd/mldsa-vet
```

### Cross-Compilation (Release)
```bash
for os in linux darwin windows; do
  for arch in amd64 arm64; do
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch \
      go build -trimpath -ldflags "-s -w" \
      -o "dist/dilivet-$os-$arch" ./cmd/dilivet
  done
done
```

## CI Workflow Commands

Based on `.github/workflows/ci.yml`:

1. **Preflight** (PRs only):
   - `go vet ./...`
   - `golangci-lint run --timeout=5m`

2. **Full Test Matrix**:
   - `go test -race -p 4 ./...`
   - `golangci-lint run` (Ubuntu only)

## Environment Requirements

- **Go Version**: 1.23.0 (or 1.24.x per CI)
- **Tools**:
  - `golangci-lint` (latest, installed via action)
  - Standard Go toolchain

## Dependencies

- `golang.org/x/crypto v0.40.0`
- `golang.org/x/sys v0.35.0` (indirect)

## Test Data

- KAT vectors: `code/clean/testdata/kats/ml-dsa/*.json`
- Test fixtures: Various `*_test.go` files with embedded test data

## Special Considerations

- Race detector is enabled for all tests
- Fuzz tests require Go 1.18+ fuzzing support
- Some tests may require specific test data files
- Cross-compilation tests require multiple GOOS/GOARCH combinations

