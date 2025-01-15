// internal/processor/processor.go
package processor

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Config struct {
	Header      string
	Footer      string
	LicenseText string
	Input       string
	Skip        string
	Prompt      bool
	DryRun      bool
	Verbose     bool
}

type FileProcessor struct {
	config Config
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

	return &FileProcessor{
		config: config,
	}
}

func (fp *FileProcessor) logVerbose(format string, args ...interface{}) {
	if fp.config.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func (fp *FileProcessor) Add() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		if license.CheckLicense(content) {
			fp.logVerbose("%s %s", successColor("License already present in:"), filename)
			return nil
		}
		newContent := license.AddLicense(content)
		fp.logVerbose("%s %s", successColor("Adding license to:"), filename)
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Remove() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		newContent := license.RemoveLicense(content)
		if newContent != content {
			fp.logVerbose("%s %s", successColor("Removing license from:"), filename)
		} else {
			fp.logVerbose("%s %s", infoColor("No license found in:"), filename)
		}
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Update() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		newContent := license.UpdateLicense(content)
		if newContent != content {
			fp.logVerbose("%s %s", successColor("Updating license in:"), filename)
		} else {
			fp.logVerbose("%s %s", infoColor("No license found to update in:"), filename)
		}
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Check() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		if !license.CheckLicense(content) {
			fmt.Printf("%s %s\n", errorColor("License missing or invalid in file:"), filename)
		} else if fp.config.Verbose {
			fmt.Printf("%s %s\n", successColor("License valid in file:"), filename)
		}
		return nil
	})
}
