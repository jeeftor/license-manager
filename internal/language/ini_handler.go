package language

import (
	"license-manager/internal/logger"
	"strings"

	"license-manager/internal/styles"
)

// INIHandler implements INI/TOML-specific license handling
type INIHandler struct {
	*GenericHandler
}

func NewINIHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *INIHandler {
	return &INIHandler{GenericHandler: NewGenericHandler(logger, style)}
}

func (h *INIHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve any initial metadata or version specifiers
		// For TOML, this includes version specifiers
		if inPreamble && (strings.HasPrefix(trimmed, ";") || strings.HasPrefix(trimmed, "#")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
