# PR #6 review and integration — 2026-09-30

[PR #6](https://github.com/leaanthony/mpress/pull/6), by
[@dogukanoklu](https://github.com/dogukanoklu), is merged at
`815fff32a34c936aabf36fc2d309e7953f349bb4`. The reviewed head was
`b3889e5e8ba939784d993260d508db6ea60d3365`, after refreshing generated viewer
files on the author's maintainer-editable branch. The original fix is
`0ba94d0e41610fe37510a20f7cfa93620a02ddb6`.

The original CSS fix correctly hides closed utility popovers from pointer and
focus interaction and restores interaction while open. The original head's full
race suite failed `TestViewerIsReproducible`: generated CSS plus two HTML CSS
cache keys were stale. Running `go run ./cmd/mpd-fixtures -root .` refreshed
exactly those three files; HTML changes were checked to be cache keys only.
The refresh commit uses the worktree-scoped taliesin-ai identity, and the squash
merge credits the original author plus taliesin-ai as co-author.

Validation of the final PR code:

- `go test -race ./...`, `go vet ./...`, stripped production build: pass.
- `govulncheck@v1.6.0`: zero reachable vulnerabilities; nine module-level findings
  without a reachable call remain the baseline, not a security clearance.
- Strict documentation build: 62 pages, 100 files; link check and release license
  collection pass. Strict build and links also pass with accessibility disabled.
- The author's new CSS test fails against original main via a Go source overlay
  (missing open-state visibility and pointer declarations) and passes with the fix.
- Actual CDP touch gestures on `/configuration/`, `/code-components/` and
  `/data-components/` at 375x667 and 375x375 now scroll the page 186–192px before
  opening and after closing the panel, without diagnostic CSS injection. The open
  accessibility body's swipe still scrolls 168–186px. Original baseline page
  scrolling was zero in all six cases; preserved in `evidence/mobile.json`.
- The focused assertion test checks both utility popovers with accessibility on,
  and the language popover with it off, at both viewport heights. Closed panels
  pass visibility, pointer, focus and hit checks immediately after closing and
  after the transition; open panels remain focusable and hit-testable; page touch
  scrolling works. Logs and touch JSON are in `evidence/pr6-*` / `pr-6-mobile.json`.
- Local private mobile-header regression passes at 320, 375 and 430px. Raw private
  test logs remain private. The full private suite's previously audited failures
  are separate backlog items and have not been resolved by this PR.
- [GitHub CI run 36701466407](https://github.com/leaanthony/mpress/actions/runs/36701466407)
  passed public-tests after first-contributor approval. Private conformance and
  documentation publication were skipped for the PR; skips are not passes.

Chromium CDP emulates touch; these are not physical-device or Safari results.
Ripwire's CSS-constant call graph cannot establish browser coverage; the actual
browser checks provide it. The new test's quality-delta has no gating regression;
its dead-code flag is a name-based false positive for a Go test entry point.

The audit branch integrated upstream through merge `e20ba8191067e2388ca7806ca4362e3cf9a73cbe`. Beads `mp-b52.2.1`
tracks this fix. The broader overflow-family matrix (`mp-b52.2.4`), short-viewport
panel geometry and platform qualification remain open. No feature was added.

[Thank-you and first-contributor congratulations](https://github.com/leaanthony/mpress/pull/6#issuecomment-5909287458)
were posted as taliesin-ai after confirming the merge.
