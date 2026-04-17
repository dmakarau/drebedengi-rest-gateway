# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Setup (after clone)

```sh
git config core.hooksPath .githooks   # activates commit message cleanup hook
cp .env.example .env                  # then fill in credentials
```

## Commands

```sh
make run       # go run ./...
make build     # go build -o drebedengi-rest ./...
make test      # go test ./...
make test-int  # go test -tags=integration ./...  (hits live API)
make lint      # golangci-lint run
```

Run a single test:
```sh
go test ./internal/soap/... -run TestBuildEnvelope
go test ./internal/handlers/... -run TestGetAccessStatus -v
go test ./... -count=1 -race -cover   # full suite with race detector
```

## Custom slash commands

- `/commit` — stage review + concise conventional commit (no AI attribution)
- `/review` — deep review of branch changes (correctness, idioms, tests, security)
- `/test [pkg|TestName]` — run tests and surface coverage gaps
- `/new-endpoint <domain> <SOAPMethod>` — scaffold all 4 layers for a new endpoint

## Configuration

Copy `.env.example` to `.env`. All config is via environment variables loaded by `godotenv`:

| Variable | Description |
|---|---|
| `DD_API_ID` | drebedengi API key (`demo_api` for dev) |
| `DD_LOGIN` | Account email |
| `DD_PASSWORD` | Account password |
| `DD_SOAP_URL` | SOAP endpoint (default: `http://www.drebedengi.ru/soap/`) |
| `PORT` | HTTP listen port (default: `8080`) |

Demo credentials: `demo_api` / `demo@example.com` / `demo` — shared sample account provided by drebedengi.

## Architecture

This is a REST/JSON gateway over the drebedengi.ru SOAP/XML API. Request flow:

```
HTTP request → chi router → Handler → soap.Account function → soap.Client.Call
                                                              ↓
                                              BuildEnvelope (manual XML)
                                              POST to SOAP endpoint
                                              ParseResponse → xml.Unmarshal
                                              ↓
                              Handler converts result → JSON response
```

**Why manual XML:** drebedengi uses RPC/encoded SOAP with Apache `ns2:Map` types and explicit `xsi:type` attributes on every value. Standard `encoding/xml` marshalling cannot handle this, so envelopes are built with `strings.Builder` in `internal/soap/encoding.go`.

### Key interfaces and patterns

- `soap.Caller` interface (`internal/soap/caller.go`) — accepts `context.Context` as first param; the only dependency handlers take on the SOAP layer. Enables handler unit tests via `mockCaller` without hitting the network.
- `soap.Param` — typed parameter passed to `BuildEnvelope`. The `Value any` field drives encoding dispatch in `encodeParam`.
- Every SOAP wrapper function follows the pattern: `func GetXxx(ctx context.Context, c Caller) (T, error)` — auth params are prepended automatically by `Client.Call`, so wrappers pass `nil` for non-auth params.
- Responses that return an `ns2:Map` array (e.g. `getBalance`) cannot be unmarshalled directly into typed structs because `encoding/xml` doesn't support attribute-driven dispatch. Instead, unmarshal into `[]balanceMapEntry` (generic key-value pairs) and switch on `kv.Key` to populate a typed struct. See `internal/soap/balance.go` for the pattern.
- SOAP errors are mapped to HTTP status by `soap.Fault.HTTPStatus()`: Client+auth faults → 401, other Client faults → 400, Server faults → 502. The `soapErr` helper in `handlers.go` does the type assertion.

### Adding a new endpoint

1. Add XML response struct(s) to `internal/soap/types_<domain>.go`
2. Add a wrapper function to `internal/soap/<domain>.go` following the existing pattern
3. Add an HTTP handler method to `internal/handlers/<domain>.go`
4. Register the route in `internal/handlers/handlers.go`

## Go standards for this project

**Error handling**
- Wrap errors with `%w` (not `%v`) whenever the caller might inspect the type or unwrap it
- Pass `r.Context()` to all SOAP wrapper functions; `Client.Call` forwards it to `http.NewRequestWithContext`
- Check all `error` returns; `respond.JSON` logs encode errors but cannot change the status code after `WriteHeader` was called

**Naming**
- SOAP wrapper functions are exported and named after the SOAP method (`GetAccessStatus`, not `FetchStatus`)
- XML response structs are unexported (`getAccessStatusResponse`), named `<camelMethod>Response`
- JSON model structs are exported and in `internal/models/`

**Testing**
- Handler tests use `mockCaller` — inject via `Handler.SOAP`; never start a real HTTP server in unit tests
- SOAP package tests (`envelope_test.go`, `encoding_test.go`) test XML output as strings; use `strings.Contains` for targeted assertions rather than full string equality when the surrounding envelope is not under test
- Prefer table-driven tests (`[]struct{ name, input, want }`) when there are three or more similar cases

**Commit messages**
- Format: `<type>(<scope>): <description>` — max 72 chars, imperative mood, lowercase after colon
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `ci`, `chore`
- Scope = package name when the change is confined to one package
- No period at end; no "Co-Authored-By"; no AI attribution
- Use `/commit` command to draft and create commits

### Drebedengi API notes

- All monetary sums are integers in hundredths (1000 = 10.00)
- Operation types: 2 = income, 3 = expense ("waste"), 4 = move, 5 = currency change
- SOAP Faults are translated to HTTP errors by `soap.Client`
