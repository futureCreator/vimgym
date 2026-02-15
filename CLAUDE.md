# VimGym

Vim keybinding learning TUI — solve "Before → After" code editing puzzles powered by embedded Neovim.
See `SPEC.md` for full project specification (Korean).

## Tech Stack

Go, Bubble Tea + Lip Gloss (TUI), Neovim `--embed` via `neovim/go-client` (msgpack-rpc), local JSON storage in `~/.vimgym/`.

## Commands

```bash
go build ./cmd/vimgym/        # build
go test ./...                  # run all tests
go test ./internal/puzzle/     # run tests for a single package
go run ./cmd/vimgym/           # run
go run ./cmd/puzzlecheck/      # validate puzzle data (par vs optimalSolution)
go run ./cmd/puzzlecheck/ -all # show per-puzzle detail
```

Prerequisite: `nvim` must be installed.

## Architecture

TUI (Bubble Tea) ↔ Neovim (`--embed`, msgpack-rpc) ↔ Buffer

1. Puzzle JSON loaded from `embed.FS` → buffer initialized with `before` text + cursor position
2. Keystrokes forwarded to Neovim → buffer changes
3. Buffer text compared against `after` target → scored against `par`

## Package Layout

- `cmd/vimgym/` — app entry point
- `cmd/puzzlecheck/` — puzzle data QA tool (validates par matches optimalSolution keystroke count)
- `internal/tui/` — Bubble Tea models/views (app, puzzle_view, track_view, progress_view, styles)
- `internal/nvim/` — Neovim embed client + buffer ops
- `internal/puzzle/` — puzzle types, loader, validator, scorer
- `internal/progress/` — local JSON persistence (`~/.vimgym/progress.json`)
- `puzzles/` — puzzle data JSON files (per track) + `embed.go` for `go:embed`

## Puzzle Data

162 puzzles total (5-7 per level × 30 levels × 4 tracks):
- Track 1 (Foundations, Lv 1-5): 29 puzzles — hjkl, word, line, block, find motions
- Track 2 (Editing, Lv 6-10): 28 puzzles — insert, delete, delete-motion, change, yank-paste
- Track 3 (Power Moves, Lv 11-15): 29 puzzles — text objects, visual, search, substitute, repeat/macros
- Track 4 (Vim Golf, Lv 16-30): 76 puzzles — combo, marks, regex, visual-block, macros, refactor, graduation

## Puzzle JSON Format

```json
{
  "id": "text-objects-01",
  "title": "String Surgery",
  "track": 3, "level": 11,
  "category": "text-objects",
  "difficulty": 2,
  "before": { "text": "...", "cursor": { "row": 0, "col": 15 } },
  "after": { "text": "..." },
  "par": 9,
  "hint": "Try using ci\" to change inside quotes",
  "optimalSolution": "ci\"Goodbye<Esc>",
  "solutionExplanation": "ci\" changes inside quotes, type Goodbye, Esc to return to normal",
  "tags": ["ci\"", "text-object", "change"]
}
```

## Conventions

- Korean comments are acceptable
- Scoring: 3-star (≤ par), 2-star (≤ 1.5× par), 1-star (cleared)
- Levels unlock sequentially — previous level must be cleared (1-star+) to unlock next
