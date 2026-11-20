# LYKOS REIMARCHIV (repo: lyricforce)

Windows desktop app: Go + Wails v2 + SQLite, plain HTML/CSS/JS frontend in `frontend/dist` (no npm).

## Start of every session

1. Read `workspace/state/Handoff.jsonl` (handoff from the previous model), `workspace/state/STATE.md` and the latest file in `workspace/logs/sessions/`.
2. Product vision, phonetic rhyme logic and roadmap: `workspace/docs/start-draft.md`. Don't change decisions recorded there or in `workspace/docs/decisions.md` without asking.

## End of every session

- Update `workspace/state/STATE.md`, append to `workspace/logs/sessions/<YYYY-MM-DD>.md`.
- Write the handoff for the next model to `workspace/state/Handoff.jsonl` (one JSON object per line; types: meta, user, done, plan, problem, todo).
- Save new work orders unchanged to `workspace/docs/build-orders/`; record decisions in `workspace/docs/decisions.md`.

## Rules

- Code, identifiers, comments, scripts and developer docs in **English**. UI text is German (product language).
- Only Reimarchiv is in active development. Don't build UI for future features; every visible element needs a concrete purpose.
- Rhyme logic is phonetic, never orthographic (see start-draft.md).
- Schema changes: append a migration to `internal/storage/sqlite.go`, never edit an existing one.

## Commands

- Dev: `powershell -ExecutionPolicy Bypass -File scripts\dev.ps1`
- Build (+ shortcut): `powershell -ExecutionPolicy Bypass -File scripts\build.ps1 [-Shortcut]` (close the app first)
- Test: `go test ./...`
