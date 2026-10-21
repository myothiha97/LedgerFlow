# Backend Roadmap 2026 (LedgerFlow edition)

> Researched 2026-10-21. Topic list taken from the roadmap.sh backend roadmap source
> (content synced 2026-10-16). Book and course choices live in
> `week-by-week-learing-roadmap.md`. This file is about topics and order.

The roadmap has about 150 topics. You do not learn them as a list. You learn them in
levels. Each level is a project that forces you to use the topics for real.

| Level | Project | Goal |
| --- | --- | --- |
| 1 | LedgerFlow v0.1.0 (10-day plan) | Core backend: HTTP, SQL, transactions, auth, testing, logs |
| 2 | LedgerFlow v0.2 production track | Run it like a real product: cache, jobs, metrics, cloud, IaC |
| 3 | A second, distributed project | Many services: events, gRPC, Kubernetes, failure handling |
| 4 | Awareness only | Know what it is and when to pick it. Learn deeply only when a job needs it |

## What changed in 2025 and 2026

| Area | Update | What it means for you |
| --- | --- | --- |
| Architecture | The industry default is now a **modular monolith**: one deployable app with strict internal boundaries. Teams split out services only with evidence (team size, scaling need). | LedgerFlow's layered monolith is the modern choice, not a beginner shortcut. Learn microservices in Level 3, on purpose. |
| Go | **Go 1.27** (2026-08-19): `encoding/json/v2`, a standard `uuid` package, generic methods, the `goroutineleak` profile, `httptest.NewTestServer`. | LedgerFlow's Go 1.26 is still supported. Upgrade after v0.1.0. |
| PostgreSQL | **18** is the current stable (async I/O, `uuidv7()`). **19** is at RC 1, GA planned 2026-11-19: `REPACK CONCURRENTLY`, `pg_plan_advice`, `WAIT` for read-your-writes on replicas. | PostgreSQL 16 is fine for v0.1.0. Upgrade in Level 2 and practice a real version upgrade. |
| Security | **OWASP Top 10:2025** (final January 2026). New: A03 Software Supply Chain Failures, A10 Mishandling of Exceptional Conditions. SSRF moved into A01. | Dependency scanning, pinned versions, and safe error handling are now top-10 risks. LedgerFlow covers them with `govulncheck`, pinned images, and one error format. |
| Observability | **OpenTelemetry** is the vendor-neutral standard. In Go, traces and metrics are stable, logs are beta (`otelslog` bridge), profiles are early. | Use `slog` for logs and OpenTelemetry for traces and metrics. LedgerFlow Day 9 adds traces. |
| Kubernetes | **ingress-nginx is retired** (no fixes after March 2026). The Gateway API replaces Ingress. | When you reach Kubernetes, learn the Gateway API, not Ingress. |
| Workflows | **Durable execution** engines (Temporal, Restate, DBOS) are growing fast. | Learn the transactional outbox and idempotent consumers first. The engines automate these ideas. |
| APIs | **RFC 9457 Problem Details** is the standard JSON error format. `Idempotency-Key` is common practice (Stripe style), but its IETF spec is still a draft. | LedgerFlow uses both. |
| AI | The roadmap now has an AI section: how LLMs work, function calling, structured outputs, embeddings, vectors, RAG, agents, MCP, AI-assisted coding. | LedgerFlow Phase 2 (text to transaction, money questions) is where you learn this, on top of the same service layer. |

## Level 1: learned in LedgerFlow v0.1.0

| Roadmap topic | Where you use it |
| --- | --- |
| HTTP, REST, JSON APIs, status codes | Days 1, 3, 4 |
| Go, Git, GitHub | Every day |
| PostgreSQL, migrations, normalization, constraints, indexes | Days 1, 4 to 8 |
| ACID, transactions, row locking | Days 6, 7 |
| N+1 problem (aggregate in SQL, not in loops) | Day 8 dashboard |
| Cookie-based auth, token auth, bcrypt, SHA hashing | Days 2, 3 |
| OWASP risks: access control, injection, auth failures | Days 2 to 8 (`user_id` scoping, `sqlc` parameters, same-error login) |
| Unit, integration, functional testing | Days 2, 6, 7, 8, 9 |
| Twelve-factor app (env config, logs to stdout, fast start and stop) | Days 1, 9 |
| Monolith (modular) | Whole project |
| Throttling and rate limiting | Day 9 |
| Instrumentation, telemetry, observability | Days 1 (`slog`), 9 (OpenTelemetry) |
| Graceful degradation, failure modes | Days 1 (`/readyz`), 9 (`429`, idempotent retry) |
| Containers, CI/CD | Day 9 |

## Level 2: LedgerFlow v0.2 production track

Do these in order after v0.1.0. Each one is 1 to 3 sessions. Security and deploy come
first because LedgerFlow has two goals: be useful in your real life, and teach you. You
only get both when you use it every day from your phone. Real daily use also gives real
data and real problems for every later step.

1. **Security for the public internet:** CSRF protection, security headers and CSP,
   secrets out of env files, an OWASP ASVS Level 1 review, SBOM and image scanning.
2. **Cloud deploy:** one provider (AWS), infrastructure as code with OpenTofu or Terraform,
   managed PostgreSQL, TLS, a CD pipeline, backups with a restore drill, and zero-downtime
   migrations (expand, migrate, contract). From here on, you use LedgerFlow daily.
3. **Metrics and alerts:** OpenTelemetry metrics to Prometheus, a Grafana dashboard with
   request rate, errors, and duration (the "RED" metrics), one alert, and one SLO.
4. **Upgrade** to Go 1.27 and PostgreSQL 18 or 19. Practice a real upgrade with a dump and
   restore, now with data you care about.
