package processor

import (
	"fmt"
	"github.com/fatih/color"
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
	// Get comment style for file type
	ext := filepath.Ext(file)
	commentStyle := styles.GetLanguageCommentStyle(ext)
	commentStyle.PreferMulti = fp.config.PreferMulti

	if fp.config.Verbose {
		fp.logger.LogInfo("Processing file: %s", file)
		fp.logger.LogInfo("  Language: %s", commentStyle.Language)
		fp.logger.LogInfo("  Comment style: %s", describeCommentStyle(commentStyle))
	}

	// Try to detect style from existing file content
	content, err := fp.fileHandler.ReadFile(file)
	var style styles.HeaderFooterStyle
	if err == nil {
		// Create a temporary manager with default style to detect existing style
		tempManager := license.NewLicenseManager("", styles.Get(fp.config.PresetStyle))
		tempManager.SetCommentStyle(commentStyle)
		if tempManager.HasLicense(content) {
			// If we found a license, detect its style for logging purposes
			style = tempManager.DetectHeaderStyle(content)
			if fp.config.Verbose {
				fp.logger.LogInfo("  Detected style: %s", style.Name)
			}
			// Always use the configured style for checking
			style = styles.Get(fp.config.PresetStyle)
		} else {
			// If no license found, use configured style
			style = styles.Get(fp.config.PresetStyle)
			if fp.config.Verbose {
				fp.logger.LogInfo("  Using configured style: %s", style.Name)
			}
		}
	} else {
		// If can't read file, use configured style
		style = styles.Get(fp.config.PresetStyle)
		if fp.config.Verbose {
			fp.logger.LogInfo("  Using configured style: %s", style.Name)
		}
	}

	// Create and configure manager
	manager := license.NewLicenseManager(fp.config.LicenseText, style)
	manager.SetCommentStyle(commentStyle)
	// Set the appropriate language handler based on the language type
	manager.SetLanguageHandler(language.GetLanguageHandler(commentStyle.Language, style))
	if fp.config.Verbose {
		manager.SetVerbose(true, fp.logger)
	}

	return manager, commentStyle
}

// resetStats resets the operation statistics
func (fp *FileProcessor) resetStats() {
	fp.stats = map[string]int{
		"added":                  0,
		"existing":               0,
		"skipped":                0,
		"failed":                 0,
		"unchanged":              0,
		"ok":                     0,
		"missing":                0,
		"mismatch":               0,
		"style_mismatch":         0,
		"content_style_mismatch": 0,
		"error":                  0,
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
	var debugMsg strings.Builder
	debugMsg.WriteString(fmt.Sprintf(" Scanning %d inputs ", len(files)))
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			coloredPattern := color.New(color.FgHiCyan).Sprint(pattern) // or any other color like "debug", "notice", etc
			debugMsg.WriteString(fmt.Sprintf("%s ", coloredPattern))
		}
	}

	// Log the entire message at once
	fp.logger.LogDebug("%s", debugMsg.String())
	// Handle skip patterns
	if fp.config.Skip != "" {
		var skipMsg strings.Builder
		skipMsg.WriteString("Skip Patterns:\n")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				coloredPattern := color.New(color.FgYellow).Sprint(pattern) // Using yellow for skip patterns
				skipMsg.WriteString(fmt.Sprintf("      %s\n", coloredPattern))
			}
		}
		fp.logger.LogInfo(skipMsg.String())
	}

	style := styles.Get(fp.config.PresetStyle)
	if fp.config.Verbose {
		fp.logger.LogDebug("Using style: %s", style.Name)
		if fp.config.PreferMulti {
			fp.logger.LogDebug("Preferring multi-line comments where supported")
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

	//fp.logger.PrintStats(fp.stats, "Added")
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

	fp.logger.PrintStats(fp.stats, "Removed")
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

		// Skip files with no license
		if status == license.NoLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s (no license)", file)
			continue
		}

		// Update files with incorrect style or content
		if status == license.StyleMismatch || status == license.ContentMismatch || status == license.ContentAndStyleMismatch {
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
		} else {
			fp.stats["unchanged"]++
			fp.logger.LogInfo("License is up-to-date in %s", file)
		}
	}

	fp.logger.PrintStats(fp.stats, "Updated")
	return nil
}

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

	hasNoLicense := false
	hasContentMismatch := false
	hasStyleMismatch := false

	for _, file := range files {
		relPath := file
		if rel, err := filepath.Rel(".", file); err == nil {
			relPath = rel
		}

		manager, _ := fp.createManager(file)
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read %s: %v", relPath, err)
			return NewCheckError(license.NoLicense, fmt.Sprintf("failed to read file: %v", err))
		}

		status := manager.CheckLicenseStatus(content)
		if status != license.FullMatch {
			fp.stats["failed"]++
			switch status {
			case license.NoLicense:
				hasNoLicense = true
				fp.stats["missing"]++
				fp.logger.LogError("%s: Missing license", relPath)
			case license.ContentMismatch:
				hasContentMismatch = true
				fp.logger.LogError("%s: License content mismatch", relPath)
			case license.StyleMismatch:
				hasStyleMismatch = true
				fp.logger.LogError("%s: License style mismatch (expected %s)", relPath, manager.GetHeaderStyle().Name)
			case license.ContentAndStyleMismatch:
				hasContentMismatch = true
				hasStyleMismatch = true
				fp.logger.LogError("%s: License content and style mismatch", relPath)
			default:
				fp.logger.LogError("%s: Unknown license error", relPath)
			}
			continue
		}
		fp.stats["passed"]++
		if fp.config.Verbose {
			fp.logger.LogSuccess("%s: License OK", relPath)
		}
	}

	if hasNoLicense {
		return NewCheckError(license.NoLicense, "license check failed: some files have missing licenses")
	}
	if hasContentMismatch && hasStyleMismatch {
		return NewCheckError(license.ContentAndStyleMismatch, "license check failed: some files have content and style mismatches")
	}
	if hasContentMismatch {
		return NewCheckError(license.ContentMismatch, "license check failed: some files have content mismatches")
	}
	if hasStyleMismatch {
		return NewCheckError(license.StyleMismatch, "license check failed: some files have style mismatches")
	}

	return nil
}

