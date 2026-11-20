# Decision log

Newest first. Each entry records what was decided and why.

## 2026-11-20 — Rhyme canvas

- **What:** the canvas lives below the list on the Reimarchiv page, with a splitter between them. It's an addition, not a replacement for the list.
- **What:** windows are per vowel class and always show the whole class, independent of the list search.
- **What:** the UI layout is persisted in SQLite (`ui_state`, a JSON document per key) instead of browser localStorage.
  - **Why:** all user data stays in one local file, so it's covered by future backups.

## 2026-11-20 — UI reduction

- **What:** the navigation is reduced to Reimarchiv, Reimzeilen, Generator and Einstellungen. Dashboard, status badges, roadmap cards and placeholder controls are removed.
  - **Why:** every visible element must serve a concrete purpose; the app should feel like a tool (order `2026-11-20-ui-reduction`).
- **What:** "Relatable-Sätze" is renamed "Reimzeilen".
- **What:** Reimarchiv is functional: add, search, filter by vowel, list, delete.
  - **Why:** the order defines this area as the working focus (search, categorise, display, add).
  - **Words schema:**
    - `vowel`: a, e, i, o, u. The user assigns it, so no spelling-based detection is needed, which follows the phonetic principle in `start-draft.md`.
    - `rhyme`: optional sound sequence from the stressed vowel. The list is grouped by it.
- **What:** code, comments, identifiers, scripts and developer docs are written in English. UI text stays German because it's the product language.
- **What:** the `workspace/` root holds documents, state and logs for both the user and the CLI. The app writes runtime logs to `%APPDATA%\LykosReimarchiv\logs\app.log`.
- **What:** `start-draft.md` moved to `workspace/docs/`.

## 2026-11-20 — MVP v0.1

- **What:** the stack is Go + Wails v2 (v2.17.1) + SQLite, as preferred in `start-draft.md`.
- **What:** the SQLite driver is `modernc.org/sqlite`.
  - **Why:** it's pure Go, so no CGO or C compiler is needed.
- **What:** the frontend is plain HTML/CSS/JS ES modules in `frontend/dist`, with no npm and no framework.
  - **Why:** it keeps the build simple.
- **What:** Wails uses `packaging/` as its build directory instead of `build/`.
  - **Why:** `build/` is excluded in `.gitignore`.
- **What:** `.gitignore` gets the exception `!frontend/dist/`.
  - **Why:** `frontend/dist` is source code here, not build output.
- **What:** schema migrations are a list in `internal/storage/sqlite.go`, versioned via `PRAGMA user_version`.
- **What:** the app works offline: it uses system fonts and loads no external resources.
