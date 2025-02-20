// python_handler.go
package language

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// unicodeEscapePattern matches \uXXXX unicode escape sequences
var unicodeEscapePattern = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)

// normalizeUnicodeEscapes converts \uXXXX sequences to their actual Unicode characters
func normalizeUnicodeEscapes(s string) string {
	return unicodeEscapePattern.ReplaceAllStringFunc(s, func(match string) string {
		// Convert the hex string to a rune
		code, _ := strconv.ParseInt(match[2:], 16, 32)
		return string(rune(code))
	})
}

// stripUnicodeZeroWidth removes zero-width spaces that might be added by formatters
func stripUnicodeZeroWidth(s string) string {
	s = strings.ReplaceAll(s, MarkerStart, "")
	s = strings.ReplaceAll(s, MarkerEnd, "")
	return s
}

// normalizeText normalizes both Unicode escapes and zero-width spaces
func normalizeText(s string) string {
	s = normalizeUnicodeEscapes(s)
	return strings.TrimSpace(s) // Always trim space after normalization
}

// extractTripleQuotedBlock extracts the content between triple quotes
func extractTripleQuotedBlock(content string) string {
	lines := strings.Split(content, "\n")
	var blockLines []string
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `"""`) || strings.HasPrefix(trimmed, `'''`) {
			if !inBlock {
				inBlock = true
				continue
			} else {
				break
			}
		}
		if inBlock {
			blockLines = append(blockLines, line)
		}
	}

	return strings.Join(blockLines, "\n")
}

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

// FormatLicense formats the license text according to Python conventions
func (h *PythonHandler) FormatLicense(
	license string,
	commentStyle styles.CommentLanguage,
	style styles.HeaderFooterStyle,
) FullLicenseBlock {
	// Normalize any Unicode escapes in the license text and style markers
	normalizedLicense := normalizeText(license)
	header := normalizeText(stripMarkers(style.Header))
	footer := normalizeText(stripMarkers(style.Footer))

	// Format the license block with triple quotes
	var result []string
	result = append(result, "'''")
	result = append(result, " * "+header)

	// Add the license body
	lines := strings.Split(normalizedLicense, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			result = append(result, " *")
		} else {
			result = append(result, " * "+line)
		}
	}

	result = append(result, " * "+footer)
	result = append(result, " '''")

	return FullLicenseBlock{
		String: strings.Join(result, "\n"),
		Header: header,
		Body:   normalizedLicense,
		Footer: footer,
	}
}
