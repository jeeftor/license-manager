package language

import (
	"strings"

	"license-manager/internal/styles"
)

// LanguageHandler defines the interface for language-specific license formatting
type LanguageHandler interface {
	// FormatLicense formats the license text according to language conventions
	FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) string
	// PreservePreamble extracts and preserves any language-specific preamble (e.g., shebang, package declaration)
	PreservePreamble(content string) (preamble, rest string)
}

// GenericHandler provides default license formatting
type GenericHandler struct {
	style styles.HeaderFooterStyle
}

func NewGenericHandler(style styles.HeaderFooterStyle) *GenericHandler {
	return &GenericHandler{style: style}
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

// GetLanguageHandler returns the appropriate handler for a given file type
func GetLanguageHandler(fileType string, style styles.HeaderFooterStyle) LanguageHandler {
	switch fileType {
	case "go":
		return NewGoHandler(style)
	case "html":
		return NewHTMLHandler(style)
	case "javascript":
		return NewJavaScriptHandler(style)
	case "yaml":
		return NewYAMLHandler(style)
	default:
		return NewGenericHandler(style)
	}
}
