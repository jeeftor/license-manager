# Testing Guide for License Manager

## Important: Unicode Markers in License Headers

### What Are Unicode Markers?

The license manager intentionally adds **zero-width unicode characters** around header and footer markers to aid in license detection and extraction.

**Defined in:** `internal/language/comment.go:17-18`

```go
const (
    MarkerStart = "​" // Zero-width space (U+200B)
    MarkerEnd   = "‌" // Zero-width non-joiner (U+200C)
)
```

### Why This Matters for Testing

When the system generates a license with the "hash" style, the actual output is:

```
/*
 *​######################################‌
 * Copyright (c) 2025
 *​######################################‌
 */
```

**NOT:**
```
/*
 * ########################################
 * Copyright (c) 2025
 * ########################################
 */
```

The unicode characters are **invisible** but present in the actual string.

### Testing Best Practices

#### ❌ DON'T: Manually Create License Blocks

```go
// This will NOT work - missing unicode markers
testContent := `/*
 * ########################################
 * Copyright (c) 2025
 * ########################################
 */

package main
`
```

#### ✅ DO: Use the System to Generate Licenses

```go
// Use the test helper
helper := NewTestHelper(t, "Copyright (c) 2025")
filePath := helper.CreateFile("test.go", "package main\n")
helper.AddLicenseToFile(filePath) // System generates proper format

// Now the file has a properly formatted license
content := helper.ReadFile(filePath)
```

#### ✅ DO: Use Partial String Matching

```go
// Check for partial content instead of exact format
if !strings.Contains(content, "Copyright (c) 2025") {
    t.Error("License text not found")
}

// Check for hash characters (unicode markers are invisible)
if !strings.Contains(content, "######") {
    t.Error("Hash markers not found")
}
```

## Test Helper Functions

### NewTestHelper

Creates a test environment with temporary directory and license file:

```go
helper := NewTestHelper(t, "Copyright (c) 2025 My Corp")
```

### CreateFile

Creates a test file with content:

```go
filePath := helper.CreateFile("main.go", "package main\n\nfunc main() {}\n")
```

### AddLicenseToFile

Adds a properly formatted license using the system:

```go
helper.AddLicenseToFile(filePath)
```

### CreateProcessor

Creates a processor with custom configuration:

```go
processor := helper.CreateProcessor(
    filepath.Join(helper.TmpDir(), "*.go"),
    force.Single, // Use single-line comments
)
```

## Complete Test Example

```go
func TestMyFeature(t *testing.T) {
    // Setup
    helper := NewTestHelper(t, "Copyright (c) 2025")
    filePath := helper.CreateFile("test.go", "package main\n")

    // Add license using system
    helper.AddLicenseToFile(filePath)

    // Verify
    content := helper.ReadFile(filePath)
    if !strings.Contains(content, "Copyright (c) 2025") {
        t.Error("License not added")
    }

    // Test update
    helper2 := NewTestHelper(t, "Copyright (c) 2026")
    processor := helper2.CreateProcessor(filePath, force.No)
    if err := processor.Update(); err != nil {
        t.Fatalf("Update failed: %v", err)
    }

    // Verify update
    content = helper.ReadFile(filePath)
    if !strings.Contains(content, "Copyright (c) 2026") {
        t.Error("License not updated")
    }
}
```

## Testing Different File Types

```go
func TestMultipleLanguages(t *testing.T) {
    helper := NewTestHelper(t, "Copyright (c) 2025")

    // Go file
    goFile := helper.CreateFile("main.go", "package main\n")
    helper.AddLicenseToFile(goFile)

    // Python file
    pyFile := helper.CreateFile("main.py", "def main():\n    pass\n")
    helper.AddLicenseToFile(pyFile)

    // JavaScript file
    jsFile := helper.CreateFile("main.js", "function main() {}\n")
    helper.AddLicenseToFile(jsFile)

    // All files now have properly formatted licenses
}
```

## Testing Comment Styles

```go
func TestCommentStyles(t *testing.T) {
    helper := NewTestHelper(t, "Copyright (c) 2025")
    filePath := helper.CreateFile("test.go", "package main\n")

    // Test single-line comments
    processor := helper.CreateProcessor(filePath, force.Single)
    if err := processor.Add(); err != nil {
        t.Fatalf("Add failed: %v", err)
    }

    content := helper.ReadFile(filePath)
    if !strings.Contains(content, "//") {
        t.Error("Expected single-line comments")
    }
}
```

## Common Pitfalls

### 1. Checking for Exact Hash Count

❌ **Don't do this:**
```go
if !strings.Contains(content, "########################################") {
    t.Error("Hash line not found")
}
```

✅ **Do this instead:**
```go
if !strings.Contains(content, "######") {
    t.Error("Hash markers not found")
}
```

### 2. Comparing Entire License Blocks

❌ **Don't do this:**
```go
expected := `/*
 * ########################################
 * Copyright
 * ########################################
 */`
if !strings.Contains(content, expected) {
    t.Error("License format mismatch")
}
```

✅ **Do this instead:**
```go
// Check for key components
if !strings.Contains(content, "Copyright") {
    t.Error("Copyright text not found")
}
if !strings.Contains(content, "/*") {
    t.Error("Multi-line comment not found")
}
```

### 3. Testing License Detection

❌ **Don't create licenses manually:**
```go
content := `/* ... manually created license ... */`
manager.SearchForLicense(content) // Will fail
```

✅ **Use system-generated licenses:**
```go
helper := NewTestHelper(t, "Copyright")
filePath := helper.CreateFile("test.go", "package main\n")
helper.AddLicenseToFile(filePath)
content := helper.ReadFile(filePath)
manager.SearchForLicense(content) // Will succeed
```

## Running Tests

```bash
# Run all tests
go test ./internal/...

# Run specific test
go test ./internal/processor -run TestMyFeature

# Run with verbose output
go test -v ./internal/processor

# Run with coverage
go test -cover ./internal/...
```

## Pre-existing Test Issues

Some tests may fail due to pre-existing issues unrelated to your changes:

- **Go handler directive scanning**: Edge cases with blank lines
- **Python handler extraction**: Pre-existing extraction issues

These are known issues and not related to the unicode marker system.
