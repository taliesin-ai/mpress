# Conversion checkpoint review

Measured against `e8eac75` before committing. No quality acknowledgements were
added. The delta exits 2: 16 gates, two minor rows and 12 new-symbol findings.

- Eight gates pair the small checked YAML marshal helper with unrelated checked
  open/read/stat helpers. They share error-return syntax, not serialization or IO
  semantics. Config Save and SaveRoot share one serializer rather than duplicate it.
- Three writer gates match wrappers around the same existing exclusive, pinned
  atomic writer. Cache intentionally omits durable sync; conversion preserves
  permissions and syncs. Sharing the implementation preserves those contracts.
- Two request-wrapper gates match action-specific forwarding to the existing shared
  authenticated request fixture. Two larger test gates match existing negative
  fixture setup/assertions. Those helpers are already shared; the conversion tests
  additionally assert that no converted target was published.
- One gate records self churn in WriteAtomic. Two minor rows record ambient churn
  in the CLI and authoring dispatcher. Their changes distinguish explicit CLI
  authority and borrow the server's pinned root.
- New-symbol complexity/length rows are moved existing conversion implementations:
  original Run ccx 64/154 lines becomes runConversion ccx 62/149 lines plus checked
  ownership wrappers; original state planning ccx 33 becomes 34 with a null-root
  guard. This is not a complexity-cleanup claim. The eight dead-code rows are
  executed Go tests. The remaining new clone is two different root-ownership and
  configuration-serialization tests, both executed.

The unused ambient PlanConversionStates wrapper was removed after confirming that
all production callers use the rooted planner. Impact review covers CLI conversion,
authoring, translation sidecars, configuration serialization and shared filesystem
writer contracts. Full race/vet, focused tests and cross-compilation cover them.
Cross-compilation is not native Windows/macOS qualification. Whole-conversion
rollback remains S06; per-file atomic replacement does not satisfy that criterion.
