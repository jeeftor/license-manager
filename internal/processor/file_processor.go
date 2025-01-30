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

	// Read File content
	content, err := fp.fileHandler.ReadFile(file)
	if err != nil {
		return nil, commentStyle, err
	}

	// Create single manager with the actual license text
	lm := license.NewLicenseManager(
		fp.logger,
		fp.config.LicenseText,
		ext,
		styles.Get(fp.config.PresetStyle),
		commentStyle,
	)

	// Set License Mangaer content
	lm.SetFileContent(content)

	// Scan for license stuff
	analysis := lm.SearchForLicense(content)

	if analysis.HasLicense && analysis.IsStyleMatch {
		// Update style if none was explicitly configured
		if fp.config.PresetStyle == "" {
			lm.SetHeaderStyle(analysis.Style)
			fp.logger.LogInfo("  Using detected style: %s", analysis.Style.Name)
		} else {
			fp.logger.LogInfo("  Using configured style: %s", fp.config.PresetStyle)
		}
	} else {
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

// logInputPatterns logs the input scan patterns
func (fp *FileProcessor) logInputPatterns(fileCount int) {
	var debugMsg strings.Builder
	debugMsg.WriteString(fmt.Sprintf(" Scanning %d inputs ", fileCount))
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			coloredPattern := color.New(color.FgHiCyan).Sprint(pattern)
			debugMsg.WriteString(fmt.Sprintf("%s ", coloredPattern))
		}
	}
	fp.logger.LogDebug("%s", debugMsg.String())

	// Handle skip patterns if any
	if fp.config.Skip != "" {
		var skipMsg strings.Builder
		skipMsg.WriteString("Skip Patterns:\n")
		for _, pattern := range strings.Split(fp.config.Skip, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				coloredPattern := color.New(color.FgYellow).Sprint(pattern)
				skipMsg.WriteString(fmt.Sprintf("      %s\n", coloredPattern))
			}
		}
		fp.logger.LogInfo(skipMsg.String())
	}
}

// confirmAction checks if an action should proceed based on dry-run and prompt settings
func (fp *FileProcessor) confirmAction(action, file string) bool {
	if fp.config.DryRun {
		fp.logger.LogInfo("Would %s license in %s", action, file)
		return false
	}

	if fp.config.Prompt && !fp.logger.Prompt(
		fp.logger.LogQuestion("%s license in %s?", action, file)) {
		fp.stats["skipped"]++
		fp.logger.LogInfo("Skipping %s", file)
		return false
	}

	return true
}

// handleFileError logs file errors and updates stats
func (fp *FileProcessor) handleFileError(file, operation string, err error) bool {
	fp.stats["failed"]++
	fp.logger.LogError("Failed to %s %s: %v", operation, file, err)
	return false
}

// prepareOperation sets up common operation requirements
func (fp *FileProcessor) prepareOperation() ([]string, error) {
	fp.resetStats()

	files, err := fp.fileHandler.FindFiles(fp.config.Input)
	if err != nil {
		return nil, err
	}

	fp.logInputPatterns(len(files))
	return files, nil
}

// Add adds license headers to files
func (fp *FileProcessor) Add() error {
	files, err := fp.prepareOperation()
	if err != nil {
		return err
	}

	for _, file := range files {
		manager, commentStyle, err := fp.createLicenseManager(file)
		if err != nil {
			fp.handleFileError(file, "process", err)
			continue
		}

		if manager.HasInitialLicense {
			fp.stats["existing"]++
			if !fp.config.IsPreCommit {
				fp.logger.LogWarning("License already exists in %s", file)
			}
			continue
		}

		newContent, err := manager.AddLicense(manager.InitialComponents.Rest, commentStyle.Language)
		if err != nil {
			fp.handleFileError(file, "add license to", err)
			continue
		}

		if manager.InitialComponents.Preamble != "" {
			newContent = manager.InitialComponents.Preamble + "\n" + newContent
		}

		// Debug the actual comment being added in verbose mode
		fp.logger.LogInfo("  License will be added as:")
		formattedLicense := manager.FormatLicenseForFile(fp.config.LicenseText)
		for _, line := range strings.Split(formattedLicense, "\n") {
			fp.logger.LogInfo("    %s", line)
		}

		if !fp.confirmAction("add", file) {
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.handleFileError(file, "write", err)
			continue
		}

		fp.stats["added"]++
		fp.logger.LogSuccess("Added license to %s", file)
	}

	return nil
}

// Remove removes license headers from files
func (fp *FileProcessor) Remove() error {
	files, err := fp.prepareOperation()
	if err != nil {
		return err
	}

	for _, file := range files {
		manager, commentStyle, err := fp.createLicenseManager(file)
		if err != nil {
			fp.handleFileError(file, "process", err)
			continue
		}

		if !manager.HasInitialLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("No license found in %s", file)
			continue
		}

		newContent, err := manager.RemoveLicense(manager.FileContent, commentStyle.Language)
		if err != nil {
			fp.handleFileError(file, "remove license from", err)
			continue
		}

		if newContent == manager.FileContent {
			fp.stats["unchanged"]++
			fp.logger.LogInfo("No changes needed for %s", file)
			continue
		}

		if !fp.confirmAction("remove", file) {
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.handleFileError(file, "write", err)
			continue
		}

		fp.stats["removed"]++
		fp.logger.LogSuccess("Removed license from %s", file)
	}

	fp.logger.PrintStats(fp.stats, "Removed")
	return nil
}

// Update updates license headers in files
func (fp *FileProcessor) Update() error {
	files, err := fp.prepareOperation()
	if err != nil {
		return err
	}

	for _, file := range files {
		manager, commentStyle, err := fp.createLicenseManager(file)
		if err != nil {
			fp.handleFileError(file, "process", err)
			continue
		}

		status := manager.CheckLicenseStatus(manager.FileContent)
		if status == license.NoLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("Skipping %s (no license)", file)
			continue
		}

		if status == license.FullMatch {
			fp.stats["unchanged"]++
			fp.logger.LogInfo("License is up-to-date in %s", file)
			continue
		}

		newContent, err := manager.UpdateLicense(manager.FileContent, commentStyle.Language)
		if err != nil {
			fp.handleFileError(file, "update license in", err)
			continue
		}

		if !fp.confirmAction("update", file) {
			continue
		}

		if err := fp.fileHandler.WriteFile(file, newContent); err != nil {
			fp.handleFileError(file, "write", err)
			continue
		}

		fp.stats["updated"]++
		fp.logger.LogSuccess("Updated license in %s", file)
	}

	fp.logger.PrintStats(fp.stats, "Updated")
	return nil
}

// Check verifies license headers in files
func (fp *FileProcessor) Check() error {
	files, err := fp.prepareOperation()
	if err != nil {
		return err
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
			fp.logger.LogError("Failed to process %s: %v", relPath, err)
			return NewCheckError(license.NoLicense, fmt.Sprintf("failed to process file: %v", err))
		}

		status := manager.CheckLicenseStatus(manager.FileContent)
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

	fp.logger.PrintStats(fp.stats, "Checked")
	return nil
}
