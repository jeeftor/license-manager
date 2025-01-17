package language

import (
	"license-manager/internal/styles"
	"strings"
)

// JavaScriptHandler implements JavaScript-specific license handling
type JavaScriptHandler struct {
	*GenericHandler
}

func NewJavaScriptHandler(style styles.HeaderFooterStyle) *JavaScriptHandler {
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
