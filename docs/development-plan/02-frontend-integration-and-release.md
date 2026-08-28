# Frontend Integration and Release Plan

> Scope: F1 through F7.
>
> Hard prerequisite: G1 in the backend and infrastructure plan must pass.

The frontend is an API consumer. It must not own or reimplement balance, budget, archive,
authorization, or dashboard business rules. It displays server results and submits user intent.

## F1. Create the frontend foundation and API client

**Goal:** A strict, tested React application can call the frozen Phase 1 API.

**Prerequisite:** G1 evidence is recorded.

**Build:**

1. Create the Vite React and TypeScript application under `frontend/`.
2. Enable strict TypeScript checks.
3. Add React Router, Mantine, TanStack Query, Mantine Form, and the chosen test tools.
4. Use a feature-first structure for `auth`, `accounts`, `categories`, `transactions`,
   `budgets`, and `dashboard`.
5. Generate the TypeScript API client from `backend/api/openapi.yaml`.
6. Keep generated client files read-only and reproducible.
7. Configure the API base URL through environment variables.
8. Create the app shell, protected-route boundary, query client, shared error handling,
   notifications, and loading and empty-state patterns.
9. Add one CI job for type checking, tests, build, and generated-client consistency.

Do not add Redux, Next.js, Tailwind, a second component system, or handwritten copies of
generated API types.

**Verify:**

- Type checking and production build pass.
- The client regenerates without unexplained changes.
- A small integration test calls liveness and handles the shared API error envelope.
- Routes render useful loading, empty, success, and error states.
- No secret is present in the frontend bundle.

**Exit gate:** The frontend foundation is stable and contains no feature business logic.

## F2. Implement authentication user flows

**Goal:** A user can register, sign in, remain signed in, and sign out in a browser.

**Prerequisite:** F1.

**Build:**

1. Add register and login routes with accessible labels and field-level errors.
2. Send cookie-authenticated requests with the correct browser credentials policy.
3. Resolve the current session from `/api/auth/me` during app startup.
4. Protect application routes and return signed-out users to login.
5. Implement logout and clear all user-specific query caches.
6. Map documented API error codes to short, safe user messages.
7. Avoid storing session tokens in local storage, session storage, or JavaScript state.

**Verify:**

- Registration, login, reload, expired session, and logout work in a real browser.
- Protected pages never flash private data before auth is known.
- Wrong credentials do not reveal whether an email exists.
- Keyboard-only use, focus order, labels, and error announcements work.
- Auth component and browser tests pass.

**Exit gate:** AUTH-01 works from the browser and no credential material is exposed to client
JavaScript.

## F3. Implement accounts and categories

**Goal:** First-time setup and archive workflows work before transaction entry is exposed.

**Prerequisite:** F2.

**Build:**

1. Add account list, create, edit, and archive flows.
2. Display starting and current balances distinctly.
3. Show account type and currency rules before submission.
4. Add category lists grouped by income and expense, plus create, edit, and archive flows.
5. Use color swatches and the established icon set for presentation choices.
6. Hide archived records from new-entry choices while keeping history views readable.
7. Invalidate only the affected account, category, and dashboard queries after mutations.
8. Display server business-rule errors without duplicating the rule in the component.

**Verify:**

- First-time account and category setup works with empty data.
- Decimal values render exactly and never pass through JavaScript arithmetic for business
  calculations.
- Archive and immutable-field errors are clear.
- Loading, empty, error, and success states work on desktop and mobile widths.
- Two browser sessions for different users never show shared cached data.

**Exit gate:** The user can create valid inputs for transactions and understand their active
and archived state.

## F4. Implement the transaction workflow

**Goal:** Daily entry and correction are fast, clear, and synchronized with server results.

**Prerequisite:** F3.

**Build:**

1. Add a transaction list with date, account, category, and type filters.
2. Add create, edit, and delete flows for income and expenses.
3. Filter category choices by transaction type and show only active records.
4. Accept amounts as decimal text, include an optional note, and submit the exact values to
   the API.
5. Confirm destructive deletion.
6. After a successful mutation, invalidate affected transactions, accounts, budgets, and
   dashboard queries.
7. Preserve entered form values after recoverable validation or network errors.
8. Show server-confirmed results. Do not optimistically calculate balances in the browser.

