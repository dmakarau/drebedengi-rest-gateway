Perform a thorough senior-engineer-level code review of this branch. Use deep reasoning — treat this like a review before merging into production.

**Steps:**
1. `git log main...HEAD --oneline` — understand what's on the branch
2. `git diff main...HEAD` — read every changed line
3. For each changed file, also read its test file and any files it imports

**Evaluate across these dimensions:**

**Correctness**
- Logic errors, unchecked error returns, off-by-one, nil dereferences
- SOAP response parsing: does the XML struct field name match the actual element name?

**Go idioms**
- Errors wrapped with `%w` (not `%v`) when they'll be inspected by callers
- `context.Context` propagated to `http.NewRequestWithContext` (not `http.NewRequest`)
- Interfaces defined at the point of use (consumer package), not the implementation package
- Unexported types for internal-only structs

**Error handling**
- All SOAP errors translated to appropriate HTTP status codes (4xx vs 5xx)
- SOAP Faults surfaced in the JSON `"error"` field

**Test quality**
- New code paths covered by `mockCaller`-based handler tests and/or soap package unit tests
- Table-driven tests where there are multiple similar cases
- No testing implementation details — test behavior through the public API

**Architecture fit**
- Changes follow handler → soap wrapper → Caller → envelope layering
- New SOAP methods have their XML response types in `types_<domain>.go`
- No business logic leaking into handlers; no HTTP concerns in the soap package

**Security**
- No credentials logged or returned in responses
- XML inputs escaped via `xmlEscape` before being embedded in envelopes

Output findings grouped by: **Critical** (must fix) → **Warning** (should fix) → **Suggestion** (optional improvement).

$ARGUMENTS
