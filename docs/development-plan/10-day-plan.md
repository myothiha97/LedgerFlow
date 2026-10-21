# LedgerFlow 10-Day Plan (v0.1.0)

> Status: Active from 2026-10-21. This plan replaces the B1 to F7 gate order for the first
> release. The older plan files stay as reference for rules and test ideas.
>
> Product source of truth: [Business requirements v3](../specifications/business-requirements-v3.md).
> When this plan and the BRD disagree, the BRD wins and this plan is fixed.

## Goal

LedgerFlow has two goals of equal weight:

1. A personal finance app that is truly useful in daily life, even for someone whose money
   habits are chaotic. That means fast entry, balances you can trust, and a dashboard you
   understand in 10 seconds.
2. A real project for learning backend engineering, DevOps, and system design.

For v0.1.0, ship BRD Phase 1 as a working local app in 10 sessions. Every requirement in BRD Section 10
(AUTH-01 to DATA-01) and every completion criterion in BRD Section 11 must pass. Nothing
outside BRD Section 4 "Included" is added as a product feature.

## How we work

- Backend is the focus. Days 1 to 9 are backend. Day 10 builds the frontend, written fully
  by Claude, with every screen the BRD requires.
- Claude writes about 90% of the code. The owner writes about 10% as a daily exercise.
- One day means one session, not one calendar day. Skipping a calendar day is fine.
- If a day overflows, the rest rolls over to the next session. Scope does not grow.
- The repository stays green at the end of every day: build, vet, and tests pass.
- Exercises are red-to-green. Claude writes the failing test first and marks it
  `t.Skip("TODO(you): ...")`, next to a stub marked `TODO(you)` that compiles and returns a
  clear "not implemented" error. You remove the skip, watch it fail, and make it pass.
  Nothing is broken while the exercise waits.
- At the end of each day Claude writes a short lesson to `docs/learning/days/day-NN.md`:
  what was built, how it works, patterns used (only if used), backend skills practiced,
  where to learn more, the exercise, and 2 or 3 quick recall questions about earlier days.
  Recalling from memory on later days is what makes the learning stick.
- Hard bugs use the `diagnosing-bugs` skill in `.claude/skills/`: first build a command
  that fails on the bug, then shrink it, list 3 to 5 testable guesses, and add a
  regression test before the fix.
- No new planning docs until v0.1.0 ships.
- For the wider backend and DevOps learning path, see
  `docs/learning/backend-roadmap-2026.md`.

## Daily build and verify loop

Every day runs these steps in order. A day is not done until step 5 has evidence.

1. **Ground.** Read the BRD rules and plan items for the day. When a library API, standard,
   or best practice is uncertain, check official docs or the web first. Do not guess.
2. **Build.** Claude implements the day's work, inner layers first: migration, domain,
   store, service, handler, wiring. Independent parts may be built by parallel sub agents.
3. **Review with fresh eyes.** Up to three sub agents, started without the builder's
   context, review the diff in parallel. Each one verifies its claims against official
   docs or the web when unsure, and reports in under 400 words:
   - spec: does the diff do what the BRD and this plan ask? Report missing or partial
     rules, wrong behavior, and work nobody asked for. Quote the BRD line for each finding.
   - standards: `CLAUDE.md`, the engineering standards below, OWASP Top 10:2025, and Go
     idioms. Also flag code smells from Fowler's _Refactoring_ (mysterious names,
     duplicated code, feature envy, primitive obsession, speculative generality, middle
     man) as judgement calls, not hard failures. Skip anything the linters already check.
   - tests: missing cases, weak assertions, and tests that pass for the wrong reason.

   Keep the reports separate. A diff can follow every standard and still build the wrong
   thing, so one report must not hide the other.
4. **Fix.** Every confirmed finding is fixed or recorded with a reason. Unconfirmed findings
   are checked again, not ignored.
