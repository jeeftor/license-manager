package license

import (
	"license-manager/internal/comment"
	"license-manager/internal/errors"
	"license-manager/internal/language"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// Status represents the status of a license in a file
type Status int

const (
	MatchingLicense Status = iota
	NoLicense
	DifferentLicense
)

// LicenseManager handles license operations
type LicenseManager struct {
	template     string
	commentStyle styles.CommentLanguage
	headerStyle  styles.HeaderFooterStyle
	langHandler  language.LanguageHandler
	verbose      bool
	logger       *logger.Logger
}

// NewLicenseManager creates a new manager
func NewLicenseManager(template string, headerStyle styles.HeaderFooterStyle, commentStyle ...styles.CommentLanguage) *LicenseManager {
	manager := &LicenseManager{
		template:    template,
		headerStyle: headerStyle,
		verbose:     false,
	}
	if len(commentStyle) > 0 {
		manager.commentStyle = commentStyle[0]
	}
	return manager
}

// SetVerbose enables verbose logging
func (m *LicenseManager) SetVerbose(verbose bool, logger *logger.Logger) {
	m.verbose = verbose
	m.logger = logger
	if m.langHandler != nil {
		m.langHandler.SetLogger(logger)
	}
}

// SetCommentStyle sets the comment style
func (m *LicenseManager) SetCommentStyle(style styles.CommentLanguage) {
	m.commentStyle = style
}

// SetLanguageHandler sets the language handler
func (m *LicenseManager) SetLanguageHandler(handler language.LanguageHandler) {
	m.langHandler = handler
	if m.verbose && m.logger != nil {
		handler.SetLogger(m.logger)
	}
}

// getLanguageHandler returns a configured language handler
func (m *LicenseManager) getLanguageHandler(fileType string) language.LanguageHandler {
	if m.langHandler != nil {
		return m.langHandler
	}
	// Fallback to creating a new handler if none is set
	handler := language.GetLanguageHandler(m.commentStyle.Language, m.headerStyle)
	if m.verbose && m.logger != nil {
		handler.SetLogger(m.logger)
	}
	return handler
}

// HasLicense checks if content contains any license block
func (m *LicenseManager) HasLicense(content string) bool {
	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Checking for existing license block...")
		m.logger.LogInfo("  Using comment style: %s", m.commentStyle.Language)
	}

	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.verbose && m.logger != nil {
		if preamble != "" {
			m.logger.LogInfo("  Found preamble (%d lines)", len(strings.Split(preamble, "\n")))
			m.logger.LogInfo("  Processing remaining content after preamble...")
		} else {
			m.logger.LogInfo("  No preamble found, processing entire content")
		}
	}

	_, body, _, success := comment.ExtractComponents(rest)
	if m.verbose && m.logger != nil {
		if success {
			m.logger.LogInfo("  Found existing license block (%d lines)", len(strings.Split(body, "\n")))
			m.logger.LogInfo("  License content preview: %s...", truncateString(body, 50))
		} else {
			m.logger.LogInfo("  No license block detected in content")
		}
	}
	return success
}

// AddLicense adds a license block to the content
func (m *LicenseManager) AddLicense(content string, fileType string) (string, error) {
	if m.verbose && m.logger != nil {
		m.logger.LogVerbose("Adding license to content...")
	}

	// Get language handler
	handler := m.getLanguageHandler(fileType)

	// Extract preamble (e.g., shebang, package declaration)
	preamble, rest := handler.PreservePreamble(content)
	if m.verbose && m.logger != nil {
		if preamble != "" {
			m.logger.LogVerbose("Found preamble: %s", truncateString(preamble, 50))
		}
	}

	// Check for existing license
	if m.HasLicense(rest) {
		return "", errors.NewLicenseError("license already exists", "add")
	}

	// Format license block with comment style
	licenseBlock := m.formatLicenseBlock(m.template)

	// Scan for build directives
	directives, endIndex := handler.ScanBuildDirectives(rest)
	if m.verbose && m.logger != nil {
		if len(directives) > 0 {
			m.logger.LogVerbose("Found %d build directives", len(directives))
		}
	}

	// Split content at build directives
	var beforeDirectives, afterDirectives string
	if len(directives) > 0 {
		lines := strings.Split(rest, "\n")
		beforeDirectives = strings.Join(lines[:endIndex], "\n")
		if endIndex < len(lines) {
			afterDirectives = strings.Join(lines[endIndex:], "\n")
		}
	} else {
		afterDirectives = rest
	}

	// Build final content
	var parts []string
	if preamble != "" {
		parts = append(parts, preamble)
	}
	if len(directives) > 0 {
		parts = append(parts, beforeDirectives)
	}
	parts = append(parts, licenseBlock)
	if afterDirectives != "" {
		if !strings.HasPrefix(afterDirectives, "\n") {
			parts = append(parts, "") // Add blank line before content
		}
		parts = append(parts, strings.TrimPrefix(afterDirectives, "\n"))
	}

	return strings.Join(parts, "\n"), nil
}

