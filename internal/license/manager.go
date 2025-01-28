package license

import (
	"license-manager/internal/comment"
	"license-manager/internal/errors"
	"license-manager/internal/language"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// Status represents the status of a license check
type Status int

const (
	// FullMatch indicates that both the license content and style match
	FullMatch Status = iota
	// NoLicense indicates that no license was found
	NoLicense
	// ContentAndStyleMismatch indicates that both the content and style are different
	ContentAndStyleMismatch
	// ContentMismatch indicates that the license content is different
	ContentMismatch
	// StyleMismatch indicates that the license content matches but the style is different
	StyleMismatch
)

func (s Status) String() string {
	switch s {
	case FullMatch:
		return "License OK"
	case NoLicense:
		return "No license found"
	case ContentMismatch:
		return "License content mismatch"
	case StyleMismatch:
		return "License style mismatch"
	case ContentAndStyleMismatch:
		return "License content and style mismatch"
	default:
		return "Unknown status"
	}
}

// LicenseManager handles license operations
type LicenseManager struct {
	template     string
	commentStyle styles.CommentLanguage
	headerStyle  styles.HeaderFooterStyle
	langHandler  language.LanguageHandler
	logger       *logger.Logger
}

// NewLicenseManager creates a new manager
func NewLicenseManager(logger *logger.Logger, template string, headerStyle styles.HeaderFooterStyle, commentStyle ...styles.CommentLanguage) *LicenseManager {
	manager := &LicenseManager{
		template:    template,
		headerStyle: headerStyle,
		logger:      logger,
	}
	if len(commentStyle) > 0 {
		manager.commentStyle = commentStyle[0]
	}
	return manager
}

// SetCommentStyle sets the comment style
func (m *LicenseManager) SetCommentStyle(style styles.CommentLanguage) {
	m.commentStyle = style
}

// SetLanguageHandler sets the language handler
func (m *LicenseManager) SetLanguageHandler(handler language.LanguageHandler) {
	m.langHandler = handler

}

// getLanguageHandler returns a configured language handler
func (m *LicenseManager) getLanguageHandler(fileType string) language.LanguageHandler {
	if m.langHandler != nil {
		return m.langHandler
	}
	// Fallback to creating a new handler if none is set
	handler := language.GetLanguageHandler(m.logger, m.commentStyle.Language, m.headerStyle)

	return handler
}

// HasLicense checks if content contains any license block
func (m *LicenseManager) HasLicense(content string) bool {

	m.logger.LogInfo("  Checking for existing license block...")
	m.logger.LogInfo("  Using comment style: %s", m.commentStyle.Language)

	h, body, f, success := comment.ExtractComponents(m.logger, content, true, m.commentStyle)
	if h == f {
		m.logger.LogInfo("  Header & Footer Match")

	}

	if success {
		m.logger.LogInfo("  Found existing license block (%d lines)", len(strings.Split(body, "\n")))
		m.logger.LogInfo("  License content preview: %s...", truncateString(body, 50))
	} else {
		m.logger.LogInfo("  No license block detected in content")
	}

	return success
}

// AddLicense adds a license block to the content
func (m *LicenseManager) AddLicense(content string, fileType string) (string, error) {
	if m.logger != nil {
		m.logger.LogVerbose("Adding license to content...")
	}

	// Get language handler
	handler := m.getLanguageHandler(fileType)

	// Extract preamble (e.g., shebang, package declaration)
	preamble, rest := handler.PreservePreamble(content)
	if m.logger != nil {
		if preamble != "" {
			m.logger.LogVerbose("Found preamble: %s", truncateString(preamble, 50))
		}
	}

	// Check for existing license
	if m.HasLicense(rest) {
		// If we're updating, allow the license to be replaced
		if status := m.CheckLicenseStatus(rest); status == FullMatch {
			return "", errors.NewLicenseError("license already exists and matches", "add")
		}
	}

	// Format license block with comment style
	licenseBlock := m.formatLicenseBlock(m.template)

	// Scan for build directives
	directives, endIndex := handler.ScanBuildDirectives(rest)
	if m.logger != nil {
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
	m.logger.LogDebug("  Attempting to remove license block...")
	m.logger.LogDebug("  Using comment style: %s", m.commentStyle.Language)

	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.logger != nil {
		if preamble != "" {
			m.logger.LogDebug("  Found preamble (%d lines)", len(strings.Split(preamble, "\n")))
			m.logger.LogDebug("  Processing remaining content...")
		} else {
			m.logger.LogDebug("  No preamble found")
		}
	}

	header, body, footer, success := comment.ExtractComponents(m.logger, rest, true, m.commentStyle)
	if !success {
		if m.logger != nil {
			m.logger.LogDebug("  No license block found to remove")
		}
		return content, nil
	}

	m.logger.LogDebug("  Found license block:")
	m.logger.LogDebug("    Header: %s", truncateString(header, 50))
	m.logger.LogDebug("    Body length: %d lines", len(strings.Split(body, "\n")))
	m.logger.LogDebug("    Footer: %s", truncateString(footer, 50))

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

	if m.logger != nil {
		m.logger.LogInfo("  License block found at lines %d-%d", startLine, endLine)
	}

	// Remove the license block and any surrounding empty lines
	result := append(lines[:startLine], lines[endLine+1:]...)

	// Clean up empty comment blocks
	var cleanedResult []string
	inEmptyComment := false
	for _, line := range result {
		trimmed := strings.TrimSpace(line)
		if trimmed == "/*" {
			inEmptyComment = true
			continue
		}
		if trimmed == "*/" && inEmptyComment {
			inEmptyComment = false
			continue
		}
		if !inEmptyComment && trimmed != "" {
			cleanedResult = append(cleanedResult, line)
		}
	}

	// Remove leading empty lines
	for len(cleanedResult) > 0 && strings.TrimSpace(cleanedResult[0]) == "" {
		cleanedResult = cleanedResult[1:]
	}

	newContent := strings.Join(cleanedResult, "\n")
	if preamble != "" {
		if m.logger != nil {
			m.logger.LogInfo("  Reconstructing content with preamble")
		}
		return preamble + "\n" + newContent, nil
	}
	return newContent, nil
}

// UpdateLicense updates the license block in the content
func (m *LicenseManager) UpdateLicense(content string) (string, error) {
	handler := m.getLanguageHandler(m.commentStyle.Language)
	_, rest := handler.PreservePreamble(content)

	if m.logger != nil {
		if rest != "" {
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
	if m.logger != nil {
		m.logger.LogInfo("  Checking license status...")
		m.logger.LogInfo("  Using comment style: %s", m.commentStyle.Language)
	}

	handler := m.getLanguageHandler(m.commentStyle.Language)
	preamble, rest := handler.PreservePreamble(content)

	if m.logger != nil {
		if preamble != "" {
			m.logger.LogInfo("  Found preamble 📝️ (%d lines)", len(strings.Split(preamble, "\n")))
			m.logger.LogInfo("  Analyzing content after preamble...")
		} else {
			m.logger.LogInfo("  No preamble found, analyzing full content")
		}
	}

	if !m.HasLicense(rest) {
		if m.logger != nil {
			m.logger.LogInfo("  Result: No license block found")
		}
		return NoLicense
	}

	// Get the current header/footer from the content
	detectedStyle := m.DetectHeaderStyle(rest)
	if m.logger != nil {
		m.logger.LogInfo("  Detected style: %s", detectedStyle.Name)
	}

	// If a specific style was requested, check if it matches
	if m.headerStyle.Name != "" && m.headerStyle.Name != detectedStyle.Name {
		if m.logger != nil {
			m.logger.LogInfo("  Style mismatch: expected [%s], found [%s]", m.headerStyle.Name, detectedStyle.Name)
		}

		// Check if content matches despite style mismatch
		currentLicense := rest
		expectedLicense := handler.FormatLicense(m.template, m.commentStyle, m.headerStyle)

		// Extract both licenses without stripping markers for accurate comparison
		_, currentBody, _, _ := comment.ExtractComponents(m.logger, currentLicense, false, m.commentStyle)
		_, expectedBody, _, _ := comment.ExtractComponents(m.logger, expectedLicense, false, m.commentStyle)

		if currentBody == expectedBody {
			return StyleMismatch
		}
		return ContentAndStyleMismatch
	}

	// Compare the license content
	currentLicense := rest
	expectedLicense := handler.FormatLicense(m.template, m.commentStyle, detectedStyle)

	// Extract both licenses without stripping markers for accurate comparison
	_, currentBody, _, _ := comment.ExtractComponents(m.logger, currentLicense, true, m.commentStyle)
	_, expectedBody, _, _ := comment.ExtractComponents(m.logger, expectedLicense, true, m.commentStyle)

	if currentBody == expectedBody {
		if m.logger != nil {
			m.logger.LogInfo("  Result: License matches expected content exactly")
		}
		return FullMatch
	}

	if m.logger != nil {
		m.logger.LogInfo("  Result: License content differs from expected")

		lines1 := strings.Split(expectedBody, "\n")
		lines2 := strings.Split(currentBody, "\n")

		m.logger.LogInfo("\nDiff: \033[31m expectedBody \033[0m\033[32m currentBody\033[0m")

		for i := 0; i < max(len(lines1), len(lines2)); i++ {
			line1 := ""
			line2 := ""
			if i < len(lines1) {
				line1 = lines1[i]
			}
			if i < len(lines2) {
				line2 = lines2[i]
			}

			if line1 != line2 {
				m.logger.LogInfo("\033[31m- %s\033[0m", line1)
				m.logger.LogInfo("\033[32m+ %s\033[0m", line2)
			} else {
				m.logger.LogInfo("  %s", line1)
			}
		}
	}
	return ContentMismatch
}

// GetLicenseComparison returns the current and expected license text for comparison
func (m *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	handler := m.getLanguageHandler(m.commentStyle.Language)
	_, rest := handler.PreservePreamble(content)

	// First detect the style from the content
	detectedStyle := m.DetectHeaderStyle(rest)
	if m.logger != nil {
		m.logger.LogInfo("  Detected style: %s", detectedStyle.Name)
	}

	// Format both licenses with the detected style
	currentLicense := rest
	expectedLicense := handler.FormatLicense(m.template, m.commentStyle, detectedStyle)

	// Extract both licenses without stripping markers for accurate comparison
	_, currentBody, _, _ := comment.ExtractComponents(m.logger, currentLicense, false, m.commentStyle)
	_, expectedBody, _, _ := comment.ExtractComponents(m.logger, expectedLicense, false, m.commentStyle)
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

// DetectHeaderStyle detects the header/footer style from the content
func (m *LicenseManager) DetectHeaderStyle(content string) styles.HeaderFooterStyle {
	handler := m.getLanguageHandler(m.commentStyle.Language)
	_, rest := handler.PreservePreamble(content)

	if m.logger != nil {
		m.logger.LogInfo("  Content to analyze:\n%s", rest)
	}

	// Get the current header/footer from the content
	header, _, footer, success := comment.ExtractComponents(m.logger, rest, true, m.commentStyle)
	if !success {
		return m.headerStyle // Return current style if extraction fails
	}

	// Split into lines and find the actual header/footer lines
	headerLines := strings.Split(header, "\n")
	footerLines := strings.Split(footer, "\n")

	var firstLine, lastLine string
	for _, line := range headerLines {
		line = strings.TrimSpace(line)
		if line == "" || line == "/*" || line == "*/" {
			continue
		}
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, " *")
		line = strings.TrimSpace(line)
		if line != "" {
			firstLine = line
			break
		}
	}

	for i := len(footerLines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(footerLines[i])
		if line == "" || line == "/*" || line == "*/" {
			continue
		}
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, " *")
		line = strings.TrimSpace(line)
		if line != "" {
			lastLine = line
			break
		}
	}

	if m.logger != nil {
		m.logger.LogInfo("  Trying to detect style from header: %q", firstLine)
		m.logger.LogInfo("  Trying to detect style from footer: %q", lastLine)
	}

	// Try to match against known styles
	headerMatch := styles.Infer(firstLine)
	footerMatch := styles.Infer(lastLine)

	m.logger.LogInfo("  Header match: [%s] (score: %.2f)", headerMatch.Style.Name, headerMatch.Score)
	m.logger.LogInfo("  Footer match: [%s] (score: %.2f)", footerMatch.Style.Name, footerMatch.Score)

	// If both header and footer match the same style with high confidence, use that style
	if headerMatch.Score > 0.8 && footerMatch.Score > 0.8 && headerMatch.Style.Name == footerMatch.Style.Name {
		return headerMatch.Style
	}

	// If just the header matches with high confidence, use that style
	if headerMatch.Score > 0.8 {
		return headerMatch.Style
	}

	// Otherwise return current style
	return m.headerStyle
}

// GetHeaderStyle returns the current header style
func (m *LicenseManager) GetHeaderStyle() styles.HeaderFooterStyle {
	return m.headerStyle
}

// helper function to truncate strings for logging
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func looksLikeLicense(body string) bool {
	// TO DO: implement this function
	return true
}
