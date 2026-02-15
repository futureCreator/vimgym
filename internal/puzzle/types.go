package puzzle

// CursorPos represents a cursor position in a buffer.
type CursorPos struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// BeforeState represents the initial state of a puzzle.
type BeforeState struct {
	Text   string    `json:"text"`
	Cursor CursorPos `json:"cursor"`
}

// AfterState represents the goal state of a puzzle.
type AfterState struct {
	Text string `json:"text"`
}

// Puzzle represents a single VimGym puzzle.
type Puzzle struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Track           int        `json:"track"`
	Level           int        `json:"level"`
	Category        string     `json:"category"`
	Difficulty      int        `json:"difficulty"`
	Before          BeforeState `json:"before"`
	After           AfterState  `json:"after"`
	Par             int        `json:"par"`
	Hint            string     `json:"hint"`
	OptimalSolution     string   `json:"optimalSolution"`
	SolutionExplanation string   `json:"solutionExplanation"`
	Tags                []string `json:"tags"`
}

// StarRating represents the score for a puzzle completion.
type StarRating int

const (
	NoStar    StarRating = 0
	OneStar   StarRating = 1
	TwoStar   StarRating = 2
	ThreeStar StarRating = 3
)

// String returns a display string for the star rating.
func (s StarRating) String() string {
	switch s {
	case ThreeStar:
		return "***"
	case TwoStar:
		return "** "
	case OneStar:
		return "*  "
	default:
		return "   "
	}
}
