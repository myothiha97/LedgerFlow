# LedgerFlow — Architecture & Engineering Guidelines (v1)

> **What this is.** A set of concrete engineering rules for building LedgerFlow, distilled from the curated best-practice sources in [`software-development-best-practices`](https://github.com/myothiha97/software-development-best-practices) and applied directly to this project's stack (Go/Gin, React/TypeScript/Mantine, PostgreSQL, REST, Docker).
>
> That repository is a *reading list* — a collection of links to authoritative guides, not a codebase. This document does the translation work: it takes the relevant entries, turns them into actionable rules for LedgerFlow, and links each section back to the source so you can read deeper when a decision needs it.
>
> **How to use it.** Treat the numbered rules as the standard; treat the "Source / further reading" links as the justification. When a rule here conflicts with the Tech Spec, the Tech Spec wins (it's project-specific); when it conflicts with a generic external guide, this document wins (it's tailored).

---

## 1. Foundational Philosophies

These are the cross-cutting principles every other section inherits.

1. **Simplicity is the goal, not a consolation.** Build for one user and five entities (LedgerFlow's actual scale), not imagined scale. Complexity must earn its place. This is the same discipline as the BRD's "over-engineering trap" warning.
2. **Twelve-Factor config.** Store all config in the environment — DB URL, secrets, ports — never in code or committed files. This is what makes the app host-portable (the entire hosting strategy depends on it). `.env.example` documents required keys; real secrets never enter git.
3. **Conventional Commits.** Use `type(scope): subject` (e.g. `feat(transactions): add reverse-then-apply on edit`). It keeps history readable and makes changelogs/automation trivial later.
4. **Clarity over cleverness.** Code is read far more than written. Prefer the obvious implementation; leave the clever one for when profiling proves it necessary.

> **Source / further reading:** [The Twelve-Factor App](https://12factor.net/) · [Conventional Commits](https://conventionalcommits.org/) · [Google Engineering Practices](https://google.github.io/eng-practices/) · [The Zen of Go](https://the-zen-of-go.netlify.app/)

---

## 2. Architecture — Layered / Clean

LedgerFlow's backend follows a layered (lightweight clean-architecture) design. This is the single most important section — it's what makes the Phase 2 AI assistant a new caller rather than a rewrite.

1. **Dependencies point inward.** `handler → service → store/domain`. Inner layers (`domain`, `service`) never import outer ones (`handler`, framework code). The `domain` layer is pure Go — no Gin, no SQL.
2. **The service layer holds all business logic.** Validation that is *business* validation (e.g. "an expense cannot exceed... ", balance reversal rules) lives in `service/`, never in handlers. Handlers do only: parse request shape → call a service → map result/error to HTTP.
3. **One service function = one business operation.** `CreateTransaction`, `UpdateTransaction`, `RecalcBalance`, `BudgetStatus`, `DashboardSummary`. These are the API that *both* the HTTP handler and the future AI module call.
4. **Apply SOLID pragmatically.** Especially the Dependency Inversion Principle: `service` depends on a `store` *interface*, not a concrete Postgres implementation. This is what lets you test services without a database and swap storage later.
5. **Don't over-abstract.** Clean architecture at this scale means three or four layers, not a dozen. Resist adding repositories-of-repositories, generic mappers, or speculative interfaces with one implementation. (The exception worth making: the `store` interface in rule 4, because it pays off in testing immediately.)

> **Source / further reading:** [How I write HTTP services after eight years](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html) · [Modern Business Software in Go](https://threedots.tech/series/modern-business-software-in-go/) · [Wild Workouts (Go DDD example)](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) · [Mastering SOLID Principles with Go](https://packagemain.tech/p/mastering-solid-principles-with-go) · [Architecture Styles (Microsoft)](https://docs.microsoft.com/en-us/azure/architecture/guide/architecture-styles/)

---

## 3. Go Backend Conventions

1. **Follow Effective Go and a single style guide.** Adopt the Uber Go Style Guide as the house standard (it's the most prescriptive and complete). Run `gofmt`/`goimports` on save in Neovim; treat formatting as non-negotiable and non-debatable.
2. **Errors are values — wrap, don't swallow.** Return errors up the stack with context (`fmt.Errorf("create transaction: %w", err)`). Handlers translate errors to HTTP status at the boundary; lower layers never call `log.Fatal` or write HTTP responses.
3. **Accept interfaces, return structs.** Services accept the `store` interface; constructors return concrete types. Keep interfaces small and defined by the *consumer*, not the implementer.
4. **`context.Context` is the first parameter** on any function doing I/O (DB calls, future HTTP/LLM calls). Thread it from the Gin handler down through service to store so cancellation and timeouts work.
5. **Table-driven tests** for service logic (see §9). The balance-lifecycle rules are the highest-value thing to test this way.
6. **Keep `main.go` thin.** `cmd/server/main.go` is wiring only: load config, open the DB pool, build the router, start the server. No business logic, ever.

> **Source / further reading:** [Effective Go](https://go.dev/doc/effective_go) · [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) · [Google Go Style Guide](https://google.github.io/styleguide/go/index) · [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) · [Clean Go Code](https://github.com/Pungyeon/clean-go-article)

---

## 4. REST API Design

LedgerFlow's API is REST/JSON. Consistency across endpoints matters more than any single rule.

1. **Resource-oriented, plural nouns.** `/api/transactions`, `/api/accounts`, `/api/budgets`. Verbs live in the HTTP method, not the path.
2. **Correct status codes.** `200` read, `201` create, `204` delete with no body, `400` bad request shape, `401` unauthenticated, `403` unauthorized, `404` missing, `409` conflict (e.g. duplicate budget for a category+month), `422` business-rule violation, `500` unexpected.
3. **One consistent error envelope.** Every error returns the same JSON shape (e.g. `{ "error": { "code": "...", "message": "..." } }`). The frontend handles one format, not five.
4. **Validate at the edge for shape, in the service for meaning.** Gin struct-tag binding checks request *shape* (required fields, types). Business validity is checked in the service so the AI caller gets the same checks.
5. **Filtering via query params** (the BRD's transaction filters): `GET /api/transactions?category_id=...&from=...&to=...&type=expense`. Keep names consistent with the data model.
6. **The OpenAPI spec is the contract.** `backend/api/openapi.yaml` is the single source of truth that generates both the Go types and the TS client. Update the spec *before* implementing an endpoint change.

> **Source / further reading:** [Google API Design Guide](https://cloud.google.com/apis/design/) · [Microsoft API Guidelines](https://github.com/microsoft/api-guidelines/blob/vNext/Guidelines.md) · [Zalando RESTful API Guidelines](https://opensource.zalando.com/restful-api-guidelines/) · [OWASP API Security Project](https://owasp.org/www-project-api-security/)

---

## 5. Database & PostgreSQL

1. **Money is `NUMERIC`/`DECIMAL`, never `float`/`double`** — at the column level and in Go (`shopspring/decimal`). This is the rule the whole app's correctness rests on.
2. **Multi-step writes run in one transaction.** The balance reverse-then-apply (and especially an edit that moves a transaction between accounts) must be a single DB transaction so an account never sits half-updated.
3. **Index for the queries you actually run.** The hot paths are "transactions for this user, this month, this category." Add composite indexes matching those filters; don't index speculatively.
4. **Write real SQL; know what it does.** Given the learning goal, prefer `sqlc` (type-safe SQL from queries you write) over an ORM — you learn SQL and keep the `store` layer transparent. If you choose an ORM later, still be able to read the SQL it emits.
5. **Migrations are versioned and forward-only in spirit.** Every schema change is a checked-in migration file, applied in order. Never edit a shipped migration; add a new one.
6. **Heed the Postgres footguns.** Read "Don't Do This" before designing the schema — it'll save you from `timestamp` vs `timestamptz`, `char(n)`, and similar traps that are painful to undo later.

> **Source / further reading:** [SQL Style Guide](https://www.sqlstyle.guide/) · [Use The Index, Luke](https://use-the-index-luke.com) · [PostgreSQL: Don't Do This](https://wiki.postgresql.org/wiki/Don%27t_Do_This) · [4 principles of high-quality DB integration tests in Go](https://threedots.tech/post/database-integration-testing/)

---

## 6. Security

A finance app is a credentials-and-money target even with one user. Get the basics right from day one.

1. **Hash passwords with bcrypt or argon2.** Never store or log plaintext. Use a vetted library, default cost factors, and per-password salts (handled by the library).
2. **Validate and sanitize all input** at the boundary; rely on parameterized queries (sqlc/prepared statements) so SQL injection is structurally impossible. Never build SQL by string concatenation.
3. **Authn vs authz.** Authenticate every `/api/*` route except register/login. Even with one user, scope every query by `user_id` — it's the habit that makes Phase 3 multi-user safe.
4. **Secrets only via environment variables.** No keys, DB URLs, or tokens in code or git. Rotate anything that leaks immediately.
5. **HTTPS everywhere in any deployed environment**; secure, http-only cookies if you choose cookie sessions.
6. **Work the OWASP API checklist** before the first public deploy — it's short and catches the common misconfigurations.

> **Source / further reading:** [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/) · [API Security Checklist](https://github.com/shieldfy/API-Security-Checklist) · [Mozilla Web Security Guidelines](https://infosec.mozilla.org/guidelines/web_security) · [Web Developer Security Checklist](https://github.com/virajkulkarni14/WebDeveloperSecurityChecklist)

---

## 7. Frontend (TypeScript / React) Conventions

1. **TypeScript strict mode on.** For a money app, `strict: true` is non-negotiable. No `any` in domain code; types for API data come from the generated client.
2. **Feature-first folders.** `features/transactions`, `features/budgets`, etc., each owning its components, hooks, and API calls. Shared, generic UI lives in `components/`; the generated client in `api/generated/` is never hand-edited.
3. **Server state via TanStack Query; local state via React.** Don't mirror server data into Redux/global state. Mutations invalidate the relevant queries so balances/dashboard refresh automatically.
4. **Components stay small and presentational where possible.** Push data-fetching into hooks (`useTransactions`), keep components focused on rendering. Clean-code naming: a function/component name should tell you what it does.
5. **One styling paradigm — Mantine's.** No Tailwind alongside it. Use Mantine's theme for tokens (colors, spacing) so the look stays consistent without ad-hoc styles.

> **Source / further reading:** [Clean Code JavaScript](https://github.com/ryanmcdermott/clean-code-javascript) · [Project Guidelines](https://github.com/elsewhencode/project-guidelines) · [Design Patterns in TypeScript](https://github.com/RefactoringGuru/design-patterns-typescript)

---

## 8. Containerization (Docker)

1. **Multi-stage builds.** Build the Go binary in a `golang` stage, copy it into a minimal final image (`distroless` or `alpine`). The final image carries the binary, not the toolchain.
2. **Run as a non-root user** in the container. Set an explicit `USER`; don't run app processes as root.
3. **One process per container.** Backend and Postgres are separate services in `docker-compose`, not crammed into one image.
4. **`.dockerignore` aggressively** — exclude `node_modules`, build artifacts, `.git`, `.env`. Smaller context, faster builds, no secret leakage into layers.
5. **Pin base image versions**, don't rely on `latest`. Reproducible builds matter even for a learning project.
6. **`CMD` vs `ENTRYPOINT` deliberately.** Use `ENTRYPOINT` for the binary, `CMD` for default args, so the image behaves predictably.

> **Source / further reading:** [Docker development best practices](https://docs.docker.com/develop/dev-best-practices/) · [Best practices for building containers (Google)](https://cloud.google.com/architecture/best-practices-for-building-containers) · [RUN vs CMD vs ENTRYPOINT](https://www.docker.com/blog/docker-best-practices-choosing-between-run-cmd-and-entrypoint/)

---

## 9. Testing

1. **Arrange–Act–Assert** structure for every test. Keep the three sections visually distinct.
2. **Table-driven tests in Go** for service logic — they make the balance-lifecycle cases (create/edit-amount/edit-account/edit-type/delete) exhaustive and readable in one place.
3. **Test the service layer hardest.** It holds the business rules and is pure enough to test without HTTP. Mock the `store` interface (from §2 rule 4).
4. **One integration test path for the DB** to catch the things mocks hide — especially that the reverse-then-apply transaction actually commits atomically against real Postgres.
5. **Tests run in CI** (`go test ./...`) on every push once you have tests worth running (mid-Phase 1, per the Tech Spec).

> **Source / further reading:** [Arrange Act Assert](http://wiki.c2.com/?ArrangeActAssert) · [Unit Testing Best Practices](https://dzone.com/articles/unit-testing-best-practices-how-to-get-the-most-ou) · [DB integration testing in Go](https://threedots.tech/post/database-integration-testing/)

---

## 10. Git & Workflow

1. **Conventional Commits** (restated from §1 because it's a daily habit): `feat`, `fix`, `refactor`, `test`, `docs`, `chore` with a scope.
2. **Small, focused commits and PRs.** Even solo, a reviewable diff is one you can reason about and revert cleanly.
3. **Self-review before merge.** Read your own diff as if it were someone else's — the Google code-review guidelines apply even to a team of one.
4. **`main` stays green.** Don't merge anything that breaks the build or tests once CI exists.

> **Source / further reading:** [Google Engineering Practices (Code Review)](https://google.github.io/eng-practices/) · [Conventional Commits](https://conventionalcommits.org/)

---

## Appendix — Source Index

The guidelines above are derived from these entries in the curated repo, filtered to LedgerFlow's stack. The full list (including languages and tools not used here — PHP, Python, Kafka, Redis, etc.) lives at the [source repository](https://github.com/myothiha97/software-development-best-practices).

| Area | Primary sources used |
|---|---|
| Philosophy | 12-Factor App, Conventional Commits, Google Eng Practices, Zen of Go |
| Architecture | How I write HTTP services, Modern Business Software in Go, Wild Workouts DDD, SOLID with Go |
| Go | Effective Go, Uber/Google Go style guides, Clean Go Code, Go Code Review Comments |
| API | Google API Design, Microsoft API Guidelines, Zalando REST, OWASP API Security |
| Database | SQL Style Guide, Use The Index Luke, Postgres Don't Do This, Go DB integration testing |
| Security | OWASP Cheat Sheets, API Security Checklist, Mozilla Web Security |
| Frontend | Clean Code JavaScript, Project Guidelines, TS Design Patterns |
| Docker | Docker dev best practices, Google build-containers, RUN/CMD/ENTRYPOINT |
| Testing | Arrange-Act-Assert, Unit Testing Best Practices, DB integration testing |

---

*This is a living document. Revise as the stack settles and as the §13 open decisions in the Tech Spec are made.*
