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
		for _, line := range strings.Split(content, "\n") {
			lines = append(lines, line)
		}
		lines = append(lines, lm.commentStyle.MultiEnd)
	} else if lm.commentStyle.Single != "" {
		// Use single-line comment style
		for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
			if line == "" {
				lines = append(lines, lm.commentStyle.Single)
			} else {
				// Add a space after the comment character only if the line is not empty
				lines = append(lines, lm.commentStyle.Single+" "+line)
			}
		}
	} else {
		// No comment style available, use the content as-is
		return content
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// formatGoLicenseBlock handles special cases for Go files like build tags
func (lm *LicenseManager) formatGoLicenseBlock(content string) string {
	var lines []string

	// For Go files, we always use the "//" style for license headers
	for _, line := range strings.Split(content, "\n") {
		if line == "" {
			lines = append(lines, "//")
		} else {
			lines = append(lines, "// "+line)
		}
	}

	return strings.Join(lines, "\n")
}

// AddLicense adds the license text to the content
func (lm *LicenseManager) AddLicense(content string) string {
	// Check if the file already has any header/footer markers
	header := lm.style.Header
	footer := lm.style.Footer

	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerbose("Starting AddLicense process...")
		fp.logVerbose("Header to look for: %q", header)
		fp.logVerbose("Footer to look for: %q", footer)
	}

	// If either header or footer exists, skip adding the license
	if strings.Contains(content, header) || strings.Contains(content, footer) {
		if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
			fp.logVerbose("Found existing header or footer - skipping license addition")
		}
		return content
	}

	// Format the license block with comments
	var formattedLicense string
	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerbose("Raw license text:")
		fp.logVerboseWithLineNumbers(lm.licenseText, 1, "")
	}

	// Add comments to the license text
	commentedLicenseText := lm.formatLicenseBlock(lm.licenseText)
	
	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerbose("Commented license text:")
		fp.logVerboseWithLineNumbers(commentedLicenseText, 1, "")
	}
	
	// Add header and footer
	if lm.commentStyle.Single != "" {
		formattedLicense = lm.commentStyle.Single + " " + header + "\n" + 
			commentedLicenseText + "\n" + 
			lm.commentStyle.Single + " " + footer + "\n\n"
	} else {
		formattedLicense = header + "\n" + 
			commentedLicenseText + "\n" + 
			footer + "\n\n"
	}

	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerbose("Final formatted license block:")
		fp.logVerboseWithLineNumbers(formattedLicense, 1, "")
	}

	var buf bytes.Buffer

	// Special handling for Go files with build tags
	if lm.commentStyle.FileType == "go" {
		if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
			fp.logVerbose("Processing Go file - checking for build tags...")
		}

		lines := strings.Split(content, "\n")
		var buildTagsEnd int

		// Find the end of build tags section
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "//") || strings.HasPrefix(line, "// +build") {
				continue
			}
			if line == "//" || line == "" {
				buildTagsEnd = i + 1
				if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
					fp.logVerbose("Found build tags section ending at line %d", buildTagsEnd)
				}
				break
			}
		}

		// Write build tags if they exist
		if buildTagsEnd > 0 {
			buildTagsSection := strings.Join(lines[:buildTagsEnd], "\n")
			if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
				fp.logVerbose("Writing build tags section:")
				fp.logVerboseWithLineNumbers(buildTagsSection, 1, "")
			}
			buf.WriteString(buildTagsSection)
			buf.WriteString("\n\n")
		}

		// Add the license
		buf.WriteString(formattedLicense)

		// Write the rest of the file
		if buildTagsEnd > 0 {
			remainingContent := strings.Join(lines[buildTagsEnd:], "\n")
			if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
				fp.logVerbose("Writing remaining content after build tags:")
				fp.logVerboseWithLineNumbers(remainingContent, buildTagsEnd+1, "")
			}
			buf.WriteString("\n")
			buf.WriteString(remainingContent)
		} else {
			if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
				fp.logVerbose("No build tags found - writing entire content:")
				fp.logVerboseWithLineNumbers(content, 1, "")
			}
			buf.WriteString(content)
		}
	} else {
		// For non-Go files, simply prepend the license
		if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
			fp.logVerbose("Processing non-Go file - prepending license and writing content")
		}
		buf.WriteString(formattedLicense)
		buf.WriteString(content)
	}

	result := buf.String()
	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerbose("Final result:")
		fp.logVerboseWithLineNumbers(result, 1, "")
	}

	return result
}

