// internal/processor/processor.go
package processor

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type FileProcessor struct {
	config Config
	style  HeaderFooterStyle
}

// Colored output helpers
var (
	errorColor   = color.New(color.FgRed).SprintFunc()
	warningColor = color.New(color.FgYellow).SprintFunc()
	successColor = color.New(color.FgGreen).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
)

func NewFileProcessor(config Config) *FileProcessor {
	// Read license text from file if provided
	if config.LicenseText != "" {
		content, err := os.ReadFile(config.LicenseText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", errorColor("Error reading license file"), err)
			os.Exit(1)
		}
		config.LicenseText = string(content)
	}

	// Get the preset style - the markers are already included in preset styles
	style := GetPresetStyle(config.PresetStyle)

	return &FileProcessor{
		config: config,
		style:  style,
	}
}

func (fp *FileProcessor) logVerbose(format string, args ...interface{}) {
	if fp.config.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func (fp *FileProcessor) logVerboseWithLineNumbers(text string, startLine int, prefix string) {
	if fp.config.Verbose {
		if prefix != "" {
			fmt.Printf("%s\n", prefix)
		}
		fmt.Printf("%s\n", formatWithLineNumbers(text, startLine))
	}
}

func (fp *FileProcessor) Add() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		if license.CheckLicense(content) {
			return NewCheckError(fmt.Sprintf("license already exists in file: %s (use 'update' command to modify existing licenses)", filename))
		}
		newContent := license.AddLicense(content)
		fp.logVerbose("%s %s", successColor("✅ Adding license to:"), filename)
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Remove() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		newContent := license.RemoveLicense(content)
		if newContent != content {
			fp.logVerbose("%s %s", successColor("✅ Removing license from:"), filename)
		} else {
			fp.logVerbose("%s %s", infoColor("ℹ️ No license found in:"), filename)
		}
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Update() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		status := license.CheckLicenseStatus(content)
		if status == NoLicense {
			return NewCheckError(fmt.Sprintf("no license found to update in file: %s", filename))
		}
		
		newContent := license.UpdateLicense(content)
		if status == DifferentLicense {
			fp.logVerbose("%s %s", warningColor("⚠️ Updating different license in:"), filename)
		} else {
			fp.logVerbose("%s %s", successColor("✅ Updating matching license in:"), filename)
		}
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Check() error {
	hasFailures := false
	err := fp.processFiles(func(filename, content string, license *LicenseManager) error {
		status := license.CheckLicenseStatus(content)
		switch status {
		case NoLicense:
			hasFailures = true
			fmt.Printf("%s %s\n", errorColor("❌ No license found in file:"), filename)
		case DifferentLicense:
			hasFailures = true
			fmt.Printf("%s %s\n", warningColor("⚠️ License doesn't match in file:"), filename)
			if fp.config.Verbose {
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
		case MatchingLicense:
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