5. **Verify gate.** All of these must pass, and the output is recorded in the day lesson:
   - `gofmt -l .` is empty, `go vet ./...` and `golangci-lint run` are clean;
   - `go test -race ./...` passes, including PostgreSQL integration tests;
   - `make generate` twice leaves no diff;
   - new migrations apply, roll back, and apply again;
   - `govulncheck ./...` reports no reachable vulnerability;
   - the server starts and a `curl` smoke test of the day's endpoints returns the
     expected status codes and bodies.

### Model choice for sub agents

The owner is on a limited plan, so each sub agent uses the cheapest model that does the
job well. Never default every agent to Opus.

| Model | Use it for |
| --- | --- |
| Haiku | Searching the repo, checking that paths and links exist, running the verify gate and summarizing its output, simple mechanical edits |
| Sonnet | Most building (CRUD store, queries, handlers, wiring, Docker, CI, frontend pages), test writing, the security and test reviews, web research |
| Opus | The hard parts: the transaction and balance lifecycle, locking, money math, auth security, and the correctness review on Days 6 to 8 |
| Fable | Never used. It needs separate API pricing. |

To keep token use low:

- Give each agent a short, exact prompt: the files, the rules to check, and the answer
  format. Reviewers read the diff and the files it touches, not the whole repo.
- Agents return short findings, not file dumps.
- On a small or low-risk day, use one Sonnet review instead of three.

## Decisions locked for v0.1.0

These settle every open item in BRD Section 13.

| ID | Decision | BRD |
| --- | --- | --- |
| D1 | Account types: `cash`, `bank`, `ewallet`, `savings`. No credit cards. | §5, §13.1 |
| D2 | One currency per user. The first account sets it; every later account must use the same uppercase ISO 4217 code. So the dashboard always shows one total. Reason: budgets and category totals add expenses from all accounts, and mixed currencies would make them wrong. | §6.5, §13.2 |
| D3 | New users get starter categories on register. Income: Salary, Other Income. Expense: Food, Transport, Rent, Utilities, Shopping, Health, Entertainment, Other. | §8, §13.3 |
| D4 | An account can be archived only when its current balance is zero. | §13.4 |
| D5 | The UI and API say "budget remaining", never "safe to spend". | §13.5 |
| D6 | Future-dated transactions are rejected. "Today" uses the app timezone (D9). | §13.6 |
| D7 | No OpenAPI codegen for v0.1.0. The frontend uses small hand-written TypeScript types. | none |
| D8 | Server-side sessions: random token in an `HttpOnly` cookie, only its SHA-256 hash in the database. Passwords use bcrypt. A session lasts 30 days. | §7, AUTH-01 |
| D9 | One app timezone from the `APP_TIMEZONE` env var decides "today" and month boundaries. | §6.5 |
| D11 | Money is `NUMERIC(19,4)` in SQL, `decimal.Decimal` in Go, and a string in JSON. Budget status is decided from the exact percentage. Rounding to 1 decimal place is for display only, so 70.04% is still "caution". | §6.4 |

The frontend calls the API through the Vite dev proxy, so the browser sees one origin.
That means no CORS setup is needed locally.

## Engineering standards (every day)

Speed never skips these. A day is not done until its code follows them. These are
non-functional quality rules. They add no product features beyond the BRD.

Architecture

- Layers: `handler -> service -> store`, plus a framework-free `domain`. Dependencies point
  inward. `main` only builds and wires the parts (dependency injection by constructor).
- Handlers parse input, call one service method, and map the result to HTTP. Business rules
  live only in services, so tests and the Phase 2 AI caller reuse the same rules (BRD §12).
- Small interfaces are defined where they are used (in the service), not next to the
  implementation. No generic repository layer.
- Prefer deep modules: a small interface that hides a lot of behavior. Example:
  `UpdateTransaction` hides the database transaction, the locks, and reverse-then-apply,
  so a caller never has to know them. Avoid shallow pass-through layers.

Data and correctness

- Every user-owned query filters by `user_id`. Another user's record returns 404, not 403,
  so IDs do not leak (AUTH-02).
- Multi-step writes run in one database transaction and roll back on any error (DATA-01).
  Shared rows are locked with `SELECT ... FOR UPDATE`, in stable ID order when there are two.
