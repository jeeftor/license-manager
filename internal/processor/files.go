package processor

import (
	"bufio"
	"fmt"
	"license-manager/internal/license"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

func (fp *FileProcessor) processFiles(action func(string, string, *license.LicenseManager) error) error {
	patterns := strings.Split(fp.config.Input, ",")
	var lastErr error

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
				if _, ok := err.(*CheckError); ok {
					// For CheckErrors (like existing license), continue processing
					fp.logVerbose("%s %s", warningColor("⚠️"), err)
					lastErr = err
					continue
				}
				// For other errors (like file system errors), stop processing
				return err
			}
		}
	}

	// Print summary
	if fp.stats.Added > 0 || fp.stats.Existing > 0 {
		fmt.Println("\nSummary:")
		if fp.stats.Added > 0 {
			fmt.Printf("%s %d files\n", successColor("✅ Added license to:"), fp.stats.Added)
		}
		if fp.stats.Existing > 0 {
			fmt.Printf("%s %d files (use 'update' command to modify)\n", warningColor("⚠️ License already exists in:"), fp.stats.Existing)
		}
		if fp.stats.Skipped > 0 {
			fmt.Printf("%s %d files\n", infoColor("ℹ️ Skipped:"), fp.stats.Skipped)
		}
		if fp.stats.Errors > 0 {
			fmt.Printf("%s %d files\n", errorColor("❌ Errors in:"), fp.stats.Errors)
		}
	}

	return lastErr
}

func (fp *FileProcessor) processFile(filename string, action func(string, string, *license.LicenseManager) error) error {
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

	// Get comment style for file extension
	commentStyle := getCommentStyle(filename)
	commentStyle.PreferMulti = fp.config.PreferMulti

	// Create license manager for this file
	license := license.NewLicenseManager(fp.style, fp.config.LicenseText, commentStyle)

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
