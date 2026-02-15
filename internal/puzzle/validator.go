package puzzle

import "strings"

// Validate checks if the current buffer text matches the target text.
func Validate(current, target string) bool {
	return strings.TrimRight(current, "\n") == strings.TrimRight(target, "\n")
}
