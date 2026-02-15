package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	nvimclient "github.com/vimgym/vimgym/internal/nvim"
	"github.com/vimgym/vimgym/internal/puzzle"
	"github.com/vimgym/vimgym/internal/progress"
)

type screen int

const (
	screenTrack screen = iota
	screenPuzzle
)

// App is the main Bubble Tea model.
type App struct {
	screen     screen
	trackView  TrackView
	puzzleView PuzzleView
	nvim       *nvimclient.Client
	puzzles    []puzzle.Puzzle
	progress   *progress.Store
	width      int
	height     int
	err        error
}

// NewApp creates the main application model.
func NewApp(puzzles []puzzle.Puzzle) (*App, error) {
	if len(puzzles) == 0 {
		return nil, fmt.Errorf("no puzzles found")
	}

	// Load progress
	prog, err := progress.New()
	if err != nil {
		return nil, fmt.Errorf("loading progress: %w", err)
	}

	app := &App{
		screen:   screenTrack,
		puzzles:  puzzles,
		progress: prog,
	}
	app.trackView = NewTrackView(puzzles, prog)

	return app, nil
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	switch a.screen {
	case screenTrack:
		return a.updateTrack(msg)
	case screenPuzzle:
		return a.updatePuzzle(msg)
	}

	return a, nil
}

func (a App) updateTrack(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case selectedPuzzle:
		// Start puzzle - create Neovim instance
		nv, err := nvimclient.New()
		if err != nil {
			a.err = fmt.Errorf("starting neovim: %w", err)
			return a, nil
		}
		a.nvim = nv
		a.screen = screenPuzzle
		a.puzzleView = NewPuzzleView(msg.puzzle, nv, a.progress, a.puzzles)
		a.puzzleView.width = a.width
		a.puzzleView.height = a.height
		return a, a.puzzleView.Init()
	default:
		var cmd tea.Cmd
		a.trackView, cmd = a.trackView.Update(msg)
		return a, cmd
	}
}

func (a App) updatePuzzle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case puzzleExitMsg:
		if msg.next {
			if next, ok := nextPuzzleInLevel(a.puzzles, a.puzzleView.puzzle); ok && a.nvim != nil {
				a.puzzleView = NewPuzzleView(next, a.nvim, a.progress, a.puzzles)
				return a, a.puzzleView.Init()
			}
			// No next puzzle in this level: go to level selection for current track.
			a.screen = screenTrack
			a.trackView = trackViewForLevel(a.puzzles, a.progress, a.puzzleView.puzzle)
			a.trackView.width = a.width
			a.trackView.height = a.height
			if a.nvim != nil {
				a.nvim.Close()
				a.nvim = nil
			}
			return a, nil
		}
		// Clean up Neovim and go back to track view
		if a.nvim != nil {
			a.nvim.Close()
			a.nvim = nil
		}
		a.screen = screenTrack
		// Refresh track view with updated progress, cursor on current level
		a.trackView = trackViewForLevel(a.puzzles, a.progress, a.puzzleView.puzzle)
		a.trackView.width = a.width
		a.trackView.height = a.height
		return a, nil
	default:
		var cmd tea.Cmd
		a.puzzleView, cmd = a.puzzleView.Update(msg)
		return a, cmd
	}
}

func nextPuzzleInLevel(all []puzzle.Puzzle, current puzzle.Puzzle) (puzzle.Puzzle, bool) {
	index := -1
	for i, p := range all {
		if p.ID == current.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return puzzle.Puzzle{}, false
	}

	for i := index + 1; i < len(all); i++ {
		p := all[i]
		if p.Track != current.Track || p.Level != current.Level {
			break
		}
		return p, true
	}

	return puzzle.Puzzle{}, false
}

func trackViewForLevel(all []puzzle.Puzzle, prog *progress.Store, current puzzle.Puzzle) TrackView {
	tv := NewTrackView(all, prog)
	tv.cursor = 0
	for i, entry := range tv.allLevels {
		if entry.level == current.Level {
			tv.cursor = i
			break
		}
	}
	return tv
}

func (a App) View() string {
	if a.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to exit.", a.err)
	}

	switch a.screen {
	case screenTrack:
		return a.trackView.View()
	case screenPuzzle:
		return a.puzzleView.View()
	}

	return ""
}
