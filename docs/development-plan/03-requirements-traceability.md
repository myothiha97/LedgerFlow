# Requirements Traceability and Verification

> Purpose: prove that each BRD requirement is designed, implemented, and tested.
>
> Current state: no Phase 1 requirement is complete yet. Migration `0001` is foundation work,
> not proof of an auth requirement.

## 1. How to maintain this document

Use these status values:

- `Planned`: the milestone and expected proof are known.
- `Implemented`: code exists, but all required proof is not complete.
- `Verified`: the listed automated and acceptance evidence passes.
- `Blocked`: an unresolved decision or defect prevents completion.

Change a row to `Verified` only after linking the exact test, CI run, or recorded acceptance
result. A passing build alone does not verify user behavior.

## 2. Requirement matrix

| Requirement | Planned milestone | Required automated evidence | Acceptance evidence | Status |
| --- | --- | --- | --- | --- |
| AUTH-01 | B4, B5, F2 | Auth service tests, session store integration tests, auth HTTP tests, browser auth flow | Register, sign in, reload with a valid session, sign out, and lose access | Planned |
| AUTH-02 | B5 through B12, F6 | Two-user ownership tests for every protected store, service, and endpoint | User A cannot read or change any User B record, including guessed UUIDs | Planned |
| ACC-01 | B6, F3 | Account validation, decimal round-trip, create integration and HTTP tests | Create an account with name, allowed type, starting balance, and currency | Planned |
| ACC-02 | B6, F3 | Account read and edit tests, starting-balance recalculation tests, immutable-field tests | View current balance and edit only fields allowed by the rules | Planned |
| ACC-03 | B6, F3, F6 | Archive policy, ownership, active-list, and history tests | Archive an allowed account; hide it from new choices without deleting history | Planned |
| CAT-01 | B7, F3 | Type, normalized uniqueness, edit, ownership, integration, and HTTP tests | Create and edit one income and one expense category | Planned |
| CAT-02 | B7, F3, F6 | Archive, active-list, immutable-type, and history tests | Archive a category; hide it from new choices without changing old records | Planned |
| TXN-01 | B8, B9, F4 | Validation, decimal, ownership, atomic create, balance reconciliation, and HTTP tests | Create valid income and expense and observe correct balances | Planned |
| TXN-02 | B9, F4 | Combined date, account, category, type, ownership, and result-bound tests | Filter a populated list and clear filters | Planned |
| TXN-03 | B8, B9, F4, F6 | Full reverse-then-apply matrix, PostgreSQL rollback, reconciliation, and HTTP tests | Edit every supported field, move accounts, change type and date, then delete | Planned |
| BUD-01 | B10, F5 | Expense-category validation, month normalization, uniqueness, CRUD, ownership, and HTTP tests | Create, edit, and delete one budget for an expense category and month | Planned |
| BUD-02 | B10, F5 | Status boundary, negative remaining, downward movement, and overall remaining tests | See amount, spent, remaining, percentage, and correct status | Planned |
| DASH-01 | B11, F5 | Aggregation, month boundary, archive, currency, timezone, ownership, and contract tests | Read account and selected-month activity in one view | Planned |
| DASH-02 | B11, F4, F5 | Mutation-to-summary integration tests and frontend query invalidation tests | Change a transaction or budget and see the dashboard refresh without reload | Planned |
| DATA-01 | B8, B9, B10, B12 | Forced-failure PostgreSQL tests for every multi-step write plus reconciliation | Trigger a rejected change and confirm records, balances, and budgets are unchanged | Planned |

## 3. Cross-cutting technical controls

| Control | Design rule | Primary proof |
| --- | --- | --- |
| Layer ownership | Handlers parse and format; services own business rules; stores own PostgreSQL access | Package review plus service tests without Gin |
| Money precision | `NUMERIC(19,4)`, Go decimal, JSON decimal strings, no floating point | Static review plus SQL-Go-JSON round-trip tests |
| Authorization | Every user-owned query includes `user_id` | Two-user tests and query review |
| Atomicity | Transaction record and all balance effects commit or roll back together | Forced-failure PostgreSQL integration tests |
| Balance integrity | Current balance equals starting balance plus income minus expense | Reconciliation after every lifecycle case |
| Archive integrity | Archive hides future choices but preserves history | Active-list and historical-report tests |
| Contract integrity | OpenAPI changes first and generated code is reproducible | Contract validation and clean generation in CI |
| Migration integrity | Versioned migrations apply, roll back, and apply again | Disposable PostgreSQL CI job |
| Session security | Password hashes, hashed random tokens, safe cookies, expiry, and revocation | Auth unit, integration, HTTP, and browser tests |
| Secret handling | Configuration comes from environment and secrets never enter logs or images | Git review, image inspection, and log tests |
| Recoverability | PostgreSQL data can be backed up and restored | Restore drill plus reconciliation |
| Reproducibility | A clean checkout builds backend, frontend, generated code, and containers | CI and release smoke test |

