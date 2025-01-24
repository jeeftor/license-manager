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
	// SetLogger sets the logger for verbose output
	SetLogger(logger *logger.Logger)
}

// GenericHandler provides default license formatting
type GenericHandler struct {
	style  styles.HeaderFooterStyle
	logger *logger.Logger
}

func NewGenericHandler(style styles.HeaderFooterStyle) *GenericHandler {
	return &GenericHandler{
		style:  style,
		logger: nil,
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
func GetLanguageHandler(fileType string, style styles.HeaderFooterStyle) LanguageHandler {
	switch fileType {
	case "go":
		return NewGoHandler(style)
	case "js", "jsx":
		return NewJavaScriptHandler(style)
	case "ts", "tsx":
		return NewTypeScriptHandler(style)
	case "yaml":
		return NewYAMLHandler(style)
	case "python", "py":
		return NewPythonHandler(style)
	case "cpp", "c", "h", "hpp":
		return NewCppHandler(style)
	case "php":
		return NewPHPHandler(style)
	case "rb", "ruby":
		return NewRubyHandler(style)
	case "lua":
		return NewLuaHandler(style)
	case "rs":
		return NewRustHandler(style)
	case "shell":
		return NewShebangHandler(style)
	case "kotlin":
		return NewKotlinHandler(style)
	case "scala":
		return NewScalaHandler(style)
	case "css":
		return NewCSSHandler(style)
	case "xml", "html":
		return NewXMLHandler(style) // XML can use HTML handler (both use <!-- -->)
	//case "markdown", "md":
	//	return NewHTMLHandler(style) // Markdown can use HTML handler (both use <!-- -->)
	case "ini", "toml":
		return NewINIHandler(style)
	case "swift":
		return NewSwiftHandler(style)
	case "csharp":
		return NewCSharpHandler(style)
	default:
		return NewGenericHandler(style)
	}
}
