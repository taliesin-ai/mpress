# Changelog

Changes to existing M-Press behaviour, tests and performance are recorded here.

## Unreleased

### Fixed — 2026-10-01

- Export now rejects external project-state parents before creating a workspace, uses rooted workspace cleanup and archive reads, and retains safe internal state aliases. An authenticated HTTP regression catches temporary external writes; full Wails ZIP entry contents remain unchanged. Beads `mp-b52.1.3`; export/snapshot qualification continues.

- Confined direct REST/MCP authoring file operations, uploads, backups, trees and preview/site serving to the project using rooted IO; protected preview cleanup from source, state and published-output overlap. Delegated source/cache/static/CSS/navigation inputs now use the same boundary. Safe aliases and ordinary artifacts remain covered. Beads `mp-b52.1.3` stays in progress for the remaining delegates; see `audit/lean-fast-2026-09-30/S03.md`.
- Extended shared build/clean/preview guards to preserve configured navigation, translation glossary, writing guide and contributor guide files outside the main content directory. Added failing-before data-preservation checks for these cases.

### Fixed — 2026-09-30

- Page routes and language components now reject traversal and nonportable path forms. Canonical route/output preflight catches duplicate, multilingual, generated-file, static and snapshot collisions before replacing the current site; cached pages and direct page writes use the same contract. Root/nested/Unicode routes and supported static overrides remain covered. Beads `mp-b52.1.2`; see `audit/lean-fast-2026-09-30/S02.md` for failing-before tests, artifact equivalence and remaining qualification.
- Build and clean now share canonical output validation and rooted removal, rejecting output that overlaps source, assets, configuration, Git or persistent state. Safe internal aliases, missing output directories, export and live preview remain supported; final output symlinks are replaced without deleting their targets. Sentinel regressions and identical ordinary-doc artifact hashes qualify Beads `mp-b52.1.1`; see `audit/lean-fast-2026-09-30/S01.md`.
- Retained literal runtime classes from inline scripts during production CSS purging, including whitespace-separated class lists and template literals. Added failing-before regression coverage at the purge and build boundaries.
- Marked both frontmatter hero theme images eager, high priority and asynchronously decoded. Prepared in [PR #7](https://github.com/leaanthony/mpress/pull/7); public CI passed, publication and field qualification remain pending.
- Prepared Wails documentation fixes for blocking CSS and reported legacy carousel/article images. The released M-Press 1.0.18 production build now safely purges CSS with page-local animation styles: 191,986 → 123,740 bytes. Repeated cold throttled v3 LCP medians improved 9–10%; this is laboratory evidence, not a claim that the reported field outliers are resolved.
- Closed utility popovers no longer intercept mobile touch scrolling or keyboard focus; open panels remain interactive. Thanks to [@dogukanoklu](https://github.com/dogukanoklu), M-Press's first community contributor, for [PR #6](https://github.com/leaanthony/mpress/pull/6).
- Regenerated the embedded MPD viewer CSS and its HTML cache keys for the popover fix, restoring viewer reproducibility tests.
- Added browser checks against generated production CSS for popover focus, hit testing, closing transitions and touch scrolling, including builds with accessibility disabled.

### Maintenance — 2026-09-30

- Audited current `main` at `f6ca003f255fe7636e97d63ba80e4bfcd9cffa32` under a feature freeze.
- Established Beads epic `mp-b52` with 48 prioritized work items, five workstreams, dependencies and acceptance criteria.
- Recorded public race/coverage, vet, vulnerability, strict-build, link-check, fuzz, browser touch, filesystem-safety and full Wails corpus results in `audit/lean-fast-2026-09-30/`.
- Established build time/allocation, generated asset size, binary size and startup baselines for measured optimization.
- Added user-reported LCP task `mp-b52.4.12`, controlled browser measurements and theme/no-JavaScript/autoplay regression checks. See `audit/lean-fast-2026-09-30/LCP.md` for implementation, evidence and remaining deployment work.

The audit baseline remains preserved; verified product fixes are recorded separately above.
