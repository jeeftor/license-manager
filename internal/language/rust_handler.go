package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"regexp"
	"strings"
)

// RustHandler implements Rust-specific license handling
type RustHandler struct {
	*GenericHandler
}

func NewRustHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *RustHandler {
	return &RustHandler{GenericHandler: NewGenericHandler(logger, style, ".rs")}
}

func (h *RustHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	// Regex for feature attributes and crate attributes
	featurePattern := regexp.MustCompile(`^#!\[feature\(.*\)\]$`)
	cratePattern := regexp.MustCompile(`^#!\[.*\]$`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skips empty lines at the start
		if trimmed == "" {
			preamble = append(preamble, line)
			continue
		}

		// Match feature or crate attributes
		if featurePattern.MatchString(trimmed) || cratePattern.MatchString(trimmed) {
			preamble = append(preamble, line)
			continue
		}

		// If we hit a non-attribute line, we're done
		rest = lines[i:]
		break
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
