package license

import (
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

// TestCheckLicenseStatus tests the license status checking functionality
func TestCheckLicenseStatus(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	tests := []struct {
		name           string
		licenseText    string
		fileContent    string
		expectedStatus Status
		description    string
	}{
		{
			name:        "full match",
			licenseText: "Copyright (c) 2025 Test Corp",
			fileContent: `/*
 * ########################################
 * Copyright (c) 2025 Test Corp
 * ########################################
 */

package main

func main() {}
`,
			expectedStatus: FullMatch,
			description:    "License content and style both match",
		},
		{
			name:        "no license",
			licenseText: "Copyright (c) 2025 Test Corp",
			fileContent: `package main

func main() {}
`,
			expectedStatus: NoLicense,
			description:    "No license header present",
		},
		{
			name:        "content mismatch",
			licenseText: "Copyright (c) 2025 Test Corp",
			fileContent: `/*
 * ########################################
 * Copyright (c) 2024 Old Corp
 * ########################################
 */

package main

func main() {}
`,
			expectedStatus: ContentMismatch,
			description:    "License exists but content differs",
		},
		{
			name:        "style mismatch - single line vs multi line",
			licenseText: "Copyright (c) 2025 Test Corp",
			fileContent: `// ########################################
// Copyright (c) 2025 Test Corp
// ########################################

package main

func main() {}
`,
			expectedStatus: StyleMismatch,
			description:    "License content matches but style differs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewLicenseManager(log, tt.licenseText, ".go", style, commentStyle)
			manager.SetFileContent(tt.fileContent)

			// First search for the license
			_ = manager.SearchForLicense(tt.fileContent)

			status := manager.CheckLicenseStatus(tt.fileContent)

			if status != tt.expectedStatus {
				t.Errorf("%s\nExpected status: %v (%s)\nGot status: %v (%s)",
					tt.description,
					tt.expectedStatus,
					tt.expectedStatus.String(),
					status,
					status.String())
			}
		})
	}
}

// TestUpdateLicense tests the license update functionality
func TestUpdateLicense(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	commentStyle := styles.GetLanguageCommentStyle(".go")
	style := styles.Get("hash")

	tests := []struct {
		name             string
		oldLicense       string
		newLicense       string
		fileContent      string
		expectError      bool
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:       "update license content",
			oldLicense: "Copyright (c) 2024 Old Corp",
			newLicense: "Copyright (c) 2025 New Corp",
			fileContent: `/*
 * ########################################
 * Copyright (c) 2024 Old Corp
 * ########################################
 */

package main

func main() {}
`,
			expectError:      false,
			shouldContain:    []string{"Copyright (c) 2025 New Corp", "package main"},
			shouldNotContain: []string{"Copyright (c) 2024 Old Corp"},
		},
		{
			name:       "update multi-line license",
			oldLicense: "Copyright (c) 2024\nOld Corp",
			newLicense: "Copyright (c) 2025\nNew Corp\nAll rights reserved",
			fileContent: `/*
 * ########################################
 * Copyright (c) 2024
 * Old Corp
 * ########################################
 */

package main
`,
			expectError:      false,
			shouldContain:    []string{"Copyright (c) 2025", "New Corp", "All rights reserved"},
			shouldNotContain: []string{"Old Corp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First create manager with old license to detect it
			oldManager := NewLicenseManager(log, tt.oldLicense, ".go", style, commentStyle)
			oldManager.SetFileContent(tt.fileContent)
			_ = oldManager.SearchForLicense(tt.fileContent)

			// Now create manager with new license to update
			manager := NewLicenseManager(log, tt.newLicense, ".go", style, commentStyle)
			manager.SetFileContent(tt.fileContent)
			_ = manager.SearchForLicense(tt.fileContent)

			newContent, err := manager.UpdateLicense(manager.InitialComponents, "go")

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				for _, expected := range tt.shouldContain {
					if !strings.Contains(newContent, expected) {
						t.Errorf("Expected content to contain %q, but it doesn't.\nContent:\n%s",
							expected, newContent)
					}
				}

				for _, notExpected := range tt.shouldNotContain {
					if strings.Contains(newContent, notExpected) {
						t.Errorf("Expected content NOT to contain %q, but it does.\nContent:\n%s",
							notExpected, newContent)
					}
				}
			}
		})
	}
}

// TestSearchForLicense tests the license search functionality
func TestSearchForLicense(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	licenseText := "Copyright (c) 2025 Test Corp"
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	tests := []struct {
		name        string
		fileContent string
		expectFound bool
		expectStyle bool
		description string
	}{
		{
			name: "finds multi-line license",
			fileContent: `/*
 * ########################################
 * Copyright (c) 2025 Test Corp
 * ########################################
 */

package main
`,
			expectFound: true,
			expectStyle: true,
			description: "Should find multi-line license with style markers",
		},
		{
			name: "finds single-line license",
			fileContent: `// ########################################
// Copyright (c) 2025 Test Corp
// ########################################

package main
`,
			expectFound: true,
			expectStyle: true,
			description: "Should find single-line license with style markers",
		},
		{
			name: "no license present",
			fileContent: `package main

func main() {}
`,
			expectFound: false,
			expectStyle: false,
			description: "Should not find license when none exists",
		},
		{
			name: "license without style markers",
			fileContent: `/*
 * Copyright (c) 2025 Test Corp
 */

package main
`,
			expectFound: true,
			expectStyle: false,
			description: "Should find license but not match style",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)

			results := manager.SearchForLicense(tt.fileContent)

			if results.HasLicense != tt.expectFound {
				t.Errorf("%s\nExpected HasLicense=%v, got %v",
					tt.description, tt.expectFound, results.HasLicense)
			}

			if results.IsStyleMatch != tt.expectStyle {
				t.Errorf("%s\nExpected IsStyleMatch=%v, got %v",
					tt.description, tt.expectStyle, results.IsStyleMatch)
			}
		})
	}
}

