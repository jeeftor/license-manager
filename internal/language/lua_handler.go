package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// LuaHandler implements Lua-specific license handling
type LuaHandler struct {
	*GenericHandler
}

func NewLuaHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *LuaHandler {
	return &LuaHandler{GenericHandler: NewGenericHandler(logger, style, "lua")}
}

func (h *LuaHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for shebang in first line
		if i == 0 && strings.HasPrefix(trimmed, "#!") {
			preamble = append(preamble, line)
			continue
		}

		// Check for module declarations
		if strings.HasPrefix(trimmed, "module") ||
			strings.HasPrefix(trimmed, "require") {
			preamble = append(preamble, line)
			continue
		}

		// If we hit other content, we're done
		if trimmed != "" {
			rest = lines[i:]
			break
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}
