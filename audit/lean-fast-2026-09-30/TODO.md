# M-Press current-feature audit TODOs — 2026-09-30

Baseline: `f6ca003f255fe7636e97d63ba80e4bfcd9cffa32`. Beads epic: `mp-b52`.

Beads is the live source of status and dependencies. This file indexes the initial scope; use `bd ready`, `bd show ID` and `bd list --parent ID` for current work. `backlog-map.json` maps audit keys to Beads IDs.

P0 = data loss / filesystem escape; P1 = major behaviour or qualification failure / high measured cost; P2 = bounded coverage, robustness or optimisation investigation.

A source risk or coverage gap is not a confirmed defect. Resolve investigations with measured evidence; do not invent fixes to satisfy a list.

## Protect source files and filesystem boundaries — `mp-b52.1`

### S01 · `mp-b52.1.1` · P0 · Prevent build and clean from deleting source, state or external directories

Evidence: **confirmed**. Sources: `internal/site/build.go:safeOutput; cmd/mpress/main.go:clean; internal/config/config.go`.

Disposable CLI reproductions: outputDir=content and outputDir=static both exit 0 and delete sentinel source/assets. A symlinked parent linked/site deletes an external sentinel. Current guard checks lexical project containment only.

Done when: Validate canonical containment and disjointness from content, static, configuration, Git and persistent state BEFORE any removal. Cover equal, ancestor, descendant, symlinked and non-existent paths, override output paths and clean. Existing sentinels survive every rejection; ordinary, export and live-preview output still work.

Resolution: implemented and qualified on 2026-09-30; see [S01.md](S01.md) for the sentinel regressions, unchanged artifact hashes and remaining filesystem-audit boundaries. Beads carries the authoritative closure and commit.

### S02 · `mp-b52.1.2` · P0 · Contain and normalise every generated page route

Evidence: **confirmed**. Sources: `internal/content/content.go:routeFor/outputFor; internal/site/render_pages.go; internal/config/config.go`.

CLI reproduction: frontmatter slug ../escaped exits 0 and creates project/escaped/index.html outside site. Reserved assets/mpress.css route fails late after clearing output. Duplicate routes are compared before filesystem normalisation.

Done when: Reject dot segments, platform separators, invalid language path components, output escapes and generated-file namespace collisions before writes. Detect equivalent routes deterministically. Cover Markdown, MPD, multilingual routes and Windows path forms; source/outside sentinels remain unchanged.

Resolution: implemented and qualified on 2026-09-30; see [S02.md](S02.md) for fresh/cached route regressions, output preflight, artifact equivalence and private conformance correction. Beads carries the authoritative closure and commit.

### S03 · `mp-b52.1.3` · P0 · Keep authoring file reads and writes inside the project through symlinks

Evidence: **confirmed**. Sources: `internal/dev/server.go:sourcePath/handleFile/writeSource/handleLivePreview; internal/dev/mcp.go`.

Black-box dev API: GET content/linked/outside.md returns an external sentinel, and PUT content/linked/new.md returns 200 and creates an external file. linked is a directory symlink. sourcePath currently validates only lexical paths.

Done when: Use a shared confined file boundary for REST and MCP reads, writes, previews and uploads. Reject external parent/final symlinks before reading or creating files; cover non-existent targets, existing links, authenticated remote mode and read-only mode. External sentinels never change.

### S04 · `mp-b52.1.4` · P1 · Prevent static copying from publishing external symlink targets

Evidence: **confirmed**. Sources: `internal/site/build.go:copyTree/copyFile; internal/version/version.go:copyTree; internal/exportzip/export.go`.

CLI reproduction: static/linked.txt points to an external sentinel; strict build succeeds and copies external-test-data into site/linked.txt. WalkDir does not follow directories, but copyFile follows file links.

Done when: Specify and enforce a consistent symlink policy across build, export and snapshots: no ambient external files in output. Cover internal/external files, directories, dangling links and non-regular files. Preserve ordinary assets and diagnose rejected entries clearly.

### S05 · `mp-b52.1.5` · P1 · Preserve the last good site when a strict build fails

Evidence: **confirmed**. Sources: `internal/site/build.go:Build; internal/dev/server.go:rebuildLocked; internal/operations/operations.go`.

Strict CLI reproduction with an unparsed directive exits 1 after replacing the prior output; an existing site sentinel is gone. Strict hasErrors runs at the very end. Navigation, rendering and I/O errors can also occur after RemoveAll.

