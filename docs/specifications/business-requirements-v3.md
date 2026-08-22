# LedgerFlow Business Requirements (v3)

> **Status:** Draft for Phase 1 review
>
> **Purpose:** This document explains what LedgerFlow must do for the user. Technical
> design belongs in the technical specification and database design documents.

**In one sentence:** LedgerFlow helps a person record everyday money activity, keep
account balances accurate, and see whether monthly spending is within budget.

## 1. Product at a Glance

LedgerFlow gives the user a clear view of their money without the complexity of an
accounting system.

| The user wants to know | LedgerFlow shows |
| --- | --- |
| How much money do I have now? | The current balance of each active account and one total when the currencies match |
| How much did I earn and spend? | Income, expenses, and net savings for the selected month |
| Where did my money go? | Spending grouped by expense category |
| Am I within budget? | Amount spent, amount remaining, percentage used, and status for each budget |
| How much budget is left? | Remaining budget by category and across all budgeted categories |

The main product goal is simple: the user should understand their money situation within
10 seconds of opening the dashboard.

LedgerFlow provides tracking information. It does not provide financial advice or
guarantee that a purchase is affordable.

## 2. User Problem

Many personal finance tools are either too complex or too general. This makes daily
entry slow, so records become incomplete and balances cannot be trusted.

The user needs a small system that makes two jobs easy:

1. Record income and expenses in a few steps.
2. Turn those records into accurate balances, monthly totals, and budget status.

## 3. Target User

Phase 1 is for one person managing their own finances. The system may have many
registered users, but each user has a private and separate set of financial records.

The target user:

- tracks money manually;
- uses a few accounts and common income and expense categories;
- wants a quick daily view and a simple monthly review;
- does not need business accounting, shared household budgets, or investment tracking.

## 4. Phase 1 Scope

Phase 1 is complete when this core flow works from start to finish:

> Record money in or out, update the correct account balance, update the relevant
> monthly totals and budgets, then show the result on the dashboard.

### Included

| Area | User outcome |
| --- | --- |
| Account access | Register, sign in, sign out, and access only personal data |
| Accounts | Track where money is held and see the current balance |
| Categories | Organize income and expenses into meaningful groups |
| Transactions | Add, view, edit, delete, and filter income and expense records |
| Budgets | Set and monitor a monthly limit for each expense category |
| Dashboard | Understand balances, monthly activity, and budget status at a glance |

### Not Included

- bank synchronization;
- transfers between accounts;
- recurring transactions;
- savings goals and yearly planning;
- CSV import or export;
- receipt scanning;
- investments;
- shared or family budgets;
- multiple-currency conversion;
- a mobile app or installable web app (PWA);
- AI entry, insights, or recommendations.

## 5. Key Terms

| Term | Plain meaning |
| --- | --- |
| Account | A place where money is tracked, such as cash, a bank account, an e-wallet, or savings |
| Starting balance | The account balance when the user begins tracking it in LedgerFlow |
| Current balance | Starting balance plus recorded income, minus recorded expenses |
| Category | The reason for money coming in or going out, such as Salary, Food, Rent, or Transport |
| Transaction | One income or expense record linked to one account and one matching category |
| Budget | A spending limit for one expense category in one calendar month |
| Dashboard | A summary of current balances and activity for the selected month |
| Archive | Remove an account or category from future use while keeping its history |

Credit cards need a separate balance rule and are not confirmed as a Phase 1 account
type. See Section 13.

## 6. Core Business Rules

### 6.1 A Valid Transaction

Every transaction must:

- belong to the signed-in user;
- be either income or expense;
- have an amount greater than zero;
- use an active account owned by the user;
- use an active category owned by the user;
- use an income category for income or an expense category for an expense;
- have a transaction date;
- allow an optional note.

Phase 1 transactions are entered by the user. LedgerFlow marks their source as `manual`.

### 6.2 How Transactions Change Balances

For every account:

```text
current balance = starting balance + total income - total expenses
```

- Creating income adds the full amount to the selected account.
- Creating an expense subtracts the full amount from the selected account.
- Editing a transaction first removes the complete effect of the old version, then
  applies the complete effect of the new version.
