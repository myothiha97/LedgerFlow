# LedgerFlow Session State

## Current Milestone

Backend foundation.

## Verified Completed

- The Go module root is `backend/`.
- The API executable remains at `backend/cmd/server/main.go`.
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
- `docker compose config --quiet` passes and the database container is running.
- An authenticated `psql` connection through the published host port returns user and
  database `ledgerflow`.
- The existing volume password was updated without deleting its data.

- Migration `0001_init` creates `users` and `sessions` and is applied. `schema_migrations`
  reports version `1`, not dirty.
- `\d users` and `\d sessions` confirm UUID primary keys, unique `email`, unique
  `token_hash`, `ON DELETE CASCADE` from `sessions` to `users`, and `sessions_user_id_idx`.
- The up and down pair is repeatable: apply, roll back to only `schema_migrations`, reapply.

## Next Outcome

Hand-write `backend/db/queries/users.sql`, add `backend/sqlc.yaml`, and generate the
type-safe store code into `backend/internal/store/gen` with `make generate`.

## Blockers

None.

## Locked Decisions

- Backend module root: `backend/`
- Executable entry point: `backend/cmd/server/`
- Database access: `sqlc`
- Migrations: `golang-migrate`
- Phase 1 excludes AI features.