Done when: Validate early and stage output for promotion only on success. Failed parsing, navigation, CSS/minification, copy, mount and finalization leave the prior site intact. Check collectors, timings, exports and dev serving work with staging; test promotion failure recovery.

### S06 · `mp-b52.1.6` · P1 · Make document conversion recoverable across partial filesystem failures

Evidence: **source-risk**. Sources: `internal/projectconvert/convert.go:Run; internal/translate/conversion.go`.

Run claims atomic conversion but sequentially renames every target, then deletes sources, then saves changed configuration. A later rename/delete/config failure leaves earlier changes in place. Temporary names are predictable target+.tmp.

Done when: Inject target write/rename, source deletion and config/state-save failures. Preserve originals and translation state or perform verified rollback; reject target collisions and existing temporary symlinks. Successful Markdown/MPD conversions remain semantically equivalent.

### S07 · `mp-b52.1.7` · P1 · Validate and protect every version snapshot operation

Evidence: **source-risk**. Sources: `internal/version/version.go:Capture/List/Verify/Remove/Mount/manifestFilePath`.

Capture uses dest+.tmp, removes it preemptively and removes the destination before rename; Verify joins a caller label without the Capture label validator. Manifest path checking is lexical. Public package coverage is 44.2%.

Done when: Exercise invalid labels, symlinked artifacts/manifest paths, simultaneous captures, overwrite failure, missing/extra/tampered files and read-only destinations. Reuse one label/path policy; preserve the last verified snapshot on failure and keep mounted navigation correct.

### S08 · `mp-b52.1.8` · P2 · Bound and confine knowledge artifact loading

Evidence: **source-risk**. Sources: `internal/knowledge/artifact.go:Load/readJSON; internal/knowledge/storage.go:readArtifact`.

Artifact filenames from a manifest are joined before digest validation; gzip decompression uses unbounded io.ReadAll. Digest checking protects consistency, but does not bound memory or constrain filenames.

Done when: Test traversal, absolute filenames, symlink escapes, unsupported schemas, duplicate identities, gzip truncation and excessive decompression with small disposable fixtures. Enforce documented limits without rejecting supported large generated bundles; preserve digest checks.

### S09 · `mp-b52.1.9` · P2 · Apply parser resource limits before large allocations

Evidence: **source-risk**. Sources: `internal/mpd/parser.go:ParseWithOptions; internal/content/d2.go; internal/components/components.go`.

MPD allocates node/attribute capacity and scans sourceIsASCII before MaxFileSize rejection. D2 has a 30s timeout and per-diagram ruler. Nested Markdown components have repeated passes. These are source risks, not measured hangs.

Done when: Add boundary workloads for oversized input, deep nesting, attribute count and invalid UTF-8; ensure rejection happens before avoidable allocations. Verify D2 import isolation and cancellation. Document supported limits and measure worst-case time/memory; no arbitrary feature removal.


## Restore mobile scrolling and qualify browser behaviour — `mp-b52.2`

### B01 · `mp-b52.2.1` · P1 · Make closed popovers inert so mobile touch scrolling works

Evidence: **confirmed**. Sources: `internal/site/theme_default.go:.utility-menu-panel; internal/site/accessibility.go:.mpress-accessibility-panel`.

Real CDP touch gestures on /configuration/, /code-components/ and /data-components/ at 375x667 and 375x375 leave scrollY=0 before and after panel close. elementFromPoint hits closed opacity:0 accessibility content. Injected inert CSS restores 186-192px scrolling. Related upstream PR #6.

Done when: Closed popovers cannot intercept pointer/touch events or keyboard focus; open popovers remain interactive. Reproduce swipe/hit-testing before and after open/close, including transitions, production purged CSS and accessibility disabled. Reuse or attribute PR #6 where appropriate.

### B02 · `mp-b52.2.2` · P1 · Keep accessibility settings reachable in short mobile viewports

Evidence: **confirmed**. Sources: `internal/site/accessibility.go:.mpress-accessibility-body/mobile rules; accessibilityJS panel positioning`.

At 375x375, panel clientHeight=357 while scrollHeight=528, and body min-height=400. The panel and body have nested scrolling. At 375x667 body scrolling works. This is a confirmed constrained-viewport overflow, requiring reachability regression coverage.

Done when: All tabs/settings/reset/close remain reachable by touch and keyboard in portrait, landscape and 200% text, including width/font preferences. Reduce inflexible minimum sizing and assert the actual intended scroll container and panel bounds.

### B03 · `mp-b52.2.3` · P1 · Verify contribution dialog scrolling and focus on phones

