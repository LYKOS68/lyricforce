# workspace/

Shared root for project documents, current state and logs, used by both the user and the CLI (Claude Code).
This folder isn't shipped with the application.

```
workspace/
  docs/              project documents
    start-draft.md   product vision, tech decision, phonetic rhyme logic, roadmap (source of truth)
    decisions.md     decision log (what was decided, when, why)
    build-orders/    work orders as given, one file per order (YYYY-MM-DD-<topic>.md)
  state/
    STATE.md         current state: what works, active area, next steps, open questions
    Handoff.jsonl    handoff to the next CLI model: done, plan, problems, todos (one JSON object per line)
  logs/
    sessions/        work log per day (YYYY-MM-DD.md): what was done, verified, left open
    build-*.log      build output written by scripts/build.ps1 (not versioned)
```

## Rules

- **Before work:** read `state/Handoff.jsonl`, `state/STATE.md` and the latest session log.
- **After work:** update `state/STATE.md` and append to `logs/sessions/<today>.md`.
- **New work order:** save it unchanged to `docs/build-orders/`.
- **Decision made or changed:** add an entry to `docs/decisions.md`.

Runtime data of the app itself (database, `app.log`) lives per user in `%APPDATA%\LykosReimarchiv\`, not here.