- Deleting a transaction removes its complete effect from the account.
- Moving a transaction to another account reverses it on the old account and applies it
  to the new account.
- The transaction record and all related balances must change together. If any part
  fails, nothing changes.

Example:

```text
Starting balance                         1,000
Add an expense of 200                      800
Edit the expense from 200 to 150           850
Delete the expense                       1,000
```

The system must never update only the difference between old and new transaction values.
It must always reverse the old version and apply the new version.

### 6.3 What Happens After a Transaction Changes

After a transaction is created, edited, or deleted, LedgerFlow must update every result
that uses it:

- old and new account balances;
- income, expense, and net savings totals for the old and new months;
- spending totals for the old and new categories;
- old and new monthly budget status;
- dashboard charts and recent transactions.

The user must see the successful result without needing to sign out or reload the app.

### 6.4 Monthly Budget Rules

- A budget belongs to one active expense category and one calendar month.
- A category can have only one budget in the same month.
- The budget amount must be greater than zero.
- Income categories cannot have budgets.
- Editing a budget recalculates its status immediately.
- Deleting a budget removes only the limit. It does not delete transactions.

For each category budget:

```text
spent = total expenses in that category during the selected month
remaining = budget amount - spent
percentage used = spent / budget amount * 100
```

`remaining` may be negative when the user is over budget.

| Status | Percentage used | Meaning |
| --- | --- | --- |
| Safe | 70% or less | On track |
| Caution | More than 70%, up to 90% | Slow down |
| High risk | More than 90%, up to 100% | Near the limit |
| Exceeded | More than 100% | Over budget |

Budget status is calculated from current transactions. It must move up or down after a
transaction is added, edited, moved, or deleted.

Overall budget remaining is the sum of all budget amounts minus spending in those
budgeted categories. Spending in a category without a budget still counts as a monthly
expense, but it does not change overall budget remaining.

### 6.5 Dashboard Rules

The dashboard opens on the current month. The user can select another month to review
its activity.

It must show:

- the current balance of each active account;
- one total balance when all active accounts use the same currency;
- income, expenses, and net savings for the selected month;
- average daily spending for the selected month;
- spending grouped by expense category;
- each budget's amount, spent amount, remaining amount, percentage, and status;
- recent transactions for the selected month.

The monthly calculations are:

```text
monthly income = total income in the selected month
monthly expenses = total expenses in the selected month
net savings = monthly income - monthly expenses
```

For the current month, average daily spending is monthly expenses divided by the number
of calendar days elapsed, including today. For a completed month, it is divided by the
number of days in that month.

Past monthly totals include transactions linked to archived accounts and categories.
Archived accounts do not appear in the current account summary.

## 7. Information LedgerFlow Keeps

| Record | User-visible information |
| --- | --- |
| User | Name, email, and password used to sign in |
| Account | Name, type, starting balance, current balance, currency, and active or archived state |
| Category | Name, income or expense type, optional color and icon, and active or archived state |
| Transaction | Account, category, income or expense type, amount, date, optional note, and source |
| Budget | Expense category, month, and limit amount |

Ownership and relationships are part of the business requirement:

- every record belongs to one user;
- one transaction belongs to one account and one category;
- one budget belongs to one expense category and one month;
- one user must never read or change another user's records.

Passwords must never be stored or displayed as plain text.

## 8. Common User Flows

### First-Time Setup

1. The user registers and signs in.
2. The user creates at least one account with a starting balance and currency.
3. The user creates or reviews income and expense categories.
4. The dashboard is ready for the first transaction.

### Record Daily Activity

1. The user chooses income or expense.
2. The user enters an amount, account, matching category, date, and optional note.
3. The user saves the transaction.
4. LedgerFlow updates balances, monthly totals, budgets, and the dashboard.

### Correct a Mistake

1. The user finds a transaction using the list or filters.
2. The user edits or deletes it.
3. LedgerFlow reverses the old effect and applies the corrected result everywhere.

### Retire an Account or Category

1. The user archives the account or category.
2. It no longer appears in new transaction or budget choices.
3. Existing transactions and past reports keep their original information.

## 9. Main Screens