Evidence: **source-risk**. Sources: `internal/site/theme_default.go:.mpress-contribute-dialog/mobile rules; contribution dialog handlers`.

The mobile dialog receives max-height:calc(100dvh - 1rem) while its base rule has overflow:hidden and content grows when setup/explanation is opened. Existing private test checks strings and click state, not reachability.

Done when: Open choices, local setup and explanation on 320/375px phones and short landscape. Swipe to every control, restore focus on Escape/close and retain source route in generated launcher arguments. Fix clipping only if reproduced; retain site-specific scripts.

### B04 · `mp-b52.2.4` · P1 · Add durable real-touch scrolling regressions for every overflow family

Evidence: **coverage-gap**. Sources: `browser test harness; internal/site/theme_default.go; internal/site/accessibility.go; internal/dev/assets/devbar.css`.

Current private mobile tests use click/hit-testing at 320/375/430px; no DispatchTouchEvent. Audit harness demonstrates why passing click tests does not protect scrolling.

Done when: Automate real swipes for document body, long sidebar, code/terminal, ordinary/data tables, tabs, search results, popovers, dialogs and dev drawers. Assert scroll deltas and background locking before/after close, with production CSS and small viewports. Keep browser dependencies out of the production binary.

### B05 · `mp-b52.2.5` · P1 · Resolve dialog Escape focus restoration failures

Evidence: **test-failure**. Sources: `internal/site/theme_default.go:image/contribution dialog handlers; private browser_test.go`.

Private image-dialog test reports FocusAfterClose empty. It opens with JS click and checks immediately after Escape, so timing and opener tracking must be distinguished. Image close handler clears content without explicit opener restoration.

Done when: Use real keyboard/pointer activation and wait for native close/focus events. Cover image and contribution dialogs, repeated opens, programmatic close, backdrop and Escape. Fix actual focus loss and repair stale/timing-sensitive expectations without weakening behaviour.

### B06 · `mp-b52.2.6` · P1 · Prevent heading anchors from colliding with shell control IDs

Evidence: **confirmed**. Sources: `internal/site/template.go; internal/site/theme_default.go selectors; internal/content/content.go heading IDs`.

Whole-output scan confirms duplicate IDs search, accessibility and theme on configuration/index.html. Private release scan also reports them. check currently stores IDs in a set and does not diagnose duplicates.

Done when: Give shell controls a reserved namespace or otherwise preserve unique IDs while retaining existing public heading anchors. Check all docs and locale pages, ARIA references and keyboard shortcuts. Add generated-output duplicate-ID regression coverage.

### B07 · `mp-b52.2.7` · P1 · Assign repeated components unique deterministic instance IDs

Evidence: **confirmed**. Sources: `internal/components/carousel.go; internal/components/components.go:Tabs; internal/components/table.go; API playground IDs`.

Two identical @carousel blocks in one page produce the same mpress-carousel-a2db93368d ID. Tabs and tables also use content hashes; identical instances need explicit verification. Meta map iteration in Tabs rendering can also affect byte determinism.

Done when: Repeated identical carousels/tabs/tables/playgrounds have unique IDs and correct local ARIA relationships. IDs and output bytes remain deterministic across concurrent cold/warm builds; explicit authored IDs get clear collision diagnostics. Test independent interaction.

### B08 · `mp-b52.2.8` · P2 · Exercise keyboard behaviour across existing interactive components

Evidence: **coverage-gap**. Sources: `internal/site/theme_default.go runtime; internal/components; internal/dev/assets/devbar.js`.

Public tests mostly assert generated strings. Private browser suite covers selected tabs/image/preferences, leaving table filters, custom selects, carousel, tutorial, calendar, audience and API playground paths largely untested.

Done when: Test keyboard navigation, Enter/Space activation, Escape, disabled states, focus order, local ARIA state and independent instances for shipped component families. Use mock HTTP for playgrounds; no new widgets or product behaviour.

### B09 · `mp-b52.2.9` · P2 · Test overlay stacking and scroll-lock restoration sequences

Evidence: **coverage-gap**. Sources: `internal/site/theme_default.go:closeNavigation/openSearch/closeSearch/popovers/dialogs; internal/dev/assets/devbar.js:lockPageScroll`.

Navigation/search use body classes; dev drawers save inline html/body overflow; native dialogs and popovers add top-layer interactions. Closing one overlay while another remains open needs sequence coverage.

