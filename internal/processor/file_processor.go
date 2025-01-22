package processor

import (
	"fmt"
	"license-manager/internal/language"
	"license-manager/internal/license"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"path/filepath"
	"strings"
)

// FileProcessor handles license operations on files
type FileProcessor struct {
	config      *Config
	fileHandler *FileHandler
	logger      *logger.Logger
	stats       map[string]int
}

// NewFileProcessor creates a new FileProcessor instance
func NewFileProcessor(cfg *Config) *FileProcessor {
	log := logger.NewLogger(cfg.Verbose)
	fh := NewFileHandler(log)
	fh.SetSkipPattern(cfg.Skip) // Set the skip pattern
	return &FileProcessor{
		config:      cfg,
		logger:      log,
		fileHandler: fh,
		stats:       make(map[string]int),
	}
}

// createManager creates and configures a license manager for a given file
func (fp *FileProcessor) createManager(file string) (*license.LicenseManager, styles.CommentLanguage) {
	style := styles.Get(fp.config.PresetStyle)

	// Get comment style for file type
	ext := filepath.Ext(file)
	commentStyle := styles.GetLanguageCommentStyle(ext)
	commentStyle.PreferMulti = fp.config.PreferMulti

	if fp.config.Verbose {
		fp.logger.LogInfo("Processing file: %s", file)
		fp.logger.LogInfo("  Language: %s", commentStyle.Language)
		fp.logger.LogInfo("  Comment style: %s", describeCommentStyle(commentStyle))
	}

	// Create and configure manager
	manager := license.NewLicenseManager(fp.config.LicenseText, style)
	manager.SetCommentStyle(commentStyle)
	if fp.config.Verbose {
		manager.SetVerbose(true, fp.logger)
	}

	return manager, commentStyle
}

// resetStats resets the operation statistics
func (fp *FileProcessor) resetStats() {
	fp.stats = map[string]int{
		"added":     0,
		"existing":  0,
		"skipped":   0,
		"failed":    0,
		"unchanged": 0,
	}
}

// Add adds license headers to files
func (fp *FileProcessor) Add() error {
	fp.resetStats()

	files, err := fp.fileHandler.FindFiles(fp.config.Input)
	if err != nil {
		return err
	}

	// Print scanning message with patterns
	fp.logger.LogInfo("Scanning %d Directories:", len(files))
	fp.logger.LogInfo("Inputs Patterns:")
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			fp.logger.LogInfo("  %s", pattern)
		}
	}
	if fp.config.Skip != "" {
		fp.logger.LogInfo("Skips Patterns:")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				fp.logger.LogInfo("  %s", pattern)
			}
		}
	}

	style := styles.Get(fp.config.PresetStyle)
	if fp.config.Verbose {
		fp.logger.LogInfo("Using style: %s", style.Name)
		if fp.config.PreferMulti {
			fp.logger.LogInfo("Preferring multi-line comments where supported")
		}
	}

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, commentStyle := fp.createManager(file)

		// Get the appropriate language handler and check for preamble
		handler := language.GetLanguageHandler(commentStyle.Language, style)
		preamble, rest := handler.PreservePreamble(content)

		if fp.config.Verbose {
			if preamble != "" {
				fp.logger.LogInfo("  Found preamble:")
				for _, line := range strings.Split(strings.TrimSpace(preamble), "\n") {
					if line != "" {
						fp.logger.LogInfo("    %s", line)
					}
				}
			} else {
				fp.logger.LogInfo("  No preamble found")
			}
		}

		if manager.HasLicense(content) {
			fp.stats["existing"]++
			fp.logger.LogWarning("License already exists in %s", file)
			continue
		}

		// Use rest instead of content to add license after preamble
		newContent, err := manager.AddLicense(rest, commentStyle.Language)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to add license to %s: %v", file, err)
			continue
		}

		// Recombine preamble with the licensed content
		if preamble != "" {
			newContent = preamble + "\n" + newContent
		}

		// Debug the actual comment being added in verbose mode
		if fp.config.Verbose {
			fp.logger.LogInfo("  License will be added as:")
			formattedLicense := manager.FormatLicenseForFile(fp.config.LicenseText)
			for _, line := range strings.Split(formattedLicense, "\n") {
				fp.logger.LogInfo("    %s", line)
			}
		}

		if fp.config.DryRun {
			fp.logger.LogInfo("Would add license to %s", file)
			continue
		}

		if fp.config.Prompt && !fp.logger.Prompt(
			fp.logger.LogQuestion("Add license to %s?", file)) {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s", file)
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to write file %s: %v", file, err)
			continue
		}

		fp.stats["added"]++
		fp.logger.LogSuccess("Added license to %s", file)
	}

	fp.logger.PrintStats(fp.stats)
	return nil
}

