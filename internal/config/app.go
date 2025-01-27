// internal/config/app.go
package config

import (
	"license-manager/internal/errors"
	"license-manager/internal/logger"
	"license-manager/internal/processor"
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
	Verbose     bool
	LogLevel    logger.LogLevel
	Interactive bool
	Force       bool

	// Style preferences
	HeaderStyle  string
	CommentStyle string
	PreferMulti  bool
	IgnoreFail   bool
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
		Verbose:     c.Verbose,
		PresetStyle: c.HeaderStyle,
		PreferMulti: c.PreferMulti,
		IgnoreFail:  c.IgnoreFail,
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
