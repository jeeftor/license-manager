package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLicenseAddRemoveCycle(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "license-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test files with different comment styles
	testCases := []struct {
		templatePath string
		filename     string
		commentStyle CommentStyle
	}{
		// C family
		{
			templatePath: "templates/c/hello.c",
			filename:     "test.c",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			templatePath: "templates/c/hello.h",
			filename:     "test.h",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			templatePath: "templates/cpp/hello.cpp",
			filename:     "test.cpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},
		{
			templatePath: "templates/cpp/hello.hpp",
			filename:     "test.hpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},
		// C#
		{
			templatePath: "templates/csharp/Hello.cs",
			filename:     "test.cs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "csharp"},
		},
		// Web
		{
			templatePath: "templates/css/style.css",
			filename:     "test.css",
			commentStyle: CommentStyle{Single: "", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "css"},
		},
		{
			templatePath: "templates/html/index.html",
			filename:     "test.html",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "html"},
		},
		{
			templatePath: "templates/xml/hello.xml",
			filename:     "test.xml",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "xml"},
		},
		// JavaScript family
		{
			templatePath: "templates/javascript/hello.js",
			filename:     "test.js",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
		},
		{
			templatePath: "templates/javascript/component.jsx",
			filename:     "test.jsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "javascript"},
		},
		{
			templatePath: "templates/typescript/hello.ts",
			filename:     "test.ts",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "typescript"},
		},
		{
			templatePath: "templates/typescript/component.tsx",
			filename:     "test.tsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "typescript"},
		},
		// Go
		{
			templatePath: "templates/go/hello.go",
			filename:     "test.go",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: false, FileType: "go"},
		},
		// Java
		{
			templatePath: "templates/java/HelloWorld.java",
			filename:     "test.java",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "java"},
		},
		// Scripting languages
		{
			templatePath: "templates/python/hello.py",
			filename:     "test.py",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "python"},
		},
		{
			templatePath: "templates/ruby/hello.rb",
			filename:     "test.rb",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "ruby"},
		},
		{
			templatePath: "templates/perl/hello.pl",
			filename:     "test.pl",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			templatePath: "templates/perl/Hello.pm",
			filename:     "test.pm",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			templatePath: "templates/php/hello.php",
			filename:     "test.php",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "php"},
		},
		{
			templatePath: "templates/lua/hello.lua",
			filename:     "test.lua",
			commentStyle: CommentStyle{Single: "--", MultiStart: "--[[", MultiEnd: "]]", PreferMulti: false, FileType: "lua"},
		},
		{
			templatePath: "templates/r/hello.r",
			filename:     "test.r",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "r"},
		},
		// Shell scripts
		{
			templatePath: "templates/shell/hello.sh",
			filename:     "test.sh",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
		},
		{
			templatePath: "templates/shell/hello.bash",
			filename:     "test.bash",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
		},
		// Other languages
		{
			templatePath: "templates/rust/hello.rs",
			filename:     "test.rs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: false, FileType: "rust"},
		},
		{
			templatePath: "templates/swift/hello.swift",
			filename:     "test.swift",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "swift"},
		},
		// Style sheets
		{
			templatePath: "templates/sass/style.sass",
			filename:     "test.sass",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "sass"},
		},
		{
			templatePath: "templates/scss/style.scss",
			filename:     "test.scss",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "scss"},
		},
		// Config files
		{
			templatePath: "templates/yaml/config.yml",
			filename:     "test.yml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
		{
			templatePath: "templates/yaml/config.yaml",
			filename:     "test.yaml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
	}

	// Read license text from LICENSE file
	licenseText, err := os.ReadFile("../../LICENSE")
	if err != nil {
		t.Fatalf("Failed to read LICENSE file: %v", err)
	}

	headerFooterStyle := HeaderFooterStyle{
		Header: "----------------------------------------",
		Footer: "----------------------------------------",
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			// Read template file
			templateContent, err := os.ReadFile(filepath.Join("../../", tc.templatePath))
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", tc.templatePath, err)
			}

			// Create test file
			testFile := filepath.Join(tempDir, tc.filename)
			originalContent := string(templateContent)
			err = os.WriteFile(testFile, templateContent, 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(headerFooterStyle, string(licenseText), tc.commentStyle)

			// Step 1: Add license for the first time
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			contentWithLicense := lm.AddLicense(string(content))
			err = os.WriteFile(testFile, []byte(contentWithLicense), 0644)
			if err != nil {
				t.Fatalf("Failed to write file with license: %v", err)
			}

			// Verify license exists
			if !lm.CheckLicense(contentWithLicense) {
				t.Error("License not found after first addition")
			}

			// Step 2: Try to add license again
			content, err = os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			contentWithSecondLicense := lm.AddLicense(string(content))

			// Verify content hasn't changed
			if contentWithSecondLicense != contentWithLicense {
				t.Error("Content changed after second license addition")
				t.Logf("First addition:\n%s", contentWithLicense)
				t.Logf("Second addition:\n%s", contentWithSecondLicense)
			}

			// Step 3: Remove license
			contentWithoutLicense := lm.RemoveLicense(contentWithSecondLicense)

			// Step 4: Verify content matches original
			if normalizeNewlines(contentWithoutLicense) != normalizeNewlines(originalContent) {
				t.Error("Content after removal doesn't match original")
				t.Logf("Original:\n%s", originalContent)
				t.Logf("After removal:\n%s", contentWithoutLicense)
			}
		})
	}
}

