package language

import (
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
	"strings"
)

// LanguageHandler defines the interface for language-specific license formatting
type LanguageHandler interface {
	// FormatLicense formats the license text according to language conventions
	FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) FullLicenseBlock

	// PreservePreamble extracts and preserves any language-specific preamble (e.g., shebang, package declaration)
	PreservePreamble(content string) (preamble, rest string)

	// ScanBuildDirectives scans the content and returns the build directives and where they end
	ScanBuildDirectives(content string) (directives []string, endIndex int)

	// ExtractComponents extracts all components from the content including preamble, license parts, and remaining content
	ExtractComponents(content string) (components ExtractedComponents, success bool)
}

// ExtractedComponents
type ExtractedComponents struct {
	Preamble         string
	Header           string
	Body             string
	Footer           string
	Rest             string
	FullLicenseBlock *FullLicenseBlock
}

type FullLicenseBlock struct {
	String string // Entire license block
	Body   string // just body
	Header string // header portion
	Footer string // footer portion
}

// CommentExtractor handles the extraction of comment blocks
type CommentExtractor struct {
	logger *logger.Logger
	style  styles.CommentLanguage
}

func NewCommentExtractor(logger *logger.Logger, style styles.CommentLanguage) *CommentExtractor {
	return &CommentExtractor{
		logger: logger,
		style:  style,
	}
}

func (ce *CommentExtractor) extractSingleLineComments(lines []string) (header string, body []string, footer string, endIndex int, success bool) {
	if ce.style.Single == "" {
		return "", nil, "", -1, false
	}

	marker := ce.style.Single

	// Look for header style in first non-empty comment
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// If first non-empty line isn't a comment, no license block exists
		if !strings.HasPrefix(trimmed, marker) {
			return "", nil, "", -1, false
		}

		// Strip marker and spaces
		content := strings.TrimPrefix(trimmed, marker)
		if ce.style.LinePrefix != "" {
			content = strings.TrimPrefix(content, ce.style.LinePrefix)
		}
		content = strings.TrimSpace(content)

		// Try to infer header style
		match := styles.Infer(content)
		if match.Score > 0 && match.IsHeader {
			header = content
			endIndex = i
			if ce.logger != nil {
				ce.logger.LogDebug("Found header style: score %.2f", match.Score)
			}
		} else {
			// First comment line isn't a header - no license block
			return "", nil, "", -1, false
		}
		break
	}

	// Collect body and look for any valid footer style
	var bodyLines []string
	for i := endIndex + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			bodyLines = append(bodyLines, "")
			continue
		}

		// If hit non-comment line, we're done - no footer found
		if !strings.HasPrefix(trimmed, marker) {
			return "", nil, "", -1, false
		}

		// Strip marker and spaces
		content := strings.TrimPrefix(trimmed, marker)
		if ce.style.LinePrefix != "" {
			content = strings.TrimPrefix(content, ce.style.LinePrefix)
		}
		content = strings.TrimSpace(content)

		// Look for any valid footer style
		match := styles.Infer(content)
		if match.Score > 0 && match.IsFooter {
			footer = content
			return header, bodyLines, footer, i, true
		}

		bodyLines = append(bodyLines, content)
	}

	// If we hit the end without finding a footer, it's not valid
	return "", nil, "", -1, false
}

