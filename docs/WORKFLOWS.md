# Build and release lifecycle (aligned with jls / jls-*)

**Branch policy:** `main` is the only persistent development branch. Commit source, experiments accepted for retention, workflows, documentation and canonical reports directly to `main`. Never make research or PR branches. Ephemeral `jls-test/*` branches exist only to dispatch testing and are deleted automatically.

Five workflows, structurally aligned with `jacoblockett/jls` and `jacoblockett/jls-template`, adapted for a standalone Go CLI rather than a skill manifest.

- **Build:** reusable `workflow_call`, Go validation then cross-compile and package a complete matrix of Linux, Windows and macOS (amd64/arm64) ZIP executables. Upload exactly one verified `release-bundle` Actions artifact. Never commit binaries.
- **Test:** manual workflow on `main`. Publishes an ephemeral rolling `test` **prerelease and tag**. Deletes any older `test` release and assets first, then re-creates from the latest verified complete bundle.
- **Ephemeral Test Trigger:** create a branch beginning with `jls-test/` (for example `jls-test/ready`). `test-trigger.yaml` dispatches the Test workflow on `main` and deletes the trigger branch even if dispatch fails. This tests **main**, not arbitrary branch contents. To test proposed source, commit it directly to main first, like the JLS lifecycle; do not merge a separate development branch.
- **Nightly:** at 23:00 America/New_York, change-gated against the existing nightly tag. Rebuild only if package inputs/workflows have changed; the rolling nightly release is **wholly replaced** and tag updated. This is not an accumulating asset repository.
- **Stable Release:** manually choose major/minor/patch or `replace current`. Version derives from most recent stable release; a new stable version cannot reuse an existing tag. `replace current` is explicitly destructive and only allowed if stable exists. Stable uses exactly the same validated bundle shape as Test and Nightly.

All publishing workflows share one concurrency group and do not cancel each other. Automation has `contents: write`, and trigger needs `actions: write`. Check GitHub Actions settings if the automatic GITHUB_TOKEN cannot create/delete releases or tags. Keep user-facing releases distinct from repository source; do not add executables or research corpora to Git history.

## Validation matrix

Run uncached Go tests, vet and appropriate race regression on Ubuntu; build all supported targets from checked-in Go source. Ensure artifact count is correct before publication. Do not publish a stable release before the first green Test; add no version file in source merely to appease release tooling.
