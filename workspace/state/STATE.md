# Current state

_Last updated: 2026-11-20_

## Works

- **Desktop app:** Windows EXE `packaging\bin\LykosReimarchiv.exe` with the desktop shortcut "LYKOS REIMARCHIV".
- **Reimarchiv:**
  - Add a word with a vowel class (A/E/I/O/U) and an optional rhyme sound.
  - Live search over words and rhyme sounds.
  - Vowel filter with counts.
  - List grouped by rhyme sound.
  - Two-click delete.
  - Duplicate check that ignores case, including umlauts.
- **Rhyme canvas** (below the list, height adjustable via the splitter):
  - Dark, dotted surface: drag the background to pan, mouse wheel to zoom at the cursor (30–200 %), "Einpassen" fits all windows.
  - One window per vowel class (A/E/I/O/U) with its words grouped by rhyme sound. Windows can be moved by the header, resized at the corner, brought to front and closed.
  - Layout (open windows, positions, sizes, zoom, pane height) is saved in SQLite (`ui_state`) and restored on start.
- **Einstellungen:** shows version, database path, schema and size, and the log file path.
- **SQLite:** stored at `%APPDATA%\LykosReimarchiv\reimarchiv.db` with schema v3 (`app_meta`, `words`, `ui_state`).
- **App log:** written to `%APPDATA%\LykosReimarchiv\logs\app.log`.
- **Tests:** `go test ./...` covers storage: migrations, validation, list/filter/search, counts, delete and UI state.

## Navigation

| Area | Status |
| --- | --- |
| Reimarchiv | active — the only area in development |
| Reimzeilen | planned (navigation target only) |
| Generator | planned (navigation target only) |
| Einstellungen | active |

## Next steps (roadmap step 1: monosyllables, simple vowels)

- [ ] Decide whether the archive is limited to monosyllables for now, or whether a syllable count is stored.
- [ ] Edit an existing word (currently only add and delete).
- [ ] Notation for the rhyme sound: free text vs. IPA vs. a fixed sound inventory.
- [ ] Short vs. long vowel (/a/ vs. /aː/) as its own field or as part of the rhyme sound.

## Canvas ideas (not built, decide first)

- Windows per rhyme sound (e.g. "-acht") in addition to vowel classes.
- Keyboard controls for the canvas (zoom +/−, arrow keys to pan).
- Should the canvas windows follow the list search/filter? Currently they always show the full vowel class.

## Open questions

- Should empty Einstellungen sections (export, backup) only come back once they work? Currently they're removed.
