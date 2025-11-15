<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

# Branch Protection Settings

This repository relies on a consolidated GitHub Actions setup. Before enforcing
branch protection (for example, on `main`), ensure the status checks line up
with the workflow names below:

- `CI` – matrix build, lint, and coverage (runs on push to main and PRs)
- `Fuzz` – Go fuzz targets (`go test -fuzz=Fuzz`) (runs on PRs)
- `Release` – reproducible release pipeline (runs on tag pushes only)

Recommended configuration:

1. Enable "Require a pull request before merging" with at least one approving
   review.
2. Require status checks to pass before merging, and add `CI` and `Fuzz` as
   required checks. Keep `Release` optional (it only runs on tag pushes).
3. Enable "Require branches to be up to date before merging" so consolidation
   changes rerun the `CI` workflow on the merge commit.
4. Optionally, enable "Require signed commits" if your workflow mandates it.

**Current Workflow Status:**
- `CI` workflow: Runs on push to `main` and on PRs (tests, lint, cross-build)
- `Fuzz` workflow: Runs on PRs (fuzz testing)
- `Release` workflow: Runs on tag pushes (builds and publishes releases)

**Note:** The workflow names in GitHub Actions are case-sensitive. Use the exact
names as shown in the workflow files: `CI`, `Fuzz`, and `Release`.

