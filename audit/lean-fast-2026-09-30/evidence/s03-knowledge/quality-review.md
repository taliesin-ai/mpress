# Quality review

The new rooted loader separates IO/aggregation from pure version sorting. That
removes its initial complexity-17 finding. The remaining twelve major heuristic
rows are six pairs each reported as duplication and clone-of-helper:

- `Load | LoadAll`: both own an opened site root, close it on all returns and call
  different loaders. Their short lifecycle duplication is intentional; a callback
  ownership framework would obscure these two straightforward entry points.
- `Load | ReadDir`, `Load | Stat`, `Load | loadKnowledge`: these share a short
  checked-result/return skeleton. The contracts differ: explicit site opening,
  rooted filesystem enumeration/stat, and borrowing/pinning an authoring output
  root. They compose the actual `projectfs` primitive already; combining the
  wrappers would mix responsibilities and ownership.
- `Load | marshalCompactArtifact`: the normalized error-return shape matches
  unrelated JSON marshaling. No directory operation can reuse a JSON encoder.
- `devSymlink | knowledgeSymlink`: both test packages need to distinguish absent
  native symlink support from a fixture failure. Go test-only private helpers
  cannot be imported across packages; duplicating this small fixture guard keeps
  production packages free of test fixture APIs/dependencies.

The two new generated knowledge fixture helpers use the existing Generate API
in distinct test packages; one controls the external manifest title, the other
checks content/search equivalence. Test entry points are executed by the Go test
runner, despite dead-code heuristics. No blanket acknowledgements or zero-debt
claim were added.
