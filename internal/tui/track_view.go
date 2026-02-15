package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vimgym/vimgym/internal/progress"
	"github.com/vimgym/vimgym/internal/puzzle"
)

// Track names
var trackNames = map[int]string{
	1: "Foundations (Basic Movement)",
	2: "Editing (Insert/Delete/Change)",
	3: "Power Moves (Advanced)",
	4: "Vim Golf (Challenge)",
}

// Level descriptions
var levelDescriptions = map[int]string{
	1:  "Basic Movement",
	2:  "Word Motion",
	3:  "Line Motion",
	4:  "Block Motion",
	5:  "Find Motion",
	6:  "Insert",
	7:  "Delete Basic",
	8:  "Motion + Delete",
	9:  "Change",
	10: "Yank & Paste",
	11: "Text Objects",
	12: "Visual + Text Objects",
	13: "Search",
	14: "Substitute",
	15: "Repeat & Macros",
	16: "Combo Basics",
	17: "Text Object Combos",
	18: "Yank Combos",
	19: "Marks & Jumps",
	20: "Regex Substitute",
	21: "Visual Block",
	22: "Indent & Format",
	23: "Advanced Macros",
	24: "Multi-line Editing",
	25: "Refactoring",
	26: "Advanced Combos",
	27: "Speed Editing",
	28: "Code Transform",
	29: "Expert Commands",
	30: "Graduation",
}

type viewMode int

const (
	viewLevels viewMode = iota
	viewPuzzles
)

// levelEntry represents a level in the flat list.
type levelEntry struct {
	track int
	level int
}

// TrackView handles the level/puzzle selection screen.
type TrackView struct {
	puzzles  []puzzle.Puzzle
	progress *progress.Store
	mode     viewMode

	allLevels  []levelEntry
	puzzleList []puzzle.Puzzle

	cursor       int
	confirmReset bool
	width        int
	height       int
}

// NewTrackView creates a new level selection view.
func NewTrackView(puzzles []puzzle.Puzzle, prog *progress.Store) TrackView {
	tracks := puzzle.GetTracks(puzzles)
	var allLevels []levelEntry
	for _, t := range tracks {
		levels := puzzle.GetLevelsForTrack(puzzles, t)
		for _, l := range levels {
			allLevels = append(allLevels, levelEntry{track: t, level: l})
		}
	}

	return TrackView{
		puzzles:   puzzles,
		progress:  prog,
		mode:      viewLevels,
		allLevels: allLevels,
	}
}

// selectedPuzzle is a message sent when a puzzle is selected.
type selectedPuzzle struct {
	puzzle puzzle.Puzzle
}

func (v TrackView) Update(msg tea.Msg) (TrackView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		return v, nil
	case tea.KeyMsg:
		if v.confirmReset {
			switch msg.String() {
			case "y", "Y":
				_ = v.progress.Reset()
				v = NewTrackView(v.puzzles, v.progress)
				return v, tea.ClearScreen
			}
			// Any other key cancels the reset prompt.
			v.confirmReset = false
			return v, tea.ClearScreen
		}
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			v.cursor = min(v.cursor+1, v.maxCursor())
		case "ctrl+r":
			v.confirmReset = true
			return v, nil
		case "enter", "l":
			return v.selectItem()
		case "esc", "h", "backspace":
			return v.back(), nil
		case "q":
			if v.mode == viewLevels {
				return v, tea.Quit
			}
			return v.back(), nil
		}
	}
	return v, nil
}

