package comment

import (
	"strings"

	"license-manager/internal/language"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
)

var log *logger.Logger

// SetVerbose enables or disables verbose logging for the comment package
func SetVerbose(verbose bool) {
	log = logger.NewLogger(verbose)
}

func init() {
	// Initialize with default (non-verbose) logging
	log = logger.NewLogger(false)
}

const (
	MarkerStart = "​" // Zero-width space
	MarkerEnd   = "‌" // Zero-width non-joiner
)

// Comment represents a complete comment block with all its components
type Comment struct {
	style           styles.CommentLanguage
	header          string
	body            string
	footer          string
	hfStyle         styles.HeaderFooterStyle
	CommentLanguage styles.CommentLanguage
	langHandler     language.LanguageHandler
}

func (c *Comment) String() string {
	var result []string
	content := []string{}

	// Add header if present
	if c.header != "" {
		content = append(content, c.header)
	}

	// Add body
	if c.body != "" {
		if len(content) > 0 {
			content = append(content, "")
		}
		content = append(content, c.body)
	}

	// Add footer if present
	if c.footer != "" {
		if len(content) > 0 {
			content = append(content, "")
		}
		content = append(content, c.footer)
	}

	// Join all content with newlines
	text := strings.Join(content, "\n")

	if c.style.PreferMulti && c.style.MultiStart != "" {
		// Multi-line comment style
		result = append(result, c.style.MultiStart)

		// Process each line
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if c.style.MultiPrefix != "" {

				if line == "" {
					result = append(result, c.style.MultiPrefix)
				} else {
					result = append(result, c.style.MultiPrefix+c.style.LinePrefix+line)
				}
			} else {
				result = append(result, line)
			}
		}

		result = append(result, c.style.MultiEnd)
	} else if c.style.Single != "" {
		// Single-line comment style
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if line == "" {
				result = append(result, "")
			} else {
				result = append(result, c.style.Single+c.style.LinePrefix+line)
			}
		}
	} else if c.style.MultiStart != "" {
		// Fallback to multi-line style for languages that only support multi-line comments
		result = append(result, c.style.MultiStart)
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			result = append(result, line)
		}
		result = append(result, c.style.MultiEnd)
	} else {
		// No comment style defined, return raw text
		return text
	}

	return strings.Join(result, "\n")
}

// CommentStyle represents how comments should be formatted for a specific language
type CommentStyle struct {
	// Language identifier (e.g., "go", "python", "javascript")
	Language string

	// Single-line comment prefix (e.g., "//", "#")
	Single string

	// Multi-line comment start marker (e.g., "/*", "'''")
	MultiStart string

	// Multi-line comment end marker (e.g., "*/", "'''")
	MultiEnd string

	// Multi-line comment line prefix (e.g., " * ")
	MultiPrefix string

	// Line prefix for content after comment marker (e.g., " ")
	LinePrefix string

	// Whether to prefer multi-line comments over single-line
	PreferMulti bool

	// Header and footer for the comment block
	Header string
	Footer string
}

func UncommentContent(content string, style styles.CommentLanguage) string {
	// Split into lines for processing
	lines := strings.Split(content, "\n")
	processedLines := make([]string, 0, len(lines))

	// Process each line
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skips comment start/end markers
		if !hasMarkers(line) {
			if strings.HasPrefix(line, style.MultiStart) {
				line = strings.TrimSpace(strings.TrimPrefix(line, style.MultiStart))
			}
			if strings.HasSuffix(line, style.MultiEnd) {
				line = strings.TrimSpace(strings.TrimSuffix(line, style.MultiEnd))
			}
			if line == "" {
				continue
			}
		}

		// Handle line prefixes while preserving markers
		if hasMarkers(line) {
			start := strings.Index(line, MarkerStart)
			end := strings.Index(line, MarkerEnd) + len(MarkerEnd)
			markers := line[start:end]

			// Keep the line as is if it only contains markers
			if start == 0 && end == len(line) {
				processedLines = append(processedLines, markers)
				continue
			}
		}

		// Remove comment prefixes
		if style.Single != "" && strings.HasPrefix(line, style.Single) {
			line = strings.TrimSpace(strings.TrimPrefix(line, style.Single))
		}
		if style.MultiPrefix != "" {
			// Handle both MultiPrefix with and without LinePrefix
			fullPrefix := style.MultiPrefix + style.LinePrefix
			if strings.HasPrefix(line, fullPrefix) {
				line = strings.TrimSpace(strings.TrimPrefix(line, fullPrefix))
			} else if strings.HasPrefix(line, style.MultiPrefix) {
				line = strings.TrimSpace(strings.TrimPrefix(line, style.MultiPrefix))
			}
		}

		// Handle empty lines specially - if they were originally commented, keep them
		if line == "" && i > 0 && i < len(lines)-1 {
			processedLines = append(processedLines, "")
			continue
		}

		// Skips empty lines at the start or end
		if (i == 0 || i == len(lines)-1) && line == "" {
			continue
		}

		processedLines = append(processedLines, line)
	}

	return strings.TrimSpace(strings.Join(processedLines, "\n"))
}

