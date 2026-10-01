# State migration review

Measured against `83faa1d`, without quality acknowledgements. Three gates remain:
MigrateState complexity 38 → 41, length 67 → 81, and an 80-token match between a
native-symlink probe and an unrelated write fixture. The method adds checked root
acquisition/cleanup and absolute snapshot normalization; triplet alignment and
approval rules are unchanged. The probe does not share write-test semantics.

Five new-symbol findings cover three executed tests and two fixture-driven test
bodies. Their guard branches qualify external sentinels, pre-write destination
preservation, original approval evidence and caller ownership. No production dead
helper or ordinary performance improvement is claimed.

Impact review covers standalone and borrowed migration, separate explicit snapshot
authority, config parsing, source selection, triplet readers and state persistence.
The genuine original migration file reproduces exactly nine failures and one
external overwrite, with all nine internal-alias controls passing. Full public
race/vet and existing migration/idempotency/dry-run/alignment tests pass. Latest
translation tests compile for Windows/macOS; Linux is the native platform.

Old review/run/audit replays now include their matching old migration dependency;
selected original handlers and failure sets are unchanged. Later replays and the
historical direct replay retain their exact failures. Multi-sidecar transaction
recovery remains separate scope; planning precedes the first write, but individual
state writes do not make the entire migration one transaction.
