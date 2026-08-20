# Agent instructions

**Read [`CLAUDE.md`](CLAUDE.md) in this folder first. It is the canonical entry point for
every AI agent working in this repository, whatever tool you are.**

This file exists only so that agents which look for `AGENTS.md` (Codex and others) find
their way there. It is deliberately a pointer and not a copy. Two files holding the same
rules drift apart, and then an agent follows the stale one.

In short, and only in short, because `CLAUDE.md` is the authority:

1. **Inspect before proposing.** Treat code, Git state, and command output as truth.
   Documentation may be stale. Never assume planned architecture is already built.
2. **One active objective.** Evaluate every detour against the scope test in `CLAUDE.md`.
   If it does not block the active objective, name it as scope drift and park it.
3. **Discovery-first coaching.** The owner is learning Go and backend engineering. Give
   the goal, the files, precise requirements, hints, and acceptance criteria. Do not lead
   with exact commands or complete code.
4. **Session continuity.** `docs/SESSION_STATE.md` is the living handoff. Compare it
   against the repository at session start and correct stale claims.

Two rules worth stating here because they are the easiest to violate by accident:

- **Do not claim success without running the relevant checks.** Separate verified facts
  from assumptions, and show the evidence.
- **Keep business rules out of Gin handlers.** Handlers parse and format; services own
  business behavior. Dependencies point inward: handler to service to store/domain.