// ExtractComponents extracts the header, body, and footer from a license block
func ExtractComponents(content string, stripMarkers bool, languageStyle styles.CommentLanguage) (header, body, footer string, success bool) {
	if content == "" {
		return "", "", "", false
	}

	// Split into lines for processing
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", "", "", false
	}

	// Find the start and end of the license block
	var startIndex, endIndex int
	var foundStart, foundEnd bool

	//for i, line := range lines {
	//	line = strings.TrimSpace(line)
	//	if line == "" {
	//		continue
	//	}
	//
	//	// Look for start markers using language style
	//	if !foundStart && languageStyle.MultiStart != "" && strings.HasPrefix(line, languageStyle.MultiStart) {
	//		startIndex = i
	//		foundStart = true
	//		continue
	//	}
	//
	//	// Look for end markers
	//	if foundStart && !foundEnd && languageStyle.MultiEnd != "" && strings.HasSuffix(line, languageStyle.MultiEnd) {
	//		endIndex = i
	//		foundEnd = true
	//		break
	//	}
	//}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !foundStart && languageStyle.MultiStart != "" && strings.HasPrefix(line, languageStyle.MultiStart) {
			startIndex = i
			foundStart = true
		} else if foundStart && languageStyle.MultiEnd != "" && strings.HasSuffix(line, languageStyle.MultiEnd) {
			endIndex = i
			foundEnd = true
			break
		}
	}
	if !foundStart || !foundEnd {
		return "", "", "", false
	}

	// Extract the header (first non-empty line after the start marker)
	var headerLines []string
	for i := startIndex + 1; i < endIndex; i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		headerLines = append(headerLines, lines[i])
		break
	}

	// Extract the footer (last non-empty line before the end marker)
	var footerLines []string
	for i := endIndex - 1; i > startIndex; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		footerLines = append(footerLines, lines[i])
		break
	}

	// Extract the body (everything between header and footer)
	var bodyLines []string
	bodyStart := startIndex + 1
	if len(headerLines) > 0 {
		bodyStart++
	}
	bodyEnd := endIndex
	if len(footerLines) > 0 {
		bodyEnd--
	}

	for i := bodyStart; i < bodyEnd; i++ {
		line := lines[i]
		if stripMarkers {
			line = strings.TrimSpace(lines[i])
			line = strings.TrimPrefix(line, languageStyle.LinePrefix)
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimPrefix(line, " *")
			line = strings.TrimSpace(line)
		}
		bodyLines = append(bodyLines, line)
	}

	header = strings.Join(headerLines, "\n")

	body = strings.Join(bodyLines, "\n")
	footer = strings.Join(footerLines, "\n")

	return header, body, footer, true
}

// extractComponentsWithMarkers attempts to extract components using Unicode markers.
// This is a simpler implementation that looks specifically for MarkerStart and MarkerEnd.
func extractComponentsWithMarkers(lines []string, shouldStrip bool) (header string, body string, footer string, success bool) {
	startIdx := -1
	endIdx := -1

	// Find lines with markers
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skips empty lines
		if line == "" {
			continue
		}

		if hasMarkers(line) {
			if startIdx == -1 {
				startIdx = i
			} else if endIdx == -1 {
				endIdx = i
				break
			}
		}
	}

	// Must have both start and end markers
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return "", "", "", false
	}

	header = strings.TrimSpace(lines[startIdx])
	footer = strings.TrimSpace(lines[endIdx])

	// Extract body (everything between header and footer)
	bodyLines := lines[startIdx+1 : endIdx]
	// Remove leading and trailing empty lines from body
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[0]) == "" {
		bodyLines = bodyLines[1:]
	}
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[len(bodyLines)-1]) == "" {
		bodyLines = bodyLines[:len(bodyLines)-1]
	}
	body = strings.TrimSpace(strings.Join(bodyLines, "\n"))

	// Strip the Unicode markers if requested
	if shouldStrip {
		header = strings.TrimPrefix(strings.TrimSuffix(header, MarkerEnd), MarkerStart)
		footer = strings.TrimPrefix(strings.TrimSuffix(footer, MarkerEnd), MarkerStart)
	}

	// Check if the content looks like a license
	if !looksLikeLicense(body) {
		return "", "", "", false
	}

	return header, body, footer, true
}

