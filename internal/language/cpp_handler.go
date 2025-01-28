package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"regexp"
	"strings"
)

// CppHandler implements C/C++-specific license handling
type CppHandler struct {
	*GenericHandler
}

func NewCppHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *CppHandler {
	return &CppHandler{GenericHandler: NewGenericHandler(logger, style, "cpp")}
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
