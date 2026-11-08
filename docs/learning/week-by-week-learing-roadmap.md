# Go Backend Roadmap for Myo (Oct 15 – early Dec 2026): Are the Official Docs, Free Resources and Alex Edwards' Two Books Enough?

**Short answer:** Both Alex Edwards books are reliable, up to date for Go 1.27, and the best purchase you can make for the "build a production JSON API in Go" part of your goal. Buy both. Together with the official docs they are **not enough** for a backend, cloud and systems foundation. They leave four gaps that matter for you: (1) Go language depth, meaning concurrency, interfaces, generics and testing an API; (2) containers and CI; (3) one real cloud provider; and (4) observability and reliability thinking. Close those gaps with one extra Go book (_Learning Go_, 2nd ed.), the free 100go.co site, free official docs (Docker, AWS, OpenTelemetry, Postgres), and the free Google SRE book. Defer everything else.

## TL;DR

- **The books are reliable and current, but neither one covers everything.** _Let's Go_ ($34.95) and *Let's Go Further* ($59.95) are both updated for Go 1.27. _Further_ was last updated Oct 3, 2026, and Edwards' FAQ says updates usually come every 6 months, timed to Go releases. They have no Docker, no OpenTelemetry, no cloud-provider content, and _Further_ has no testing chapter. _Further_ also still uses the third-party `httprouter` rather than the Go 1.22+ standard library `ServeMux`.
- **Buy _Let's Go_ even though you're building a JSON API, but read it selectively.** _Further_ says it is "not designed for newcomers to Go" and expects you to already know how to build and test database-backed web apps. Testing is only taught in _Let's Go_ (Ch. 13). Skip its HTML templates, forms and sessions chapters.
- **Plan: seven weeks, about 75/25 backend-cloud-systems vs. DevOps, with LedgerFlow as the only project.** Weeks 1–3 cover Go plus Postgres API fundamentals and tests. Week 4 is Docker and GitHub Actions CI. Week 5 is auth, security and an AWS deploy. Week 6 is observability and reliability. Week 7 is a no-AI rebuild and consolidation. Paid resources come to $151.89 (*Let's Go* $34.95 + _Let's Go Further_ $59.95 + the *Learning Go* ebook at $56.99). Kubernetes, microservices, a second cloud and the AI/RAG layer wait until December.

## Key Findings

### 1. How reliable and current are the two books?

