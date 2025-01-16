package license

import (
	"strings"

	"license-manager/internal/styles"
)

const (
	markerStart = "​" // Zero-width space
	markerEnd   = "‌" // Zero-width non-joiner
)

// LicenseBlock represents a complete license block with style information
type LicenseBlock struct {
	Style  styles.CommentLanguage
	Header string
	Body   string
	Footer string
}

// String returns the complete license block as a string
func (lb *LicenseBlock) String() string {
	var result []string

	// Helper function to add markers if needed
	addMarkersIfNeeded := func(text string) string {
		if hasMarkers(text) {
			return text
		}
		return addMarkers(text)
	}

	if lb.Style.PreferMulti && lb.Style.MultiStart != "" {
		// Multi-line comment style
		result = append(result, lb.Style.MultiStart)
		result = append(result, " * "+addMarkersIfNeeded(lb.Header))

		// Add body with comment markers
		for _, line := range strings.Split(lb.Body, "\n") {
			if line == "" {
				result = append(result, " *")
			} else {
				result = append(result, " * "+line)
			}
		}

		result = append(result, " * "+addMarkersIfNeeded(lb.Footer))
		result = append(result, " "+lb.Style.MultiEnd)
	} else if lb.Style.Single != "" {
		// Single-line comment style
		result = append(result, lb.Style.Single+" "+addMarkersIfNeeded(lb.Header))

		// Add body with comment markers
		for _, line := range strings.Split(lb.Body, "\n") {
			if line == "" {
				result = append(result, lb.Style.Single)
			} else {
				result = append(result, lb.Style.Single+" "+line)
			}
		}

		result = append(result, lb.Style.Single+" "+addMarkersIfNeeded(lb.Footer))
	} else {
		// No comment style (e.g., for text files)
		result = append(result, addMarkersIfNeeded(lb.Header))
		result = append(result, lb.Body)
		result = append(result, addMarkersIfNeeded(lb.Footer))
	}

	return strings.Join(result, "\n")
}

// Helper functions for working with markers
func hasMarkers(text string) bool {
	return strings.Contains(text, markerStart) && strings.Contains(text, markerEnd)
}

func addMarkers(text string) string {
	if hasMarkers(text) {
		return text
	}
	return markerStart + text + markerEnd
}

func stripMarkers(text string) string {
	text = strings.ReplaceAll(text, markerStart, "")
	text = strings.ReplaceAll(text, markerEnd, "")
	return text
}
