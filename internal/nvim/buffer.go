package nvim

import (
	"fmt"
	"strings"

	"github.com/vimgym/vimgym/internal/puzzle"
)

// LoadPuzzle sets up the buffer with the puzzle's before state.
func (c *Client) LoadPuzzle(p puzzle.Puzzle) error {
	buf, err := c.nv.CurrentBuffer()
	if err != nil {
		return fmt.Errorf("getting current buffer: %w", err)
	}

	// Split the before text into lines
	lines := strings.Split(p.Before.Text, "\n")
	byteLines := make([][]byte, len(lines))
	for i, l := range lines {
		byteLines[i] = []byte(l)
	}

	// Set buffer contents
	if err := c.nv.SetBufferLines(buf, 0, -1, false, byteLines); err != nil {
		return fmt.Errorf("setting buffer lines: %w", err)
	}

	// Set cursor position (Neovim uses 1-indexed rows)
	win, err := c.nv.CurrentWindow()
	if err != nil {
		return fmt.Errorf("getting current window: %w", err)
	}

	row := p.Before.Cursor.Row + 1 // convert to 1-indexed
	col := p.Before.Cursor.Col
	if err := c.nv.SetWindowCursor(win, [2]int{row, col}); err != nil {
		return fmt.Errorf("setting cursor: %w", err)
	}

	// Ensure we're in normal mode
	c.Input("\x1b") // Esc

	return nil
}

// ResetPuzzle reloads the puzzle state (same as LoadPuzzle).
func (c *Client) ResetPuzzle(p puzzle.Puzzle) error {
	c.Input("\x1b\x1b") // Esc Esc
	return c.LoadPuzzle(p)
}
