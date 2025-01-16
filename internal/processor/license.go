// internal/processor/license.go
package processor

import (
	"strings"
	"unicode"
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
func (lm *LicenseManager) formatLicenseBlock(license string) string {
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// For multi-line comments, just add indentation
		lines := strings.Split(license, "\n")
		var builder strings.Builder

		// Add license text with proper indentation
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				builder.WriteString(" " + line + "\n")
			} else {
				builder.WriteString(" \n")
			}
		}
		return builder.String()
	}

	// For single-line comments
	lines := strings.Split(license, "\n")
	var builder strings.Builder
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			builder.WriteString(lm.commentStyle.Single + " " + line)
		} else {
			builder.WriteString(lm.commentStyle.Single)
		}
		if i < len(lines)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
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

// AddLicense adds a license header to the content if one doesn't already exist
func (lm *LicenseManager) AddLicense(content string) string {
	// Check if the content already has the license
	if lm.CheckLicense(content) {
		if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
			fp.logVerbose("Found existing license - skipping license addition")
		}
		return content
	}

	// Format the license text with appropriate comment style
	formattedLicense := lm.formatLicenseBlock(lm.licenseText)

	// For multi-line comments, wrap the formatted license in comment markers
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		formattedLicense = lm.commentStyle.MultiStart + "\n" + formattedLicense + lm.commentStyle.MultiEnd + "\n"
	}

	// Add an extra newline after the license block
	formattedLicense = formattedLicense + "\n"

	// If the content is empty, just return the license
	if strings.TrimSpace(content) == "" {
		return formattedLicense
	}

	// Otherwise, add the license before the content
	return formattedLicense + content
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// For multi-line comments
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		start := strings.Index(content, lm.commentStyle.MultiStart)
		if start != -1 {
			end := strings.Index(content[start:], lm.commentStyle.MultiEnd)
			if end != -1 {
				// Extract the comment block without the comment markers
				commentBlock := content[start+len(lm.commentStyle.MultiStart) : start+end]
				commentBlock = strings.TrimSpace(commentBlock)

				// Format the license text without comment markers for comparison
				formattedLicense := lm.formatLicenseBlock(lm.licenseText)
				formattedLicense = strings.TrimSpace(formattedLicense)

				// Compare normalized text to check if this is our license
				normalizedComment := lm.normalizeText(commentBlock)
				normalizedLicense := lm.normalizeText(formattedLicense)

				if strings.Contains(normalizedComment, normalizedLicense) {
					// Keep content before and after the license block
					beforeLicense := content[:start]
					afterLicense := content[start+end+len(lm.commentStyle.MultiEnd):]

					// Remove any leading newlines from afterLicense
					afterLicense = strings.TrimLeft(afterLicense, "\n")
					
					if strings.TrimSpace(beforeLicense) == "" {
						return afterLicense
					}
					if strings.TrimSpace(afterLicense) == "" {
						return beforeLicense
					}
					return beforeLicense + afterLicense
				}
			}
		}
		return content
	}

	// For single-line comments
	if lm.commentStyle.Single != "" {
		formattedLicense := lm.formatLicenseBlock(lm.licenseText)
		lines := strings.Split(content, "\n")
		var result []string
		var inLicense bool
		var licenseStarted bool
		var shebangLine string

		// Extract shebang line if present
		if len(lines) > 0 && strings.HasPrefix(lines[0], "#!") {
			shebangLine = lines[0]
			lines = lines[1:]
		}

		// Process lines
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)

			// Skip license block
			if !licenseStarted && strings.HasPrefix(trimmedLine, lm.commentStyle.Single) {
				// Check if this is the start of our license block
				remainingLines := make([]string, 0)
				for i := len(result); i < len(lines); i++ {
					if strings.HasPrefix(strings.TrimSpace(lines[i]), lm.commentStyle.Single) {
						remainingLines = append(remainingLines, lines[i])
					} else {
						break
					}
				}
				potentialBlock := strings.Join(remainingLines, "\n")
				if strings.Contains(potentialBlock, formattedLicense) {
					inLicense = true
					licenseStarted = true
					continue
				}
			}

			if inLicense {
				if !strings.HasPrefix(trimmedLine, lm.commentStyle.Single) {
					inLicense = false
				} else {
					continue
				}
			}

			result = append(result, line)
		}

		// Add shebang line back to the beginning if present
		if shebangLine != "" {
			result = append([]string{shebangLine}, result...)
		}

		// Join all lines and preserve newlines
		return preserveNewlines(strings.Join(result, "\n"), content)
	}

	// For files without comments, just look for the raw license text
	startIdx := strings.Index(content, lm.licenseText)
	if startIdx == -1 {
		return content
	}

	endIdx := startIdx + len(lm.licenseText)
	beforeLicense := strings.TrimSpace(content[:startIdx])
	afterLicense := strings.TrimSpace(content[endIdx:])
	
	if beforeLicense == "" {
		return afterLicense
	}
	if afterLicense == "" {
		return beforeLicense
	}
	return beforeLicense + "\n\n" + afterLicense
}

