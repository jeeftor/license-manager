// internal/config/app.go
package config

import (
	"license-manager/internal/errors"
	"license-manager/internal/processor"
	"os"
	"path/filepath"
)

// AppConfig holds CLI and application-level configuration
type AppConfig struct {
	// File paths
	LicenseFile string // Path to license template file
	Input       string // Input file patterns
	Skip        string // Skip patterns

	// UI/Behavior settings
	Verbose     bool
	Interactive bool
	DryRun      bool
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
		HeaderStyle:  "simple",
		CommentStyle: "go",
	}
}

// ToProcessorConfig converts AppConfig to processor.Config
func (c *AppConfig) ToProcessorConfig() (*processor.Config, error) {
	// Validate and load license file
	if c.LicenseFile == "" {
		return nil, errors.NewValidationError("license file is required", "LicenseFile")
	}

	licenseText, err := c.loadLicenseFile()
	if err != nil {
		return nil, err
	}

	// Convert to processor config
	return &processor.Config{
		LicenseText: licenseText,
		Input:       c.Input,
		Skip:        c.Skip,
		Prompt:      c.Interactive,
		DryRun:      c.DryRun,
		Verbose:     c.Verbose,
		PresetStyle: c.HeaderStyle,
		PreferMulti: c.PreferMulti,
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