Done when: Cover nav to search, accessibility to dialog, drawer transitions, resize to desktop, Escape and repeated close. No stale scroll locks, hidden interceptors or background swipe leak; focus returns to a visible control.

### B10 · `mp-b52.2.10` · P2 · Qualify reader preferences and accessibility on representative pages

Evidence: **coverage-gap**. Sources: `internal/site/accessibility.go; internal/site/theme_default.go; Axe browser qualification`.

Axe passes one fixture in dark/light, but docs duplicate IDs exist and actual phone overlays were broken. Large text, high contrast, bionic text, widths, colours and reduced motion mutate the live DOM/layout.

Done when: Run Axe and keyboard/reachability checks over docs, landing, blog, data components and dialogs under dark/light, 200% text and reduced motion. Verify preference persistence/reset and storage-disabled/corrupt storage fallback. Record the limits of automated checks.

### B11 · `mp-b52.2.11` · P2 · Validate progressive rendering with JavaScript or optional features disabled

Evidence: **coverage-gap**. Sources: `internal/site/template.go; internal/site/theme_default.go; config Search/Accessibility/Contribution; component fallback CSS`.

The template uses an inline js enhancement marker and conditional reader code. Real no-JS rendering and disabled-feature combinations are not qualified by the public string tests.

Done when: Content/navigation remain readable with JS blocked; unenhanced tab panels and expanded content are accessible. Search/accessibility/contribution disabled do not leave dead controls, hidden content or runtime errors. Verify production purging retains fallback styles.


## Qualify build, authoring, translation and routing correctness — `mp-b52.3`

### C01 · `mp-b52.3.1` · P1 · Make the dev watcher detect deletions, renames and non-monotonic mtimes

Evidence: **source-risk**. Sources: `internal/dev/server.go:watch/latest`.

watch compares only maximum regular-file ModTime and rebuilds only when next.After(last). Deleting a non-newest page or changing an older file to a timestamp below the maximum can be missed. It scans the entire project every 500ms.

Done when: Black-box edit/delete/rename/create and preserved/older-mtime changes trigger exactly the needed rebuild. Ignore output/cache/VCS/unrelated files, handle errors explicitly and debounce bursts. Measure idle work and rebuild latency on large sites.

### C02 · `mp-b52.3.2` · P1 · Refresh dev configuration safely and remove shared-config races

Evidence: **source-risk**. Sources: `internal/dev/server.go:Server.cfg/watch/rebuildLocked/handleConfig/handleTranslations; newMCPHandler`.

s.cfg is read across concurrent handlers/watch and assigned by config/translation operations without a uniform lock. rebuildLocked reloads config inside site.Build but does not refresh s.cfg or the file server. Public race tests do not stress these concurrent paths.

Done when: Run concurrent requests/config writes under race detection and external config-edit black-box tests. Serve the correct output after outputDir/language/theme changes; invalid config keeps the prior usable state. Use immutable snapshots or a clear lock contract and avoid deadlocks.

### C03 · `mp-b52.3.3` · P2 · Verify parse-cache invalidation and corruption recovery end to end

Evidence: **coverage-gap**. Sources: `internal/site/parse.go:parseCacheVersion/key/decode/parseOne; quickedit integration`.

Glint cache keys include parser version, language, path and source; some production/development metadata is shared. Decode recovers panics but cache migrations, stale configuration, source overrides and concurrent writes need broad contract coverage.

Done when: Compare cached/uncached output after source, mtime, locale, format and config changes; production/development switch and overrides stay correct. Truncated/invalid cache safely misses and rebuilds. Bound stale cache growth if measurements justify it.

### C04 · `mp-b52.3.4` · P2 · Qualify routing, locales, base paths and version links together

Evidence: **coverage-gap**. Sources: `internal/site/build.go:localizePageLinks/languageLinks/sitePageURL; internal/site/nav_compile.go; internal/version/version.go`.

Standalone routes and multilingual/version tests exist, but interacting baseURL subpaths, missing translations, defaultLanguageAtRoot=false, encoded Unicode routes and deep mounted version links are a broad correctness surface.

Done when: Build a compact route matrix and run link/asset/fragment checks on every output. Canonical/alternate/search/nav/version URLs agree; no double encoding or false fallback. Cover spaces, CJK, root/deep paths and missing locale pages without adding routing features.

### C05 · `mp-b52.3.5` · P1 · Protect custom runtime CSS classes during production purging

Evidence: **confirmed**. Sources: `internal/site/csspurge.go:addRuntimeTokens/purgeUnusedCSSWithTokens; internal/site/build.go:custom CSS/static copy`.

