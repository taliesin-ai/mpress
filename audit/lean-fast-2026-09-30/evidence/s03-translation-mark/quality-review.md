# Translation review quality and contract review

The authenticated handler borrows the server's existing pinned project root.
Standalone Mark owns and closes one project root. Rooted discovery retains
logical paths and shares language/file selection with the original contract;
all Engine source enumeration now uses that rooted discovery. State decoding
and page-key lookup share implementations with existing translation operations.
Those other operations retain their separate, still-open IO qualification.

The initial moved review method measured complexity 31 and length 100. Selection,
rooted document/state loading and segment approval are now separate routines;
the initial negative fixture's setup was also extracted. Final measurements add
no complexity/length finding. Existing reordered/relocated/human-edited/stale
approval tests pass, along with five real authenticated boundary failures and
eight persisted safe-alias controls. State schema and serialization remain the
same; writes reuse projectfs's exclusive temporary-file/atomic-write contract.
Cleanup of a relocated owned sidecar now reports errors rather than hiding them.

Thirteen quality gate rows remain seven short clone pairs and six reuse rows:

- Load/sourceFiles, Mark/loadKnowledge, Mark/sourceFiles and loadKnowledge/sourceFiles
  are checked root-open/defer-close/forward wrappers across distinct roots and
  return types. Their 41–48 token matches are retained rather than introducing
  a generic callback abstraction solely for this lifecycle boilerplate.
- loadState/loadStateRoot are 30-token reader wrappers that forward to the same
  decoder. They deliberately choose ambient or rooted IO for separately scoped
  operations; removing the root variant would undo the safety change.
- marshalArtifact/saveStateRoot share 45 tokens of checked indented JSON encoding
  in different packages. Knowledge's encoder is private to its package, and
  making translation depend on knowledge to reuse this encoding shape would
  introduce unrelated coupling. Review state schema handling belongs here.
- readIdentityFixture/readTranslationFixture share 33 tokens of checked test
  file reading in different package compilation units. The knowledge fixture
  helper is inaccessible from dev; no cross-package test-support API is added.

Three dead-code entries are Go tests executed by the recorded runs. The
handleTranslations churn row is minor and reflects the borrowed-root call.
No acknowledgements are added, and the quality command's exit 2 is retained.
