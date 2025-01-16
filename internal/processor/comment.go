package processor

import (
	"path/filepath"
	"strings"

	"license-manager/internal/styles"
)

const (
	markerStart = "​" // Zero-width space
	markerEnd   = "‌" // Zero-width non-joiner
)

// Comment represents a complete comment block with all its components
type Comment struct {
	Style  styles.CommentStyle
	Header string
	Body   string
	Footer string
}

func (c *Comment) String() string {
	var result []string
	content := []string{}

	// Add header if present
	if c.Header != "" {
		content = append(content, c.Header)
	}

	// Add an empty line after header if we have more content
	if c.Header != "" && (c.Body != "" || c.Footer != "") {
		content = append(content, "")
	}

	// Add body
	if c.Body != "" {
		content = append(content, c.Body)
	}

	// Add an empty line before footer if we have one
	if c.Footer != "" && (c.Header != "" || c.Body != "") {
		content = append(content, "")
	}

	// Add footer if present
	if c.Footer != "" {
		content = append(content, c.Footer)
	}

	// Join all content with newlines
	text := strings.Join(content, "\n")

	if c.Style.PreferMulti && c.Style.MultiStart != "" {
		// Multi-line comment style
		result = append(result, c.Style.MultiStart)

		// Process each line
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if c.Style.MultiPrefix != "" {
				result = append(result, c.Style.MultiPrefix+c.Style.LinePrefix+line)
			} else {
				result = append(result, line)
			}
		}

		result = append(result, c.Style.MultiEnd)
	} else if c.Style.Single != "" {
		// Single-line comment style
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if line == "" {
				result = append(result, c.Style.Single)
			} else {
				result = append(result, c.Style.Single+c.Style.LinePrefix+line)
			}
		}
	} else if c.Style.MultiStart != "" {
		// Fallback to multi-line style for languages that only support multi-line comments
		result = append(result, c.Style.MultiStart)
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			result = append(result, line)
		}
		result = append(result, c.Style.MultiEnd)
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

var extensionStyles = map[string]CommentStyle{
	".rb":    {Language: "ruby", Single: "#", MultiStart: "=begin", MultiEnd: "=end", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".js":    {Language: "javascript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".jsx":   {Language: "javascript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".ts":    {Language: "typescript", Single: "//", MultiStart: "/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".tsx":   {Language: "typescript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".java":  {Language: "java", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".c":     {Language: "c", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".cpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".hpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".cs":    {Language: "csharp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".php":   {Language: "php", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".swift": {Language: "swift", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".rs":    {Language: "rust", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: false},
	".sh":    {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".bash":  {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".zsh":   {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yml":   {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yaml":  {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".py":    {Language: "python", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".go":    {Language: "go", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".html":  {Language: "html", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
	".css":   {Language: "css", Single: "", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".md":    {Language: "markdown", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
}

func getCommentStyle(filename string) styles.CommentStyle {
	ext := strings.ToLower(filepath.Ext(filename))
	if style, ok := extensionStyles[ext]; ok {
		return style
	}

	// Default to no comments for unknown file types
	return styles.CommentStyle{
		Language: "text",
	}
}

func uncommentContent(content string, style styles.CommentStyle) string {
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
			start := strings.Index(line, markerStart)
			end := strings.Index(line, markerEnd) + len(markerEnd)
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

// ParseLicenseComponents extracts the header, body, and footer from a license block.
// Returns the extracted components and a success flag.
func ParseLicenseComponents(content string) (header, body, footer string, success bool) {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return "", "", "", false
	}

	// Find the first non-empty line after comment start
	startIdx := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "/*") && !strings.HasPrefix(line, "<!--") {
			startIdx = i
			break
		}
	}

	// Find the last non-empty line before comment end
	endIdx := len(lines) - 1
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasSuffix(line, "*/") && !strings.HasSuffix(line, "-->") {
			endIdx = i
			break
		}
	}

	if startIdx >= endIdx {
		return "", "", "", false
	}

	// Extract header, body, and footer
	header = strings.TrimSpace(lines[startIdx])
	footer = strings.TrimSpace(lines[endIdx])

	// Extract body (everything between header and footer)
	bodyLines := lines[startIdx+1 : endIdx]
	body = strings.TrimSpace(strings.Join(bodyLines, "\n"))

	return header, body, footer, true
}

// Helper functions for working with markers
func stripMarkers(line string) string {
	line = strings.ReplaceAll(line, markerStart, "")
	line = strings.ReplaceAll(line, markerEnd, "")
	return line
}

func hasMarkers(text string) bool {
	return strings.Contains(text, markerStart) && strings.Contains(text, markerEnd)
}

func addMarkers(text string) string {
	if hasMarkers(text) {
		return text
	}
	return markerStart + text + markerEnd
}

func NewComment(style styles.CommentStyle, header, body, footer string) *Comment {
	return &Comment{
		Style:  style,
		Header: header,
		Body:   body,
		Footer: footer,
	}
}

func (c *Comment) Clone() *Comment {
	return &Comment{
		Style:  c.Style,
		Header: c.Header,
		Body:   c.Body,
		Footer: c.Footer,
	}
}

func (c *Comment) SetBody(body string) {
	c.Body = body
}

func (c *Comment) SetStyle(style styles.CommentStyle) {
	c.Style = style
}

func (c *Comment) SetHeaderAndFooterStyle(styleName string) {
	style := styles.Get(styleName)
	c.Header = style.Header
	c.Footer = style.Footer
}
