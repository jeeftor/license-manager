package processor

import (
	"path/filepath"
	"strings"
)

const (
	// Invisible markers for license blocks
	markerStart = "\u200B" // Zero-Width Space
	markerEnd   = "\u200C" // Zero-Width Non-Joiner
)

// Comment represents a complete comment block with all its components
type Comment struct {
	// Style defines how the comment should be formatted
	Style CommentStyle

	// The actual content components
	Header string
	Body   string
	Footer string
}

// String returns the complete comment as a formatted string
func (c *Comment) String() string {
	var result []string

	// Always add markers to header and footer if they don't have them
	header := addMarkersIfNeeded(c.Header)
	footer := addMarkersIfNeeded(c.Footer)

	if c.Style.PreferMulti && c.Style.MultiStart != "" {
		// Multi-line comment style
		result = append(result, c.Style.MultiStart)
		if c.Header != "" {
			result = append(result, c.Style.MultiPrefix+header)
		}

		// Add body with comment markers
		for _, line := range strings.Split(c.Body, "\n") {
			if line == "" {
				result = append(result, c.Style.MultiPrefix)
			} else {
				result = append(result, c.Style.MultiPrefix+line)
			}
		}

		if c.Footer != "" {
			result = append(result, c.Style.MultiPrefix+footer)
		}
		result = append(result, c.Style.MultiEnd)
	} else if c.Style.Single != "" {
		// Single-line comment style
		if c.Header != "" {
			result = append(result, c.Style.Single+c.Style.LinePrefix+header)
		}

		// Add body with comment markers
		for _, line := range strings.Split(c.Body, "\n") {
			if line == "" {
				result = append(result, c.Style.Single)
			} else {
				result = append(result, c.Style.Single+c.Style.LinePrefix+line)
			}
		}

		if c.Footer != "" {
			result = append(result, c.Style.Single+c.Style.LinePrefix+footer)
		}
	} else {
		// No comment style (e.g., text files)
		if c.Header != "" {
			result = append(result, header)
		}
		result = append(result, c.Body)
		if c.Footer != "" {
			result = append(result, footer)
		}
	}

	return strings.Join(result, "\n")
}

// Parse attempts to parse a comment from the given content
func Parse(content string, style CommentStyle) (*Comment, bool) {
	// First uncomment the content
	content = uncommentContent(content, style)

	// Look for markers
	start, end := findLicenseBlock(content)
	if start == -1 || end == -1 {
		return nil, false
	}

	// Split the content into lines
	lines := strings.Split(content[start:end], "\n")
	if len(lines) < 3 { // Need at least header, body, footer
		return nil, false
	}

	// Extract header and footer (keep markers)
	header := strings.TrimSpace(lines[0])
	footer := strings.TrimSpace(lines[len(lines)-1])

	// Everything in between is the body
	body := strings.Join(lines[1:len(lines)-1], "\n")

	// Trim any extra whitespace from the body
	body = strings.TrimSpace(body)

	return &Comment{
		Style:  style,
		Header: header,
		Body:   body,
		Footer: footer,
	}, true
}

// CommentStyle represents how comments should be formatted for a specific language
type CommentStyle struct {
	// Language identifier (e.g., "go", "python", "javascript")
	// This is also used as the file type identifier
	Language string

	// Single-line comment prefix (e.g., "//", "#", "--")
	Single string

	// Multi-line comment markers
	MultiStart  string // e.g., "/*"
	MultiEnd    string // e.g., "*/"
	MultiPrefix string // e.g., " * " for multi-line comment body
	LinePrefix  string // e.g., " " for single-line comments

	// Whether to prefer multi-line comments over single-line
	PreferMulti bool
}

// Helper function to add markers if they're not already present
func addMarkersIfNeeded(text string) string {
	if hasMarkers(text) {
		return text
	}
	return addMarkers(text)
}