Disposable CLI reproduction: a custom script adds audit-runtime-created to body; strict production build copies the script but removes its custom.css rule from mpress.css. Wails currently uses --no-purge-css for its imported homepage animation.

Done when: Reproduce a custom script-created class, escaped selectors and modern at-rules; retain required styling conservatively while keeping demonstrably unused built-in rules removable. Preserve custom CSS contract and test purged/unpurged pages by computed styles.

### C06 · `mp-b52.3.6` · P1 · Prove JavaScript and CSS minifier semantics with execution checks

Evidence: **source-risk**. Sources: `internal/site/minify_javascript.go; internal/site/minify_css.go; internal/site/assets.go; assets_test.go`.

assets.go describes lexical minifiers without renaming, but JS now mangles names using token heuristics. Template expressions are scanned as one token. Existing fuzzers check only panic freedom. Initial JS fuzz run terminated during minimization; saved seed replays successfully, so this is not a confirmed minifier hang.

Done when: Execute source and minified code against an engine for scoping, closures, template interpolation, shorthand keys, methods, labels, ASI, regex/division and modern syntax. Compare CSS computed styles. Reproduce any failure before classifying it; rerun isolated fuzz and preserve true regression seeds.

### C07 · `mp-b52.3.7` · P2 · Stress nested Markdown/MPD components and literal-code preservation

Evidence: **coverage-gap**. Sources: `internal/content/components.go; internal/components/components.go; internal/mpd parser/markdown; codefence_renderer.go`.

There are two fence-protection implementations and inside-out component processing; existing tests cover selected nesting and fuzz seeds, but malformed mixed directives and large repeated/nested blocks need adversarial coverage.

Done when: Exercise literal directives in fences/inline code, lists/quotes, duplicate headings, mixed native syntax, CRLF/BOM/Unicode and deep nesting. No silent loss or expanded literal examples; diagnostics remain deterministic and bounded. Round-trip supported MPD semantics.

### C08 · `mp-b52.3.8` · P2 · Run full import and round-trip qualification on pinned corpora

Evidence: **coverage-gap**. Sources: `internal/importer/starlight.go; internal/importer/markdown_mpd.go; internal/mpdconformance; projectconvert`.

Public migration fixtures pass. A pinned clean Wails native corpus builds 2,010 pages and passes links; the old private Starlight full-Wails test is skipped without WAILS_DOCS_SOURCE and no longer models the current native corpus.

Done when: Pin provenance for the existing native and Starlight fixtures; qualify imports, loss diagnostics, slugs, headings, assets, translations and build output. Keep private fixtures private. Make the full-corpus gate explicit instead of silently skipped.

### C09 · `mp-b52.3.9` · P2 · Resolve translation/navigation conformance and failure recovery

Evidence: **test-failure**. Sources: `internal/translate/engine.go/document.go/state.go; internal/navigation/navigation.go; private conformance_test.go`.

Private translation test fails expecting label: FR Start here in locale nav. Current code supports shared nav labels maps, so distinguish stale fixture assumptions from a product regression. Public provenance/target-safety tests pass.

Done when: Test translated nav rendering rather than old serialization assumptions. Qualify partial provider responses, cancellation, retries, human-edited targets and source/target/state save failures with fake providers. Preserve glossary/code/links and manual review states without calling paid services.


## Optimise measured build, output and runtime costs — `mp-b52.4`

### P01 · `mp-b52.4.1` · P1 · Maintain repeatable build, startup and size benchmark baselines

Evidence: **measured**. Sources: `internal/site/build_bench_test.go; internal/mpd/benchmark_test.go; audit evidence and workload harness`.

Baseline 201-page synthetic builds: cold 559/581/597ms, warm 541/517/501ms; cold ~484MB and warm ~248MB allocations. Docs 62 pages builds in ~199-243ms warm internally. Linux stripped binary is 50,131,209 bytes; version uses ~37.5MB peak RSS in one sample.

Done when: Provide pinned small/201-page/large multilingual/component-heavy workloads, command startup, wall/RSS and per-phase/allocation/size measurements. Record toolchain, CPU, options, warm/cold cache and repeated medians; keep correctness/determinism gates and comparisons on the same machine.

### P02 · `mp-b52.4.2` · P1 · Reduce knowledge HTML section parsing cost

Evidence: **measured**. Sources: `internal/knowledge/artifact.go:makeChunks/sectionsFromHTML/nodeText`.

