package language

import (
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"strings"
)

// XMLHandler implements XML-specific license handling
type XMLHandler struct {
	*GenericHandler
}

func NewXMLHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *XMLHandler {
	h := &XMLHandler{
		GenericHandler: NewGenericHandler(logger, style, "xml"),
	}
	h.GenericHandler.subclassHandler = h
	return h

}

func (h *XMLHandler) PreservePreamble(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var preamble []string
	var rest []string
	inPreamble := true

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for XML declaration and DOCTYPE
		if inPreamble && (strings.HasPrefix(trimmed, "<?xml") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "<!DOCTYPE") ||
			strings.HasPrefix(trimmed, "<?xml-stylesheet")) {
			preamble = append(preamble, line)
		} else {
			inPreamble = false
			rest = append(rest, lines[i:]...)
			break
		}
	}

	if len(preamble) == 0 {
		return "", content
	}

	return strings.Join(preamble, "\n"), strings.Join(rest, "\n")
}

func (h *XMLHandler) FormatLicense(license string, commentStyle styles.CommentLanguage, style styles.HeaderFooterStyle) FullLicenseBlock {
	header := strings.TrimSpace(style.Header)
	footer := strings.TrimSpace(style.Footer)

	var result []string
	result = append(result, "<!--")
	result = append(result, header)
	result = append(result, license)
	result = append(result, footer)
	result = append(result, "-->")

	return FullLicenseBlock{
		String: strings.Join(result, "\n"),
		Header: header,
		Body:   license,
		Footer: footer,
	}
}
