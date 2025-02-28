package license

import (
	"reflect"
	"strings"

	"github.com/jeeftor/license-manager/internal/errors"
	"github.com/jeeftor/license-manager/internal/language"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
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
func NewLicenseManager(
	logger *logger.Logger,
	licenseTemplate, fileExtension string,
	headerStyle styles.HeaderFooterStyle,
	commentStyle styles.CommentLanguage,
) *LicenseManager {

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

type ScanResults struct {
	HasLicense   bool
	Components   *language.ExtractedComponents
	Style        styles.HeaderFooterStyle
	IsStyleMatch bool
}

func (m *LicenseManager) SearchForLicense(content string) ScanResults {
	actualType := reflect.TypeOf(m.langHandler)
	m.logger.LogInfo("  SearchForLicense::Handler type: %v", actualType)
	m.logger.LogInfo("  SearchForLicense::Analyzing license block...")
	m.logger.LogInfo("  SearchForLicense::Using comment style: %s", m.commentStyle.Language)

	components, success := m.langHandler.ExtractComponents(content)

	// Handle special case for c-go style headers first
	if m.commentStyle.Language == "go" && strings.HasPrefix(components.Header, "#include") {
		m.logger.LogInfo("  SearchForLicense::Detected c-go style header... skipping")

		// Preserve preamble handling for C-Go files
		//preamble, rest := m.langHandler.PreservePreamble(content)
		components = language.ExtractedComponents{
			Preamble:         components.Preamble,
			Header:           "",
			Body:             "",
			Footer:           "",
			Rest:             components.Header + components.Body + components.Footer + components.Body,
			FullLicenseBlock: nil,
		}

		// Store results and return early
		m.InitialComponents = &components
		m.HasInitialLicense = false

		return ScanResults{
			HasLicense:   false,
			Components:   &components,
			Style:        m.headerStyle,
			IsStyleMatch: false,
		}
	}

	// Store results for later use (maintaining backward compatibility)
	m.InitialComponents = &components
	m.HasInitialLicense = success

	analysis := ScanResults{
		HasLicense:   success,
		Components:   &components,
		Style:        m.headerStyle, // default to current style
		IsStyleMatch: false,
	}

	// Handle special case for c-go style headers
	if m.commentStyle.Language == "go" && strings.HasPrefix(components.Header, "#include") {
		m.logger.LogInfo("  SearchForLicense::Detected c-go style header... skipping")
		analysis.HasLicense = false
		return analysis
	}

	if !success {
		m.logger.LogInfo("  SearchForLicense::No license block detected")
		return analysis
	}

	// If we found a license, detect its style
	if components.Header != "" {
		// Try to match against known styles
		headerMatch := styles.Infer(components.Header)
		footerMatch := styles.Infer(components.Footer)

		m.logger.LogInfo(
			"Header match: [%s] (score: %.2f)",
			headerMatch.Style.Name,
			headerMatch.Score,
		)
		m.logger.LogInfo(
			"Footer match: [%s] (score: %.2f)",
			footerMatch.Style.Name,
			footerMatch.Score,
		)

		// If both header and footer match with high confidence
		if headerMatch.Score > 0.8 && footerMatch.Score > 0.8 &&
			headerMatch.Style.Name == footerMatch.Style.Name {
			analysis.Style = headerMatch.Style
			analysis.IsStyleMatch = true
		} else if headerMatch.Score > 0.8 {
			// If just header matches with high confidence
			analysis.Style = headerMatch.Style
			analysis.IsStyleMatch = true
		}
	}

	return analysis
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

// Common utility function to rebuild file content from components
func (m *LicenseManager) rebuildContent(
	components *language.ExtractedComponents,
	newLicenseBlock string,
) string {
	var parts []string

	if components.Preamble != "" {
		parts = append(parts, components.Preamble)
	}

	if newLicenseBlock != "" {
		parts = append(parts, newLicenseBlock)
	}

	if components.Rest != "" {
		parts = append(parts, components.Rest)
		//if !strings.HasPrefix(components.Rest, "\n") && len(parts) > 0 {
		//	parts = append(parts, "") // Add blank line before content
		//}
		//parts = append(parts, strings.TrimPrefix(components.Rest, "\n"))
	}

	return strings.Join(parts, "\n")
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

func (m *LicenseManager) AddLicense(
	components *language.ExtractedComponents,
	fileType string,
) (string, error) {
	m.logger.LogVerbose("Adding license to content...")

	// Get language handler
	handler := m.getLanguageHandler(fileType)

	// We already know the license status from SearchForLicense
	if m.HasInitialLicense {
		if status := m.CheckLicenseStatus(m.FileContent); status == FullMatch {
			return "", errors.NewLicenseError("license already exists and matches", "add")
		}
	}

	// Format license block with comment style
	licenseBlock := m.formatLicenseBlock(m.licenseTemplate)

	// Handle special case for c-go style headers
	if m.commentStyle.Language == "go" && strings.HasPrefix(components.Header, "#include") {
		m.logger.LogVerbose("Detected c-go style header... skipping")
		return m.FileContent, nil
	}

	// Scan for build directives in the Rest part
	directives, endIndex := handler.ScanBuildDirectives(components.Rest)
	if len(directives) > 0 {
		m.logger.LogVerbose("Found %d build directives", len(directives))
	}

	// Split the Rest content at build directives
	var beforeDirectives, afterDirectives string
	if len(directives) > 0 {
		lines := strings.Split(components.Rest, "\n")
		beforeDirectives = strings.Join(lines[:endIndex], "\n")
		if endIndex < len(lines) {
			afterDirectives = strings.Join(lines[endIndex:], "\n")
		}
	} else {
		afterDirectives = components.Rest
	}

	// Build final content
	var parts []string
	if components.Preamble != "" {
		parts = append(parts, components.Preamble)
	}
	if len(directives) > 0 {
		parts = append(parts, beforeDirectives)
	}
	parts = append(parts, licenseBlock)
	if afterDirectives != "" {
		//if !strings.HasPrefix(afterDirectives, "\n") {
		//	parts = append(parts, "") // Add blank line before content
		//}
		parts = append(parts, afterDirectives)
	}
	return strings.Join(parts, "\n"), nil

	//return strings.TrimSuffix(strings.Join(parts, "\n"), "\n"), nil
}

// RemoveLicense removes the license block from the content using existing components
func (m *LicenseManager) RemoveLicense(
	components *language.ExtractedComponents,
	fileType string,
) (string, error) {
	m.logger.LogDebug("  Attempting to remove license block...")

	if !m.HasInitialLicense {
		m.logger.LogDebug("  RemoveLicense::No license block detected in content")
		return m.FileContent, nil
	}

	// Rebuild content without the license block
	return m.rebuildContent(components, ""), nil
}

// UpdateLicense updates the license block using existing components
func (m *LicenseManager) UpdateLicense(
	components *language.ExtractedComponents,
	fileType string,
) (string, error) {
	m.logger.LogDebug("Attempting to update license block...")

	if !m.HasInitialLicense {
		return "", errors.NewLicenseError("content has no license to update", "")
	}

	// Format the new license block
	newLicenseBlock := m.formatLicenseBlock(m.licenseTemplate)

	// Rebuild the content with the new license block
	return m.rebuildContent(components, newLicenseBlock), nil
}

//// UpdateLicense updates the license block in the content
//func (m *LicenseManager) UpdateLicense(content string, fileType string) (string, error) {
//	m.logger.LogDebug("Attempting to update license block...")
//	handler := m.getLanguageHandler(fileType)
//
//	// First extract any preamble (build directives, etc.)
//	preamble, rest := handler.PreservePreamble(content)
//
//	// Check if rest has a license block
//	hasLicense, _ := m.HasLicense(rest)
//	if !hasLicense {
//		return "", errors.NewLicenseError("content has no license to update", "")
//	}
//
//	// Extract components from the rest of the content
//	extract, _ := handler.ExtractComponents(rest)
//
//	// Build content without license:
//	// 1. Start with preamble if it exists
//	// 2. Add rest of content after license
//	var parts []string
//	if preamble != "" {
//		parts = append(parts, preamble)
//	}
//	if extract.Rest != "" {
//		parts = append(parts, strings.TrimSpace(extract.Rest))
//	}
//
//	// Add the new license
//	return m.AddLicense(strings.Join(parts, "\n\n"), fileType)
//}

func (m *LicenseManager) CheckLicenseStatus(content string) Status {

	handler := m.langHandler
	actualType := reflect.TypeOf(handler)
	m.logger.LogInfo("CheckLicenseStatus::Handler type: %v", actualType)
	m.logger.LogInfo("CheckLicenseStatus:Checking license status...")
	m.logger.LogInfo("CheckLicenseStatus:Using comment style: %s", m.commentStyle.Language)

	actualExtract, success := handler.ExtractComponents(content)

	if actualExtract.Preamble != "" {
		m.logger.LogInfo(
			"Found preamble 📝️ (%d lines)",
			len(strings.Split(actualExtract.Preamble, "\n")),
		)
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
	detectedStyle, foundHeaderMatch := m.DetectHeaderAndFooterStyle(
		actualExtract.Header,
		actualExtract.Footer,
	)

	if foundHeaderMatch {
		m.logger.LogInfo("Detected style: %s", detectedStyle.Name)
	} else {
		return NoLicense
	}

	// Use the detected style + the license template to generate the expected styles - and start comparing
	expectedLicenseText := language.FormatComment(m.licenseTemplate, m.commentStyle, m.headerStyle)
	expectedExtract, _ := handler.ExtractComponents(expectedLicenseText)

	actualBody := actualExtract.Body
	expectedBody := expectedExtract.Body

	// If headers don't match
	if m.headerStyle.Name != "" && m.headerStyle.Name != detectedStyle.Name {
		// Mismatch of headers - but if license ist he same
		m.logger.LogInfo(
			"Style mismatch: expected [%s], found [%s]",
			m.headerStyle.Name,
			detectedStyle.Name,
		)

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

func (m *LicenseManager) DetectHeaderAndFooterStyle(
	header, footer string,
) (styles.HeaderFooterStyle, bool) {
	// Try to match against known styles
	headerMatch := styles.Infer(header)
	footerMatch := styles.Infer(footer)

	m.logger.LogInfo("Header match: [%s] (score: %.2f)", headerMatch.Style.Name, headerMatch.Score)
	m.logger.LogInfo("Footer match: [%s] (score: %.2f)", footerMatch.Style.Name, footerMatch.Score)

	// If both header and footer match the same style with high confidence, use that style
	if headerMatch.Score > 0.8 && footerMatch.Score > 0.8 &&
		headerMatch.Style.Name == footerMatch.Style.Name {
		return headerMatch.Style, true
	}

	// If just the header matches with high confidence, use that style
	if headerMatch.Score > 0.8 {
		return headerMatch.Style, true
	}

	// Otherwise return current style
	return m.headerStyle, false
}

func (m *LicenseManager) detectHeaderStyle(
	components language.ExtractedComponents,
) styles.HeaderFooterStyle {

	m.logger.LogInfo("Trying to detect style from header: %q", components.Header)
	m.logger.LogInfo("Trying to detect style from footer: %q", components.Footer)

	// Try to match against known styles
	headerMatch := styles.Infer(components.Header)
	footerMatch := styles.Infer(components.Footer)

	m.logger.LogInfo("Header match: [%s] (score: %.2f)", headerMatch.Style.Name, headerMatch.Score)
	m.logger.LogInfo("Footer match: [%s] (score: %.2f)", footerMatch.Style.Name, footerMatch.Score)

	// If both header and footer match the same style with high confidence, use that style
	if headerMatch.Score > 0.8 && footerMatch.Score > 0.8 &&
		headerMatch.Style.Name == footerMatch.Style.Name {
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
