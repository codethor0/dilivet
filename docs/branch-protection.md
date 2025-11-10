# Branch Protection Settings

This repository relies on a consolidated GitHub Actions setup. Before enforcing
branch protection (for example, on `main`), ensure the status checks line up
with the workflow names below:

- `ci` – matrix build, lint, and coverage
- `fuzz` – Go fuzz targets (`go test -fuzz=Fuzz`)
- `release` – reproducible release pipeline (tags only; optional for PRs)
- `scorecard` – OpenSSF Scorecard evaluation

Recommended configuration:

1. Enable “Require a pull request before merging” with at least one approving
   review.
2. Require status checks to pass before merging, and add `ci`, `fuzz`, and
   `scorecard` as required checks. Keep `release` optional unless you ship
   release candidates from PR branches.
3. Enable “Require branches to be up to date before merging” so consolidation
   changes rerun the `ci` workflow on the merge commit.
4. Optionally, enable “Require signed commits” if your workflow mandates it.

These names replace the older `Go Test`, `Go Lint`, and `Release (multiplatform)`
workflows. Update any existing protection rules to avoid stuck merges.