- PostgreSQL constraints (`NOT NULL`, `CHECK`, `UNIQUE`, foreign keys) are the last guard.
  Services still validate first, so users get clear errors.
- Migrations are versioned, have a working down file, and are never edited after being
  applied outside throwaway local data.
- Money never touches `float64`. Lists have a limit so one request cannot load everything.

API

- REST resources under `/api`, correct status codes (`201`, `204`, `400`, `401`, `404`,
  `409`, `422`, `429`, `500`), and one error format: RFC 9457 Problem Details
  (`application/problem+json` with `type`, `title`, `status`, `detail`, plus a stable
  `code` field the frontend can switch on).
- Internal errors are logged with detail but returned as a generic message. No stack
  traces, SQL, or hashes in responses.
- Request body size is limited. Unknown JSON fields are rejected.
- Every mutation response returns the changed records, so the client can refresh every
  affected view (BRD §6.3).

Reliability and operations

- Config comes from env vars, is validated at startup, and the app fails fast on bad config.
- Structured logs with `log/slog`, one line per request with method, path, status, duration,
  and request ID. Never log passwords, tokens, cookies, notes, or the database URL.
- Every I/O call takes `context.Context` so timeouts and cancelled requests stop work.
- HTTP server timeouts, graceful shutdown on `SIGINT`/`SIGTERM`, liveness (`/healthz`) and
  readiness (`/readyz`) endpoints.

Security

- bcrypt for passwords, crypto-random session tokens, only token hashes stored (BRD §7).
- Cookie: `HttpOnly`, `SameSite=Lax`, `Secure` outside local dev, expiry equal to the session.
- Login failures return the same error for unknown email and wrong password.
- `govulncheck` runs in CI.

Testing and code quality

- Test behavior through public interfaces (service methods, HTTP endpoints), not private
  helpers. A test must survive a refactor that does not change behavior.
- Fake only at system boundaries: the database (for pure rule tests and forced failures),
  time, and randomness. Anything touching data behavior runs against real PostgreSQL. No
  assertions on call counts or call order.
- Expected values come from an independent source: the BRD examples or known literals.
  Never recompute the expected value the way the code does.
- Build in vertical slices: one test, then the code that makes it pass, then the next.
- Table-driven tests named `TestThing_Scenario`. Every BRD rule gets at least one test that
  names the rule it proves.
- `gofmt`, `go vet`, and `go test -race` pass. Errors are wrapped with `%w` and context.
- One focused Conventional Commit per finished piece of work.

## Day by day

### Day 1: Backend foundation

- Claude: finish task 3 (`users.sql` fixes and a clean `sqlc` output), typed config from env,
  PostgreSQL pool, shared Problem Details errors, router, request logging with `slog` and a
  request ID, `GET /healthz`, server timeouts, and graceful shutdown. `main.go` becomes
  wiring only. Set up the verify gate tools: a `.golangci.yml` config, `govulncheck`, and
  a `make check` target that runs the whole gate.
- You: write the `GET /readyz` handler body. Ping the database with a short timeout and
  return 200 or 503.
- Done means: the server starts, `/healthz` returns 200, `/readyz` returns 200 with the
  database up and 503 with it down.

### Day 2: Auth rules without HTTP (AUTH-01)

- Claude: auth domain types and errors, store adapter over `sqlc`, auth service (register,
  login, logout, validate session), bcrypt hashing, random token plus SHA-256 hash, 30-day
  expiry, unit tests with a fake store.
- You: write `normalizeEmail` (trim and lowercase) and add its table-driven test cases.
- Done means: `make test` passes and covers duplicate email, wrong password, and expired
  session.

### Day 3: Auth over HTTP (AUTH-01, AUTH-02)

- Claude: register, login, and me handlers, session cookie, auth middleware that puts the
  user ID in the request context, protected route group.
- You: write the logout handler (call the service, clear the cookie), then run the full
  register, login, me, logout flow with `curl`.
- Done means: `me` returns the user with a valid cookie and 401 without one or after
  logout.

### Day 4: Accounts (ACC-01, ACC-02, ACC-03)

