package styles

import "strings"

// HeaderFooterStyle represents a header/footer style for license comments
type HeaderFooterStyle struct {
	Name        string
	Description string
	Header      string
	Footer      string
}

// Match represents a matched style with confidence score
type Match struct {
	Style    HeaderFooterStyle
	Score    float64
	IsHeader bool // true if matched against header
	IsFooter bool // true if matched against footer
}

var presetStyles = map[string]HeaderFooterStyle{
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
	"swords": {
		Name:        "Swords",
		Description: "Swords border",
		Header:      "⚔️══✦══✦══ LICENSE ══✦══✦══⚔️",
		Footer:      "⚔️══✦══✦═ END LICENSE ═✦══✦══⚔️",
	},
	"scrolls": {
		Name:        "Scrolls",
		Description: "Scroll pattern border",
		Header:      "📜 ∽∽∽ LICENSE ∽∽∽ 📜",
		Footer:      "📜 ∽∽∽ END LICENSE ∽∽∽ 📜",
	},
	"waves": {
		Name:        "Waves",
		Description: "Wave pattern border",
		Header:      "〰️〰️〰️〰️ LICENSE BEGIN 〰️〰️〰️〰️",
		Footer:      "〰️〰️〰️〰️ LICENSE END 〰️〰️〰️〰️",
	},
}

// Get returns a HeaderFooterStyle by name, or a default if not found
func Get(name string) HeaderFooterStyle {
	if style, ok := presetStyles[name]; ok {
		return style
	}
	return presetStyles["hash"] // Return simple style as default
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

	// Clean up the line by removing common comment markers and spaces
	cleanLine := strings.TrimSpace(line)
	cleanLine = strings.TrimPrefix(cleanLine, "/*")
	cleanLine = strings.TrimSuffix(cleanLine, "*/")
	cleanLine = strings.TrimPrefix(cleanLine, "//")
	cleanLine = strings.TrimPrefix(cleanLine, "*")
	cleanLine = strings.TrimSpace(cleanLine)

	// Remove zero-width spaces used as markers
	cleanLine = strings.ReplaceAll(cleanLine, "​", "") // Zero-width space
	cleanLine = strings.ReplaceAll(cleanLine, "‌", "") // Zero-width non-joiner

	if cleanLine == "" {
		return bestMatch
	}

	for _, style := range presetStyles {
		// Clean up the style headers/footers the same way
		cleanHeader := strings.TrimSpace(style.Header)
		cleanFooter := strings.TrimSpace(style.Footer)

		// Try matching against header
		if cleanLine == cleanHeader {
			return Match{
				Style:    style,
				Score:    1.0,
				IsHeader: true,
				IsFooter: style.Header == style.Footer,
			}
		}

		// Try matching against footer
		if cleanLine == cleanFooter {
			return Match{
				Style:    style,
				Score:    1.0,
				IsHeader: style.Header == style.Footer,
				IsFooter: true,
			}
		}

		// If no exact match, try similarity matching
		headerScore := calculateSimilarity(cleanLine, cleanHeader)
		footerScore := calculateSimilarity(cleanLine, cleanFooter)

		if headerScore > bestMatch.Score {
			bestMatch = Match{
				Style:    style,
				Score:    headerScore,
				IsHeader: true,
				IsFooter: style.Header == style.Footer,
			}
		}

		if footerScore > bestMatch.Score {
			bestMatch = Match{
				Style:    style,
				Score:    footerScore,
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
