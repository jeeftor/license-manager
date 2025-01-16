package processor

import (
	"strings"

	"license-manager/internal/styles"
)

// LicenseChecker handles the detection and verification of license blocks
type LicenseChecker struct {
	style        styles.CommentStyle
	headerFooter styles.HeaderFooterStyle
}

// SearchPattern represents a pattern to look for in license headers
type SearchPattern struct {
	Start     string // Start of the license block
	End       string // End of the license block
	Priority  int    // Higher priority patterns should be checked first
	NeedsUTF8 bool   // Whether this pattern requires UTF-8 markers
}

// NewLicenseChecker creates a new LicenseChecker instance
func NewLicenseChecker(style styles.CommentStyle, headerFooter styles.HeaderFooterStyle) *LicenseChecker {
	return &LicenseChecker{
		style:        style,
		headerFooter: headerFooter,
	}
}

// generatePatterns creates all possible combinations of comment styles and markers
func (lc *LicenseChecker) generatePatterns() []SearchPattern {
	patterns := make([]SearchPattern, 0)
	priority := 1000 // Start with high priority

	// Helper to add patterns with proper priority
	addPattern := func(start, end string, needsUTF8 bool) {
		patterns = append(patterns, SearchPattern{
			Start:     start,
			End:       end,
			Priority:  priority,
			NeedsUTF8: needsUTF8,
		})
		priority-- // Each subsequent pattern has lower priority
	}

	// Unicode marker patterns (highest priority)
	if lc.style.MultiStart != "" {
		// Multi-line with unicode markers
		addPattern(
			lc.style.MultiStart+"\n"+lc.style.MultiPrefix+markerStart+lc.headerFooter.Header+markerEnd,
			lc.style.MultiPrefix+markerStart+lc.headerFooter.Footer+markerEnd+"\n"+lc.style.MultiEnd,
			true,
		)
		// Multi-line with unicode markers and space
		addPattern(
			lc.style.MultiStart+"\n"+lc.style.MultiPrefix+" "+markerStart+lc.headerFooter.Header+markerEnd,
			lc.style.MultiPrefix+" "+markerStart+lc.headerFooter.Footer+markerEnd+"\n"+lc.style.MultiEnd,
			true,
		)
	}

	if lc.style.Single != "" {
		// Single-line with unicode markers
		addPattern(
			lc.style.Single+" "+markerStart+lc.headerFooter.Header+markerEnd,
			lc.style.Single+" "+markerStart+lc.headerFooter.Footer+markerEnd,
			true,
		)
		// Single-line with unicode markers, no space
		addPattern(
			lc.style.Single+markerStart+lc.headerFooter.Header+markerEnd,
			lc.style.Single+markerStart+lc.headerFooter.Footer+markerEnd,
			true,
		)
	}

	// Fallback patterns without unicode markers (lower priority)
	if lc.style.MultiStart != "" {
		// Multi-line with header/footer only
		addPattern(
			lc.style.MultiStart+"\n"+lc.style.MultiPrefix+lc.headerFooter.Header,
			lc.style.MultiPrefix+lc.headerFooter.Footer+"\n"+lc.style.MultiEnd,
			false,
		)
		// Multi-line with comment markers only
		addPattern(
			lc.style.MultiStart,
			lc.style.MultiEnd,
			false,
		)
	}

	if lc.style.Single != "" {
		// Single-line with header/footer only
		addPattern(
			lc.style.Single+" "+lc.headerFooter.Header,
			lc.style.Single+" "+lc.headerFooter.Footer,
			false,
		)
		// Single-line with common license text
		addPattern(
			lc.style.Single+" Copyright",
			lc.style.Single+" SOFTWARE.",
			false,
		)
	}

	return patterns
}

// FindLicenseBlock attempts to find a license block in the content
func (lc *LicenseChecker) FindLicenseBlock(content string) (start, end int) {
	lines := strings.Split(content, "\n")
	patterns := lc.generatePatterns()

	// Try each pattern in priority order
	for _, pattern := range patterns {
		startLine, endLine := -1, -1

		// Find start
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, pattern.Start) {
				startLine = i
				break
			}
		}

		if startLine != -1 {
			// Found start, find the end
			for i := len(lines) - 1; i > startLine; i-- {
				line := strings.TrimSpace(lines[i])
				if strings.Contains(line, pattern.End) {
					endLine = i
					break
				}
			}

			if endLine != -1 {
				// Convert line numbers to character positions
				start = 0
				for i := 0; i < startLine; i++ {
					start += len(lines[i])
					if i < len(lines)-1 {
						start++ // +1 for newline
					}
				}

				end = 0
				for i := 0; i <= endLine; i++ {
					end += len(lines[i])
					if i < len(lines)-1 {
						end++ // +1 for newline
					}
				}
				return start, end
			}
		}
	}

	return -1, -1
}

// ParseLicenseBlock attempts to parse a license block from content
func ParseLicenseBlock(content string, style styles.CommentStyle) (*LicenseBlock, bool) {
	// Create a checker instance with empty style since we're just parsing
	checker := NewLicenseChecker(style, styles.HeaderFooterStyle{})

	// Find the license block boundaries
	start, end := checker.FindLicenseBlock(content)
	if start == -1 || end == -1 {
		return nil, false
	}

	// Extract the content between the boundaries
	blockContent := content[start:end]

	// Uncomment the content
	uncommented := uncommentContent(blockContent, style)

	// Split into lines
	lines := strings.Split(uncommented, "\n")
	if len(lines) < 3 { // Need at least header, body, footer
		return nil, false
	}

	// First non-empty line should be header with markers
	var header string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			header = stripMarkers(line)
			break
		}
	}

	// Last non-empty line should be footer with markers
	var footer string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			footer = stripMarkers(lines[i])
			break
		}
	}

	// Everything in between is the body
	bodyLines := lines[1 : len(lines)-1]
	// Trim any empty lines at start and end of body
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[0]) == "" {
		bodyLines = bodyLines[1:]
	}
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[len(bodyLines)-1]) == "" {
		bodyLines = bodyLines[:len(bodyLines)-1]
	}
	body := strings.Join(bodyLines, "\n")

	return &LicenseBlock{
		CommentStyle: style,
		Header:       header,
		Body:         body,
		Footer:       footer,
	}, true
}

// CheckLicense verifies if the content contains a valid license block
func (lc *LicenseChecker) CheckLicense(content string) bool {
	start, end := lc.FindLicenseBlock(content)
	return start != -1 && end != -1
}
