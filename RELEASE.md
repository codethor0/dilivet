<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

# Release Process

DiliVet uses semantic versioning (`vMAJOR.MINOR.PATCH`) and GitHub Releases.

## Pre-Release Checklist

1. **Update CHANGELOG.md**
   - Move items from `## Unreleased` to `## [vX.Y.Z] - YYYY-MM-DD`
   - Group by: Added, Changed, Fixed, Security
   - Link to PR numbers where relevant

2. **Run full tests**
   ```bash
   go test -race ./...
   CGO_ENABLED=0 go build -trimpath ./cmd/dilivet
   CGO_ENABLED=0 go build -trimpath ./cmd/mldsa-vet
   ```

3. **Update version docs**
   - Ensure README examples use `@latest` or document version pinning
   - Check that badges point to correct workflows

## Tagging & Publishing

1. **Create and push tag**
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0: <one-line summary>"
   git push origin v0.2.0
   ```

2. **GitHub Actions auto-builds**
   - `.github/workflows/release-dilivet.yml` triggers on `v*` tags
   - Builds for Linux/macOS/Windows (amd64/arm64)
   - Uploads `dilivet-$OS-$ARCH.zip` and `mldsa-vet-$OS-$ARCH.zip`
   - Generates `SHA256SUMS.txt`

3. **Draft release notes**
   - Go to [Releases](https://github.com/codethor0/dilivet/releases)
   - Auto-generate from commits or paste CHANGELOG section
   - Mark as "Latest release"
   - Publish

## Versioning Rules

- **Patch (v0.1.X):** Bug fixes, doc updates, test improvements
- **Minor (v0.X.0):** New features, backward-compatible changes
- **Major (vX.0.0):** Breaking API/CLI changes (rare for CLIs)

## Hotfix Process

For urgent security fixes:
1. Branch from latest release tag: `git checkout -b hotfix/v0.1.3 v0.1.2`
2. Cherry-pick fix commit
3. Update CHANGELOG with `[v0.1.3] - Security` section
4. Tag and release immediately
5. Merge hotfix branch back to `main`

## Post-Release

- Tweet or announce in community channels
- Close any issues fixed by the release
- Update homebrew tap or other package managers (if applicable)