8x warm-build CPU profile including priming: knowledge generation is 53.8% cumulative CPU; sectionsFromHTML is 34.0%. Full HTML DOM construction and Token allocations dominate part of this path.

Done when: Add a focused representative benchmark and preserve chunk text/headings/IDs/URLs. Compare streaming section extraction or reuse of already-parsed information against current output. Land only a measurable CPU/allocation improvement with multilingual and mixed HTML regressions.

### P03 · `mp-b52.4.3` · P1 · Reduce knowledge inverted-index token and posting allocations

Evidence: **measured**. Sources: `internal/knowledge/artifact.go:makeIndex/tokens; internal/knowledge/store.go`.

Allocation profile: makeIndex flat 187.7MB of 3,047MB across profiled run. It retokenizes page title, route and tags per chunk, creates pointer postings and a new route replacer per chunk, and sorts per-chunk term keys.

Done when: Benchmark component/long-page/CJK corpora; reuse invariant tokens/replacer and cheaper posting storage where safe. Preserve exact postings, deterministic bytes, relevance evaluation and filters. Show before/after allocations/time, without speculative rewrites.

### P04 · `mp-b52.4.4` · P2 · Reduce knowledge artifact serialization and peak buffers

Evidence: **measured**. Sources: `internal/knowledge/artifact.go:marshalArtifact/marshalCompactArtifact/generate; internal/knowledge/storage.go`.

Profile shows MarshalIndent/appendIndent and simultaneously retained page/chunk/index buffers. Knowledge generate accounts for ~52% cumulative allocation space; index alone is already compact.

Done when: Measure serialized bytes, peak RSS and runtime; assess compact page/chunk/manifest encoding or streaming without changing decoded schemas/digests unexpectedly. Preserve deterministic gzip and old loaders. Land a proven reduction and document compatibility decisions.

### P05 · `mp-b52.4.5` · P2 · Avoid rebuilding unchanged knowledge work on warm builds

Evidence: **source-hypothesis**. Sources: `internal/knowledge/artifact.go:generate; internal/site/parse.go; internal/site/build.go`.

Every warm build regenerates all sections/index/JSON even when parse entries hit and output is byte-identical. This is a measured dominant subsystem with an unmeasured caching hypothesis.

Done when: Measure reuse granularity, key every source/config/version/language dependency and prove cache invalidation/corruption fallback. Reuse only safe immutable intermediate work; retain complete output and deterministic manifests. If no justified win, close with measurements explaining why.

### P06 · `mp-b52.4.6` · P1 · Avoid the second full HTML pass for final asset hashes

Evidence: **measured**. Sources: `internal/site/build.go:rewriteAssetVersions/Build; internal/site/render_pages.go`.

Allocation profile: rewriteAssetVersions callback flat 276.6MB and cumulative 678.5MB, 22.3% of total run allocations. Build writes provisional hashes, then rereads/rewrites every HTML file after CSS optimisation.

Done when: Measure a bounded strategy for computing/filling final hashes without rereading and copying every page. Preserve custom/static HTML, version snapshots, cache URLs, purging and collector output. Prove allocations/I/O/time reduction and byte equivalence except intended hashes.

### P07 · `mp-b52.4.7` · P2 · Reduce shipped CSS and JavaScript for pages without optional widgets

Evidence: **measured**. Sources: `internal/site/build.go:asset assembly/csspurge; internal/site/theme_default.go; internal/site/reader_locale.go`.

Documentation bundle: CSS 120,963 bytes (21,388 gzip), JS 73,458 (21,066 gzip). One default runtime includes all interactive families and locale strings even for small pages; no browser coverage measurement yet.

Done when: Measure unused rules/functions on minimal and rich sites first. Remove provably unneeded output or initialization without breaking later dialogs/DOM-added classes, custom content or localization. Keep one Go binary, no Node build dependency, and all current features available.

### P08 · `mp-b52.4.8` · P2 · Reduce repeated browser search scoring work and handle index failures

Evidence: **source-hypothesis**. Sources: `internal/site/theme_default.go:loadIndex/normalise/scoreItem/runSearch; internal/site/build.go:writeSearchIndex`.

Each query renormalizes every item, headings and tags and builds/sorts all matches. loadIndex caches a failure as an empty array permanently. English docs index is 189,925 bytes (59,637 gzip). Large-index runtime cost is not yet measured.

Done when: Benchmark realistic 100/1,000/10,000-page multilingual indexes and input latency. Precompute invariant normalized data or bounded top results only if justified; preserve ranking/snippets/keyboard semantics. Test failed fetch/retry, stale queries and malformed indexes.

