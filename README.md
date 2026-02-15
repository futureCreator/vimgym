# VimGym

> Master Vim through real code editing puzzles in your terminal.

VimGym is a TUI trainer that teaches Vim keybindings by presenting **Before → After** code transformation puzzles. Powered by an embedded Neovim instance, every keystroke behaves exactly like real Vim.

```
┌─ VimGym ─────────────────────────────────────────┐
│                                                   │
│  Level 12: Object Surgery                         │
│  Category: Text Objects  |  Par: 9 keystrokes     │
│                                                   │
│  ┌─ GOAL ───────────────────────────────────────┐ │
│  │ const name = "Goodbye";                      │ │
│  └──────────────────────────────────────────────┘ │
│                                                   │
│  ┌─ EDITOR (Neovim) ───────────────────────────┐ │
│  │ const name = "Hello, World!";                │ │
│  │                ^ cursor                      │ │
│  └──────────────────────────────────────────────┘ │
│                                                   │
│  Keystrokes: 0  |  Mode: NORMAL                   │
│                                                   │
│  Ctrl+H Hint  Ctrl+O Solution  Ctrl+R Reset       │
└───────────────────────────────────────────────────┘
```

## Features

- **Real Vim engine** — Neovim runs in `--embed` mode, so every command works exactly as expected. No incomplete emulation.
- **162 puzzles across 30 levels** — Progressive curriculum from `hjkl` basics to Vim Golf challenges.
- **Vim Golf scoring** — Each puzzle has a par (minimum keystrokes). Earn up to 3 stars by matching or beating it.
- **4 learning tracks** — Foundations, Editing, Power Moves, and Vim Golf.
- **Hints & solutions** — Get unstuck with hints or view the optimal solution with explanation.
- **Local progress** — Your results are saved locally in `~/.vimgym/`. No account required.
- **Modern TUI** — Built with Bubble Tea and Lip Gloss for a polished terminal experience.

## Learning Tracks

| Track | Levels | Topics |
|-------|--------|--------|
| **Foundations** | 1–5 | `hjkl`, word/line/block motions, `f`/`t` find |
| **Editing** | 6–10 | Insert, delete, change, yank & paste |
| **Power Moves** | 11–15 | Text objects, visual mode, search, substitute, macros |
| **Vim Golf** | 16–30 | Real-world combos, refactoring, speed challenges |

Levels unlock sequentially — clear all puzzles in a level (1 star or above) to unlock the next.

## Scoring

| Rating | Condition |
|--------|-----------|
| 3 stars | At or under par |
| 2 stars | At or under 1.5x par |
| 1 star | Cleared |

## Prerequisites

- **Go** 1.24+
- **Neovim** (`nvim` must be in your PATH)

## Install & Run

```bash
git clone https://github.com/futureCreator/vimgym.git
cd vimgym
go build ./cmd/vimgym/
./vimgym
```

Or run directly:

```bash
go run ./cmd/vimgym/
```

## Controls

| Key | Action |
|-----|--------|
| `Ctrl+H` | Toggle hint |
| `Ctrl+O` | Toggle optimal solution |
| `Ctrl+R` | Reset puzzle |
| `Ctrl+Q` | Quit to level select |

## Architecture

```
TUI (Bubble Tea) ←→ Neovim (--embed, msgpack-rpc) ←→ Buffer
```

1. Puzzle JSON loaded from embedded filesystem → buffer initialized with `before` text + cursor position
2. Keystrokes forwarded to Neovim → buffer changes in real time
3. Buffer text compared against `after` target → scored against par

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go |
| TUI | Bubble Tea + Lip Gloss |
| Vim Engine | Neovim `--embed` |
| Neovim Client | neovim/go-client (msgpack-rpc) |
| Storage | `~/.vimgym/` (JSON) |
| Puzzle Data | `embed.FS` (built into binary) |

## Project Structure

```
vimgym/
├── cmd/
│   ├── vimgym/          # App entry point
│   └── puzzlecheck/     # Puzzle data validation tool
├── internal/
│   ├── tui/             # Bubble Tea models & views
│   ├── nvim/            # Neovim embed client & buffer ops
│   ├── puzzle/          # Puzzle types, loader, validator, scorer
│   └── progress/        # Local JSON persistence
├── puzzles/             # Puzzle data (JSON, embedded at build)
├── go.mod
└── go.sum
```

## License

MIT
