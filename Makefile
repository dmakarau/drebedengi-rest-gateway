.PHONY: run build test test-int lint

run:
	go run ./...

build:
	go build -o drebedengi-rest ./...

test:
	go test ./...

test-int:
	go test -tags=integration ./...

lint:
	golangci-lint run
