package processor

import (
	"bytes"
	"fmt"
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
	/* Temporarily disabled Ruby and Shell support
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
	*/
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

			// Create license manager with preset header/footer style
			style := GetPresetStyle("brackets")
			lm := NewLicenseManager(style, string(licenseText), tc.commentStyle)

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

			// Check if the license was added correctly with proper header/footer
			if !strings.Contains(contentWithLicense, style.Header) {
				t.Errorf("Header marker not found in modified file.\nExpected: %q", style.Header)
			}
			if !strings.Contains(contentWithLicense, style.Footer) {
				t.Errorf("Footer marker not found in modified file.\nExpected: %q", style.Footer)
			}

			// Also verify the license text itself is present
			if !lm.CheckLicense(contentWithLicense, false) {
				t.Error("License text not found in modified file")
			}

			// Check if the original content is preserved (normalize newlines first)
			normalizedContent := normalizeNewlines(contentWithLicense)
			normalizedTemplate := normalizeNewlines(string(templateContent))
			if !strings.Contains(normalizedContent, normalizedTemplate) {
				t.Errorf("Original content not preserved in modified file.\nExpected to find:\n%s\nIn:\n%s", normalizedTemplate, normalizedContent)
			}
		})
	}
}

// TestAddLicenseTwice tests that adding a license twice doesn't create duplicates
func TestAddLicenseTwice(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		licenseText string
		style       HeaderFooterStyle
		commentStyle CommentStyle
	}{
		{
			name: "Go file",
			content: `package main

func main() {
	println("Hello, World!")
}`,
			licenseText: "Copyright 2025",
			style: HeaderFooterStyle{
				Header: "⚡️ LICENSE-START ⚡️",
				Footer: "⚡️ LICENSE-END ⚡️",
			},
			commentStyle: CommentStyle{
				Single:      "//",
				MultiStart:  "/*",
				MultiEnd:    "*/",
				PreferMulti: true,
				FileType:    "go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := NewLicenseManager(tt.style, tt.licenseText, tt.commentStyle)

			// Add license first time
			contentWithLicense := lm.AddLicense(tt.content)

			// Verify license was added
			if !lm.CheckLicense(contentWithLicense, false) {
				t.Error("License not found after first addition")
			}

			// Add license second time
			contentWithLicenseTwice := lm.AddLicense(contentWithLicense)

			// Verify content hasn't changed
			if contentWithLicenseTwice != contentWithLicense {
				t.Error("Content changed after adding license twice")
			}
		})
	}
}

// TestAddLicenseMatrix tests adding a license with each header style to each file type
func TestAddLicenseMatrix(t *testing.T) {
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

	// Get all preset styles
	styles := ListPresetStyles()

	for _, tc := range testFiles {
		for _, styleName := range styles {
			t.Run(fmt.Sprintf("%s/%s", tc.name, styleName), func(t *testing.T) {
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

				// Create license manager with preset header/footer style
				style := GetPresetStyle(styleName)
				lm := NewLicenseManager(style, string(licenseText), tc.commentStyle)

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

				// Check if the license was added correctly with proper header/footer
				if !strings.Contains(contentWithLicense, style.Header) {
					t.Errorf("Header marker not found in modified file.\nExpected: %q", style.Header)
				}
				if !strings.Contains(contentWithLicense, style.Footer) {
					t.Errorf("Footer marker not found in modified file.\nExpected: %q", style.Footer)
				}

				// Also verify the license text itself is present
				if !lm.CheckLicense(contentWithLicense, false) {
					t.Error("License text not found in modified file")
				}

				// Check if the original content is preserved (normalize newlines first)
				normalizedContent := normalizeNewlines(contentWithLicense)
				normalizedTemplate := normalizeNewlines(string(templateContent))
				if !strings.Contains(normalizedContent, normalizedTemplate) {
					t.Errorf("Original content not preserved in modified file.\nExpected to find:\n%s\nIn:\n%s", normalizedTemplate, normalizedContent)
				}

				// Try to add the license again and verify it doesn't change
				contentWithSecondLicense := lm.AddLicense(contentWithLicense)
				if contentWithSecondLicense != contentWithLicense {
					t.Error("Content changed after adding license twice")
				}
			})
		}
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

			// Create license manager with header/footer style
			var header, footer string
			if tt.commentStyle.Single != "" {
				header = tt.commentStyle.Single + " ⚡️ LICENSE-START ⚡️\n"
				footer = tt.commentStyle.Single + " ⚡️ LICENSE-END ⚡️"
			} else {
				header = tt.commentStyle.MultiStart + " ⚡️ LICENSE-START ⚡️ " + tt.commentStyle.MultiEnd + "\n"
				footer = tt.commentStyle.MultiStart + " ⚡️ LICENSE-END ⚡️ " + tt.commentStyle.MultiEnd
			}
			lm := NewLicenseManager(HeaderFooterStyle{
				Header: header,
				Footer: footer,
			}, mitLicense, tt.commentStyle)

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
				t.Errorf("Content after license removal does not match original content.\nExpected:\n%s\nGot:\n%s\nExpected bytes: %v\nGot bytes: %v",
					originalContentStr, contentAfterRemoval, []byte(originalContentStr), []byte(contentAfterRemoval))
			}
		})
	}
}
