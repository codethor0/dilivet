# ML-DSA Debug Whitepaper

A systematic framework for AI-assisted bug detection in post-quantum cryptographic implementations.

## Usage

```bash
cd code/clean
go test -v
```

### Alias
```bash
go install github.com/codethor0/dilivet/cmd/dilivet@latest
```

## Install
```bash
go install github.com/codethor0/dilivet/cmd/dilivet@latest
```

## Verify
```bash
dilivet -version
# After downloading release zips + SHA256SUMS.txt:
# GOOS="$(uname -s | tr A-Z a-z)"; GOARCH="$(uname -m | sed 's/aarch64/arm64/;s/x86_64/amd64/')"
# grep -E "(dilivet|mldsa-vet)-${GOOS}-${GOARCH}\.zip" SHA256SUMS.txt | shasum -a 256 -c
```
