package license

import (
	"github.com/jeeftor/license-manager/internal/errors"
	"github.com/jeeftor/license-manager/internal/language"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
	"reflect"
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
	licenseTemplate   string
	commentStyle      styles.CommentLanguage
	headerStyle       styles.HeaderFooterStyle
	langHandler       language.LanguageHandler
	logger            *logger.Logger
	InitialComponents *language.ExtractedComponents
	HasInitialLicense bool // did we detect a license at startup of the manager
	//todo: Should we rename this variable later
	FileContent string
}

// NewLicenseManager creates a new manager
func NewLicenseManager(logger *logger.Logger, licenseTemplate, fileExtension string, headerStyle styles.HeaderFooterStyle, commentStyle styles.CommentLanguage) *LicenseManager {

	// Determine Language Handler
	langHandler := language.GetLanguageHandler(logger, fileExtension, headerStyle)

	manager := &LicenseManager{
		licenseTemplate: licenseTemplate,
		langHandler:     langHandler,
		headerStyle:     headerStyle,
		commentStyle:    commentStyle,
		logger:          logger,
	}

	return manager
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
func (m *LicenseManager) HasLicense(content string) (bool, *language.ExtractedComponents) {
	actualType := reflect.TypeOf(m.langHandler)
	m.logger.LogInfo("  HasLicense::Handler type: %v", actualType)
	m.logger.LogInfo("  HasLicense::Checking for existing license block...")
	m.logger.LogInfo("  HasLicense::Using comment style: %s", m.commentStyle.Language)

	//components, success := m.langHandler.ExtractComponents(content)
	components, success := m.langHandler.ExtractComponents(content)

	// Store results for later use
	m.InitialComponents = &components
	m.HasInitialLicense = success

	if components.Header == components.Footer {
		m.logger.LogInfo("  HasLicense::Header & Footer Match")
	}

	// Handle Special Cases
	// TODO: Come up with a better way
	if m.commentStyle.Language == "go" && strings.HasPrefix(components.Header, "#include") {
		m.logger.LogInfo("  HasLicense::Detected c-go style header... skipping")

		preamble, rest := m.langHandler.PreservePreamble(content)
		m.InitialComponents = &language.ExtractedComponents{
			Preamble:         preamble,
			Header:           "",
			Body:             "",
			Footer:           "",
			Rest:             rest,
			FullLicenseBlock: nil,
		}
		m.HasInitialLicense = false

		return false, m.InitialComponents

	}

	if success {
		m.logger.LogInfo("  HasLicense::Found existing license block")
	} else {
		m.logger.LogInfo("  HasLicense::No license block detected in content")
	}

	return success, &components
}

// AddLicense adds a license block to the content
func (m *LicenseManager) AddLicense(content string, fileType string) (string, error) {
	m.logger.LogVerbose("Adding license to content...")

	// Get language handler
	handler := m.getLanguageHandler(fileType)

	// Extract preamble (e.g., shebang, package declaration)
	preamble, rest := handler.PreservePreamble(content)
	if m.logger != nil {
		if preamble != "" {
			m.logger.LogVerbose("Found preamble: %s", truncateString(preamble, 50))
		}
	}

	hasLicense, _ := m.HasLicense(rest)
	// Check for existing license
	if hasLicense {
		// If we're updating, allow the license to be replaced
		if status := m.CheckLicenseStatus(rest); status == FullMatch {
			return "", errors.NewLicenseError("license already exists and matches", "add")
		}
	}

	// Format license block with comment style
	licenseBlock := m.formatLicenseBlock(m.licenseTemplate)

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
func (m *LicenseManager) RemoveLicense(content string, fileType string) (string, error) {

	m.logger.LogDebug("  Attempting to remove license block...")
	handler := m.getLanguageHandler(fileType)
	_, rest := handler.PreservePreamble(content)
	hasLicense, _ := m.HasLicense(rest)
	if !hasLicense {
		m.logger.LogDebug("  RemoveLicense::No license block detected in content")
		return content, nil
	}

	// Removal Logic - Can be inlined I think
	extract, _ := handler.ExtractComponents(content)

	// Remove existing license -> including a \n usually added after it
	contentWithoutLicense := extract.Preamble + "\n" + strings.TrimLeft(extract.Rest, "\n")
	return contentWithoutLicense, nil
}

// UpdateLicense updates the license block in the content
func (m *LicenseManager) UpdateLicense(content string, fileType string) (string, error) {

	m.logger.LogDebug("Attempting to update license block...")
	// Use the filetype to identify the correct langauge Handler
	handler := m.getLanguageHandler(fileType)

	// Determine whether there IS a license
	_, rest := handler.PreservePreamble(content)
	hasLicense, _ := m.HasLicense(rest)
	if !hasLicense {
		return "", errors.NewLicenseError("content has no license to update", "")
	}

	// Removal Logic - Can be inlined I think
	extract, _ := handler.ExtractComponents(content)

	// Remove existing license -> including a \n usually added after it
	contentWithoutLicense := extract.Preamble + strings.TrimLeft(extract.Rest, "\n")

	// Add the new license
	return m.AddLicense(contentWithoutLicense, m.commentStyle.Language)
}

func (m *LicenseManager) CheckLicenseStatus(content string) Status {

	handler := m.langHandler
	actualType := reflect.TypeOf(handler)
	m.logger.LogInfo("CheckLicenseStatus::Handler type: %v", actualType)
	m.logger.LogInfo("CheckLicenseStatus:Checking license status...")
	m.logger.LogInfo("CheckLicenseStatus:Using comment style: %s", m.commentStyle.Language)

	actualExtract, success := handler.ExtractComponents(content)

	if actualExtract.Preamble != "" {
		m.logger.LogInfo("Found preamble 📝️ (%d lines)", len(strings.Split(actualExtract.Preamble, "\n")))
	}

	if m.InitialComponents != nil {
		actualExtract = *m.InitialComponents
		success = m.HasInitialLicense
	} else {
		// Only extract if we don't have stored results
		actualExtract, success = handler.ExtractComponents(content)
		m.InitialComponents = &actualExtract
		m.HasInitialLicense = success
	}

	// Extract license components
	if !success {
		m.logger.LogInfo("Result: No license block found")
		return NoLicense
	}

	// Detect and validate style
	detectedStyle := m.DetectHeaderAndFooterStyle(actualExtract.Header, actualExtract.Footer)
	m.logger.LogInfo("Detected style: %s", detectedStyle.Name)

	// Use the detected style + the license template to generate the expected styles - and start comparing
	expectedLicenseText := language.FormatComment(m.licenseTemplate, m.commentStyle, m.headerStyle)
	expectedExtract, _ := handler.ExtractComponents(expectedLicenseText)

	actualBody := actualExtract.Body
	expectedBody := expectedExtract.Body

	// If headers don't match
	if m.headerStyle.Name != "" && m.headerStyle.Name != detectedStyle.Name {
		// Mismatch of headers - but if license ist he same
		m.logger.LogInfo("Style mismatch: expected [%s], found [%s]", m.headerStyle.Name, detectedStyle.Name)

		if actualBody == expectedBody {
			// We body match w/out same headers
			return StyleMismatch
		}

		return ContentAndStyleMismatch
	}

	if actualBody == expectedBody {
		return FullMatch
	}

	// Handle mismatch
	m.logger.LogInfo("Result: License content differs from expected")
	m.logDiff(expectedBody, actualBody) // Extract diff logging to separate method
	return ContentMismatch
}

func (m *LicenseManager) logDiff(expected, current string) {
	m.logger.LogInfo("BodySize current[%d] expected[%d]",
		strings.Count(current, "\n"),
		strings.Count(expected, "\n"))
	m.logger.LogInfo(" current: %x", []byte(current))
	m.logger.LogInfo("expected: %x", []byte(expected))

	lines1 := strings.Split(expected, "\n")
	lines2 := strings.Split(current, "\n")

	m.logger.LogInfo("\nDiff: \033[31m expected \033[0m\033[32m current\033[0m")
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

// FormatLicenseForFile formats the license text with the current comment style
func (m *LicenseManager) FormatLicenseForFile(text string) string {
	if m.commentStyle.Language == "" {
		return "No comment style set - cannot format license"
	}
	return m.formatLicenseBlock(text)
}

// formatLicenseBlock formats the license text with the appropriate comment style
func (m *LicenseManager) formatLicenseBlock(text string) string {
	return language.FormatComment(text, m.commentStyle, m.headerStyle)
}

func (m *LicenseManager) DetectHeaderAndFooterStyle(header, footer string) styles.HeaderFooterStyle {
	// Try to match against known styles
	headerMatch := styles.Infer(header)
	footerMatch := styles.Infer(footer)

	m.logger.LogInfo("Header match: [%s] (score: %.2f)", headerMatch.Style.Name, headerMatch.Score)
	m.logger.LogInfo("Footer match: [%s] (score: %.2f)", footerMatch.Style.Name, footerMatch.Score)

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

func (m *LicenseManager) detectHeaderStyle(components language.ExtractedComponents) styles.HeaderFooterStyle {

	m.logger.LogInfo("Trying to detect style from header: %q", components.Header)
	m.logger.LogInfo("Trying to detect style from footer: %q", components.Footer)

	// Try to match against known styles
	headerMatch := styles.Infer(components.Header)
	footerMatch := styles.Infer(components.Footer)

	m.logger.LogInfo("Header match: [%s] (score: %.2f)", headerMatch.Style.Name, headerMatch.Score)
	m.logger.LogInfo("Footer match: [%s] (score: %.2f)", footerMatch.Style.Name, footerMatch.Score)

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

func (m *LicenseManager) SetHeaderStyle(style styles.HeaderFooterStyle) {
	m.headerStyle = style
}

func (m *LicenseManager) SetFileContent(content string) {
	m.FileContent = content
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
