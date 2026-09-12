# Backend and Infrastructure Plan

> Scope: B1 through I2 and readiness gate G1.
>
> Start here: B1.
>
> Frontend source work is forbidden until G1 passes.

## How to read each milestone

- `Goal` states the usable outcome.
- `Build` lists the required work in dependency order.
- `Verify` lists evidence, not assumptions.
- `Exit gate` is the rule for moving to the next milestone.

Each milestone should be split into small session tasks in `backend/tasks/`. Keep one task
active until it is reviewed and verified.

## B1. Finish the `sqlc` foundation

**Goal:** Hand-written user and session SQL generates type-safe Go that compiles.

**Why now:** Migration `0001_init` exists, but no Go code can use it yet. This is the
current objective recorded in `backend/SESSION_STATE.md` and `backend/tasks/2.md`.

**Build:**

1. Add `backend/db/queries/users.sql` with the user and session operations needed by auth.
2. Add `backend/sqlc.yaml` at the Go module root.
3. Configure PostgreSQL generation for `pgx/v5` into `backend/internal/store/gen`.
4. Generate models and query methods. Never edit generated files by hand.
5. Add only the direct dependencies required by generated code, then tidy the module.
6. Check that query names, parameter types, row counts, and returned columns match migration
   `0001`.

**Required queries:**

- create a user;
- get a user by normalized email;
- get a user by ID;
- create a session with a hashed token;
- get a valid session and its user by token hash;
- delete one session by token hash and user;
- optionally delete expired sessions only when cleanup behavior is implemented.

Every ownership-sensitive delete or lookup must include enough data to prevent changing
another user's session.

**Verify:**

- `make generate` completes and can be repeated without changing output.
- Generated `User` and `Session` fields match migration `0001`.
- Go format, test, vet, and build pass.
- No generated file has a manual edit.

**Exit gate:** A reviewed generated store package exists and the module is green.

## B2. Establish the runtime architecture foundation

**Goal:** Replace the one-file server with the agreed layered runtime without adding feature
business logic.

**Prerequisite:** B1.

**Build:**

1. Add framework-free auth entities and shared domain errors under
   `backend/internal/domain/`.
2. Add typed environment loading under `backend/internal/config/`. Validate required values
   at startup and keep secrets out of logs.
3. Add a PostgreSQL pool and concrete store under `backend/internal/store/`.
4. Add one consistent JSON success and error format under `backend/internal/httpx/` or the
   equivalent established package.
5. Add router construction under `backend/internal/handler/`.
6. Add a liveness endpoint that proves the process is running and a readiness endpoint that
   proves PostgreSQL is reachable.
7. Keep `backend/server/main.go` limited to configuration, dependency wiring, server
   startup, signal handling, and graceful shutdown.
8. Set explicit HTTP read, write, idle, and shutdown timeouts.
9. Replace `/ping` only after the new health endpoints are verified.

**Layer rules:**

```text
handler -> service -> store/domain
main    -> constructs concrete dependencies
domain -> imports no Gin or database package
```

Interfaces live in the package that consumes them. Do not add generic repositories or an
interface until a service needs that boundary for real work or testing.

**Verify:**

- The server starts with valid configuration and fails clearly with invalid configuration.
- Liveness succeeds without using business data.
- Readiness succeeds with PostgreSQL available and fails when it is unavailable.
- An interrupt causes graceful shutdown without a panic.
- Package imports follow the dependency direction.
- Go format, test, vet, and build pass.

**Exit gate:** The application skeleton matches the technical architecture and contains no
business rules in handlers or `main`.

## B3. Create the API contract foundation

**Goal:** Make `backend/api/openapi.yaml` the reviewed contract for all Phase 1 HTTP work.

**Prerequisite:** B2 and decision D7.

**Build:**

1. Add an OpenAPI document with server information, reusable schemas, authentication, and the
   shared error envelope.
