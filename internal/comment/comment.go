package comment

import (
	"strings"

	"license-manager/internal/language"
	"license-manager/internal/styles"
)

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

		// Skip comment start/end markers
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

		// Skip empty lines at the start or end
		if (i == 0 || i == len(lines)-1) && line == "" {
			continue
		}

		processedLines = append(processedLines, line)
	}

	return strings.TrimSpace(strings.Join(processedLines, "\n"))
}

// ExtractComponents extracts the header, body, and footer from a license block.
// Returns the extracted components and a success flag.
func ExtractComponents(content string, stripMarkers ...bool) (header string, body string, footer string, success bool) {
	shouldStrip := false
	if len(stripMarkers) > 0 {
		shouldStrip = stripMarkers[0]
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return "", "", "", false
	}

	// Find the first non-empty line after comment start
	startIdx := -1
	endIdx := -1
	hasCommentMarkers := false
	foundKnownStyle := false
	var commentStyle *styles.CommentLanguage
	var knownStyle styles.HeaderFooterStyle

	// First try to find comment markers
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if startIdx == -1 {
			// Try to detect the comment style from the first non-empty line
			if strings.HasPrefix(line, "/*") {
				commentStyle = &styles.CommentLanguage{MultiStart: "/*", MultiEnd: "*/"}
				hasCommentMarkers = true
			} else if strings.HasPrefix(line, "<!--") {
				commentStyle = &styles.CommentLanguage{MultiStart: "<!--", MultiEnd: "-->"}
				hasCommentMarkers = true
			} else {
				// If no comment markers, try to match against known header patterns
				match := styles.Infer(line)
				if match.Score > 0 && match.IsHeader {
					foundKnownStyle = true
					knownStyle = match.Style
				}
			}
			startIdx = i
			continue
		}
		if commentStyle != nil && strings.HasSuffix(line, commentStyle.MultiEnd) {
			hasCommentMarkers = true
			endIdx = i
		}
	}

	// If no comment markers found, try to find a matching footer
	if !hasCommentMarkers {
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			match := styles.Infer(line)
			if match.Score > 0 && match.IsFooter {
				if !foundKnownStyle {
					foundKnownStyle = true
					knownStyle = match.Style
				} else if match.Style.Name != knownStyle.Name {
					// If footer doesn't match the header style, reject it
					return "", "", "", false
				}
				endIdx = i
				break
			}
		}

		// If we didn't find any known style patterns, reject it
		if !foundKnownStyle {
			return "", "", "", false
		}

		// If we found a known style but no footer, reject it
		if endIdx == -1 {
			return "", "", "", false
		}
	}

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return "", "", "", false
	}

	// Extract header, body, and footer
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

	if shouldStrip {
		if commentStyle != nil {
			header = commentStyle.StripCommentMarkers(header)
			footer = commentStyle.StripCommentMarkers(footer)
		}
	}

	return header, body, footer, true
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

		// Skip empty lines
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
