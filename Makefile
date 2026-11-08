# LedgerFlow task runner. The Go module lives in ./backend.
# Requires sqlc and golang-migrate on PATH (installed via `go install`; ensure
# $(go env GOPATH)/bin is on your PATH).
BACKEND := backend
MIGRATIONS := $(BACKEND)/db/migrations
# Load .env when it exists. Its values OVERRIDE shell env vars (make rule), so edit .env
# rather than prefixing `PORT=... make dev`. Make parses it: no quotes and no `$` in values.
-include .env
DATABASE_URL ?= postgres://ledgerflow:ledgerflow@localhost:5432/ledgerflow?sslmode=disable
# Pass the settings the server reads (internal/config) to every command.
export PORT APP_ENV DATABASE_URL

.PHONY: dev watch build test lint vuln check generate tidy db-up db-down migrate-up migrate-down

dev: ## Run the API on the host (Postgres must be up: make db-up)
	cd $(BACKEND) && go run ./server

watch: ## Run the API with live reload (requires air: go install github.com/air-verse/air@latest)
	cd $(BACKEND) && air

build: ## Compile the server binary to backend/bin/server
	cd $(BACKEND) && go build -o bin/server ./server

test: ## Run all Go tests
	cd $(BACKEND) && go test ./...

lint: ## Run golangci-lint (formatters + linters, config in backend/.golangci.yml)
	cd $(BACKEND) && golangci-lint fmt --diff && golangci-lint run ./...

vuln: ## Check dependencies for known vulnerabilities (go install golang.org/x/vuln/cmd/govulncheck@latest)
	cd $(BACKEND) && govulncheck ./...

check: lint ## The verify gate: lint, vet, build, race tests, vulnerability scan
	cd $(BACKEND) && go vet ./... && go build ./... && go test -race ./...
	$(MAKE) vuln

generate: ## Regenerate type-safe Go from SQL (sqlc)
	cd $(BACKEND) && sqlc generate

tidy: ## Sync go.mod / go.sum
	cd $(BACKEND) && go mod tidy

db-up: ## Start only Postgres via docker-compose
	docker compose up -d db

db-down: ## Stop all docker-compose services
	docker compose down

migrate-up: ## Apply all up migrations
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" up

migrate-down: ## Roll back the most recent migration
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" down 1
