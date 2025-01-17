// internal/processor/file_processor.go
package processor

import (
	"license-manager/internal/config"
	"license-manager/internal/errors"
	"license-manager/internal/license"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
)

// FileProcessor handles license operations on files
type FileProcessor struct {
	config      *config.Config
	license     string
	style       styles.HeaderFooterStyle
	fileHandler *FileHandler
	logger      *logger.Logger
	stats       map[string]int
}

// NewFileProcessor creates a new FileProcessor instance
func NewFileProcessor(config *config.Config, license string, style styles.HeaderFooterStyle) *FileProcessor {
	log := logger.NewLogger(config.Verbose)
	return &FileProcessor{
		config:      config,
		license:     license,
		style:       style,
		logger:      log,
		fileHandler: NewFileHandler(log),
		stats:       make(map[string]int),
	}
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

	for _, file := range files {
		fp.logger.LogVerbose("Processing file: %s", file)

		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager := license.NewManager(fp.license, fp.style)
		if manager.HasLicense(content) {
			fp.stats["existing"]++
			fp.logger.LogWarning("License already exists in %s", file)
			continue
		}

		newContent, err := manager.AddLicense(content)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to add license to %s: %v", file, err)
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

	for _, file := range files {
		fp.logger.LogVerbose("Processing file: %s", file)

		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager := license.NewManager(fp.license, fp.style)
		newContent, err := manager.RemoveLicense(content)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to remove license from %s: %v", file, err)
			continue
		}

		if newContent == content {
			fp.stats["skipped"]++
			fp.logger.LogInfo("No license found in %s", file)
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

	for _, file := range files {
		fp.logger.LogVerbose("Processing file: %s", file)

		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			continue
		}

		manager := license.NewManager(fp.license, fp.style)
		status := manager.CheckLicenseStatus(content)

		if status == license.NoLicense {
			fp.stats["skipped"]++
			fp.logger.LogInfo("No license found in %s", file)
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

	hasFailures := false
	for _, file := range files {
		fp.logger.LogVerbose("Processing file: %s", file)

		content, err := fp.fileHandler.ReadFile(file)
		if err != nil {
			fp.stats["failed"]++
			fp.logger.LogError("Failed to read file %s: %v", file, err)
			hasFailures = true
			continue
		}

		manager := license.NewManager(fp.license, fp.style)
		status := manager.CheckLicenseStatus(content)

		switch status {
		case license.NoLicense:
			fp.stats["skipped"]++
			fp.logger.LogError("No license found in %s", file)
			hasFailures = true

		case license.DifferentLicense:
			fp.stats["failed"]++
			fp.logger.LogError("License doesn't match in %s", file)
			hasFailures = true

			if fp.config.Verbose {
				current, expected := manager.GetLicenseComparison(content)
				fp.logger.LogInfo("Current license in %s:\n%s", file, current)
				fp.logger.LogInfo("Expected license:\n%s", expected)
			}

		case license.MatchingLicense:
			fp.stats["unchanged"]++
			fp.logger.LogSuccess("License matches in %s", file)
		}
	}

	fp.logger.PrintStats(fp.stats)
	if hasFailures {
		return errors.NewLicenseError("one or more files have missing or incorrect licenses", "")
	}

	return nil
}
