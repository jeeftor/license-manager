package language

import (
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// PythonHandler implements Python-specific license handling
type PythonHandler struct {
	*GenericHandler
}

func NewPythonHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *PythonHandler {
	h := &PythonHandler{
		GenericHandler: NewGenericHandler(logger, style, "py"),
	}
	h.GenericHandler.subclassHandler = h
	return h
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

func (h *PythonHandler) FormatLicense(
	license string,
	commentStyle styles.CommentLanguage,
	style styles.HeaderFooterStyle,
) FullLicenseBlock {
	// First try to detect if there's already a comment style
	lines := strings.Split(license, "\n")
	hasTripleQuotes := false
	hasHashComments := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `'''`) || strings.HasPrefix(trimmed, `"""`) {
			hasTripleQuotes = true
			break
		}
		if strings.HasPrefix(trimmed, "#") {
			hasHashComments = true
		}
	}

	header := strings.TrimSpace(style.Header)
	footer := strings.TrimSpace(style.Footer)
	var result []string

	// If it's already using triple quotes, keep that style
	if hasTripleQuotes {
		return FullLicenseBlock{
			String: license,
			Header: header,
			Body:   license,
			Footer: footer,
		}
	}

	// If it's using hash comments, keep that style
	if hasHashComments {
		commentStyle.Single = "#"
		commentStyle.MultiStart = ""
		commentStyle.MultiEnd = ""
		commentStyle.LinePrefix = " "
		commentStyle.PreferMulti = false
	} else {
		// Otherwise use triple quotes by default for multi-line licenses
		commentStyle.Single = ""
		commentStyle.MultiStart = `'''`
		commentStyle.MultiEnd = `'''`
		commentStyle.MultiPrefix = "" // No prefix needed for Python triple quotes
		commentStyle.LinePrefix = ""  // No prefix needed for Python triple quotes
		commentStyle.PreferMulti = true
	}

	var bodyLines []string
	if commentStyle.PreferMulti {
		// Multi-line comment style with triple quotes
		result = append(result, commentStyle.MultiStart)
		if header != "" {
			result = append(result, header)
		}
		bodyLines = lines // Store raw body lines
		result = append(result, license)
		if footer != "" {
			result = append(result, footer)
		}
		result = append(result, commentStyle.MultiEnd)
	} else {
		// Single-line comment style with hash
		if header != "" {
			result = append(result, commentStyle.Single+commentStyle.LinePrefix+header)
		}
		for _, line := range lines {
			if line == "" {
				result = append(result, "")
				bodyLines = append(bodyLines, "")
			} else {
				result = append(result, commentStyle.Single+commentStyle.LinePrefix+line)
				bodyLines = append(bodyLines, line)
			}
		}
		if footer != "" {
			result = append(result, commentStyle.Single+commentStyle.LinePrefix+footer)
		}
	}

	return FullLicenseBlock{
		String: strings.Join(result, "\n"),
		Header: header,
		Body:   strings.Join(bodyLines, "\n"),
		Footer: footer,
	}
}
