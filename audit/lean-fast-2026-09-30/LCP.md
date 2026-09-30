# Documentation LCP fixes — 2026-09-30

Beads `mp-b52.4.12` and [wailsapp/wails#6198](https://github.com/wailsapp/wails/issues/6198) track the user's v3 hero and legacy documentation image report. CSS, hero and carousel fixes are prepared and lab-qualified; animated Portfall, deployment and subsequent field qualification remain pending. No product features were added.

## Changes prepared

M-Press [PR #7](https://github.com/leaanthony/mpress/pull/7), commit `77565ee3da7fcade836695864f064beeab7b7cfc`:

- Production CSS purge conservatively retains literal class tokens from inline scripts, including whitespace-separated strings and backtick literals. This fixes a reason sites had to disable purging to preserve their animation classes. It does not claim to evaluate arbitrary JavaScript or dynamically constructed names.
- Both frontmatter hero theme images use eager loading, high fetch priority and asynchronous decoding. Theme selection and the SVG artwork/animation are unchanged.
- Three focused regressions fail against the original source via a Go overlay and pass against the candidate; public race tests, vet, strict build, links and the reachable-vulnerability scan pass. Public GitHub CI passed; private conformance was skipped and does not count as passing.

Wails local branch `perf/docs-lcp-20260930`, commit `2a0aa66dc5244e33410afb7ff0a1548806e0c902`:

- Restored production CSS purging. The existing morph animation styles are page-local in all ten landing translations, so the change works with the actual released M-Press **1.0.18** before PR #7 is released. Its generated CSS is **191,986 → 123,740 bytes (35.5% less)**. The controlled local server's gzip encoded bodies are **33,029 → 22,287 bytes (32.5% less)**.
- Refreshed 13 hash-bound translation exception entries after those CSS-only source changes. Exception identities, translation segments and reasons are unchanged; page prose matches the original source outside the inserted style rules. Full translation validation passes.
- Legacy homepage prioritizes the first screenshot and defers later carousel slides. Explicit image dimensions and one stable aspect ratio reserve the frame; all screenshots remain visible in full and autoplay remains enabled.
- The first screenshot is delivered as lossless WebP, **91,573 → 83,722 bytes (8.6% less)**, with byte-identical decoded RGBA pixels. The PNG remains available at its existing URL.
- The reported primary article images receive eager/high-priority loading, asynchronous decoding and their natural dimensions in current and default v2.15.0 sources for the configured locales. Older versions and inactive translations are unchanged. Animated WebP artwork is preserved.

Wails work is committed locally and **not pushed**: its `AGENTS.md` requires explicit manual confirmation. Both original user checkouts and unrelated changes were preserved in separate worktrees.

## Controlled before/after evidence

Both Wails builds start at `08bdd6dac3ac02be9985170dd220afa066230635`. Baseline v3 uses released 1.0.18 with `--strict --no-purge-css`; candidate v3 runs the real checksum-verified `docs/mpress/scripts/build.sh` with released 1.0.18 and the source changes above. Baseline and candidate English legacy builds use separate `npm ci` installations from the same committed lockfile.

Chrome **151.0.7922.71**, fresh context for every sample, cache disabled, **150ms latency, 200 KiB/s download and 4x CPU slowdown**. Local text assets are gzipped in both builds; this approximates compression, not Cloudflare's complete delivery configuration. Viewports are **1365×844** and **390×844**. Three samples per theme/viewport/build; medians below use the last observed LCP candidate. Observers start before navigation and run until four seconds after load completes. No parallel build or other browser measurement ran during the formal pairs.

| Site | Viewport/theme | Before LCP | After LCP | Improvement |
| --- | --- | ---: | ---: | ---: |
| v3 | Desktop/light | 828ms | 752ms | 9.2% |
| v3 | Desktop/dark | 832ms | 760ms | 8.7% |
| v3 | Mobile/light | 804ms | 728ms | 9.5% |
| v3 | Mobile/dark | 808ms | 724ms | 10.4% |
| Legacy | Desktop/light | 2,960ms | 1,788ms | 39.6% |
| Legacy | Desktop/dark | 2,956ms | 1,792ms | 39.4% |
| Legacy | Mobile/light | 3,016ms | 2,068ms | 31.4% |
| Legacy | Mobile/dark | 3,016ms | 2,076ms | 31.2% |

These v3 samples qualify the **Wails CSS integration**, not an unreleased M-Press binary: hero priority hints in PR #7 are independently regression-tested, with no extra performance gain attributed to them in this table.

The legacy first screenshot finishes downloading substantially earlier; the original initiates all eight carousel assets together, whereas the candidate initially prioritizes the first screenshot and loads nearby lazy slides later. Chrome can prefetch nearby lazy slides; the claim is deferred loading and reduced contention, not zero requests for every off-screen image.

A separate two-sample-per-condition Portfall article pair reproduces a slow animated-image load: desktop light/dark medians **7,690/7,694 → 7,686/7,686ms**, mobile **7,034/7,028 → 7,004/7,000ms**. Loading hints and dimensions make no material LCP improvement there. This remains an unresolved part of issue #6198, not a fixed score. Recompression pilots preserving all 234 frames and 9,920ms timing produced larger files: pixel-preserving lossless 4,010,258 bytes; quality-80 mixed encoding 1,136,204; a fixed-256-colour palette with lossless encoding 1,621,356. All were rejected; the original 1,060,030-byte animation is unchanged. No frame dropping, static replacement or lower-resolution artwork was shipped to obtain a better metric.

Observed legacy shift sums fall from approximately **0.0083 → 0.00018 desktop** and **0.0106 → 0 mobile**. v3 shift sums remain small: desktop roughly 0.0011–0.0052 and mobile roughly 0.012; dark-desktop samples vary within that range. These are sums over this bounded lab capture, not full-visit field CLS session windows. Raw entries are retained.

## Behaviour and production qualification

- Released production v3 script: 2,030 pages, link check, nine-language translation validation, 16 Python tests and generated-site validation of 2,031 HTML files/eight redirects all pass.
- All seven configured Docusaurus locales build: English, French, Japanese, Korean, Portuguese, Russian and Simplified Chinese. Existing broken-anchor warnings remain baseline warnings; they are not described as new fixes.
- Real browser checks, using generated production styles without injections, pass on desktop/mobile: both theme logos, animation styles, translated no-JavaScript rendering and absence of horizontal overflow. Legacy carousel rotates to a fully loaded lazy slide and maintains heights of 330.9px/193.1px.
- Wails CodeRabbit reviews report no findings for v3, legacy sources and the final carousel changes. The mandated `coderabbit --plain` invocation was attempted; installed 0.7.6 removed that option, so its default plain `review` command was used. Separate bounded reviews avoid the service's 150-file limit. Subsequent edits only restore inactive files and refresh translation digests; complete production checks qualify the final tree.
- Audit browser dependencies stay in the separate tools module. No production dependency or image transformation pipeline was added.
- Final Ripwire quality delta has no unacknowledged preexisting code regression. One explicit acknowledgement records the JSON backlog-array growth as a non-executable verbosity false positive; four new Go test/CLI entry points are reported as dead-code false positives and were actually executed. The broad static test gate still identifies existing uncovered dev handlers and manual CLI entry points; it is not reported as passing coverage. Those broader gaps remain in the audit backlog. Browser checks were rerun after separating hero/carousel test helpers.

## Field qualification remains open

The user supplied v3 light-logo **31,280ms/27 observations** and dark-logo **4,635ms/8**, plus legacy image rows. The report provider, date range and statistic were not supplied. If this is Cloudflare's Core Web Vitals debug table, its LCP column is **P75**, not a sum or mean ([Cloudflare metric documentation](https://developers.cloudflare.com/web-analytics/data-metrics/core-web-vitals/)). Do not silently reinterpret those values.

Direct requests to both custom domains returned **403 with `cf-mitigated: challenge`** during this run. A cold-load probe of the public `wails-v3-site.pages.dev` deployment, using the same stated throttling, gave roughly **0.95–1.0s** LCP before these changes. That does not reproduce or explain the 31-second field observation. A challenge/delivery delay is a plausible separate contributor, not an established cause; no Cloudflare security settings were changed.

The hero SVGs are about 12KB each. The 721KB homepage background starts after first paint in these traces and is not the measured hero LCP gate. Its artwork was preserved. Portfall is a roughly 1MB, 234-frame animated WebP; dropping animation to shrink it would violate the feature freeze. Article hints alone do not establish that all reported image outliers are resolved.

Remaining acceptance: publish the Wails changes after the required confirmation, run remote CI and deploy; release/integrate the separate M-Press fix through its normal gates; collect new same-window field percentiles/counts for both v3 themes and the reported legacy URLs. Inspect challenge/TTFB/resource/render phases when a field slow load is reproduced. Physical iOS/Safari and full real-user field results were not tested here. `mp-b52.4.12`, the final audit qualification gate and the persistent goal remain open.

## Reproduce and inspect

Commands and browser contract prerequisites are in [tools/README.md](tools/README.md). Public raw results and CodeRabbit/browser logs are in [evidence/lcp](evidence/lcp/). `v3-formal-*.json` and `v2-formal-*.json` are the qualified homepage pairs; `portfall-formal-*.json` are the article pair with hints only and `live-pages.json` is a limited live baseline probe. Exploratory uncompressed/purged pilots were excluded from the formal comparison.