// extractComponentsWithoutMarkers attempts to extract components using comment syntax and style inference.
// This is a more intensive scan that tries to identify the comment type based on available styles.
func extractComponentsWithoutMarkers(lines []string, shouldStrip bool) (header string, body string, footer string, success bool) {
	startIdx := -1
	endIdx := -1
	hasCommentMarkers := false
	foundKnownStyle := false
	var commentStyle *styles.CommentLanguage
	var knownStyle styles.HeaderFooterStyle

	log.LogVerbose("Starting style detection...")

	// Scan for comment markers and known styles
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if startIdx == -1 {
			// Try to detect the comment style from the first non-empty line
			if strings.HasPrefix(line, "/*") {
				log.LogVerbose("Found C-style comment: %s", line)
				commentStyle = &styles.CommentLanguage{MultiStart: "/*", MultiEnd: "*/"}
				hasCommentMarkers = true
			} else if strings.HasPrefix(line, "<!--") {
				log.LogVerbose("Found HTML-style comment: %s", line)
				commentStyle = &styles.CommentLanguage{MultiStart: "<!--", MultiEnd: "-->"}
				hasCommentMarkers = true
			} else if strings.HasPrefix(line, "#") {
				log.LogVerbose("Found Shell-style comment: %s", line)
				commentStyle = &styles.CommentLanguage{Single: "#"}
				hasCommentMarkers = true
			} else if strings.HasPrefix(line, "//") {
				log.LogVerbose("Found C++-style comment: %s", line)
				commentStyle = &styles.CommentLanguage{Single: "//"}
				hasCommentMarkers = true
			} else if line == "'''" || line == `"""` {
				log.LogVerbose("Found Python-style comment: %s", line)
				commentStyle = &styles.CommentLanguage{MultiStart: line, MultiEnd: line}
				hasCommentMarkers = true
			} else {
				// If no comment markers, try to match against known header patterns
				match := styles.Infer(line)
				if match.Score > 0 && match.IsHeader {
					log.LogVerbose("Found known header style: %s (score: %.2f)", match.Style.Name, match.Score)
					foundKnownStyle = true
					knownStyle = match.Style
					startIdx = i
				} else {
					log.LogVerbose("No style match for line: %s", line)
				}
			}
			if hasCommentMarkers {
				startIdx = i
				log.LogVerbose("Set start index to %d", i)
			}
			continue
		}

		// If we have comment markers, look for the end marker
		if hasCommentMarkers {
			if commentStyle.MultiEnd != "" {
				if line == commentStyle.MultiEnd || strings.HasSuffix(line, commentStyle.MultiEnd) {
					log.LogVerbose("Found matching end marker at line %d: %s", i, line)
					endIdx = i
					break
				}
			} else if commentStyle.Single != "" {
				// For single-line comments, look for the first non-comment line
				if !strings.HasPrefix(strings.TrimSpace(line), commentStyle.Single) {
					log.LogVerbose("Found end of single-line comments at line %d", i)
					endIdx = i - 1 // End at the last comment line
					break
				}
			}
		} else {
			// Look for a matching footer
			match := styles.Infer(line)
			if match.Score > 0 && match.IsFooter {
				if !foundKnownStyle {
					log.LogVerbose("Found footer style before header: %s (score: %.2f)", match.Style.Name, match.Score)
					foundKnownStyle = true
					knownStyle = match.Style
				} else if match.Style.Name != knownStyle.Name {
					log.LogVerbose("Footer style mismatch: found %s but expected %s", match.Style.Name, knownStyle.Name)
					// If footer doesn't match the header style, reject it
					return "", "", "", false
				}
				log.LogVerbose("Found matching footer at line %d: %s", i, line)
				endIdx = i
				break
			} else if match.Score > 0 {
				log.LogVerbose("Found potential footer match: %s (score: %.2f) but not a footer pattern", match.Style.Name, match.Score)
			}
		}
	}

	// For single-line comments, if we haven't found an end, use the last line
	if hasCommentMarkers && commentStyle.Single != "" && endIdx == -1 {
		// Find the last line that starts with the comment marker
		for i := len(lines) - 1; i > startIdx; i-- {
			line := strings.TrimSpace(lines[i])
			if line != "" && strings.HasPrefix(line, commentStyle.Single) {
				log.LogVerbose("Found last single-line comment at line %d", i)
				endIdx = i
				break
			}
		}
	}

	if !foundKnownStyle && !hasCommentMarkers {
		log.LogVerbose("No known style patterns or comment markers found")
	}
	if !hasCommentMarkers && endIdx == -1 {
		log.LogVerbose("Found known style but no footer")
	}
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		log.LogVerbose("Invalid start/end indices: start=%d, end=%d", startIdx, endIdx)
	}

	// If we didn't find any known style patterns and no comment markers, reject it
	if !foundKnownStyle && !hasCommentMarkers {
		log.LogVerbose("No known style patterns or comment markers found")
		return "", "", "", false
	}

	// If we found a known style but no footer, reject it
	if !hasCommentMarkers && endIdx == -1 {
		log.LogVerbose("Found known style but no footer")
		return "", "", "", false
	}

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		log.LogVerbose("Invalid start/end indices: start=%d, end=%d", startIdx, endIdx)
		return "", "", "", false
	}

	// Extract header, body, and footer
	header = strings.TrimSpace(lines[startIdx])
	footer = strings.TrimSpace(lines[endIdx])

	log.LogVerbose("Extracted header: %s", header)
	log.LogVerbose("Extracted footer: %s", footer)

	if foundKnownStyle {
		log.LogInfo("Found matching style: %s", knownStyle.Name)
		log.LogVerbose("  Header: %s", knownStyle.Header)
		log.LogVerbose("  Footer: %s", knownStyle.Footer)
	}

	// Extract body (everything between header and footer)
	bodyLines := lines[startIdx+1 : endIdx]
	// Remove leading and trailing empty lines from body
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[0]) == "" {
		bodyLines = bodyLines[1:]
	}
	for len(bodyLines) > 0 && strings.TrimSpace(bodyLines[len(bodyLines)-1]) == "" {
		bodyLines = bodyLines[:len(bodyLines)-1]
	}
	body = strings.TrimSpace(strings.Join(bodyLines, "\n"))

	if shouldStrip {
		if commentStyle != nil {
			header = commentStyle.StripCommentMarkers(header)
			footer = commentStyle.StripCommentMarkers(footer)
			log.LogVerbose("Stripped header: %s", header)
			log.LogVerbose("Stripped footer: %s", footer)
		}
	}

	// Check if the content looks like a license
	if !looksLikeLicense(body) {
		log.LogVerbose("Content does not look like a license")
		return "", "", "", false
	}

	log.LogVerbose("Successfully extracted license components")

	return header, body, footer, true
}

