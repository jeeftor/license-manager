package language

import (
	"license-manager/internal/styles"
	"regexp"
	"strings"
)

// PythonHandler implements Python-specific license handling
type PythonHandler struct {
	*GenericHandler
}

func NewPythonHandler(style styles.HeaderFooterStyle) *PythonHandler {
	return &PythonHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *PythonHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	seenShebang := false
	seenEncoding := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for shebang in first line
		if i == 0 && strings.HasPrefix(trimmed, "#!") {
			preamble = append(preamble, line)
			seenShebang = true
			continue
		}

		// Check for encoding declaration (must be in first two lines)
		if (i == 0 || (i == 1 && seenShebang)) &&
			strings.Contains(trimmed, "coding:") {
			preamble = append(preamble, line)
			seenEncoding = true
			continue
		}

		// If we've seen either directive, add rest of file
		if seenShebang || seenEncoding {
			rest = lines[i:]
			break
		}

		// If we haven't seen any directives by line 2, no preamble
		if i > 1 {
			return "", content
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

// CppHandler implements C/C++-specific license handling
type CppHandler struct {
	*GenericHandler
}

func NewCppHandler(style styles.HeaderFooterStyle) *CppHandler {
	return &CppHandler{GenericHandler: NewGenericHandler(style)}
}

func (h *CppHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	guardPattern := regexp.MustCompile(`^#ifndef\s+\w+$`)
	definePattern := regexp.MustCompile(`^#define\s+\w+$`)

	// Look for include guards
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Match #ifndef guard
		if guardPattern.MatchString(trimmed) {
			// Look for matching #define on next line
			if i+1 < len(lines) && definePattern.MatchString(strings.TrimSpace(lines[i+1])) {
				preamble = append(preamble, line, lines[i+1])
				rest = lines[i+2:]
				break
			}
		}

		// If we see any other non-whitespace content, no guards found
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			return "", content
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
