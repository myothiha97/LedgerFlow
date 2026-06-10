# LedgerFlow — Backend Project Setup (Phase 1 Foundation)

> **Status:** planned, not yet implemented — this is the resume point for backend work.
> Companion to `LedgerFlow-TechSpec-v1.md` and `LedgerFlow-Architecture-Guidelines-v1.md`.
> Implementation begins next session.

---

## Gist (TL;DR)

**What this builds:** the LedgerFlow backend foundation as a **draft scaffold for learning Go** —
not a finished implementation. The boring 50% + two worked references are generated for you; you
implement the conceptually rich 50%.

**Generated for you** (boring boilerplate + worked references):

- All infra: `docker-compose.yml`, `Dockerfile`, `Makefile`, `.env.example`, `.gitignore`,
  `sqlc.yaml`, `go.mod`
- Plumbing: `config/` (env loading), `httpx/` (error envelope), `domain/` (User + errors),
  `store.Store` interface + pgxpool, thin `main.go`, `router.go`, migration `0001_init`
- **Worked reference #1** — `/health` (proves wiring → DB ping)
- **Worked reference #2** — `POST /api/auth/register` *end-to-end through every layer*
  (DTO → validation → bcrypt → sqlc store → 201): a real template for the CRUD pattern you mirror

**You implement** (the learning 50%, marked `// TODO(you):`):

- `sessions.sql` + the session store methods
- `Login` / `Logout` / `ValidateSession` service logic (the session + cookie mechanics)
- `login` / `logout` / `me` handler bodies + the auth middleware
- the remaining test cases

Everything compiles and the server boots from day one; unbuilt routes return `501` until you fill them.

**Flags before resuming:**

1. **Git** — don't `git commit` yet; it would land in the `portfolio-v3` repo. Initialize a
   dedicated repo for `ledgerflow/` first.
