# Workflow Hardening Summary

This document summarizes the workflow hardening pass applied to DiliVet's GitHub Actions workflows.

## Date
2025-11-15

## Overview

The GitHub Actions workflows were cleaned up and hardened to ensure:
- Minimal set of workflows (only what's needed)
- Strict, predictable behavior (no random red runs)
- Clear separation of concerns
- Professional documentation (no emojis)

## Changes Made

### Phase 1: Minimal Workflow Set

**Status**: Already minimal - only 3 workflows exist:
- `.github/workflows/ci.yml` - Core CI
- `.github/workflows/fuzz.yml` - Fuzz testing
- `.github/workflows/release.yml` - Release automation

No additional workflows were found or needed to be disabled.

### Phase 2: CI Workflow Simplification

**File**: `.github/workflows/ci.yml`

**Changes**:
- Removed path filters from PR triggers (simplified to all PRs to main)
- Removed `feat/**` branch trigger (only `main` now)
- Consolidated jobs:
  - Removed `preflight` job (was PR-only)
  - Created `go-ci` job that runs `./scripts/check-all.sh` (includes vet, lint, tests, fuzz smoke, cross-builds)
  - Created `web-ci` job for web stack tests and build
- Both jobs run in parallel on `ubuntu-latest`
- All test/build steps are strict (no `continue-on-error`)
- Removed unnecessary matrix strategy

**Result**: Simpler, faster CI that runs on every push/PR to main.

### Phase 3: Fuzz Workflow Isolation

**File**: `.github/workflows/fuzz.yml`

**Changes**:
- Removed path filters from PR triggers
- Now triggers on all PRs to main (or manual dispatch)
- Kept 20-minute timeout
- Kept strict failure behavior

**Result**: Fuzz tests run on all PRs, not just when fuzz files change.

### Phase 4: Release Workflow Verification

**File**: `.github/workflows/release.yml`

**Status**: Already properly configured
- Only runs on `v*` tags or `workflow_dispatch`
- All jobs have proper `if:` guards
- Optional steps (provenance, signing) use `continue-on-error: true`
- Strict steps (build, test, checksums, sbom) are required

**No changes needed**.

### Phase 5: Emoji Removal

**Status**: Workflows already clean
- No emojis found in workflow files
- Documentation already states "no emojis" policy
- All workflow files use professional, clear language

### Phase 6: Documentation Updates

**File**: `docs/CI_AND_WORKFLOW_POLICY.md`

**Changes**:
- Updated CI workflow description to match new job structure
- Updated Fuzz workflow triggers description
- Confirmed Release workflow documentation is accurate

## Workflow Structure (Final State)

### CI Workflow (`ci.yml`)

**Triggers**:
- Push to `main`
- Pull requests to `main`

**Jobs**:
- `go-ci`: Runs `./scripts/check-all.sh` (vet, lint, tests, fuzz smoke, cross-builds)
- `web-ci`: Web backend tests, frontend tests, frontend build

**Policy**: All steps are strict - failures cause workflow to fail.

### Fuzz Workflow (`fuzz.yml`)

**Triggers**:
- Pull requests to `main`
- Manual dispatch

**Jobs**:
- `fuzz`: Runs `FuzzDecodePublicKey` and `FuzzVerify` (1 minute each)

**Policy**: Strict - failures cause workflow to fail.

### Release Workflow (`release.yml`)

**Triggers**:
- Push of `v*` tags
- Manual dispatch

**Jobs**:
- `build-and-test`: Strict validation
- `build`: Strict cross-platform builds
- `checksums`: Strict checksum generation
- `sbom`: Strict SBOM generation
- `provenance`: Optional (may fail)
- `release`: Optional signing steps (may fail)

**Policy**: Core jobs are strict, optional jobs use `continue-on-error`.

## Cache Strategy

All workflows use versioned cache keys:
- CI: `dilivet-v1-ci-*`
- Fuzz: `dilivet-v1-fuzz-*`
- Release: `dilivet-v1-release-*`

This allows easy cache invalidation by updating the version prefix.

## Verification

After these changes:
- All workflows are minimal and focused
- No `continue-on-error` on critical test/build steps
- Clear separation between strict and optional steps
- Professional documentation without emojis
- Workflows are predictable and maintainable

## Maintenance

Going forward:
- Keep workflow count minimal (only add if truly necessary)
- Maintain strict behavior for tests and builds
- Use `continue-on-error` only for truly optional steps (signing, provenance)
- Keep documentation updated when workflows change
- No emojis in workflows or documentation

