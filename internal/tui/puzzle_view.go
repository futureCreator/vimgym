package tui

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	nvimclient "github.com/vimgym/vimgym/internal/nvim"
	"github.com/vimgym/vimgym/internal/puzzle"
	"github.com/vimgym/vimgym/internal/progress"
)

var debugKeysEnabled = os.Getenv("VIMGYM_DEBUG_KEYS") != ""

type puzzleState int

const (
	statePlaying puzzleState = iota
	stateCleared
)

// puzzleExitMsg is sent when leaving puzzle view.
type puzzleExitMsg struct {
	next bool
}

// PuzzleView handles the puzzle solving screen.
type PuzzleView struct {
	puzzle     puzzle.Puzzle
	allPuzzles []puzzle.Puzzle
	nvim       *nvimclient.Client
	progress   *progress.Store
	state      puzzleState

	// Runtime state
	keystrokes int
	mode       string
	lines      []string
	cursorRow  int
	cursorCol  int
	showHint   bool
	showSolution bool
	stars      puzzle.StarRating
	width      int
	height     int
	// pendingKeys holds a prefix command waiting for the next key (ex: "r").
	pendingKeys string
	// pendingOperator indicates we're waiting for a motion/text object after an operator (d/c/y).
	pendingOperator bool
	// pendingNeedsChar indicates the pending keys need one more character (r/f/t/F/T).
	pendingNeedsChar bool
	// pendingTextObject indicates we're waiting for a text object after i/a.
	pendingTextObject bool
	// pendingHasCount tracks whether an operator has started a count (e.g. d2w).
	pendingHasCount bool
	// pendingCount buffers a leading count before a motion/operator (e.g. 4w, 3dw).
	pendingCount string
}

// NewPuzzleView creates a new puzzle view.
func NewPuzzleView(p puzzle.Puzzle, nv *nvimclient.Client, prog *progress.Store, all []puzzle.Puzzle) PuzzleView {
	return PuzzleView{
		puzzle:     p,
		allPuzzles: all,
		nvim:       nv,
		progress:   prog,
		state:      statePlaying,
		mode:       "NORMAL",
	}
}

type initPuzzleMsg struct{}
type nvimSyncMsg struct{}

// Init initializes the puzzle view by loading the puzzle into Neovim.
func (v PuzzleView) Init() tea.Cmd {
	return func() tea.Msg {
		return initPuzzleMsg{}
	}
}

// syncReadBuffer reads buffer state from Neovim synchronously.
func (v *PuzzleView) syncReadBuffer() {
	if v.nvim == nil {
		return
	}

	lines, err := v.nvim.GetLines()
	if err != nil {
		return
	}
	v.lines = lines

	row, col, err := v.nvim.GetCursor()
	if err == nil {
		v.cursorRow = row
		v.cursorCol = col
	}

	modeStr, err := v.nvim.GetMode()
	if err == nil {
		v.mode = nvimclient.ModeDisplayName(modeStr)
	}
}

// syncCheckClear reads buffer text and checks for puzzle completion.
func (v *PuzzleView) syncCheckClear() {
	if v.state != statePlaying || v.nvim == nil {
		return
	}
	text, err := v.nvim.GetBufferText()
	if err != nil {
		return
	}
	if puzzle.Validate(text, v.puzzle.After.Text) {
		v.state = stateCleared
		v.stars = puzzle.Score(v.keystrokes, v.puzzle.Par)
		v.progress.SetBest(v.puzzle.ID, v.stars, v.keystrokes)
		v.progress.Save()
	}
}

