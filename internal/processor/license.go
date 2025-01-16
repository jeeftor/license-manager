// internal/processor/license.go
package processor

import (
	"bytes"
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
func (lm *LicenseManager) formatLicenseBlock(content string) string {
	var lines []string

	// Special handling for Go files
	if lm.commentStyle.FileType == "go" {
		return lm.formatGoLicenseBlock(content)
	}

	// For multi-line comments, we need to:
	// 1. Add each line with a * prefix
	// 2. The multi-line comment markers will be added in AddLicense
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// For all multi-line comment styles
		for _, line := range strings.Split(content, "\n") {
			if line == "" {
				lines = append(lines, " *")
			} else {
				lines = append(lines, " * "+line)
			}
		}
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
	// Check if the file already has a license
	if lm.CheckLicense(content) {
		if fp, ok := interface{}(lm).(*FileProcessor); ok && fp.config.Verbose {
			fp.logVerbose("Found existing license - skipping license addition")
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
	
	// Add header, footer, and comment markers
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		// Handle JSX-style comments
		if strings.HasSuffix(lm.commentStyle.MultiEnd, "}") {
			formattedLicense = lm.commentStyle.MultiStart + "\n" +
				commentedLicenseText + "\n" +
				lm.commentStyle.MultiEnd + "\n\n"
		} else {
			formattedLicense = lm.commentStyle.MultiStart + "\n" +
				commentedLicenseText + "\n" +
				lm.commentStyle.MultiEnd + "\n\n"
		}
	} else if lm.commentStyle.Single != "" {
		formattedLicense = commentedLicenseText + "\n\n"
	} else {
		formattedLicense = commentedLicenseText + "\n\n"
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
			buf.WriteString(strings.Join(lines[:buildTagsEnd], "\n") + "\n\n")
			content = strings.Join(lines[buildTagsEnd:], "\n")
		}
	}

	buf.WriteString(formattedLicense)
	buf.WriteString(content)

	return buf.String()
}

// RemoveLicense removes the license text from the content
func (lm *LicenseManager) RemoveLicense(content string) string {
	// For JSX files, we need to find the license within JSX comments
	if lm.commentStyle.FileType == "javascript" && strings.HasSuffix(lm.commentStyle.MultiEnd, "}") {
		// Find all JSX comment blocks
		var result []string
		parts := strings.Split(content, lm.commentStyle.MultiStart)
		for i, part := range parts {
			if i == 0 {
				// Keep everything before the first comment block
				result = append(result, strings.TrimSpace(part))
				continue
			}

			// Find the end of this comment block
			commentParts := strings.SplitN(part, lm.commentStyle.MultiEnd, 2)
			if len(commentParts) != 2 {
				// If we can't find the end, keep the part as is
				result = append(result, lm.commentStyle.MultiStart+part)
				continue
			}

			// Check if this comment block contains our formatted license text
			commentContent := commentParts[0]
			formattedLicense := lm.formatLicenseBlock(lm.licenseText)
			if !strings.Contains(commentContent, formattedLicense) {
				// If it's not a license block, keep it
				result = append(result, lm.commentStyle.MultiStart+part)
			} else {
				// If it is a license block, only keep what comes after it
				afterLicense := strings.TrimSpace(commentParts[1])
				if afterLicense != "" {
					result = append(result, afterLicense)
				}
			}
		}
		return strings.Join(result, "\n")
	}

	// For other files with multi-line comments
	if lm.commentStyle.PreferMulti && lm.commentStyle.MultiStart != "" {
		start := strings.Index(content, lm.commentStyle.MultiStart)
		if start != -1 {
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
				if strings.Contains(cleanedContent.String(), lm.licenseText) {
					// Keep content before the license and after the footer
					beforeLicense := strings.TrimSpace(content[:start])
					afterLicense := strings.TrimSpace(content[start+end+len(lm.commentStyle.MultiEnd):])
					
					if beforeLicense == "" {
						return afterLicense
					}
					if afterLicense == "" {
						return beforeLicense
					}
					return beforeLicense + "\n\n" + afterLicense
				}
			}
		}
		return content
	}

	// For files with single-line comments
	if lm.commentStyle.Single != "" {
		formattedLicense := lm.formatLicenseBlock(lm.licenseText)
		lines := strings.Split(content, "\n")
		var result []string
		var inLicense bool
		var licenseStarted bool

		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
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

			if strings.TrimSpace(line) != "" {
				result = append(result, line)
			}
		}

		return strings.Join(result, "\n")
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

// CheckLicense verifies if the content contains a matching license
func (lm *LicenseManager) CheckLicense(content string) bool {
	normalizedLicense := normalizeText(lm.licenseText)

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
				normalizedComment := normalizeText(cleanedContent.String())
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
				normalizedComment := normalizeText(commentBlock.String())
				if strings.Contains(normalizedComment, normalizedLicense) {
					return true
				}
				commentBlock.Reset()
				inComment = false
			} else if inComment {
				// Non-empty line after comment block
				normalizedComment := normalizeText(commentBlock.String())
				if strings.Contains(normalizedComment, normalizedLicense) {
					return true
				}
				commentBlock.Reset()
				inComment = false
			}
		}
		// Check last comment block
		if inComment {
			normalizedComment := normalizeText(commentBlock.String())
			if strings.Contains(normalizedComment, normalizedLicense) {
				return true
			}
		}
	}

	// For files without comments
	normalizedContent := normalizeText(content)
	return strings.Contains(normalizedContent, normalizedLicense)
}

// normalizeText removes all whitespace and punctuation from text to make it easier to compare
func normalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	// Remove all whitespace
	text = strings.Join(strings.Fields(text), "")
	// Remove all punctuation
	text = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)
	return text
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
