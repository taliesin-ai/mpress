# M-Press current-feature audit — 2026-09-30

The audit found reproducible data-loss/path-escape defects and the cause of the
reported mobile scrolling problem. Public unit/race checks pass, but they do not
cover these behaviours. Fix correctness first, then optimize measured costs.
The complete actionable scope is the [TODO index](TODO.md), managed in Beads
under epic `mp-b52` with 48 work items and five workstreams.

## Provenance and scope

- Source: `leaanthony/mpress`, current GitHub `main` at
  `f6ca003f255fe7636e97d63ba80e4bfcd9cffa32`.
- Branch/worktree: `audit/lean-fast-20260930`,
  `/var/home/lea/projects/mpress-audit-20260930`.
- The original `/var/home/lea/projects/mpress` worktree was on `mpress-pro` at
  `ba22d9d`; its existing work was preserved.
- Toolchain: source-selected Go 1.26.6, Linux amd64; Ryzen 9 3950X;
  benchmarks reported GOMAXPROCS 32. Host Go 1.26.5 downloaded the selected toolchain.
- Browser: installed Google Chrome 151.0.7922.71, headless CDP; actual dispatched touch
  sequences at 375×667 and 375×375, not simulated `scrollTop` changes.
- Private conformance: local `mpress-tests` revision
  `e0b55f42f6e8e50a511592b015dfecc3b3e249a8`, with explicit source/binary selection.
  Its expectations predate the latest contribution-script changes. GitHub access
  to that private repository was unavailable under `taliesin-ai` during the audit.
- Large corpus: clean Git archive of Wails revision
  `337a7571b2b78063bdb73c3efb99aae1a1e5e34b`, `docs/mpress/` native documents.
  The working Wails checkout had unrelated changes, so it was not built in place.
- No new product features, real hosting deployments or paid translation calls.

This is a source-guided, system-level audit of the existing feature boundary, not
a claim that every line or every browser/device is verified. Source risks and
coverage gaps are labelled separately from reproducible bugs in Beads.

## Validation results

| Check | Result | Evidence |
|---|---|---|
| Public `go test -race -coverprofile ... ./...` | Pass; 73.6% overall statement coverage | [public-tests.log](evidence/public-tests.log) |
| `go vet ./...` | Pass | [vet.log](evidence/vet.log) |
| `govulncheck@v1.6.0 ./...` | No reachable vulnerabilities; 9 findings in required modules outside imported/called paths | [vuln.log](evidence/vuln.log) |
| Strict self-build | Pass; 62 pages, 100 files, no diagnostics | [docs-build.json](evidence/docs-build.json) |
| Self-site link/asset check | Pass | [docs-check.log](evidence/docs-check.log) |
| Release license collection | Pass | Ran `scripts/collect-release-licenses.sh` into disposable output |
| Private conformance | Four top-level tests fail; one optional full-Wails test skips | Classified below; raw private logs remain outside this public repo |
| Real touch scroll audit | Fail before fix; an injected CSS correction restores scrolling | [mobile.json](evidence/mobile.json) |
| Native Wails strict build, `--no-purge-css` | Pass; 2,010 pages, 2,165 files, no diagnostics | [wails-build.json](evidence/wails-build.json) |
| Native Wails link/asset check | Pass | [wails-check.log](evidence/wails-check.log) |
| MPD parser fuzz, 15s | Pass; 32,931 executions | [fuzz-mpd.log](evidence/fuzz-mpd.log) |
| JSON scanner differential fuzz, 10s | Pass; 74,354 executions | [fuzz-json.log](evidence/fuzz-json.log) |
| CSS minifier fuzz, 10s | Pass; 44,192 executions | [fuzz-css.log](evidence/fuzz-css.log) |
| JS minifier initial fuzz, 15s | Worker terminated during minimization; saved seed replay passes | [initial run](evidence/fuzz-js.log), [replay](evidence/fuzz-js-repro.log) |
| Isolated JS minifier fuzz, 8s, one worker | Pass; 2,621 executions | [fuzz-js-isolated.log](evidence/fuzz-js-isolated.log) |
| Repeated docs output hashes | Identical across four warm builds | [output-metrics.json](evidence/output-metrics.json) |

The JS fuzz interruption is an investigation, not a confirmed minifier hang.
The saved reduced input (`const valuz`) returned immediately on replay. Existing
panic-only fuzz tests also do not establish minifier semantic equivalence.

## Confirmed defects and release risks

