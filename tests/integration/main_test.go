// tests/integration/main_test.go
package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// testCase defines a test configuration for a specific language
type testCase struct {
	language string
	patterns []string
}

var testCases = []testCase{
	{"batch", []string{"test_data/batch/*.bat"}},
	{"c", []string{"test_data/c/*.c", "test_data/c/*.h"}},
	{"cpp", []string{"test_data/cpp/*.cpp", "test_data/cpp/*.hpp"}},
	{"csharp", []string{"test_data/csharp/*.cs"}},
	{"css", []string{"test_data/css/*.css"}},
	{"go", []string{"test_data/go/*.go"}},
	{"html", []string{"test_data/html/*.html"}},
	{"java", []string{"test_data/java/*.java"}},
	{"javascript", []string{"test_data/javascript/*.js", "test_data/javascript/*.jsx"}},
	{"lua", []string{"test_data/lua/*.lua"}},
	{"perl", []string{"test_data/perl/*.pl", "test_data/perl/*.pm"}},
	{"php", []string{"test_data/php/*.php"}},
	{"python", []string{"test_data/python/*.py"}},
	{"r", []string{"test_data/r/*.r"}},
	{"ruby", []string{"test_data/ruby/*.rb"}},
	{"rust", []string{"test_data/rust/*.rs"}},
	{"sass", []string{"test_data/sass/*.sass"}},
	{"scss", []string{"test_data/scss/*.scss"}},
	{"shell", []string{"test_data/shell/*.sh", "test_data/shell/*.bash"}},
	{"swift", []string{"test_data/swift/*.swift"}},
	{"typescript", []string{"test_data/typescript/*.ts", "test_data/typescript/*.tsx"}},
	{"xml", []string{"test_data/xml/*.xml"}},
	{"yaml", []string{"test_data/yaml/*.yaml", "test_data/yaml/*.yml"}},
}

var projectRoot string

func init() {
	// Get project root directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectRoot = filepath.Dir(filepath.Dir(wd))
}

func TestSetup(t *testing.T) {
	_, _, err := runCommand(t, "build-test-data")
	if err != nil {
		t.Fatalf("Failed to build test data: %v", err)
	}
}

func runCommand(t *testing.T, args ...string) (string, string, error) {
	cmdArgs := append([]string{"run", "main.go"}, args...)
	//t.Logf("Running command: go %v", cmdArgs)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	//t.Logf("Command output:\nStdout: %s\nStderr: %s\nError: %v", stdout.String(), stderr.String(), err)
	return stdout.String(), stderr.String(), err
}

type testHelper struct {
	t            *testing.T
	projectRoot  string
	licenseFile  string
	mu           sync.RWMutex
	originalData map[string][]byte
}

func newTestHelper(t *testing.T) *testHelper {
	t.Helper()

	// Create a temporary license file
	tmpDir := t.TempDir()
	licenseFile := filepath.Join(tmpDir, "license.txt")
	licenseContent := "MIT License\nCopyright (c) 2024 Test Author\nPermission is hereby granted..."
	if err := os.WriteFile(licenseFile, []byte(licenseContent), 0644); err != nil {
		t.Fatal(err)
	}

	return &testHelper{
		t:            t,
		projectRoot:  projectRoot,
		licenseFile:  licenseFile,
		originalData: make(map[string][]byte),
	}
}

func (h *testHelper) getPattern(tc testCase) []string {
	patterns := make([]string, len(tc.patterns))
	for i, pattern := range tc.patterns {
		patterns[i] = filepath.Join(h.projectRoot, pattern)
	}
	return patterns
}

func (h *testHelper) getTemplateFile(testFile string) string {
	// Convert test_data path to templates path
	rel, err := filepath.Rel(filepath.Join(h.projectRoot, "test_data"), testFile)
	if err != nil {
		h.t.Fatal(err)
	}
	return filepath.Join(h.projectRoot, "templates", rel)
}