// TestAddLicenseOnce tests adding a license once to each file type
func TestAddLicenseOnce(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "license-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test files with different comment styles
	testCases := []struct {
		name         string
		templatePath string
		filename     string
		commentStyle CommentStyle
	}{
		// C family
		{
			name:         "C",
			templatePath: "templates/c/hello.c",
			filename:     "test.c",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			name:         "C Header",
			templatePath: "templates/c/hello.h",
			filename:     "test.h",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			name:         "C++",
			templatePath: "templates/cpp/hello.cpp",
			filename:     "test.cpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},
		{
			name:         "C++ Header",
			templatePath: "templates/cpp/hello.hpp",
			filename:     "test.hpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},

		// Web languages
		{
			name:         "JavaScript",
			templatePath: "templates/javascript/hello.js",
			filename:     "test.js",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
		},
		{
			name:         "JSX",
			templatePath: "templates/javascript/component.jsx",
			filename:     "test.jsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "javascript"},
		},
		{
			name:         "TypeScript",
			templatePath: "templates/typescript/hello.ts",
			filename:     "test.ts",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "typescript"},
		},
		{
			name:         "TSX",
			templatePath: "templates/typescript/component.tsx",
			filename:     "test.tsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "typescript"},
		},
		{
			name:         "CSS",
			templatePath: "templates/css/style.css",
			filename:     "test.css",
			commentStyle: CommentStyle{Single: "", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "css"},
		},
		{
			name:         "SCSS",
			templatePath: "templates/scss/style.scss",
			filename:     "test.scss",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "scss"},
		},
		{
			name:         "SASS",
			templatePath: "templates/sass/style.sass",
			filename:     "test.sass",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "sass"},
		},
		{
			name:         "HTML",
			templatePath: "templates/html/index.html",
			filename:     "test.html",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "html"},
		},
		{
			name:         "XML",
			templatePath: "templates/xml/hello.xml",
			filename:     "test.xml",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "xml"},
		},

		// Scripting languages
		{
			name:         "Python",
			templatePath: "templates/python/hello.py",
			filename:     "test.py",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "python"},
		},
		{
			name:         "Ruby",
			templatePath: "templates/ruby/hello.rb",
			filename:     "test.rb",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "ruby"},
		},
		{
			name:         "Perl",
			templatePath: "templates/perl/hello.pl",
			filename:     "test.pl",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			name:         "Perl Module",
			templatePath: "templates/perl/Hello.pm",
			filename:     "test.pm",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			name:         "PHP",
			templatePath: "templates/php/hello.php",
			filename:     "test.php",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "php"},
		},
		{
			name:         "Lua",
			templatePath: "templates/lua/hello.lua",
			filename:     "test.lua",
			commentStyle: CommentStyle{Single: "--", MultiStart: "--[[", MultiEnd: "]]", PreferMulti: false, FileType: "lua"},
		},
		{
			name:         "R",
			templatePath: "templates/r/hello.r",
			filename:     "test.r",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "r"},
		},

		// Shell scripts
		{
			name:         "Shell",
			templatePath: "templates/shell/hello.sh",
			filename:     "test.sh",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
		},
		{
			name:         "Bash",
			templatePath: "templates/shell/hello.bash",
			filename:     "test.bash",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
		},

		// Compiled languages
		{
			name:         "Go",
			templatePath: "templates/go/hello.go",
			filename:     "test.go",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: false, FileType: "go"},
		},
		{
			name:         "Java",
			templatePath: "templates/java/HelloWorld.java",
			filename:     "test.java",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "java"},
		},
		{
			name:         "C#",
			templatePath: "templates/csharp/Hello.cs",
			filename:     "test.cs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "csharp"},
		},
		{
			name:         "Rust",
			templatePath: "templates/rust/hello.rs",
			filename:     "test.rs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "rust"},
		},
		{
			name:         "Swift",
			templatePath: "templates/swift/hello.swift",
			filename:     "test.swift",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "swift"},
		},

		// Config files
		{
			name:         "YAML",
			templatePath: "templates/yaml/config.yaml",
			filename:     "test.yaml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
		{
			name:         "YML",
			templatePath: "templates/yaml/config.yml",
			filename:     "test.yml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
	}

	// Read license text from LICENSE file
	licenseText, err := os.ReadFile("../../LICENSE")
	if err != nil {
		t.Fatalf("Failed to read LICENSE file: %v", err)
	}

	headerFooterStyle := HeaderFooterStyle{
		Header: "----------------------------------------",
		Footer: "----------------------------------------",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read template file
			templateContent, err := os.ReadFile(filepath.Join("../../", tc.templatePath))
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", tc.templatePath, err)
			}

			// Create test file
			testFile := filepath.Join(tempDir, tc.filename)
			err = os.WriteFile(testFile, templateContent, 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(headerFooterStyle, string(licenseText), tc.commentStyle)

			// Step 1: Add license
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			// Log original content
			t.Logf("Original content:\n%s", content)

			// Add license
			contentWithLicense := lm.AddLicense(string(content))

			// Log content with license
			t.Logf("Content with license:\n%s", contentWithLicense)

			// Write back to file
			err = os.WriteFile(testFile, []byte(contentWithLicense), 0644)
			if err != nil {
				t.Fatalf("Failed to write file with license: %v", err)
			}

			// Verify license exists
			if !lm.CheckLicense(contentWithLicense) {
				t.Error("License not found after addition")
				t.Log("Expected formatted license:")
				t.Log(lm.formatLicenseBlock(string(licenseText)))
			}
		})
	}
}