// Check checks license headers in files
//func (fp *FileProcessor) Check() error {
//	fp.resetStats()
//
//	files, err := fp.fileHandler.FindFiles(fp.config.Input)
//	if err != nil {
//		return err
//	}
//
//	// Print scanning message with patterns
//	fp.logger.LogInfo("Scanning %d Directories:", len(files))
//	fp.logger.LogInfo("Inputs Patterns:")
//	for _, pattern := range strings.Split(fp.config.Input, ",") {
//		pattern = strings.TrimSpace(pattern)
//		if pattern != "" {
//			fp.logger.LogInfo("  %s", pattern)
//		}
//	}
//	if fp.config.Skip != "" {
//		fp.logger.LogInfo("Skips Patterns:")
//		for _, pattern := range strings.Split(fp.config.Skip, ",") {
//			pattern = strings.TrimSpace(pattern)
//			if pattern != "" {
//				fp.logger.LogInfo("  %s", pattern)
//			}
//		}
//	}
//
//	hasFailures := false
//	for _, file := range files {
//		// Get relative path for logging
//		relPath := file
//		if rel, err := filepath.Rel(".", file); err == nil {
//			relPath = rel
//		}
//
//		// Check license status
//		manager, _ := fp.createManager(file)
//		content, err := fp.fileHandler.ReadFile(file)
//		if err != nil {
//			fp.stats["failed"]++
//			fp.logger.LogError("Failed to read %s: %v", relPath, err)
//			return NewCheckError(license.NoLicense, fmt.Sprintf("failed to read file: %v", err))
//		}
//
//		status := manager.CheckLicenseStatus(content)
//		if status != license.FullMatch {
//			fp.stats["failed"]++
//			switch status {
//			case license.NoLicense:
//				fp.logger.LogError("%s: Missing license", relPath)
//			case license.ContentMismatch:
//				fp.logger.LogError("%s: License content mismatch", relPath)
//			case license.StyleMismatch:
//				fp.logger.LogError("%s: License style mismatch (expected %s)", relPath, manager.GetHeaderStyle().Name)
//			case license.ContentAndStyleMismatch:
//				fp.logger.LogError("%s: License content and style mismatch", relPath)
//			default:
//				fp.logger.LogError("%s: Unknown license error", relPath)
//			}
//			hasFailures = true
//			continue
//		} else {
//			fp.stats["passed"]++
//			if fp.config.Verbose {
//				fp.logger.LogSuccess("%s: License OK", relPath)
//			}
//		}
//	}
//	if hasFailures {
//		if fp.stats["missing"] > 0 {
//			return NewCheckError(license.NoLicense, "license check failed: some files have missing licenses")
//		}
//		return NewCheckError(license.ContentMismatch, "license check failed: some files have incorrect licenses")
//	}
//
//	return nil
//}

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
