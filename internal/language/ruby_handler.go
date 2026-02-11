package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
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

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for shebang
		if i == 0 && strings.HasPrefix(trimmed, "#!") {
			preamble = append(preamble, line)
			continue
		}

		// Check for magic comments (must be consecutive with shebang)
		if len(preamble) > 0 &&
			(strings.HasPrefix(trimmed, "# frozen_string_literal:") ||
				strings.HasPrefix(trimmed, "# encoding:") ||
				strings.HasPrefix(trimmed, "# warn_indent:")) {
			preamble = append(preamble, line)
			continue
		}

		// Once we've collected preamble, take the rest from here (preserving blank lines)
		if len(preamble) > 0 {
			return strings.Join(preamble, "\n"), strings.Join(lines[i:], "\n")
		}

		// No preamble found at line 0
		break
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), ""
}
