// internal/processor/processor.go
package processor

import (
	"fmt"
	"license-manager/internal/license"
	"os"
	"path/filepath"
	"strings"

	"license-manager/internal/styles"

	"github.com/fatih/color"
)

type FileProcessor struct {
	config  *Config
	license string
	style   styles.HeaderFooterStyle
	stats   struct {
		Added     int
		Existing  int
		Skipped   int
		Failed    int
		Unchanged int
		Errors    int
	}
	verbose bool
}

// Colored output helpers
var (
	errorColor   = color.New(color.FgRed).SprintFunc()
	warningColor = color.New(color.FgYellow).SprintFunc()
	successColor = color.New(color.FgGreen).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
)

func NewFileProcessor(config *Config, license string, style styles.HeaderFooterStyle) *FileProcessor {
	return &FileProcessor{
		config:  config,
		license: license,
		style:   style,
	}
}

func (fp *FileProcessor) logVerbose(format string, args ...interface{}) {
	if fp.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func (fp *FileProcessor) logVerboseWithLineNumbers(text string, startLine int, prefix string) {
	if fp.verbose {
		if prefix != "" {
			fmt.Printf("%s\n", prefix)
		}
		fmt.Printf("%s\n", formatWithLineNumbers(text, startLine))
	}
}

func (fp *FileProcessor) resetStats() {
	fp.stats.Added = 0
	fp.stats.Existing = 0
	fp.stats.Skipped = 0
	fp.stats.Failed = 0
	fp.stats.Unchanged = 0
	fp.stats.Errors = 0
}

func (fp *FileProcessor) Add() error {
	fp.resetStats()
	err := fp.processFiles(func(filename, content string, license *license.LicenseManager) error {
		if license.CheckLicense(content, true) {
			fp.stats.Unchanged++
			return NewCheckError(fmt.Sprintf("license already exists in file: %s (use 'update' command to modify existing licenses)", filename))
		}
		newContent := license.AddLicense(content)
		fp.logVerbose("%s %s", successColor("✅ Adding license to:"), filename)
		fp.stats.Added++
		return os.WriteFile(filename, []byte(newContent), 0644)
	})

	// If the only errors were CheckErrors (existing licenses), return nil
	if _, ok := err.(*CheckError); ok {
		return nil
	}
	return err
}

func (fp *FileProcessor) Remove() error {
	fp.resetStats()
	return fp.processFiles(func(filename, content string, license *license.LicenseManager) error {
		newContent := license.RemoveLicense(content)
		if newContent != content {
			fp.logVerbose("%s %s", successColor("✅ Removing license from:"), filename)
			fp.stats.Added++
		} else {
			fp.logVerbose("%s %s", infoColor("ℹ️ No license found in:"), filename)
			fp.stats.Skipped++
		}
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Update() error {
	fp.resetStats()
	err := fp.processFiles(func(filename, content string, license *license.LicenseManager) error {
		status := license.CheckLicenseStatus(content)
		if status == license.NoLicense {
			fp.stats.Skipped++
			return NewCheckError(fmt.Sprintf("no license found to update in file: %s", filename))
		}

		licenseText, err := fp.readLicenseText()
		if err != nil {
			return err
		}
		newContent := license.UpdateLicense(string(content), licenseText)
		if newContent != content {
			if status == license.DifferentLicense {
				fp.logVerbose("%s %s", warningColor("⚠️ Updating different license in:"), filename)
				fp.stats.Added++
			} else {
				fp.logVerbose("%s %s", successColor("✅ Updating matching license in:"), filename)
				fp.stats.Unchanged++
			}
			return os.WriteFile(filename, []byte(newContent), 0644)
		}
		return nil
	})

	return err
}

func (fp *FileProcessor) Check() error {
	fp.resetStats()
	var hasFailures bool
	err := fp.processFiles(func(filename, content string, license *license.LicenseManager) error {
		status := license.CheckLicenseStatus(content)
		switch status {
		case license.NoLicense:
			hasFailures = true
			fp.stats.Skipped++
			fmt.Printf("%s %s\n", errorColor("❌ No license found in file:"), filename)
		case license.DifferentLicense:
			hasFailures = true
			fp.stats.Failed++
			fmt.Printf("%s %s\n", warningColor("⚠️ License doesn't match in file:"), filename)
			if fp.verbose {
				current, expected := license.GetLicenseComparison(content)
				fmt.Printf("\nCurrent license in %s:\n%s\n", filename, infoColor(current))
				fmt.Printf("\nExpected license:\n%s\n", successColor(expected))
				fmt.Println("\nDifferences:")
				// Print a simple character-based diff
				currentLines := strings.Split(current, "\n")
				expectedLines := strings.Split(expected, "\n")
				for i := 0; i < len(currentLines) || i < len(expectedLines); i++ {
					var currentLine, expectedLine string
					if i < len(currentLines) {
						currentLine = currentLines[i]
					}
					if i < len(expectedLines) {
						expectedLine = expectedLines[i]
					}
					if currentLine != expectedLine {
						if currentLine == "" {
							fmt.Printf("%s %s\n", errorColor("-"), expectedLine)
						} else if expectedLine == "" {
							fmt.Printf("%s %s\n", warningColor("+"), currentLine)
						} else {
							fmt.Printf("%s %s\n%s %s\n", errorColor("-"), expectedLine, warningColor("+"), currentLine)
						}
					}
				}
				fmt.Println()
			}
		case license.MatchingLicense:
			fp.stats.Unchanged++
			fp.logVerbose("%s %s", successColor("✅ License matches in:"), filename)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if hasFailures {
		return fmt.Errorf("one or more files have missing or incorrect licenses")
	}

	return nil
}

// readLicenseText reads the license text from the license file
func (fp *FileProcessor) readLicenseText() (string, error) {
	// content, err := os.ReadFile(fp.licenseFile)
	// if err != nil {
	// 	return "", fmt.Errorf("error reading license file: %v", err)
	// }
	return fp.license, nil
}

// formatWithLineNumbers adds line numbers to a block of text
func formatWithLineNumbers(text string, startLine int) string {
	lines := strings.Split(text, "\n")
	var result []string

	// Find the width needed for line numbers
	width := len(fmt.Sprintf("%d", startLine+len(lines)))

	// Format each line with its number
	for i, line := range lines {
		lineNum := fmt.Sprintf("%*d", width, startLine+i)
		result = append(result, fmt.Sprintf("%s: %s", lineNum, line))
	}

	return strings.Join(result, "\n")
}

func getCommentStyle(filename string) styles.CommentLanguage {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go":
		return styles.CommentLanguage{
			Language:    "go",
			Single:      "//",
			MultiStart:  "/*",
			MultiEnd:    "*/",
			LinePrefix:  " ",
			PreferMulti: true,
		}
	case ".py":
		return styles.CommentLanguage{
			Language:   "python",
			Single:     "#",
			LinePrefix: " ",
		}
	case ".js", ".ts", ".jsx", ".tsx":
		return styles.CommentLanguage{
			Language:   "javascript",
			Single:     "//",
			MultiStart: "/*",
			MultiEnd:   "*/",
			LinePrefix: " ",
		}
	case ".html", ".htm":
		return styles.CommentLanguage{
			Language:   "html",
			MultiStart: "<!--",
			MultiEnd:   "-->",
			LinePrefix: " ",
		}
	default:
		return styles.CommentLanguage{
			Language: "text",
		}
	}
}
