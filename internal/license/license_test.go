package license

import (
	"bytes"
	"license-manager/internal/processor"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommentStylePreference tests that the PreferMulti flag correctly affects the comment style
func TestCommentStylePreference(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "license-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Read license text from LICENSE file
	licenseText, err := os.ReadFile("../../templates/licenses/mit.txt")
	if err != nil {
		t.Fatalf("Failed to read LICENSE file: %v", err)
	}

	// Replace placeholders in license text
	licenseText = bytes.ReplaceAll(licenseText, []byte("[year]"), []byte("2025"))
	licenseText = bytes.ReplaceAll(licenseText, []byte("[fullname]"), []byte("Test User"))

	tests := []struct {
		name        string
		filename    string
		preferMulti bool
		wantStyle   string
	}{
		{
			name:        "Go with multi-line comments",
			filename:    "test.go",
			preferMulti: true,
			wantStyle:   "/*",
		},
		{
			name:        "Go with single-line comments",
			filename:    "test.go",
			preferMulti: false,
			wantStyle:   "//",
		},
		{
			name:        "Python always uses single-line",
			filename:    "test.py",
			preferMulti: true, // Should be ignored for Python
			wantStyle:   "#",
		},
		{
			name:        "HTML always uses multi-line",
			filename:    "test.html",
			preferMulti: false, // Should be ignored for HTML
			wantStyle:   "<!--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.filename)
			err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create license manager with comment style preference
			style := GetPresetStyle("simple")
			commentStyle := processor.getCommentStyle(tt.filename)
			commentStyle.PreferMulti = tt.preferMulti
			lm := NewLicenseManager(style, string(licenseText), commentStyle)

			// Add license to file
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			// Add license
			result := lm.AddLicense(string(content))

			// Verify comment style
			if !strings.Contains(result, tt.wantStyle) {
				t.Errorf("Expected comment style %q not found in result:\n%s", tt.wantStyle, result)
			}

			// For Go files, verify the opposite style is not present
			if tt.filename == "test.go" {
				oppositeStyle := "//"
				if !tt.preferMulti {
					oppositeStyle = "/*"
				}
				if strings.Contains(result, oppositeStyle) {
					t.Errorf("Unexpected comment style %q found in result:\n%s", oppositeStyle, result)
				}
			}
		})
	}
}

// TestCommentStyleUpdate tests that updating a license preserves the original comment style
func TestCommentStyleUpdate(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "license-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Read initial license text
	licenseText, err := os.ReadFile("../../templates/licenses/mit.txt")
	if err != nil {
		t.Fatalf("Failed to read LICENSE file: %v", err)
	}

	// Replace placeholders in initial license text
	licenseText = bytes.ReplaceAll(licenseText, []byte("[year]"), []byte("2025"))
	licenseText = bytes.ReplaceAll(licenseText, []byte("[fullname]"), []byte("Test User"))

	// Create updated license text with different year
	updatedLicenseText := bytes.ReplaceAll(licenseText, []byte("2025"), []byte("2026"))

	tests := []struct {
		name          string
		filename      string
		initialMulti  bool
		updateMulti   bool
		wantStyle     string
		wantPreserved bool
	}{
		{
			name:          "Update Go multi-line to multi-line",
			filename:      "test.go",
			initialMulti:  true,
			updateMulti:   true,
			wantStyle:     "/*",
			wantPreserved: true,
		},
		{
			name:          "Update Go single-line to single-line",
			filename:      "test.go",
			initialMulti:  false,
			updateMulti:   false,
			wantStyle:     "//",
			wantPreserved: true,
		},
		{
			name:          "Update Go multi-line to single-line",
			filename:      "test.go",
			initialMulti:  true,
			updateMulti:   false,
			wantStyle:     "//",
			wantPreserved: false,
		},
		{
			name:          "Update Go single-line to multi-line",
			filename:      "test.go",
			initialMulti:  false,
			updateMulti:   true,
			wantStyle:     "/*",
			wantPreserved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.filename)
			err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create initial license manager
			style := GetPresetStyle("simple")
			commentStyle := processor.getCommentStyle(tt.filename)
			commentStyle.PreferMulti = tt.initialMulti
			lm := NewLicenseManager(style, string(licenseText), commentStyle)

			// Add initial license
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			initial := lm.AddLicense(string(content))

			// Create update license manager with different preference
			commentStyle.PreferMulti = tt.updateMulti
			lm = NewLicenseManager(style, string(licenseText), commentStyle)

			// Update license with new text
			updated := lm.UpdateLicense(initial, string(updatedLicenseText))

			// Verify final comment style
			if !strings.Contains(updated, tt.wantStyle) {
				t.Errorf("Expected comment style %q not found in result:\n%s", tt.wantStyle, updated)
			}

			// Verify the license was actually updated
			if !strings.Contains(updated, "2026") {
				t.Errorf("License text was not updated properly, expected year 2026 not found in:\n%s", updated)
			}

			// Verify preservation of style if expected
			initialStyle := "/*"
			if !tt.initialMulti {
				initialStyle = "//"
			}
			hasInitialStyle := strings.Contains(updated, initialStyle)
			if tt.wantPreserved && !hasInitialStyle {
				t.Errorf("Expected original style %q to be preserved in result:\n%s", initialStyle, updated)
			} else if !tt.wantPreserved && hasInitialStyle && initialStyle != tt.wantStyle {
				t.Errorf("Expected original style %q to be changed in result:\n%s", initialStyle, updated)
			}
		})
	}
}
