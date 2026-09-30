#!/usr/bin/env python3
"""Replay negative confinement checks against a06977e production paths."""
import json
import pathlib
import re
import subprocess
import tempfile

root = pathlib.Path(__file__).resolve().parents[4]
baseline = "a06977e"
with tempfile.TemporaryDirectory(prefix="mpress-s03-replay-") as work:
    temporary = pathlib.Path(work)
    def original(names):
        replacement = {}
        for name in names:
            target = temporary / name.replace("/", "-")
            target.write_bytes(subprocess.check_output(["git", "show", baseline + ":" + name], cwd=root))
            replacement[str(root / name)] = str(target)
        return replacement

    cases = [
        ("direct", original(["internal/dev/server.go", "internal/dev/mcp.go", "internal/dev/translation_models.go", "internal/config/output.go"]), "./internal/dev", "^Test(AuthoringRESTRejectsExternalSymlinkReadsAndWrites|AuthoringMCPRejectsExternalSymlinkPaths|AuthoringUploadRejectsExternalStaticAndContentParents|AuthoringBackupsRejectExternalParentWithoutChangingSource|NewServerPreservesUnsafePreviewPools|AuthoringConfigurationAndTreeDoNotFollowExternalLinks)$"),
        ("delegated", original(["internal/site/build.go", "internal/site/parse.go"]), "./internal/dev", "^TestAuthoring(RebuildDoesNotWriteCacheThroughExternalParent|PreviewRejectsOtherExternalInputs)$"),
        ("output", original(["internal/config/output.go"]), "./internal/config", "TestRemoveOutputProtectsConfiguredDataAndAncestors/(navigation|localized|translation_glossary|translation_style_guide|contributor_guide)"),
    ]
    # The old Server owns no root handle/Close method. Remove new handle cleanup
    # and the one new cleanup-method test from the baseline's compilation unit.
    current = (root / "internal/dev/confinement_test.go").read_text()
    current = re.sub(r"^\s*t.Cleanup\(func\(\) \{ _ = server.Close\(\) \}\)\n", "\n", current, flags=re.M)
    current = re.sub(r"\nfunc TestAuthoringPreviewCleanupRechecksRedirectedPool\(.*?(?=\nfunc )", "", current, flags=re.S)
    adapted = temporary / "direct-confinement_test.go"
    adapted.write_text(current)
    cases[0][1][str(root / "internal/dev/confinement_test.go")] = str(adapted)
    for label, replacement, package, pattern in cases:
        overlay = temporary / (label + ".json")
        overlay.write_text(json.dumps({"Replace": replacement}))
        print("Baseline replay:", label, flush=True)
        result = subprocess.run(["go", "test", "-overlay=" + str(overlay), package, "-run", pattern, "-count=1"], cwd=root, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        print(result.stdout, end="", flush=True)
        failures = re.findall(r"^--- FAIL: (\S+)", result.stdout, re.M)
        if len(failures) != {"direct": 6, "delegated": 2, "output": 1}[label] or "[build failed]" in result.stdout:
            raise SystemExit("Baseline replay did not reproduce the expected entry points: " + label)
        if result.returncode == 0:
            raise SystemExit("Expected failing-before regressions did not fail: " + label)
