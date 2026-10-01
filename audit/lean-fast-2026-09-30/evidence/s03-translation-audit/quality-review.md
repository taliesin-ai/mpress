# Audit and refinement review

The authenticated handler borrows its pinned project root on an Engine copy.
Static audit, independent review, refinement and final re-audit share it.
Standalone operations own one root. Source enumeration/reads, target reads,
state reads/writes and provider input reads all use that authority. RunRoot
reuses the same copy operation. Ownership/reuse tests cover all three methods.
The obsolete ambient atomic writer is removed; two initially separate public
root wrappers were replaced by one setup method. No dependency/product feature
was added. Existing clean-segment, damaged-placeholder, formatting, navigation,
human-edited and stale-state contracts remain covered by the translation suite.

Shared fixtures install provider-input links and perform callback-time directory
swaps. The third copy of provider-input rejection checks now uses one case runner,
while operation-specific call assertions and the original named tests are retained.

The final quality delta has 22 gates, six minor rows and seven new-symbol rows.
No acknowledgements are added; exit 2 is retained. Eight short clone pairs
(35–47 tokens) plus eight reuse rows compare checked RunRoot forwarding to
unrelated root wrappers/parsers/JSON encoders. They have different authority and
return types; no generic callback abstraction is added for that syntax. Two churn
rows concern Engine and Run after successive P0 fixes.

Four other gate pairs are retained: plan/review HTTP rejection (157 tokens) tests
different endpoints, refinement/run parent swaps (173 tokens) exercise different
provider operations and consequences, and two 23–25-token table wrappers forward
to the shared case runner. The new-symbol clone between the two new wrappers is
25 tokens. Six dead-code rows are five executed Go tests and the reviewer interface
implementation. Minor changes are refinement complexity 98 to 99/length 190 to
195, handler complexity 100 to 102/length 275 to 279 and two ambient churn rows.
This is explicit review, not a green quality or final audit gate.

Target and sidecar remain individually atomic rather than a transactional pair.
Global per-user provider credentials remain intentionally outside the project.
Comparison/check/conversion/deploy and wider S07 recovery remain unqualified.
