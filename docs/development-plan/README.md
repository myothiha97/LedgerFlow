# LedgerFlow Phase 1 Development Plan

> Status: Active
>
> Baseline date: 2026-08-28
>
> Current milestone: B1, finish the `sqlc` foundation for users and sessions.

## 1. Purpose

This plan turns the business requirements and technical specification into an ordered,
testable path for building LedgerFlow Phase 1.

The delivery model is intentionally backend-first:

```text
requirements and decisions
        -> backend architecture and database
        -> backend features and automated tests
        -> local infrastructure and CI
        -> backend and infrastructure readiness gate
        -> frontend implementation and integration
        -> Phase 1 acceptance and release
```

Frontend development must not begin before the readiness gate in Section 8 passes. This
keeps the frontend from defining unfinished business rules or depending on unstable APIs.

## 2. Source of truth

When documents disagree, use this order:

| Priority | Source | What it controls |
| --- | --- | --- |
| 1 | Running code, migrations, tests, and command output | What exists now and what is verified |
| 2 | [Business requirements v3](../specifications/business-requirements-v3.md) | Product scope, user outcomes, and business rules |
| 3 | [Technical specification v1](../specifications/technical-specification-v1.md) | Stack, API style, repository shape, and delivery constraints |
| 4 | [Database design v1](../specifications/database-design-v1.md) | Tables, constraints, indexes, and database transactions |
| 5 | [Architecture guidelines v1](../specifications/architecture-guidelines-v1.md) | Engineering conventions and quality rules |
| 6 | This development plan | Work order, gates, and required evidence |
| 7 | Session notes and older guides | Historical context only |

The BRD uses the user-facing term `starting balance`. The database design uses the column
name `initial_balance`. They mean the same value. Use `starting balance` in the UI and
API documentation, and follow the database design for SQL names.

## 3. Phase 1 boundary

Phase 1 includes only:

- registration, login, logout, and session restoration;
- accounts;
- income and expense categories;
- income and expense transactions;
- monthly category budgets;
- the dashboard;
- the correct transaction, balance, budget, and dashboard lifecycle;
- local Docker infrastructure and continuous integration.

Phase 1 excludes transfers, recurring transactions, goals, imports, currency conversion,
mobile or PWA work, AI, subscriptions, microservices, Kubernetes, and speculative scale
infrastructure.

## 4. Verified project baseline

This table records the repository state found on 2026-08-28.

| Area | Current evidence | State |
| --- | --- | --- |
| Git | `main` matches `origin/main`; unrelated local edits exist in `example.go` | Preserve the local edit |
| Go module | `backend/go.mod`, Go 1.26.4, Gin 1.12.0 | Done |
| API process | `backend/cmd/server/main.go` serves `GET /ping`; port defaults to 4000 | Minimal foundation |
| Go checks | Test, vet, and build pass | Verified, but there are no test files |
| Local database config | PostgreSQL 16 Compose service uses `ledgerflow` credentials and port 5432 | Configured |
| Database runtime | Docker Desktop is currently stopped | Live database state not reverified |
| Migrations | `0001_init` defines `users` and `sessions`; apply, rollback, and reapply were verified on 2026-08-22 | Done based on recorded evidence |
| Database generation | `sqlc` 1.31.1 is installed; no `sqlc.yaml`, query directory, or generated package exists | Current work |
| Backend layers | `domain`, `service`, and `store` packages do not exist; `handler` is empty | Not started |
| API contract | `backend/api/openapi.yaml` does not exist | Not started |
| Backend image | No backend `Dockerfile` or `.dockerignore` exists | Not started |
| CI | No workflow exists | Not started |
| Frontend | The `frontend/` directory is empty | Correctly deferred |

The June 2026 session notes describe a larger generated scaffold. That scaffold is not in
the current tree because the backend was later reset for an incremental rebuild. Do not
treat those removed files as completed work.

### Session history reconciliation

| Session | How this plan treats it |
| --- | --- |
| 2026-06-10 backend setup | Historical architecture and learning intent; several referenced files no longer exist |
| 2026-06-10 backend scaffold | Superseded by the later rebuild; do not mark its removed packages as complete |
| 2026-08-13 backend foundation | Current module layout, configurable port, and minimal Gin server are present and verified |
| 2026-08-14 PostgreSQL ports | Temporary learning configuration; superseded by the next session |
| 2026-08-15 PostgreSQL configuration | Current files use PostgreSQL port 5432 and `ledgerflow` credentials |
| 2026-08-22 first migration | Migration files are present and recorded as applied, rolled back, and reapplied successfully |

The latest recorded next outcome remains `backend/tasks/2.md`: add SQL queries and `sqlc`
configuration, generate the store package, and restore a green build.

## 5. Delivery method

LedgerFlow uses a lightweight, stage-gated software development lifecycle. It keeps the
useful controls from larger teams without adding process that does not help a solo project.

