# LedgerFlow — Backend Walkthrough & Go Primer (v1)

> **Who this is for:** a senior frontend engineer (React/Next/TypeScript) learning Go
> backend. It explains the layered architecture, the Go idioms you need to *read* any file,
> and traces one real request end-to-end. Every Go concept is mapped to a TypeScript
> equivalent you already know.
>
> **How to use it while rebuilding:** read Parts 1–3 to build the mental model, then use
> Part 4 as your build order. Write the substance yourself; the committed scaffold (commit
> `d94d09f`) is there as a reference to peek at when stuck, not to copy.

---

## Part 1 — The mental model

The backend is four layers, and **dependencies only point inward**:

```text
HTTP request
   │
   ▼
handler/   ── parse the request, call a service, shape the response   (knows HTTP/Gin)
   │
   ▼
service/   ── the business logic. THE important layer.                (knows nothing about HTTP)
   │
   ▼
store/     ── talk to Postgres                                        (knows nothing about HTTP or rules)
   │
   ▼
domain/    ── plain data types + error definitions                   (pure — knows nothing about anything)
```

**The rule:** `domain` and `service` never import `handler` or framework code.

**Why it matters for LedgerFlow specifically:** in Phase 2 the AI assistant becomes a
*second caller* of the same `service` functions. If business logic lived in handlers, the
AI couldn't reuse it — you'd rewrite everything. Logic in `service/` makes the AI "just
another caller," not a rewrite.

> **FE analogy:** the same instinct as keeping business logic out of React components and in
> hooks/services — a component (`handler`) orchestrates and renders; it doesn't hold the
> rules. Here the Go *compiler* enforces it via package import rules, not just convention.

---

## Part 2 — Go reading-vocabulary (mapped to TypeScript)

Seven concepts let you read every file in the scaffold.

### 1. `package` + folder = unit of code
Every `.go` file starts with `package domain` (etc.). **One folder = one package.** Imports
use the full module path, e.g. `github.com/myothiha97/ledgerflow/backend/internal/domain`.
Like a TS module, but the *folder* is the unit, not the file.

### 2. `struct` = a TS `type` for data
```go
type User struct {
    ID           uuid.UUID
    Name         string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
}
```
Just `type User = { id: UUID; name: string; ... }` — a bag of typed fields.

### 3. Method with a *receiver* = a method on a class, written outside it
```go
func (s *Store) CreateUser(ctx context.Context, /* ... */) (domain.User, error) {
```
`(s *Store)` is the **receiver**: "this is a method on `Store`, and inside, `s` is `this`."
Go writes `this` explicitly and puts it before the name. So `s.queries` == `this.queries`.

### 4. `*` = pointer = "reference, not a copy"
`*Store` = "a pointer to a Store." Go copies values by default (like JS primitives); `*`
opts into reference semantics (like JS objects). You'll see `*Store`, `*AuthService`
everywhere because we want one shared instance, not copies.