// TestCGoHeaderHandling tests special handling of C-Go style headers
func TestCGoHeaderHandling(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	licenseText := "Copyright (c) 2025 Test Corp"
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	cgoContent := `#include <stdio.h>
#include <stdlib.h>

package main

func main() {}
`

	manager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
	manager.SetFileContent(cgoContent)

	results := manager.SearchForLicense(cgoContent)

	if results.HasLicense {
		t.Errorf("C-Go headers should not be detected as licenses")
	}

	// Try to add license - should return original content unchanged
	newContent, err := manager.AddLicense(manager.InitialComponents, "go")
	if err != nil {
		t.Errorf("AddLicense should not error on C-Go files: %v", err)
	}

	if newContent != cgoContent {
		t.Errorf("C-Go files should not be modified")
	}
}

// TestEmptyFileHandling tests handling of empty or whitespace-only files
func TestEmptyFileHandling(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	licenseText := "Copyright (c) 2025 Test Corp"
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	tests := []struct {
		name        string
		fileContent string
		description string
	}{
		{
			name:        "completely empty",
			fileContent: "",
			description: "Empty file should be handled gracefully",
		},
		{
			name:        "only whitespace",
			fileContent: "   \n\n\t\n   ",
			description: "Whitespace-only file should be handled gracefully",
		},
		{
			name:        "only newlines",
			fileContent: "\n\n\n",
			description: "Newline-only file should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
			manager.SetFileContent(tt.fileContent)

			results := manager.SearchForLicense(tt.fileContent)

			if results.HasLicense {
				t.Errorf("%s: Should not find license in empty/whitespace file", tt.description)
			}

			// Should be able to add license without error
			_, err := manager.AddLicense(manager.InitialComponents, "go")
			if err != nil {
				t.Errorf("%s: AddLicense failed: %v", tt.description, err)
			}
		})
	}
}

// TestStatusString tests the Status.String() method
func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{FullMatch, "License OK"},
		{NoLicense, "No license found"},
		{ContentMismatch, "License content mismatch"},
		{StyleMismatch, "License style mismatch"},
		{ContentAndStyleMismatch, "License content and style mismatch"},
		{Status(999), "Unknown status"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Status(%d).String() = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

// TestFormatLicenseForFile tests the public API for formatting licenses
func TestFormatLicenseForFile(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	tests := []struct {
		name          string
		licenseText   string
		shouldContain []string
	}{
		{
			name:        "single line license",
			licenseText: "Copyright (c) 2025 Test Corp",
			shouldContain: []string{
				"########################################",
				"Copyright (c) 2025 Test Corp",
			},
		},
		{
			name: "multi-line license",
			licenseText: `Copyright (c) 2025 Test Corp
All rights reserved.

Licensed under MIT License`,
			shouldContain: []string{
				"########################################",
				"Copyright (c) 2025 Test Corp",
				"All rights reserved.",
				"Licensed under MIT License",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewLicenseManager(log, tt.licenseText, ".go", style, commentStyle)

			formatted := manager.FormatLicenseForFile(tt.licenseText)

			for _, expected := range tt.shouldContain {
				if !strings.Contains(formatted, expected) {
					t.Errorf("Expected formatted license to contain %q\nGot:\n%s",
						expected, formatted)
				}
			}
		})
	}
}

// TestBuildDirectivesPreservation tests that Go build directives are preserved
func TestBuildDirectivesPreservation(t *testing.T) {
	log := logger.NewLogger(logger.ErrorLevel)
	licenseText := "Copyright (c) 2025 Test Corp"
	style := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".go")

	fileWithDirectives := `//go:build linux
// +build linux

package main

func main() {}
`

	manager := NewLicenseManager(log, licenseText, ".go", style, commentStyle)
	manager.SetFileContent(fileWithDirectives)

	_ = manager.SearchForLicense(fileWithDirectives)

	newContent, err := manager.AddLicense(manager.InitialComponents, "go")
	if err != nil {
		t.Fatalf("AddLicense failed: %v", err)
	}

	// Verify directives are preserved and come before license
	lines := strings.Split(newContent, "\n")
	foundDirective := false
	foundLicense := false
	directiveBeforeLicense := false

	for i, line := range lines {
		if strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build") {
			foundDirective = true
			if !foundLicense {
				directiveBeforeLicense = true
			}
		}
		if strings.Contains(line, "Copyright") {
			foundLicense = true
			if i > 0 && !foundDirective {
				t.Errorf("License should come after build directives")
			}
		}
	}

	if !foundDirective {
		t.Errorf("Build directives were not preserved")
	}
	if !foundLicense {
		t.Errorf("License was not added")
	}
	if !directiveBeforeLicense {
		t.Errorf("Build directives should come before license")
	}
}