## 4. Test levels

Use the lowest test level that proves the behavior, then add a higher-level test only for a
boundary or critical user flow.

| Test level | Main purpose | LedgerFlow focus |
| --- | --- | --- |
| Domain or pure function | Prove deterministic calculations | Balance effects, budget thresholds, date and month helpers |
| Service unit | Prove business decisions and operation sequence | Auth rules, ownership, archive rules, reverse-then-apply, dashboard assembly |
| Store integration | Prove SQL, constraints, locks, transactions, and mappings | CRUD, atomic transaction lifecycle, indexes, reconciliation |
| Handler or contract | Prove HTTP parsing, auth, status, and JSON shape | Every endpoint and shared error envelope |
| Browser component | Prove rendering and user interaction in isolation | Forms, tables, statuses, loading and error states |
| End to end | Prove a small number of complete user journeys | First-time setup, daily entry, correction, budget review, two-user isolation |
| Manual acceptance | Prove human understanding and visual quality | Dashboard comprehension, keyboard use, responsive layout |

Critical service rules need many focused unit cases. Browser tests should cover important
journeys, not duplicate every lower-level edge case.

## 5. Canonical acceptance data

Use deterministic data so results are easy to inspect.

| Record | Suggested value |
| --- | --- |
| User A | `alice@example.test` |
| User B | `bob@example.test` |
| Account A1 | Cash, starting balance 1000.0000 |
| Account A2 | Bank, starting balance 500.0000 |
| Categories | Salary income; Food and Rent expenses |
| Food budget | 300.0000 for the selected month |

Use a fixed past month for month calculations and a fixed application timezone. This avoids
tests changing with today's date.

## 6. Critical acceptance scenarios

### A. Authentication and isolation

1. Register User A and User B.
2. Sign in as User A and create one record of every user-owned type.
3. Sign in as User B and request each User A UUID through read and mutation endpoints.
4. Expect no data disclosure and no changed User A record.
5. Sign out User A and confirm the old session no longer works.

### B. Transaction and balance lifecycle

Starting state: Cash 1000.0000 and Bank 500.0000.

| Action | Expected Cash | Expected Bank |
| --- | ---: | ---: |
| Create Cash expense 200 | 800.0000 | 500.0000 |
| Edit expense to 150 | 850.0000 | 500.0000 |
| Move expense to Bank | 1000.0000 | 350.0000 |
| Change Bank transaction from expense to income | 1000.0000 | 650.0000 |
| Delete transaction | 1000.0000 | 500.0000 |

Run reconciliation after every row. Also move a transaction between months and categories and
verify both the old and new reporting groups.

### C. Atomic failure

1. Start an edit that would reverse an old transaction and apply a new one.
2. Force a database failure after the reverse step but before completion.
3. Expect the old transaction and all old balances to remain unchanged.
4. Run reconciliation and expect no difference.

### D. Budget boundaries and reversal

For a Food budget of 300.0000:

| Spent | Used | Expected status |
| ---: | ---: | --- |
| 0 | 0% | Safe |
| 210 | 70% | Safe |
| Above 210, through 270 | Above 70%, through 90% | Caution |
| Above 270, through 300 | Above 90%, through 100% | High risk |
| Above 300 | Above 100% | Exceeded |

After reaching `Exceeded`, edit or delete spending until the result becomes `Safe`. This proves
status moves down as well as up.

### E. Archive and history

1. Create transactions using an account and category.
2. Archive each record according to the chosen policy.
3. Confirm neither appears in new transaction or budget choices.
4. Confirm past transactions and selected-month totals still include their history.
5. Confirm archived accounts do not appear in the current active-account summary.

### F. Dashboard consistency

1. Create known income and expenses in one fixed month.
2. Compare each dashboard number with direct transaction and budget queries.
3. Edit amount, account, category, type, and date one at a time.
4. Delete the transaction.
5. After each action, confirm affected old and new months, categories, accounts, budgets, and
   recent lists agree.

### G. Browser completion flow

1. Register and sign in.
2. Create an account and categories.
3. Record income and expense.
4. Correct and delete a mistake.
5. Set and review a monthly budget.
6. Read the dashboard.
7. Sign out.

The user should identify current balance, monthly expenses, and budget status within 10 seconds
of opening a populated dashboard.

## 7. Phase 1 release evidence

Before F7 can pass, record:

- the commit and release tag;
- the CI run covering that exact commit;
- migration apply, rollback, and reapply result;
- automated test summary by level;
- transaction reconciliation result;
- backup and restore result;
- backend image and frontend build identifiers;
- manual dashboard comprehension result;
- unresolved limitations that do not violate Phase 1 requirements.

Missing evidence means the requirement remains `Implemented`, not `Verified`.
