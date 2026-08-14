# Session Log: PostgreSQL Container Ports

- Date: 2026-08-14
- Objective: make local PostgreSQL configuration consistent and verify readiness.
- Outcome: PostgreSQL started and readiness was verified; connection settings still need
  normalization.

## Activities

- Started the PostgreSQL 16 service with Docker Compose.
- Practiced Docker port publishing and the `host:container` port format.
- Configured PostgreSQL to listen on container port `5000` and published it on host port
  `6000` as a learning exercise.
- Compared attached `docker compose up` with detached `docker compose up -d`.
- Learned that Compose recreates a container when its startup command changes.
- Briefly learned how software architecture differs from infrastructure, documented in
  `docs/architecture-vs-infrastructure.md`.

## Verification

- `docker compose ps` reports the database container running with
  `0.0.0.0:6000->5000/tcp`.
- `pg_isready` reports the `ledgerflow` database accepts connections on port `5000`.
- `make test` passes. There are no Go tests yet.

## Concepts Practiced

- Container ports versus host ports
- Publishing ports with Docker Compose
- PostgreSQL startup configuration
- Attached versus detached containers
- Container recreation after configuration changes
- Architecture as system design and responsibility boundaries
- Infrastructure as the environment and tools that run and support the system

## Next

Normalize Docker Compose, the Makefile DSN, and `.env.example` to one set of PostgreSQL
connection settings, then verify a TCP connection from the host.

## Blocked By

Nothing.
