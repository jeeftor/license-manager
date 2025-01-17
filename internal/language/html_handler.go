package language

import (
	"license-manager/internal/styles"
	"strings"
)

// HTMLHandler implements HTML-specific license handling
type HTMLHandler struct {
	*GenericHandler
}

func NewHTMLHandler(style styles.HeaderFooterStyle) *HTMLHandler {
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