### 5. Interfaces are *structural and implicit* — exactly like TS
```go
type AuthStore interface {
    CreateUser(ctx context.Context, name, email, passwordHash string) (domain.User, error)
    GetUserByEmail(ctx context.Context, email string) (domain.User, error)
    // ...
}
```
There is **no `implements` keyword**. `*store.Store` satisfies `AuthStore` *automatically*
because it has those methods — same as TS structural typing ("if it has the right shape, it
fits"). This is exactly why the service can be unit-tested with a fake: the test's
`mockStore` also has the right shape, so it fits the same interface with no database.

### 6. Errors are *return values*, not exceptions
Go has no `try/catch` for normal flow. Functions return `(result, error)`; you check it:
```go
user, err := h.auth.Register(ctx, req.Name, req.Email, req.Password)
if err != nil {
    // handle it
}
```
Like returning `{ data, error }` everywhere instead of throwing. `err != nil` = "something
went wrong." The repetitive `if err != nil` is the Go tax — embrace it.

### 7. Capitalization = visibility
`CreateUser` (capital) is **exported** (public, importable by other packages).
`toDomainUser` (lowercase) is **private** to its package. No `public`/`private` keywords —
the *case of the first letter* is the access modifier.

### Bonus: `context.Context` = an `AbortSignal` you thread through I/O
The first argument to anything doing I/O (DB calls, future HTTP/LLM calls). It carries
cancellation/timeout/deadline down the call chain. Pass it from the HTTP request all the way
to the DB query, so if the client disconnects, the query can be cancelled. That's why every
store/service method takes `ctx` first.

---

## Part 3 — Tracing `POST /api/auth/register` end-to-end

This one request touches six files. Following it is the fastest way to understand how the
layers connect. Each hop shows the **key idea**, not the full code.

### Hop 0 — Startup wiring · `cmd/server/main.go`
Before any request, `main` builds the object graph by hand (this is dependency injection):
```go
pool, err := store.NewPool(ctx, cfg.DatabaseURL) // open DB pool
st := store.New(pool)                            // concrete store over the pool
authService := service.NewAuthService(st)        // service depends on the store (as an interface)
router := handler.NewRouter(pool, authService, cfg.CookieSecure)
```
Each layer is constructed and handed to the next. `main` is *wiring only* — no logic ever.
(It also does graceful shutdown: start the server in a goroutine, block on an OS signal,
then `srv.Shutdown` to drain in-flight requests.)

### Hop 1 — Routing · `internal/handler/router.go`
```go
authGroup.POST("/register", authHandler.Register) // worked reference #2
```
Gin matches the method + path and calls `Register`. Every route in the API surface is
registered here, so the whole API is visible in one file.

### Hop 2 — Handler · `internal/handler/auth.go` → `Register`
The handler does exactly three things — parse, call, respond:
```go
var req registerRequest
if err := c.ShouldBindJSON(&req); err != nil {            // 1. parse + validate SHAPE
    httpx.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
    return
}
user, err := h.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password) // 2. call service
if err != nil {                                           // 3. map error → HTTP status
    if errors.Is(err, domain.ErrEmailTaken) {
        httpx.Error(c, http.StatusConflict, "email_taken", "that email is already registered")
        return
    }
    httpx.Error(c, http.StatusInternalServerError, "internal_error", "could not register user")
    return
}
httpx.JSON(c, http.StatusCreated, toUserResponse(user))   // success → 201
```
Key points:
- **Shape validation only** here, via struct tags: `binding:"required,email,min=8"`. *Business*
  validity belongs in the service.
- `c.Request.Context()` passes the request's `context` (the AbortSignal) down the chain.
- The response type `userResponse` deliberately omits `PasswordHash` — the hash never leaves
  the server.

### Hop 3 — Service · `internal/service/auth.go` → `Register`
The business logic. No HTTP types in sight:
```go
email = strings.ToLower(strings.TrimSpace(email))                   // normalize (a business rule)
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // hashing is a rule → lives here
if err != nil {
    return domain.User{}, fmt.Errorf("hash password: %w", err)
}
user, err := s.store.CreateUser(ctx, name, email, string(hash))     // call the store INTERFACE
if err != nil {
    return domain.User{}, fmt.Errorf("register: %w", err)
}
return user, nil
```
- `s.store` is the **interface** `AuthStore`, not the concrete store — that's the seam that
  makes this testable.
- `fmt.Errorf("register: %w", err)` **wraps** the error with context. The `%w` verb preserves
  the original so a caller can still `errors.Is(err, domain.ErrEmailTaken)` further up. (Wrap,
  don't swallow.)

### Hop 4 — Store · `internal/store/users_store.go` → `CreateUser`
Talks to Postgres and **translates errors**:
```go
row, err := s.queries.CreateUser(ctx, gen.CreateUserParams{Name: name, Email: email, PasswordHash: passwordHash})
if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation { // "23505" = unique violation
        return domain.User{}, domain.ErrEmailTaken               // raw DB error → domain sentinel
    }
    return domain.User{}, fmt.Errorf("create user: %w", err)
}
return toDomainUser(row), nil                                    // gen.User row → domain.User
```
This is the **translation boundary**: raw infrastructure errors (Postgres SQLSTATE codes)
become clean `domain` errors here, so the service/handler never deal with Postgres specifics.

### Hop 5 — Generated SQL · `internal/store/gen/users.sql.go`
You wrote the SQL in `db/queries/users.sql`; `sqlc generate` produced the type-safe Go:
```go
// generated — do not edit
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email, arg.PasswordHash)
    var i User
    err := row.Scan(&i.ID, &i.Name, &i.Email, &i.PasswordHash, &i.CreatedAt)
    return i, err
}
```
**You write SQL; sqlc writes the Go.** Parameterized (`$1, $2, …`) so SQL injection is
structurally impossible.

### Hop 6 — Back out
`store` returns `domain.User` → `service` returns it → `handler` renders `201` JSON.

### The one big takeaway
**Every layer translates, and `domain` is the shared vocabulary that flows through all of
them:**

| Layer | Translates |
|---|---|
| `handler` | HTTP ⇄ Go calls |
| `service` | raw inputs → business operations (hashing, normalization, rules) |
| `store` | Go ⇄ SQL, and infra-errors → `domain` errors |
| `domain` | the common types (`domain.User`, `domain.ErrEmailTaken`) everyone speaks |

---

## Part 4 — Rebuild roadmap (build it yourself, in this order)

Dependency order — build inner layers first so each new file only imports things that
already exist. For each file: its single job + the governing doc section. Write the code
yourself; ask for review.

| # | File(s) | Its one job | Governing doc |
|---|---|---|---|
| 1 | `go.mod` + deps | module + dependencies (gin, pgx, bcrypt, uuid, godotenv, decimal) | — *(mechanical; ask me)* |
| 2 | `db/migrations/0001_init.*.sql` | `users` + `sessions` tables; `timestamptz` not `timestamp` | TechSpec §7, Guidelines §5.6 |
| 3 | `db/queries/users.sql` + `sqlc.yaml` → `make generate` | hand-write SQL; generate `store/gen/` | Guidelines §5.4 |
| 4 | `internal/domain/` | pure `User`/`Session` structs + sentinel errors | Guidelines §2.1 |
| 5 | `internal/config/config.go` | env → typed `Config` (Twelve-Factor) | Guidelines §1.2 |
| 6 | `internal/httpx/response.go` | the one consistent error envelope + JSON helpers | Guidelines §4.3 |
| 7 | `internal/store/` | `NewPool`, concrete `Store`, `users_store.go` (map rows + errors) | Guidelines §5 |
| 8 | `internal/service/auth.go` | `AuthStore` interface + `Register` (bcrypt, normalize, wrap errors) | Guidelines §2.2 |
| 9 | `internal/service/auth_test.go` | table-driven test of `Register` with a mock store (AAA) | Guidelines §9 |
| 10 | `internal/handler/` | `health`, `auth` (register), `router` (all routes), `middleware` skeleton | Guidelines §4 |
| 11 | `cmd/server/main.go` | wiring only: config → pool → store → service → router → serve | Guidelines §3.6 |

**The mechanical 20% I can still generate for you on request** (no Go-learning value in
hand-typing these): `go.mod`, `Dockerfile`, `docker-compose.yml`, `Makefile`, `.env.example`,
`.gitignore`, `sqlc.yaml`. Everything else — the four layers — you write.

**The highest-value Go to write yourself** (don't let anyone hand you these):
- The `service/` functions — business logic is the whole point of the layer.
- The `store/` error translation (mapping `pgx.ErrNoRows` / `23505` to `domain` errors).
- Later: the **transaction balance lifecycle** (reverse-then-apply, BRD §6.2) — the single
  richest piece of Go in the project, and the thing LedgerFlow lives or dies on.

---

## Quick Go syntax cheatsheet

```go
package store                       // every file declares its package

import (                            // grouped imports; full module paths
    "context"
    "github.com/google/uuid"
)

type Config struct { Port string }  // a data type

func (s *Store) Foo(ctx context.Context, id uuid.UUID) (domain.User, error) {
    //   ^receiver        ^ctx first for I/O      ^multiple return values (result, error)
    user, err := s.something(ctx, id)
    if err != nil {                 // always check errors immediately
        return domain.User{}, fmt.Errorf("foo: %w", err) // wrap with %w
    }
    return user, nil                // nil = "no error"
}

var ErrThing = errors.New("thing")  // a sentinel error (compare with errors.Is)
```

- `:=` declares + infers type (`x := 5`). `=` assigns to an existing variable.
- `domain.User{}` is the **zero value** — an empty struct you return alongside an error.
- A lowercase name is package-private; an Uppercase name is exported.

---

*Living document — extend it as you build accounts, categories, transactions, budgets, and
the dashboard. The patterns above repeat for every entity.*
