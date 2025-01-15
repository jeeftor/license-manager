// internal/processor/license.go
package processor

import (
	"bytes"
	"strings"
)

type LicenseManager struct {
	style        HeaderFooterStyle // Replace separate header/footer with HeaderFooterStyle
	licenseText  string
	commentStyle CommentStyle
}

func NewLicenseManager(style HeaderFooterStyle, licenseText string, commentStyle CommentStyle) *LicenseManager {
	return &LicenseManager{
		style:        style,
		licenseText:  licenseText,
		commentStyle: commentStyle,
	}
}

// formatLicenseBlock formats the license text with appropriate comment styles
func (lm *LicenseManager) formatLicenseBlock(content string) string {
	var lines []string

	// Special handling for Go files
	if lm.commentStyle.FileType == "go" {
		return lm.formatGoLicenseBlock(content)
	}

	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// Use multi-line comment style
		lines = append(lines, lm.commentStyle.MultiStart)
		lines = append(lines, stripMarkers(lm.style.Header))
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			lines = append(lines, line)
		}
		lines = append(lines, "")
		lines = append(lines, stripMarkers(lm.style.Footer))
		lines = append(lines, lm.commentStyle.MultiEnd)
	} else if lm.commentStyle.Single != "" {
		// Use single-line comment style
		lines = append(lines, lm.commentStyle.Single+" "+stripMarkers(lm.style.Header))
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			if line == "" {
				lines = append(lines, lm.commentStyle.Single)
			} else {
				lines = append(lines, lm.commentStyle.Single+" "+line)
			}
		}
		lines = append(lines, "")
		lines = append(lines, lm.commentStyle.Single+" "+stripMarkers(lm.style.Footer))
	} else {
		// HTML/XML-style or similar
		lines = append(lines, lm.commentStyle.MultiStart)
		lines = append(lines, stripMarkers(lm.style.Header))
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			lines = append(lines, line)
		}
		lines = append(lines, "")
		lines = append(lines, stripMarkers(lm.style.Footer))
		lines = append(lines, lm.commentStyle.MultiEnd)
	}

	return strings.Join(lines, "\n")
}

// formatGoLicenseBlock handles special cases for Go files like build tags
func (lm *LicenseManager) formatGoLicenseBlock(content string) string {
	var lines []string

	// For Go files, we always use the "//" style for license headers
	lines = append(lines, "// "+stripMarkers(lm.style.Header))
	lines = append(lines, "//")
	for _, line := range strings.Split(lm.licenseText, "\n") {
		if line == "" {
			lines = append(lines, "//")
		} else {
			lines = append(lines, "// "+line)
		}
	}
	lines = append(lines, "//")
	lines = append(lines, "// "+stripMarkers(lm.style.Footer))

	return strings.Join(lines, "\n")
}