func (v PuzzleView) Update(msg tea.Msg) (PuzzleView, tea.Cmd) {
	switch msg := msg.(type) {
	case initPuzzleMsg:
		v.nvim.LoadPuzzle(v.puzzle)
		v.syncReadBuffer()
		return v, nil
	case nvimSyncMsg:
		v.syncReadBuffer()
		v.syncCheckClear()
		return v, nil

	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		if v.nvim != nil {
			v.nvim.ResizeUI(msg.Width, msg.Height)
		}
		return v, nil

	case tea.KeyMsg:
		if v.state == stateCleared {
			switch msg.String() {
			case "enter", "q", "esc":
				v.clearPending()
				return v, func() tea.Msg { return puzzleExitMsg{next: msg.String() == "enter"} }
			case "r", "ctrl+r":
				v.state = statePlaying
				v.keystrokes = 0
				v.showHint = false
				v.showSolution = false
				v.clearPending()
				v.nvim.LoadPuzzle(v.puzzle)
				v.syncReadBuffer()
				return v, nil
			}
			return v, nil
		}

		// Playing state controls
		switch msg.String() {
		case "ctrl+q":
			v.clearPending()
			return v, func() tea.Msg { return puzzleExitMsg{next: false} }
		case "ctrl+r":
			v.keystrokes = 0
			v.showHint = false
			v.showSolution = false
			v.clearPending()
			v.nvim.LoadPuzzle(v.puzzle)
			v.syncReadBuffer()
			return v, nil
		case "ctrl+h":
			v.showHint = !v.showHint
			return v, nil
		case "ctrl+o":
			v.showSolution = !v.showSolution
			return v, nil
		default:
			keys := translateKey(msg)
			debugKeyInput(msg, keys)
			return v.handleNvimInput(keys)
		}
	default:
		if keys := translateCSIu(msg); keys != "" {
			debugKeyInput(msg, keys)
			return v.handleNvimInput(keys)
		}
	}

	return v, nil
}

func (v PuzzleView) View() string {
	width := v.width
	if width <= 0 {
		width = 80
	}
	height := v.height
	if height <= 0 {
		height = 24
	}

	contentWidth := width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}
	innerWidth := contentWidth - 4
	if innerWidth < 1 {
		innerWidth = 1
	}

	maxEditorLines := len(v.lines)
	if maxEditorLines == 0 {
		maxEditorLines = 1
	}
	if height > 0 {
		maxEditorLines = min(maxEditorLines, height)
	}

	maxGoalLines := countLines(v.puzzle.After.Text)
	if height <= 0 {
		return v.renderView(contentWidth, innerWidth, maxGoalLines, maxEditorLines)
	}

	for editorLines := maxEditorLines; editorLines >= 1; editorLines-- {
		for goalLines := maxGoalLines; goalLines >= 1; goalLines-- {
			view := v.renderView(contentWidth, innerWidth, goalLines, editorLines)
			if lipgloss.Height(view) <= height {
				return view
			}
		}
	}

	return v.renderView(contentWidth, innerWidth, 1, 1)
}

func (v PuzzleView) renderView(contentWidth, innerWidth, goalLines, editorLines int) string {
	header := fmt.Sprintf("Level %d: %s", v.puzzle.Level, v.puzzle.Title)
	info := fmt.Sprintf("Category: %s", v.puzzle.Category)
	if progressText := overallProgressText(v.progress, v.allPuzzles); progressText != "" {
		info += "  " + progressText
	}
	headerBlock := titleStyle.MaxWidth(contentWidth).Render(header)
	infoBlock := mutedStyle.MaxWidth(contentWidth).Render(info)

	goalLabel := labelStyle.Render(" GOAL ")
	goalContent := v.renderGoalContent(goalLines)
	goalBox := goalBoxStyle.Width(contentWidth).Render(goalContent)

	editorLabel := labelStyle.Render(" EDITOR ")
	editorContent := v.renderBuffer(innerWidth, editorLines)
	editorBox := editorBoxStyle.Width(contentWidth).Render(editorContent)

	modeDisplay := ModeStyle(v.mode).Render(fmt.Sprintf(" %s ", v.mode))
	keystrokeDisplay := fmt.Sprintf("Keystrokes: %d", v.keystrokes)
	parDisplay := mutedStyle.Render(fmt.Sprintf("(par: %d)", v.puzzle.Par))
	statusLine := fmt.Sprintf("%s  %s %s", modeDisplay, keystrokeDisplay, parDisplay)
	statusBlock := statusBarStyle.MaxWidth(contentWidth).Render(statusLine)

	parts := []string{
		headerBlock,
		infoBlock,
		"",
		goalLabel,
		goalBox,
		"",
		editorLabel,
		editorBox,
		"",
		statusBlock,
	}

	if v.showHint {
		parts = append(parts, hintStyle.Width(contentWidth).Render("Hint: "+v.puzzle.Hint))
	}
	if v.state == statePlaying && v.showSolution && v.puzzle.OptimalSolution != "" {
		parts = append(parts, solutionStyle.Width(contentWidth).Render("Solution: "+v.puzzle.OptimalSolution))
		if v.puzzle.SolutionExplanation != "" {
			parts = append(parts, explanationStyle.Width(contentWidth).Render(v.puzzle.SolutionExplanation))
		}
	}

	if v.state == stateCleared {
		starDisplay := FormatStars(int(v.stars))
		clearMsg := fmt.Sprintf(
			"Cleared! %s\n\nKeystrokes: %d  (par: %d)\nOptimal: %s\n\n[enter] next  [r] retry  [q] back",
			starDisplay, v.keystrokes, v.puzzle.Par, v.puzzle.OptimalSolution,
		)
		parts = append(parts, "", successStyle.MaxWidth(contentWidth).Render(clearMsg))
	} else {
		helpLine := "Ctrl+H: hint  Ctrl+O: solution  Ctrl+R: reset  Ctrl+Q: quit"
		parts = append(parts, helpStyle.MaxWidth(contentWidth).Render(helpLine))
	}

	return strings.Join(parts, "\n")
}

