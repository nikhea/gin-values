.PHONY: fmt fmt-check vet lint test build swagger run migrate help

GOLANGCI ?= golangci-lint

help: ## Show this help
	@grep -E '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS=":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

fmt: ## Format all hand-written Go files
	gofmt -l -w $$(git ls-files '*.go' | grep -v '^docs/')

fmt-check: ## Fail if any hand-written Go file is unformatted
	@test -z "$$(gofmt -l $$(git ls-files '*.go' | grep -v '^docs/'))" || (echo "Unformatted files:"; gofmt -l $$(git ls-files '*.go' | grep -v '^docs/'); exit 1)

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint v2 (installs it if missing)
	@command -v $(GOLANGCI) >/dev/null || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
	$(GOLANGCI) run ./...

test: ## Run the full test suite (needs Postgres, see TEST_DATABASE_URL in tests/setup_test.go)
	go test -count=1 ./...

build: ## Build all binaries
	go build ./...

swagger: ## Regenerate Swagger docs
	swag init

migrate: ## Run DB migrations manually: make migrate ARGS="up|down|version"
	go run ./cmd/migrate $(ARGS)

check: fmt-check vet lint test ## Run the full gate: format, vet, lint, test
