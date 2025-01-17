package language

import (
	"license-manager/internal/styles"
	"strings"
)

// ShebangHandler implements shebang-specific license handling
type ShebangHandler struct {
	*GenericHandler
}

func NewShebangHandler(style styles.HeaderFooterStyle) *ShebangHandler {
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