// RemoveLicense removes the license block from the content
func (m *LicenseManager) RemoveLicense(content string) (string, error) {
	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Attempting to remove license block...")
		m.logger.LogInfo("  Using comment style: %s", m.commentStyle.Language)
	}

	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.verbose && m.logger != nil {
		if preamble != "" {
			m.logger.LogInfo("  Found preamble (%d lines)", len(strings.Split(preamble, "\n")))
			m.logger.LogInfo("  Processing remaining content...")
		} else {
			m.logger.LogInfo("  No preamble found")
		}
	}

	header, body, footer, success := comment.ExtractComponents(rest)
	if !success {
		if m.verbose && m.logger != nil {
			m.logger.LogInfo("  No license block found to remove")
		}
		return content, nil
	}

	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Found license block:")
		m.logger.LogInfo("    Header: %s", truncateString(header, 50))
		m.logger.LogInfo("    Body length: %d lines", len(strings.Split(body, "\n")))
		m.logger.LogInfo("    Footer: %s", truncateString(footer, 50))
	}

	// Find the start and end of the license block
	lines := strings.Split(rest, "\n")
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

	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  License block found at lines %d-%d", startLine, endLine)
	}

	// Remove the license block and any surrounding empty lines
	result := append(lines[:startLine], lines[endLine+1:]...)
	for len(result) > 0 && strings.TrimSpace(result[0]) == "" {
		result = result[1:]
	}

	newContent := strings.Join(result, "\n")
	if preamble != "" {
		if m.verbose && m.logger != nil {
			m.logger.LogInfo("  Reconstructing content with preamble")
		}
		return preamble + "\n" + newContent, nil
	}
	return newContent, nil
}

// UpdateLicense updates the license block in the content
func (m *LicenseManager) UpdateLicense(content string) (string, error) {
	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.verbose && m.logger != nil {
		if preamble != "" {
			m.logger.LogInfo("  Found preamble when updating license")
		}
	}

	if !m.HasLicense(rest) {
		return "", errors.NewLicenseError("content has no license to update", "")
	}

	// Remove the old license
	contentWithoutLicense, err := m.RemoveLicense(content)
	if err != nil {
		return "", err
	}

	// Add the new license
	return m.AddLicense(contentWithoutLicense, m.commentStyle.Language)
}

// CheckLicenseStatus checks the status of the license in the content
func (m *LicenseManager) CheckLicenseStatus(content string) Status {
	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Checking license status...")
		m.logger.LogInfo("  Using comment style: %s", m.commentStyle.Language)
	}

	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.verbose && m.logger != nil {
		if preamble != "" {
			m.logger.LogInfo("  Found preamble (%d lines)", len(strings.Split(preamble, "\n")))
			m.logger.LogInfo("  Analyzing content after preamble...")
		} else {
			m.logger.LogInfo("  No preamble found, analyzing full content")
		}
	}

	if !m.HasLicense(rest) {
		if m.verbose && m.logger != nil {
			m.logger.LogInfo("  Result: No license block found")
		}
		return NoLicense
	}

	_, body, _, _ := comment.ExtractComponents(rest, true)
	_, expectedBody, _, _ := comment.ExtractComponents(m.formatLicenseBlock(m.template), true)

	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Found license block (%d lines)", len(strings.Split(body, "\n")))
		m.logger.LogInfo("  Expected license (%d lines)", len(strings.Split(expectedBody, "\n")))
		m.logger.LogInfo("  Comparing license content...")
	}

	if strings.TrimSpace(body) == strings.TrimSpace(expectedBody) {
		if m.verbose && m.logger != nil {
			m.logger.LogInfo("  Result: License matches expected content exactly")
		}
		return MatchingLicense
	}

	if m.verbose && m.logger != nil {
		m.logger.LogInfo("  Result: License differs from expected content")
	}
	return DifferentLicense
}

// GetLicenseComparison returns the current and expected license text for comparison
func (m *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	handler := m.getLanguageHandler(m.commentStyle.Language)
	_, rest := handler.PreservePreamble(content)

	_, currentBody, _, _ := comment.ExtractComponents(rest, true)
	_, expectedBody, _, _ := comment.ExtractComponents(m.formatLicenseBlock(m.template), true)
	return currentBody, expectedBody
}

// FormatLicenseForFile formats the license text with the current comment style
func (m *LicenseManager) FormatLicenseForFile(text string) string {
	if m.commentStyle.Language == "" {
		return "No comment style set - cannot format license"
	}
	return m.formatLicenseBlock(text)
}

// formatLicenseBlock formats the license text with the appropriate comment style
func (m *LicenseManager) formatLicenseBlock(text string) string {
	return comment.FormatComment(text, m.commentStyle, m.headerStyle)
}

// helper function to truncate strings for logging
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
