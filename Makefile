.DEFAULT_GOAL := build

BIN_DIR := $(CURDIR)/bin
export GOBIN := $(BIN_DIR)

GOLANGCI_LINT_VERSION := v2.12.2
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

COVERAGE_MIN := 80
COVERAGE_FILE := coverage.out

.PHONY: clean lint lint-fix build test test-coverage

clean:
	go clean ./...

lint: clean
	$(GOLANGCI_LINT) run ./...

lint-fix: clean
	$(GOLANGCI_LINT) run --fix ./...

build: lint-fix
	go build -o bin/gendiff ./cmd/gendiff

test-coverage:
	go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) code
	@go tool cover -func=$(COVERAGE_FILE) | awk -v min=$(COVERAGE_MIN) \
		'/^total:/ { printf "Coverage: %s\n", $$3; if ($$3 + 0 < min) { printf "FAIL: coverage %s is below %d%%\n", $$3, min; exit 1 } }'