2. **Install** Docker Desktop, `sqlc`, and `golang-migrate` to *run* it (Go 1.26 + Make are present).
3. **Locked choices** — sqlc + golang-migrate, Postgres via docker-compose, session-cookie auth,
   Makefile (not Taskfile, since `go-task` isn't installed).

---

## Context

LedgerFlow is a personal-finance tracker, currently **docs-only**: `backend/` and
`frontend/` are empty; the substance lives in three planning docs (`docs/`):
**BRD v2**, **Tech Spec v1**, **Architecture Guidelines v1**.

**Primary goal of this project (per Tech Spec §1): the owner is learning full-stack
development — specifically Go backend.** The owner is an experienced FE engineer
(React/Next/TypeScript, ~5 yrs) and wants to deepen backend skills with Go. So this
task is **not** a turnkey implementation. The deliverable is a **draft scaffold**:

- **I generate the "boring 50%"** — infra, config, wiring, repetitive boilerplate, and
  **one fully-worked vertical reference** (`/health`) that demonstrates the layering once.
- **The owner implements the "learning 50%"** — the Go business logic in `service/`,
  the `store/` method bodies, and the auth handlers — left as compiling stubs marked
  `// TODO(you):` with doc-comment contracts that cite the relevant doc section.

Frontend is explicitly **out of scope for now** (`frontend/` untouched).

Everything below follows **Tech Spec v1** and **Architecture Guidelines v1**; where they
conflict, the Tech Spec wins (Guidelines §0 precedence rule).

### Decisions locked (from clarifying questions)

| Decision | Choice | Source |
|---|---|---|
| DB access | **sqlc** (hand-written SQL → type-safe Go) + **golang-migrate** | TechSpec §13.2 |
| Local DB | **Postgres via docker-compose** (owner installs Docker Desktop) | TechSpec §9 |
| Vertical slice | **Auth** (register/login/logout/me), **session cookies** | TechSpec §13.3 |
| Task runner | **Makefile** (GNU Make present; `go-task` not installed) | TechSpec §9 |
| Git | **Leave as-is** — do NOT run git commands (see caveat) | — |

---

## ⚠️ Git caveat (must respect)

`ledgerflow/` has **no git repo of its own**. The enclosing repo root is `/Users/mtkh97`
(home dir), whose remote is **`portfolio-v3.git`**. Per the owner's choice I will **not run
any git commands**. **Do not `git commit` backend code yet** — it would land in the
portfolio-v3 repo. The Tech Spec (§3) calls for a dedicated repo; initialize one before the
first commit. I will create a `.gitignore` (inert until `git init`) so it's ready.

---

## Division of labor

### I build (the boring 50% + two worked references)

> Kept deliberately **minimal** — initial setup to start coding *from*, not a finished
> skeleton. Stubs are thin; the interesting code is left for you.

Infra / boilerplate:
- All infra/config: `docker-compose.yml`, `Dockerfile`, `.dockerignore`, `Makefile`,
  `.env.example`, `.gitignore`, `sqlc.yaml`.
- `go.mod` + dependency selection.
- `internal/config/` — env loading (Twelve-Factor, Guidelines §1.2).
- `internal/httpx/` — the **consistent error envelope** + response helpers (Guidelines §4.3).
- `internal/domain/` — `User` entity (pure struct) + sentinel errors.
- `db/migrations/0001_init.*.sql` — `users` + `sessions` tables (mundane transcription of TechSpec §7).
- `internal/store/db.go` (pgxpool init) + `store.Store` **interface** (consumer-defined,
  Guidelines §2.4) so services can be tested without a DB.
- `cmd/server/main.go` — wiring only; wires **concrete** store + service structs (which satisfy
  the interfaces) so everything compiles and boots even with stubbed bodies (Guidelines §3.6).
- `internal/handler/router.go` — **all routes registered** (full surface visible).

**Worked reference #1 — `/health`** (proves connectivity + wiring): `internal/handler/health.go`
fully working (pings the pool).

**Worked reference #2 — `POST /api/auth/register` end-to-end through every layer** (the real
CRUD template you'll mirror): the full `handler → service → store → sqlc → domain` flow with a
DTO, shape validation, bcrypt hashing, a sqlc-backed insert, and a `201` response. Concretely:
- `db/queries/users.sql` — `CreateUser`, `GetUserByEmail` (real sqlc queries).
- `internal/store/users_store.go` — implemented against sqlc-generated code.
- `internal/service/auth.go` → `Register(...)` implemented; other methods stubbed (see below).
- `internal/handler/auth.go` → register DTO + validation + handler body implemented; other
  handlers stubbed.
- `internal/service/auth_test.go` — a real table-driven `Register` test (mock store, AAA) as a
  test template; remaining cases `t.Skip("TODO(you)")`.
- `api/openapi.yaml` — documents `/health` + `/api/auth/*` (contract-first, Guidelines §4.6).

### You build (the learning 50% — guided stubs, marked `// TODO(you):`)

With the worked `Register` vertical + `/health` as your templates, you implement the
conceptually rich Go — including the more interesting session/cookie mechanics:

1. `db/queries/sessions.sql` — write the session SQL, then `make generate` (sqlc practice).
2. `internal/store/sessions_store.go` + the session methods — implement against sqlc
   (learn pgx + context + error wrapping). Ships as a **compiling stub** (returns
   not-implemented) so the build stays green until you fill it.
3. `internal/service/auth.go` → `Login`, `Logout`, `ValidateSession` — verify password,
   create/validate/revoke sessions. **This is the service-layer rule the project rests on**
   (BRD §12, Guidelines §2.2). Stubs return not-implemented so it compiles.
4. `internal/handler/auth.go` bodies — `login` / `logout` / `me`: map request → service → HTTP
   via the `httpx` envelope (mirror the `register` handler).
5. `internal/handler/middleware.go` — session-cookie lookup/attach logic (I provide the skeleton).
6. Fill in the remaining `auth_test.go` table cases.

> Boundary is a proposal — adjust if you'd rather own more (e.g. write `users.sql`
> + the `Register` logic yourself too) or have me do more.

---

## Build sequence (order I'll generate in)

1. **Tooling/infra files**: `.gitignore`, `.dockerignore`, `.env.example`, `sqlc.yaml`,
   `docker-compose.yml` (postgres + backend), `backend/Dockerfile` (multi-stage, non-root,
   pinned base — Guidelines §8), `Makefile` (`dev`, `test`, `generate`, `migrate-up/down`,
   `build`).
2. **Go module**: `go mod init` + add deps (gin, pgx/pgxpool, bcrypt, uuid, godotenv).
3. **config**: `internal/config/config.go` (env → typed struct; `Secure` cookie flag is
   env-driven so login works over local http — off in dev, on in prod).