func (v TrackView) View() string {
	var b strings.Builder

	width := v.width
	if width <= 0 {
		width = 80
	}
	height := v.height
	if height <= 0 {
		height = 24
	}

	switch v.mode {
	case viewLevels:
		headerLines := []string{
			titleStyle.MaxWidth(width).Render("VimGym - Select Level"),
		}
		if progressText := overallProgressText(v.progress, v.puzzles); progressText != "" {
			headerLines = append(headerLines, mutedStyle.MaxWidth(width).Render(progressText))
		}
		header := strings.Join(headerLines, "\n")

		lastTrack := 0
		itemIndex := 0
		var lines []string
		cursorLine := 0
		for _, entry := range v.allLevels {
			// Track header
			if entry.track != lastTrack {
				if lastTrack != 0 {
					lines = append(lines, "")
				}
				name := trackNames[entry.track]
				if name == "" {
					name = fmt.Sprintf("Track %d", entry.track)
				}
				lines = append(lines, trackHeaderStyle.Render(fmt.Sprintf("── Track %d: %s ──", entry.track, name)))
				lastTrack = entry.track
			}

			unlocked := v.progress.IsLevelUnlocked(entry.level, v.puzzles)
			desc := levelDescriptions[entry.level]
			if desc == "" {
				desc = fmt.Sprintf("Level %d", entry.level)
			}

			prefix := "  "
			style := unselectedStyle
			if itemIndex == v.cursor {
				prefix = "> "
				style = selectedStyle
				cursorLine = len(lines)
			}

			if !unlocked {
				style = lockedStyle
				lockIcon := " [locked]"
				lines = append(lines, fmt.Sprintf("%s%s%s", prefix, style.Render(fmt.Sprintf("Lv %d: %s", entry.level, desc)), style.Render(lockIcon)))
			} else {
				stars := v.progress.GetLevelStars(entry.level, v.puzzles)
				starStr := FormatStars(int(stars))
				lines = append(lines, fmt.Sprintf("%s%s  %s", prefix, style.Render(fmt.Sprintf("Lv %d: %s", entry.level, desc)), starStr))
			}
			itemIndex++
		}
		helpLine := "  j/k: navigate  enter: select  q: quit  Ctrl+R: reset progress"
		footer := helpStyle.MaxWidth(width).Render(helpLine)
		if v.confirmReset {
			footer = footer + "\n" + dangerStyle.MaxWidth(width).Render("Reset all progress? [y]es / [n]o")
		}
		available := height - lipgloss.Height(header) - lipgloss.Height(footer) - 2
		if available < 1 {
			available = 1
		}
		visible := windowLines(lines, cursorLine, available)

		b.WriteString(header)
		b.WriteString("\n\n")
		for i, line := range visible {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(fitWidth(line, width))
		}
		b.WriteString("\n\n")
		b.WriteString(footer)

	case viewPuzzles:
		level := v.puzzleList[0].Level
		desc := levelDescriptions[level]
		headerLines := []string{
			titleStyle.MaxWidth(width).Render(fmt.Sprintf("Level %d: %s", level, desc)),
		}
		if progressText := overallProgressText(v.progress, v.puzzles); progressText != "" {
			headerLines = append(headerLines, mutedStyle.MaxWidth(width).Render(progressText))
		}
		header := strings.Join(headerLines, "\n")

		var lines []string
		cursorLine := 0
		for i, p := range v.puzzleList {
			prefix := "  "
			style := unselectedStyle
			if i == v.cursor {
				prefix = "> "
				style = selectedStyle
				cursorLine = len(lines)
			}

			result := v.progress.GetBest(p.ID)
			starStr := FormatStars(int(result.Stars))
			keystrokeInfo := ""
			if result.Keystrokes > 0 {
				keystrokeInfo = mutedStyle.Render(fmt.Sprintf(" (%d keys, par %d)", result.Keystrokes, p.Par))
			} else {
				keystrokeInfo = mutedStyle.Render(fmt.Sprintf(" (par %d)", p.Par))
			}

			lines = append(lines, fmt.Sprintf("%s%s  %s%s", prefix, style.Render(p.Title), starStr, keystrokeInfo))
		}
		helpLine := "  j/k: navigate  enter: start  esc: back  Ctrl+R: reset progress"
		footer := helpStyle.MaxWidth(width).Render(helpLine)
		if v.confirmReset {
			footer = footer + "\n" + dangerStyle.MaxWidth(width).Render("Reset all progress? [y]es / [n]o")
		}
		available := height - lipgloss.Height(header) - lipgloss.Height(footer) - 2
		if available < 1 {
			available = 1
		}
		visible := windowLines(lines, cursorLine, available)

		b.WriteString(header)
		b.WriteString("\n\n")
		for i, line := range visible {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(fitWidth(line, width))
		}
		b.WriteString("\n\n")
		b.WriteString(footer)
	}

	return b.String()
}

func (v TrackView) maxCursor() int {
	switch v.mode {
	case viewLevels:
		return max(0, len(v.allLevels)-1)
	case viewPuzzles:
		return max(0, len(v.puzzleList)-1)
	}
	return 0
}

func (v TrackView) selectItem() (TrackView, tea.Cmd) {
	switch v.mode {
	case viewLevels:
		if v.cursor < len(v.allLevels) {
			entry := v.allLevels[v.cursor]
			if !v.progress.IsLevelUnlocked(entry.level, v.puzzles) {
				return v, nil
			}
			v.puzzleList = puzzle.GetPuzzlesForLevel(v.puzzles, entry.level)
			v.mode = viewPuzzles
			v.cursor = 0
		}
	case viewPuzzles:
		if v.cursor < len(v.puzzleList) {
			p := v.puzzleList[v.cursor]
			return v, func() tea.Msg { return selectedPuzzle{puzzle: p} }
		}
	}
	return v, nil
}

func (v TrackView) back() TrackView {
	switch v.mode {
	case viewPuzzles:
		// Go back to levels, restore cursor to the level we came from
		if len(v.puzzleList) > 0 {
			targetLevel := v.puzzleList[0].Level
			for i, entry := range v.allLevels {
				if entry.level == targetLevel {
					v.cursor = i
					break
				}
			}
		} else {
			v.cursor = 0
		}
		v.mode = viewLevels
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func windowLines(lines []string, cursor, height int) []string {
	if height <= 0 || len(lines) <= height {
		return lines
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(lines) {
		cursor = len(lines) - 1
	}
	start := cursor - height/2
	if start < 0 {
		start = 0
	}
	if start+height > len(lines) {
		start = len(lines) - height
	}
	return lines[start : start+height]
}

func fitWidth(line string, width int) string {
	if width <= 0 {
		return line
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(line)
}
