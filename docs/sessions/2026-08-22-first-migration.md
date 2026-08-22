# Session Log: First Database Migration

- Date: 2026-08-22
- Objective: create and verify the first database migration.
- Outcome: completed.

## Activities

- Wrote `backend/db/migrations/0001_init.up.sql` with the `users` and `sessions` tables
  from `docs/specifications/database-design-v1.md` section 4.
- Wrote the matching `0001_init.down.sql`, dropping `sessions` before `users`.
- Applied, rolled back, and reapplied the pair to prove it is repeatable.
- Confirmed the schema against the running container rather than against the files.

## Verification

- `make migrate-up` applies cleanly.
- `\d users` shows the primary key, the unique `email` constraint, and the incoming
  cascading foreign key from `sessions`.
- `\d sessions` shows the unique `token_hash` constraint, `sessions_user_id_idx`, and the
  `ON DELETE CASCADE` foreign key to `users(id)`.
- `make migrate-down` leaves only `schema_migrations` in `\dt`.
- Reapplying `make migrate-up` succeeds. `schema_migrations` reports version `1`, not dirty.

## Concepts Practiced

- Versioned up and down migration pairs in golang-migrate
- The `schema_migrations` version and dirty flag
- Foreign key dependency fixing both creation and drop order
- PostgreSQL not indexing foreign keys automatically
- `gen_random_uuid()` as a built-in in PostgreSQL 16
- `TIMESTAMPTZ` over `TIMESTAMP`

## Next

Hand-write `backend/db/queries/users.sql` and add `backend/sqlc.yaml`, then generate the
type-safe store code with `make generate`.

## Blocked By

Nothing.
