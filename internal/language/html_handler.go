package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// HTMLHandler implements HTML-specific license handling
type HTMLHandler struct {
	*GenericHandler
}

func NewHTMLHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *HTMLHandler {
	h := &HTMLHandler{GenericHandler: NewGenericHandler(logger, style, "html")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *HTMLHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "<!DOCTYPE") {
			return line, strings.Join(lines[i+1:], "\n")
		}
	}
	return "", content
}
