package language

import (
	"strings"

	"license-manager/internal/styles"
)

// SwiftHandler implements Swift-specific license handling
type SwiftHandler struct {
	*GenericHandler
}

func NewSwiftHandler(style styles.HeaderFooterStyle) *SwiftHandler {
	return &SwiftHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *SwiftHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve compiler directives and imports at the top
		if inPreamble && (strings.HasPrefix(trimmed, "#if") ||
			strings.HasPrefix(trimmed, "#else") ||
			strings.HasPrefix(trimmed, "#endif") ||
			strings.HasPrefix(trimmed, "#elseif") ||
			strings.HasPrefix(trimmed, "import ") ||
			strings.HasPrefix(trimmed, "@_exported")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

func (h *SwiftHandler) ScanBuildDirectives(content string) ([]string, int) {
	lines := strings.Split(content, "\n")
	var directives []string
	lastDirectiveIndex := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#if") ||
			strings.HasPrefix(trimmed, "#else") ||
			strings.HasPrefix(trimmed, "#endif") ||
			strings.HasPrefix(trimmed, "#elseif") {
			directives = append(directives, line)
			lastDirectiveIndex = i
		}
	}

	return directives, lastDirectiveIndex + 1
}
