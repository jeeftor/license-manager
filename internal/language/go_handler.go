package language

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

type GoHandler struct {
	*GenericHandler
}

func NewGoHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *GoHandler {
	h := &GoHandler{GenericHandler: NewGenericHandler(logger, style, "go")}
	h.GenericHandler.subclassHandler = h
	return h
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
	var blankLinesAfterDirectives int

	if h.logger != nil {
		h.logger.LogVerbose("Go handler: Scanning 📡️ for directives...")
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for package declaration first
		if strings.HasPrefix(trimmed, "package ") {
			if h.logger != nil {
				h.logger.LogVerbose("  Found 📦️ package declaration at line %d", i)
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
				if inBuildSection && len(directives) > 0 && directives[len(directives)-1] != "" {
					// Add a blank line only if the last line wasn’t already blank
					directives = append(directives, "")
				}
			}
			if !inBuildSection && !isGenerate {
				inBuildSection = true
			}

			directives = append(directives, line)
			lastWasDirective = true
			if h.logger != nil {
				h.logger.LogVerbose(
					"  Found 🔧 directive: %s",
					color.New(color.FgHiYellow).Sprint(line),
				)
			}
		} else if trimmed == "" {
			// Keep up to two blank lines after the last directive
			if lastWasDirective {
				if blankLinesAfterDirectives < 2 {
					directives = append(directives, line)
					blankLinesAfterDirectives++
					if h.logger != nil {
						h.logger.LogVerbose("  Found ↕ blank line after directive")
					}
				}
			} else if blankLinesAfterDirectives > 0 {
				// Stop collecting blank lines if we’re past the directive section
				return directives, i, true
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

	// If we found directives, return them with the correct end index
	if len(directives) > 0 {
		// Trim any trailing empty strings beyond two
		for len(directives) > 0 && directives[len(directives)-1] == "" && blankLinesAfterDirectives > 2 {
			directives = directives[:len(directives)-1]
			blankLinesAfterDirectives--
		}
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
