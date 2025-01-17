package language

import (
	"license-manager/internal/styles"
	"strings"
)

// PHPHandler implements PHP-specific license handling
type PHPHandler struct {
	*GenericHandler
}

func NewPHPHandler(style styles.HeaderFooterStyle) *PHPHandler {
	return &PHPHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *PHPHandler) PreservePreamble(content string) (string, string) {
	// Look for opening PHP tag
	idx := strings.Index(content, "<?php")
	if idx == -1 {
		return "", content
	}

	// Include the opening tag in preamble
	return content[:idx+5], strings.TrimSpace(content[idx+5:])
}