func (v PuzzleView) renderGoalContent(height int) string {
	if height < 1 {
		height = 1
	}
	afterLines := strings.Split(v.puzzle.After.Text, "\n")
	beforeLines := strings.Split(v.puzzle.Before.Text, "\n")
	focusRow := goalFocusRow(beforeLines, afterLines)
	start, end := windowRange(len(afterLines), focusRow, height)
	return strings.Join(afterLines[start:end], "\n")
}

func goalFocusRow(before, after []string) int {
	maxLines := max(len(before), len(after))
	if maxLines == 0 {
		return 0
	}
	for i := 0; i < maxLines; i++ {
		var b string
		if i < len(before) {
			b = before[i]
		}
		var a string
		if i < len(after) {
			a = after[i]
		}
		if b != a {
			if len(after) == 0 {
				return 0
			}
			if i >= len(after) {
				return len(after) - 1
			}
			return i
		}
	}
	if len(after) == 0 {
		return 0
	}
	return len(after) - 1
}

func countLines(text string) int {
	if text == "" {
		return 1
	}
	return strings.Count(text, "\n") + 1
}

// renderBuffer renders the Neovim buffer content with cursor highlight.
func (v PuzzleView) renderBuffer(width, height int) string {
	if len(v.lines) == 0 {
		return mutedStyle.Render("(loading...)")
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	start, end := windowRange(len(v.lines), v.cursorRow, height)
	rendered := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		line := v.lines[i]
		if i == v.cursorRow {
			rendered = append(rendered, v.renderLineWithCursor(line, v.cursorCol, width))
		} else {
			rendered = append(rendered, truncateLine(line, width))
		}
	}

	return strings.Join(rendered, "\n")
}

// renderLineWithCursor renders a line with the cursor position highlighted.
func (v PuzzleView) renderLineWithCursor(line string, col int, width int) string {
	if width < 1 {
		width = 1
	}
	runes := []rune(line)
	if col < 0 {
		col = 0
	}
	if col > len(runes) {
		col = len(runes)
	}

	if len(runes) <= width {
		cursorIdx := -1
		if col < len(runes) {
			cursorIdx = col
		}
		return renderCursorInRunes(runes, cursorIdx, col == len(runes), width)
	}

	start := col - width/2
	if start < 0 {
		start = 0
	}
	if start+width > len(runes) {
		start = len(runes) - width
	}
	end := start + width
	visible := runes[start:end]
	cursorIdx := -1
	if col >= start && col < end {
		cursorIdx = col - start
	}
	if start > 0 && cursorIdx != 0 && len(visible) > 0 {
		visible[0] = '~'
	}
	if end < len(runes) && cursorIdx != len(visible)-1 && len(visible) > 0 {
		visible[len(visible)-1] = '~'
	}
	return renderCursorInRunes(visible, cursorIdx, col == len(runes) && end == len(runes), width)
}

