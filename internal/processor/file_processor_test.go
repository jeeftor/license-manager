package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
)

// TestAddLicenseToMultipleFiles tests adding licenses to multiple files
func TestAddLicenseToMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test files
	files := map[string]string{
		"file1.go": `package main

func main() {
	println("file1")
}
`,
		"file2.go": `package utils

func Helper() string {
	return "helper"
}
`,
		"subdir/file3.go": `package subpkg

func SubFunc() {}
`,
	}

	// Create files
	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	// Create license file
	licenseFile := filepath.Join(tmpDir, "LICENSE")
	licenseText := "Copyright (c) 2025 Test Corp"
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	// Configure processor
	cfg := &Config{
		LicenseText:       licenseText,
		Input:             filepath.Join(tmpDir, "**/*.go"),
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify all files have licenses
	for path := range files {
		fullPath := filepath.Join(tmpDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", path, err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "Copyright (c) 2025 Test Corp") {
			t.Errorf("File %s does not contain license", path)
		}
	}

	// Verify stats
	if processor.stats["added"] != 3 {
		t.Errorf("Expected 3 files to be processed, got %d", processor.stats["added"])
	}
}

// TestSkipPatterns tests that skip patterns work correctly
func TestSkipPatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"include.go":       "package main\n",
		"skip_me.go":       "package skip\n",
		"vendor/vendor.go": "package vendor\n",
		"normal/normal.go": "package normal\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	licenseText := "Copyright (c) 2025"
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText: licenseText,
		Input:       filepath.Join(tmpDir, "**/*.go"),
		Skip: filepath.Join(
			tmpDir,
			"skip_*.go",
		) + "," + filepath.Join(
			tmpDir,
			"vendor/**",
		),
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify included files have licenses
	for _, path := range []string{"include.go", "normal/normal.go"} {
		fullPath := filepath.Join(tmpDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", path, err)
		}
		if !strings.Contains(string(content), "Copyright") {
			t.Errorf("File %s should have license", path)
		}
	}

	// Verify skipped files don't have licenses
	for _, path := range []string{"skip_me.go", "vendor/vendor.go"} {
		fullPath := filepath.Join(tmpDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", path, err)
		}
		if strings.Contains(string(content), "Copyright") {
			t.Errorf("File %s should NOT have license (should be skipped)", path)
		}
	}
}