2. Define auth endpoints first, then extend the contract before each later feature.
3. Represent money as decimal strings in JSON. Do not use JSON floating-point numbers.
4. Represent financial dates as `YYYY-MM-DD` and instants as RFC 3339 strings.
5. Define UUIDs, validation limits, status codes, and error codes consistently.
6. Define transaction filter names and bounded list behavior before implementing the list.
7. Generate backend contract types or interfaces with the chosen OpenAPI tool.
8. Add one project command that regenerates all contract output.
9. Document that generated files are read-only.

Do not place domain behavior in generated API types. Handlers translate between contract
types and the service vocabulary.

**Verify:**

- The OpenAPI document validates.
- Generation succeeds twice with no unexplained diff.
- Auth routes, request shapes, responses, cookies, and errors are described.
- A contract test proves a decimal amount remains exact through JSON serialization.
- Go format, test, vet, build, and generation checks pass.

**Exit gate:** Auth HTTP implementation can be written against a stable, validated contract.

## B4. Build auth domain, store, and service behavior

**Goal:** Registration and session behavior work without depending on Gin.

**Prerequisites:** B1 through B3 and locked decision D8.

**Build:**

1. Define the service's small `AuthStore` interface from the methods auth actually uses.
2. Translate generated rows and PostgreSQL errors into domain types and errors in the store.
3. Implement registration with trimmed input, normalized lowercase email, password policy,
   and a vetted password hash.
4. Implement login with constant-behavior credential failure, a cryptographically random
   session token, stored token hash, and expiry.
5. Implement logout by deleting the stored session.
6. Implement session validation without returning password hashes or stored token hashes.
7. Pass `context.Context` from every service method to every store operation.
8. Wrap unexpected errors with operation context while preserving domain errors.

The raw session token may exist only long enough to return it to the HTTP boundary. Store and
log only its hash.

**Unit-test cases:**

- valid registration;
- trimmed and normalized email;
- invalid name, email, and password;
- duplicate email;
- valid login;
- unknown email and wrong password return the same public error;
- valid, expired, missing, and revoked sessions;
- store and hashing failures are returned safely.

**Verify:**

- Service tests use a small fake or mock store and contain no Gin or PostgreSQL dependency.
- Password hashes differ from plaintext and validate only the correct password.
- Session tokens are not persisted raw.
- Go format, test, vet, and build pass.

**Exit gate:** Auth business behavior is fully tested below the HTTP layer.

## B5. Complete auth HTTP and middleware

**Goal:** AUTH-01 and the auth part of AUTH-02 work through real HTTP and PostgreSQL.

**Prerequisite:** B4.

**Build:**

1. Implement `POST /api/auth/register`.
2. Implement `POST /api/auth/login` and set an `HttpOnly` session cookie.
3. Implement `POST /api/auth/logout` and clear the cookie.
4. Implement `GET /api/auth/me`.
5. Add authentication middleware that resolves the session once and places the user identity
   in request context.
6. Apply middleware to every protected route group.
7. Use `Secure` cookies outside local HTTP, an appropriate `SameSite` policy, a narrow path,
   and expiry that matches the server session.
8. Add narrow CORS configuration only if the SPA uses a different local origin.
9. Add CSRF protection or strict origin checks before any cookie-authenticated public deploy.
10. Keep shape validation and HTTP mapping in handlers. Keep auth rules in the service.

**HTTP and integration checks:**

- register returns a safe user response and never exposes a hash;
- duplicate registration returns the documented conflict;
- login creates a usable session cookie;
- `me` rejects missing, invalid, expired, and revoked sessions;
- logout invalidates the server session and clears the browser cookie;
- protected routes reject anonymous requests;
- malformed JSON and wrong content types return the shared error shape.

**Exit gate:** The complete auth flow works through HTTP and a real database, and its contract
tests pass.

## B6. Build accounts end to end

**Goal:** ACC-01 through ACC-03 work with correct starting and current balances.

**Prerequisites:** B5 and decisions D1, D2, and D4.

**Design first:**

1. Add migration `0002_accounts` using the database design as the target.
2. Record the supported Phase 1 account types. Do not include credit cards unless D1 is
   explicitly changed with a debt balance rule.
3. Record the Phase 1 currency policy and archive policy.
4. Extend OpenAPI before adding handlers.

**Build:**

