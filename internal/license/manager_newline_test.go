package license

import (
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

func TestPreserveNewlinesBeforeEOF(t *testing.T) {
	tests := []struct {
		name           string
		inputContent   string
		licenseText    string
		expectedEnding string
		description    string
	}{
		{
			name: "preserve single newline before EOF",
			inputContent: `package main

func main() {
	println("hello")
}
`,
			licenseText: "Copyright (c) 2025",
			expectedEnding: `func main() {
	println("hello")
}
`,
			description: "File with single trailing newline should preserve it",
		},
		{
			name: "preserve double newline before EOF",
			inputContent: `package main

func main() {
	println("hello")
}


`,
			licenseText: "Copyright (c) 2025",
			expectedEnding: `func main() {
	println("hello")
}


`,
			description: "File with double trailing newlines should preserve them",
		},
		{
			name: "preserve triple newline before EOF",
			inputContent: `package main

func main() {
	println("hello")
}



`,
			licenseText: "Copyright (c) 2025",
			expectedEnding: `func main() {
	println("hello")
}



`,
			description: "File with triple trailing newlines should preserve them",
		},
		{
			name: "no newline before EOF",
			inputContent: `package main

func main() {
	println("hello")
}`,
			licenseText: "Copyright (c) 2025",
			expectedEnding: `func main() {
	println("hello")
}`,
			description: "File with no trailing newline should preserve that",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewLogger(logger.ErrorLevel)
			commentStyle := styles.GetLanguageCommentStyle(".go")
			style := styles.Get("hash")

			manager := NewLicenseManager(log, tt.licenseText, ".go", style, commentStyle)
			manager.SetFileContent(tt.inputContent)

			hasLicense, components := manager.HasLicense(tt.inputContent)
			if hasLicense {
				t.Fatalf("Expected no license in input content")
			}

			newContent, err := manager.AddLicense(components, "go")
			if err != nil {
				t.Fatalf("AddLicense failed: %v", err)
			}

			if !strings.HasSuffix(newContent, tt.expectedEnding) {
				t.Errorf("%s\nExpected content to end with:\n%q\nBut got:\n%q\n\nFull content:\n%s",
					tt.description,
					tt.expectedEnding,
					getLastNChars(newContent, len(tt.expectedEnding)+50),
					newContent)
			}

			// Also verify the exact number of trailing newlines
			expectedTrailingNewlines := countTrailingNewlines(tt.inputContent)
			actualTrailingNewlines := countTrailingNewlines(newContent)
			if expectedTrailingNewlines != actualTrailingNewlines {
				t.Errorf(
					"Expected %d trailing newlines, got %d\nInput ends with: %q\nOutput ends with: %q",
					expectedTrailingNewlines,
					actualTrailingNewlines,
					getLastNChars(tt.inputContent, 20),
					getLastNChars(newContent, 20),
				)
			}
		})
	}
}

func TestUpdateLicensePreservesNewlines(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	commentStyle := styles.GetLanguageCommentStyle(".go")
	style := styles.Get("hash")

	// Start with a file that has no license and a trailing newline
	originalContent := "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n\n"

	// Step 1: Add a license using the system so it has proper unicode markers
	oldLicense := "Copyright (c) 2024"
	addManager := NewLicenseManager(log, oldLicense, ".go", style, commentStyle)
	addManager.SetFileContent(originalContent)
	_, addComponents := addManager.HasLicense(originalContent)
	contentWithLicense, err := addManager.AddLicense(addComponents, "go")
	if err != nil {
		t.Fatalf("AddLicense failed: %v", err)
	}

	// Step 2: Now update the license
	newLicense := "Copyright (c) 2025"
	updateManager := NewLicenseManager(log, newLicense, ".go", style, commentStyle)
	updateManager.SetFileContent(contentWithLicense)

	hasLicense, components := updateManager.HasLicense(contentWithLicense)
	if !hasLicense {
		t.Fatalf("Expected to find license in content after add")
	}

	updatedContent, err := updateManager.UpdateLicense(components, "go")
	if err != nil {
		t.Fatalf("UpdateLicense failed: %v", err)
	}

	// Verify trailing newlines are preserved
	expectedTrailingNewlines := countTrailingNewlines(contentWithLicense)
	actualTrailingNewlines := countTrailingNewlines(updatedContent)
	if expectedTrailingNewlines != actualTrailingNewlines {
		t.Errorf("Expected %d trailing newlines, got %d",
			expectedTrailingNewlines,
			actualTrailingNewlines)
	}

	// Verify the new license text is present
	if !strings.Contains(updatedContent, "Copyright (c) 2025") {
		t.Errorf("Expected updated content to contain new license text")
	}
	if strings.Contains(updatedContent, "Copyright (c) 2024") {
		t.Errorf("Expected updated content to NOT contain old license text")
	}
}

func TestRemoveLicensePreservesNewlines(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	commentStyle := styles.GetLanguageCommentStyle(".go")
	style := styles.Get("hash")

	licenseText := "Copyright (c) 2025"

	// Start with a file that has no license and two trailing newlines
	originalContent := "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n\n"

	// Step 1: Add a license using the system so it has proper unicode markers
	addManager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
	addManager.SetFileContent(originalContent)
	_, addComponents := addManager.HasLicense(originalContent)
	contentWithLicense, err := addManager.AddLicense(addComponents, "go")
	if err != nil {
		t.Fatalf("AddLicense failed: %v", err)
	}

	// Step 2: Now remove the license
	removeManager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
	removeManager.SetFileContent(contentWithLicense)

	hasLicense, components := removeManager.HasLicense(contentWithLicense)
	if !hasLicense {
		t.Fatalf("Expected to find license in content after add")
	}

	removedContent, err := removeManager.RemoveLicense(components, "go")
	if err != nil {
		t.Fatalf("RemoveLicense failed: %v", err)
	}

	// Verify the code is preserved
	if !strings.Contains(removedContent, "package main") {
		t.Errorf("Expected removed content to contain package declaration")
	}
	if !strings.Contains(removedContent, "println(\"hello\")") {
		t.Errorf("Expected removed content to contain function body")
	}

	// Verify the license is gone
	if strings.Contains(removedContent, "Copyright") {
		t.Errorf("Expected removed content to NOT contain license text")
	}

	// Verify trailing newlines are preserved
	expectedTrailingNewlines := countTrailingNewlines(contentWithLicense)
	actualTrailingNewlines := countTrailingNewlines(removedContent)
	if expectedTrailingNewlines != actualTrailingNewlines {
		t.Errorf("Expected %d trailing newlines, got %d",
			expectedTrailingNewlines,
			actualTrailingNewlines)
	}
}

func countTrailingNewlines(s string) int {
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			count++
		} else {
			break
		}
	}
	return count
}

func getLastNChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
