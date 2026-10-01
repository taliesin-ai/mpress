# Translation check review

Check owns or borrows one project root and uses a local Engine copy for discovery,
coverage, audit and exception hashing. The original Engine never retains that
handle. Source enumeration now honors the borrowed root; helper readers use the
same manager, so a whole check does not reopen the project through an ambient path.
Missing navigation still produces ordinary coverage diagnostics; unsafe/unreadable
navigation errors are surfaced. Exceptions preserve their accepted codes, reason,
exact source/target hashes and lexical traversal guards. Existing coverage,
exclusion, stale-hash and warning tests pass. Check remains read-only and invokes
no provider. Borrowed ownership/reuse tests now include Check.

Original helper regressions reproduce two external hash reads accepted as valid
exceptions and one external coverage read reported as valid. These are core reader
contract tests, not three HTTP requests. Both safe internal alias controls pass.
Each native symlink fixture probes support first; subsequent fixture failures are
fatal rather than skipped.

The final quality delta retains three gates: a 53-token checked root-lifecycle
wrapper match between randomComparisonSample and sourceFiles, and two self-churn
rows for successive source-enumeration and ownership qualifications. Introducing
a generic callback solely for that wrapper would obscure return types and purpose.
The one minor row is acceptedFinding complexity 15 to 16 for checking the root-open
error. Three new dead-code heuristic rows are Go tests executed by recorded runs.
No acknowledgements are added; exit 2 is retained rather than reported green.

This checkpoint does not qualify migration/conversion, deployment or remaining
metadata IO, and does not close the P0 authoring issue or broader audit gates.
