# CI and Workflow Policy

This document describes the GitHub Actions workflows for DiliVet and the policies governing their use.

## Overview

DiliVet uses a minimal set of three GitHub Actions workflows:

1. **CI** (`ci.yml`) - Core continuous integration
2. **Fuzz** (`fuzz.yml`) - Fuzz testing
3. **Release** (`release.yml`) - Release automation

## Workflow Details

### CI Workflow (`ci.yml`)

**Purpose**: Run tests, linting, and builds on every push and pull request.

**Triggers**:
- Push to `main` branch
- Pull requests targeting `main`

**Jobs**:
- `go-ci`: Core Go CLI tests and builds
  - Runs `./scripts/check-all.sh` which includes:
    - Go vet
    - golangci-lint (if available)
    - Go tests with race detector
    - Fuzz smoke tests (short runs)
    - Cross-platform builds
- `web-ci`: Web stack tests and build
  - Web backend tests (`go test ./web/server/...`)
  - Frontend tests (`npm test`)
  - Frontend build (`npm run build`)

**Policy**:
- All test steps are **strict** - failures cause the workflow to fail
- No `continue-on-error` on any test or build steps
- Runs on `ubuntu-latest` only
- Both jobs run in parallel

**Cache Strategy**:
- Cache keys prefixed with `dilivet-v1-ci-` to invalidate old caches
- Caches Go modules and build cache based on `go.sum` hash

### Fuzz Workflow (`fuzz.yml`)

**Purpose**: Run fuzz tests to discover edge cases and potential bugs.

**Triggers**:
- Pull requests targeting `main`
- Manual dispatch via `workflow_dispatch`

**Jobs**:
- `fuzz`: Runs fuzz tests for:
  - `FuzzDecodePublicKey`
  - `FuzzVerify`

**Policy**:
- Fuzz tests are **strict** - failures cause the workflow to fail
- Timeout: 20 minutes
- Only runs on `ubuntu-latest`

**Cache Strategy**:
- Cache keys prefixed with `dilivet-v1-fuzz-` to invalidate old caches

### Release Workflow (`release.yml`)

**Purpose**: Build and publish release artifacts when version tags are created.

**Triggers**:
- Push of tags matching `v*` pattern (e.g., `v0.3.0`, `v1.0.0`)
- Manual dispatch via `workflow_dispatch`

**Jobs**:
- `build-and-test`: Validates the codebase before building (strict)
- `build`: Cross-compiles binaries for multiple platforms (strict)
- `checksums`: Generates SHA256 checksums (strict)
- `sbom`: Generates Software Bill of Materials (strict)
- `provenance`: Generates SLSA provenance (optional, may fail)
- `release`: Creates GitHub Release and signs artifacts (optional steps may fail)

**Policy**:
- **Strict jobs**: `build-and-test`, `build`, `checksums`, `sbom` must pass
- **Optional jobs/steps**: `provenance` and signing steps use `continue-on-error: true`
  - These steps may fail if OIDC is not configured or signing keys are missing
  - Failures in optional steps do not block the release

**Cache Strategy**:
- Cache keys prefixed with `dilivet-v1-release-` to invalidate old caches

## Workflow Naming and Organization

All workflows follow these conventions:

- **File names**: Use lowercase with hyphens (e.g., `ci.yml`, `fuzz.yml`, `release.yml`)
- **Workflow names**: Use title case (e.g., `CI`, `Fuzz`, `Release`)
- **Job names**: Use lowercase with hyphens
- **Step names**: Use descriptive, lowercase names

## Adding New Workflows

Before adding a new workflow, consider:

1. **Is it truly necessary?** Can the functionality be added to an existing workflow?
2. **Does it follow the minimal principle?** We aim to keep the workflow count low
3. **Is it well-documented?** Add clear comments explaining purpose and triggers
4. **Does it have proper guards?** Use `if:` conditions to prevent accidental runs

If a new workflow is needed:

1. Create the file in `.github/workflows/`
2. Use the standard header comment with copyright and project info
3. Document the workflow in this file
4. Test locally if possible (using `act` or similar tools)
5. Ensure it follows the cache key naming convention (`dilivet-v1-<workflow-name>-`)

## Disabling Workflows

If a workflow needs to be temporarily disabled:

1. Rename it to `disabled-<original-name>.yml`
2. Change the `on:` trigger to only `workflow_dispatch:`
3. Add a comment at the top explaining why it's disabled and when it might be re-enabled

## Documentation Standards

- **No emojis**: Workflow files, documentation, and commit messages should not contain emojis
- **Professional tone**: Keep descriptions clear and professional
- **Clear comments**: Add comments explaining non-obvious logic or optional steps

## Cache Management

Cache keys are versioned with a prefix (`dilivet-v1-`) to allow invalidation of old caches when needed. To invalidate all caches:

1. Update the prefix in all workflow files (e.g., `dilivet-v2-`)
2. Commit and push
3. Old caches will naturally expire or be replaced

## Troubleshooting

### Workflow not running

- Check that the trigger conditions are met (branch name, file paths, etc.)
- Verify the workflow file is in `.github/workflows/`
- Check GitHub Actions settings for the repository

### Workflow failing unexpectedly

- Review the workflow logs for specific error messages
- Ensure all required secrets are configured
- Verify Go version and dependencies are correct
- Check that test files and scripts are executable

### Release workflow not triggering

- Ensure the tag matches the `v*` pattern
- Verify the tag is pushed to the remote repository
- Check that `workflow_dispatch` is available for manual triggers

## Maintenance

This policy should be reviewed and updated when:

- New workflows are added
- Workflow behavior changes significantly
- New CI/CD patterns are adopted
- Cache strategies are updated

Last updated: 2025-11-14

