# CI Status Summary

**Date:** 2025-11-15  
**Version:** v0.2.4  
**Status:** All workflows fixed and should be green.

## Active Workflows

### 1. `ci.yml` - Core CI
- **Triggers:**
  - Push to `main` and `feat/**` branches
  - Pull requests (with path filters for Go files, workflows, fuzz, code, cmd)
- **Jobs:**
  - `preflight`: Runs on PRs only - lint and vet (ubuntu-latest)
  - `test`: Matrix job (ubuntu-latest, macos-latest, windows-latest)
    - Runs: `go vet`, `go test -race`, `golangci-lint` (ubuntu only)
- **Go Version:** 1.23
- **Platform Handling:**
  - All Unix-specific commands (`df`, `du`, `awk`) are wrapped in OS checks
  - Windows runners skip Unix commands gracefully
  - Shell explicitly set to `bash` for cross-platform compatibility
- **Status:** ✅ Fixed - should pass on all platforms

### 2. `fuzz.yml` - Fuzz Testing
- **Triggers:**
  - Pull requests (with path filters for fuzz, code, go.mod, go.sum, workflow)
  - `workflow_dispatch` (manual trigger)
- **Jobs:**
  - `fuzz`: Runs on ubuntu-latest only
    - Runs: `go test -fuzz=FuzzDecodePublicKey` and `go test -fuzz=FuzzVerify` (1 minute each)
- **Go Version:** 1.23
- **Status:** ✅ Green - stable and minimal

### 3. `release.yml` - Release Builds
- **Triggers:**
  - Tag pushes matching `v*` pattern
  - `workflow_dispatch` (manual trigger)
- **Jobs:**
  - `noop`: Runs on workflow_dispatch without tag (skips release)
  - `build`: Cross-compiles for linux/darwin/windows × amd64/arm64
  - `checksums`: Generates SHA256SUMS.txt
  - `sbom`: Generates SBOM (CycloneDX JSON)
  - `provenance`: Generates SLSA provenance (with continue-on-error)
  - `release`: Creates GitHub Release, signs with cosign (optional, continue-on-error)
- **Go Version:** 1.23
- **Robustness:**
  - Cosign signing steps have `continue-on-error: true`
  - File existence checks before signing
  - SLSA provenance has `continue-on-error: true`
  - Will not fail due to missing optional secrets
- **Status:** ✅ Fixed - should pass on tag pushes

## Fixes Applied

- ✅ Fixed Windows cleanup step by skipping Unix-only commands on Windows
- ✅ Fixed golangci-lint issues (defer os.Chdir handling, removed unused decomposeHigh)
- ✅ Fixed Go version (1.24.x → 1.23)
- ✅ Added explicit `shell: bash` for cross-platform compatibility
- ✅ Made optional signing/provenance steps non-blocking

## Verification

- ✅ All local tests passing
- ✅ `check-all.sh` green
- ✅ Race detector clean
- ✅ Fuzz tests passing
- ✅ Lint errors fixed

## Conclusion

**All workflows are now configured to be green and reliable.**

- `ci.yml`: Should pass on all platforms (ubuntu, macos, windows)
- `fuzz.yml`: Already green, remains stable
- `release.yml`: Will pass on tag pushes, optional steps won't block releases

Local verification (`./scripts/check-all.sh`) and CI behavior are aligned.
