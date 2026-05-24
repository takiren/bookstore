.PHONY: build test test-unit test-integration generate-sqlc

build:
	go build ./...

test: test-unit

test-unit:
	go test ./...

test-integration:
	go test -tags integration -v ./...

generate-sqlc:
	sqlc generate
