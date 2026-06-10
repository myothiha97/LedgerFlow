# LedgerFlow

> A personal finance tracker that turns scattered daily spending into a single, honest answer to **"Can I afford this right now?"**

LedgerFlow keeps an accurate running picture of your money and your budgets, so financial
decisions stop being guesswork. It is not an accounting suite — it's a discipline tool. Every
feature exists to answer five practical questions:

| The user asks… | LedgerFlow answers with… |
|---|---|
| How much money do I have? | Total balance across all accounts |
| How much did I spend this month? | Monthly expense summary |
| Where is my money going? | Top spending categories |
| Am I still within budget? | Per-category budget status |
| How much can I safely spend? | Remaining budget + daily average |

**Success = open the app and understand your financial position in under 10 seconds.**

---

## Status

🚧 **Phase 1 (Core MVP) — in active development.**

This is primarily a **full-stack learning project** (deepening Go backend skills). The backend
was reset to be rebuilt from scratch, so `backend/` currently boots a minimal Gin server while
the layered architecture below is implemented incrementally. The frontend has not been started.

| Area | State |
|---|---|
| Planning docs (BRD, Tech Spec, Architecture) | ✅ Complete — see [`docs/`](docs/) |
| Backend foundation (Gin + Postgres) | 🚧 In progress |
| Auth · Accounts · Categories · Transactions · Budgets · Dashboard | ⏳ Planned (Phase 1) |
| Frontend (React SPA) | ⏳ Not started |
| AI assistant | 🔮 Phase 2 |

---

## Tech Stack

| Layer | Choice | Why |
|---|---|---|
| Backend language | **Go** | Single dependency-free binary; compiler enforces the service-layer boundary |
| Backend framework | **Gin** | Fast, minimal; keeps HTTP handlers thin |
| Database | **PostgreSQL** | Relational fit; `NUMERIC`/`DECIMAL` for money |
| DB access | **sqlc** + **golang-migrate** | Hand-written SQL → type-safe Go |
| Money type | **shopspring/decimal** | Money never touches `float64` |
| API style | **REST (JSON)** | Simple to consume and test |
| Frontend | **React + Vite + TypeScript** | SPA, no SSR (auth-walled, no SEO surface) |
| UI / state | **Mantine** + **TanStack Query** | Batteries-included UI; server-state caching |
| Containerization | **Docker + docker-compose** | Local dev + host portability |

> **Architectural rule that everything rests on:** all business logic lives in a **service layer**
> that HTTP handlers call — handlers contain no logic. This is what lets the Phase 2 AI assistant
> become *just another caller*, not a rewrite.

---

## Repository Structure

A monorepo with two independently-tooled sibling projects (Go and JS don't share a package manager):

```text
ledgerflow/
├── backend/              # Go module (go.mod lives here, not at root)
│   ├── cmd/server/       # wiring only: config, router, DB, start
│   └── internal/
│       ├── handler/      # thin Gin handlers: parse → call service → respond
│       ├── service/      # the business logic (balance lifecycle, budget status…)
│       ├── domain/       # entities, money type, lifecycle rules
│       └── store/        # DB access; balance recalc on write
├── frontend/             # React SPA (package.json lives here, not at root)
├── docs/                 # BRD, Tech Spec, Architecture Guidelines, session notes
├── docker-compose.yml    # local Postgres (+ optionally backend)
├── Makefile              # dev · test · generate · migrate · build
└── .env.example
```

---

## Core Concepts

Five entities, one mental model:

- **Accounts** — where money lives (cash, bank, e-wallet, savings, credit card). Each has a balance.
- **Categories** — what money is for (income: Salary, Freelance; expense: Food, Rent, Transport…).
- **Transactions** — the core record. Income or expense, tied to one account and one category.
- **Budgets** — a monthly spending cap on an expense category.
- **Dashboard** — the read layer that aggregates everything into a snapshot.

**The critical rule:** a transaction's effect on an account must be **reversed before any edit or
delete**, or balances drift out of sync. Edits always fully reverse the old version, then apply the
new one cleanly — never patch a balance incrementally. (See [BRD §6.2](docs/LedgerFlow-BRD-v2.md).)

---

## Getting Started

### Prerequisites

- **Go** 1.26+
- **Docker Desktop** (for Postgres)
- **sqlc** and **golang-migrate** on your `PATH` (used by `make generate` / `make migrate-*`):
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

### Setup

```bash
# 1. Configure environment
cp .env.example .env

# 2. Start Postgres
make db-up

# 3. Run the backend
make dev
```

### Common tasks (Makefile)

| Command | Description |
|---|---|
| `make dev` | Run the API on the host (Postgres must be up) |
| `make watch` | Run the API with live reload (requires [air](https://github.com/air-verse/air)) |
| `make test` | Run all Go tests |
| `make generate` | Regenerate type-safe Go from SQL (sqlc) |
| `make migrate-up` / `make migrate-down` | Apply / roll back migrations |
| `make db-up` / `make db-down` | Start / stop Docker services |
| `make build` | Compile the server binary to `backend/bin/server` |

> **Note:** the backend is mid-rebuild. The current `backend/main.go` boots a minimal Gin server
> on `:4000` (`cd backend && go run .`); the `make`-based workflow and `cmd/server` layout above
> describe the target structure being built out per [the Tech Spec](docs/LedgerFlow-TechSpec-v1.md).

---

## API Overview

```http
# Auth
POST /api/auth/register
POST /api/auth/login
POST /api/auth/logout
GET  /api/auth/me

# Accounts
GET|POST       /api/accounts
PATCH|DELETE   /api/accounts/:id

# Categories
GET|POST       /api/categories
PATCH|DELETE   /api/categories/:id

# Transactions
GET|POST       /api/transactions
PATCH|DELETE   /api/transactions/:id

# Budgets
GET  /api/budgets?month=6&year=2026
POST /api/budgets
PATCH|DELETE   /api/budgets/:id

# Dashboard
GET  /api/dashboard/summary?month=6&year=2026
```

---

## Roadmap

```text
Phase 1 — Core MVP (current)      Auth, accounts, categories, transactions,
  local + Docker                  budgets, dashboard. Correct balance lifecycle.

Phase 2 — AI Assistant            Natural-language entry, conversational queries,
  still local + Docker            auto-categorization — built on the same services.

Phase 3 — Real Business App       Subscriptions, roles/permissions, paid hosting + CD.
  hosted, paid
```

Phase 1 ships zero AI but is architected so the assistant slots in without a rewrite. Deferred
features (transfers, recurring transactions, saving goals, CSV import, multi-currency, mobile/PWA)
are explicitly out of Phase 1 scope.

---

## Documentation

| Doc | Purpose |
|---|---|
| [Business Requirements (BRD v2)](docs/LedgerFlow-BRD-v2.md) | What LedgerFlow is and why — product vision, business logic, data model |
| [Technical Specification (v1)](docs/LedgerFlow-TechSpec-v1.md) | How it's built — stack, architecture, delivery phases |
| [Architecture Guidelines (v1)](docs/LedgerFlow-Architecture-Guidelines-v1.md) | Coding conventions and layering rules |
| [Backend Setup (v1)](docs/LedgerFlow-Backend-Setup-v1.md) | Backend foundation plan |
| [Backend Walkthrough (v1)](docs/LedgerFlow-Backend-Walkthrough-v1.md) | Guided tour of the backend |