// TestAddLicenseTwice tests that adding a license twice doesn't create duplicates
func TestAddLicenseTwice(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "license-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test files with different comment styles
	testCases := []struct {
		name         string
		templatePath string
		filename     string
		commentStyle CommentStyle
	}{
		// C family
		{
			name:         "C",
			templatePath: "templates/c/hello.c",
			filename:     "test.c",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			name:         "C Header",
			templatePath: "templates/c/hello.h",
			filename:     "test.h",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
		},
		{
			name:         "C++",
			templatePath: "templates/cpp/hello.cpp",
			filename:     "test.cpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},
		{
			name:         "C++ Header",
			templatePath: "templates/cpp/hello.hpp",
			filename:     "test.hpp",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
		},

		// Web languages
		{
			name:         "JavaScript",
			templatePath: "templates/javascript/hello.js",
			filename:     "test.js",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
		},
		{
			name:         "JSX",
			templatePath: "templates/javascript/component.jsx",
			filename:     "test.jsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "javascript"},
		},
		{
			name:         "TypeScript",
			templatePath: "templates/typescript/hello.ts",
			filename:     "test.ts",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "typescript"},
		},
		{
			name:         "TSX",
			templatePath: "templates/typescript/component.tsx",
			filename:     "test.tsx",
			commentStyle: CommentStyle{Single: "//", MultiStart: "{/*", MultiEnd: "*/}", PreferMulti: true, FileType: "typescript"},
		},
		{
			name:         "CSS",
			templatePath: "templates/css/style.css",
			filename:     "test.css",
			commentStyle: CommentStyle{Single: "", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "css"},
		},
		{
			name:         "SCSS",
			templatePath: "templates/scss/style.scss",
			filename:     "test.scss",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "scss"},
		},
		{
			name:         "SASS",
			templatePath: "templates/sass/style.sass",
			filename:     "test.sass",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "sass"},
		},
		{
			name:         "HTML",
			templatePath: "templates/html/index.html",
			filename:     "test.html",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "html"},
		},
		{
			name:         "XML",
			templatePath: "templates/xml/hello.xml",
			filename:     "test.xml",
			commentStyle: CommentStyle{Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "xml"},
		},

		// Scripting languages
		{
			name:         "Python",
			templatePath: "templates/python/hello.py",
			filename:     "test.py",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "python"},
		},
		{
			name:         "Ruby",
			templatePath: "templates/ruby/hello.rb",
			filename:     "test.rb",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "ruby"},
		},
		{
			name:         "Perl",
			templatePath: "templates/perl/hello.pl",
			filename:     "test.pl",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			name:         "Perl Module",
			templatePath: "templates/perl/Hello.pm",
			filename:     "test.pm",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "perl"},
		},
		{
			name:         "PHP",
			templatePath: "templates/php/hello.php",
			filename:     "test.php",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "php"},
		},
		{
			name:         "Lua",
			templatePath: "templates/lua/hello.lua",
			filename:     "test.lua",
			commentStyle: CommentStyle{Single: "--", MultiStart: "--[[", MultiEnd: "]]", PreferMulti: true, FileType: "lua"},
		},
		{
			name:         "R",
			templatePath: "templates/r/hello.r",
			filename:     "test.r",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "r"},
		},
		{
			name:         "Shell",
			templatePath: "templates/shell/hello.sh",
			filename:     "test.sh",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "shell"},
		},
		{
			name:         "Batch",
			templatePath: "templates/batch/hello.bat",
			filename:     "test.bat",
			commentStyle: CommentStyle{Single: "REM", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "batch"},
		},

		// Systems languages
		{
			name:         "Go",
			templatePath: "templates/go/hello.go",
			filename:     "test.go",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: false, FileType: "go"},
		},
		{
			name:         "Java",
			templatePath: "templates/java/Hello.java",
			filename:     "test.java",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "java"},
		},
		{
			name:         "C#",
			templatePath: "templates/csharp/Hello.cs",
			filename:     "test.cs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "csharp"},
		},
		{
			name:         "Rust",
			templatePath: "templates/rust/hello.rs",
			filename:     "test.rs",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "rust"},
		},
		{
			name:         "Swift",
			templatePath: "templates/swift/hello.swift",
			filename:     "test.swift",
			commentStyle: CommentStyle{Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "swift"},
		},

		// Config languages
		{
			name:         "YAML",
			templatePath: "templates/yaml/config.yaml",
			filename:     "test.yaml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
		{
			name:         "YML",
			templatePath: "templates/yaml/config.yml",
			filename:     "test.yml",
			commentStyle: CommentStyle{Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
		},
	}

	// Read license text from LICENSE file
	licenseText, err := os.ReadFile("../../LICENSE")
	if err != nil {
		t.Fatalf("Failed to read LICENSE file: %v", err)
	}

	headerFooterStyle := HeaderFooterStyle{
		Header: "----------------------------------------",
		Footer: "----------------------------------------",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read template file
			templateContent, err := os.ReadFile(filepath.Join("../../", tc.templatePath))
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", tc.templatePath, err)
			}

			// Create test file
			testFile := filepath.Join(tempDir, tc.filename)
			err = os.WriteFile(testFile, templateContent, 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create license manager
			lm := NewLicenseManager(headerFooterStyle, string(licenseText), tc.commentStyle)

			// Step 1: Add license first time
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			// Log original content
			t.Logf("Original content:\n%s", content)

			// Add license first time
			contentWithOneLicense := lm.AddLicense(string(content))

			// Log content after first addition
			t.Logf("Content after first license addition:\n%s", contentWithOneLicense)

			// Add license second time
			contentWithTwoLicenses := lm.AddLicense(contentWithOneLicense)

			// Log content after second addition
			t.Logf("Content after second license addition:\n%s", contentWithTwoLicenses)

			// Verify content hasn't changed after second addition
			if contentWithOneLicense != contentWithTwoLicenses {
				t.Error("Content changed after adding license twice")
				t.Log("Expected (after first addition):")
				t.Log(contentWithOneLicense)
				t.Log("Got (after second addition):")
				t.Log(contentWithTwoLicenses)
				t.Log("Difference in lengths:", len(contentWithTwoLicenses)-len(contentWithOneLicense))
			}

			// Count occurrences of license text
			formattedLicense := lm.formatLicenseBlock(string(licenseText))
			count := strings.Count(contentWithTwoLicenses, formattedLicense)
			if count > 1 {
				t.Errorf("Found %d licenses in the file, expected 1", count)
				t.Log("License text being searched for:")
				t.Log(formattedLicense)
				t.Log("Full file content:")
				t.Log(contentWithTwoLicenses)
			} else if count == 0 {
				t.Error("No license found in file")
				t.Log("License text being searched for:")
				t.Log(formattedLicense)
				t.Log("Full file content:")
				t.Log(contentWithTwoLicenses)
			}
		})
	}
}

// normalizeNewlines replaces all newlines with \n for consistent comparison
func normalizeNewlines(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}
