# Review of quality heuristics

The initially duplicated native swipe/retry logic in contribution reachability
and accessibility touch checks now uses `browserActions.reachByTouch`. Each caller
still checks its own scroll owner on every probe; neither assertion was relaxed.

Three remaining major heuristics compare the small test-only `browserQuote`
JSON-string conversion with production `readBody`/`writeJSONString`. The shared
normalized shape is JSON marshaling/error/return syntax, with different contracts
and packages. Production IO functions are unsuitable for generating JavaScript
string literals in the separate browser-test module. The existing `json.Marshal`
primitive is reused; no production dependency or extra JSON implementation exists.
A string has no marshal error, and the ignored error here cannot hide an IO failure.
No blanket acknowledgements were added. The Go test entry point is invoked by the
Go runner and verified in the recorded native test run; it is not dead code.
Ambient churn on the CSS constant records recent independent correctness fixes.
