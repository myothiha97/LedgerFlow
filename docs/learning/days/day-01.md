# Day 1: Backend foundation

## What was built

- `internal/config`: reads `PORT`, `APP_ENV`, `DATABASE_URL` from env vars, checks them,
  and reports every problem at once. `APP_ENV` defaults to `prod` (fail closed).
- `internal/store/postgres.go`: opens a `pgxpool` and pings it, so the app refuses to
  start without a database.
- `internal/handler`: Problem Details errors, request ID, `slog` request logging, panic
  recovery, `GET /healthz`, a `GET /readyz` stub (your exercise), and a JSON 404.
- `server/main.go`: wiring only. Server timeouts and graceful shutdown on Ctrl+C or
  `SIGTERM`.
- Verify gate: `backend/.golangci.yml`, `govulncheck`, and `make check`.
- Security upkeep: `govulncheck` found 15 reachable vulnerabilities. Fixed by pinning
  `toolchain go1.26.10` in `go.mod` and upgrading `x/net`, `x/text`, and `quic-go`.

## How it works

A request flows through middleware in a fixed order:

```
RequestID -> Logger -> Recovery -> route handler (or NoRoute 404)
```

`RequestID` runs first, so the logger and recovery can both print the same ID. `Logger`
calls `c.Next()`, waits for the handler, then writes one line with the final status.
`Recovery` uses `defer` + `recover()`, so a panic becomes a 500 Problem response
instead of a crash.

`main` calls `run()`, which returns an error instead of calling `os.Exit`. That way
`defer pool.Close()` always runs. `signal.NotifyContext` gives a context that is
cancelled on Ctrl+C. `main` then calls `srv.Shutdown`, which stops new connections and
waits up to 10 seconds for running requests to finish.

## Patterns used

- **Dependency injection by constructor.** `NewRouter(logger, db)` gets what it needs as
  arguments. Tests pass a fake database. `main` passes the real pool.
- **Small interface where it is used.** `handler.Pinger` has one method. `*pgxpool.Pool`
  satisfies it without knowing it exists.
- **Fail fast, fail closed.** Bad config stops the app at startup. A missing `APP_ENV`
  means `prod`, never debug mode.

## Backend skills practiced

- Liveness and readiness are different: `/healthz` says "the process runs", `/readyz`
  says "I can serve traffic now". A database outage should fail readiness only, or the
  platform restarts a healthy process for nothing.
- RFC 9457 Problem Details: one error shape with a stable `code` for the frontend.
- Never leak internals: no panic values, SQL, or ping errors in responses. Logs drop
  the query string, and config errors never print `DATABASE_URL`.
- Mutation testing by hand: a reviewer broke the code on purpose. 4 of 6 breaks went
  unnoticed, and tests were added for each one.

## Where to learn more

- Liveness vs readiness probes:
  https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/
- `context.WithTimeout` and why every I/O call takes a context:
  https://go.dev/blog/context

## Your exercise: `GET /readyz`

Open `backend/internal/handler/health.go` and replace the body of `Ready`:

1. Remove the `t.Skip` line in `TestReady_DependsOnDatabase` (`handler_test.go`).
   Run `make test` and watch it fail.
2. Ping `db` with a 2-second timeout made from the request's context.
3. On success, answer 200 with `{"status":"ready"}`.
4. On failure, call `writeProblem` with 503 and `CodeUnavailable`. Do not put the ping
   error text in the response.
5. Delete the now-unused `CodeNotImplemented` constant.

Done means: `make check` passes with the test unskipped. With `make dev` running,
`curl -i localhost:8080/readyz` gives 200. After `docker compose stop db` it gives 503.

Hint: the test's fake database fails if the context has no deadline. That is how it
checks your timeout.

## Recall questions

1. Why does `/healthz` stay 200 when the database is down?
2. Why does `main` call `run()` instead of doing the work itself?
3. Why does `RequestID` ignore an `X-Request-ID` header sent by the client?

## Verify gate (2026-10-22)

```
make check
  golangci-lint fmt --diff + run      0 issues
  go vet, go build                     ok
  go test -race ./...                  config ok, handler ok
  govulncheck ./...                    Your code is affected by 0 vulnerabilities.
Smoke test (PORT=4123, real Postgres)
  GET /healthz       200 {"status":"ok"}, X-Request-Id set
  GET /readyz        501 application/problem+json (stub, until the exercise)
  GET /nope          404 application/problem+json, code not_found
  SIGINT             "shutting down" -> "server stopped cleanly"
  bad env            exit 1, lists all 3 problems
  DB unreachable     exit 1, "ping database: ...", password not in the error
Exercise test check
  Ready without timeout -> FAIL (status = 503, want 200)
  Ready with timeout    -> ok
```