func (h *testHelper) verifyContentMatchesTemplate(patterns []string) {
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			h.t.Fatal(err)
		}

		for _, testFile := range matches {
			templateFile := h.getTemplateFile(testFile)
			testContent, err := os.ReadFile(testFile)
			if err != nil {
				h.t.Fatal(err)
			}

			templateContent, err := os.ReadFile(templateFile)
			if err != nil {
				h.t.Fatal(err)
			}

			if !bytes.Equal(normalizeContent(testContent), normalizeContent(templateContent)) {
				h.t.Errorf("Content mismatch between test file %s and template %s", testFile, templateFile)
			}
		}
	}
}

func (h *testHelper) runLicenseCommand(cmd string, patterns []string, expectError bool) (string, string, error) {
	args := []string{cmd}
	for _, pattern := range patterns {
		args = append(args, "--input", pattern)
	}
	if cmd != "remove" {
		args = append(args, "--license", h.licenseFile)
	}
	// args = append(args, "--verbose")

	stdout, stderr, err := runCommand(h.t, args...)

	if expectError && err == nil {
		h.t.Errorf("Expected %s command to fail, but it succeeded", cmd)
	} else if !expectError && err != nil {
		h.t.Errorf("%s command failed: %v\nStdout: %s\nStderr: %s", cmd, err, stdout, stderr)
	}

	return stdout, stderr, err
}

func (h *testHelper) verifyLicenseMissing(patterns []string) {
	// Run a license check command and expect it to fail
	for _, pattern := range patterns {
		_, stderr, err := h.runLicenseCommand("check", []string{pattern}, true)
		if err == nil {
			h.t.Error("Expected check to fail for files without license, but it passed")
		}
		if !strings.Contains(stderr, "missing") && !strings.Contains(stderr, "incorrect") {
			h.t.Errorf("Expected missing/incorrect license error, got: %s", stderr)
		}
	}
}

func (h *testHelper) verifyLicensePresent(patterns []string) {
	for _, pattern := range patterns {
		stdout, stderr, err := h.runLicenseCommand("check", []string{pattern}, false)
		if err != nil {
			h.t.Errorf("License check failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
		}
	}
}

func TestAddCheck(t *testing.T) {
	resetTestData(t)
	h := newTestHelper(t)

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.language, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel
			patterns := h.getPattern(tc)

			// Verify no license initially
			h.verifyLicenseMissing(patterns)

			// Add license
			h.runLicenseCommand("add", patterns, false)

			// Verify license present
			h.verifyLicensePresent(patterns)
		})
	}
}

func TestAddCheckFail(t *testing.T) {
	resetTestData(t)
	h := newTestHelper(t)

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.language, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel
			patterns := h.getPattern(tc)

			// Create a license file with incorrect content
			licenseFile := createLicenseFile(t, "incorrect license")
			defer os.Remove(licenseFile)

			// Add license
			h.runLicenseCommand("add", patterns, false)

			// Verify license present
			h.verifyLicensePresent(patterns)

			// Create a new license file with different content
			newLicenseFile := createLicenseFile(t, "new license")
			defer os.Remove(newLicenseFile)

			// Check should fail since license content doesn't match
			h.runLicenseCommand("check", patterns, true)
		})
	}
}

func TestAddUpdateCheck(t *testing.T) {
	resetTestData(t)
	h := newTestHelper(t)

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.language, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel
			patterns := h.getPattern(tc)

			// Create a license file
			licenseFile := createLicenseFile(t, "original license")
			defer os.Remove(licenseFile)

			// Add license
			h.runLicenseCommand("add", patterns, false)

			// Verify license present
			h.verifyLicensePresent(patterns)

			// Create a new license file with different content
			newLicenseFile := createLicenseFile(t, "new license")
			defer os.Remove(newLicenseFile)

			// Update license
			h.runLicenseCommand("update", patterns, false)

			// Check should pass with new license
			h.runLicenseCommand("check", patterns, false)
		})
	}
}

func TestAddRemove(t *testing.T) {
	resetTestData(t)
	h := newTestHelper(t)

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.language, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel
			patterns := h.getPattern(tc)

			// Add license
			h.runLicenseCommand("add", patterns, false)

			// Verify license present
			h.verifyLicensePresent(patterns)

			// Remove license
			h.runLicenseCommand("remove", patterns, false)

			// Verify content matches original
			h.verifyContentMatchesTemplate(patterns)

			// Verify license is missing
			h.verifyLicenseMissing(patterns)
		})
	}
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
	// Normalize content for comparison
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}
