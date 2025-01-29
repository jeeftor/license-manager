package language

import (
	"github.com/jeeftor/license-manager/internal/logger"
	"strings"

	"github.com/jeeftor/license-manager/internal/styles"
)

// CSharpHandler implements C#-specific license handling
type CSharpHandler struct {
	*GenericHandler
}

func NewCSharpHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *CSharpHandler {
	h := &CSharpHandler{GenericHandler: NewGenericHandler(logger, style, "cs")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *CSharpHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve using statements, assembly attributes, and preprocessor directives
		if inPreamble && (strings.HasPrefix(trimmed, "using ") ||
			strings.HasPrefix(trimmed, "[assembly:") ||
			strings.HasPrefix(trimmed, "#if") ||
			strings.HasPrefix(trimmed, "#else") ||
			strings.HasPrefix(trimmed, "#endif") ||
			strings.HasPrefix(trimmed, "#define") ||
			strings.HasPrefix(trimmed, "#undef") ||
			strings.HasPrefix(trimmed, "#region") ||
			strings.HasPrefix(trimmed, "#endregion") ||
			strings.HasPrefix(trimmed, "#pragma")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

func (h *CSharpHandler) ScanBuildDirectives(content string) ([]string, int) {
	lines := strings.Split(content, "\n")
	var directives []string
	lastDirectiveIndex := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#if") ||
			strings.HasPrefix(trimmed, "#else") ||
			strings.HasPrefix(trimmed, "#endif") ||
			strings.HasPrefix(trimmed, "#define") ||
			strings.HasPrefix(trimmed, "#undef") ||
			strings.HasPrefix(trimmed, "#region") ||
			strings.HasPrefix(trimmed, "#endregion") ||
			strings.HasPrefix(trimmed, "#pragma") {
			directives = append(directives, line)
			lastDirectiveIndex = i
		}
	}

	return directives, lastDirectiveIndex + 1
}
