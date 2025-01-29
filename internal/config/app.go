// internal/config/app.go
package config

import (
	"github.com/jeeftor/license-manager/internal/errors"
	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/processor"
	"os"
	"path/filepath"
)

// AppConfig holds CLI and application-level configuration
type AppConfig struct {
	// File paths
	LicenseFile string // Path to license template file
	Inputs      string // Inputs file patterns
	Skips       string // Skips patterns

	// UI/Behavior settings
	LogLevel    logger.LogLevel
	Interactive bool
	Force       bool

	// Style preferences
	HeaderStyle       string
	CommentStyle      string
	PreferMulti       *bool
	IgnoreFail        bool
	ForceCommentStyle force.ForceCommentStyle
}

// NewAppConfig returns default application config
func NewAppConfig() AppConfig {
	return AppConfig{
		HeaderStyle:  "hash",
		CommentStyle: "go",
	}
}

// ToProcessorConfig converts AppConfig to processor.Config
func (c *AppConfig) ToProcessorConfig() (*processor.Config, error) {
	var licenseText string
	var err error

	// Load license file if provided
	if c.LicenseFile != "" {
		licenseText, err = c.loadLicenseFile()
		if err != nil {
			return nil, err
		}
	}

	// Convert to processor config
	return &processor.Config{
		LicenseText: licenseText,
		Input:       c.Inputs,
		Skip:        c.Skips,
		Prompt:      c.Interactive,

		PresetStyle:       c.HeaderStyle,
		ForceCommentStyle: c.ForceCommentStyle,
		IgnoreFail:        c.IgnoreFail,
		LogLevel:          c.LogLevel,
	}, nil
}

func (c *AppConfig) loadLicenseFile() (string, error) {
	if !filepath.IsAbs(c.LicenseFile) {
		abs, err := filepath.Abs(c.LicenseFile)
		if err != nil {
			return "", errors.NewValidationError("invalid license file path", "LicenseFile")
		}
		c.LicenseFile = abs
	}

	content, err := os.ReadFile(c.LicenseFile)
	if err != nil {
		return "", errors.NewValidationError("failed to read license file", "LicenseFile")
	}

	return string(content), nil
}
