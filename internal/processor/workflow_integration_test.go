package processor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
)

// TestFullWorkflow tests the complete workflow: add → check → update → remove
func TestFullWorkflow(t *testing.T) {
	helper := NewTestHelper(t, "Copyright (c) 2025 Original Corp")

	// Create test files
	file1 := helper.CreateFile("main.go", "package main\n\nfunc main() {}\n")
	file2 := helper.CreateFile("utils/helper.go", "package utils\n\nfunc Helper() {}\n")

	pattern := filepath.Join(helper.TmpDir(), "**/*.go")

	// Step 1: Add licenses
	t.Run("Add", func(t *testing.T) {
		processor := helper.CreateProcessor(pattern, force.No)
		if err := processor.Add(); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Verify both files have licenses
		content1 := helper.ReadFile(file1)
		content2 := helper.ReadFile(file2)

		if !strings.Contains(content1, "Copyright (c) 2025 Original Corp") {
			t.Error("File 1 missing license")
		}
		if !strings.Contains(content2, "Copyright (c) 2025 Original Corp") {
			t.Error("File 2 missing license")
		}

		// Verify stats
		if processor.stats["added"] != 2 {
			t.Errorf("Expected 2 files added, got %d", processor.stats["added"])
		}
	})

	// Step 2: Check licenses (should pass)
	t.Run("Check_Pass", func(t *testing.T) {
		processor := helper.CreateProcessor(pattern, force.No)

		if err := processor.Check(); err != nil {
			t.Fatalf("Check failed: %v", err)
		}

		// Verify stats - the key is "passed" per file_processor.go:417
		if processor.stats["passed"] != 2 {
			t.Errorf(
				"Expected 2 files passed, got %d (stats: %+v)",
				processor.stats["passed"],
				processor.stats,
			)
		}
	})

	// Step 3: Update licenses
	t.Run("Update", func(t *testing.T) {
		// Create new helper with updated license
		helper2 := NewTestHelper(t, "Copyright (c) 2026 Updated Corp")
		processor := helper2.CreateProcessor(pattern, force.No)

		if err := processor.Update(); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Verify both files have updated licenses
		content1 := helper.ReadFile(file1)
		content2 := helper.ReadFile(file2)

		if !strings.Contains(content1, "Copyright (c) 2026 Updated Corp") {
			t.Error("File 1 not updated")
		}
		if !strings.Contains(content2, "Copyright (c) 2026 Updated Corp") {
			t.Error("File 2 not updated")
		}
		if strings.Contains(content1, "Copyright (c) 2025 Original Corp") {
			t.Error("File 1 still has old license")
		}

		// Verify stats
		if processor.stats["updated"] != 2 {
			t.Errorf("Expected 2 files updated, got %d", processor.stats["updated"])
		}
	})

	// Step 4: Check with wrong license (should report mismatch)
	t.Run("Check_Fail", func(t *testing.T) {
		// Check with original license text (should mismatch since we updated)
		processor := helper.CreateProcessor(pattern, force.No)

		// Check() returns an error when licenses don't match - that's expected
		err := processor.Check()
		if err == nil {
			t.Error("Expected Check to return error for mismatched licenses")
		}

		// Should have failed files
		if processor.stats["failed"] < 2 {
			t.Errorf(
				"Expected at least 2 failed files, got %d (stats: %+v)",
				processor.stats["failed"],
				processor.stats,
			)
		}
	})

	// Step 5: Remove licenses
	t.Run("Remove", func(t *testing.T) {
		// Use the updated license for removal
		helper2 := NewTestHelper(t, "Copyright (c) 2026 Updated Corp")
		processor := helper2.CreateProcessor(pattern, force.No)

		if err := processor.Remove(); err != nil {
			t.Fatalf("Remove failed: %v", err)
		}

		// Verify licenses are removed
		content1 := helper.ReadFile(file1)
		content2 := helper.ReadFile(file2)

		if strings.Contains(content1, "Copyright") {
			t.Error("File 1 still has license")
		}
		if strings.Contains(content2, "Copyright") {
			t.Error("File 2 still has license")
		}

		// Verify code is preserved
		if !strings.Contains(content1, "package main") {
			t.Error("File 1 code not preserved")
		}
		if !strings.Contains(content2, "package utils") {
			t.Error("File 2 code not preserved")
		}

		// Verify stats
		if processor.stats["removed"] != 2 {
			t.Errorf("Expected 2 files removed, got %d", processor.stats["removed"])
		}
	})

	// Step 6: Check after removal (should report missing)
	t.Run("Check_Missing", func(t *testing.T) {
		helper2 := NewTestHelper(t, "Copyright (c) 2026 Updated Corp")
		processor := helper2.CreateProcessor(pattern, force.No)

		// Check() returns an error when licenses are missing - that's expected
		err := processor.Check()
		if err == nil {
			t.Error("Expected Check to return error for missing licenses")
		}

		// Should report missing licenses
		if processor.stats["missing"] != 2 {
			t.Errorf(
				"Expected 2 files missing licenses, got %d (stats: %+v)",
				processor.stats["missing"],
				processor.stats,
			)
		}
	})
}

