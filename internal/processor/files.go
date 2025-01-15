// internal/processor/files.go
package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

func (fp *FileProcessor) processFiles(action func(string, string, *LicenseManager) error) error {
	patterns := strings.Split(fp.config.Input, ",")

	for _, basePattern := range patterns {
		if basePattern == "" {
			continue
		}

		matches, err := doublestar.Glob(basePattern)
		if err != nil {
			return fmt.Errorf("%s %s: %v", errorColor("Error with pattern"), basePattern, err)
		}
		for _, match := range matches {
			fp.logVerbose("%s %s", successColor("Processing file:"), match)
			if err := fp.processFile(match, action); err != nil {
				return err
			}
		}
	}
	return nil
}

func (fp *FileProcessor) processFile(filename string, action func(string, string, *LicenseManager) error) error {
	// Check if path is a directory first
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return NewCheckError(fmt.Sprintf("error accessing file %s: %v", filename, err))
	}
	if fileInfo.IsDir() {
		fp.logVerbose("%s %s", warningColor("Skipping directory:"), filename)
		return nil
	}

	// Skip files that match skip patterns
	for _, pattern := range strings.Split(fp.config.Skip, ",") {
		if pattern != "" {
			matched, err := filepath.Match(pattern, filepath.Base(filename))
			if err != nil {
				return fmt.Errorf("%s %s: %v", errorColor("Invalid skip pattern"), pattern, err)
			}
			if matched {
				fp.logVerbose("%s %s", warningColor("Skipping file:"), filename)
				return nil
			}
		}
	}

	if fp.config.DryRun {
		fmt.Printf("%s %s\n", infoColor("Would process file:"), filename)
		return nil
	}

	if fp.config.Prompt {
		if !promptUser(fmt.Sprintf("Process file %s?", filename)) {
			fp.logVerbose("%s %s", warningColor("Skipping file (user choice):"), filename)
			return nil
		}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return NewCheckError(fmt.Sprintf("%s: %v", errorColor("Error reading file "+filename), err))
	}

	commentStyle := getCommentStyle(filename)
	fp.logVerbose("%s %s: %s", infoColor("Using comment style for"), filename, commentStyle.FileType)

	// Create LicenseManager with the HeaderFooterStyle
	license := NewLicenseManager(fp.style, fp.config.LicenseText, commentStyle)

	return action(filename, string(content), license)
}

func promptUser(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", message)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
