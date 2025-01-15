// internal/processor/license.go
package processor

import (
	"bytes"
	"strings"
)

type LicenseManager struct {
	header      string
	footer      string
	licenseText string
	style       CommentStyle
}

func NewLicenseManager(header, footer, licenseText string, style CommentStyle) *LicenseManager {
	return &LicenseManager{
		header:      header,
		footer:      footer,
		licenseText: licenseText,
		style:       style,
	}
}

// formatLicenseBlock formats the license text with appropriate comment styles
func (lm *LicenseManager) formatLicenseBlock(content string) string {
	var lines []string

	if lm.style.PreferMulti && lm.style.MultiStart != "" {
		// Use multi-line comment style
		lines = append(lines, lm.style.MultiStart)
		lines = append(lines, lm.header)
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			lines = append(lines, line)
		}
		lines = append(lines, "")
		lines = append(lines, lm.footer)
		lines = append(lines, lm.style.MultiEnd)
	} else if lm.style.Single != "" {
		// Use single-line comment style
		lines = append(lines, lm.style.Single+" "+lm.header)
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			if line == "" {
				lines = append(lines, lm.style.Single)
			} else {
				lines = append(lines, lm.style.Single+" "+line)
			}
		}
		lines = append(lines, "")
		lines = append(lines, lm.style.Single+" "+lm.footer)
	} else {
		// HTML/XML-style or similar
		lines = append(lines, lm.style.MultiStart)
		lines = append(lines, lm.header)
		lines = append(lines, "")
		for _, line := range strings.Split(lm.licenseText, "\n") {
			lines = append(lines, line)
		}
		lines = append(lines, "")
		lines = append(lines, lm.footer)
		lines = append(lines, lm.style.MultiEnd)
	}

	return strings.Join(lines, "\n")
}

// AddLicense adds the license text to the content, respecting build tags for Go files
func (lm *LicenseManager) AddLicense(content string) string {
	// If content already has a license, don't add another one
	if lm.CheckLicense(content) {
		return content
	}

	var buf bytes.Buffer

	// Special handling for Go files
	if lm.style.FileType == "go" {
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
		buf.WriteString(lm.formatLicenseBlock(lm.licenseText))
		buf.WriteString("\n\n")

		// Write the rest of the file
		if buildTagsEnd > 0 {
			buf.WriteString(strings.Join(lines[buildTagsEnd:], "\n"))
		} else {
			buf.WriteString(content)
		}
	} else {
		// Normal handling for non-Go files
		buf.WriteString(lm.formatLicenseBlock(lm.licenseText))
		buf.WriteString("\n\n")
		buf.WriteString(content)
	}

	return buf.String()
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// Try to find the complete license block
	formattedLicense := lm.formatLicenseBlock(lm.licenseText)
	if strings.Contains(content, formattedLicense) {
		return strings.Replace(content, formattedLicense, "", 1)
	}

	// If exact match not found, try to find the block between header and footer
	if lm.style.PreferMulti && lm.style.MultiStart != "" {
		// For multi-line comments, find the block including comment markers
		start := strings.Index(content, lm.style.MultiStart)
		if start != -1 {
			end := strings.Index(content[start:], lm.style.MultiEnd)
			if end != -1 {
				end += start + len(lm.style.MultiEnd)
				// Remove the license block and any following whitespace
				remainder := content[end:]
				return strings.TrimLeft(remainder, "\n\r\t ")
			}
		}
	} else {
		// For single-line comments, find the block between header and footer
		headerLine := lm.style.Single + " " + lm.header
		footerLine := lm.style.Single + " " + lm.footer

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

	// If no license block found, return original content
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
	formattedLicense := lm.formatLicenseBlock(lm.licenseText)
	if strings.Contains(content, formattedLicense) {
		return true
	}

	// Then check for header and footer with any content between
	if lm.style.PreferMulti && lm.style.MultiStart != "" {
		// For multi-line comments
		hasStart := strings.Contains(content, lm.style.MultiStart)
		hasEnd := strings.Contains(content, lm.style.MultiEnd)
		hasHeader := strings.Contains(content, lm.header)
		hasFooter := strings.Contains(content, lm.footer)
		hasLicense := strings.Contains(content, lm.licenseText)

		return hasStart && hasEnd && hasHeader && hasFooter && hasLicense
	} else {
		// For single-line comments
		headerLine := lm.style.Single + " " + lm.header
		footerLine := lm.style.Single + " " + lm.footer

		hasHeader := strings.Contains(content, headerLine)
		hasFooter := strings.Contains(content, footerLine)

		// Check if license text exists with comment prefixes
		licenseLines := strings.Split(lm.licenseText, "\n")
		for _, line := range licenseLines {
			if line != "" && !strings.Contains(content, lm.style.Single+" "+line) {
				return false
			}
		}

		return hasHeader && hasFooter
	}
}

// FormatLicenseText returns the formatted license text without adding it to any content
// Useful for preview or verification purposes
func (lm *LicenseManager) FormatLicenseText() string {
	return lm.formatLicenseBlock(lm.licenseText)
}

// RemoveComments removes comment markers from license text
// Useful when comparing license content regardless of comment style
func (lm *LicenseManager) RemoveComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		if lm.style.Single != "" {
			line = strings.TrimPrefix(line, lm.style.Single)
		}
		if lm.style.MultiStart != "" {
			line = strings.TrimPrefix(line, lm.style.MultiStart)
		}
		if lm.style.MultiEnd != "" {
			line = strings.TrimSuffix(line, lm.style.MultiEnd)
		}
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