// preserveNewlines ensures that the content has the same number of leading and trailing newlines as the original
func preserveNewlines(content, original string) string {
	// Count leading newlines in original
	originalLeadingCount := 0
	for i := 0; i < len(original); i++ {
		if original[i] != '\n' {
			break
		}
		originalLeadingCount++
	}

	// Count trailing newlines in original
	originalTrailingCount := 0
	for i := len(original) - 1; i >= 0; i-- {
		if original[i] != '\n' {
			break
		}
		originalTrailingCount++
	}

	// Trim all newlines from content
	trimmedContent := strings.Trim(content, "\n")

	// Add the same number of newlines as in original
	result := strings.Repeat("\n", originalLeadingCount) + trimmedContent + strings.Repeat("\n", originalTrailingCount)
	return result
}

// normalizeText removes whitespace and comment markers from text
func (lm *LicenseManager) normalizeText(text string) string {
	// Remove all whitespace
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)

	// Remove all punctuation
	text = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)

	return strings.ToLower(text)
}

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string) bool {
	normalizedLicense := lm.normalizeText(lm.licenseText)

	// For files with multi-line comments
	if lm.commentStyle.MultiStart != "" {
		start := strings.Index(content, lm.commentStyle.MultiStart)
		for start != -1 {
			end := strings.Index(content[start:], lm.commentStyle.MultiEnd)
			if end != -1 {
				commentContent := content[start : start+end+len(lm.commentStyle.MultiEnd)]
				// Clean up the comment content
				commentContent = strings.TrimSpace(commentContent)
				commentContent = strings.TrimPrefix(commentContent, lm.commentStyle.MultiStart)
				commentContent = strings.TrimSuffix(commentContent, lm.commentStyle.MultiEnd)
				commentContent = strings.TrimSpace(commentContent)
				// For each line, remove leading asterisk and spaces
				lines := strings.Split(commentContent, "\n")
				var cleanedContent strings.Builder
				for _, line := range lines {
					line = strings.TrimSpace(line)
					line = strings.TrimPrefix(line, "*")
					line = strings.TrimSpace(line)
					cleanedContent.WriteString(line + "\n")
				}
				normalizedComment := lm.normalizeText(cleanedContent.String())
				if strings.Contains(normalizedComment, normalizedLicense) {
					return true
				}
				// Move to next comment block
				start = strings.Index(content[start+end+len(lm.commentStyle.MultiEnd):], lm.commentStyle.MultiStart)
				if start != -1 {
					start += end + len(lm.commentStyle.MultiEnd)
				}
			} else {
				break
			}
		}
	}

	// For files with single-line comments
	if lm.commentStyle.Single != "" {
		lines := strings.Split(content, "\n")
		var commentBlock strings.Builder
		inComment := false
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, lm.commentStyle.Single) {
				commentLine := strings.TrimPrefix(trimmedLine, lm.commentStyle.Single)
				commentLine = strings.TrimSpace(commentLine)
				commentBlock.WriteString(commentLine + "\n")
				inComment = true
			} else if inComment && trimmedLine == "" {
				// Empty line after comment block
				normalizedComment := lm.normalizeText(commentBlock.String())
				if strings.Contains(normalizedComment, normalizedLicense) {
					return true
				}
				commentBlock.Reset()
				inComment = false
			} else if inComment {
				// Non-empty line after comment block
				normalizedComment := lm.normalizeText(commentBlock.String())
				if strings.Contains(normalizedComment, normalizedLicense) {
					return true
				}
				commentBlock.Reset()
				inComment = false
			}
		}
		// Check last comment block
		if inComment {
			normalizedComment := lm.normalizeText(commentBlock.String())
			if strings.Contains(normalizedComment, normalizedLicense) {
				return true
			}
		}
	}

	// For files without comments
	normalizedContent := lm.normalizeText(content)
	return strings.Contains(normalizedContent, normalizedLicense)
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
	} else if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// For multi-line comments, remove the start/end markers and any leading asterisks
		licenseBlock = strings.TrimSpace(licenseBlock)
		
		// Handle JSX-style comments and regular multi-line comments
		if strings.HasPrefix(licenseBlock, lm.commentStyle.MultiStart) {
			licenseBlock = licenseBlock[len(lm.commentStyle.MultiStart):]
		}
		if strings.HasSuffix(licenseBlock, lm.commentStyle.MultiEnd) {
			licenseBlock = licenseBlock[:len(licenseBlock)-len(lm.commentStyle.MultiEnd)]
		}
		
		// Clean up any leading asterisks that are commonly used in multi-line comments
		var cleanedLines []string
		for _, line := range strings.Split(licenseBlock, "\n") {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "*") {
				cleanedLine := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "*"))
				cleanedLines = append(cleanedLines, cleanedLine)
			} else {
				cleanedLines = append(cleanedLines, trimmedLine)
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

	// For JSX files, we need to handle the special comment style
	if lm.commentStyle.FileType == "javascript" && strings.HasSuffix(lm.commentStyle.MultiEnd, "}") {
		// Remove any remaining JSX comment markers that might have been missed
		cleanExtracted = strings.TrimPrefix(cleanExtracted, "{")
		cleanExtracted = strings.TrimSuffix(cleanExtracted, "}")
		cleanExtracted = strings.TrimSpace(cleanExtracted)
	}

	if cleanExtracted == cleanLicense {
		return MatchingLicense
	}
	return DifferentLicense
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
