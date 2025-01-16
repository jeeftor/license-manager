package license

import (
	"fmt"
	"strings"

	"license-manager/internal/comment"
	"license-manager/internal/styles"
)

// LicenseManager handles license operations for a specific style and text
type LicenseManager struct {
	style        styles.HeaderFooterStyle
	licenseText  string
	commentStyle styles.CommentLanguage
	checker      *LicenseChecker
}

// NewLicenseManager creates a new LicenseManager instance
func NewLicenseManager(style styles.HeaderFooterStyle, licenseText string, commentStyle styles.CommentLanguage) *LicenseManager {
	return &LicenseManager{
		style:        style,
		licenseText:  licenseText,
		commentStyle: commentStyle,
		checker:      NewLicenseChecker(commentStyle, style),
	}
}

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string, verbose bool) bool {
	if verbose {
		lm.debugLicenseCheck(content)
	}

	// First check if there's a license block
	start, end := lm.checker.FindLicenseBlock(content)
	if start == -1 || end == -1 {
		return false
	}

	// Extract the license text and compare
	licenseBlock := content[start:end]
	licenseBlock = comment.uncommentContent(licenseBlock, lm.commentStyle)

	// Create the expected comment for comparison
	expected := comment.NewComment(
		lm.commentStyle,
		lm.style.Header,
		lm.licenseText,
		lm.style.Footer,
	)

	// Compare the uncommented text
	return strings.TrimSpace(licenseBlock) == strings.TrimSpace(expected.String())
}

// AddLicense adds the license text to the content
func (lm *LicenseManager) AddLicense(content string) string {
	// If there's already a license, return as is
	if lm.CheckLicense(content, false) {
		return content
	}

	// Create the comment with license
	comment := comment.NewComment(
		lm.commentStyle,
		lm.style.Header,
		lm.licenseText,
		lm.style.Footer,
	)

	// Add the comment to the content
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		return lm.addMultiLineComment(content, comment)
	}
	return lm.addSingleLineComment(content, comment)
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// Find the license block
	start, end := lm.checker.FindLicenseBlock(content)
	if start == -1 || end == -1 {
		return content
	}

	// Remove the license block
	return content[:start] + content[end:]
}

// ReplaceLicense replaces an existing license with a new one
func (lm *LicenseManager) ReplaceLicense(content string) string {
	// Remove existing license if present
	content = lm.RemoveLicense(content)

	// Add the new license
	return lm.AddLicense(content)
}

// UpdateLicense updates the existing license text with new content
func (lm *LicenseManager) UpdateLicense(content, newLicenseText string) string {
	// Save the current license text
	oldLicenseText := lm.licenseText

	// Set the new license text
	lm.licenseText = newLicenseText

	// Replace the license
	result := lm.ReplaceLicense(content)

	// Restore the original license text
	lm.licenseText = oldLicenseText

	return result
}

// GetLicenseComparison returns the current and expected license text for comparison
func (lm *LicenseManager) GetLicenseComparison(content string) (string, string) {
	start, end := lm.checker.FindLicenseBlock(content)
	if start == -1 || end == -1 {
		return "", lm.licenseText
	}

	// Extract the current license text
	currentLicense := content[start:end]
	currentLicense = comment.uncommentContent(currentLicense, lm.commentStyle)

	return currentLicense, lm.licenseText
}

// CheckLicenseStatus returns the status of the license in the content
func (lm *LicenseManager) CheckLicenseStatus(content string) LicenseStatus {
	if !lm.CheckLicense(content, false) {
		return NoLicense
	}

	current, expected := lm.GetLicenseComparison(content)
	if strings.TrimSpace(current) == strings.TrimSpace(expected) {
		return MatchingLicense
	}

	return DifferentLicense
}

// LicenseStatus represents the status of a license in a file
type LicenseStatus int

const (
	NoLicense LicenseStatus = iota
	DifferentLicense
	MatchingLicense
)

// debugLicenseCheck prints debug information about license checking
func (lm *LicenseManager) debugLicenseCheck(content string) {
	start, end := lm.checker.FindLicenseBlock(content)
	if start == -1 || end == -1 {
		fmt.Printf("No license block found in content\n")
		return
	}

	fmt.Printf("Found license block from position %d to %d\n", start, end)
	fmt.Printf("License block content:\n%s\n", content[start:end])
}

// addMultiLineComment adds a multi-line comment to the content
func (lm *LicenseManager) addMultiLineComment(content string, comment *comment.Comment) string {
	var result strings.Builder

	// Add comment start
	result.WriteString(lm.commentStyle.MultiStart)
	result.WriteString("\n")

	// Add header with prefix
	if comment.Header != "" {
		result.WriteString(lm.commentStyle.MultiPrefix)
		result.WriteString(" ")
		result.WriteString(comment.addMarkers(comment.Header))
		result.WriteString("\n")
	}

	// Add body with prefix for each line
	lines := strings.Split(comment.Body, "\n")
	for _, line := range lines {
		if line != "" {
			result.WriteString(lm.commentStyle.MultiPrefix)
			result.WriteString(lm.commentStyle.LinePrefix)
			result.WriteString(line)
		}
		result.WriteString("\n")
	}

	// Add footer with prefix
	if comment.Footer != "" {
		result.WriteString(lm.commentStyle.MultiPrefix)
		result.WriteString(" ")
		result.WriteString(comment.addMarkers(comment.Footer))
		result.WriteString("\n")
	}

	// Add comment end
	result.WriteString(lm.commentStyle.MultiEnd)
	result.WriteString("\n\n")

	// Add the original content
	result.WriteString(content)

	return result.String()
}

// addSingleLineComment adds a single-line comment to the content
func (lm *LicenseManager) addSingleLineComment(content string, comment *comment.Comment) string {
	var result strings.Builder

	// Add header
	if comment.Header != "" {
		result.WriteString(lm.commentStyle.Single)
		result.WriteString(lm.commentStyle.LinePrefix)
		result.WriteString(comment.addMarkers(comment.Header))
		result.WriteString("\n")
	}

	// Add body
	lines := strings.Split(comment.Body, "\n")
	for _, line := range lines {
		result.WriteString(lm.commentStyle.Single)
		result.WriteString(lm.commentStyle.LinePrefix)
		result.WriteString(line)
		result.WriteString("\n")
	}

	// Add footer
	if comment.Footer != "" {
		result.WriteString(lm.commentStyle.Single)
		result.WriteString(lm.commentStyle.LinePrefix)
		result.WriteString(comment.addMarkers(comment.Footer))
		result.WriteString("\n")
	}

	// Add a blank line
	result.WriteString("\n")

	// Add the original content
	result.WriteString(content)

	return result.String()
}
