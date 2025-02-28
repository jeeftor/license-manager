package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// TypeScriptHandler extends JavaScript handler with additional TypeScript-specific features
type TypeScriptHandler struct {
	*JavaScriptHandler
}

func NewTypeScriptHandler(
	logger *logger.Logger,
	style styles.HeaderFooterStyle,
) *TypeScriptHandler {
	h := &TypeScriptHandler{JavaScriptHandler: NewJavaScriptHandler(logger, style)}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *TypeScriptHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skips empty lines at start
		if trimmed == "" {
			if len(preamble) > 0 || len(rest) == 0 {
				preamble = append(preamble, line)
			}
			continue
		}

		// Handle TypeScript-specific directives
		if strings.HasPrefix(trimmed, "/// <reference") ||
			strings.HasPrefix(trimmed, "@ts-") ||
			strings.HasPrefix(trimmed, "// @ts-") {
			preamble = append(preamble, line)
			continue
		}

		// If not a TypeScript directive, delegate to JavaScript handler
		if len(rest) == 0 {
			jsPreample, jsRest := h.JavaScriptHandler.PreservePreamble(
				strings.Join(lines[i:], "\n"),
			)
			if jsPreample != "" {
				preamble = append(preamble, jsPreample)
			}
			if jsRest != "" {
				rest = strings.Split(jsRest, "\n")
			}
			break
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