### P09 · `mp-b52.4.9` · P2 · Reduce command startup and binary size while preserving shipped features

Evidence: **measured**. Sources: `cmd/mpress; go.mod; internal/content/d2.go; dependency initialization`.

Stripped binary is 50.1MB and version command takes ~60ms/37.5MB RSS in one sample. Symbol inspection shows linked D2/Chroma/collation code; symbol/BSS size is not itself file-size attribution.

Done when: Measure multiple cold/warm command starts and package-init time; attribute on-disk binary growth. Reduce eager work or unnecessary linked paths/dependencies where supported by evidence. Preserve D2, translation, MCP and all existing commands; update licenses if dependencies change.

### P10 · `mp-b52.4.10` · P2 · Measure D2 repeated-diagram rendering and cache reuse

Evidence: **source-hypothesis**. Sources: `internal/content/d2.go:renderD2Diagram; internal/site/parse.go`.

Each rendered D2 fence builds a text ruler and compiles/renders with embedded fonts. Per-page parse cache helps unchanged pages but identical diagrams across pages are not separately reused. No focused D2 baseline yet.

Done when: Benchmark repeated/unique diagrams, many small graphs and parallel parses. Reuse only thread-safe immutable setup or keyed results if measured beneficial; preserve import isolation, SVG accessibility labels and timeout/error behaviour.

### P11 · `mp-b52.4.11` · P2 · Tune worker pools and asset transfer memory using actual workloads

Evidence: **source-hypothesis**. Sources: `internal/site/parse.go/render_pages.go; internal/deploy/cloudflare.go:collectAssets/uploadAssets; internal/deploy/netlify.go:zipOutput`.

Parser/render caps are both 8; serial finalization may dominate. Cloudflare hashes a base64 string and rereads uploads; Netlify builds a complete ZIP in a bytes.Buffer. These source patterns need workload measurements before changes.

Done when: Compare constrained and multi-core build settings, static-heavy sites and mock deployment uploads. Measure peak RSS/bytes/time; use bounded streaming/reuse where justified and preserve hashes, cancellation, retry/error handling and file limits. No real deployment needed.


### P12 · `mp-b52.4.12` · P1 · Fix measured v3 landing-page LCP delays without adding features

Added after the initial 48-task audit in response to the user's field report. Evidence: **user-reported / measured**. Sources: `internal/site/csspurge.go`, `internal/site/template.go`, Wails documentation builds and [wailsapp/wails#6198](https://github.com/wailsapp/wails/issues/6198).

The report includes v3 hero SVG LCP of 31,280ms and 4,635ms, plus legacy Docusaurus carousel/article images. Blocking v3 CSS is measurable; direct custom-domain requests also receive a Cloudflare managed challenge. The 31-second field outcome has not been reproduced or attributed.

Done when: repeated controlled before/after mobile/desktop and both-theme measurements establish improvements; regression checks preserve animation, theme selection, no-JavaScript rendering and carousel autoplay. Update implementation/evidence, CHANGELOG and Beads; qualify deployment and new field data separately. Source work and lab results are in [LCP.md](LCP.md); this task remains in progress pending publication and field qualification.

## Make qualification repeatable and close the audit — `mp-b52.5`

### V01 · `mp-b52.5.1` · P1 · Reconcile the private conformance suite with current contracts

Evidence: **test-failure**. Sources: `sibling mpress-tests at e0b55f42f6e8e50a511592b015dfecc3b3e249a8; release workflows`.

Four private top-level tests fail: launcher expectation, image focus, translation nav serialization, docs duplicate IDs/missing carousel (one test covers two findings). Current main merged site-specific contribution scripts; private expectations predate that merge. Private repo was inaccessible to taliesin-ai via gh in this session.

Done when: Separate real defects from obsolete assertions, fix defects, and update local private expectations to semantic current behaviour in an isolated branch if needed. Run entire suite with explicit candidate source/bin; retain privacy and flag genuine access limits without substituting weak public checks.

### V02 · `mp-b52.5.2` · P1 · Make browser and corpus gates execute explicitly in CI

Evidence: **coverage-gap**. Sources: `.github/workflows/ci.yml/release.yml; private browser_test.go/wails_qualification_test.go`.

CI calls private go test but does not install Chrome or axe-core and does not set a Wails corpus. Browser tests and Axe can skip; public qualification has no real swipe tests. PR forks cannot read private fixtures.