// AddLicense adds the license text to the content
func (lm *LicenseManager) AddLicense(content string) string {
	// If content already has a license, don't add another one
	if lm.CheckLicense(content) {
		return content
	}

	var buf bytes.Buffer
	if lm.commentStyle.FileType == "go" {
		// Handle build tags for Go files
		lines := strings.Split(content, "\n")
		buildTagsEnd := 0

		// Find where build tags end
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") &&
				(strings.Contains(trimmed, "+build") || strings.Contains(trimmed, "go:build")) {
				buildTagsEnd = i + 1
				// Skip the required empty line after build tags
				if len(lines) > buildTagsEnd && len(strings.TrimSpace(lines[buildTagsEnd])) == 0 {
					buildTagsEnd++
				}
			} else if buildTagsEnd > 0 && len(trimmed) == 0 {
				// Keep the empty line after build tags
				buildTagsEnd = i + 1
			} else if buildTagsEnd > 0 {
				// We've found the end of the build tags section
				break
			} else if len(trimmed) > 0 {
				// No build tags found, stop looking
				break
			}
		}

		// Write any build tags first
		if buildTagsEnd > 0 {
			buf.WriteString(strings.Join(lines[:buildTagsEnd], "\n"))
			buf.WriteString("\n")
		}

		// Add the license
		buf.WriteString(lm.formatLicenseBlock(content))
		buf.WriteString("\n\n")

		// Write the rest of the file
		if buildTagsEnd > 0 {
			buf.WriteString(strings.Join(lines[buildTagsEnd:], "\n"))
		} else {
			buf.WriteString(content)
		}
	} else {
		buf.WriteString(lm.formatLicenseBlock(content))
		buf.WriteString("\n\n")
		buf.WriteString(content)
	}

	return buf.String()
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// Try to find the complete license block
	formattedLicense := lm.formatLicenseBlock("")
	if strings.Contains(content, formattedLicense) {
		return strings.Replace(content, formattedLicense+"\n\n", "", 1)
	}

	// If exact match not found, try to find the block between header and footer
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// For multi-line comments, find the block including comment markers
		start := strings.Index(content, lm.commentStyle.MultiStart)
		if start != -1 {
			end := strings.Index(content[start:], lm.commentStyle.MultiEnd)
			if end != -1 {
				end += start + len(lm.commentStyle.MultiEnd)
				// Remove the license block and any following whitespace
				remainder := content[end:]
				return strings.TrimLeft(remainder, "\n\r\t ")
			}
		}
	} else {
		// For single-line comments, find the block between header and footer
		headerLine := lm.commentStyle.Single + " " + stripMarkers(lm.style.Header)
		footerLine := lm.commentStyle.Single + " " + stripMarkers(lm.style.Footer)

		start := strings.Index(content, headerLine)
		if start != -1 {
			end := strings.Index(content[start:], footerLine)
			if end != -1 {
				end += start + len(footerLine)
				// Remove the license block and any following whitespace
				remainder := content[end:]
				return strings.TrimLeft(remainder, "\n\r\t ")
			}
		}
	}

	return content
}

// UpdateLicense updates the existing license text with new content
func (lm *LicenseManager) UpdateLicense(content string) string {
	// First remove the existing license
	content = lm.RemoveLicense(content)
	// Then add the new license
	return lm.AddLicense(content)
}

// CheckLicense verifies if the content contains the license text
func (lm *LicenseManager) CheckLicense(content string) bool {
	// First check for exact formatted license
	formattedLicense := lm.formatLicenseBlock("")
	if strings.Contains(content, formattedLicense) {
		return true
	}

	// Handle different comment styles
	if lm.commentStyle.FileType == "go" {
		headerLine := "// " + stripMarkers(lm.style.Header)
		footerLine := "// " + stripMarkers(lm.style.Footer)

		hasHeader := strings.Contains(content, headerLine)
		hasFooter := strings.Contains(content, footerLine)

		// Check if license text exists with comment prefixes
		licenseLines := strings.Split(lm.licenseText, "\n")
		for _, line := range licenseLines {
			if line != "" && !strings.Contains(content, "// "+line) {
				return false
			}
		}

		return hasHeader && hasFooter
	} else if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// For multi-line comments
		hasStart := strings.Contains(content, lm.commentStyle.MultiStart)
		hasEnd := strings.Contains(content, lm.commentStyle.MultiEnd)
		hasHeader := strings.Contains(content, stripMarkers(lm.style.Header))
		hasFooter := strings.Contains(content, stripMarkers(lm.style.Footer))

		// Check if license text exists between the multi-line comments
		if hasStart && hasEnd && hasHeader && hasFooter {
			start := strings.Index(content, lm.commentStyle.MultiStart)
			end := strings.Index(content[start:], lm.commentStyle.MultiEnd) + start
			block := content[start:end]
			return strings.Contains(block, lm.licenseText)
		}
	} else {
		// For single-line comments
		headerLine := lm.commentStyle.Single + " " + stripMarkers(lm.style.Header)
		footerLine := lm.commentStyle.Single + " " + stripMarkers(lm.style.Footer)

		hasHeader := strings.Contains(content, headerLine)
		hasFooter := strings.Contains(content, footerLine)

		// Check if license text exists with comment prefixes
		licenseLines := strings.Split(lm.licenseText, "\n")
		for _, line := range licenseLines {
			if line != "" && !strings.Contains(content, lm.commentStyle.Single+" "+line) {
				return false
			}
		}

		return hasHeader && hasFooter
	}

	return false
}
