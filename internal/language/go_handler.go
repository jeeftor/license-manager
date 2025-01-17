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

// isDirective checks if a line is a Go directive
func (h *GoHandler) isDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "//go:") ||
		strings.HasPrefix(trimmed, "// +build") ||
		strings.HasPrefix(trimmed, "//+build")
}

// isGenerateDirective checks if a line is a Go generate directive
func (h *GoHandler) isGenerateDirective(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "//go:generate")
}

// scanDirectives scans content for Go directives and returns the directives, their end index,
// and whether any directives were found
func (h *GoHandler) scanDirectives(content string) ([]string, int, bool) {
	lines := strings.Split(content, "\n")
	var directives []string
	var lastWasDirective bool
	var inBuildSection bool
	var inGenerateSection bool

	if h.logger != nil {
		h.logger.LogVerbose("Go handler: Scanning for directives...")
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for package declaration first
		if strings.HasPrefix(trimmed, "package ") {
			if h.logger != nil {
				h.logger.LogVerbose("Found package declaration at line %d", i)
			}
			if len(directives) > 0 {
				return directives, i, true
			}
			return nil, i, false
		}

		// Check for Go directives
		if h.isDirective(line) {
			// Check if we're switching from build to generate directives
			isGenerate := h.isGenerateDirective(line)
			if !inGenerateSection && isGenerate {
				inGenerateSection = true
				if inBuildSection {
					// Add a blank line between build and generate sections
					directives = append(directives, "")
				}
			}
			if !inBuildSection && !isGenerate {
				inBuildSection = true
			}

			directives = append(directives, line)
			lastWasDirective = true
			if h.logger != nil {
				h.logger.LogVerbose("Found directive: %s", line)
			}
		} else if trimmed == "" {
			// Keep blank lines if we're still in a directive section
			if lastWasDirective {
				directives = append(directives, line)
				if h.logger != nil {
					h.logger.LogVerbose("Keeping blank line after directive")
				}
			}
		} else {
			// If we hit a non-directive line that's not a package declaration
			// and we were in a directive section, we're done
			if lastWasDirective {
				if h.logger != nil {
					h.logger.LogVerbose("Found non-directive line at %d: %s", i, trimmed)
				}
				return directives, i, true
			}
			lastWasDirective = false
		}
	}

	if h.logger != nil {
		if len(directives) > 0 {
			h.logger.LogVerbose("Reached end of file, found %d directives", len(directives))
		} else {
			h.logger.LogVerbose("No directives found in file")
		}
	}

	// If we found directives, return them
	if len(directives) > 0 {
		return directives, len(lines), true
	}
	return nil, len(lines), false
}

func (h *GoHandler) PreservePreamble(content string) (string, string) {
	directives, endIndex, found := h.scanDirectives(content)
	if !found {
		return "", strings.TrimSpace(content)
	}

	lines := strings.Split(content, "\n")
	
	// Ensure there's a blank line after directives
	if !strings.HasSuffix(strings.Join(directives, "\n"), "\n\n") {
		directives = append(directives, "")
	}

	// Return preamble and rest
	return strings.Join(directives, "\n"), strings.TrimSpace(strings.Join(lines[endIndex:], "\n"))
}

func (h *GoHandler) ScanBuildDirectives(content string) ([]string, int) {
	directives, endIndex, found := h.scanDirectives(content)
	if !found {
		return nil, endIndex
	}
	return directives, endIndex
}