func (ce *CommentExtractor) extractMultiLineComments(lines []string) (header string, body []string, footer string, endIndex int, success bool) {
	if ce.style.MultiStart == "" || ce.style.MultiEnd == "" {
		return "", nil, "", -1, false
	}

	var startIndex int
	var foundStart, headerFound bool
	var rawBodyLines []string     // Original lines with comment markers
	var trimmedBodyLines []string // Lines with markers stripped

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if headerFound {
				rawBodyLines = append(rawBodyLines, line)       // Keep original empty line
				trimmedBodyLines = append(trimmedBodyLines, "") // Clean empty line
			}
			continue
		}

		if !foundStart {
			if strings.HasPrefix(trimmed, strings.TrimSpace(ce.style.MultiStart)) {
				startIndex = i
				foundStart = true
				ce.logger.LogDebug("Found multi-line comment start at line %d", startIndex)

			}
			continue
		}

		// Check for end marker
		if strings.HasSuffix(trimmed, strings.TrimSpace(ce.style.MultiEnd)) {
			// Trim any suffix we might have
			trimmed = strings.TrimSuffix(trimmed, strings.TrimSpace(ce.style.MultiEnd))
			if len(rawBodyLines) > 0 {
				// Use last non-empty line as footer
				footer = rawBodyLines[len(rawBodyLines)-1]

				// Trim the footer
				footer = strings.TrimPrefix(footer, strings.TrimSpace(ce.style.LinePrefix))
				footer = strings.TrimPrefix(footer, strings.TrimSpace(ce.style.MultiPrefix))

				ce.logger.LogDebug("Found multi-line comment end at line %d", startIndex)
				ce.logger.LogDebug("Footer: %s", footer)
				rawBodyLines = rawBodyLines[:len(rawBodyLines)-1]
				trimmedBodyLines = trimmedBodyLines[:len(trimmedBodyLines)-1]
			}
			return header, trimmedBodyLines, footer, i, true
		}

		// Keep original line
		rawBodyLines = append(rawBodyLines, line)

		// Strip comment markers and prefixes for trimmed version
		content := trimmed
		if ce.style.MultiPrefix != "" {
			content = strings.TrimPrefix(content, strings.TrimSpace(ce.style.MultiPrefix))
		}
		if ce.style.LinePrefix != "" {
			content = strings.TrimPrefix(content, strings.TrimSpace(ce.style.LinePrefix))
		}
		content = strings.TrimSpace(content)
		trimmedBodyLines = append(trimmedBodyLines, content)

		if !headerFound {
			header = content // Use trimmed version for header
			headerFound = true
			if ce.logger != nil {
				ce.logger.LogDebug("Found license header: %s", content)
			}
			// Remove the header line from both arrays
			rawBodyLines = rawBodyLines[:len(rawBodyLines)-1]
			trimmedBodyLines = trimmedBodyLines[:len(trimmedBodyLines)-1]
			continue
		}
	}

	return "", nil, "", -1, false
}

// GenericHandler provides default license formatting
type GenericHandler struct {
	style         styles.HeaderFooterStyle
	logger        *logger.Logger
	languageStyle styles.CommentLanguage
	// the subclassHandler lets us call into subclasses w/out having to duplicate methods
	// it seems very hacky but works - i think we can re-use the existing interface
	subclassHandler LanguageHandler // New field

}

func NewGenericHandler(logger *logger.Logger, style styles.HeaderFooterStyle, extension string) *GenericHandler {
	h := &GenericHandler{
		style:         style,
		logger:        logger,
		languageStyle: styles.GetLanguageCommentStyle(extension),
	}
	h.subclassHandler = h // Default to using itself
	return h
}

func (h *GenericHandler) ExtractComponents(content string) (components ExtractedComponents, success bool) {
	components = ExtractedComponents{
		Preamble:         "",
		Header:           "",
		Body:             "",
		Footer:           "",
		Rest:             "",
		FullLicenseBlock: nil,
	}

	if content == "" {
		return components, false
	}

	// Use the subclass handler
	preamble, remainingContent := h.subclassHandler.PreservePreamble(content)
	components.Preamble = preamble
	remainingLines := strings.Split(remainingContent, "\n")

	// Create extractor
	extractor := NewCommentExtractor(h.logger, h.languageStyle)

	// Try multi-line extraction first if preferred
	if h.languageStyle.PreferMulti {
		header, bodyLines, footer, endIndex, success := extractor.extractMultiLineComments(remainingLines)
		components.Header = header
		components.Footer = footer

		if success {
			components.Body = strings.Join(bodyLines, "\n")
			if endIndex < len(remainingLines)-1 {
				components.Rest = strings.Join(remainingLines[endIndex+1:], "\n")
			}
			lb := h.FormatLicense(components.Body, h.languageStyle, h.style)
			components.FullLicenseBlock = &lb

			return components, true
		} else {
			h.logger.LogDebug("Multi-line extraction failed -- trying single-line")
		}
	}

	// Try single-line extraction
	if h.languageStyle.Single != "" {
		header, bodyLines, footer, endIndex, success := extractor.extractSingleLineComments(remainingLines)
		components.Header = header
		components.Footer = footer
		if success {
			components.Body = strings.Join(bodyLines, "\n")
			if endIndex < len(remainingLines)-1 {
				components.Rest = strings.Join(remainingLines[endIndex+1:], "\n")
			}
			lb := h.FormatLicense(components.Body, h.languageStyle, h.style)
			components.FullLicenseBlock = &lb
			return components, true
		}
	}

	// If single-line failed and multi-line wasn't preferred, try multi-line
	if !h.languageStyle.PreferMulti && h.languageStyle.MultiStart != "" {
		header, bodyLines, footer, endIndex, success := extractor.extractMultiLineComments(remainingLines)
		components.Header = header
		components.Footer = footer
		if success {
			components.Body = strings.Join(bodyLines, "\n")
			if endIndex < len(remainingLines)-1 {
				components.Rest = strings.Join(remainingLines[endIndex+1:], "\n")
			}
			lb := h.FormatLicense(components.Body, h.languageStyle, h.style)
			components.FullLicenseBlock = &lb
			return components, true
		}
	}

	components.Rest = remainingContent
	return components, false
}