| Priority / task | Evidence and impact |
|---|---|
| P0 · S01 | `outputDir: content` or `static` succeeds after deleting source/assets. An output under a symlinked parent deletes an external sentinel. `clean` duplicates the same lexical guard. |
| P0 · S02 | `slug: ../escaped` succeeds and writes `escaped/index.html` outside `site/`. Reserved generated-file collisions fail after output has been removed. |
| P0 · S03 | Dev API reads an external file and creates another through a symlinked `content/` parent, both HTTP 200. The path check is lexical. |
| P1 · S04 | A static file symlink is dereferenced and external sentinel content is published into output. |
| P1 · S05 | A strict build with an unparsed directive replaces the previous site before exiting unsuccessfully. |
| P1 · B01 | A closed, transparent accessibility popover is still `display:grid`, with `pointer-events:auto`, and receives touch gestures over most of the page. |
| P1 · B02 | In a 375px-high viewport, the accessibility panel is 357px high with 528px scroll content; its inner body has a 400px minimum, creating constrained nested scrolling. |
| P1 · B06 | Configuration page has duplicate `search`, `accessibility` and `theme` IDs from headings and shell controls. The standalone checker currently collapses IDs into a set. |
| P1 · B07 | Two identical carousels on one page share an ID. Content-hashed tabs/tables/playgrounds need the same instance check. |
| P1 · C05 | A custom script-created class is removed from the merged production stylesheet, while the script is copied correctly. |

Filesystem reproductions operated exclusively on disposable test projects and
sentinel files. See [basic cases](evidence/reproductions.json),
[route/static/component cases](evidence/reproductions-extra.json),
[authoring API cases](evidence/dev-reproductions.json), and
[custom CSS case](evidence/css-reproduction.json).

