# Changelog

Changes to existing M-Press behaviour, tests and performance are recorded here.

## Unreleased

### Fixed — 2026-09-30

- Closed utility popovers no longer intercept mobile touch scrolling or keyboard focus; open panels remain interactive. Thanks to [@dogukanoklu](https://github.com/dogukanoklu), M-Press's first community contributor, for [PR #6](https://github.com/leaanthony/mpress/pull/6).
- Regenerated the embedded MPD viewer CSS and its HTML cache keys for the popover fix, restoring viewer reproducibility tests.
- Added browser checks against generated production CSS for popover focus, hit testing, closing transitions and touch scrolling, including builds with accessibility disabled.

### Maintenance — 2026-09-30

- Audited current `main` at `f6ca003f255fe7636e97d63ba80e4bfcd9cffa32` under a feature freeze.
- Established Beads epic `mp-b52` with 48 prioritized work items, five workstreams, dependencies and acceptance criteria.
- Recorded public race/coverage, vet, vulnerability, strict-build, link-check, fuzz, browser touch, filesystem-safety and full Wails corpus results in `audit/lean-fast-2026-09-30/`.
- Established build time/allocation, generated asset size, binary size and startup baselines for measured optimization.

The audit baseline remains preserved; verified product fixes are recorded separately above.
