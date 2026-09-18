.PHONY: fmt fmt-check vet lint test build swagger run migrate vendor help

GOLANGCI ?= golangci-lint

help: ## Show this help
	@grep -E '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS=":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

fmt: ## Format all hand-written Go files
	gofmt -l -w $$(git ls-files '*.go' | grep -v '^docs/' | grep -v '^vendor/')

fmt-check: ## Fail if any hand-written Go file is unformatted
	@test -z "$$(gofmt -l $$(git ls-files '*.go' | grep -v '^docs/' | grep -v '^vendor/'))" || (echo "Unformatted files:"; gofmt -l $$(git ls-files '*.go' | grep -v '^docs/' | grep -v '^vendor/'); exit 1)

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

vendor: ## Re-vendor modules AND re-sanitize known dummy secrets (always use this, never raw go mod vendor)
	go mod vendor
	sed -i 's#https://hooks.slack.com/services/[A-Za-z0-9/]*#https://hooks.slack.com/services/REDACTED#' vendor/github.com/go-openapi/spec/appveyor.yml
	@test -z "$$(grep -rEn 'hooks.slack.com/services/[A-Za-z0-9/]+' vendor/ | grep -v REDACTED)" || (echo "Secret pattern reappeared in vendor/"; exit 1)

check: fmt-check vet lint test ## Run the full gate: format, vet, lint, test