1. Add account SQL and regenerate `sqlc`.
2. Map PostgreSQL `NUMERIC(19,4)` to `shopspring/decimal`.
3. Add account domain types, errors, store translation, service logic, handlers, and routes.
4. Set `current_balance` equal to `initial_balance` on create.
5. Recalculate current balance when the starting balance changes.
6. Prevent account type and currency changes after transactions exist.
7. Implement archive as `archived_at`, not a physical delete.
8. Exclude archived accounts from new transaction choices and the active account summary.
9. Keep ownership checks in every account operation.

**Verify:**

- Migration apply, inspection, rollback, and reapply pass.
- Amounts preserve four decimal places through SQL, Go, and JSON.
- Two users cannot read, edit, or archive each other's accounts.
- Creating and editing a starting balance produces the expected current balance.
- Archive behavior matches D4 and preserves history.
- Unit, store integration, handler, and contract tests pass.

**Exit gate:** Accounts work end to end and their rules are protected before transactions can
depend on them.

## B7. Build categories end to end

**Goal:** CAT-01 and CAT-02 work with stable income and expense meaning.

**Prerequisite:** B6 and decision D3.

**Design first:**

1. Add migration `0003_categories` with type, optional color and icon, timestamps, archive
   state, ownership, and normalized uniqueness.
2. Decide whether registration or category setup creates the starter set.
3. Extend OpenAPI before implementation.

**Build:**

1. Add category SQL and regenerate `sqlc`.
2. Add domain, store, service, handler, and route behavior.
3. Allow only `income` or `expense` types.
4. Normalize names consistently before the unique constraint is reached.
5. Prevent category type changes after transactions or budgets exist.
6. Archive instead of deleting and hide archived categories from new choices.
7. Preserve category information in past transaction and dashboard queries.

**Verify:**

- Migration apply, inspection, rollback, and reapply pass.
- Duplicate normalized names are handled predictably.
- Two users cannot access each other's categories.
- Archived categories disappear from active lists but remain usable for history reads.
- Type-change rules pass before and after dependent records exist.
- All relevant automated checks pass.

**Exit gate:** Accounts and categories provide valid, owned, active inputs for transactions.

## B8. Build the transaction schema and atomic store

**Goal:** PostgreSQL can persist transaction changes and account balance effects as one atomic
operation.

**Prerequisites:** B6 and B7.

**Design first:**

1. Add migration `0004_transactions` with the checks, foreign keys, `source`, `note`, date,
   timestamps, and indexes from the database design.
2. Extend OpenAPI with transaction schemas, mutation inputs, filters, and errors.
3. Define the exact database transaction boundary before writing store code.

**Build:**

1. Add read queries scoped by user, date range, account, category, and type.
2. Add transaction-capable store operations for create, edit, and delete.
3. Lock affected account rows before changing balances.
4. When two accounts are affected, lock them in stable UUID order.
5. Store transaction amounts as positive decimals. The type controls the sign effect.
6. Keep every transaction record change and every balance change in the same PostgreSQL
   transaction.
7. Roll back on validation, constraint, query, or commit failure.
8. Add a reconciliation query that derives each account balance from its starting balance and
   transactions.

This milestone builds persistence mechanics only. B9 owns the business decision sequence.

**Integration-test cases:**

- create income and expense;
- edit amount;
- edit type;
- move between accounts;
- delete;
- forced failure after the old effect is reversed;
- concurrent writes to the same account;
- reconciliation after every operation.

**Exit gate:** The store proves atomicity and balance integrity against real PostgreSQL.

## B9. Build transaction service and API behavior

**Goal:** TXN-01 through TXN-03 and DATA-01 work through the complete stack.

**Prerequisite:** B8 and decision D6.

**Build:**

1. Implement `CreateTransaction`, `UpdateTransaction`, and `DeleteTransaction` in the service.
2. Validate positive amount, date policy, ownership, active account, active category, and
   category type compatibility.
3. Accept an optional note and set `source` to `manual` on the server. Phase 1 clients cannot
   claim another source.
