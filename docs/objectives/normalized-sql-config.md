Active objective: normalize PostgreSQL configuration. This unlocks migrations and database code.

Done means:

- Compose matches the DATABASE_URL in .env.example and Makefile.
- PostgreSQL uses its default internal port.
- docker compose config passes.
- An authenticated host connection succeeds.
- make test passes.

Your step: update docker-compose.yml. Compare every port and credential with .env.example. Remove any custom PostgreSQL port override.

Footgun: the existing pgdata volume may retain the old password. Do not delete it yet.

Docker Desktop is currently stopped. Make the edit, then tell me when ready for review.
