package processor

import "strings"

// LanguageHandler defines the interface for language-specific license handling
type LanguageHandler interface {
	// FormatLicense formats the license text according to language conventions
	FormatLicense(license string, commentStyle CommentStyle, style HeaderFooterStyle) string
	// PreservePreamble extracts and preserves any language-specific preamble (e.g., shebang, package declaration)
	PreservePreamble(content string) (preamble, rest string)
}

// GenericHandler implements basic license handling for most languages
type GenericHandler struct {
	style HeaderFooterStyle
}

func NewGenericHandler(style HeaderFooterStyle) *GenericHandler {
	return &GenericHandler{style: style}
}

func (h *GenericHandler) FormatLicense(license string, commentStyle CommentStyle, style HeaderFooterStyle) string {
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

// GoHandler implements Go-specific license handling
type GoHandler struct {
	*GenericHandler
}

func NewGoHandler(style HeaderFooterStyle) *GoHandler {
	return &GoHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *GoHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var buildTags []string
	var rest []string
	
	// Collect build tags
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:build") || strings.HasPrefix(trimmed, "// +build") {
			buildTags = append(buildTags, line)
		} else {
			rest = append(rest, line)
		}
	}
	
	if len(buildTags) > 0 {
		return strings.Join(buildTags, "\n"), strings.Join(rest, "\n")
	}
	return "", content
}

// HTMLHandler implements HTML-specific license handling
type HTMLHandler struct {
	*GenericHandler
}

func NewHTMLHandler(style HeaderFooterStyle) *HTMLHandler {
	return &HTMLHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *HTMLHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "<!DOCTYPE") {
			return line, strings.Join(lines[i+1:], "\n")
		}
	}
	return "", content
}

// JavaScriptHandler implements JavaScript-specific license handling
type JavaScriptHandler struct {
	*GenericHandler
}

func NewJavaScriptHandler(style HeaderFooterStyle) *JavaScriptHandler {
	return &JavaScriptHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *JavaScriptHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "'use strict'") || strings.Contains(trimmed, "\"use strict\"") {
			return line, strings.Join(lines[i+1:], "\n")
		}
	}
	return "", content
}

// YAMLHandler implements YAML-specific license handling
type YAMLHandler struct {
	*GenericHandler
}

func NewYAMLHandler(style HeaderFooterStyle) *YAMLHandler {
	return &YAMLHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *YAMLHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var directives []string
	var i int
	
	for i = 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "%") || trimmed == "---" {
			directives = append(directives, lines[i])
		} else {
			break
		}
	}
	
	if len(directives) > 0 {
		return strings.Join(directives, "\n"), strings.Join(lines[i:], "\n")
	}
	return "", content
}

// ShebangHandler implements shebang-specific license handling
type ShebangHandler struct {
	*GenericHandler
}

func NewShebangHandler(style HeaderFooterStyle) *ShebangHandler {
	return &ShebangHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *ShebangHandler) PreservePreamble(content string) (string, string) {
	lines := strings.SplitN(content, "\n", 2)
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "#!") {
		return "", content
	}
	if len(lines) == 1 {
		return lines[0], ""
	}
	return lines[0], lines[1]
}

// GetLanguageHandler returns the appropriate handler for a given file type
func GetLanguageHandler(fileType string, style HeaderFooterStyle) LanguageHandler {
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
