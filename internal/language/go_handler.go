package language

import (
	"license-manager/internal/styles"
	"strings"
)

type GoHandler struct {
	*GenericHandler
}

func NewGoHandler(style styles.HeaderFooterStyle) *GoHandler {
	return &GoHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *GoHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var buildTags []string
	var rest []string
	var foundBuildTag bool

	// First pass: collect build tags only
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:build") || strings.HasPrefix(trimmed, "// +build") {
			buildTags = append(buildTags, line)
			foundBuildTag = true
		} else if foundBuildTag && trimmed == "" {
			// Keep one blank line after build tags
			buildTags = append(buildTags, line)
		} else {
			rest = append(rest, line)
		}
	}

	if len(buildTags) > 0 {
		return strings.Join(buildTags, "\n"), strings.Join(rest, "\n")
	}

	// If no build tags found, return empty preamble
	return "", content
}
