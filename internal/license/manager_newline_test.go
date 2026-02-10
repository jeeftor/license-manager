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

	newLicense := "Copyright (c) 2025"

	inputContent := `/*
 * ########################################
 * Copyright (c) 2024
 * ########################################
 */

package main

func main() {
	println("hello")
}

`

	expectedEnding := `func main() {
	println("hello")
}

`

	manager := NewLicenseManager(log, newLicense, ".go", style, commentStyle)
	manager.SetFileContent(inputContent)

	hasLicense, components := manager.HasLicense(inputContent)
	if !hasLicense {
		t.Fatalf("Expected to find license in input content")
	}

	newContent, err := manager.UpdateLicense(components, "go")
	if err != nil {
		t.Fatalf("UpdateLicense failed: %v", err)
	}

	if !strings.HasSuffix(newContent, expectedEnding) {
		t.Errorf("Expected content to end with:\n%q\nBut got:\n%q",
			expectedEnding,
			getLastNChars(newContent, len(expectedEnding)+50))
	}

	expectedTrailingNewlines := countTrailingNewlines(inputContent)
	actualTrailingNewlines := countTrailingNewlines(newContent)
	if expectedTrailingNewlines != actualTrailingNewlines {
		t.Errorf("Expected %d trailing newlines, got %d",
			expectedTrailingNewlines,
			actualTrailingNewlines)
	}
}

func TestRemoveLicensePreservesNewlines(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	commentStyle := styles.GetLanguageCommentStyle(".go")
	style := styles.Get("hash")

	licenseText := "Copyright (c) 2025"

	inputContent := `/*
 * ########################################
 * Copyright (c) 2025
 * ########################################
 */

package main

func main() {
	println("hello")
}


`

	expectedContent := `package main

func main() {
	println("hello")
}


`

	manager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
	manager.SetFileContent(inputContent)

	hasLicense, components := manager.HasLicense(inputContent)
	if !hasLicense {
		t.Fatalf("Expected to find license in input content")
	}

	newContent, err := manager.RemoveLicense(components, "go")
	if err != nil {
		t.Fatalf("RemoveLicense failed: %v", err)
	}

	if newContent != expectedContent {
		t.Errorf("Expected content:\n%q\nGot:\n%q",
			expectedContent,
			newContent)
	}

	expectedTrailingNewlines := countTrailingNewlines(inputContent)
	actualTrailingNewlines := countTrailingNewlines(newContent)
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
