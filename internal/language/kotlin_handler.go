package language

import (
	"license-manager/internal/logger"
	"strings"

	"license-manager/internal/styles"
)

// KotlinHandler implements Kotlin-specific license handling
type KotlinHandler struct {
	*GenericHandler
}

func NewKotlinHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *KotlinHandler {
	h := &KotlinHandler{GenericHandler: NewGenericHandler(logger, style, ".kt")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *KotlinHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve package declaration, imports, file annotations, and multiplatform annotations
		if inPreamble && (strings.HasPrefix(trimmed, "package ") ||
			strings.HasPrefix(trimmed, "import ") ||
			strings.HasPrefix(trimmed, "@file:") ||
			strings.HasPrefix(trimmed, "@Target") ||
			strings.HasPrefix(trimmed, "@Retention") ||
			strings.HasPrefix(trimmed, "@OptIn") ||
			strings.HasPrefix(trimmed, "@Suppress") ||
			strings.HasPrefix(trimmed, "@JvmName") ||
			strings.HasPrefix(trimmed, "@kotlin.") ||
			// Multiplatform annotations
			strings.HasPrefix(trimmed, "@SharedImmutable") ||
			strings.HasPrefix(trimmed, "@ThreadLocal") ||
			strings.HasPrefix(trimmed, "@NativeThread") ||
			strings.HasPrefix(trimmed, "@MainThread") ||
			strings.HasPrefix(trimmed, "@WorkerThread") ||
			strings.HasPrefix(trimmed, "@UIThread") ||
			strings.HasPrefix(trimmed, "@AnyThread") ||
			// Common multiplatform expect/actual declarations
			strings.HasPrefix(trimmed, "expect ") ||
			strings.HasPrefix(trimmed, "actual ")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

func (h *KotlinHandler) ScanBuildDirectives(content string) ([]string, int) {
	lines := strings.Split(content, "\n")
	var directives []string
	lastDirectiveIndex := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check for Kotlin-specific compiler directives and annotations
		if strings.HasPrefix(trimmed, "@file:JvmName") ||
			strings.HasPrefix(trimmed, "@file:JvmMultifileClass") ||
			strings.HasPrefix(trimmed, "@file:OptIn") ||
			strings.HasPrefix(trimmed, "@file:Suppress") {
			directives = append(directives, line)
			lastDirectiveIndex = i
		}
	}

	return directives, lastDirectiveIndex + 1
}
