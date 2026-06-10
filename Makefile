# LedgerFlow task runner. The Go module lives in ./backend.
# Requires sqlc and golang-migrate on PATH (installed via `go install`; ensure
# $(go env GOPATH)/bin is on your PATH).

BACKEND := backend
MIGRATIONS := $(BACKEND)/db/migrations
# Read from your shell env / .env; falls back to the local docker-compose DSN.
DATABASE_URL ?= postgres://ledgerflow:ledgerflow@localhost:5432/ledgerflow?sslmode=disable

.PHONY: dev watch build test generate tidy db-up db-down migrate-up migrate-down

dev: ## Run the API on the host (Postgres must be up: make db-up)
	cd $(BACKEND) && go run ./cmd/server

watch: ## Run the API with live reload (requires air: go install github.com/air-verse/air@latest)
	cd $(BACKEND) && air

build: ## Compile the server binary to backend/bin/server
	cd $(BACKEND) && go build -o bin/server ./cmd/server

test: ## Run all Go tests
	cd $(BACKEND) && go test ./...

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