**Currency (verified on the books' own sites):**

|            | _Let's Go_                                                      | _Let's Go Further_                                                                             |
| ---------- | --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| Price      | $34.95                                                          | $59.95                                                                                         |
| Go version | "correct for Go version 1.27 and newer"                         | Written for Go 1.27; last updated Oct 3rd, 2026; version 1.27.1                                |
| Size       | —                                                               | 614 pages (PDF)                                                                                |
| App built  | _Snippetbox_, a server-rendered HTML app (MySQL)                | _Greenlight_, a JSON API (PostgreSQL)                                                          |
| Routing    | stdlib `http.NewServeMux` with Go 1.22 method/wildcard patterns | third-party `httprouter`                                                                       |
| Updates    | included                                                        | "usually sent out once every 6 months, to coincide with new major releases of the Go language" |
| Refund     | 30-day                                                          | 30-day money-back guarantee                                                                    |

These are actively maintained books, not stale ones. The _Further_ sample already uses `log/slog`, and its JSON chapter now teaches `encoding/json/v2`, which became generally available in Go 1.27. On Sept 18, Edwards announced on X that both books were "fully updated for Go 1.27 – including coverage of json/v2 in Let's Go Further."

**Reputation:**

- _Let's Go Further_ has a 4.73 average on Goodreads (97 ratings).
- John Arundel (Bitfield Consulting) calls it "both an essential reference for Go developers and a valuable compendium of battle-tested knowledge." He says its biggest value is the advice on "common traps, gotchas, and footguns." He also lists the pair as the best books for web services in his 2026 roundup.
- Hacker News "how do I learn Go" threads regularly recommend both books.

**Criticisms worth taking seriously:**

- **Audience.** Arundel says _Further_ is aimed "squarely at fairly experienced developers: those with at least one or two years of commercial Go development," and that you "should read [Let's Go] first." You have neither, which is why the plan sequences the books carefully.
- **Code style.** Arundel found the code "heavily commented: so much so, indeed, that it's hard to follow the code sometimes." The book is also clearly built for copy-paste-along reading. That is a real risk for someone who wants genuine understanding.
- **Structure.** A 2019 HN commenter complained that "the structure is so so and everything is in package main." That's an older edition, but the `cmd/api` + `internal/` layout is still deliberately flat and pragmatic. It is not a "clean architecture" or domain-driven layout.
- **Router choice.** _Further_ §2.3 still uses `httprouter` instead of `http.ServeMux`. Edwards gives two reasons: ServeMux's plaintext 404/405 responses are hard to customize (there is an open Go proposal about this), and it doesn't handle OPTIONS automatically. This is a defensible choice, not an outdated one. **For LedgerFlow, use the stdlib ServeMux anyway** and write JSON 404/405 handling yourself. Translating the book's code is exactly the kind of friction that forces understanding.
- **Page-count discrepancy.** Arundel's older review says "well over 1000 pages." The current site lists 614 PDF pages. The difference is probably formatting or edition; it doesn't change the recommendation.

**Coverage gaps (checked against the current tables of contents):**

- _Further_ has no Docker chapter. Deployment is a DigitalOcean droplet running the API as a systemd service behind Caddy (Ch. 20).
- No OpenTelemetry. Observability is limited to `expvar` metrics (Ch. 18).
- No testing chapter. The closest is §19.3 "Quality-controlling code."
- Neither book covers Kubernetes, a cloud provider's managed services, distributed systems, Go runtime internals or deep concurrency.
- _Let's Go_ uses MySQL, not Postgres.

**Do you need _Let's Go_ before _Further_? Yes, but read it as a fast 2–3 day skim-and-build.** Three reasons:

1. _Further_ explicitly reuses _Let's Go_'s patterns.
2. The only testing instruction in either book is _Let's Go_ Ch. 13: unit and sub-tests, testing handlers and middleware, end-to-end tests, mocking, integration testing and coverage.
3. _Let's Go_ teaches the Go 1.22+ ServeMux, middleware chaining, request context and server security hardening, which _Further_ assumes you know.

Skip _Let's Go_ Ch. 5 (HTML templates), Ch. 7 (forms), most of Ch. 8 (sessions) and Ch. 12 (file embedding).

### 2. Are the docs, free resources and two books enough? No. Here's what's missing and what fills each gap

| Gap                                                                             | Why it matters for you                                                                        | Best fix (reading-first)                                                                                                                                                                                                                                                                                                    | Cost                                                                         |
| ------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Go language depth (interfaces, errors, generics, concurrency, context, testing) | Tour + part of Effective Go isn't enough to understand _why_ Edwards' code is shaped as it is | **_Learning Go_ 2nd ed. (Bodner, O'Reilly, Jan 2024, 494 pp.)**. Arundel calls it the best one-volume reference. Key chapters: 7 Types/Methods/Interfaces, 9 Errors, 12 Concurrency, 13 Standard Library (incl. net/http and structured logging), 14 Context, 15 Writing Tests (httptest, integration tests, race detector) | Paid: $56.99 DRM-free ebook (ebooks.com); $65.99 print list price (O'Reilly) |
| Idiom and bug-avoidance review                                                  | You need an external yardstick for code review when not leaning on AI                         | **100go.co** (free summary of _100 Go Mistakes_, Manning 2022). Book is optional                                                                                                                                                                                                                                            | Free                                                                         |
| Containers                                                                      | Required for CI, deployment, and later Kubernetes                                             | Docker official docs (free) + optionally _**Docker Deep Dive**_ (Poulton, May 2025 ed.)                                                                                                                                                                                                                                     | Free / $9.99 Kindle, $19.99 minimum on Leanpub, $49.99 paperback list        |
| CI/CD                                                                           | DevOps share of your goal                                                                     | GitHub Actions docs (free) + _Further_ Ch. 19 (Makefile, audit, versioning)                                                                                                                                                                                                                                                 | Free                                                                         |
| One cloud                                                                       | Your stated spine                                                                             | **AWS** official docs/tutorials (free). See cloud choice below                                                                                                                                                                                                                                                              | Free tier                                                                    |
| Observability                                                                   | "Highly reliable production systems" is impossible without it                                 | OpenTelemetry Go docs, Prometheus docs, `log/slog` docs, `net/http/pprof` (all free)                                                                                                                                                                                                                                        | Free                                                                         |
| Reliability thinking                                                            | The tech-lead trajectory                                                                      | **Google SRE book** (free online): Ch. 4 SLOs, Ch. 6 Monitoring, Ch. 21–22 Overload/Cascading failures                                                                                                                                                                                                                      | Free                                                                         |
| Postgres                                                                        | Ledger correctness = transactions, constraints, indexes                                       | Postgres official docs (tutorial, transaction isolation, indexes) + **use-the-index-luke.com** (free)                                                                                                                                                                                                                       | Free                                                                         |
| Linux/Bash                                                                      | Spine skill; you deploy to a VM                                                               | _The Linux Command Line_ (Shotts, free PDF) – shell, permissions, processes, env, scripting basics                                                                                                                                                                                                                          | Free                                                                         |
| Security                                                                        | Money data + auth                                                                             | OWASP API Security Top 10 (free) + _Further_ Ch. 15–17                                                                                                                                                                                                                                                                      | Free                                                                         |
| Data-system design                                                              | Long-term systems depth                                                                       | **_Designing Data-Intensive Applications_ 2nd ed.** (Kleppmann & Riccomini; Kleppmann's site dates it March 2026, 670 pp.; retailers list Feb–Mar 2026). Start in Week 7, finish after November                                                                                                                             | Paid: $59.99 ebook (Apple Books); $69.99 paperback (Books-A-Million)         |

**Must-have (buy or use now):** _Let's Go_, _Let's Go Further_, _Learning Go_ 2e, 100go.co, official Go docs and release notes, Docker docs, AWS docs, OpenTelemetry Go docs, Postgres docs + use-the-index-luke, SRE book (selected chapters), OWASP API Top 10.

**Nice-to-have (only if a specific gap hurts):**

- _Docker Deep Dive_ 2025 edition.
- _100 Go Mistakes_ in print, if you prefer the full text to the website.
- _Cloud Native Go_ 2nd ed. (Titmus, O'Reilly, Oct–Dec 2024, 539 pp.). It's the best "resilience patterns in Go" book (circuit breakers, retries, a distributed key-value store), but it overlaps heavily with what you'll build. Read it in December.

**Skip for now:**

- _**The Go Programming Language**_ (Donovan & Kernighan, 2015). Still excellent prose, but it predates modules, generics, `slog`, iterators and the new ServeMux. An HN asker in 2024 pointed out it was already 8 years old. Use _Learning Go_ instead.
- _**Concurrency in Go**_ (Cox-Buday, 1st ed., July 2017, 238 pp.; no newer edition found). The patterns are timeless, but the code predates generics and modern `context`/`errgroup` usage. Defer it until you have a real concurrency problem.
- _**Efficient Go**_ (Płotka). This is performance engineering. Premature for you; defer to 2027.
- **Ardan Labs Ultimate Go / gotraining.** Your decision to skip the $699 video bundle is correct. Bill Kennedy's free written material is a fine reference on mechanics, but don't make it a track.
- **ByteByteGo / _System Design Interview_.** That's interview prep, not foundations. DDIA 2e comes first.
- **Kubernetes books, _Kubernetes Up & Running_, any second cloud.** See the defer list.

**Cloud choice: AWS.** JobStreet Singapore showed 7,796 "AWS cloud engineer" results in Oct 2026. A 2026 analysis of 200k backend listings found AWS in 55% of postings. That source is weaker: a vendor blog with unaudited methodology. One Medium piece claims AWS and GCP are "near parity" in Southeast Asia, but only for _data engineering_ postings, and it's a low-quality source. GCP is a valid choice only if your employer runs on it. Check that first. If they're on GCP, pick GCP, because daily exposure beats market share.

### 3. Recent Go changes that define "modern Go" in late 2026

Go 1.27.1 (Sept 1, 2026) is the current stable release; Go 1.27 was released on August 19, 2026, per the official Go Blog announcement. What matters for a backend learner:

| Version             | Change                                                                                                                                                                                                                                                                                                                                                                | Practical impact                                                                                                                             |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| 1.21                | `log/slog` structured logging                                                                                                                                                                                                                                                                                                                                         | Use it from day 1 (both books do)                                                                                                            |
| 1.22                | Method + wildcard routing in `http.ServeMux` (`"GET /v1/tx/{id}"`); per-iteration loop variables                                                                                                                                                                                                                                                                      | Third-party routers optional; old "capture the loop var" advice is obsolete                                                                  |
| 1.23                | Range-over-func iterators                                                                                                                                                                                                                                                                                                                                             | Know how to read them; rarely needed in handlers                                                                                             |
| 1.25                | Container-aware `GOMAXPROCS` on Linux (reads cgroup CPU _limits_, adjusts periodically); `testing/synctest` GA; `sync.WaitGroup.Go`; `net/http.CrossOriginProtection`                                                                                                                                                                                                 | `uber-go/automaxprocs` is mostly redundant. Set CPU _limits_, not just requests, or Go can't see them                                        |
| 1.26 (Feb 10, 2026) | `new(expr)`; Green Tea GC on by default; rewritten `go fix` with "modernizers"; experimental goroutine-leak profile                                                                                                                                                                                                                                                   | Run `go fix ./...` to modernize code; treat AI-generated older-style Go with suspicion                                                       |
| 1.27 (Aug 2026)     | **Generic methods**; `encoding/json/v2` GA, and `encoding/json` itself is now backed by v2 (unmarshal "significantly faster"; error message text may differ); goroutine-leak profile GA at `/debug/pprof/goroutineleak`; `go test` runs the `stdversion` vet check; `go vet` printf check flags `%w` wrapping a pointer-to-error-type; `go doc pkg@version` and `-ex` | Learn json/v2 semantics (stricter: rejects invalid UTF-8 and duplicate names; case-sensitive field matching). Use the leak profile in Week 6 |

Go 1.27 also adds a standard-library `uuid` package. The official Go Blog release post says it "provides native support for generating and parsing UUIDs," and VictoriaMetrics' Go 1.27 tour notes that it offers `uuid.New()`, `NewV4()` and `NewV7()` per RFC 9562. For new code you can drop `google/uuid`.

**Outdated-resource flags:**

- Anything teaching `gorilla/mux` as necessary, `automaxprocs` as mandatory, the loop-variable capture idiom, `ioutil`, or `log` instead of `slog` is pre-2024 thinking.
- _The Go Programming Language_ and _Concurrency in Go_ predate generics.
- _100 Go Mistakes_ (2022) predates 1.22. Its loop-variable mistake no longer applies to modules declaring `go 1.22+`, and a few stdlib recommendations have newer alternatives.
- AI assistants commonly produce Go that "looks a few years old" (JetBrains' phrasing). That's another reason to run `go fix` and read the release notes yourself.

## Week-by-Week Roadmap (LedgerFlow as the only build vehicle)

**Assumptions.** About 5–6 focused hours per day. Each week is roughly 60% reading-while-building, 30% building alone, and 10% checkpoints. DevOps-heavy work is Week 4 plus parts of Weeks 5–6, which comes to about 25% overall.

**LedgerFlow domain rules, fixed from day 1 because they force real backend thinking:**

- Money is stored as `BIGINT` minor units, never floats.
- Every balance-affecting write happens inside a DB transaction.
- Transfers are double-entry: two rows that must sum to zero.
- `POST` endpoints accept an `Idempotency-Key` header.

### Week 1 (Oct 15–21): Go language depth and the HTTP basics

- **Read:**
  - _Learning Go_ 2e Ch. 1–11. Skim Ch. 1–4. Read Ch. 5–10 carefully: functions, pointers, types/methods/interfaces, generics, errors, modules.
  - _Let's Go_ Ch. 2 (Foundations), Ch. 3 (Configuration and error handling), Ch. 4 (Database-driven responses). Do Ch. 4 with **Postgres + pgx** instead of MySQL.
  - Go 1.22 and 1.27 release notes.
- **Linux:** _The Linux Command Line_ Part 1 (shell, navigation, permissions, processes, environment), about 30 min per day.
- **Build:** LedgerFlow skeleton.
  - `cmd/api` + `internal/` layout.
  - Config from flags/env.
  - `slog` JSON logger.
  - Stdlib ServeMux with `GET /v1/healthcheck`.
  - Postgres running locally (Docker is fine as a black box for now), plus one `accounts` table read path.
- **Checkpoint (no AI):** Explain on paper: interface satisfaction, pointer vs. value receivers, `errors.Is`/`As` with wrapping, and what a `*sql.DB` actually is (a pool, not a connection).

### Week 2 (Oct 22–28): JSON API and Postgres properly

- **Read:**
  - _Further_ Ch. 2–7: getting started, sending JSON, parsing JSON, database setup and connection-pool configuration, SQL migrations, CRUD.
  - Appendices 21.4–21.5 (JSON encoding/decoding nuances).
  - _Learning Go_ Ch. 13 (Standard Library) and Ch. 14 (Context).
  - Postgres docs: tutorial, plus the "Transaction Isolation" and "Indexes" chapters.
  - use-the-index-luke: the B-tree and WHERE-clause chapters.
- **Build:**
  - Accounts and transactions CRUD with JSON envelopes and consistent JSON error responses, including your own JSON 404/405 handlers.
  - Input validation and request size limits.
  - Migrations (`golang-migrate` or `goose`).
  - Query timeouts via `context`.
  - A `POST /v1/transfers` that does a double-entry insert in one DB transaction.
- **Checkpoint:** Explain the difference between READ COMMITTED and SERIALIZABLE using your own transfer endpoint as the example. Explain when `SetMaxOpenConns`/`SetMaxIdleConns` matter, and why.

### Week 3 (Oct 29–Nov 4): Correctness under load and testing

- **Read:**
  - _Further_ Ch. 8 (Advanced CRUD, especially optimistic concurrency control), Ch. 9 (filtering/sorting/pagination), Ch. 10 (rate limiting), Ch. 11 (graceful shutdown).
  - _Let's Go_ Ch. 6 (Middleware), Ch. 9 (Server and security improvements), Ch. 11 (Request context), **Ch. 13 (Testing)**.
  - _Learning Go_ Ch. 12 (Concurrency) and Ch. 15 (Writing Tests).
- **Build:**
  - Version-column optimistic locking on account updates.
  - Paginated transaction listing.
  - Per-IP rate limiter middleware.
  - Panic recovery and request-ID middleware.
  - Graceful shutdown that drains in-flight requests and background goroutines.
  - Idempotency keys stored in a table with a unique constraint.
  - **Tests:**
    - table-driven unit tests;
    - `httptest` handler tests;
    - integration tests against a real Postgres (Testcontainers or a compose service);
    - a concurrent-transfer test run with `go test -race` that proves balances never go wrong.
- **Checkpoint:** First pass of your code against the 100go.co chapters on error management, concurrency and testing. Log every hit in a `MISTAKES.md`.

### Week 4 (Nov 5–11): DevOps block, containers and CI (~80% DevOps this week)

- **Read:**
  - Docker docs: "Get started," then multi-stage builds, Compose, and build best practices.
  - Optionally _Docker Deep Dive_ chapters on images, containers, containerizing apps and Compose.
  - _Further_ Ch. 19: Makefiles, quality control (`go vet`, `staticcheck`, `govulncheck`), versioning via Git.
  - GitHub Actions docs: workflow syntax, service containers, caching.
- **Build:**
  - Multi-stage Dockerfile producing a static binary (`CGO_ENABLED=0`) on a distroless or scratch base, running as non-root.
  - `compose.yaml` with api + postgres + migrations.
  - GitHub Actions pipeline: `go vet` → `staticcheck` → `govulncheck` → `go test -race` with a Postgres service container → build and push the image to GHCR on `main`.
- **Checkpoint:** Answer without AI or notes: why multi-stage, why non-root, what a layer cache invalidation is, and what `GOMAXPROCS` your container will see with and without a CPU limit (the Go 1.25 behaviour).

### Week 5 (Nov 12–18): Auth, security and the first real cloud deploy

- **Read:**
  - _Further_ Ch. 12 (user model), Ch. 15 (authentication), Ch. 16 (permission-based authorization), Ch. 17 (CORS), Ch. 20 (deployment and hosting). Skim Ch. 13–14 on email and activation.
  - OWASP API Security Top 10.
  - AWS docs: IAM basics (users vs. roles, least privilege), VPC and security groups, EC2, RDS for PostgreSQL, SSM Parameter Store/Secrets Manager.
- **Build:**
  - Token auth, and per-user data isolation on every query (the classic BOLA/IDOR bug in a finance app).
  - Deploy on **AWS** following _Further_ Ch. 20's VM pattern, adapted from DigitalOcean: one EC2 instance running the container (or systemd + binary) behind Caddy for TLS, with **RDS Postgres** in a private subnet and secrets in Parameter Store.
  - Set an AWS billing alarm on day 1.
- **Checkpoint:** Draw the network path of one request from your laptop to RDS and back: DNS, TLS, security groups, subnets, pool. If you can't draw it, you don't understand it yet.

### Week 6 (Nov 19–25): Observability and reliability

- **Read:**
  - _Further_ Ch. 18 (Metrics).
  - OpenTelemetry Go "Getting Started" plus the `otelhttp` instrumentation docs.
  - Prometheus "Overview" and "Instrumentation" best-practices pages.
  - `net/http/pprof` docs and the Go 1.27 goroutine-leak profile notes.
  - SRE book Ch. 4 (Service Level Objectives), Ch. 6 (Monitoring Distributed Systems), Ch. 21–22 (Handling Overload; Addressing Cascading Failures).
- **Build:**
  - Replace or augment `expvar` with OTel traces (HTTP → DB spans) and Prometheus RED metrics (rate, errors, duration).
  - Correlate `slog` logs with trace IDs.
  - Define two SLOs for LedgerFlow (e.g., 99.5% of `POST /transfers` succeed in under 300 ms).
  - Add server timeouts (`ReadHeaderTimeout`, etc.) and DB timeouts. Load-test with `hey` as the book does.
  - Profile once with pprof and inspect `/debug/pprof/goroutineleak` after deliberately leaking a goroutine.
  - Add one background worker (e.g., recurring-transaction scheduler) using `context` cancellation and `errgroup`.
- **Checkpoint:** Second full 100go.co review pass. Write a one-page postmortem of an outage you inject yourself, such as killing the DB mid-load.

### Week 7 (Nov 26–Dec 2): Consolidation, no-AI rebuild, and starting systems thinking

- **No-AI rebuild:** From an empty directory, in one day, with only `pkg.go.dev` and the books, rebuild a mini-service: one resource, Postgres, migrations, tests, Dockerfile, graceful shutdown. Every place you get stuck goes into your study list for December.
- **Read:**
  - Start DDIA 2e: the opening chapters on trade-offs, reliability/scalability/maintainability, and data models/storage.
  - Re-read your `MISTAKES.md`.
- **Write:**
  - An architecture doc / ADRs for LedgerFlow: why stdlib routing, why BIGINT money, why optimistic locking, why EC2 over ECS for now.
  - A README with a runbook.
- **Optional:** Read the `net/http` server source path from `ListenAndServe` → `conn.serve` → `ServeHTTP`. Reading stdlib source is now a habit you can sustain.

### Using Claude Code without hollowing out your learning

- **First implementation is always yours.** For each book chapter, type the code yourself, adapted to LedgerFlow. Don't paste. Edwards' heavy inline comments make copy-paste tempting, and that's the main way these books fail learners.
- **Use AI as a reviewer and examiner, not an author.** Good prompts: "Review this handler against 100 Go Mistakes; list issues but don't fix them," "Quiz me on Postgres isolation levels," "What failure modes does this transfer code have?"
- **AI is allowed to generate** boilerplate you've already written once by hand (second migration, second CRUD resource), CI YAML after you've written the first version yourself, and test-data fixtures.
- **Verify AI-generated Go against 1.25–1.27 idioms.** Run `go fix ./...` and `go vet`, and reject pre-1.22 patterns.
- **Weekly no-AI hour:** a fixed slot to write a feature or fix a bug with AI disabled. It's your honest progress meter.
- **Explain-it-back:** End each week with a 300-word written explanation of the hardest concept. If you can't write it, re-read the chapter.

### What NOT to do in these seven weeks

- No Kubernetes, Helm, service mesh or EKS. Containers and one VM are enough to learn the fundamentals. Go 1.25's cgroup awareness is the only K8s-adjacent fact you need right now.
- No microservices, Kafka or event sourcing. One well-built monolith teaches more.
- No second cloud, no Terraform deep dive. At most, write a small Terraform file for the EC2 + RDS setup in Week 7 if you're ahead.
- No tutorial hopping or video courses. No switching to Gin, Echo or Fiber.
- No sqlc/ORM rabbit holes before you've written raw SQL with `database/sql`/pgx by hand. Adding sqlc in December is fine.
- No AI/RAG features in LedgerFlow yet. A reliable data layer has to come first.
- No reading _Concurrency in Go_, _Efficient Go_, _The Go Programming Language_, or system-design interview books cover to cover.

## What to Defer Until After November

- **December–January:**
  - finish DDIA 2e;
  - _Cloud Native Go_ 2e (resilience patterns);
  - sqlc;
  - Terraform for your AWS setup;
  - moving from EC2 to ECS Fargate;
  - an AWS Solutions Architect Associate–level study, if credentials matter for your job search.
- **Q1 2027:**
  - Kubernetes (official docs + _Kubernetes Up & Running_, deploying LedgerFlow to a local kind cluster, then EKS);
  - _Concurrency in Go_ alongside a real concurrent feature;
  - the AI-systems layer on LedgerFlow (RAG over transactions, evals, a Python sidecar where it fits).
- **Later:**
  - _Efficient Go_, Go runtime internals, distributed-systems theory beyond DDIA;
  - system-design interview prep (Alex Xu/ByteByteGo) only when you're actually interviewing;
  - IoT/robotics.

## Total Cost

| Item                                                                                                                       | Status       | Price                                                                                                 |
| -------------------------------------------------------------------------------------------------------------------------- | ------------ | ----------------------------------------------------------------------------------------------------- |
| _Let's Go_                                                                                                                 | Must         | $34.95                                                                                                |
| _Let's Go Further_                                                                                                         | Must         | $59.95                                                                                                |
| _Learning Go_ 2e (ebook or print)                                                                                          | Must         | $56.99 DRM-free ebook (ebooks.com); $65.99 print list price (O'Reilly)                                |
| 100go.co, Go docs, Docker/AWS/OTel/Prometheus/Postgres docs, use-the-index-luke, SRE book, _The Linux Command Line_, OWASP | Must         | Free                                                                                                  |
| _Docker Deep Dive_ (2025 ed.)                                                                                              | Optional     | $9.99 Kindle; $19.99 minimum on Leanpub; $49.99 paperback list (Amazon)                               |
| DDIA 2e                                                                                                                    | Start Week 7 | $59.99 ebook (Apple Books); $69.99 paperback (Books-A-Million)                                        |
| AWS usage                                                                                                                  | Budget       | Free tier where eligible; RDS/EC2 beyond free tier can cost a few dollars a week. Set a billing alarm |

**Core spend: $94.90 for the Edwards books, or $151.89 with the _Learning Go_ ebook ($160.89 with the print edition).** Adding the DDIA 2e ebook ($59.99) brings the total to about $212. That is still well under a third of the Ardan bundle you rejected. Free alternatives of comparable quality: 100go.co instead of the _100 Go Mistakes_ book; the Go blog plus release notes plus Effective Go instead of _The Go Programming Language_. Check whether your employer's learning budget covers O'Reilly online access. If it does, _Learning Go_, DDIA 2e, _Cloud Native Go_ and _Concurrency in Go_ cost you nothing.

## Caveats

- Prices were checked at one retailer each in October 2026 (Edwards' own sites, ebooks.com, Apple Books, Books-A-Million, Leanpub, Amazon) and vary by store and region. AWS costs are estimates.
- DDIA 2e publication dates differ by source (Feb 17–Mar 24, 2026). Kleppmann's own page says March 2026.
- The Go 1.27 standard-library `uuid` package is confirmed in the official Go Blog release post.
- Cloud job-share figures come from vendor and blog sources with unclear methodology. Treat "AWS has the most demand" as directionally right, not precise.
- The seven-week plan is aggressive. If you slip, cut Week 6's OpenTelemetry tracing (keep metrics and logs) and the Week 7 DDIA start, not the testing or the no-AI rebuild.