// TestWorkflowWithCommentStyles tests workflow with different comment styles
func TestWorkflowWithCommentStyles(t *testing.T) {
	helper := NewTestHelper(t, "Copyright (c) 2025")
	file := helper.CreateFile("test.go", "package main\n\nfunc main() {}\n")

	// Add with multi-line comments (default)
	t.Run("Add_MultiLine", func(t *testing.T) {
		processor := helper.CreateProcessor(file, force.Multi)
		if err := processor.Add(); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		content := helper.ReadFile(file)
		if !strings.Contains(content, "/*") {
			t.Error("Expected multi-line comments")
		}
	})

	// Remove
	t.Run("Remove", func(t *testing.T) {
		processor := helper.CreateProcessor(file, force.Multi)
		if err := processor.Remove(); err != nil {
			t.Fatalf("Remove failed: %v", err)
		}
	})

	// Add with single-line comments
	t.Run("Add_SingleLine", func(t *testing.T) {
		processor := helper.CreateProcessor(file, force.Single)
		if err := processor.Add(); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		content := helper.ReadFile(file)
		if !strings.Contains(content, "//") {
			t.Error("Expected single-line comments")
		}
		if strings.Contains(content, "/*") {
			t.Error("Should not have multi-line comments")
		}
	})
}

// TestWorkflowWithSkipPatterns tests workflow with skip patterns
func TestWorkflowWithSkipPatterns(t *testing.T) {
	helper := NewTestHelper(t, "Copyright (c) 2025")

	// Create files
	include := helper.CreateFile("src/main.go", "package main\n")
	skip := helper.CreateFile("vendor/lib.go", "package vendor\n")

	// Add licenses with skip pattern
	t.Run("Add_WithSkip", func(t *testing.T) {
		cfg := &Config{
			LicenseText:       helper.LicenseText(),
			Input:             filepath.Join(helper.TmpDir(), "**/*.go"),
			Skip:              filepath.Join(helper.TmpDir(), "vendor/**"),
			Prompt:            false,
			PresetStyle:       "hash",
			ForceCommentStyle: force.No,
			LogLevel:          logger.ErrorLevel,
		}
		processor := NewFileProcessor(cfg)

		if err := processor.Add(); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Verify included file has license
		includeContent := helper.ReadFile(include)
		if !strings.Contains(includeContent, "Copyright") {
			t.Error("Included file should have license")
		}

		// Verify skipped file doesn't have license
		skipContent := helper.ReadFile(skip)
		if strings.Contains(skipContent, "Copyright") {
			t.Error("Skipped file should not have license")
		}

		// Verify stats
		if processor.stats["added"] != 1 {
			t.Errorf("Expected 1 file added, got %d", processor.stats["added"])
		}
	})
}

// TestWorkflowIdempotency tests that operations are idempotent
func TestWorkflowIdempotency(t *testing.T) {
	helper := NewTestHelper(t, "Copyright (c) 2025")
	file := helper.CreateFile("test.go", "package main\n")

	// Add license
	processor := helper.CreateProcessor(file, force.No)
	if err := processor.Add(); err != nil {
		t.Fatalf("First add failed: %v", err)
	}

	content1 := helper.ReadFile(file)

	// Try to add again - should be skipped
	processor = helper.CreateProcessor(file, force.No)
	if err := processor.Add(); err != nil {
		t.Fatalf("Second add failed: %v", err)
	}

	content2 := helper.ReadFile(file)

	// Content should be unchanged
	if content1 != content2 {
		t.Error("File should not be modified on second add")
	}

	// Stats should show existing
	if processor.stats["existing"] != 1 {
		t.Errorf("Expected 1 existing license, got %d", processor.stats["existing"])
	}
	if processor.stats["added"] != 0 {
		t.Errorf("Expected 0 added, got %d", processor.stats["added"])
	}
}

// TestWorkflowPreservesContent tests that operations preserve file content
func TestWorkflowPreservesContent(t *testing.T) {
	helper := NewTestHelper(t, "Copyright (c) 2025")

	originalContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	file := helper.CreateFile("test.go", originalContent)

	// Add license
	processor := helper.CreateProcessor(file, force.No)
	if err := processor.Add(); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	contentWithLicense := helper.ReadFile(file)

	// Verify original content is preserved
	if !strings.Contains(contentWithLicense, "package main") {
		t.Error("Package declaration not preserved")
	}
	if !strings.Contains(contentWithLicense, `fmt.Println("Hello, World!")`) {
		t.Error("Function content not preserved")
	}

	// Remove license
	processor = helper.CreateProcessor(file, force.No)
	if err := processor.Remove(); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	contentAfterRemove := helper.ReadFile(file)

	// Content should match original (minus potential whitespace differences)
	if !strings.Contains(contentAfterRemove, "package main") {
		t.Error("Package declaration not preserved after removal")
	}
	if !strings.Contains(contentAfterRemove, `fmt.Println("Hello, World!")`) {
		t.Error("Function content not preserved after removal")
	}
}
