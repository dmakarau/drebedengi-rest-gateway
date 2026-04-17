Scaffold a new REST endpoint following the project's 4-layer pattern.

**Required argument:** `<domain> <SOAPMethodName>` — e.g. `records GetRecordList`

**Layers to create/extend (in order):**

1. **XML response type** — add to `internal/soap/types_<domain>.go` (create file if needed):
   ```go
   type <camelMethod>Response struct {
       XMLName xml.Name `xml:"<SOAPMethodName>Response"`
       Return  <type>   `xml:"<SOAPMethodName>Return"`
   }
   ```

2. **SOAP wrapper** — add to `internal/soap/<domain>.go` (create file if needed), following `internal/soap/account.go` exactly:
   - Accept `(c Caller, apiId, login, pass string)` — no extra params unless the method requires them
   - Wrap parsing errors with `fmt.Errorf("parsing <method> response: %w", err)`

3. **JSON model** — add struct to `internal/models/models.go` with `json:"snake_case"` tags

4. **HTTP handler** — add method to `internal/handlers/<domain>.go` following `internal/handlers/account.go`:
   - Call SOAP wrapper, return `502` on error via `respond.Error`, return `200` with model via `respond.JSON`

5. **Route registration** — add `r.Get("/<path>", h.<HandlerMethod>)` in `internal/handlers/handlers.go`

6. **Tests** — add to `internal/handlers/<domain>_test.go` (create file if needed):
   - Happy path: mock returns valid XML, assert JSON field value
   - SOAP error path: mock returns error, assert `502` and `"error"` key in body

Before writing, read the existing files for the domain if they exist.

$ARGUMENTS