// Remove removes license headers from files
func (fp *FileProcessor) Remove() error {
	fp.resetStats()

	files, err := fp.fileHandler.FindFiles(fp.config.Input)
	if err != nil {
		return err
	}

	// Print scanning message with patterns
	fp.logger.LogInfo("Scanning %d Directories:", len(files))
	fp.logger.LogInfo("Inputs Patterns:")
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			fp.logger.LogInfo("  %s", pattern)
		}
	}
	if fp.config.Skip != "" {
		fp.logger.LogInfo("Skips Patterns:")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				fp.logger.LogInfo("  %s", pattern)
			}
		}
	}

	if fp.config.Verbose {
		fp.logger.LogInfo("Using style: %s", styles.Get(fp.config.PresetStyle).Name)
	}

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, _ := fp.createManager(file)
		newContent, err := manager.RemoveLicense(content)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to remove license from %s: %v", file, err)
			continue
		}

		if newContent == content {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s", file)
			continue
		}

		if fp.config.DryRun {
			fp.logger.LogInfo("Would remove license from %s", file)
			continue
		}

		if fp.config.Prompt && !fp.logger.Prompt(
			fp.logger.LogQuestion("Remove license from %s?", file)) {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s", file)
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to write file %s: %v", file, err)
			continue
		}

		fp.stats["added"]++
		fp.logger.LogSuccess("Removed license from %s", file)
	}

	fp.logger.PrintStats(fp.stats)
	return nil
}

// Update updates license headers in files
func (fp *FileProcessor) Update() error {
	fp.resetStats()

	files, err := fp.fileHandler.FindFiles(fp.config.Input)
	if err != nil {
		return err
	}

	// Print scanning message with patterns
	fp.logger.LogInfo("Scanning %d Directories:", len(files))
	fp.logger.LogInfo("Inputs Patterns:")
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			fp.logger.LogInfo("  %s", pattern)
		}
	}
	if fp.config.Skip != "" {
		fp.logger.LogInfo("Skips Patterns:")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				fp.logger.LogInfo("  %s", pattern)
			}
		}
	}

	if fp.config.Verbose {
		fp.logger.LogInfo("Using style: %s", styles.Get(fp.config.PresetStyle).Name)
	}

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, _ := fp.createManager(file)
		status := manager.CheckLicenseStatus(content)

		if status == license.NoLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s", file)
			continue
		}

		newContent, err := manager.UpdateLicense(content)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to update license in %s: %v", file, err)
			continue
		}

		if newContent == content {
			fp.stats["unchanged"]++
			fp.logger.LogInfo("License is up-to-date in %s", file)
			continue
		}

		if fp.config.DryRun {
			fp.logger.LogInfo("Would update license in %s", file)
			continue
		}

		if fp.config.Prompt && !fp.logger.Prompt(
			fp.logger.LogQuestion("Update license in %s?", file)) {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s", file)
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to write file %s: %v", file, err)
			continue
		}

		fp.stats["added"]++
		fp.logger.LogSuccess("Updated license in %s", file)
	}

	fp.logger.PrintStats(fp.stats)
	return nil
}

// Check checks license headers in files
func (fp *FileProcessor) Check() error {
	fp.resetStats()

	files, err := fp.fileHandler.FindFiles(fp.config.Input)
	if err != nil {
		return err
	}

	// Print scanning message with patterns
	fp.logger.LogInfo("Scanning %d Directories:", len(files))
	fp.logger.LogInfo("Inputs Patterns:")
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			fp.logger.LogInfo("  %s", pattern)
		}
	}
	if fp.config.Skip != "" {
		fp.logger.LogInfo("Skips Patterns:")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				fp.logger.LogInfo("  %s", pattern)
			}
		}
	}

	hasFailures := false
	for _, file := range files {
		// Get relative path for cleaner output
		relPath, err := filepath.Rel(".", file)
		if err != nil {
			relPath = file // Fallback to full path
		}

		manager, _ := fp.createManager(file)

		// Print file path first for better log readability
		if fp.config.Verbose {
			fp.logger.LogInfo("Checking %s...", relPath)
		}

		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read %s: %v", relPath, err)
			hasFailures = true
			continue
		}

		status := manager.CheckLicenseStatus(content)

		// Clear status indication for each file
		switch status {
		case license.NoLicense:
			fp.stats["skipped"]++
			fp.logger.LogError("%s: No license found", relPath)
			hasFailures = true

		case license.DifferentLicense:
			fp.stats["failed"]++
			fp.logger.LogError("%s: License text differs from expected", relPath)
			hasFailures = true

			if fp.config.Verbose {
				// In verbose mode, show actual differences
				current, expected := manager.GetLicenseComparison(content)
				fp.logger.LogInfo("Current license in %s:", relPath)
				fp.logger.LogInfo("%s", current)
				fp.logger.LogInfo("Expected:")
				fp.logger.LogInfo("%s", expected)
			}

		case license.MatchingLicense:
			fp.stats["unchanged"]++
			fp.logger.LogSuccess("%s: License OK", relPath)
		}
	}
	if hasFailures {
		if !fp.config.IgnoreFail {
			return fmt.Errorf("license check failed: some files have missing or incorrect licenses")
		}
	}
	//if hasFailures {
	//	if !fp.config.IgnoreFail {
	//		// Print a final summary without all the cobra usage info
	//		return NewCheckError("license check failed: some files have missing or incorrect licenses")
	//	}
	//}

	return nil
}

// describeCommentStyle returns a human-readable description of the comment style
func describeCommentStyle(cs styles.CommentLanguage) string {
	var parts []string
	if cs.MultiStart != "" {
		parts = append(parts, fmt.Sprintf("multi-line (%s...%s)", cs.MultiStart, cs.MultiEnd))
	}
	if cs.Single != "" {
		parts = append(parts, fmt.Sprintf("single-line (%s)", cs.Single))
	}
	if cs.MultiPrefix != "" {
		parts = append(parts, fmt.Sprintf("prefix: %q", cs.MultiPrefix))
	}
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.Join(parts, ", ")
}
