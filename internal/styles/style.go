package styles

import "strings"

// Style represents a header/footer style for license comments
type Style struct {
	Name        string
	Description string
	Header      string
	Footer      string
}

// Match represents a matched style with confidence score
type Match struct {
	Style    Style
	Score    float64
	IsHeader bool // true if matched against header
	IsFooter bool // true if matched against footer
}

var presetStyles = map[string]Style{
	"simple": {
		Name:        "Simple",
		Description: "Simple lines before and after",
		Header:      "----------------------------------------",
		Footer:      "----------------------------------------",
	},
	"hash": {
		Name:        "Hash",
		Description: "Hash symbol borders",
		Header:      "######################################",
		Footer:      "######################################",
	},
	"box": {
		Name:        "Box",
		Description: "Box with plus signs",
		Header:      "+------------------------------------+",
		Footer:      "+------------------------------------+",
	},
	"brackets": {
		Name:        "Brackets",
		Description: "Brackets with descriptive text",
		Header:      "[ License Start ]----------------------",
		Footer:      "[ License End ]------------------------",
	},
	"equals": {
		Name:        "Equals",
		Description: "Equal signs border",
		Header:      "========================================",
		Footer:      "========================================",
	},
	"stars": {
		Name:        "Stars",
		Description: "Asterisk border",
		Header:      "****************************************",
		Footer:      "****************************************",
	},
	"slashes": {
		Name:        "Slashes",
		Description: "Forward slash border",
		Header:      "////////////////////////////////////////",
		Footer:      "////////////////////////////////////////",
	},
	"dots": {
		Name:        "Dots",
		Description: "Dotted border",
		Header:      "........................................",
		Footer:      "........................................",
	},
}

// Get returns a Style by name, or a default if not found
func Get(name string) Style {
	if style, ok := presetStyles[name]; ok {
		return style
	}
	return presetStyles["simple"] // Return simple style as default
}

// List returns a slice of all available style names
func List() []string {
	var names []string
	for name := range presetStyles {
		names = append(names, name)
	}
	return names
}

// Infer attempts to match a line against known header/footer patterns
// and returns the best matching style with confidence score
func Infer(line string) Match {
	var bestMatch Match
	bestMatch.Score = 0.0

	cleanLine := strings.TrimSpace(line)
	if cleanLine == "" {
		return bestMatch
	}

	for _, style := range presetStyles {
		// Try matching against header
		if cleanLine == strings.TrimSpace(style.Header) {
			return Match{
				Style:    style,
				Score:    1.0,
				IsHeader: true,
				IsFooter: style.Header == style.Footer,
			}
		}

		// Try matching against footer
		if cleanLine == strings.TrimSpace(style.Footer) {
			return Match{
				Style:    style,
				Score:    1.0,
				IsHeader: style.Header == style.Footer,
				IsFooter: true,
			}
		}
	}

	return bestMatch
}

// calculateSimilarity returns a similarity score between 0 and 1
// where 1 means exact match and 0 means no match
func calculateSimilarity(input, pattern string) float64 {
	input = strings.TrimSpace(strings.ToLower(input))
	pattern = strings.TrimSpace(strings.ToLower(pattern))

	if input == "" || pattern == "" {
		return 0.0
	}

	if input == pattern {
		return 1.0
	}

	return 0.0
}

// isRepeatedPattern checks if a string consists of a single character repeated
func isRepeatedPattern(s string) bool {
	if len(s) == 0 {
		return false
	}
	char := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != char && s[i] != ' ' {
			return false
		}
	}
	return true
}

// getRepeatedChar returns the character that is repeated in the pattern
func getRepeatedChar(s string) byte {
	if len(s) == 0 {
		return 0
	}
	return s[0]
}

// getDistinctPatternChars returns a string of unique pattern characters
func getDistinctPatternChars(s string) string {
	seen := make(map[rune]bool)
	var result []rune

	for _, r := range s {
		if isPatternChar(r) && !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}

	return string(result)
}

// isPatternChar returns true for characters commonly used in patterns
func isPatternChar(r rune) bool {
	switch r {
	case '-', '=', '*', '#', '+', '/', '.', '_', '~':
		return true
	default:
		return false
	}
}
