// internal/processor/config.go
package processor

import (
	"license-manager/internal/force"
	"license-manager/internal/logger"
)

// Config holds the configuration for the file processor
type Config struct {
	// Core processing needs
	LicenseText string // The actual license text content
	Input       string // Input file patterns
	Skip        string // Patterns to skip
	PresetStyle string // Header/Footer style to use

	// Processing behavior
	Prompt            bool // Whether to prompt before changes
	DryRun            bool // Whether to show what would be done without doing it
	LogLevel          logger.LogLevel
	IgnoreFail        bool // Whether to return success even if checks fail
	ForceCommentStyle force.ForceCommentStyle
}
