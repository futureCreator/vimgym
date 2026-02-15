package nvim

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

// Client wraps a Neovim embedded instance.
type Client struct {
	nv *nvim.Nvim
}

// New starts a new embedded Neovim process and connects via msgpack-rpc.
func New() (*Client, error) {
	nv, err := nvim.NewChildProcess(
		nvim.ChildProcessArgs("--embed", "--clean", "-n"),
		nvim.ChildProcessServe(false),
	)
	if err != nil {
		return nil, fmt.Errorf("starting nvim: %w", err)
	}

	// Register no-op handler for UI "redraw" notifications to avoid log spam
	nv.RegisterHandler("redraw", func(...[]interface{}) {})
	go nv.Serve()

	// Set some sensible defaults for puzzle mode
	batch := nv.NewBatch()
	batch.Command("set noswapfile")
	batch.Command("set nobackup")
	batch.Command("set nowritebackup")
	batch.Command("set noundofile")
	batch.Command("set shortmess+=I") // no intro message
	if err := batch.Execute(); err != nil {
		nv.Close()
		return nil, fmt.Errorf("configuring nvim: %w", err)
	}

	// Attach a minimal UI so nvim_input processes keys through the event loop.
	// Without UI, nvim_input keys sit in the input buffer unprocessed.
	if err := nv.AttachUI(80, 24, map[string]interface{}{"rgb": true}); err != nil {
		nv.Close()
		return nil, fmt.Errorf("attaching UI: %w", err)
	}

	return &Client{nv: nv}, nil
}

// Close shuts down the Neovim process.
func (c *Client) Close() error {
	if c.nv != nil {
		c.nv.DetachUI()
		return c.nv.Close()
	}
	return nil
}

// ResizeUI updates the attached UI size to match the terminal.
func (c *Client) ResizeUI(width, height int) {
	if c.nv == nil {
		return
	}
	if width <= 0 || height <= 0 {
		return
	}
	_ = c.nv.TryResizeUI(width, height)
}

// Input sends key input to Neovim via feedkeys.
// Use feedkeys with immediate execution so buffer reads reflect the keystroke.
func (c *Client) Input(keys string) error {
	if c.nv == nil {
		return nil
	}

	// Queue raw input; nvim_input parses keycodes like <Esc> directly.
	_, err := c.nv.Input(keys)
	return err
}

// GetBufferText returns the full text content of the current buffer.
func (c *Client) GetBufferText() (string, error) {
	lines, err := c.GetLines()
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

// GetLines returns the buffer content as a slice of strings.
func (c *Client) GetLines() ([]string, error) {
	buf, err := c.nv.CurrentBuffer()
	if err != nil {
		return nil, fmt.Errorf("getting current buffer: %w", err)
	}

	lines, err := c.nv.BufferLines(buf, 0, -1, false)
	if err != nil {
		return nil, fmt.Errorf("getting buffer lines: %w", err)
	}

	strLines := make([]string, len(lines))
	for i, l := range lines {
		strLines[i] = string(l)
	}
	return strLines, nil
}

// GetCursor returns the current cursor position (0-indexed row, col).
func (c *Client) GetCursor() (int, int, error) {
	win, err := c.nv.CurrentWindow()
	if err != nil {
		return 0, 0, fmt.Errorf("getting current window: %w", err)
	}

	pos, err := c.nv.WindowCursor(win)
	if err != nil {
		return 0, 0, fmt.Errorf("getting cursor: %w", err)
	}

	// Neovim returns 1-indexed row, 0-indexed col
	return pos[0] - 1, pos[1], nil
}

// GetMode returns the current Neovim mode string.
func (c *Client) GetMode() (string, error) {
	var mode string
	if err := c.nv.Eval("mode()", &mode); err != nil {
		return "", fmt.Errorf("getting mode: %w", err)
	}
	return mode, nil
}

// ModeDisplayName converts a Neovim mode string to a display name.
func ModeDisplayName(mode string) string {
	switch {
	case strings.HasPrefix(mode, "n"):
		return "NORMAL"
	case strings.HasPrefix(mode, "i"):
		return "INSERT"
	case strings.HasPrefix(mode, "v"):
		return "VISUAL"
	case strings.HasPrefix(mode, "V"):
		return "V-LINE"
	case mode == "\x16": // Ctrl-V
		return "V-BLOCK"
	case strings.HasPrefix(mode, "c"):
		return "COMMAND"
	case strings.HasPrefix(mode, "R"):
		return "REPLACE"
	default:
		return strings.ToUpper(mode)
	}
}