- Claude: migration `0002_accounts`, `sqlc` queries, decimal mapping, store, service,
  handlers for list, get, create, edit, and archive.
- Rules: types from D1, currency from D2, current balance starts equal to the starting
  balance, editing the starting balance recalculates the current balance, name is always
  editable, archived accounts leave the active list but stay readable.
- You: write the D4 archive rule in the service (block when balance is not zero) and its
  test.
- Done means: a second user gets 404 for the first user's account, and money keeps 4
  decimal places through SQL, Go, and JSON.

### Day 5: Categories (CAT-01, CAT-02)

- Claude: migration `0003_categories`, store, service, handlers for list, create, edit, and
  archive, optional color and icon, normalized unique name per user and type, starter
  categories (D3) created inside the register transaction.
- You: write the `sqlc` SQL query that lists active categories by user and type.
- Done means: a new user sees the starter categories, cannot see another user's, and an
  archived category leaves the active list but stays readable.

### Day 6: Transactions part 1, create and delete (TXN-01, DATA-01)

- Claude: migration `0004_transactions`, a store that runs several writes in one database
  transaction, row lock on the account (`SELECT ... FOR UPDATE`), create and delete in the
  service, and a PostgreSQL integration test setup.
- Rules from BRD §6.1: amount above zero, active account and category owned by the user,
  category type matches transaction type, date required and not in the future (D6),
  optional note, `source` always set to `manual` by the server.
- Now that transactions exist, also lock these (BRD §10): account type and currency cannot
  change once the account has transactions, and category type cannot change once the
  category has transactions.
- You: write the domain function that turns a transaction type and amount into a balance
  effect (income adds, expense subtracts), plus its unit test.
- Done means: create and delete change the account balance, and a forced error rolls back
  everything.

### Day 7: Transactions part 2, edit and list (TXN-02, TXN-03)

- Claude: edit with reverse-then-apply for amount, type, account, category, and date;
  moving between accounts with stable lock order; list with date range, account, category,
  and type filters plus a limit; a reconciliation test that recomputes every balance from
  history.
- The BRD §6.2 example is a test: start 1,000, expense 200 gives 800, edit to 150 gives
  850, delete gives 1,000.
- You: add the "move transaction to another account" case to the test matrix.
- Done means: after create, edit amount, change type, move, and delete, every stored
  balance equals the recomputed balance.

### Day 8: Budgets and dashboard (BUD-01, BUD-02, DASH-01, DASH-02)

- Claude: migration `0005_budgets`, budget create, edit, delete, and monthly list, the
  dashboard service and endpoint.
- Budget rules from BRD §6.4: one active expense category per budget, one budget per
  category per month, amount above zero, no income categories, deleting a budget never
  touches transactions, `remaining` may be negative, overall budget remaining counts only
  budgeted categories. Category type is also locked once the category has budgets.
- Dashboard rules from BRD §6.5: defaults to the current month, any month can be selected,
  active account balances and the total, income, expenses, net savings, average daily
  spending (days elapsed including today for the current month, all days for a past month),
  spending by expense category, every budget with status, and recent transactions for that
  month. Past totals include archived accounts and categories. Archived accounts are left
  out of the account summary. Totals are computed in SQL, not in Go loops.
- You: write the status function (`safe`, `caution`, `high_risk`, `exceeded`) and its
  boundary test cases (0, 70, 70.01, 90, 90.01, 100, 100.01).
- Done means: editing, moving, or deleting a transaction moves budget status both up and
  down, and dashboard figures match direct SQL totals.

### Day 9: Production hardening and infrastructure (backend)

- Claude: login rate limiting (token bucket, `429` with `Retry-After`), `Idempotency-Key`
  support on `POST /api/transactions` so a retried request never creates a second record
  (migration `0006_idempotency_keys`), OpenTelemetry traces for HTTP and SQL with Jaeger in
  Compose, multi-stage Dockerfile that runs as a non-root user, Compose service for the API,
  and a GitHub Actions workflow (format, vet, `go test -race` with a PostgreSQL service,
  `govulncheck`, image build).
