package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// RubyHandler implements Ruby-specific license handling
type RubyHandler struct {
	*GenericHandler
}

func NewRubyHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *RubyHandler {
	h := &RubyHandler{GenericHandler: NewGenericHandler(logger, style, "rb")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *RubyHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for shebang
		if i == 0 && strings.HasPrefix(trimmed, "#!") {
			preamble = append(preamble, line)
			continue
		}

		// Check for magic comments
		if strings.HasPrefix(trimmed, "# frozen_string_literal:") ||
			strings.HasPrefix(trimmed, "# encoding:") ||
			strings.HasPrefix(trimmed, "# warn_indent:") {
			preamble = append(preamble, line)
			continue
		}

		// If we hit any other content, we're done with preamble
		if trimmed != "" {
			rest = lines[i:]
			break
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