// Common words that indicate a block of text is likely a license
var licenseIndicators = []string{
	"copyright",
	"license",
	"permission",
	"permitted",
	"granted",
	"rights",
	"reserved",
	"warranties",
	"liability",
	"contributors",
	"apache",
	"mit ",
	"bsd ",
	"gpl",
	"lgpl",
	"mozilla",
	"boost",
}

// looksLikeLicense checks if the content appears to be a license by looking for common license-related terms
// and checking the length of the content
func looksLikeLicense(content string) bool {
	if len(content) < 50 { // Most licenses are longer than 50 characters
		return false
	}

	lowerContent := strings.ToLower(content)

	// Check for common license indicators
	for _, indicator := range licenseIndicators {
		if strings.Contains(lowerContent, indicator) {
			return true
		}
	}

	return false
}

// FormatComment formats text with the given comment style and header/footer style
func FormatComment(text string, commentStyle styles.CommentLanguage, headerStyle styles.HeaderFooterStyle) string {
	lines := strings.Split(text, "\n")
	var result []string

	// Add header
	if commentStyle.MultiStart != "" {
		result = append(result, commentStyle.MultiStart)
		// Add marker to header
		result = append(result, commentStyle.MultiPrefix+MarkerStart+headerStyle.Header+MarkerEnd)
	} else {
		// Add marker to header
		result = append(result, commentStyle.Single+MarkerStart+headerStyle.Header+MarkerEnd)
	}

	// Add body with proper comment prefixes
	for _, line := range lines {
		if commentStyle.MultiStart != "" {
			if line == "" {
				result = append(result, commentStyle.MultiPrefix)
			} else {
				result = append(result, commentStyle.MultiPrefix+commentStyle.LinePrefix+line)
			}
		} else {
			if line == "" {
				result = append(result, commentStyle.Single)
			} else {
				result = append(result, commentStyle.Single+commentStyle.LinePrefix+line)
			}
		}
	}

	// Add footer
	if commentStyle.MultiStart != "" {
		// Add marker to footer
		result = append(result, commentStyle.MultiPrefix+MarkerStart+headerStyle.Footer+MarkerEnd)
		result = append(result, commentStyle.MultiEnd)
	} else {
		// Add marker to footer
		result = append(result, commentStyle.Single+MarkerStart+headerStyle.Footer+MarkerEnd)
	}

	return strings.Join(result, "\n")
}

