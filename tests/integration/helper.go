// helper.go
package integration

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testHelper struct {
	t           *testing.T
	projectRoot string
	licenseFile string
}

func newTestHelper(t *testing.T) *testHelper {
	t.Helper()
	tmpDir := t.TempDir()
	licenseFile := filepath.Join(tmpDir, "license.txt")
	licenseContent := "MIT License\nCopyright (c) 2024 Test Author\nPermission is hereby granted..."

	if err := os.WriteFile(licenseFile, []byte(licenseContent), 0644); err != nil {
		t.Fatal(err)
	}

	return &testHelper{
		t:           t,
		projectRoot: projectRoot,
		licenseFile: licenseFile,
	}
}

func (h *testHelper) getPattern(tc languageData) []string {
	patterns := make([]string, len(tc.patterns))
	for i, pattern := range tc.patterns {
		patterns[i] = filepath.Join(h.projectRoot, pattern)
	}
	return patterns
}

func (h *testHelper) getTemplateFile(testFile string) string {
	rel, err := filepath.Rel(filepath.Join(h.projectRoot, "test_data"), testFile)
	if err != nil {
		h.t.Fatal(err)
	}
	return filepath.Join(h.projectRoot, "templates", rel)
}

func (h *testHelper) verifyContentMatchesTemplate(patterns []string) error {
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to glob pattern %s: %v", pattern, err)
		}

		for _, testFile := range matches {
			templateFile := h.getTemplateFile(testFile)
			testContent, err := os.ReadFile(testFile)
			if err != nil {
				return fmt.Errorf("failed to read test file %s: %v", testFile, err)
			}

			templateContent, err := os.ReadFile(templateFile)
			if err != nil {
				return fmt.Errorf("failed to read template file %s: %v", templateFile, err)
			}

			if !bytes.Equal(normalizeContent(testContent), normalizeContent(templateContent)) {
				return fmt.Errorf("content mismatch between test file %s and template %s", testFile, templateFile)
			}
		}
	}
	return nil
}

func (h *testHelper) runLicenseCommand(cmd string, patterns []string) (string, string, error) {
	args := []string{cmd}
	for _, pattern := range patterns {
		args = append(args, "--input", pattern)
	}
	if cmd != "remove" {
		args = append(args, "--license", h.licenseFile)
	}
	stdout, stderr, err := runCommand(h.t, args...)
	return stdout, stderr, err
}

func (h *testHelper) verifyLicenseMissing(patterns []string) error {
	for _, pattern := range patterns {
		_, stderr, err := h.runLicenseCommand("check", []string{pattern})

		if strings.Contains(stderr, "exit status 2") {
			// This is the case we want.
			return fmt.Errorf("at least one non-matching license was found")
		}
		if err == nil {
			return fmt.Errorf("expected check to fail for files without license, but it passed")
		}
	}
	return nil
}

func (h *testHelper) verifyLicensePresent(patterns []string) error {
	for _, pattern := range patterns {
		stdout, stderr, err := h.runLicenseCommand("check", []string{pattern})
		if err != nil {
			return fmt.Errorf("license check failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
		}
	}
	return nil
}

func resetTestData(t *testing.T) {
	testDataDir := filepath.Join(projectRoot, "test_data")
	if err := os.RemoveAll(testDataDir); err != nil {
		t.Fatal(err)
	}

	_, _, err := runCommand(t, "build-test-data")
	if err != nil {
		t.Fatal(err)
	}
}

func createLicenseFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	licenseFile := filepath.Join(tmpDir, "license.txt")
	if err := os.WriteFile(licenseFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return licenseFile
}

func normalizeContent(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}
