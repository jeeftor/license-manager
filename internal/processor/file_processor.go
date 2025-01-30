package processor

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/license"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
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
	log := logger.NewLogger(cfg.LogLevel)
	fh := NewFileHandler(log)
	fh.SetSkipPattern(cfg.Skip) // Set the skip pattern
	return &FileProcessor{
		config:      cfg,
		logger:      log,
		fileHandler: fh,
		stats:       make(map[string]int),
	}
}

func (fp *FileProcessor) createLicenseManager(file string) (*license.LicenseManager, styles.CommentLanguage, error) {
	// Get comment headerFooterStyle for file type
	ext := filepath.Ext(file)
	commentStyle := styles.GetLanguageCommentStyle(ext)

	if fp.config.ForceCommentStyle == force.Single {
		fp.logger.LogWarning("Overriding default comment headerFooterStyle to %s", fp.config.ForceCommentStyle)
		commentStyle.PreferMulti = false
	} else if fp.config.ForceCommentStyle == force.Multi {
		fp.logger.LogWarning("Overriding default comment headerFooterStyle to %s", fp.config.ForceCommentStyle)
		commentStyle.PreferMulti = true
	}

	fp.logger.LogInfo("Processing file: %s", file)
	fp.logger.LogInfo("  Language: %s", commentStyle.Language)
	fp.logger.LogInfo("  Comment headerFooterStyle: %s", describeCommentStyle(commentStyle))

	// Create single manager with the actual license text
	lm := license.NewLicenseManager(
		fp.logger,
		fp.config.LicenseText,
		ext,
		styles.Get(fp.config.PresetStyle),
		commentStyle,
	)

	// Read content into return structure
	content, err := fp.fileHandler.ReadFile(file)
	if err != nil {
		fp.logger.LogInfo("  Using configured style: %s", fp.config.PresetStyle)
		return lm, commentStyle, err
	}

	// Set content
	lm.SetFileContent(content)

	// Scan for existing license code in the content
	success, components := lm.HasLicense(content)
	if success {
		// Detect header and footer style since we found a license block
		headerFooterStyle, foundMatch := lm.DetectHeaderAndFooterStyle(components.Header, components.Footer)
		if !foundMatch {
			return lm, commentStyle, nil
		}
		fp.logger.LogInfo("  Detected style: %s", headerFooterStyle.Name)
		// Update manager's style if no style was explicitly configured
		if fp.config.PresetStyle == "" {
			lm.SetHeaderStyle(headerFooterStyle)
			fp.logger.LogInfo("  Using detected style: %s", headerFooterStyle.Name)
		} else {
			fp.logger.LogInfo("  Using configured style: %s", fp.config.PresetStyle)
		}
	} else {
		// No license block detected.
		fp.logger.LogInfo("  Using configured style: %s", fp.config.PresetStyle)
	}

	return lm, commentStyle, nil
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

	fp.logger.LogDebug("Using style: %s", style.Name)

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, commentStyle, _ := fp.createLicenseManager(file)

		hasLicense, extract := manager.HasLicense(content)
		if hasLicense {
			fp.stats["existing"]++
			if !fp.config.IsPreCommit {
				fp.logger.LogWarning("License already exists in %s", file)
			}
			continue
		}

		// Use rest instead of content to add license after preamble
		newContent, err := manager.AddLicense(extract.Rest, commentStyle.Language)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to add license to %s: %v", file, err)
			continue
		}

		// Recombine preamble with the licensed content
		if extract.Preamble != "" {
			newContent = extract.Preamble + "\n" + newContent
		}

		// Debug the actual comment being added in verbose mode
		fp.logger.LogInfo("  License will be added as:")
		formattedLicense := manager.FormatLicenseForFile(fp.config.LicenseText)
		for _, line := range strings.Split(formattedLicense, "\n") {
			fp.logger.LogInfo("    %s", line)
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
	fp.logger.LogInfo("Input Patterns:")
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

	fp.logger.LogInfo("Using style: %s", styles.Get(fp.config.PresetStyle).Name)

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, commentStyle, _ := fp.createLicenseManager(file)

		newContent, err := manager.RemoveLicense(content, commentStyle.Language)
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
	fp.logger.LogInfo("Inputs Patterns: %s", fp.config.Input)
	if fp.config.Skip != "" {
		fp.logger.LogInfo("Skips Patterns: %s", fp.config.Skip)
	}

	fp.logger.LogInfo("Using style: %s", styles.Get(fp.config.PresetStyle).Name)

	for _, file := range files {
		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager, commentStyle, _ := fp.createLicenseManager(file)

		status := manager.CheckLicenseStatus(content)

		// Skip files with no license
		if status == license.NoLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s (no license)", file)
			continue
		}

		// Update files with incorrect style or content
		if status == license.StyleMismatch || status == license.ContentMismatch || status == license.ContentAndStyleMismatch {
			newContent, err := manager.UpdateLicense(content, commentStyle.Language)
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
	fp.logger.LogInfo("Input Patterns:")
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

		manager, _, err := fp.createLicenseManager(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read %s: %v", relPath, err)
			return NewCheckError(license.NoLicense, fmt.Sprintf("failed to read file: %v", err))
		}
		content := manager.FileContent

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
		fp.logger.LogSuccess("%s: License OK", relPath)

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
