Run the test suite and surface coverage gaps.

**Steps:**
1. Run `go test ./... -count=1 -cover -coverprofile=/tmp/drebedengi-coverage.out 2>&1`
2. Run `go tool cover -func=/tmp/drebedengi-coverage.out | sort -k3 -n` to rank functions by coverage
3. List functions at 0% coverage and explain what cases are missing
4. If any test failed, show the full failure output and identify the root cause

**If $ARGUMENTS specifies a scope:**
- Single package: `go test ./internal/<package>/... -count=1 -v -cover`
- Single test: `go test ./... -run <TestName> -v -count=1`
- With race detector: add `-race` flag

$ARGUMENTS