// extractLicenseText extracts the license text between header and footer
func (lm *LicenseManager) extractLicenseText(content string) (string, bool) {
	header := lm.style.Header
	footer := lm.style.Footer

	startIdx := strings.Index(content, header)
	if startIdx == -1 {
		return "", false
	}

	// Look for footer after the header
	remainingContent := content[startIdx+len(header):]
	endIdx := strings.Index(remainingContent, footer)
	if endIdx == -1 {
		return "", false
	}

	// Extract the text between header and footer, skipping the header/footer lines
	licenseBlock := remainingContent[:endIdx]

	// For single-line comment styles, we need to clean up the extracted text
	if lm.commentStyle.Single != "" {
		var cleanedLines []string
		for _, line := range strings.Split(licenseBlock, "\n") {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, lm.commentStyle.Single) {
				// Remove the comment character and one optional space after it
				cleanedLine := strings.TrimPrefix(trimmedLine, lm.commentStyle.Single)
				if strings.HasPrefix(cleanedLine, " ") {
					cleanedLine = cleanedLine[1:]
				}
				cleanedLines = append(cleanedLines, cleanedLine)
			} else if trimmedLine == "" {
				cleanedLines = append(cleanedLines, "")
			}
		}
		licenseBlock = strings.Join(cleanedLines, "\n")
	}

	// Trim any extra whitespace
	return strings.TrimSpace(licenseBlock), true
}

// CheckLicenseStatus verifies the license status of the content
func (lm *LicenseManager) CheckLicenseStatus(content string) LicenseStatus {
	extractedText, found := lm.extractLicenseText(content)
	if !found {
		return NoLicense
	}

	// Clean both texts for comparison (remove empty lines and whitespace)
	cleanExtracted := strings.TrimSpace(extractedText)
	cleanLicense := strings.TrimSpace(lm.licenseText)

	if cleanExtracted == cleanLicense {
		return MatchingLicense
	}
	return DifferentLicense
}

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string) bool {
	return lm.CheckLicenseStatus(content) == MatchingLicense
}

// UpdateLicense updates the existing license text with new content
func (lm *LicenseManager) UpdateLicense(content string) string {
	status := lm.CheckLicenseStatus(content)
	if status == NoLicense {
		return content // No license found to update
	}

	header := lm.style.Header
	footer := lm.style.Footer

	// Find the start of the license block
	startIdx := strings.Index(content, header)
	if startIdx == -1 {
		return content
	}

	// Look for footer after the header
	afterHeader := content[startIdx+len(header):]
	endIdx := strings.Index(afterHeader, footer)
	if endIdx == -1 {
		return content
	}

	// Get the content before the header to check for comment characters
	beforeHeader := strings.TrimSpace(content[:startIdx])
	commentPrefix := ""
	if lm.commentStyle.Single != "" {
		// If there's a comment character before the header, use it
		if strings.HasSuffix(beforeHeader, lm.commentStyle.Single) {
			commentPrefix = lm.commentStyle.Single + " "
		} else {
			commentPrefix = lm.commentStyle.Single + " "
		}
	}

	// Format the new license block
	formattedLicense := commentPrefix + lm.style.Header + "\n\n" + 
		lm.formatLicenseBlock(lm.licenseText) + "\n" + 
		commentPrefix + lm.style.Footer
	
	// Construct the updated content by preserving everything before and after the license block
	beforeLicense := content[:startIdx]
	if idx := strings.LastIndex(beforeLicense, lm.commentStyle.Single); idx != -1 {
		// Remove any existing comment characters before the header
		beforeLicense = strings.TrimSpace(beforeLicense[:idx]) + "\n"
	}
	afterLicense := afterHeader[endIdx+len(footer):]
	
	return beforeLicense + formattedLicense + afterLicense
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	header := lm.style.Header
	footer := lm.style.Footer

	startIdx := strings.Index(content, header)
	if startIdx == -1 {
		return content
	}

	// Look for footer after the header
	remainingContent := content[startIdx:]
	endIdx := strings.Index(remainingContent, footer)
	if endIdx == -1 {
		return content
	}

	// Get the line number information
	preContent := content[:startIdx]
	startLine := len(strings.Split(preContent, "\n"))
	licenseBlock := content[startIdx:startIdx+endIdx+len(footer)]
	
	// Remove everything up to the end of the footer
	result := strings.TrimLeft(content[startIdx+endIdx+len(footer):], "\n\r\t ")

	// Log the removed content with line numbers if in verbose mode
	if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
		fp.logVerboseWithLineNumbers(licenseBlock, startLine, "Removing license block:")
	}

	return result
}

// GetLicenseComparison returns the current and expected license text for comparison
func (lm *LicenseManager) GetLicenseComparison(content string) (current, expected string) {
	// Extract current license text
	currentText, found := lm.extractLicenseText(content)
	if !found {
		return "", ""
	}

	// For expected text, we need to format it the same way as the current text
	// First format with comments, then extract the license text to ensure consistent comparison
	formattedLicense := ""
	if lm.commentStyle.Single != "" {
		formattedLicense = lm.commentStyle.Single + " " + lm.style.Header + "\n\n" + 
			lm.formatLicenseBlock(lm.licenseText) + "\n" + 
			lm.commentStyle.Single + " " + lm.style.Footer
	} else {
		formattedLicense = lm.style.Header + "\n\n" + 
			lm.formatLicenseBlock(lm.licenseText) + "\n" + 
			lm.style.Footer
	}

	// Extract the expected text using the same extraction process
	expectedText, _ := lm.extractLicenseText(formattedLicense)

	return currentText, expectedText
}

// LicenseStatus represents the status of a license check
type LicenseStatus int

const (
	NoLicense LicenseStatus = iota
	DifferentLicense
	MatchingLicense
)
