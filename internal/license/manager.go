package license

import (
	"license-manager/internal/comment"
	"license-manager/internal/errors"
	"license-manager/internal/styles"
	"strings"
)

// Status represents the status of a license in a file
type Status int

const (
	NoLicense Status = iota
	MatchingLicense
	DifferentLicense
)

// LicenseManager handles license operations
type LicenseManager struct {
	template     string
	commentStyle styles.CommentLanguage
	headerStyle  styles.HeaderFooterStyle
}

func NewLicenseManager(template string, headerStyle styles.HeaderFooterStyle, commentStyle ...styles.CommentLanguage) *LicenseManager {
	manager := &LicenseManager{
		template:    template,
		headerStyle: headerStyle,
	}
	if len(commentStyle) > 0 {
		manager.commentStyle = commentStyle[0]
	}
	return manager
}

// FormatLicenseForFile formats the license text with the current comment style
// This is useful for debugging and preview purposes
func (m *LicenseManager) FormatLicenseForFile(text string) string {
	if m.commentStyle.Language == "" {
		return "No comment style set - cannot format license"
	}
	return m.formatLicenseBlock(text)
}

// GetCurrentStyle returns the current comment style
// Useful for debugging and logging
func (m *LicenseManager) GetCurrentStyle() styles.CommentLanguage {
	return m.commentStyle
}

// SetCommentStyle sets the comment style for the manager
func (m *LicenseManager) SetCommentStyle(style styles.CommentLanguage) {
	m.commentStyle = style
}

// HasLicense checks if content contains any license block
func (m *LicenseManager) HasLicense(content string) bool {
	_, _, _, success := comment.ExtractComponents(content)
	return success
}

// AddLicense adds a license block to the content
func (m *LicenseManager) AddLicense(content string) (string, error) {
	if m.HasLicense(content) {
		return "", errors.NewLicenseError("content already has a license", "")
	}

	// Format the license block
	licenseBlock := comment.FormatComment(m.template, m.commentStyle, m.headerStyle)

	// If content is empty or only whitespace, just return the license
	if strings.TrimSpace(content) == "" {
		return licenseBlock, nil
	}

	// Add a newline between license and content
	return licenseBlock + "\n\n" + content, nil
}

// RemoveLicense removes the license block from the content
func (m *LicenseManager) RemoveLicense(content string) (string, error) {
	header, _, footer, success := comment.ExtractComponents(content)
	if !success {
		return content, nil
	}

	// Find the start and end of the license block
	lines := strings.Split(content, "\n")
	var startLine, endLine int
	for i, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(header) {
			startLine = i
			break
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == strings.TrimSpace(footer) {
			endLine = i
			break
		}
	}

	// Remove the license block and any surrounding empty lines
	result := append(lines[:startLine], lines[endLine+1:]...)
	for len(result) > 0 && strings.TrimSpace(result[0]) == "" {
		result = result[1:]
	}
	return strings.Join(result, "\n"), nil
}

// UpdateLicense updates the license block in the content
func (m *LicenseManager) UpdateLicense(content string) (string, error) {
	if !m.HasLicense(content) {
		return "", errors.NewLicenseError("content has no license to update", "")
	}

	// Remove the old license
	contentWithoutLicense, err := m.RemoveLicense(content)
	if err != nil {
		return "", err
	}

	// Add the new license
	return m.AddLicense(contentWithoutLicense)
}

// CheckLicenseStatus checks the status of the license in the content
func (m *LicenseManager) CheckLicenseStatus(content string) Status {
	if !m.HasLicense(content) {
		return NoLicense
	}

	_, body, _, _ := comment.ExtractComponents(content, true)
	_, expectedBody, _, _ := comment.ExtractComponents(m.formatLicenseBlock(m.template), true)

	if strings.TrimSpace(body) == strings.TrimSpace(expectedBody) {
		return MatchingLicense
	}
	return DifferentLicense
}

// GetLicenseComparison returns the current and expected license text for comparison
func (m *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	_, currentBody, _, _ := comment.ExtractComponents(content, true)
	_, expectedBody, _, _ := comment.ExtractComponents(m.formatLicenseBlock(m.template), true)
	return currentBody, expectedBody
}

// formatLicenseBlock formats the license text with the appropriate comment style
func (m *LicenseManager) formatLicenseBlock(text string) string {
	return comment.FormatComment(text, m.commentStyle, m.headerStyle)
}
