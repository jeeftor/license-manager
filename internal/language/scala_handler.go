package language

import (
	"license-manager/internal/logger"
	"strings"

	"license-manager/internal/styles"
)

// ScalaHandler implements Scala-specific license handling
type ScalaHandler struct {
	*GenericHandler
}

func NewScalaHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *ScalaHandler {
	h := &ScalaHandler{GenericHandler: NewGenericHandler(logger, style, "scala")}
	h.GenericHandler.subclassHandler = h
	return h
}

func (h *ScalaHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve package declaration, imports, package objects, and annotations
		if inPreamble && (strings.HasPrefix(trimmed, "package ") ||
			strings.HasPrefix(trimmed, "package object ") ||
			strings.HasPrefix(trimmed, "import ") ||
			// Common Scala annotations
			strings.HasPrefix(trimmed, "@main") ||
			strings.HasPrefix(trimmed, "@native") ||
			strings.HasPrefix(trimmed, "@inline") ||
			strings.HasPrefix(trimmed, "@throws") ||
			strings.HasPrefix(trimmed, "@tailrec") ||
			strings.HasPrefix(trimmed, "@switch") ||
			strings.HasPrefix(trimmed, "@specialized") ||
			strings.HasPrefix(trimmed, "@transient") ||
			strings.HasPrefix(trimmed, "@volatile") ||
			// Scala 3 specific syntax
			strings.HasPrefix(trimmed, "given ") ||
			strings.HasPrefix(trimmed, "export ") ||
			strings.HasPrefix(trimmed, "transparent ") ||
			strings.HasPrefix(trimmed, "opaque ") ||
			// Common compiler directives
			strings.HasPrefix(trimmed, "import scala.language.") ||
			strings.HasPrefix(trimmed, "import scala.annotation.")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, line)
		}
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

func (h *ScalaHandler) ScanBuildDirectives(content string) ([]string, int) {
	lines := strings.Split(content, "\n")
	var directives []string
	lastDirectiveIndex := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check for Scala-specific compiler directives and language imports
		if strings.HasPrefix(trimmed, "import scala.language.") ||
			strings.HasPrefix(trimmed, "import scala.annotation.") ||
			strings.HasPrefix(trimmed, "@scala.annotation.") {
			directives = append(directives, line)
			lastDirectiveIndex = i
		}
	}

	return directives, lastDirectiveIndex + 1
}
