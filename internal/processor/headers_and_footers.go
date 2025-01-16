package processor

import "strings"

type HeaderFooterStyle struct {
	Name        string
	Description string
	Header      string
	Footer      string
}

var PresetStyles = map[string]HeaderFooterStyle{
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
		Description: "Box style with corners",
		Header:      "+------------------------------------+",
		Footer:      "+------------------------------------+",
	},
	"brackets": {
		Name:        "Brackets",
		Description: "Square bracket style",
		Header:      "[ License Start ]----------------------",
		Footer:      "[ License End ]------------------------",
	},
	"stars": {
		Name:        "Stars",
		Description: "Star border style",
		Header:      "****************************************",
		Footer:      "****************************************",
	},
	"modern": {
		Name:        "Modern",
		Description: "Modern style with unicode blocks",
		Header:      "■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■",
		Footer:      "■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■",
	},
	"arrows": {
		Name:        "Arrows",
		Description: "Arrow style borders",
		Header:      ">>>>>>> LICENSE HEADER >>>>>>>>>>>>>>>>>>",
		Footer:      "<<<<<<<<<<<<<<<<<<<< LICENSE END <<<<<<<",
	},
	"elegant": {
		Name:        "Elegant",
		Description: "Elegant style with dashes and brackets",
		Header:      "---[ Begin License ]-------------------",
		Footer:      "---[ End License ]---------------------",
	},
	"wave": {
		Name:        "Wave",
		Description: "Wavy style border",
		Header:      "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		Footer:      "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
	},
	"dots": {
		Name:        "Dots",
		Description: "Dotted border style",
		Header:      "......................................",
		Footer:      "......................................",
	},
	"equals": {
		Name:        "Equals",
		Description: "Equal signs border",
		Header:      "======================================",
		Footer:      "======================================",
	},
	"decorative": {
		Name:        "Decorative",
		Description: "Decorative style with unicode symbols",
		Header:      "♦═══════════[ LICENSE ]══════════════♦",
		Footer:      "♦══════════[ END LICENSE ]═══════════♦",
	},
	"minimal": {
		Name:        "Minimal",
		Description: "Minimal style with thin lines",
		Header:      "─────────────────────────────────────",
		Footer:      "─────────────────────────────────────",
	},
	"classic": {
		Name:        "Classic",
		Description: "Classic style with asterisks and title",
		Header:      "**** BEGIN LICENSE BLOCK ****",
		Footer:      "**** END LICENSE BLOCK ****",
	},
	"angular": {
		Name:        "Angular",
		Description: "Angular style with forward slashes",
		Header:      "//////////// LICENSE START ////////////",
		Footer:      "//////////// LICENSE END //////////////",
	},
	"banner": {
		Name:        "Banner",
		Description: "Banner style with pipe borders",
		Header:      "|====================================|",
		Footer:      "|====================================|",
	},
	"retro": {
		Name:        "Retro",
		Description: "Retro style with plus signs",
		Header:      "++++++++++++++++++++++++++++++++++++",
		Footer:      "++++++++++++++++++++++++++++++++++++",
	},
	"clean": {
		Name:        "Clean",
		Description: "Clean style with triple dashes",
		Header:      "-----------------------------------",
		Footer:      "-----------------------------------",
	},
	"double": {
		Name:        "Double",
		Description: "Double line border",
		Header:      "══════════════════════════════════",
		Footer:      "══════════════════════════════════",
	},
	"branded": {
		Name:        "Branded",
		Description: "Corporate style with copyright symbol",
		Header:      "© ─────────[ LICENSE ]──────────── ©",
		Footer:      "© ────────[ END LICENSE ]────────── ©",
	},
}

// GetPresetStyle returns a HeaderFooterStyle by name, or a default if not found
func GetPresetStyle(name string) HeaderFooterStyle {
	if style, ok := PresetStyles[name]; ok {
		return style
	}
	return PresetStyles["simple"]
}

// ListPresetStyles returns a slice of all available style names
func ListPresetStyles() []string {
	var styles []string
	for name := range PresetStyles {
		styles = append(styles, name)
	}
	return styles
}

// HeaderAndFooterStyleMatch represents a matched style with confidence score
type HeaderAndFooterStyleMatch struct {
	Style    HeaderFooterStyle
	Score    float64
	IsHeader bool // true if matched against header
	IsFooter bool // true if matched against footer
}

// InferHeaderAndFooterStyle attempts to match a line against known header/footer patterns
// and returns the best matching style with confidence score
func InferHeaderAndFooterStyle(line string) HeaderAndFooterStyleMatch {
	var bestMatch HeaderAndFooterStyleMatch
	bestMatch.Score = 0.0

	cleanLine := strings.TrimSpace(line)
	if cleanLine == "" {
		return bestMatch
	}

	for _, style := range PresetStyles {
		// Try matching against header
		headerScore := calculateSimilarity(cleanLine, style.Header)
		if headerScore > bestMatch.Score {
			bestMatch = HeaderAndFooterStyleMatch{
				Style:    style,
				Score:    headerScore,
				IsHeader: true,
				IsFooter: false,
			}
		}

		// Try matching against footer
		footerScore := calculateSimilarity(cleanLine, style.Footer)
		if footerScore > bestMatch.Score {
			bestMatch = HeaderAndFooterStyleMatch{
				Style:    style,
				Score:    footerScore,
				IsHeader: false,
				IsFooter: true,
			}
		}
	}

	return bestMatch
}

// calculateSimilarity returns a similarity score between 0 and 1
// using a combination of exact matching and pattern matching
func calculateSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}

	// Remove all spaces and case sensitivity for base comparison
	cleanA := strings.ToLower(strings.ReplaceAll(a, " ", ""))
	cleanB := strings.ToLower(strings.ReplaceAll(b, " ", ""))

	if cleanA == cleanB {
		return 0.9
	}

	// Check for pattern matching by looking at distinct characters
	patternA := getDistinctPatternChars(cleanA)
	patternB := getDistinctPatternChars(cleanB)

	if patternA == patternB {
		return 0.8
	}

	// Calculate character overlap
	commonChars := 0
	for _, char := range patternA {
		if strings.ContainsRune(patternB, char) {
			commonChars++
		}
	}

	// Return a score based on character overlap
	overlap := float64(commonChars*2) / float64(len(patternA)+len(patternB))
	return overlap * 0.7 // Scale down pattern-only matches
}

// getDistinctPatternChars returns a string of unique pattern characters
func getDistinctPatternChars(s string) string {
	seen := make(map[rune]bool)
	var result strings.Builder

	for _, char := range s {
		// Only include pattern characters
		if isPatternChar(char) && !seen[char] {
			seen[char] = true
			result.WriteRune(char)
		}
	}

	return result.String()
}

// isPatternChar returns true for characters commonly used in patterns
func isPatternChar(r rune) bool {
	patterns := "-=~#*+[].><!|/♦©"
	return strings.ContainsRune(patterns, r)
}
