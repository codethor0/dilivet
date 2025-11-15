<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

# Branch and CI Protection Policy

This document outlines the recommended branch protection and CI policy for the DiliVet repository, specifically for protecting the `main` branch.

## Recommended GitHub Branch Protection for `main`

### 1. Require Pull Requests

- **Enable**: "Require a pull request before merging"
- **Settings**:
  - Require approvals: **1** (at minimum)
  - Dismiss stale pull request approvals when new commits are pushed: **Enabled**
  - Require review from Code Owners: **Optional** (if CODEOWNERS file exists)

### 2. Require Status Checks to Pass

- **Enable**: "Require status checks to pass before merging"
- **Required checks** (minimum):
  - `CI` (the main CI job from `ci.yml`)
  - `Fuzz` (optional but recommended for PRs that touch fuzz targets)

- **Optional checks** (not required but run):
  - `Release` (only runs on tag pushes, not required for PRs)

- **Settings**:
  - Require branches to be up to date before merging: **Enabled**
  - This ensures the `CI` workflow reruns on the merge commit

### 3. Disallow Force Pushes

- **Enable**: "Do not allow bypassing the above settings"
- **Disable**: "Allow force pushes" (should be unchecked)
- **Disable**: "Allow deletions" (should be unchecked)

### 4. Require Signed Commits (Optional)

- **Enable**: "Require signed commits" (if your workflow mandates it)
- This is optional and depends on your organization's security requirements

## Workflow Status Checks

The following workflows correspond to the status checks:

### CI Workflow

- **Name**: `CI`
- **Runs on**: 
  - Push to `main` and `feat/**` branches
  - Pull requests (when Go files, web files, or workflow files change)
- **Jobs**:
  - `preflight` (lint + types) - runs on PRs only
  - `test` (go 1.23 • ubuntu-latest/macos-latest/windows-latest) - runs on all events
- **Required for merge**: **Yes** (at least the `test` job on `ubuntu-latest`)

### Fuzz Workflow

- **Name**: `Fuzz`
- **Runs on**: 
  - Pull requests (when fuzz targets or code changes)
  - Manual dispatch (`workflow_dispatch`)
- **Jobs**:
  - `fuzz` (runs fuzz targets for 1 minute each)
- **Required for merge**: **Optional** (recommended for PRs touching fuzz targets)

### Release Workflow

- **Name**: `Release`
- **Runs on**: 
  - Tag pushes (`v*` tags)
  - Manual dispatch (`workflow_dispatch`)
- **Jobs**:
  - `build` (cross-platform builds)
  - `checksums` (SHA256 checksums)
  - `sbom` (SBOM generation)
  - `provenance` (SLSA provenance)
  - `sign` (cosign signing)
- **Required for merge**: **No** (only runs on tags, not PRs)

## Configuration Steps

1. Navigate to: `Settings` → `Branches` → `Branch protection rules` → `Add rule`
2. Branch name pattern: `main`
3. Configure the settings as described above
4. In "Require status checks to pass", add:
   - `CI` (specifically the `test` job on `ubuntu-latest`)
   - `Fuzz` (optional)
5. Save the rule

## Heavy Tests (Optional)

The following tests are available via scripts but are **not required** for every PR:

- **E2E tests**: `./scripts/test-web-e2e.sh` (requires Docker)
- **Load tests**: `./scripts/test-web-load.sh` (requires k6 and running server)

These can be run manually or via optional scheduled workflows, but should not block PR merges due to their resource requirements and dependencies.

## Quick Reference

**Minimum required checks for PR merge:**
- `CI` (test job on ubuntu-latest)

**Recommended additional checks:**
- `Fuzz` (for PRs touching fuzz targets or core code)

**Not required:**
- `Release` (only runs on tags)
- E2E tests (manual/optional)
- Load tests (manual/optional)

## Related Documentation

- [Branch Protection Settings](./branch-protection.md) - Detailed workflow configuration
- [CI Status](./CI_STATUS.md) - Current CI workflow status (if exists)
- [Web UI Documentation](./WEB_UI.md) - Web UI development and testing
- [Web Security Review](./WEB_SECURITY_REVIEW.md) - Security considerations

## Notes

- Workflow names in GitHub Actions are case-sensitive. Use the exact names as shown in the workflow files: `CI`, `Fuzz`, and `Release`.
- The `preflight` job in `CI` only runs on PRs and is not required for merge (it's a fast check before the full test suite).
- The `test` job in `CI` runs on multiple OSes (ubuntu, macos, windows). For branch protection, require at least the `ubuntu-latest` variant.
- If you need to bypass protection rules in an emergency, repository administrators can do so, but this should be logged and reviewed.