4. Apply the full new effect for create.
5. Reverse the full old effect, then apply the full new effect for every edit.
6. Reverse the full old effect for delete.
7. Never update only the difference between old and new amounts.
8. Implement listing and filters by date, account, category, and type with bounded results.
9. Add thin handlers and contract-aligned status and error mapping.
10. Ensure successful mutations return enough information for a client to refresh affected
   accounts, months, categories, budgets, and dashboard data.

**Required service test matrix:**

| Change | Old effect | New effect | Accounts affected |
| --- | --- | --- | --- |
| Create income | None | Add | New account |
| Create expense | None | Subtract | New account |
| Edit amount | Reverse old | Apply new | Same account |
| Edit type | Reverse old type | Apply new type | Same account |
| Move account | Reverse old | Apply new | Old and new accounts |
| Change category or date | Reverse old record effect | Apply new record effect | Account plus old and new reporting groups |
| Delete | Reverse old | None | Old account |

Also test archived inputs, mismatched category type, other-user IDs, zero and negative amount,
future dates according to D6, missing records, and injected failures.

**Verify:**

- Unit tests prove the service sequence independently of PostgreSQL.
- Integration tests prove the same sequence commits or rolls back as one unit.
- HTTP tests prove documented status codes and response shapes.
- Reconciliation remains exact after the full matrix.
- Race and concurrency checks pass for shared account updates.

**Exit gate:** Transaction history and every stored account balance remain correct after all
supported lifecycle changes.

## B10. Build budgets end to end

**Goal:** BUD-01 and BUD-02 work for one expense category per calendar month.

**Prerequisite:** B9 and decision D11.

**Design first:**

1. Add migration `0005_budgets` using `month_start` and the required unique constraint.
2. Extend OpenAPI with budget mutations and monthly reads.
3. Keep `spent`, `remaining`, percentage, and status out of stored budget columns. They are
   derived from current expense transactions.
4. Record percentage and display rounding rules before implementing calculations.

**Build:**

1. Add budget and monthly spending queries, then regenerate `sqlc`.
2. Validate ownership, active expense category, first-day month normalization, positive amount,
   and uniqueness in the service.
3. Implement create, edit, delete, and monthly list operations.
4. Calculate status from current spending every time it is read.
5. Calculate overall remaining only from budgeted categories.
6. Deleting a budget must not change transactions.

**Boundary tests:**

- exactly 0%, 70%, above 70%, 90%, above 90%, 100%, and above 100%;
- negative remaining when exceeded;
- movement to a lower status after a transaction edit or delete;
- duplicate category and month;
- income, archived, missing, and other-user categories;
- spending in an unbudgeted category affects monthly expense but not overall budget remaining.

**Exit gate:** Budget values and statuses always reflect current transactions in both directions.

## B11. Build the dashboard end to end

**Goal:** DASH-01 and DASH-02 return a trustworthy selected-month summary.

**Prerequisites:** B10 and decisions D2, D9, and D11.

**Design first:**

1. Define the dashboard response in OpenAPI before SQL or handler implementation.
2. Lock the month boundary and application timezone rule.
3. Lock the mixed-currency presentation rule. Never sum unlike currencies.

**Build:**

1. Add queries for active account balances, monthly income, expenses, category spending,
   budget status inputs, and recent transactions.
2. Implement one dashboard service that assembles those results.
3. Calculate net savings as income minus expenses.
4. For the current month, divide average daily spending by calendar days elapsed, including
   today. For past months, divide by days in that month.
5. Include archived accounts and categories in past transaction totals, but exclude archived
   accounts from the current account summary.
6. Keep dashboard calculations reusable outside HTTP.

**Verify:**

- Empty month, current month, completed month, leap February, and month boundaries pass.
- Income, expenses, net savings, category totals, budget summaries, and recent transactions
  match direct source queries.
- Archived-history and mixed-currency rules pass.
- A transaction create, edit, move, type change, date change, or delete changes every affected
  dashboard value on the next read.
- Two users receive completely separate dashboards.

**Exit gate:** A backend consumer can answer every dashboard question in BRD Section 6.5 from
one stable endpoint.

## B12. Harden and accept the backend

