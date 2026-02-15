package puzzle

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		target   string
		expected bool
	}{
		{"exact match", "hello world", "hello world", true},
		{"trailing newline current", "hello world\n", "hello world", true},
		{"trailing newline target", "hello world", "hello world\n", true},
		{"no match", "hello", "world", false},
		{"empty strings", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.current, tt.target); got != tt.expected {
				t.Errorf("Validate(%q, %q) = %v, want %v", tt.current, tt.target, got, tt.expected)
			}
		})
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		name       string
		keystrokes int
		par        int
		expected   StarRating
	}{
		{"at par", 9, 9, ThreeStar},
		{"under par", 5, 9, ThreeStar},
		{"at 1.5x par", 13, 9, TwoStar},
		{"between par and 1.5x", 11, 9, TwoStar},
		{"over 1.5x par", 14, 9, OneStar},
		{"way over par", 100, 9, OneStar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Score(tt.keystrokes, tt.par); got != tt.expected {
				t.Errorf("Score(%d, %d) = %v, want %v", tt.keystrokes, tt.par, got, tt.expected)
			}
		})
	}
}