| Lifecycle activity | How LedgerFlow applies it | Required evidence |
| --- | --- | --- |
| Requirements | Use BRD requirement IDs and settle open business choices before affected code | Decision recorded and acceptance behavior stated |
| Design | Update the database design and OpenAPI contract before implementation changes their shape | Reviewed migration or API contract |
| Implementation | Build one small vertical slice from domain to HTTP | Focused diff with clear layer boundaries |
| Verification | Test business rules, database behavior, HTTP behavior, and full user flows at the correct level | Passing automated checks and recorded manual evidence where needed |
| Release | Build from a clean state, apply migrations safely, and verify the packaged application | Reproducible build, smoke test, and release checklist |
| Operations | Keep configuration external, protect secrets, provide health checks, and prove backup restoration | Runbook and successful restore drill before public use |

Security work is part of every stage, not a final cleanup. This follows the intent of the
[NIST Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final).
Web security acceptance should use the relevant Level 1 controls from the
[OWASP Application Security Verification Standard](https://owasp.org/www-project-application-security-verification-standard/).

## 6. Work cycle for every milestone

Use this loop for B1 through F7:

1. Confirm the requirement IDs and prerequisite decisions.
2. Create or reuse one small task in `backend/tasks/<n>.md` while backend work is active.
3. Change the contract or migration first when the external or stored data shape changes.
4. Implement from inner layers to outer layers: domain, store, service, handler, wiring.
5. Add the smallest meaningful tests, then add integration coverage for database boundaries.
6. Run the milestone checks and inspect the result.
7. Review the diff for scope, secrets, generated files, and accidental changes.
8. Update the plan status and session state only after the exit gate passes.
9. Make one focused commit while the branch is green.

Only one milestone is active at a time. New ideas that do not block its exit gate go to a
parking lot instead of interrupting the work.

## 7. Ready and done rules

### Definition of ready

A milestone is ready when:

- its business behavior is clear;
- all required earlier milestones are complete;
- open decisions that affect it are resolved;
- the expected files and layer ownership are known;
- acceptance checks can be stated before coding.

### Definition of done

A milestone is done when:

- the stated behavior works;
- business logic is in services, not handlers;
- every user-owned query is scoped by `user_id`;
- money never uses a floating-point type;
- automated tests cover success, invalid input, ownership, and failure rollback where relevant;
- migrations apply, roll back, and apply again when the schema changed;
- generated files can be reproduced and were not edited by hand;
- formatting, test, vet, build, and relevant integration checks pass;
- documentation and requirement traceability are current;
- no secrets, debug artifacts, or unrelated files entered the change.

## 8. Ordered roadmap and gates

| Order | Milestone | Result | Status |
| --- | --- | --- | --- |
| B0 | Baseline and decision register | Current state and decision deadlines are explicit | Done by this plan |
| B1 | `sqlc` foundation | Type-safe user and session queries generate and compile | Current |
| B2 | Runtime architecture foundation | Config, DB pool, responses, router, health, and shutdown follow the target layers | Planned |
| B3 | API contract foundation | OpenAPI becomes the versioned HTTP contract | Planned |
| B4 | Auth domain, store, and service | Auth rules work without HTTP and have unit tests | Planned |
| B5 | Auth HTTP and middleware | Register, login, logout, and `me` work end to end | Planned |
| B6 | Accounts vertical slice | Accounts work end to end with correct balances and archive rules | Planned |
| B7 | Categories vertical slice | Categories work end to end with type and archive rules | Planned |
| B8 | Transaction schema and atomic store | Database operations can change transactions and balances atomically | Planned |
| B9 | Transaction service and API | Create, edit, move, delete, list, and filters preserve all invariants | Planned |
| B10 | Budgets vertical slice | Monthly budgets and derived status work end to end | Planned |
| B11 | Dashboard vertical slice | Selected-month summaries match source transactions | Planned |
| B12 | Backend hardening | Cross-user isolation, security, errors, and failure behavior are verified | Planned |
| I1 | Backend container and local operations | The API and database run reproducibly through Docker | Planned |
| I2 | Continuous integration | Every change is automatically checked from a clean environment | Planned |
| G1 | Backend and infrastructure readiness gate | Backend contract, architecture, tests, and local infrastructure are stable | Blocked by B1 through I2 |
| F1 | Frontend foundation and generated client | React can safely consume the frozen Phase 1 API | Blocked by G1 |
| F2 | Auth user interface | Session flows work in the browser | Blocked by F1 |
| F3 | Accounts and categories user interface | Setup and archive flows work in the browser | Blocked by F2 |
| F4 | Transactions user interface | Daily entry, editing, deleting, and filtering work in the browser | Blocked by F3 |
| F5 | Budgets and dashboard user interface | Monthly review works and updates after changes | Blocked by F4 |
| F6 | Full-stack stabilization | End-to-end, accessibility, and responsive checks pass | Blocked by F5 |
| F7 | Phase 1 acceptance and release | Every Phase 1 requirement has evidence | Blocked by F6 |

Detailed work is in:

- [Backend and infrastructure plan](01-backend-and-infrastructure.md)
- [Frontend integration and release plan](02-frontend-integration-and-release.md)
- [Requirements traceability and verification](03-requirements-traceability.md)

## 9. Frontend start gate

G1 passes only when every item below is true:

- [ ] All Phase 1 backend endpoints are implemented and documented in OpenAPI.
- [ ] Auth, ownership, accounts, categories, transactions, budgets, and dashboard tests pass.
- [ ] Transaction edit and delete behavior is proven against PostgreSQL, including rollback.
- [ ] A reconciliation test proves stored account balances match transaction-derived balances.
- [ ] Migrations build a clean database and can be rolled back and reapplied in development.
- [ ] The backend container runs as a non-root user and passes health checks.
- [ ] Compose starts the complete backend environment from documented configuration.
- [ ] CI runs formatting, generation, unit, integration, build, and image checks.
- [ ] No unresolved product decision can change an API response or frontend workflow.
- [ ] The traceability matrix has planned or passing evidence for every BRD requirement.

Passing G1 authorizes creating frontend source files. It does not require a public cloud
deployment. Local Docker infrastructure is enough for Phase 1.

## 10. Decision register

`Recommended` means the plan's default, not an accepted product decision. Record the final
choice in the BRD or a small decision log before its deadline.

| ID | Decision | Recommended Phase 1 choice | Decide before |
| --- | --- | --- | --- |
| D1 | Credit cards | Exclude credit cards until the debt balance meaning is specified | B6 accounts |
| D2 | Currency | Use one currency for Phase 1; never add unlike currencies | B6 accounts |
| D3 | Starter categories | Create a small default set for each new user | B7 categories |
| D4 | Archiving a non-zero account | Block archive until its balance is zero | B6 accounts |
| D5 | Budget wording | Use `budget remaining`, not `safe to spend` | Locked by BRD v3 |
| D6 | Future-dated transactions | Reject them in Phase 1 | B9 transactions |
| D7 | API generation | Use OpenAPI, generated Go contract types, and a generated TypeScript client | B3 API contract |
| D8 | Authentication | Use server-side sessions with secure cookie handling | Locked by project history |
| D9 | Calendar timezone | Use one configured application timezone for Phase 1 | B11 dashboard |
| D10 | First public host | Defer selection until a public URL is required | After F7 |
| D11 | Decimal display and rounding | Keep exact four-decimal values in storage and APIs; define display and percentage rounding explicitly | B10 budgets |

## 11. Main risks

| Risk | Prevention | Proof |
| --- | --- | --- |
| Balance drift | Reverse old effect, apply new effect, use one database transaction | Unit matrix, PostgreSQL integration tests, reconciliation query |
| Cross-user data access | Scope every query by `user_id`; never authorize by record ID alone | Two-user tests for every feature |
| Decimal loss | PostgreSQL `NUMERIC(19,4)`, Go decimal, JSON decimal strings | Contract tests and round-trip tests |
| Partial writes | Use explicit database transactions and rollback on every error | Forced-failure integration tests |
| API and frontend drift | OpenAPI is the contract; generated code is reproducible | Generation check in CI |
| Schema drift | Versioned migrations and clean-database CI | Apply, inspect, rollback, reapply |
| Frontend rework | Do not start frontend before G1 | Gate checklist |
| Lost local data | Document and test backup and restore before relying on the data | Restore drill |
| Scope expansion | Enforce the Phase 1 boundary and one active milestone | Review checklist and parking lot |

## 12. Change control

- A business-rule change updates the BRD, acceptance examples, and affected tests first.
- A stored-data change adds a new migration and updates the database design. Do not edit a
  migration that has been shared or applied to non-disposable data.
- An HTTP contract change updates OpenAPI before implementation and regenerates consumers.
- An architecture change must identify a real blocker and update the technical specification.
- A newly requested feature is not part of Phase 1 until scope, cost, and displaced work are
  recorded.

This keeps the plan useful as the repository changes without turning it into a second source
of business truth.

## 13. Beginner glossary

| Term | Plain meaning |
| --- | --- |
| Acceptance criteria | Observable facts that prove a requirement works |
| API contract | The agreed request and response shapes used by backend and frontend |
| CI | Automated checks that run from a clean environment after a code change |
| Database migration | A numbered, repeatable change to the database structure |
| Exit gate | Checks that must pass before the next milestone starts |
| Generated code | Code produced from SQL or OpenAPI and never edited by hand |
| Integration test | A test that checks two real parts together, such as Go and PostgreSQL |
| Reconciliation | Recomputing balances from source transactions and comparing the result |
| Rollback | Undoing a failed database operation or migration safely |
| Unit test | A focused test of one rule without external systems |
| Vertical slice | One feature built through every needed layer, from storage to HTTP |
