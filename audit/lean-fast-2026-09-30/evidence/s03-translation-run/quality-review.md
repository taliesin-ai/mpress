# Translation execution review

RunRoot borrows the server root on a shallow Engine copy. The original Engine
never retains that handle; standalone Run opens and closes its own root after
workers join. Ownership/reuse and safe internal aliases have regression tests.
Planning and execution share rooted source, target, state and provider-input IO.
Existing manual/stale/migration/relocation and provider contracts remain covered.

The final working-tree quality command exits 2: eight gates remain, without
acknowledgements. Two clone pairs and their two reuse rows are the 65-token
glossary reader wrappers and 38-token optional reader wrappers. They select
different filesystem authority for separately scoped operations; glossary
parsing is shared. Four self-churn rows are findStateByPageKeyRoot,
handleTranslations, readTranslationTarget and translateFile, reflecting successive
P0 boundary fixes. These are reviewed findings, not a green quality gate.

Three minor findings are Run complexity 40 to 41, length 118 to 123 and Engine
ambient churn. Initial moved-method and cleanup complexity findings were reduced
by retaining Run's implementation, extracting relocation cleanup and removing the
dead ambient page-key finder. Ten new dead-code heuristics are executed Go tests,
a benchmark and interface implementations; recorded runs exercise those entries.

Target and state writes are individually atomic, not an atomic pair. A rejected
state-parent replacement can occur after an internal target write. C09/S06 recovery
qualification remains open. External sentinels survive provider-time target/state
parent rearrangement. Remaining audit/refinement/comparison delegates are not
covered by this checkpoint. Native execution is Linux; platform compilation does
not establish native Windows/macOS behavior.