func renderCursorInRunes(runes []rune, cursorIdx int, showCursorSpace bool, width int) string {
	var b strings.Builder
	for i, r := range runes {
		if i == cursorIdx {
			b.WriteString(cursorStyle.Render(string(r)))
			continue
		}
		b.WriteRune(r)
	}
	if cursorIdx == -1 && showCursorSpace && len(runes) < width {
		b.WriteString(cursorStyle.Render(" "))
	}
	return b.String()
}

func truncateLine(line string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(line)
	if len(runes) <= width {
		return line
	}
	if width == 1 {
		return "~"
	}
	return string(runes[:width-1]) + "~"
}

func windowRange(total, cursor, height int) (int, int) {
	if height <= 0 || total <= height {
		return 0, total
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= total {
		cursor = total - 1
	}
	start := cursor - height/2
	if start < 0 {
		start = 0
	}
	if start+height > total {
		start = total - height
	}
	return start, start + height
}

// translateKey converts Bubble Tea key messages to Neovim input strings.
func translateKey(msg tea.KeyMsg) string {
	// Special keys
	switch msg.Type {
	case tea.KeyEsc:
		return "<Esc>"
	case tea.KeyEnter:
		return "<CR>"
	case tea.KeyTab:
		return "<Tab>"
	case tea.KeyBackspace:
		return "<BS>"
	case tea.KeyDelete:
		return "<Del>"
	case tea.KeyUp:
		return "<Up>"
	case tea.KeyDown:
		return "<Down>"
	case tea.KeyLeft:
		return "<Left>"
	case tea.KeyRight:
		return "<Right>"
	case tea.KeyHome:
		return "<Home>"
	case tea.KeyEnd:
		return "<End>"
	case tea.KeyPgUp:
		return "<PageUp>"
	case tea.KeyPgDown:
		return "<PageDown>"
	case tea.KeySpace:
		return " "
	case tea.KeyRunes:
		// Escape literal "<" so nvim_input doesn't treat it as a keycode.
		return strings.ReplaceAll(string(msg.Runes), "<", "<LT>")
	}

	// ctrl combinations handled by Bubble Tea
	str := msg.String()
	if strings.HasPrefix(str, "ctrl+") {
		letter := strings.TrimPrefix(str, "ctrl+")
		return "<C-" + letter + ">"
	}

	return ""
}

// shouldBufferKey returns true for commands that require a following key.
func shouldBufferKey(keys string) bool {
	switch keys {
	case "r", "f", "t", "F", "T", "m", "'", "`", "g", "z", "@", "q", "[", "]":
		return true
	}
	return false
}

func shouldStartOperator(keys string) bool {
	switch keys {
	case "d", "c", "y":
		return true
	}
	return false
}

func keyNeedsChar(keys string) bool {
	switch keys {
	case "r", "f", "t", "F", "T":
		return true
	}
	return false
}

func isDigitKey(keys string) bool {
	if len(keys) != 1 {
		return false
	}
	return keys[0] >= '0' && keys[0] <= '9'
}

func isTextObjectPrefix(keys string) bool {
	switch keys {
	case "i", "a":
		return true
	}
	return false
}

func isMotionCharPrefix(keys string) bool {
	switch keys {
	case "f", "t", "F", "T":
		return true
	}
	return false
}

func (v PuzzleView) handleNvimInput(keys string) (PuzzleView, tea.Cmd) {
	if keys == "" {
		return v, nil
	}

	v.keystrokes++

	// Do not buffer in insert/replace/command mode.
	if v.mode != "NORMAL" {
		v.clearPending()
		v.applyImmediateMode(keys)
		return v, v.inputAndSync(keys)
	}

	// If we have a pending prefix command, combine and send together.
	if v.pendingKeys != "" {
		if keys == "<Esc>" {
			v.clearPending()
			v.applyImmediateMode(keys)
			return v, v.inputAndSync(keys)
		}

		if v.pendingOperator {
			combined := v.pendingKeys + keys
			sent := v.handleOperatorPending(keys)
			if sent {
				v.applyImmediateMode(combined)
				return v, v.scheduleSync()
			}
			return v, nil
		}

		if v.pendingNeedsChar || v.pendingTextObject {
			combined := v.pendingKeys + keys
			v.clearPending()
			v.applyImmediateMode(combined)
			return v, v.inputAndSync(combined)
		}

		combined := v.pendingKeys + keys
		v.clearPending()
		v.applyImmediateMode(combined)
		return v, v.inputAndSync(combined)
	}

	// Handle leading counts (e.g. 4w, 3dw, 10fX).
	if v.pendingCount != "" {
		if keys == "<Esc>" {
			v.pendingCount = ""
			v.applyImmediateMode(keys)
			return v, v.inputAndSync(keys)
		}

		if isDigitKey(keys) {
			v.pendingCount += keys
			return v, nil
		}

		if shouldStartOperator(keys) {
			v.pendingKeys = v.pendingCount + keys
			v.pendingOperator = true
			v.pendingCount = ""
			return v, nil
		}

		if shouldBufferKey(keys) {
			v.pendingKeys = v.pendingCount + keys
			v.pendingNeedsChar = keyNeedsChar(keys)
			v.pendingCount = ""
			return v, nil
		}

		combined := v.pendingCount + keys
		v.pendingCount = ""
		v.applyImmediateMode(combined)
		return v, v.inputAndSync(combined)
	}

	if isDigitKey(keys) {
		if keys != "0" {
			v.pendingCount = keys
			return v, nil
		}
		// "0" is a motion when it's the first digit.
		v.applyImmediateMode(keys)
		return v, v.inputAndSync(keys)
	}

	// Buffer prefix commands that require a following key.
	if shouldStartOperator(keys) {
		v.pendingKeys = keys
		v.pendingOperator = true
		return v, nil
	}

	if shouldBufferKey(keys) {
		v.pendingKeys = keys
		v.pendingNeedsChar = keyNeedsChar(keys)
		return v, nil
	}

	v.applyImmediateMode(keys)
	return v, v.inputAndSync(keys)
}

// applyImmediateMode updates the local mode when a key deterministically changes it.
// This avoids misclassifying fast follow-up keys before the next nvim sync.
func (v *PuzzleView) applyImmediateMode(keys string) {
	if keys == "<Esc>" {
		v.mode = "NORMAL"
		return
	}

	if v.mode != "NORMAL" {
		return
	}

	switch keys {
	case "i", "I", "a", "A", "o", "O", "s", "S", "C":
		v.mode = "INSERT"
		return
	case "R":
		v.mode = "REPLACE"
		return
	case "v":
		v.mode = "VISUAL"
		return
	case "V":
		v.mode = "V-LINE"
		return
	case "<C-v>":
		v.mode = "V-BLOCK"
		return
	case ":", "/", "?":
		v.mode = "COMMAND"
		return
	}

	if entersInsertAfterChange(keys) {
		v.mode = "INSERT"
	}
}

func entersInsertAfterChange(keys string) bool {
	if keys == "" {
		return false
	}

	i := 0
	for i < len(keys) && keys[i] >= '0' && keys[i] <= '9' {
		i++
	}
	if i >= len(keys) {
		return false
	}
	if keys[i] != 'c' {
		return false
	}
	// Bare "c" (or count + "c") waits for a motion; no mode change yet.
	return i+1 < len(keys)
}

func translateCSIu(msg tea.Msg) string {
	bytes, ok := csiBytes(msg)
	if !ok || len(bytes) < 3 {
		return ""
	}
	if bytes[0] != 0x1b || bytes[1] != '[' {
		return ""
	}
	if bytes[len(bytes)-1] != 'u' {
		return ""
	}

	body := string(bytes[2 : len(bytes)-1])
	parts := strings.Split(body, ";")
	if len(parts) == 0 {
		return ""
	}

	code, err := strconv.Atoi(parts[0])
	if err != nil {
		return ""
	}

	r := rune(code)
	if !utf8.ValidRune(r) {
		return ""
	}

	mod := 1
	if len(parts) > 1 {
		if parsed, err := strconv.Atoi(parts[1]); err == nil {
			mod = parsed
		}
	}

	mask := mod - 1
	if mask&4 != 0 { // ctrl modifier
		return ctrlKeyString(r)
	}

	return string(r)
}

func debugKeyInput(msg tea.Msg, keys string) {
	if !debugKeysEnabled {
		return
	}

	f, err := os.OpenFile("/tmp/vimgym-keys.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	msgStr := ""
	if s, ok := msg.(fmt.Stringer); ok {
		msgStr = s.String()
	}

	raw := debugRawBytes(msg)
	fmt.Fprintf(
		f,
		"%s type=%T msg=%q keys=%q raw=% x\n",
		time.Now().Format(time.RFC3339Nano),
		msg,
		msgStr,
		keys,
		raw,
	)
}

func debugRawBytes(msg tea.Msg) []byte {
	if b, ok := csiBytes(msg); ok {
		return b
	}

	v := reflect.ValueOf(msg)
	if v.Kind() == reflect.Uint8 {
		return []byte{byte(v.Uint())}
	}

	if km, ok := msg.(tea.KeyMsg); ok {
		if km.Type == tea.KeyRunes && len(km.Runes) > 0 {
			return []byte(string(km.Runes))
		}
	}

	return nil
}

func csiBytes(msg tea.Msg) ([]byte, bool) {
	v := reflect.ValueOf(msg)
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
		return nil, false
	}

	out := make([]byte, v.Len())
	reflect.Copy(reflect.ValueOf(out), v)
	return out, true
}

func ctrlKeyString(r rune) string {
	switch {
	case r >= 'a' && r <= 'z':
		return "<C-" + string(r) + ">"
	case r >= 'A' && r <= 'Z':
		return "<C-" + strings.ToLower(string(r)) + ">"
	default:
		return "<C-" + string(r) + ">"
	}
}

func (v *PuzzleView) handleOperatorPending(keys string) bool {
	// Waiting for text object (diw, ci", etc)
	if v.pendingTextObject {
		combined := v.pendingKeys + keys
		v.clearPending()
		v.nvim.Input(combined)
		return true
	}

	// Waiting for character after f/t/F/T (dfx, cty, etc)
	if v.pendingNeedsChar {
		combined := v.pendingKeys + keys
		v.clearPending()
		v.nvim.Input(combined)
		return true
	}

	// Double-operator (dd/cc/yy)
	if len(v.pendingKeys) == 1 && keys == v.pendingKeys && !v.pendingHasCount {
		combined := v.pendingKeys + keys
		v.clearPending()
		v.nvim.Input(combined)
		return true
	}

	// Counts (d2w, c3e, etc)
	if isDigitKey(keys) {
		if keys == "0" && !v.pendingHasCount {
			combined := v.pendingKeys + keys
			v.clearPending()
			v.nvim.Input(combined)
			return true
		}
		v.pendingKeys += keys
		v.pendingHasCount = true
		return false
	}

	// Text objects (diw, da")
	if isTextObjectPrefix(keys) {
		v.pendingKeys += keys
		v.pendingTextObject = true
		return false
	}

	// Motions requiring a character (dfx)
	if isMotionCharPrefix(keys) {
		v.pendingKeys += keys
		v.pendingNeedsChar = true
		return false
	}

	combined := v.pendingKeys + keys
	v.clearPending()
	v.nvim.Input(combined)
	return true
}

func (v *PuzzleView) clearPending() {
	v.pendingKeys = ""
	v.pendingOperator = false
	v.pendingNeedsChar = false
	v.pendingTextObject = false
	v.pendingHasCount = false
	v.pendingCount = ""
}

func (v *PuzzleView) scheduleSync() tea.Cmd {
	return tea.Tick(10*time.Millisecond, func(time.Time) tea.Msg {
		return nvimSyncMsg{}
	})
}

func (v *PuzzleView) inputAndSync(keys string) tea.Cmd {
	if v.nvim != nil {
		v.nvim.Input(keys)
	}
	return v.scheduleSync()
}

// renderBufferWithHighlight applies diff highlighting compared to goal.
func renderBufferWithHighlight(lines []string, goalText string, cursorRow, cursorCol int) string {
	goalLines := strings.Split(goalText, "\n")
	var rendered []string

	for i, line := range lines {
		var goalLine string
		if i < len(goalLines) {
			goalLine = goalLines[i]
		}

		if line == goalLine {
			// Line matches goal - render normally
			rendered = append(rendered, line)
		} else {
			// Line differs - render with dim
			style := lipgloss.NewStyle().Foreground(colorWarning)
			rendered = append(rendered, style.Render(line))
		}
	}

	return strings.Join(rendered, "\n")
}