| Screen | Main purpose | Required content or actions |
| --- | --- | --- |
| Register and sign in | Secure account access | Register, sign in, sign out, and clear errors |
| Dashboard | Quick financial review | Account balances, monthly totals, category spending, budget status, and recent transactions |
| Accounts | Manage places where money is tracked | List, create, edit, archive, and view current balance |
| Categories | Manage income and expense labels | List, create, edit, and archive |
| Transactions | Record and review activity | List, add, edit, delete, and filter by date, account, category, and type |
| Budgets | Plan and monitor monthly limits | Select month, create, edit, delete, and view usage and status |

## 10. Phase 1 Requirement Checklist

| ID | Requirement |
| --- | --- |
| AUTH-01 | The user can register, sign in, sign out, and stay signed in while their login is valid. |
| AUTH-02 | The user can access only their own financial records. |
| ACC-01 | The user can create an account with a name, type, starting balance, and currency. |
| ACC-02 | The user can view an account's current balance and edit details allowed by the rules below. |
| ACC-03 | The user can archive an account without deleting its transaction history. |
| CAT-01 | The user can create and edit income and expense categories. |
| CAT-02 | The user can archive a category without changing past transactions. |
| TXN-01 | The user can create a valid income or expense transaction. |
| TXN-02 | The user can list and filter transactions by date, account, category, and type. |
| TXN-03 | The user can edit or delete a transaction and see the correct result everywhere. |
| BUD-01 | The user can create, edit, and delete one monthly budget per expense category. |
| BUD-02 | The user can see spent, remaining, percentage used, and status for every budget. |
| DASH-01 | The user can see the current account summary and selected-month activity in one view. |
| DASH-02 | The dashboard reflects every successful transaction and budget change. |
| DATA-01 | A failed change leaves transactions, balances, and budgets unchanged. |

Details that change how money is calculated need clear rules:

- Editing the starting balance must recalculate the current balance.
- Account type and currency cannot change after the account has transactions.
- Category type cannot change after the category has transactions or budgets.
- Account and category names and display details may still be edited.

## 11. Phase 1 Completion Criteria

Phase 1 is done only when all requirements above are implemented and verified. At a
minimum, verification must prove that:

1. A registered user cannot access another user's data.
2. Creating, editing, moving, changing the type of, and deleting a transaction always
   produces the correct account balances.
3. Budget usage and status move both up and down when related transactions change.
4. Archived accounts and categories disappear from new-entry choices but remain in
   history.
5. Dashboard figures match the underlying transactions for the selected month.
6. A user can open the dashboard and understand balances, monthly spending, and budget
   status within 10 seconds.

## 12. Future AI Support

AI is planned for Phase 2, not Phase 1. Future capabilities may include:

- turning everyday text or voice into proposed transactions for user confirmation;
- answering questions about spending and income;
- suggesting categories;
- warning that a budget may be exceeded.

Phase 1 must preserve the transaction note and source so future input methods can use
the same records. AI must follow the same validation, balance, and budget rules as manual
entry. It must never bypass them.

Phase 1 includes no AI features or AI setup.

## 13. Decisions Still Needed

These choices were unclear in v2. They must be settled before the related features are
built:

1. **Credit-card meaning:** Decide whether a credit-card balance means amount owed or
   available credit. The normal income-adds and expense-subtracts rule does not model
   debt correctly. Until this is decided, do not claim credit-card support.
2. **Phase 1 currency rule:** Either require all accounts to use one currency or group
   dashboard totals by currency. Never add unlike currencies into one total. Requiring
   one currency is the simpler Phase 1 choice.
3. **Starting categories:** Decide whether a new user starts with common categories or an
   empty list. Starter categories make first-time setup faster.
4. **Non-zero archived accounts:** Decide whether archiving is blocked until the balance
   reaches zero or whether the archived balance remains visible separately. Hiding a
   non-zero balance would make the user's total misleading.
5. **Budget language:** Use "budget remaining" until a true "safe to spend" calculation
   includes available cash, unbudgeted expenses, and reserved money.
6. **Future-dated transactions:** Decide whether Phase 1 allows them and, if so, when they
   change the current balance. Rejecting future dates is the simpler Phase 1 choice.

These decisions do not expand Phase 1. They make existing Phase 1 behavior precise.
