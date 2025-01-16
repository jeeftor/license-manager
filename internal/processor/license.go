// internal/processor/license.go
package processor

import (
	"fmt"
	"strings"
	"unicode"
)

// LicenseBlock represents a complete license block with its components
type LicenseBlock struct {
	CommentStyle CommentStyle
	Header       string
	Body         string
	Footer       string
}

// String returns the complete license block as a string
func (lb *LicenseBlock) String() string {
	var result []string

	// Helper function to add markers if needed
	addMarkersIfNeeded := func(text string) string {
		if hasMarkers(text) {
			return text
		}
		return addMarkers(text)
	}

	if lb.CommentStyle.PreferMulti && lb.CommentStyle.MultiStart != "" {
		// Multi-line comment style
		result = append(result, lb.CommentStyle.MultiStart)
		result = append(result, " * "+addMarkersIfNeeded(lb.Header))
		
		// Add body with comment markers
		for _, line := range strings.Split(lb.Body, "\n") {
			if line == "" {
				result = append(result, " *")
			} else {
				result = append(result, " * "+line)
			}
		}
		
		result = append(result, " * "+addMarkersIfNeeded(lb.Footer))
		result = append(result, " "+lb.CommentStyle.MultiEnd)
	} else if lb.CommentStyle.Single != "" {
		// Single-line comment style
		result = append(result, lb.CommentStyle.Single+" "+addMarkersIfNeeded(lb.Header))
		
		// Add body with comment markers
		for _, line := range strings.Split(lb.Body, "\n") {
			if line == "" {
				result = append(result, lb.CommentStyle.Single)
			} else {
				result = append(result, lb.CommentStyle.Single+" "+line)
			}
		}
		
		result = append(result, lb.CommentStyle.Single+" "+addMarkersIfNeeded(lb.Footer))
	} else {
		// No comment style (e.g., for text files)
		result = append(result, addMarkersIfNeeded(lb.Header))
		result = append(result, lb.Body)
		result = append(result, addMarkersIfNeeded(lb.Footer))
	}

	return strings.Join(result, "\n")
}

// ParseLicenseBlock attempts to parse a license block from content
func ParseLicenseBlock(content string, style CommentStyle) (*LicenseBlock, bool) {
	// First uncomment the content if needed
	content = uncommentContent(content, style)

	// Look for markers
	start, end := findLicenseBlock(content)
	if start == -1 || end == -1 {
		return nil, false
	}

	// Split the content into lines
	lines := strings.Split(content[start:end], "\n")
	if len(lines) < 3 { // Need at least header, body, footer
		return nil, false
	}

	// Extract header and footer (they should have markers)
	header := stripMarkers(lines[0])
	footer := stripMarkers(lines[len(lines)-1])

	// Everything in between is the body
	body := strings.Join(lines[1:len(lines)-1], "\n")

	return &LicenseBlock{
		CommentStyle: style,
		Header:       header,
		Body:         body,
		Footer:       footer,
	}, true
}

// NewLicenseBlock creates a new license block with the given style and content
func NewLicenseBlock(style HeaderFooterStyle, licenseText string, commentStyle CommentStyle) *LicenseBlock {
	return &LicenseBlock{
		CommentStyle: commentStyle,
		Header:       style.Header,
		Body:         licenseText,
		Footer:       style.Footer,
	}
}

// LicenseManager handles license text operations
type LicenseManager struct {
	style        HeaderFooterStyle
	licenseText  string
	commentStyle CommentStyle
	handler      LanguageHandler
}

// NewLicenseManager creates a new instance of LicenseManager
func NewLicenseManager(style HeaderFooterStyle, licenseText string, commentStyle CommentStyle) *LicenseManager {
	return &LicenseManager{
		style:        style,
		licenseText:  licenseText,
		commentStyle: commentStyle,
		handler:      GetLanguageHandler(commentStyle.Language, style),
	}
}

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string, verbose bool) bool {
	if verbose {
		lm.debugLicenseCheck(content)
	}

	// Create the expected comment
	expected := NewComment(
		lm.commentStyle,
		lm.style.Header,
		lm.licenseText,
		lm.style.Footer,
	)

	// Try to parse the actual comment from content
	actual, found := Parse(content, lm.commentStyle)
	if !found {
		return false
	}

	// Compare the comments
	return expected.String() == actual.String()
}