var extensionStyles = map[string]CommentStyle{
	".rb":    {Language: "ruby", Single: "#", MultiStart: "=begin", MultiEnd: "=end", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".js":    {Language: "javascript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: true},
	".jsx":   {Language: "javascript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: true},
	".ts":    {Language: "typescript", Single: "//", MultiStart: "/*", MultiEnd: "*/}", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".tsx":   {Language: "typescript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".java":  {Language: "java", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".c":     {Language: "c", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".cpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".hpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".cs":    {Language: "csharp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".php":   {Language: "php", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".swift": {Language: "swift", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".rs":    {Language: "rust", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".sh":    {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".bash":  {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".zsh":   {Language: "shell", Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yml":   {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yaml":  {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".pl":    {Language: "perl", Single: "#", MultiStart: "=pod", MultiEnd: "=cut", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".pm":    {Language: "perl", Single: "#", MultiStart: "=pod", MultiEnd: "=cut", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".r":     {Language: "r", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".html":  {Language: "html", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: "", PreferMulti: true},
	".xml":   {Language: "xml", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: "", PreferMulti: false},
	".css":   {Language: "css", Single: "", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: "", PreferMulti: true},
	".scss":  {Language: "scss", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".sass":  {Language: "sass", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " * ", LinePrefix: " ", PreferMulti: false},
	".lua":   {Language: "lua", Single: "--", MultiStart: "--[[", MultiEnd: "--]]", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
}

func getCommentStyle(filename string) CommentStyle {
	ext := filepath.Ext(filename)
	switch ext {
	case ".go":
		return CommentStyle{
			Language:    "go",
			Single:      "//",
			MultiStart:  "/*",
			MultiEnd:    "*/",
			MultiPrefix: " * ",
			LinePrefix:  " ",
			PreferMulti: true,
		}
	case ".py":
		return CommentStyle{
			Language:    "python",
			Single:      "#",
			MultiStart:  "",
			MultiEnd:    "",
			MultiPrefix: "",
			LinePrefix:  " ",
			PreferMulti: false,
		}
	default:
		if style, ok := extensionStyles[ext]; ok {
			return style
		}
		// Default to C-style comments if unknown
		return CommentStyle{
			Language:    "unknown",
			Single:      "//",
			MultiStart:  "/*",
			MultiEnd:    "*/",
			MultiPrefix: " * ",
			LinePrefix:  " ",
			PreferMulti: false,
		}
	}
}

func uncommentContent(content string, style CommentStyle) string {
	// Remove single-line comments
	if style.Single != "" {
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, style.Single) {
				// Remove the comment prefix and any extra space
				line = strings.TrimSpace(strings.TrimPrefix(line, style.Single))
				
				// Preserve any unicode markers
				if hasMarkers(line) {
					start := strings.Index(line, markerStart)
					end := strings.Index(line, markerEnd) + len(markerEnd)
					markers := line[start:end]
					if start > 0 {
						// Extract the part before markers and the part after markers separately
						beforeMarkers := line[:start]
						afterMarkers := ""
						if end < len(line) {
							afterMarkers = line[end:]
						}
						line = beforeMarkers + markers + afterMarkers
					} else {
						// If markers are at the start, just append the rest of the line
						if end < len(line) {
							line = markers + line[end:]
						} else {
							line = markers
						}
					}
				}
				lines[i] = line
			}
		}
		content = strings.Join(lines, "\n")
	}

	// Remove multi-line comments
	if style.MultiStart != "" && style.MultiEnd != "" {
		// First handle special cases
		if style.Language == "javascript" && strings.HasPrefix(style.MultiStart, "{") {
			// JSX-style comments
			content = strings.TrimSpace(strings.TrimPrefix(content, style.MultiStart))
			content = strings.TrimSpace(strings.TrimSuffix(content, style.MultiEnd))
		} else {
			// Standard multi-line comments
			content = strings.TrimSpace(strings.TrimPrefix(content, style.MultiStart))
			content = strings.TrimSpace(strings.TrimSuffix(content, style.MultiEnd))
		}

		// Remove line prefixes while preserving markers
		if style.MultiPrefix != "" {
			lines := strings.Split(content, "\n")
			for i, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, style.MultiPrefix) || strings.HasPrefix(line, "*") {
					// Remove the prefix and any extra space
					if strings.HasPrefix(line, style.MultiPrefix) {
						line = strings.TrimSpace(strings.TrimPrefix(line, style.MultiPrefix))
					} else {
						line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
					}
					
					// Preserve any unicode markers
					if hasMarkers(line) {
						start := strings.Index(line, markerStart)
						end := strings.Index(line, markerEnd) + len(markerEnd)
						markers := line[start:end]
						if start > 0 {
							// Extract the part before markers and the part after markers separately
							beforeMarkers := line[:start]
							afterMarkers := ""
							if end < len(line) {
								afterMarkers = line[end:]
							}
							line = beforeMarkers + markers + afterMarkers
						} else {
							// If markers are at the start, just append the rest of the line
							if end < len(line) {
								line = markers + line[end:]
							} else {
								line = markers
							}
						}
					}
					lines[i] = line
				}
			}
			content = strings.Join(lines, "\n")
		}
	}

	return strings.TrimSpace(content)
}

func findLicenseBlock(content string) (start, end int) {
	lines := strings.Split(content, "\n")
	startLine := -1
	endLine := -1

	// Find the first line with markers (header)
	for i, line := range lines {
		if hasMarkers(line) {
			startLine = i
			break
		}
	}

	// Find the last line with markers (footer)
	for i := len(lines) - 1; i >= 0; i-- {
		if hasMarkers(lines[i]) {
			endLine = i
			break
		}
	}

	if startLine == -1 || endLine == -1 || startLine >= endLine {
		return -1, -1
	}

	// Convert line numbers to character positions
	start = 0
	for i := 0; i < startLine; i++ {
		start += len(lines[i])
		if i < len(lines)-1 { // Only add newline if not the last line
			start++ // +1 for newline
		}
	}

	end = 0
	for i := 0; i <= endLine; i++ {
		end += len(lines[i])
		if i < len(lines)-1 { // Only add newline if not the last line
			end++ // +1 for newline
		}
	}

	return start, end
}

func stripMarkers(line string) string {
	line = strings.ReplaceAll(line, markerStart, "")
	line = strings.ReplaceAll(line, markerEnd, "")
	return strings.TrimSpace(line)
}

func hasMarkers(text string) bool {
	return strings.Contains(text, markerStart) && strings.Contains(text, markerEnd)
}

func addMarkers(text string) string {
	return markerStart + text + markerEnd
}

func NewComment(style CommentStyle, header, body, footer string) *Comment {
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

func (c *Comment) SetHeader(header string) {
	c.Header = header
}

func (c *Comment) SetFooter(footer string) {
	c.Footer = footer
}

func (c *Comment) SetStyle(style CommentStyle) {
	c.Style = style
}