// TestUpdateExistingLicense tests updating licenses that already exist
func TestUpdateExistingLicense(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	oldContent := `/*
 * ########################################
 * Copyright (c) 2024 Old Corp
 * ########################################
 */

package main

func main() {}
`

	if err := os.WriteFile(testFile, []byte(oldContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	newLicenseText := "Copyright (c) 2025 New Corp"
	if err := os.WriteFile(licenseFile, []byte(newLicenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText:       newLicenseText,
		Input:             testFile,
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Update(); err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Copyright (c) 2025 New Corp") {
		t.Errorf("File should contain new license")
	}
	if strings.Contains(contentStr, "Copyright (c) 2024 Old Corp") {
		t.Errorf("File should not contain old license")
	}
}

// TestRemoveLicense tests removing licenses from files
func TestRemoveLicense(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	contentWithLicense := `/*
 * ########################################
 * Copyright (c) 2025 Test Corp
 * ########################################
 */

package main

func main() {}
`

	if err := os.WriteFile(testFile, []byte(contentWithLicense), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	licenseText := "Copyright (c) 2025 Test Corp"
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText:       licenseText,
		Input:             testFile,
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Remove(); err != nil {
		t.Fatalf("Remove() failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "Copyright (c) 2025 Test Corp") {
		t.Errorf("License should be removed from file")
	}
	if !strings.Contains(contentStr, "package main") {
		t.Errorf("File content should be preserved")
	}
	if !strings.Contains(contentStr, "func main()") {
		t.Errorf("File content should be preserved")
	}
}

// TestAddLicenseToFileWithExistingLicense tests that files with licenses are skipped
func TestAddLicenseToFileWithExistingLicense(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	contentWithLicense := `/*
 * ########################################
 * Copyright (c) 2025 Test Corp
 * ########################################
 */

package main

func main() {}
`

	if err := os.WriteFile(testFile, []byte(contentWithLicense), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	licenseText := "Copyright (c) 2025 Test Corp"
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText:       licenseText,
		Input:             testFile,
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify file wasn't modified (should be skipped)
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != contentWithLicense {
		t.Errorf("File with existing license should not be modified")
	}

	// Check stats
	if processor.stats["existing"] != 1 {
		t.Errorf("Expected 1 file with existing license, got %d", processor.stats["existing"])
	}
	if processor.stats["added"] != 0 {
		t.Errorf("Expected 0 files added, got %d", processor.stats["added"])
	}
}

// TestCheckLicenseOperation tests the check operation
func TestCheckLicenseOperation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with different license states
	files := map[string]struct {
		content        string
		expectedStatus string
	}{
		"good.go": {
			content: `/*
 * ########################################
 * Copyright (c) 2025 Test Corp
 * ########################################
 */

package main
`,
			expectedStatus: "ok",
		},
		"missing.go": {
			content: `package main

func main() {}
`,
			expectedStatus: "missing",
		},
		"wrong_content.go": {
			content: `/*
 * ########################################
 * Copyright (c) 2024 Old Corp
 * ########################################
 */

package main
`,
			expectedStatus: "mismatch",
		},
	}

	for filename, fileData := range files {
		fullPath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(fullPath, []byte(fileData.content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", filename, err)
		}
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	licenseText := "Copyright (c) 2025 Test Corp"
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText:       licenseText,
		Input:             filepath.Join(tmpDir, "*.go"),
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
		IgnoreFail:        true, // Don't fail on check errors for this test
	}

	processor := NewFileProcessor(cfg)
	_ = processor.Check() // Ignore error for stats checking

	// Verify stats
	if processor.stats["ok"] < 1 {
		t.Errorf("Expected at least 1 file with correct license")
	}
	if processor.stats["missing"] < 1 {
		t.Errorf("Expected at least 1 file with missing license")
	}
}

// TestDifferentFileTypes tests processing different file types
func TestDifferentFileTypes(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]struct {
		content         string
		expectedComment string
	}{
		"test.go": {
			content:         "package main\n",
			expectedComment: "/*",
		},
		"test.py": {
			content:         "def main():\n    pass\n",
			expectedComment: "'''",
		},
		"test.js": {
			content:         "function main() {}\n",
			expectedComment: "/*",
		},
	}

	licenseText := "Copyright (c) 2025 Test Corp"
	licenseFile := filepath.Join(tmpDir, "LICENSE")
	if err := os.WriteFile(licenseFile, []byte(licenseText), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	for filename, fileData := range files {
		fullPath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(fullPath, []byte(fileData.content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", filename, err)
		}

		cfg := &Config{
			LicenseText:       licenseText,
			Input:             fullPath,
			Skip:              "",
			Prompt:            false,
			PresetStyle:       "hash",
			ForceCommentStyle: force.No,
			LogLevel:          logger.ErrorLevel,
		}

		processor := NewFileProcessor(cfg)
		if err := processor.Add(); err != nil {
			t.Fatalf("Add() failed for %s: %v", filename, err)
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", filename, err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, fileData.expectedComment) {
			t.Errorf("File %s should use %s comment style", filename, fileData.expectedComment)
		}
		if !strings.Contains(contentStr, "Copyright (c) 2025 Test Corp") {
			t.Errorf("File %s should contain license", filename)
		}
	}
}

// TestEmptyInputPattern tests handling of empty input patterns
func TestEmptyInputPattern(t *testing.T) {
	cfg := &Config{
		LicenseText:       "Copyright",
		Input:             "",
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	err := processor.Add()

	if err == nil {
		t.Errorf("Expected error for empty input pattern, got nil")
	}
}

// TestStatsTracking tests that statistics are tracked correctly
func TestStatsTracking(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with different states
	files := map[string]string{
		"new.go":      "package main\n",
		"existing.go": "/*\n * ########################################\n * Copyright (c) 2025\n * ########################################\n */\n\npackage main\n",
	}

	for filename, content := range files {
		fullPath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	licenseFile := filepath.Join(tmpDir, "LICENSE")
	if err := os.WriteFile(licenseFile, []byte("Copyright (c) 2025"), 0644); err != nil {
		t.Fatalf("Failed to write license file: %v", err)
	}

	cfg := &Config{
		LicenseText:       "Copyright (c) 2025",
		Input:             filepath.Join(tmpDir, "*.go"),
		Skip:              "",
		Prompt:            false,
		PresetStyle:       "hash",
		ForceCommentStyle: force.No,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify stats
	if processor.stats["added"] != 1 {
		t.Errorf("Expected 1 file added, got %d", processor.stats["added"])
	}
	if processor.stats["existing"] != 1 {
		t.Errorf("Expected 1 file with existing license, got %d", processor.stats["existing"])
	}
}
