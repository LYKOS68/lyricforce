# LYKOS REIMARCHIV

Local Windows desktop app for German rap and songwriting: a personal, phonetically organised rhyme archive.
Vision, rhyme logic and roadmap: [`workspace/docs/start-draft.md`](workspace/docs/start-draft.md). Current state: [`workspace/state/STATE.md`](workspace/state/STATE.md).

## Requirements

- Windows 10/11 with the WebView2 runtime (preinstalled on Windows 11)
- Go ≥ 1.24
- Wails CLI v2: `go install github.com/wailsapp/wails/v2/cmd/wails@v2.17.1`

Node.js/npm aren't needed.

## Commands

| Purpose | Command |
| --- | --- |
| Development (live reload) | `powershell -ExecutionPolicy Bypass -File scripts\dev.ps1` |
| Build executable | `powershell -ExecutionPolicy Bypass -File scripts\build.ps1` |
| Build + desktop shortcut | `powershell -ExecutionPolicy Bypass -File scripts\build.ps1 -Shortcut` |
| Desktop shortcut only | `powershell -ExecutionPolicy Bypass -File scripts\create-shortcut.ps1` |
| Tests | `go test ./...` |

- **Build output:** `packaging\bin\LykosReimarchiv.exe`. Close the app before building, because Windows locks the running executable.
- **Build logs:** `workspace\logs\build-*.log`.
- **Desktop shortcut:** `LYKOS REIMARCHIV.lnk`, pointing at that executable.

## Data locations

| What | Where |
| --- | --- |
| Database | `%APPDATA%\LykosReimarchiv\reimarchiv.db` (created on first start) |
| App log | `%APPDATA%\LykosReimarchiv\logs\app.log` |
| Docs, state, work logs | `workspace/` (see [`workspace/README.md`](workspace/README.md)) |

## Project structure

```
main.go                  Wails bootstrap, window, embedded UI, log file
app.go                   Go → UI bindings
internal/appdir/         per-user data paths
internal/modules/        navigation targets (active / planned)
internal/storage/        SQLite: migrations, words, UI state (canvas layout)
frontend/dist/
  index.html, styles.css shell and design tokens
  app.js                 navigation and routing (#/<id>)
  api.js                 the only access point to the Go backend
  components.js          shared helpers
  views/                 one file per area/feature (rhymes.js, canvas.js, settings.js)
packaging/               Wails build dir (icon, manifest); bin/ isn't versioned
scripts/                 dev, build, shortcut
workspace/               docs, state, logs for the user and the CLI
```

Areas without a view file are planned and render as a one-line page.
