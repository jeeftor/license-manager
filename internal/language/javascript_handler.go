package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

//
//// JavaScriptHandler implements JavaScript-specific license handling
//type JavaScriptHandler struct {
//	*GenericHandler
//}
//
//func NewJavaScriptHandler(style styles.HeaderFooterStyle) *JavaScriptHandler {
//	return &JavaScriptHandler{GenericHandler: NewGenericHandler(style)}
//}
//
//func (h *JavaScriptHandler) PreservePreamble(content string) (string, string) {
//	lines := strings.Split(content, "\n")
//	for i, line := range lines {
//		trimmed := strings.TrimSpace(line)
//		if strings.Contains(trimmed, "'use strict'") || strings.Contains(trimmed, "\"use strict\"") {
//			return line, strings.Join(lines[i+1:], "\n")
//		}
//	}
//	return "", content
//}
//

// Enhanced JavaScriptHandler with additional features
type JavaScriptHandler struct {
	*GenericHandler
}

func NewJavaScriptHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *JavaScriptHandler {
	return &JavaScriptHandler{GenericHandler: NewGenericHandler(logger, style, ".js")}
}

func (h *JavaScriptHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	// Track if we've seen any actual code
	seenCode := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skips empty lines
		if trimmed == "" {
			if !seenCode {
				preamble = append(preamble, line)
			}
			continue
		}

		// Check for 'use strict'
		if strings.Contains(trimmed, "'use strict'") ||
			strings.Contains(trimmed, "\"use strict\"") {
			preamble = append(preamble, line)
			continue
		}

		// Check for shebang (for Node.js scripts)
		if i == 0 && strings.HasPrefix(trimmed, "#!") {
			preamble = append(preamble, line)
			continue
		}

		// Check for ES6 imports
		if strings.HasPrefix(trimmed, "import ") ||
			strings.HasPrefix(trimmed, "export ") {
			preamble = append(preamble, line)
			continue
		}

		// If we hit other content, mark that we've seen code
		seenCode = true
		rest = lines[i:]
		break
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