// AddLicense adds the license text to the content
func (lm *LicenseManager) AddLicense(content string) string {
	// Check if license is already present
	if lm.CheckLicense(content, false) {
		return content
	}

	// Special handling for Go files
	if lm.commentStyle.Language == "go" {
		lines := strings.Split(content, "\n")
		directiveEnd := -1
		
		// Find the last build directive
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "// +build") || strings.HasPrefix(trimmed, "//go:") {
				directiveEnd = i
				// Skip blank line after directive
				if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
					directiveEnd = i + 1
				}
			}
		}

		// Create and format the license comment
		comment := NewComment(
			lm.commentStyle,
			lm.style.Header,
			lm.licenseText,
			lm.style.Footer,
		)
		licenseComment := comment.String()
		if !strings.HasSuffix(licenseComment, "\n") {
			licenseComment += "\n"
		}

		// Insert license after directives or at start
		if directiveEnd >= 0 {
			// Add after directives
			return strings.Join(lines[:directiveEnd+1], "\n") + "\n" + licenseComment + strings.Join(lines[directiveEnd+1:], "\n")
		} else {
			// Add at very start
			return licenseComment + content
		}
	}

	// For non-Go files, preserve any preamble
	preamble, rest := lm.handler.PreservePreamble(content)

	// Create and format the license comment
	comment := NewComment(
		lm.commentStyle,
		lm.style.Header,
		lm.licenseText,
		lm.style.Footer,
	)
	
	// Count newlines in original content after preamble
	preambleNewlines := 0
	if preamble != "" {
		// Find where preamble ends in original content
		preambleEnd := strings.Index(content, preamble) + len(preamble)
		for i := preambleEnd; i < len(content); i++ {
			if content[i] != '\n' {
				break
			}
			preambleNewlines++
		}
	}

	// Add one newline after license comment
	licenseComment := comment.String()
	if !strings.HasSuffix(licenseComment, "\n") {
		licenseComment += "\n"
	}

	// Build the result with proper newlines
	var result strings.Builder
	if preamble != "" {
		result.WriteString(preamble)
		result.WriteString(strings.Repeat("\n", preambleNewlines))
	}
	result.WriteString(licenseComment)
	if rest != "" {
		result.WriteString(rest)
	}

	return result.String()
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// Preserve any language-specific preamble
	preamble, rest := lm.handler.PreservePreamble(content)

	// Try to parse the comment
	comment, found := Parse(rest, lm.commentStyle)
	if !found {
		return content
	}

	// Remove the comment and reconstruct the content
	var result []string
	if preamble != "" {
		result = append(result, preamble)
		result = append(result, "")
	}

	// Get the content after the comment
	start := strings.Index(rest, comment.String())
	if start != -1 {
		end := start + len(comment.String())
		result = append(result, rest[:start]+rest[end:])
	}

	return strings.Join(result, "\n")
}

// UpdateLicense updates the existing license text with new content
func (lm *LicenseManager) UpdateLicense(content string) string {
	// Try to parse the existing comment
	comment, found := Parse(content, lm.commentStyle)
	if !found {
		return content // No license found to update
	}

	// Update the comment's body with the new license text
	comment.SetBody(lm.licenseText)

	// Remove the old comment and add the updated one
	content = lm.RemoveLicense(content)
	return lm.AddLicense(content)
}

// GetLicenseComparison returns the current and expected license text for comparison
func (lm *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	// Create the expected comment
	expectedComment := NewComment(
		lm.commentStyle,
		lm.style.Header,
		lm.licenseText,
		lm.style.Footer,
	)

	// Try to parse the actual comment
	actualComment, found := Parse(content, lm.commentStyle)
	if !found {
		return "", expectedComment.String()
	}

	return actualComment.String(), expectedComment.String()
}

// CheckLicenseStatus verifies the license status of the content
func (lm *LicenseManager) CheckLicenseStatus(content string) LicenseStatus {
	if lm.CheckLicense(content, false) {
		return MatchingLicense
	}
	return NoLicense
}

// LicenseStatus represents the status of a license check
type LicenseStatus int

const (
	NoLicense LicenseStatus = iota
	DifferentLicense
	MatchingLicense
)

// debugLicenseCheck performs detailed analysis of license markers in the content
func (lm *LicenseManager) debugLicenseCheck(content string) {
	fmt.Printf("\u001b[1;34m=== License Check Debug Information ===\u001b[0m\n")

	// Check for invisible markers
	fmt.Printf("\u001b[1;36mInvisible Markers Check:\u001b[0m\n")
	fmt.Printf("  Found Start Marker (%s): %v\n", markerStart, strings.Contains(content, markerStart))
	fmt.Printf("  Found End Marker (%s): %v\n", markerEnd, strings.Contains(content, markerEnd))

	// Try to parse the comment
	fmt.Printf("\n\u001b[1;36mAttempting to infer license block:\u001b[0m\n")
	if comment, found := Parse(content, lm.commentStyle); found {
		fmt.Printf("  Found license block:\n")
		fmt.Printf("    Header: %s\n", comment.Header)
		fmt.Printf("    Footer: %s\n", comment.Footer)
	} else {
		fmt.Printf("  No license block found\n")
	}

	fmt.Printf("\n\u001b[1;34m=== End Debug Information ===\u001b[0m\n")
}

// extractLicenseText extracts the license text between header and footer
func (lm *LicenseManager) extractLicenseText(content string) (string, bool) {
	// First check for the header and footer markers
	if !hasMarkers(content) {
		return "", false
	}

	start, end := findLicenseBlock(content)
	if start == -1 || end == -1 {
		return "", false
	}

	// Extract the text between the markers
	licenseBlock := content[start:end]
	licenseBlock = stripMarkers(licenseBlock)

	// Clean up the extracted text
	licenseBlock = uncommentContent(licenseBlock, lm.commentStyle)
	return strings.TrimSpace(licenseBlock), true
}

// normalizeText removes whitespace and comment markers from text
func (lm *LicenseManager) normalizeText(text string) string {
	// Convert to lowercase for case-insensitive comparison
	text = strings.ToLower(text)

	// Remove all whitespace
	var result strings.Builder
	for _, ch := range text {
		if !unicode.IsSpace(ch) {
			result.WriteRune(ch)
		}
	}

	// Remove comment markers
	text = result.String()
	if lm.commentStyle.Single != "" {
		text = strings.ReplaceAll(text, strings.ToLower(lm.commentStyle.Single), "")
	}
	if lm.commentStyle.MultiStart != "" {
		text = strings.ReplaceAll(text, strings.ToLower(lm.commentStyle.MultiStart), "")
		text = strings.ReplaceAll(text, strings.ToLower(lm.commentStyle.MultiEnd), "")
	}

	return text
}
