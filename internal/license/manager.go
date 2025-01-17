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

// Manager handles license operations
type Manager struct {
	template     string
	commentStyle styles.CommentLanguage
	headerStyle  styles.HeaderFooterStyle
}

// NewManager creates a new license manager
func NewManager(template string, headerStyle styles.HeaderFooterStyle) *Manager {
	return &Manager{
		template:    template,
		headerStyle: headerStyle,
	}
}

// SetCommentStyle sets the comment style for the manager
func (m *Manager) SetCommentStyle(style styles.CommentLanguage) {
	m.commentStyle = style
}

// HasLicense checks if content contains any license block
func (m *Manager) HasLicense(content string) bool {
	_, _, _, success := comment.ExtractComponents(content)
	return success
}

// AddLicense adds a license block to the content
func (m *Manager) AddLicense(content string) (string, error) {
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
func (m *Manager) RemoveLicense(content string) (string, error) {
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
func (m *Manager) UpdateLicense(content string) (string, error) {
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
func (m *Manager) CheckLicenseStatus(content string) Status {
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
func (m *Manager) GetLicenseComparison(content string) (current, expected string) {
	_, currentBody, _, _ := comment.ExtractComponents(content, true)
	_, expectedBody, _, _ := comment.ExtractComponents(m.formatLicenseBlock(m.template), true)
	return currentBody, expectedBody
}

// formatLicenseBlock formats the license text with the appropriate comment style
func (m *Manager) formatLicenseBlock(text string) string {
	return comment.FormatComment(text, m.commentStyle, m.headerStyle)
}
