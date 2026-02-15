package tui

import (
	"fmt"

	"github.com/vimgym/vimgym/internal/progress"
	"github.com/vimgym/vimgym/internal/puzzle"
)

func overallProgressText(prog *progress.Store, puzzles []puzzle.Puzzle) string {
	if prog == nil || len(puzzles) == 0 {
		return ""
	}

	solved, total, percent := prog.OverallProgress(puzzles)
	return fmt.Sprintf("Progress: %d/%d (%d%%)", solved, total, percent)
}
