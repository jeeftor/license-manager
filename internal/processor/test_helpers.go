package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
)

// TestHelper provides utilities for testing license operations
type TestHelper struct {
	t           *testing.T
	tmpDir      string
	licenseText string
	licenseFile string
}

// NewTestHelper creates a new test helper with a temporary directory
func NewTestHelper(t *testing.T, licenseText string) *TestHelper {
	tmpDir := t.TempDir()
	licenseFile := filepath.Join(tmpDir, "LICENSE")

	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to create license file: %v", err)
	}

	return &TestHelper{
		t:           t,
		tmpDir:      tmpDir,
		licenseText: licenseText,
		licenseFile: licenseFile,
	}
}

// CreateFile creates a test file with the given content
func (h *TestHelper) CreateFile(filename, content string) string {
	fullPath := filepath.Join(h.tmpDir, filename)

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		h.t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		h.t.Fatalf("Failed to create file %s: %v", filename, err)
	}

	return fullPath
}

// AddLicenseToFile adds a license to a file using the system
// This ensures the license is in the correct format with unicode markers
func (h *TestHelper) AddLicenseToFile(filePath string) {
	cfg := &Config{
		LicenseText:       h.licenseText,
		Input:             filePath,
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		h.t.Fatalf("Failed to add license to %s: %v", filePath, err)
	}
}

// ReadFile reads the content of a file
func (h *TestHelper) ReadFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.t.Fatalf("Failed to read file %s: %v", filePath, err)
	}
	return string(content)
}

// TmpDir returns the temporary directory path
func (h *TestHelper) TmpDir() string {
	return h.tmpDir
}

// LicenseText returns the license text
func (h *TestHelper) LicenseText() string {
	return h.licenseText
}

// CreateProcessor creates a new processor with the helper's configuration
func (h *TestHelper) CreateProcessor(
	input string,
	commentStyle force.ForceCommentStyle,
) *FileProcessor {
	cfg := &Config{
		LicenseText:       h.licenseText,
		Input:             input,
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: commentStyle,
		LogLevel:          logger.ErrorLevel,
	}
	return NewFileProcessor(cfg)
}
