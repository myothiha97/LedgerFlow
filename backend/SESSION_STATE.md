# LedgerFlow Session State

## Current Milestone

v0.1.0 through the [10-day plan](../docs/development-plan/10-day-plan.md). Momentum Mode
is active: Claude writes about 90%, the owner does the daily exercise.

## Verified Completed

- The Go module root is `backend/`.
- The API executable remains at `backend/server/main.go`.
- `make test` passes. There are currently no Go tests.
- `make build` passes and produces the ignored `backend/bin/server` binary.
- The server reads `PORT` from the environment and defaults to `4000`.
- A runtime request using `PORT=41873` verified `GET /ping` returns HTTP `200` with
  `{"message":"pong"}`.
- Removed duplicate Gin logger and recovery registration. A runtime request now produces
  one access log entry.
- PostgreSQL 16 uses the standard `5432:5432` host-to-container mapping.
- Compose, `.env.example`, the local `.env`, and the Makefile use database name,
  username, and password `ledgerflow` on port `5432`.
- `docker compose config --quiet` passes. The database container was previously verified
  running, but Docker Desktop was stopped when checked on 2026-08-28.
- An authenticated `psql` connection through the published host port returns user and
  database `ledgerflow`.
- The existing volume password was updated without deleting its data.

- Migration `0001_init` creates `users` and `sessions` and is applied. `schema_migrations`
  reports version `1`, not dirty.
- `\d users` and `\d sessions` confirm UUID primary keys, unique `email`, unique
  `token_hash`, `ON DELETE CASCADE` from `sessions` to `users`, and `sessions_user_id_idx`.
- The up and down pair is repeatable: apply, roll back to only `schema_migrations`, reapply.

- Task 3 done (2026-10-22): `users.sql` restores `CreateUser`, adds
  `GetValidSessionWithUser` (`sqlc.embed(sessions)` plus user name and email, no
  `password_hash`, expiry by database `now()`), and `DeleteSession(token_hash, user_id)`.
  `sqlc generate` twice gives no diff, no `todo` in `gen/`, build, vet, and gofmt clean.
  Checked on real Postgres in a rolled-back transaction: valid token returns 1 row,
  expired returns 0, another user's delete removes 0, owner's delete removes 1.

## Next Outcome

Day 1: backend foundation, minus task 3. Typed config from env, pgx pool, Problem
Details errors, router, `slog` request logging with request ID, `/healthz`, server
timeouts, graceful shutdown, `main.go` as wiring only, `.golangci.yml`, `govulncheck`,
and `make check`. Owner exercise: the `GET /readyz` handler body.

## Handoff (2026-10-22)

Planning was finished in a session on the work Claude account. Work continues on the
personal account. Everything decided is on disk:

- `docs/development-plan/10-day-plan.md`: goals, decisions D1 to D11, engineering
  standards, the daily build and verify loop, sub agent model choice, the days, and BRD
  coverage.
- `CLAUDE.md`, section "Momentum Mode": how Claude works in this repo now.
- `docs/learning/backend-roadmap-2026.md`: the learning path after v0.1.0.
- `.claude/skills/diagnosing-bugs/`: the debugging skill.

These planning changes are not committed yet. First step of the next session: ask the
owner to approve one `docs:` commit for them. Leave the owner's own files alone
(`.idea/`, `backend/tasks/3.md`, `docs/learning/week-by-week-learing-roadmap.md`, and the
`users.sql`, `gen/`, and `main.go` edits that Day 1 replaces).

## Blockers

None. Docker Desktop must be running for the database: start it with `open -a Docker`,
then `make db-up`. `govulncheck` is not installed yet; Day 1 installs it with
`go install golang.org/x/vuln/cmd/govulncheck@latest`. `golangci-lint`, `migrate`, and
`air` are installed.

## Locked Decisions

- Backend module root: `backend/`
- Executable entry point: `backend/server/`
- Database access: `sqlc`
- Migrations: `golang-migrate`
- Phase 1 excludes AI features.
- v0.1.0 decisions D1 to D11 are in `docs/development-plan/10-day-plan.md`.

## Parking Lot

- Owner's AI workflow notes (from the `main.go` comment): check the Cursor cloud product,
  review progress and architecture, continue an advanced Go tutorial.
- AWS track after v0.1.0: $1000 Activate credits expire 2027-04-30, main region
  `ap-southeast-1`. Read-only CLI profile `agent-readonly` exists for Claude.
