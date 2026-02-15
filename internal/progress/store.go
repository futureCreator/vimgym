package progress

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/vimgym/vimgym/internal/puzzle"
)

const progressFile = "progress.json"

// PuzzleResult stores the best result for a puzzle.
type PuzzleResult struct {
	Stars      puzzle.StarRating `json:"stars"`
	Keystrokes int               `json:"keystrokes"`
}

// Store manages progress persistence.
type Store struct {
	dir     string
	Results map[string]PuzzleResult `json:"results"` // keyed by puzzle ID
}

// New creates a new progress store.
func New() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home dir: %w", err)
	}

	dir := filepath.Join(home, ".vimgym")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating vimgym dir: %w", err)
	}

	s := &Store{
		dir:     dir,
		Results: make(map[string]PuzzleResult),
	}

	// Load existing progress if available
	s.Load()

	return s, nil
}

// Load reads progress from disk.
func (s *Store) Load() error {
	path := filepath.Join(s.dir, progressFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading progress: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("parsing progress: %w", err)
	}
	if s.Results == nil {
		s.Results = make(map[string]PuzzleResult)
	}
	return nil
}

// Save writes progress to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling progress: %w", err)
	}

	path := filepath.Join(s.dir, progressFile)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing progress: %w", err)
	}
	return nil
}

// Reset clears all progress and persists the empty state.
func (s *Store) Reset() error {
	s.Results = make(map[string]PuzzleResult)
	return s.Save()
}

// GetBest returns the best result for a puzzle, or zero value if not attempted.
func (s *Store) GetBest(puzzleID string) PuzzleResult {
	return s.Results[puzzleID]
}

// SetBest updates the best result for a puzzle if it's better than existing.
func (s *Store) SetBest(puzzleID string, stars puzzle.StarRating, keystrokes int) {
	existing, ok := s.Results[puzzleID]
	if !ok || stars > existing.Stars || (stars == existing.Stars && keystrokes < existing.Keystrokes) {
		s.Results[puzzleID] = PuzzleResult{
			Stars:      stars,
			Keystrokes: keystrokes,
		}
	}
}

// IsLevelUnlocked checks if a level is unlocked.
// Level 1 is always unlocked. Other levels require all puzzles in the previous level
// to have at least 1 star.
func (s *Store) IsLevelUnlocked(level int, allPuzzles []puzzle.Puzzle) bool {
	if level <= 1 {
		return true
	}

	// Find all puzzles in the previous level
	prevLevel := level - 1
	prevPuzzles := puzzle.GetPuzzlesForLevel(allPuzzles, prevLevel)

	if len(prevPuzzles) == 0 {
		// No puzzles in previous level means this level is unlocked
		return true
	}

	// All puzzles in previous level must have at least 1 star
	for _, p := range prevPuzzles {
		result := s.GetBest(p.ID)
		if result.Stars < puzzle.OneStar {
			return false
		}
	}
	return true
}

// GetLevelStars returns the minimum star rating across all puzzles in a level.
func (s *Store) GetLevelStars(level int, allPuzzles []puzzle.Puzzle) puzzle.StarRating {
	puzzles := puzzle.GetPuzzlesForLevel(allPuzzles, level)
	if len(puzzles) == 0 {
		return puzzle.NoStar
	}

	minStars := puzzle.ThreeStar
	for _, p := range puzzles {
		result := s.GetBest(p.ID)
		if result.Stars < minStars {
			minStars = result.Stars
		}
	}
	return minStars
}

// GetTrackStars returns the best completed level star rating within a track.
func (s *Store) GetTrackStars(track int, allPuzzles []puzzle.Puzzle) puzzle.StarRating {
	levels := puzzle.GetLevelsForTrack(allPuzzles, track)
	if len(levels) == 0 {
		return puzzle.NoStar
	}

	maxStars := puzzle.NoStar
	for _, level := range levels {
		levelStars := s.GetLevelStars(level, allPuzzles)
		if levelStars > maxStars {
			maxStars = levelStars
		}
	}
	return maxStars
}

// OverallProgress returns solved count, total puzzles, and percent solved.
func (s *Store) OverallProgress(allPuzzles []puzzle.Puzzle) (int, int, int) {
	total := len(allPuzzles)
	if total == 0 {
		return 0, 0, 0
	}

	solved := 0
	for _, p := range allPuzzles {
		if s.GetBest(p.ID).Stars >= puzzle.OneStar {
			solved++
		}
	}

	percent := int(math.Round(float64(solved) * 100 / float64(total)))
	return solved, total, percent
}
