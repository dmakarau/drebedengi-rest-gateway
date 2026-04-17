# drebedengi-rest

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Chi](https://img.shields.io/badge/Chi-v5-purple)
![SOAP](https://img.shields.io/badge/SOAP-XML-orange)
![REST](https://img.shields.io/badge/REST-JSON-green)
![Status](https://img.shields.io/badge/Status-WIP-yellow)
![Purpose](https://img.shields.io/badge/Purpose-Personal_Project-blue)

A REST/JSON gateway for the [drebedengi.ru](https://www.drebedengi.ru) personal finance service. Drebedengi exposes a SOAP/XML API; this project wraps it in a clean HTTP+JSON interface suitable for mobile apps, scripts, or anything that would rather not deal with SOAP.

Written in Go. No external SOAP libraries -- requests are built with `strings.Builder`, responses are parsed with `encoding/xml`.

## Why

Drebedengi is a solid personal finance tracker, but the only integration path is a SOAP API from the PHP era. Consuming SOAP from a mobile app (or anything modern) is painful. This service sits in front of the SOAP endpoint and speaks REST+JSON, so your client code stays simple.

It also keeps your API credentials server-side instead of embedding them in the app binary.

## Status

Work in progress. The following endpoints are implemented and tested against the live API:

### Account

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/account/status` | Payment/access status (1 = OK) |
| GET | `/api/v1/account/revision` | Current server revision number |
| GET | `/api/v1/account/expire` | Subscription expiry date (YYYY-MM-DD, or "0") |
| GET | `/api/v1/account/subscription` | Premium status (1 = active, 0 = inactive) |
| GET | `/api/v1/account/access` | Access level (0 = full, 1 = limited) |
| GET | `/api/v1/account/userid` | Internal user ID for the configured login |

### Planned

- Balance (grouped by place/currency)
- Records (expenses, incomes, moves, currency changes) with filtering
- Categories, Sources, Places, Currencies, Tags -- full CRUD
- Change log (sync support via revision tracking)
- Purchases/accumulators, checks

See [PLAN.md](PLAN.md) for the full implementation roadmap.

## Getting started

### Prerequisites

- Go 1.21+
- A drebedengi.ru account (free registration at [drebedengi.ru](https://www.drebedengi.ru))
- An API ID (request one from drebedengi, or use `demo_api` for development)

### Setup

Clone the repo and copy the example environment file:

```sh
git clone https://github.com/youruser/drebedengi-rest.git
cd drebedengi-rest
cp .env.example .env
```

Edit `.env` with your credentials:

```
DD_API_ID=demo_api
DD_LOGIN=demo@example.com
DD_PASSWORD=demo
DD_SOAP_URL=http://www.drebedengi.ru/soap/
PORT=8080
```

The demo credentials (`demo_api` / `demo@example.com` / `demo`) are provided by drebedengi for development and testing. They give access to a shared demo account with sample data.

### Run

```sh
go run .
```

Or use the Makefile:

```sh
make run
```

The server starts on the configured port (default 8080).

### Test it

```sh
curl http://localhost:8080/api/v1/account/status
# {"status":1}

curl http://localhost:8080/api/v1/account/revision
# {"revision":1234908678}
```

### Build

```sh
make build
# produces ./drebedengi-rest binary
```

## Configuration

All configuration is via environment variables. A `.env` file in the project root is loaded automatically (via godotenv).

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DD_API_ID` | yes | -- | Your drebedengi API key |
| `DD_LOGIN` | yes | -- | Account email |
| `DD_PASSWORD` | yes | -- | Account password |
| `DD_SOAP_URL` | no | `http://www.drebedengi.ru/soap/` | SOAP endpoint URL |
| `PORT` | no | `8080` | HTTP listen port |

## Project structure

```
main.go                          Entry point, router setup
internal/
    config/config.go             Environment-based configuration
    respond/respond.go           JSON response helpers
    soap/
        caller.go                Caller interface (for mocking)
        client.go                HTTP client, implements Caller
        envelope.go              SOAP envelope building and response parsing
        encoding.go              Value encoding (simple types, ns2:Map, SOAP-ENC:Array)
        account.go               Account-related SOAP method wrappers
        types_account.go         XML response structs for account methods
    handlers/
        handlers.go              Handler struct, router wiring
        account.go               Account HTTP handlers
    models/
        models.go                JSON response structs
```

## How it works

1. An HTTP request hits the chi router
2. The handler builds SOAP parameters and calls the appropriate wrapper function
3. The wrapper uses the `Caller` interface to build a SOAP XML envelope, POST it to drebedengi, and parse the response
4. The handler converts the SOAP response into a JSON model and writes it back

The SOAP XML is built by hand because drebedengi uses RPC/encoded style with Apache XML-SOAP Map types (`ns2:Map`) and explicit `xsi:type` attributes on every value. Standard Go XML marshalling does not handle this well.

## Drebedengi API notes

The upstream API is documented in their [WSDL](https://www.drebedengi.ru/soap/dd.wsdl) and [API integration page](https://www.drebedengi.ru/faq11.html). A few things worth knowing:

- All monetary sums are integers in hundredths (e.g. 1000 = 10.00 in the currency)
- Operation types: 2 = income, 3 = expense ("waste"), 4 = move between places, 5 = currency change
- The API uses "waste" to mean "expense" throughout
- Errors come back as SOAP Faults, which this service translates to appropriate HTTP status codes

## License

MIT