4. **Migrations**: `0001_init` — `users` (id, name, email, password_hash, created_at) +
   `sessions` (id token, user_id, expires_at, created_at). `timestamptz` not `timestamp`
   (Guidelines §5.6).
5. **sqlc**: `sqlc.yaml` + `db/queries/users.sql` (`CreateUser`, `GetUserByEmail`) →
   `make generate` → `internal/store/gen/`.
6. **domain**: `User` struct + sentinel errors (`ErrNotFound`, `ErrInvalidCredentials`, …).
7. **store**: `db.go` (pgxpool), `store.Store` interface; `users_store.go` **implemented**
   (CreateUser/GetUserByEmail); `sessions_store.go` compiling stub (`// TODO(you)`).
8. **service**: `auth.go` — `Register` **implemented**; `Login`/`Logout`/`ValidateSession`
   compile-safe stubs (`// TODO(you)`); `auth_test.go` with a real `Register` test + skipped cases.
9. **httpx**: error envelope + JSON helpers.
10. **handler**: `health.go` (working) + `auth.go` `register` handler (working: DTO + validation);
    `login`/`logout`/`me` stub bodies; `middleware.go` skeleton; `router.go` all routes wired.
11. **main.go**: wiring.
12. **api/openapi.yaml**: document implemented + planned endpoints.

### Target structure

```text
ledgerflow/
├── .env.example  .gitignore  docker-compose.yml  Makefile        [ME]
└── backend/
    ├── go.mod  go.sum  sqlc.yaml  Dockerfile  .dockerignore       [ME]
    ├── api/openapi.yaml                                           [ME]
    ├── db/
    │   ├── migrations/0001_init.{up,down}.sql                     [ME]
    │   └── queries/users.sql [ME]   sessions.sql [YOU]
    ├── cmd/server/main.go                                         [ME]
    └── internal/
        ├── config/config.go                                      [ME]
        ├── httpx/response.go                                      [ME]
        ├── domain/{user.go,errors.go}                            [ME]
        ├── store/{db.go,store.go,gen/}[ME]  users_store.go[ME]  sessions_store.go[YOU]
        ├── service/auth.go[ME:Register / YOU:Login,Logout,Validate]  auth_test.go[ME:1 / YOU:rest]
        └── handler/{router.go,health.go}[ME]  middleware.go[ME skeleton]  auth.go[ME:register / YOU:rest]
```

---

## Tooling to install (owner)

- **Docker Desktop** — required to run Postgres + backend via compose (not installed).
- **sqlc** — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` (wired into `make generate`).
- **golang-migrate** — `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
  (wired into `make migrate-*`); or run via the postgres container.
- Go 1.26.3 ✓ and GNU Make ✓ already present.

Not built now (per docs): CI / GitHub Actions (TechSpec §10 — mid-Phase 1), any frontend,
OpenAPI→code generation (TechSpec §13.1 — deferred; yaml is documentation for now).

---

## Verification

**What I can self-verify** (after `go install`-ing sqlc):
- `cd backend && go build ./... && go vet ./...` pass (all stubs compile).
- `make generate` runs sqlc cleanly against `users.sql`.
- `go test ./...` runs the real `Register` service test green (other cases skipped).

**User-gated — needs Docker installed; I will NOT claim these pass, you run them:**
- `docker compose up` starts Postgres + backend; `make migrate-up` applies `0001_init`.
- `curl localhost:8080/health` → `200` with a DB-ping `ok`.
- `POST /api/auth/register` → `201` (the worked reference, end-to-end).
- Unbuilt routes (`login`/`logout`/`me`) return `501` via the error envelope until you build them.

**Your work (the learning 50%) — verify as you implement each layer:**
- Re-run `go test ./...`, then curl:
  - `POST /api/auth/login` → `200` + `Set-Cookie: session=…` (HttpOnly)
  - `GET /api/auth/me` with the cookie → `200`; without it → `401`
  - `POST /api/auth/logout` → `204`, cookie cleared

---

## Next increments (not now)

Accounts → Categories → **Transactions (the reverse-then-apply balance lifecycle, BRD §6.2
— the highest-value Go logic for you to own)** → Budgets → Dashboard. The `shopspring/decimal`
money-type convention (TechSpec §4.3) gets established with the accounts/transactions slice.
