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
		handler:      GetLanguageHandler(commentStyle.FileType, style),
	}
}

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string, verbose bool) bool {
	if verbose {
		lm.debugLicenseCheck(content)
	}

	// Try to parse the license block
	expectedBlock := NewLicenseBlock(lm.style, lm.licenseText, lm.commentStyle)
	actualBlock, found := ParseLicenseBlock(content, lm.commentStyle)
	
	if !found {
		return false
	}

	// Compare the blocks
	return expectedBlock.String() == actualBlock.String()
}

// AddLicense adds the license text to the content
func (lm *LicenseManager) AddLicense(content string) string {
	// Check if license is already present
	if lm.CheckLicense(content, false) {
		return content
	}

	// Preserve any preamble
	preamble, rest := lm.handler.PreservePreamble(content)

	// Create and format the license block
	block := NewLicenseBlock(lm.style, lm.licenseText, lm.commentStyle)
	
	// Combine the parts
	var result []string
	if preamble != "" {
		result = append(result, preamble)
	}
	result = append(result, block.String())
	if rest != "" {
		result = append(result, "\n"+strings.TrimSpace(rest))
	}

	return strings.Join(result, "\n")
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// Preserve any language-specific preamble
	preamble, rest := lm.handler.PreservePreamble(content)

	// Find the license block in the rest of the content
	start, end := findLicenseBlock(rest)
	if start == -1 || end == -1 {
		return content
	}

	// Remove the license block and reconstruct the content
	var result []string
	if preamble != "" {
		result = append(result, preamble)
		result = append(result, "")
	}
	result = append(result, rest[:start]+rest[end:])

	return strings.Join(result, "\n")
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

// UpdateLicense updates the existing license text with new content
func (lm *LicenseManager) UpdateLicense(content string) string {
	if !lm.CheckLicense(content, false) {
		return content // No license found to update
	}

	// Remove the old license
	content = lm.RemoveLicense(content)

	// Add the new license
	return lm.AddLicense(content)
}

// GetLicenseComparison returns the current and expected license text for comparison
func (lm *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	current, found := lm.extractLicenseText(content)
	if !found {
		return "", ""
	}

	// Format the expected license text
	expectedContent := lm.AddLicense("")
	expected, _ = lm.extractLicenseText(expectedContent)

	return current, expected
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
	// Step 1: Look for invisible markers
	hasStart := strings.Contains(content, markerStart)
	hasEnd := strings.Contains(content, markerEnd)

	fmt.Printf("\n\033[1;34m=== License Check Debug Information ===\033[0m\n")

	// Check for invisible markers
	fmt.Printf("\033[1;36mInvisible Markers Check:\033[0m\n")
	fmt.Printf("  Found Start Marker (%s): %v\n", markerStart, hasStart)
	fmt.Printf("  Found End Marker (%s): %v\n", markerEnd, hasEnd)

	// If we found markers, show the lines containing them
	if hasStart || hasEnd {
		fmt.Printf("\n\033[1;36mLines containing markers:\033[0m\n")
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(line, markerStart) {
				fmt.Printf("  \033[32mLine %d: Start Marker in: %s\033[0m\n", i+1, line)
			}
			if strings.Contains(line, markerEnd) {
				fmt.Printf("  \033[32mLine %d: End Marker in: %s\033[0m\n", i+1, line)
			}
		}
	}

	// If we didn't find markers, try to infer them
	if !hasStart || !hasEnd {
		fmt.Printf("\n\033[1;36mAttempting to infer license block:\033[0m\n")
		// Check each preset style
		for name, style := range PresetStyles {
			headerWithoutMarkers := stripMarkers(style.Header)
			footerWithoutMarkers := stripMarkers(style.Footer)

			lines := strings.Split(content, "\n")
			for i, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if strings.Contains(trimmedLine, headerWithoutMarkers) {
					fmt.Printf("  \033[33mPossible header match (style: %s) at line %d: %s\033[0m\n",
						name, i+1, trimmedLine)
				}
				if strings.Contains(trimmedLine, footerWithoutMarkers) {
					fmt.Printf("  \033[33mPossible footer match (style: %s) at line %d: %s\033[0m\n",
						name, i+1, trimmedLine)
				}
			}
		}
	}

	// Extract and show the potential license text
	if hasStart && hasEnd {
		fmt.Printf("\n\033[1;36mExtracted License Text:\033[0m\n")
		if text, found := lm.extractLicenseText(content); found {
			fmt.Printf("  \033[32mFound license text (%d lines)\033[0m\n", len(strings.Split(text, "\n")))
		} else {
			fmt.Printf("  \033[31mCould not extract license text despite finding markers\033[0m\n")
		}
	}

	fmt.Printf("\n\033[1;34m=== End Debug Information ===\033[0m\n\n")
}

// uncommentContent removes comment markers from text while preserving the content
func uncommentContent(content string, style CommentStyle) string {
	// For single-line comments
	if style.Single != "" {
		var lines []string
		for _, line := range strings.Split(content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, style.Single) {
				// Remove the comment marker and any leading/trailing whitespace
				line = strings.TrimSpace(strings.TrimPrefix(trimmed, style.Single))
			}
			lines = append(lines, line)
		}
		content = strings.Join(lines, "\n")
	}

	// For multi-line comments
	if style.MultiStart != "" && style.MultiEnd != "" {
		// Remove multi-line comment markers
		content = strings.TrimSpace(content)
		if strings.HasPrefix(content, style.MultiStart) {
			content = strings.TrimSpace(strings.TrimPrefix(content, style.MultiStart))
		}
		if strings.HasSuffix(content, style.MultiEnd) {
			content = strings.TrimSpace(strings.TrimSuffix(content, style.MultiEnd))
		}

		// Remove asterisk prefixes common in multi-line comments
		var lines []string
		for _, line := range strings.Split(content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "*") {
				line = strings.TrimSpace(strings.TrimPrefix(trimmed, "*"))
			}
			lines = append(lines, line)
		}
		content = strings.Join(lines, "\n")
	}

	return content
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
