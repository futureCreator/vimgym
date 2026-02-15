package puzzle

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sort"
)

// LoadFromFile loads puzzles from a JSON file on disk.
func LoadFromFile(path string) ([]Puzzle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading puzzle file: %w", err)
	}
	var puzzles []Puzzle
	if err := json.Unmarshal(data, &puzzles); err != nil {
		return nil, fmt.Errorf("parsing puzzle file: %w", err)
	}
	return puzzles, nil
}

// LoadFromFS loads all puzzles from an fs.FS (e.g., embed.FS).
func LoadFromFS(fsys fs.FS, dir string) ([]Puzzle, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("reading puzzles dir: %w", err)
	}

	var all []Puzzle
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := entry.Name()
		if dir != "." {
			path = dir + "/" + path
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", entry.Name(), err)
		}
		var puzzles []Puzzle
		if err := json.Unmarshal(data, &puzzles); err != nil {
			return nil, fmt.Errorf("parsing file %s: %w", entry.Name(), err)
		}
		all = append(all, puzzles...)
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Track != all[j].Track {
			return all[i].Track < all[j].Track
		}
		return all[i].Level < all[j].Level
	})

	return all, nil
}

// GroupByLevel groups puzzles by their level number.
func GroupByLevel(puzzles []Puzzle) map[int][]Puzzle {
	m := make(map[int][]Puzzle)
	for _, p := range puzzles {
		m[p.Level] = append(m[p.Level], p)
	}
	return m
}

// GetTracks returns a sorted list of unique track numbers.
func GetTracks(puzzles []Puzzle) []int {
	seen := make(map[int]bool)
	for _, p := range puzzles {
		seen[p.Track] = true
	}
	tracks := make([]int, 0, len(seen))
	for t := range seen {
		tracks = append(tracks, t)
	}
	sort.Ints(tracks)
	return tracks
}

// GetLevelsForTrack returns sorted levels for a given track.
func GetLevelsForTrack(puzzles []Puzzle, track int) []int {
	seen := make(map[int]bool)
	for _, p := range puzzles {
		if p.Track == track {
			seen[p.Level] = true
		}
	}
	levels := make([]int, 0, len(seen))
	for l := range seen {
		levels = append(levels, l)
	}
	sort.Ints(levels)
	return levels
}

// GetPuzzlesForLevel returns puzzles for a specific level.
func GetPuzzlesForLevel(puzzles []Puzzle, level int) []Puzzle {
	var result []Puzzle
	for _, p := range puzzles {
		if p.Level == level {
			result = append(result, p)
		}
	}
	return result
}