Done when: Add an explicit public browser regression gate and enforce browser/Axe prerequisites for owned release qualification. Install pinned tooling/test dependencies without adding production dependencies; fail required checks instead of skipping silently. Keep private corpus out of public logs.

### V03 · `mp-b52.5.3` · P2 · Run browser load and interaction traces with reproducible conditions

Evidence: **coverage-gap**. Sources: `generated minimal/docs/landing/blog pages; browser audit tooling`.

Current audit has real touch reproduction and payload sizes, but no Lighthouse score or field CWV claims. Browser performance needs cold cache, throttling, parse/execute, long-task and layout-shift measurements.

Done when: Collect repeatable load/interaction traces and resource timings on representative pages under stated device/network conditions. Report LCP/CLS/long tasks and search/overlay interactions with raw provenance; prioritize measured bottlenecks and record unmeasured field/device limitations.

### V04 · `mp-b52.5.4` · P2 · Expand meaningful navigation, CLI, export and version contract tests

Evidence: **coverage-gap**. Sources: `internal/navigation (0% coverage); internal/version (44.2%); internal/lighthouse (31.8%); cmd/mpress (45.1%); internal/exportzip`.

Coverage highlights untested package boundaries. Existing package tests cover ordinary output, but not every filesystem failure and command error. Percentages are prioritization signals, not goals to game.

Done when: Add contract tests for nav generation/autogeneration/localized labels, CLI exit/JSON diagnostics, clean safety, export overwrite/cancellation/permissions and version failures. Test real observable boundaries and errors, not implementation-shaped assertions.

### V05 · `mp-b52.5.5` · P2 · Run and harden bounded fuzz and differential parser/minifier checks

Evidence: **coverage-gap**. Sources: `internal/mpd/parser_test.go/json_test.go; internal/site/minify_fuzz_test.go; component fence tests`.

MPD 15s fuzz: 32,931 executions pass; JSON scanner 10s: 74,354 pass; CSS 10s: 44,192 pass. JS fuzz worker terminated during minimization, but saved seed replay passes. Panic-only fuzzing does not check semantic equivalence.

Done when: Rerun JS in isolation and classify timeout/resource effects before declaring a code bug. Add bounded semantic/round-trip invariants where meaningful and preserve actual failures. Run short deterministic CI seeds and documented longer local fuzz budgets.

### V06 · `mp-b52.5.6` · P2 · Verify deterministic artifacts, asset budgets and release binaries

Evidence: **coverage-gap**. Sources: `internal/site assets/build; internal/version; exportzip; .github/workflows/release.yml`.

62-page docs output hashes are identical over four rebuilds. Identical components expose duplicate IDs, and unordered metadata map output needs broader coverage. Only Linux amd64 was executed locally; release workflow has native smoke matrices.

Done when: Hash cold/warm/concurrent builds across representative corpora; cover archives, timestamps that are intentionally variable, D2 and metadata order. Record raw/gzip/Brotli output and binary budgets with provenance. Build supported targets and run available native smoke tests; report unavailable execution honestly.

### V07 · `mp-b52.5.7` · P2 · Reduce measured maintenance risks only where they affect fixes

Evidence: **structural**. Sources: `internal/dev/assets/devbar.js:showConfig; internal/dev/server.go; internal/translate/engine.go; cmd/mpress/main.go; internal/site/build.go`.

Ripwire maintenance snapshot: devbar.js churn 3/ccx 2221 (showConfig ccx 556), translate engine churn 13/ccx 435, Build ccx 185. These are structural risks, not runtime bottlenecks; clone scan includes test and idiomatic duplication.

Done when: When fixing a hotspot, extract a bounded pure seam or reuse existing validation/I/O helpers where it lowers actual complexity and avoids duplicated policy. Keep contracts unchanged, run affected tests and quality-delta. Do not undertake broad aesthetic rewrites or delete unverified dead code.

### V08 · `mp-b52.5.8` · P1 · Complete final feature-freeze qualification and publish measured results

Evidence: **coverage-gap**. Sources: `all audit tasks; CHANGELOG.md; audit evidence; Beads ledger`.

The goal is completion of the current-feature audit, not a claim of being globally fastest. Fixes and optimizations must be reviewable and tied to evidence.

Done when: All child tasks resolved with evidence; race/vet/security/public/browser/strict/link/corpus gates pass or explicit external blockers remain open. Repeat baseline workloads and report measured improvements and regressions, update CHANGELOG and Beads, commit as taliesin-ai. Never close this gate with required work outstanding.