// BuildDirective represents a Go build directive
type BuildDirective struct {
	Type    string // "go" or "plus" for //go: or // + style
	Content string // The actual directive content
}

// ExtractBuildDirectives extracts all build directives from the given content.
// It handles both //go: style directives and // + build style directives.
func ExtractBuildDirectives(content string, langHandler language.LanguageHandler) []BuildDirective {
	var directives []BuildDirective

	// Use language handler to get build directives
	directiveLines, _ := langHandler.ScanBuildDirectives(content)

	for _, line := range directiveLines {
		line = strings.TrimSpace(line)

		// Skips empty lines
		if line == "" {
			continue
		}

		// Check for //go: directives
		if strings.HasPrefix(line, "//go:") {
			directive := strings.TrimPrefix(line, "//go:")
			directives = append(directives, BuildDirective{
				Type:    "go",
				Content: strings.TrimSpace(directive),
			})
			continue
		}

		// Check for // +build directives
		if strings.HasPrefix(line, "// +") || strings.HasPrefix(line, "//+") {
			directive := strings.TrimPrefix(strings.TrimPrefix(line, "// +"), "//+")
			directives = append(directives, BuildDirective{
				Type:    "plus",
				Content: strings.TrimSpace(directive),
			})
		}
	}

	return directives
}

// Internal helper functions for working with zero-width space markers
func hasMarkers(text string) bool {
	return strings.Contains(text, MarkerStart) && strings.Contains(text, MarkerEnd)
}

func addMarkers(text string) string {
	if hasMarkers(text) {
		return text
	}
	return MarkerStart + text + MarkerEnd
}

func NewComment(style styles.CommentLanguage, hfStyle styles.HeaderFooterStyle, body string, langHandler language.LanguageHandler) *Comment {
	return &Comment{
		style:       style,
		body:        body,
		hfStyle:     hfStyle,
		header:      hfStyle.Header,
		footer:      hfStyle.Footer,
		langHandler: langHandler,
	}
}

func (c *Comment) Clone() *Comment {

	return &Comment{
		style:       c.style,
		header:      c.header,
		body:        c.body,
		footer:      c.footer,
		hfStyle:     c.hfStyle,
		langHandler: c.langHandler,
	}
}

func (c *Comment) SetBody(body string) {
	c.body = body
}

//func (c *Comment) SetStyle(style styles.CommentLanguage) {
//	c.style = style
//}

func (c *Comment) SetHeaderFooterStyle(hfStyle styles.HeaderFooterStyle) {
	c.hfStyle = hfStyle
	c.header = hfStyle.Header
	c.footer = hfStyle.Footer
}

//func (c *Comment) GetStyle() styles.CommentLanguage {
//	return c.style
//}
//
//func (c *Comment) GetHeader() string {
//	return c.header
//}
//
//func (c *Comment) GetFooter() string {
//	return c.footer
//}
//
//func (c *Comment) GetBody() string {
//	return c.body
//}
