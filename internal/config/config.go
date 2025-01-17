package config

import (
	"license-manager/internal/errors"
	"license-manager/internal/styles"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// Input patterns for files to process (comma-separated)
	Input string
	// Path to the license file
	LicenseFile string
	// Comment style to use
	CommentStyle string
	// Header/footer style to use
	HeaderStyle string
	// Whether to run in verbose mode
	Verbose bool
	// Whether to prompt for confirmation
	Interactive bool
	// Whether to force operations without checking
	Force bool
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Input == "" {
		return errors.NewValidationError("input pattern is required", "Input")
	}

	if c.LicenseFile != "" {
		if !filepath.IsAbs(c.LicenseFile) {
		    var err error
		    c.LicenseFile, err = filepath.Abs(c.LicenseFile)
		    if err != nil {
			    return errors.NewValidationError("invalid license file path", "LicenseFile")
		    }
		}
	}

	// Validate comment style
	if c.CommentStyle != "" {
		validStyles := []string{"c", "cpp", "python", "html", "shell"}
		found := false
		for _, style := range validStyles {
			if strings.EqualFold(c.CommentStyle, style) {
				found = true
				break
			}
		}
		if !found {
			return errors.NewValidationError("invalid comment style", "CommentStyle")
		}
	}

	// Validate header style
	if c.HeaderStyle != "" {
		style := styles.Get(c.HeaderStyle)
		if style.Name == "" {
			return errors.NewValidationError("invalid header style", "HeaderStyle")
		}
	}

	return nil
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		CommentStyle: "c", // Default to C-style comments
		HeaderStyle:  "simple",
		Interactive: true,
		Force:       false,
	}
}
