package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
)

func TestForceCommentStyleSingle(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	initialContent := `package main

func main() {
	println("hello")
}
`

	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
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
		ForceCommentStyle: force.Single,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "//") {
		t.Errorf(
			"Expected single-line comments (//), but file doesn't contain them.\nContent:\n%s",
			contentStr,
		)
	}

	if strings.Contains(contentStr, "/*") || strings.Contains(contentStr, "*/") {
		t.Errorf(
			"Expected single-line comments (//), but found multi-line comment markers (/* */).\nContent:\n%s",
			contentStr,
		)
	}

	lines := strings.Split(contentStr, "\n")
	foundCommentLine := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") && strings.Contains(trimmed, "Copyright") {
			foundCommentLine = true
			break
		}
	}

	if !foundCommentLine {
		t.Errorf(
			"Expected to find single-line comment with 'Copyright', but didn't find it.\nContent:\n%s",
			contentStr,
		)
	}
}

func TestForceCommentStyleMulti(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	initialContent := `package main

func main() {
	println("hello")
}
`

	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
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
		ForceCommentStyle: force.Multi,
		LogLevel:          logger.ErrorLevel,
	}

	processor := NewFileProcessor(cfg)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "/*") || !strings.Contains(contentStr, "*/") {
		t.Errorf(
			"Expected multi-line comments (/* */), but didn't find them.\nContent:\n%s",
			contentStr,
		)
	}

	if strings.Contains(contentStr, "// ########################################") {
		t.Errorf(
			"Expected multi-line comments (/* */), but found single-line comment style.\nContent:\n%s",
			contentStr,
		)
	}
}

func TestForceCommentStyleNo(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	initialContent := `package main

func main() {
	println("hello")
}
`

	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
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

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "/*") || !strings.Contains(contentStr, "*/") {
		t.Errorf(
			"Expected default multi-line comments for Go (/* */), but didn't find them.\nContent:\n%s",
			contentStr,
		)
	}
}
