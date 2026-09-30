# Reproducing the audit

Run the source-selected Go toolchain from the repository root:

```sh
go test -race ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
go build -trimpath -ldflags='-s -w' -o /tmp/mpress-audit-bin ./cmd/mpress
/tmp/mpress-audit-bin build --strict --json
/tmp/mpress-audit-bin check
go test ./internal/site -run '^$' -bench '^BenchmarkBuild(Cold|Warm)$' -benchmem -benchtime=3x -count=3
```

The `BenchmarkBuild*` workload writes 200 numbered pages plus a home page.
Its cold case clears the M-Press parse cache and output; the warm case clears
output after priming the parse cache. Both enable asset minification/CSS purge.

For the touch reproduction, this directory is a separate Go module so browser
tooling does not enter the production module or binary. It serves already-built
site output through a temporary local HTTP server and uses installed Chrome:

```sh
cd audit/lean-fast-2026-09-30/tools
go run . /absolute/path/to/mpress/site > mobile-results.json
```

`mobile.go` compares real touch swipes before and after an in-browser diagnostic
CSS injection. The injection does not edit repository source or site output.
The original failing baseline is preserved in `../evidence/mobile.json`.
The browser currently defaults to `/usr/bin/google-chrome`; adjust that local
path for another installation. This harness is Chromium-only.

Public audit evidence is committed. Raw private conformance logs and private
fixtures are deliberately kept in the private test environment. To run that
suite, select the exact candidate explicitly with `MPRESS_SOURCE`/`MPRESS_BIN`
and supply installed browser/Axe prerequisites. A skipped check is not a pass.

To recreate the large workload, archive the pinned Wails `docs/mpress/` revision
listed in `AUDIT.md` into a disposable directory, then build there with
`--strict --no-purge-css` and run `check`. Do not modify the user's Wails checkout
or substitute a different revision without recording its provenance.
