package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// YAMLHandler implements YAML-specific license handling
type YAMLHandler struct {
	*GenericHandler
}

func NewYAMLHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *YAMLHandler {
	h := &YAMLHandler{GenericHandler: NewGenericHandler(logger, style, ".yml")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *YAMLHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var directives []string
	var i int

	for i = 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "%") || trimmed == "---" {
			directives = append(directives, lines[i])
		} else {
			break
		}
	}

	if len(directives) > 0 {
		return strings.Join(directives, "\n"), strings.Join(lines[i:], "\n")
	}
	return "", content
}