**Goal:** Verify cross-cutting behavior that feature-local tests can miss.

**Prerequisite:** B11.

**Build and verify:**

1. Run a two-user authorization suite against every protected endpoint.
2. Run the complete transaction lifecycle and balance reconciliation suite.
3. Force failures inside multi-step database writes and prove full rollback.
4. Confirm all handlers use the same error envelope and do not leak internal details.
5. Add request body limits, server timeouts, safe logging, and panic recovery.
6. Confirm passwords, raw tokens, cookies, database URLs, and personal notes are not logged.
7. Review cookie, CORS, CSRF, and origin behavior for the planned frontend topology.
8. Check direct and transitive dependencies for known Go vulnerabilities.
9. Run Go tests with the race detector where practical.
10. Verify every generated artifact can be reproduced from source.
11. Review relevant OWASP ASVS Level 1 authentication, session, access-control, validation,
    and API requirements.

**Exit gate:** The backend feature gate is green and no known correctness or security defect
blocks local Phase 1 use.

## I1. Package the backend and local infrastructure

**Goal:** A clean machine can run the API and PostgreSQL from documented, reproducible inputs.

**Prerequisite:** B12.

**Build:**

1. Add a multi-stage backend Dockerfile with pinned base versions.
2. Run the final image as a non-root user and include only runtime files.
3. Add `.dockerignore` entries for Git data, secrets, local builds, caches, and frontend output.
4. Extend Compose with the backend service, health checks, dependency readiness, environment
   values, and the existing PostgreSQL volume.
5. Keep migrations an explicit release or development step. Do not silently edit schema at
   application startup.
6. Keep the host-based Go development path available for fast learning cycles.
7. Emit logs to standard output and standard error.
8. Document setup, migration, reset of disposable data, backup, and restore.
9. Perform one local PostgreSQL backup and restore drill before relying on local data.

Follow Docker's official
[building best practices](https://docs.docker.com/build/building/best-practices/), while keeping
the solution limited to one API container and one database container.

**Verify:**

- Compose configuration validates.
- A clean build produces the backend image without copying `.env` or development tools.
- The process inside the final image is non-root.
- A fresh database accepts all migrations.
- Liveness and readiness behave correctly during startup and database interruption.
- Auth and one complete transaction flow work through the containerized API.
- Backup data can be restored and reconciled.

**Exit gate:** Local infrastructure is reproducible, observable, and recoverable.

## I2. Add continuous integration

**Goal:** Every push receives the same checks from a clean environment.

**Prerequisite:** I1.

**Build:**

1. Add a focused GitHub Actions workflow for backend and infrastructure checks.
2. Pin action versions and use minimum required permissions.
3. Cache dependencies without caching generated project truth.
4. Start PostgreSQL as a CI service for migration and integration checks.
5. Fail when formatting, module files, generated code, or the OpenAPI contract is stale.
6. Build the backend image without publishing it.
7. Keep deployment automation out of Phase 1 unless a real public target is selected.

**Required CI checks:**

- formatting;
- module tidiness;
- OpenAPI validation and generation consistency;
- `sqlc` generation consistency;
- unit tests;
- PostgreSQL integration tests;
- migration apply, rollback, and reapply on a disposable database;
- vet and build;
- Go vulnerability scan;
- Docker image build.

**Exit gate:** The workflow passes on the current commit and fails on a deliberately broken
test or stale generated artifact.

## G1. Backend and infrastructure readiness gate

**Goal:** Formally decide that frontend implementation may start.

**Prerequisite:** B1 through I2.

Review the gate checklist in the [master plan](README.md#9-frontend-start-gate), then run the
backend scenarios in the
[traceability matrix](03-requirements-traceability.md).

Record this evidence:

- clean Git commit tested by CI;
- migration version and clean-database result;
- backend image identifier;
- automated test summary;
- OpenAPI validation and generation result;
- transaction reconciliation result;
- remaining known limitations;
- final decisions D1 through D11.

**Exit gate:** Every checklist item passes, no frontend-affecting decision remains open, and the
Phase 1 API contract is stable. Only then move to F1.
