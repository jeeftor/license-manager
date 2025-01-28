package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
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

// extractSingleLineComments handles extraction from single-line comment blocks
func (ce *CommentExtractor) extractSingleLineComments(lines []string) (header string, body []string, footer string, endIndex int, success bool) {
	if ce.style.Single == "" {
		return "", nil, "", -1, false
	}

	var headerFound, inBody bool
	var bodyLines []string
	marker := ce.style.Single

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if inBody {
				bodyLines = append(bodyLines, "")
			}
			continue
		}

		if !strings.HasPrefix(trimmed, marker) {
			if inBody {
				// End of comment block
				return header, bodyLines, footer, i - 1, true
			}
			continue
		}

		// Strip comment marker and any line prefix
		content := strings.TrimPrefix(trimmed, marker)
		if ce.style.LinePrefix != "" {
			content = strings.TrimPrefix(content, ce.style.LinePrefix)
		}
		content = strings.TrimSpace(content)

		if !headerFound {
			header = content
			headerFound = true
			inBody = true
			if ce.logger != nil {
				ce.logger.LogDebug("Found license header: %s", content)
			}
			continue
		}

		// Check for footer pattern
		if looksLikeFooter(content) {
			footer = content
			if ce.logger != nil {
				ce.logger.LogDebug("Found license footer: %s", content)
			}
			return header, bodyLines, footer, i, true
		}

		if inBody {
			bodyLines = append(bodyLines, content)
		}
	}

	// If we reach here and were in a comment block, use last non-empty line as footer
	if inBody && len(bodyLines) > 0 {
		footer = bodyLines[len(bodyLines)-1]
		bodyLines = bodyLines[:len(bodyLines)-1]
		return header, bodyLines, footer, len(lines) - 1, true
	}

	return "", nil, "", -1, false
}

// extractMultiLineComments handles extraction from multi-line comment blocks
func (ce *CommentExtractor) extractMultiLineComments(lines []string) (header string, body []string, footer string, endIndex int, success bool) {
	if ce.style.MultiStart == "" || ce.style.MultiEnd == "" {
		return "", nil, "", -1, false
	}

	var startIndex int
	var foundStart, headerFound bool
	var bodyLines []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if headerFound {
				bodyLines = append(bodyLines, "")
			}
			continue
		}

		if !foundStart {
			if strings.HasPrefix(trimmed, ce.style.MultiStart) {
				startIndex = i
				foundStart = true
				if ce.logger != nil {
					ce.logger.LogDebug("Found multi-line comment start at line %d", startIndex)
				}
			}
			continue
		}

		// Check for end marker
		if strings.HasSuffix(trimmed, ce.style.MultiEnd) {
			if len(bodyLines) > 0 {
				// Use last non-empty line as footer
				footer = bodyLines[len(bodyLines)-1]
				bodyLines = bodyLines[:len(bodyLines)-1]
			}
			return header, bodyLines, footer, i, true
		}

		// Strip comment markers and prefixes
		content := trimmed
		if ce.style.MultiPrefix != "" {
			content = strings.TrimPrefix(content, ce.style.MultiPrefix)
		}
		if ce.style.LinePrefix != "" {
			content = strings.TrimPrefix(content, ce.style.LinePrefix)
		}
		content = strings.TrimSpace(content)

		if !headerFound {
			header = content
			headerFound = true
			if ce.logger != nil {
				ce.logger.LogDebug("Found license header: %s", content)
			}
			continue
		}

		bodyLines = append(bodyLines, content)
	}

	return "", nil, "", -1, false
}

// This is a workaround due to how Go does inheritence so we dont' have to implement everything in subclasses
//type SubclassHandler interface {
//	PreservePreamble(content string) (string, string)
//	FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) FullLicenseBlock
//}

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
	header := strings.TrimSpace(style.Header)
	footer := strings.TrimSpace(style.Footer)

	var result []string
	var headerFormatted, bodyFormatted, footerFormatted string

	if commentStyle.PreferMulti && commentStyle.MultiStart != "" {
		// Start the block
		result = append(result, commentStyle.MultiStart)

		// Handle header
		if header != "" {
			headerLine := commentStyle.MultiPrefix + commentStyle.LinePrefix + header
			result = append(result, headerLine)
			headerFormatted = header
		}

		// Handle body
		var bodyLines []string
		for _, line := range strings.Split(license, "\n") {
			if line == "" {
				result = append(result, commentStyle.MultiPrefix)
				bodyLines = append(bodyLines, "")
			} else {
				result = append(result, commentStyle.MultiPrefix+commentStyle.LinePrefix+line)
				bodyLines = append(bodyLines, line)
			}
		}
		bodyFormatted = strings.Join(bodyLines, "\n")

		// Handle footer
		if footer != "" {
			footerLine := commentStyle.MultiPrefix + commentStyle.LinePrefix + footer
			result = append(result, footerLine)
			footerFormatted = footer
		}

		// Close the block
		result = append(result, commentStyle.MultiEnd)

	} else if commentStyle.Single != "" {
		// Handle header
		if header != "" {
			headerLine := commentStyle.Single + commentStyle.LinePrefix + header
			result = append(result, headerLine)
			headerFormatted = header
		}

		// Handle body
		var bodyLines []string
		for _, line := range strings.Split(license, "\n") {
			if line == "" {
				result = append(result, commentStyle.Single)
				bodyLines = append(bodyLines, "")
			} else {
				result = append(result, commentStyle.Single+commentStyle.LinePrefix+line)
				bodyLines = append(bodyLines, line)
			}
		}
		bodyFormatted = strings.Join(bodyLines, "\n")

		// Handle footer
		if footer != "" {
			footerLine := commentStyle.Single + commentStyle.LinePrefix + footer
			result = append(result, footerLine)
			footerFormatted = footer
		}

	} else {
		// No comment style, store raw text
		if header != "" {
			result = append(result, header)
			headerFormatted = header
		}
		result = append(result, license)
		bodyFormatted = license
		if footer != "" {
			result = append(result, footer)
			footerFormatted = footer
		}
	}

	return FullLicenseBlock{
		String: strings.Join(result, "\n"),
		Body:   bodyFormatted,
		Header: headerFormatted,
		Footer: footerFormatted,
	}
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
	case "yaml":
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
	case "kotlin":
		return NewKotlinHandler(logger, style)
	case "scala":
		return NewScalaHandler(logger, style)
	case "css":
		return NewCSSHandler(logger, style)
	case "xml", "html":
		return NewXMLHandler(logger, style)
	case "ini", "toml":
		return NewINIHandler(logger, style)
	case "swift":
		return NewSwiftHandler(logger, style)
	case "csharp":
		return NewCSharpHandler(logger, style)
	default:
		logger.LogWarning("Unknown file type for language handler: %s", fileType)
		return NewGenericHandler(logger, style, "GENERIC")
	}
}