- You: write the test that sends the same `Idempotency-Key` twice and proves only one
  transaction and one balance change exist.
- Done means: a retried create is safe, the 6th fast wrong login gets `429`, one request
  shows as a trace in Jaeger, and CI is green.

### Day 10: Frontend and release

- Claude: React, Vite, TypeScript, Mantine, TanStack Query, and React Router. It builds the
  six BRD §9 screens with every required action:
  - Register and sign in: register, sign in, sign out, clear errors.
  - Dashboard: month selector and every BRD §6.5 figure.
  - Accounts: list, create, edit, archive, current balance.
  - Categories: list, create, edit, archive, grouped by type.
  - Transactions: list, add, edit, delete, filter by date, account, category, and type.
  - Budgets: select month, create, edit, delete, usage and status shown with a label, not
    only a color.

  Archived records are hidden from new-entry choices. Every mutation refetches all affected
  data, so no reload is ever needed (BRD §6.3). For fast daily entry (BRD §2), "add
  transaction" is one tap away on every screen, defaults to today and the last-used
  account, and every page works at phone width. README setup steps are included.
- You: run the acceptance checklist below yourself, then tag `v0.1.0`.
- Done means: every checklist item passes.

## BRD coverage

| BRD ID | Built on | Proven by |
| --- | --- | --- |
| AUTH-01 | Days 2, 3, 10 | Service tests, `curl` flow, browser reload stays signed in |
| AUTH-02 | Days 3 to 8 | Two-user test on every resource returns 404 |
| ACC-01 | Day 4 | Create test with all fields and D1, D2 rules |
| ACC-02 | Days 4, 6 | Starting balance edit test, type and currency lock test |
| ACC-03 | Day 4 | Archive test (D4) and history still readable |
| CAT-01 | Days 5, 6, 8 | Create and edit tests, type lock tests |
| CAT-02 | Day 5 | Archive test, past transactions unchanged |
| TXN-01 | Day 6 | BRD §6.1 validation table test |
| TXN-02 | Day 7 | Filter tests for each filter and combined |
| TXN-03 | Day 7 | BRD §6.2 example test, reconciliation test |
| BUD-01 | Day 8 | Create, edit, delete, duplicate month tests |
| BUD-02 | Day 8 | Status boundary tests, negative remaining test |
| DASH-01 | Days 8, 10 | Dashboard test against direct SQL totals |
| DASH-02 | Days 8, 10 | Status moves both ways after edits; UI refetch |
| DATA-01 | Days 6, 7, 9 | Forced-failure rollback tests, idempotent retry test |

## Acceptance checklist (v0.1.0)

BRD Section 11 completion criteria:

- [ ] A registered user cannot access another user's data.
- [ ] Create, edit, move, change type, and delete always produce the correct balances.
- [ ] Budget usage and status move both up and down when related transactions change.
- [ ] Archived accounts and categories leave new-entry choices but remain in history.
- [ ] Dashboard figures match the transactions for the selected month.
- [ ] You can open the dashboard and understand balances, monthly spending, and budget
      status within 10 seconds.

User flows from BRD Section 8:

- [ ] First-time setup: register, create an account, review starter categories.
- [ ] Record daily activity without reloading the page.
- [ ] Correct a mistake through the list and filters.
- [ ] Retire an account and a category.

Engineering checks:

- [ ] Sending the same create request twice with one `Idempotency-Key` creates one record.
- [ ] Repeated wrong logins are rate limited with `429`.
- [ ] One request is visible as a trace in Jaeger.
- [ ] The API and database start with `docker compose up`.
- [ ] CI is green on the release commit.

## Moved out of v0.1.0 (Parking Lot)

None of these are BRD Phase 1 requirements. OpenAPI and generated clients, browser
end-to-end test suite, backup and restore drill, CSRF tokens (needed before a public
deploy), public hosting, upgrade to Go 1.27 and PostgreSQL 18 or 19, and all Phase 2 AI
work (BRD §12). Pick these up after v0.1.0. The order is in
`docs/learning/backend-roadmap-2026.md`.
