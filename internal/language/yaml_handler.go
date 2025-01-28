package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// YAMLHandler implements YAML-specific license handling
type YAMLHandler struct {
	*GenericHandler
}

func NewYAMLHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *YAMLHandler {
	return &YAMLHandler{GenericHandler: NewGenericHandler(logger, style)}
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