For the scrolling defect, hit testing at `(185,300)` identified accessibility
content while its popover was closed and transparent. A swipe of 180px produced
zero document scroll on configuration, code-component and data-component pages,
both initially and after opening/closing the panel. Applying closed-state
`visibility:hidden; pointer-events:none` and open-state inverses restored
186–192px document scrolling on the same six cases. The related
[existing PR #6](https://github.com/leaanthony/mpress/pull/6) should be reused or
attributed, with behaviour tests added beyond its CSS-string assertion.

## Private suite failure classification

Four top-level tests fail, covering more than four findings:

- Contribution handoff expects the old launcher shape. Current output uses the
  site-specific script and source path merged into `main`; verify the semantic
  handoff before changing the implementation.
- Image dialog test reports missing restored focus. It opens through a JS click
  and checks immediately after Escape; distinguish timing from an opener/focus bug.
- Translation navigation test expects an older locale-nav serialization.
  Current shared-navigation label maps require semantic rendered-output checks.
- Documentation qualification detects the confirmed shell/heading ID collisions
  and missing reference-corpus exercise of the existing carousel family.

Mobile header click/hit-testing and the selected Axe dark/light fixture pass.
Those checks did not detect touch scrolling. The old full-Wails import gate
skips without its source variable; this audit separately ran the current clean
native Wails corpus instead of claiming that skipped test passed.

## Performance baseline and targets

| Workload | Baseline |
|---|---|
| Synthetic 201-page cold build, 3 iterations × 3 repeats | 559 / 581 / 597 ms per build; median 581 ms; ~483–484 MB allocated, ~2.86M allocations |
| Same workload warm parse-cache build | 541 / 517 / 501 ms; median 517 ms; ~248 MB allocated, ~2.39M allocations |
| Self documentation warm build, four repeats | 199–243 ms internal build duration; 266–313 ms process wall time |
| Clean 2,010-page native Wails cold build | 12.44 s internal, including 4.94 s parsing; CSS purging disabled as required by its current project |
| Same Wails corpus warm build | 7.24 / 7.26 / 7.35 s internal; [three repeated measurements](evidence/wails-warm.json) |
| Linux amd64 stripped binary, `-trimpath -ldflags='-s -w'` | 50,131,209 bytes |
| `mpress version`, one `/usr/bin/time -v` sample | About 60 ms wall; 37,508 KiB peak RSS |

The benchmark's “cold” means cleared M-Press parse/output caches, not a flushed
OS filesystem cache. Initial microbenchmarks ran during other audit work; use
isolated repeated runs for final comparisons. `version` startup is a single
sample, not a statistically established latency result. Wails and synthetic
options differ, so their times cannot be compared as equivalent workloads.

The 8-iteration warm CPU/allocation profile includes one priming build and
fixture setup. It nevertheless identifies bounded candidates to benchmark:

- Knowledge generation: 53.8% cumulative CPU samples; section extraction/full
  HTML DOM parsing is 34.0%. Knowledge generation is ~52% cumulative allocation space.
- Asset-version rewrite callback: 22.3% cumulative allocation space, from
  rereading/copying/replacing every rendered HTML file after CSS optimization.
- Knowledge `makeIndex`: 187.7 MB flat allocations across the profiled run,
  including repeated page-title/route/tag tokenization and pointer postings.
- JSON indentation and simultaneously retained artifact buffers also appear in
  the profile. Optimize compatibility and memory together; do not remove the
  knowledge feature merely to make the benchmark faster.

Evidence: [benchmark](evidence/build-bench.log), [profile workload](evidence/profile-bench.log),
[CPU cumulative](evidence/cpu-cumulative.txt), [allocation top](evidence/alloc-top.txt).
Binary symbol sizes include BSS and are not on-disk size attribution; the large
DRBG BSS symbol must not be reported as 32 MB of removable binary payload.

| Generated self-site asset | Raw bytes | Gzip bytes, deterministic local compression |
|---|---:|---:|
| `assets/mpress.css` | 120,963 | 21,388 |
| `assets/mpress.js` | 73,458 | 21,066 |
| English search index | 189,925 | 59,637 |
| French search index | 186,089 | 58,493 |
| Welsh search index | 8,212 | 3,220 |

Total self-site output: 9,242,821 bytes, including static images and knowledge
artifacts. These gzip sizes describe local compression potential, not measured
host transfer encoding. Actual browser unused-code/load/interaction tracing is
tracked as V03; no Lighthouse score, field CWV or real-iOS/Safari claim was made.

## Source review coverage and follow-up

| Area reviewed | Findings / follow-up keys |
|---|---|
| CLI, configuration, build discovery/parse/render/output/finalization, static copying, link collector | S01–S05, C03–C05, P01/P06/P07, V04/V06 |
| Markdown/MPD parsing, components, literal fences, syntax highlighting and D2 | S02/S09, B07/B08/B11, C06–C08, P09/P10, V05 |
| Reader runtime/CSS, accessibility, mobile overlays, search, component interactions | B01–B11, C05/C06, P07/P08, V03 |
| Dev REST/MCP, authoring paths, previews, watcher, shared config, drawers and contribution | S03, B03/B04/B09, C01/C02, V01/V04/V07 |
| Translation state/provenance/provider failures, import/conversion and navigation | S06, C04/C08/C09, V01/V04/V07 |
| Version capture/verify/mount, export, deploy adapters and licenses | S04/S07, P11, V04/V06 |
| Portable knowledge generation/load/search/MCP and compression | S08, P02–P05, P11 |
| Public/private CI and release gates | V01–V08 |

Structural review used Ripwire notes, hotspots, clones, co-change, quality and
dead-code lenses. Top maintenance risks are devbar `showConfig`, translation
engine, CLI dispatch, dev handlers and `Build`. Complexity/churn are change-risk
signals, not measured performance. The dead-code pass reported zero high-confidence
internal-linkage candidates; that does not establish absence of all unused code.
Clone reports include tests and idioms, so they do not authorize blanket deletion.

Coverage gaps deserve behaviour tests: navigation 0%, version 44.2%, CLI 45.1%,
dev 53.3%, Lighthouse 31.8%. Do not add trivial tests to chase percentages.
Source risks include watcher max-mtime deletion blindness, shared-config races,
conversion rollback, version replacement failure, and minifier heuristics;
they require reproduction before being called confirmed bugs.

## Work management and completion

Beads 1.3.0 was installed from a verified release checksum using the
[official installation/release sources](https://github.com/gastownhall/beads).
Its embedded Dolt database lives in the primary worktree's `.beads/` and is
automatically shared by Git worktrees. `bd doctor` is unsupported in this release's
embedded mode; direct info/stats/ready/export and native backup are used instead.

Run `bd ready --exclude-type epic`, claim one bounded task, update notes with
reproduction/tests/measurements/commit, and close only after its acceptance
criteria are met. Update `CHANGELOG.md` with each implementation change. The
repository JSONL is an issue interchange snapshot; Dolt is the live database.

Execution order: data safety first, mobile/browser blockers and correctness next,
then profiled optimization candidates; qualification work proceeds where its
dependencies allow. The final V08 task depends on all other work items.
An optimization investigation can close with measured no-change evidence;
confirmed bugs and failed required gates cannot be declared complete that way.

Complete the goal only when the Beads backlog and final qualification are
resolved, with before/after results and all required work accounted for. Report
real external access/device limits without pretending they passed.