5. **API design first:** an OpenAPI contract and a generated TypeScript client.
6. **Caching:** Redis cache-aside for the dashboard, with invalidation on transaction
   writes. Then HTTP caching with `ETag`.
7. **Background work:** a session cleanup job, then a transactional outbox that sends a
   "budget exceeded" notification with at-least-once delivery and an idempotent consumer.
8. **Performance:** load test with k6, profile with `pprof`, read `EXPLAIN ANALYZE`, tune
   the connection pool.
9. **Phase 2 AI:** LLM function calling and structured outputs to turn text into
   transactions through the same service, then expose the services as an MCP server.
   For a chaotic money routine, typing "lunch 120" beats filling a form.

## Level 3: a second, distributed project

Build a small system of 3 or 4 services (for example orders, payments, and notifications).
It covers the roadmap topics a single monolith cannot:

- microservices, SOA, and service boundaries
- gRPC between services, REST at the edge
- message brokers: Kafka or RabbitMQ, events, the outbox pattern, sagas
- failure handling: timeouts, retries with backoff, circuit breakers, backpressure, load
  shifting
- data at scale: replication, read replicas, sharding concepts, the CAP theorem
- real time: WebSockets or server-sent events
- Kubernetes with the Gateway API, Helm, GitOps (Argo CD), distributed tracing across
  services
- service mesh: understand what it solves, adopt only if needed

## Level 4: awareness only

Know what each one is, what problem it solves, and when to choose it. Go deep only when a
job requires it.

- GraphQL, SOAP, SAML, OpenID Connect details
- NoSQL families: learn Redis deeply, one document database (MongoDB or DynamoDB), and know
  wide-column (Cassandra, ScyllaDB), graph (Neo4j), time-series (TimescaleDB, InfluxDB),
  and analytics (ClickHouse) by their use cases
- search engines (Elasticsearch, Solr): try PostgreSQL full-text search first
- web servers (Nginx, Apache, Caddy, IIS) and LXC
- serverless

## Engineering practices to master

roadmap.sh lists tools and topics. Senior engineers are also judged on these practices,
which work in any language or stack. Ideas taken from
[Matt Pocock's skills](https://github.com/mattpocock/skills) and the books behind them.

| Practice | What it means | Where you practice it | Go deeper |
| --- | --- | --- | --- |
| Test-driven development | Write a failing test, make it pass, repeat in small slices. Test behavior through public interfaces, and fake only at system boundaries. | Every daily exercise is red to green | Kent Beck, _Test-Driven Development by Example_ |
| Disciplined debugging | Build a command that fails on the bug first, then shrink it, list testable guesses, and add a regression test before the fix. | `.claude/skills/diagnosing-bugs` in LedgerFlow | Andreas Zeller, _Why Programs Fail_ |
| Code review | Review two things separately: does it match the spec, and does it follow the standards? Name code smells as judgement calls. | Daily review step in the 10-day plan | Martin Fowler, _Refactoring_ (chapter 3, code smells) |
| Module design | Deep modules: a small interface that hides a lot of behavior. Put seams where tests and changes need them. | LedgerFlow services, then the Level 3 service boundaries | John Ousterhout, _A Philosophy of Software Design_ |
| Domain language | One shared set of terms in docs, code, and talk. | BRD §5 Key Terms is LedgerFlow's glossary | Eric Evans, _Domain-Driven Design_ |

Skills from that repo worth adding later, when the matching level starts:

- `wizard` (Level 2): a guided terminal script for steps only a human can do, like cloud
  account setup and CI secrets.
- `to-tickets` and `wayfinder` (Level 3): break a large multi-service project into
  ordered, blocking tickets.
- `teach` (Levels 2 to 4): a multi-session study workspace with short lessons and recall
  practice. It writes HTML lessons, so it costs more tokens. Use it for topics that
  LedgerFlow does not cover.

## How to go deep on any topic

1. Learn the problem it solves, in one sentence.
2. Build it in your current project. No toy demos.
3. Break it on purpose: kill the database, send a request twice, run two writes at once.
4. Measure it: a test, a trace, a load test, or a query plan.
5. Write 3 lines on what you learned in that day's lesson file.

## Sources

- [roadmap.sh backend roadmap](https://roadmap.sh/backend) and its
  [content repository](https://github.com/nilbuild/developer-roadmap)
- [Go 1.27 release blog](https://go.dev/blog/go1.27)
- [PostgreSQL 19 RC 1 announcement](https://www.postgresql.org/about/news/postgresql-19-rc-1-released-3409/)
- [OWASP Top 10:2025](https://owasp.org/Top10/2025/)
- [OpenTelemetry status](https://opentelemetry.io/status/)
- [Ingress NGINX retirement](https://kubernetes.io/blog/2025/11/11/ingress-nginx-retirement/)
- [RFC 9457 Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)
- [Idempotency-Key header draft](https://datatracker.ietf.org/doc/draft-ietf-httpapi-idempotency-key-header/)
- [DBOS July 2026 update](https://dbos.dev/blog/new-in-dbos-july-2026) and
  [durable execution overview](https://www.kai-waehner.de/blog/2025/06/05/the-rise-of-the-durable-execution-engine-temporal-restate-in-an-event-driven-architecture-apache-kafka/)
- Modular monolith trend (opinion pieces, not surveys):
  [DEV Community](https://dev.to/codewithamrendra/microservices-vs-modular-monolith-in-2026-which-architecture-actually-scales-1iie),
  [Java Code Geeks](https://www.javacodegeeks.com/2025/12/microservices-vs-monoliths-in-2026-when-each-architecture-wins.html)
