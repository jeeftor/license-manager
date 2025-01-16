package processor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Define shared test data
type testFile struct {
	name         string
	templatePath string
	commentStyle CommentStyle
}

var testFiles = []testFile{
	{
		name:         "Go - Multi-line Comments",
		templatePath: "../../templates/go/hello.go",
		commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "go"},
	},
	{
		name:         "Python - Single-line Comments",
		templatePath: "../../templates/python/hello.py",
		commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "python"},
	},
	{
		name:         "JavaScript - Multi-line Comments",
		templatePath: "../../templates/javascript/hello.js",
		commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
	},
	{
		name:         "JSX - Multi-line Comments",
		templatePath: "../../templates/javascript/component.jsx",
		commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "javascript"},
	},
	{
		name:         "HTML - Multi-line Comments",
		templatePath: "../../templates/html/index.html",
		commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "html"},
	},
	{
		name:         "CSS - Multi-line Comments",
		templatePath: "../../templates/css/style.css",
		commentStyle: CommentStyle{Single: "", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "css"},
	},
	{
		name:         "Ruby - Single-line Comments",
		templatePath: "../../templates/ruby/hello.rb",
		commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "ruby"},
	},
	{
		name:         "Shell - Single-line Comments",
		templatePath: "../../templates/shell/hello.sh",
		commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
	},
	{
		name:         "C++ - Multi-line Comments",
		templatePath: "../../templates/cpp/hello.cpp",
		commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
	},
	{
		name:         "Java - Multi-line Comments",
		templatePath: "../../templates/java/HelloWorld.java",
		commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "java"},
	},
	{
		name:         "YAML - Single-line Comments",
		templatePath: "../../templates/yaml/config.yaml",
		commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
	},
	{
		name:         "Lua - Both Comment Styles",
		templatePath: "../../templates/lua/hello.lua",
		commentStyle: CommentStyle{Single: "--", MultiStart: "--[[", MultiEnd: "]]", PreferMulti: true, FileType: "lua"},
	},
}

// TestAddLicenseOnce tests adding a license once to each file type
func TestAddLicenseOnce(t *testing.T) {
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

	for _, tc := range testFiles {
		t.Run(tc.name, func(t *testing.T) {
			// Read template file
			templateContent, err := os.ReadFile(tc.templatePath)
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", tc.templatePath, err)
			}

			// Create a temporary file
			tempFile := filepath.Join(tempDir, filepath.Base(tc.templatePath))
			err = os.WriteFile(tempFile, templateContent, 0644)
			if err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(HeaderFooterStyle{}, string(licenseText), tc.commentStyle)

			// Read the file content
			content, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			// Add license to the content
			contentWithLicense := lm.AddLicense(string(content))

			// Write back to file
			err = os.WriteFile(tempFile, []byte(contentWithLicense), 0644)
			if err != nil {
				t.Fatalf("Failed to write modified file: %v", err)
			}

			// Check if the license was added correctly
			if !lm.CheckLicense(contentWithLicense) {
				t.Error("License text not found in modified file")
			}

			// Check if the original content is preserved
			if !bytes.Contains([]byte(contentWithLicense), templateContent) {
				t.Error("Original content not preserved in modified file")
			}
		})
	}
}

// TestAddLicenseTwice tests that adding a license twice doesn't create duplicates
func TestAddLicenseTwice(t *testing.T) {
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

	for _, tc := range testFiles {
		t.Run(tc.name, func(t *testing.T) {
			// Read template file
			templateContent, err := os.ReadFile(tc.templatePath)
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", tc.templatePath, err)
			}

			// Create a temporary file
			tempFile := filepath.Join(tempDir, filepath.Base(tc.templatePath))
			err = os.WriteFile(tempFile, templateContent, 0644)
			if err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(HeaderFooterStyle{}, string(licenseText), tc.commentStyle)

			// Read the file content
			content, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			// Add license first time
			contentWithOneLicense := lm.AddLicense(string(content))

			// Write back to file
			err = os.WriteFile(tempFile, []byte(contentWithOneLicense), 0644)
			if err != nil {
				t.Fatalf("Failed to write modified file: %v", err)
			}

			// Add license second time
			contentWithTwoLicenses := lm.AddLicense(contentWithOneLicense)

			// The content should be identical after adding the license twice
			if contentWithOneLicense != contentWithTwoLicenses {
				t.Error("File content changed after adding license twice")
			}
		})
	}
}

// TestRemoveLicenseMultipleStyles tests removing licenses from files with different comment styles
func TestRemoveLicenseMultipleStyles(t *testing.T) {
	// Read the MIT license
	mitLicenseBytes, err := os.ReadFile("../../templates/licenses/mit.txt")
	if err != nil {
		t.Fatalf("Failed to read license file: %v", err)
	}
	mitLicense := string(mitLicenseBytes)
	mitLicense = strings.ReplaceAll(mitLicense, "[year]", "2025")
	mitLicense = strings.ReplaceAll(mitLicense, "[fullname]", "Test User")

	for _, tt := range testFiles {
		t.Run(tt.name, func(t *testing.T) {
			// Read the original content
			originalContent, err := os.ReadFile(tt.templatePath)
			if err != nil {
				t.Fatalf("Failed to read template file: %v", err)
			}
			originalContentStr := string(originalContent)

			// Create a temporary directory
			tempDir, err := os.MkdirTemp("", "license-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create a temporary file with the template content
			tempFile := filepath.Join(tempDir, filepath.Base(tt.templatePath))
			if err := os.WriteFile(tempFile, originalContent, 0644); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(HeaderFooterStyle{}, mitLicense, tt.commentStyle)

			// Read the file content
			content, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			// Add the license to the content
			contentWithLicense := lm.AddLicense(string(content))

			// Write back to file
			err = os.WriteFile(tempFile, []byte(contentWithLicense), 0644)
			if err != nil {
				t.Fatalf("Failed to write modified file: %v", err)
			}

			// Remove the license
			contentAfterRemoval := lm.RemoveLicense(contentWithLicense)

			// Write back to file
			err = os.WriteFile(tempFile, []byte(contentAfterRemoval), 0644)
			if err != nil {
				t.Fatalf("Failed to write modified file: %v", err)
			}

			// Compare the content after removal with the original content
			if contentAfterRemoval != originalContentStr {
				t.Errorf("Content after license removal does not match original content.\nExpected:\n%s\nGot:\n%s",
					originalContentStr, contentAfterRemoval)
			}
		})
	}
}

func normalizeNewlines(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}
