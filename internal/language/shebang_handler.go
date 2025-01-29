package language

import (
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
	"strings"
)

// ShebangHandler implements shebang-specific license handling
type ShebangHandler struct {
	*GenericHandler
}

func NewShebangHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *ShebangHandler {
	h := &ShebangHandler{
		GenericHandler: NewGenericHandler(logger, style, "sh"),
	}
	h.GenericHandler.subclassHandler = h // Set ShebangHandler as the preamble handler
	return h
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
