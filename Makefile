.DEFAULT_GOAL := build

COVERAGE_MIN := 80
COVERAGE_FILE := coverage.out

.PHONY: clean lint lint-fix build test test-coverage

clean:
	go clean ./...

lint: clean
	golangci-lint run ./...

lint-fix: clean
	golangci-lint run --fix ./...

build: lint-fix
	go build -o bin/gendiff ./cmd/gendiff

test-coverage:
	go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) code
	@go tool cover -func=$(COVERAGE_FILE) | awk -v min=$(COVERAGE_MIN) \
		'/^total:/ { printf "Coverage: %s\n", $$3; if ($$3 + 0 < min) { printf "FAIL: coverage %s is below %d%%\n", $$3, min; exit 1 } }'
