package language

import (
	"strings"

	"license-manager/internal/logger"
	"license-manager/internal/styles"
)

// LanguageHandler defines the interface for language-specific license formatting
type LanguageHandler interface {
	// FormatLicense formats the license text according to language conventions
	FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) string
	// PreservePreamble extracts and preserves any language-specific preamble (e.g., shebang, package declaration)
	PreservePreamble(content string) (preamble, rest string)
	// ScanBuildDirectives scans the content and returns the build directives and where they end
	ScanBuildDirectives(content string) (directives []string, endIndex int)
}

// GenericHandler provides default license formatting
type GenericHandler struct {
	style  styles.HeaderFooterStyle
	logger *logger.Logger
}

func NewGenericHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *GenericHandler {
	return &GenericHandler{
		style:  style,
		logger: logger,
	}
}

func (h *GenericHandler) FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) string {
	header := strings.TrimSpace(style.Header)
	footer := strings.TrimSpace(style.Footer)

	var result []string

	// Add the header marker
	if commentStyle.PreferMulti && commentStyle.MultiStart != "" {
		result = append(result, commentStyle.MultiStart)
		result = append(result, " * "+header)
		for _, line := range strings.Split(license, "\n") {
			if line == "" {
				result = append(result, " *")
			} else {
				result = append(result, " * "+line)
			}
		}
		result = append(result, " * "+footer)
		result = append(result, " "+commentStyle.MultiEnd)
	} else if commentStyle.Single != "" {
		result = append(result, commentStyle.Single+" "+header)
		for _, line := range strings.Split(license, "\n") {
			if line == "" {
				result = append(result, commentStyle.Single)
			} else {
				result = append(result, commentStyle.Single+" "+line)
			}
		}
		result = append(result, commentStyle.Single+" "+footer)
	} else {
		result = append(result, header)
		result = append(result, license)
		result = append(result, footer)
	}

	return strings.Join(result, "\n")
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

// GetLanguageHandler returns the appropriate handler for a given file type
func GetLanguageHandler(logger *logger.Logger, fileType string, style styles.HeaderFooterStyle) LanguageHandler {
	switch fileType {
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
	case "shell":
		return NewShebangHandler(logger, style)
	case "kotlin":
		return NewKotlinHandler(logger, style)
	case "scala":
		return NewScalaHandler(logger, style)
	case "css":
		return NewCSSHandler(logger, style)
	case "xml", "html":
		return NewXMLHandler(logger, style) // XML can use HTML handler (both use <!-- -->)
	//case "markdown", "md":
	//	return NewHTMLHandler(logger, style) // Markdown can use HTML handler (both use <!-- -->)
	case "ini", "toml":
		return NewINIHandler(logger, style)
	case "swift":
		return NewSwiftHandler(logger, style)
	case "csharp":
		return NewCSharpHandler(logger, style)
	default:
		return NewGenericHandler(logger, style)
	}
}