func (h *GenericHandler) FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) FullLicenseBlock {
	// Use FormatComment to do the heavy lifting
	formatted := FormatComment(license, commentStyle, style)

	// Extract the parts we need for the FullLicenseBlock
	lines := strings.Split(formatted, "\n")
	var headerFormatted, bodyFormatted, footerFormatted string

	if len(lines) > 0 {
		// Find header (first line with markers)
		for _, line := range lines {
			if hasMarkers(line) {
				headerFormatted = strings.TrimSpace(stripMarkers(line))
				break
			}
		}

		// Find footer (last line with markers)
		for i := len(lines) - 1; i >= 0; i-- {
			if hasMarkers(lines[i]) {
				footerFormatted = strings.TrimSpace(stripMarkers(lines[i]))
				break
			}
		}

		// Body is everything that isn't header/footer
		var bodyLines []string
		inBody := false
		for _, line := range lines {
			if hasMarkers(line) {
				if !inBody {
					inBody = true
					continue
				}
				break
			}
			if inBody {
				bodyLines = append(bodyLines, stripMarkers(line))
			}
		}
		bodyFormatted = strings.Join(bodyLines, "\n")
	}

	return FullLicenseBlock{
		String: formatted,
		Body:   bodyFormatted,
		Header: headerFormatted,
		Footer: footerFormatted,
	}
}

func stripMarkers(text string) string {
	text = strings.ReplaceAll(text, MarkerStart, "")
	text = strings.ReplaceAll(text, MarkerEnd, "")
	return text
}

func (h *GenericHandler) PreservePreamble(content string) (string, string) {
	return "", content
}

func (h *GenericHandler) ScanBuildDirectives(content string) ([]string, int) {
	if h.logger != nil {
		h.logger.LogVerbose("Generic handler: No build directives to scan")
	}
	return nil, 0
}

func (h *GenericHandler) SetLogger(logger *logger.Logger) {
	h.logger = logger
}

// Helper function to detect footer patterns
func looksLikeFooter(line string) bool {
	line = strings.ToLower(strings.TrimSpace(line))

	// Common footer patterns
	footerPatterns := []string{
		"end license",
		"license end",
		"end of license",
		"----",
		"====",
		"####",
	}

	for _, pattern := range footerPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	// Check for repeated characters (e.g., ------, ======, ######)
	if len(line) > 4 {
		isRepeating := true
		char := line[0]
		for i := 1; i < len(line); i++ {
			if line[i] != char {
				isRepeating = false
				break
			}
		}
		if isRepeating {
			return true
		}
	}

	return false
}

// GetLanguageHandler returns the appropriate handler for a given file type
func GetLanguageHandler(logger *logger.Logger,
	fileType string,
	style styles.HeaderFooterStyle) LanguageHandler {

	switch strings.TrimPrefix(strings.ToLower(fileType), ".") {
	case "go":
		return NewGoHandler(logger, style)
	case "js", "jsx":
		return NewJavaScriptHandler(logger, style)
	case "ts", "tsx":
		return NewTypeScriptHandler(logger, style)
	case "yaml", "yml":
		return NewYAMLHandler(logger, style)
	case "python", "py":
		return NewPythonHandler(logger, style)
	case "cpp", "c", "h", "hpp":
		return NewCppHandler(logger, style)
	case "php":
		return NewPHPHandler(logger, style)
	case "rb", "ruby":
		return NewRubyHandler(logger, style)
	case "lua":
		return NewLuaHandler(logger, style)
	case "rs":
		return NewRustHandler(logger, style)
	case "shell", "sh", "bash":
		return NewShebangHandler(logger, style)
	case "kotlin", "kt":
		return NewKotlinHandler(logger, style)
	case "scala":
		return NewScalaHandler(logger, style)
	case "css":
		return NewCSSHandler(logger, style)
	case "xml", "html", "htm":
		return NewXMLHandler(logger, style)
	case "ini", "toml":
		return NewINIHandler(logger, style)
	case "swift":
		return NewSwiftHandler(logger, style)
	case "csharp", "cs":
		return NewCSharpHandler(logger, style)
	default:
		logger.LogWarning("Unknown file type for language handler: %s", fileType)
		return NewGenericHandler(logger, style, "GENERIC")
	}
}
