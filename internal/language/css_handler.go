package language

import (
	"strings"

	"license-manager/internal/styles"
)

// CSSHandler implements CSS-specific license handling
type CSSHandler struct {
	*GenericHandler
}

func NewCSSHandler(style styles.HeaderFooterStyle) *CSSHandler {
	return &CSSHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *CSSHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check for at-rules that should be preserved at the top
		if inPreamble && (strings.HasPrefix(trimmed, "@charset") ||
			strings.HasPrefix(trimmed, "@import") ||
			strings.HasPrefix(trimmed, "@namespace")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
