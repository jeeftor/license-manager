package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// PHPHandler implements PHP-specific license handling
type PHPHandler struct {
	*GenericHandler
}

func NewPHPHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *PHPHandler {
	h := &PHPHandler{GenericHandler: NewGenericHandler(logger, style, "php")}
	h.GenericHandler.subclassHandler = h
	return h
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