**Verify:**

- Create income and expense updates the visible account balance after refetch.
- Edit amount, type, account, category, and date updates every affected view.
- Delete restores the correct balance and derived values.
- Filters can be combined and cleared.
- Invalid, archived, mismatched, and other-user references are handled safely.
- Browser tests cover normal entry and correction of a mistake.

**Exit gate:** The BRD daily activity and correct-a-mistake flows work from the browser without
manual reload.

## F5. Implement budgets and dashboard

**Goal:** The selected-month review answers the user's money questions within one view.

**Prerequisite:** F4.

**Build budgets:**

1. Add month selection and monthly budget list.
2. Add create, edit, and delete flows limited to active expense categories.
3. Display amount, spent, remaining, percentage used, and status from API values.
4. Make safe, caution, high-risk, and exceeded states understandable without relying only on
   color.
5. Use the wording `budget remaining`.

**Build dashboard:**

1. Open on the current month and allow another month to be selected.
2. Show active account balances and the allowed total or currency grouping.
3. Show income, expenses, net savings, and average daily spending.
4. Show spending by expense category, budget status, and recent transactions.
5. Use a chart only where it improves comparison and keep the numeric values available as
   text.
6. Keep the layout dense, calm, and easy to scan. The financial state, not decoration, is the
   first viewport signal.
7. Refetch affected queries after transaction and budget mutations.

**Verify:**

- Budget thresholds and negative remaining display correctly.
- Editing or deleting a related transaction moves budget status both up and down.
- Dashboard values match direct API responses for empty, current, and past months.
- Archived-history and currency rules display correctly.
- The longest labels and money values fit on supported mobile and desktop widths.
- A short usability check confirms the owner can identify balance, monthly expenses, and
  budget status within 10 seconds.

**Exit gate:** BUD-01, BUD-02, DASH-01, and DASH-02 work in the browser.

## F6. Stabilize the complete application

**Goal:** Prove the whole Phase 1 flow across the browser, API, and database.

**Prerequisite:** F5.

**Automated browser scenarios:**

1. Register, sign in, reload, and sign out.
2. Create an account and categories.
3. Create income and expense transactions.
4. Edit amount, type, account, category, and date.
5. Delete a transaction.
6. Create, edit, exceed, recover, and delete a budget.
7. Review current and past dashboard months.
8. Archive accounts and categories and verify history remains.
9. Use two users and prove data isolation.
10. Simulate API validation, auth expiry, network loss, and server failure.

**Quality checks:**

- accessibility checks plus keyboard-only review;
- responsive screenshots at representative mobile and desktop widths;
- no overlapping, clipped, or unreadable text;
- browser console has no unexpected errors or warnings;
- no `any` in domain-facing frontend code;
- loading and mutation states do not cause duplicate submissions;
- direct navigation and refresh work on every route;
- frontend build output contains no private configuration;
- full CI is green.

**Exit gate:** The application behaves correctly under normal use, expected invalid input, and
recoverable failures.

## F7. Complete Phase 1 acceptance and release

**Goal:** Close Phase 1 with evidence that every stated requirement works.

**Prerequisite:** F6.

**Acceptance:**

1. Run every scenario in the
   [requirements traceability document](03-requirements-traceability.md).
2. Record the test or observation that proves each BRD requirement.
3. Resolve every failed requirement or explicitly keep Phase 1 open.
4. Review all open product decisions and known limitations.
5. Run a clean database migration and a backup and restore drill.
6. Build the backend image and frontend production bundle from a clean checkout.
7. Run the full application from release artifacts and complete a smoke test.
8. Review relevant OWASP ASVS Level 1 controls before any public deployment.
9. Update the README, setup guide, API contract, database design, and session state.
10. Tag the verified local MVP as `v0.1.0` only when all Phase 1 acceptance criteria pass.

**Public deployment is optional:**

Phase 1 can finish as a verified local application. If a public URL is needed, add a separate
deployment milestone with HTTPS, secure cookies, narrow origins, managed secrets, database
backups, migration procedure, rollback procedure, logging, and uptime checks. Verify provider
terms at that time instead of relying on old free-tier assumptions.

**Exit gate:** Every BRD Phase 1 requirement has passing evidence, release artifacts are
reproducible, and no known critical correctness or security issue remains.
