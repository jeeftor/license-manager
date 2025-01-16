package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readMITLicense(t *testing.T) string {
	t.Helper()
	// Get the absolute path to the project root
	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir("."))))
	mitPath := filepath.Join(projectRoot, "..", "..", "templates", "licenses", "mit.txt")

	content, err := os.ReadFile(mitPath)
	if err != nil {
		t.Fatalf("Failed to read MIT license template: %v", err)
	}
	return strings.TrimSpace(string(content))
}

func TestCommentString(t *testing.T) {
	mitLicense := readMITLicense(t)
	tests := []struct {
		name     string
		comment  *Comment
		expected string
	}{
		{
			name: "Go style multi-line comment",
			comment: &Comment{
				Style:  getCommentStyle("test.go"),
				Header: "Copyright (c) 2025 Test User",
				Body:   mitLicense,
				Footer: "End of MIT License",
			},
			expected: "/*\n * Copyright (c) 2025 Test User\n *\n * MIT License\n *\n * " +
				strings.ReplaceAll(mitLicense, "\n", "\n * ") +
				"\n *\n * End of MIT License\n */",
		},
		{
			name: "Python style single-line comment",
			comment: &Comment{
				Style:  getCommentStyle("test.py"),
				Header: "Copyright (c) 2025 Test User",
				Body:   mitLicense,
				Footer: "End of MIT License",
			},
			expected: "# Copyright (c) 2025 Test User\n#\n# MIT License\n#\n# " +
				strings.ReplaceAll(mitLicense, "\n", "\n# ") +
				"\n#\n# End of MIT License",
		},
		{
			name: "HTML style comment",
			comment: &Comment{
				Style:  getCommentStyle("test.html"),
				Header: "Copyright (c) 2025 Test User",
				Body:   mitLicense,
				Footer: "End of MIT License",
			},
			expected: "<!--\nCopyright (c) 2025 Test User\n\nMIT License\n\n" +
				mitLicense +
				"\n\nEnd of MIT License\n-->",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.String()
			// Normalize newlines for comparison
			result = strings.ReplaceAll(result, "\r\n", "\n")
			expected := strings.ReplaceAll(tt.expected, "\r\n", "\n")
			if result != expected {
				t.Errorf("Comment.String() got:\n%v\nwant:\n%v", result, expected)
			}
		})
	}
}

func TestUncommentContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		style    CommentStyle
		expected string
	}{
		{
			name:     "Uncomment Go style",
			content:  "/*\n * Test header\n * Test body\n * Test footer\n */",
			style:    getCommentStyle("test.go"),
			expected: "Test header\nTest body\nTest footer",
		},
		{
			name:     "Uncomment Python style",
			content:  "# Test header\n# Test body\n# Test footer",
			style:    getCommentStyle("test.py"),
			expected: "Test header\nTest body\nTest footer",
		},
		{
			name:     "Uncomment HTML style",
			content:  "<!-- Test header\nTest body\nTest footer -->",
			style:    getCommentStyle("test.html"),
			expected: "Test header\nTest body\nTest footer",
		},
		{
			name:     "Uncomment with empty lines",
			content:  "/*\n * Test header\n *\n * Test body\n *\n * Test footer\n */",
			style:    getCommentStyle("test.go"),
			expected: "Test header\n\nTest body\n\nTest footer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uncommentContent(tt.content, tt.style)
			// Normalize newlines and whitespace for comparison
			result = strings.TrimSpace(strings.ReplaceAll(result, "\r\n", "\n"))
			expected := strings.TrimSpace(strings.ReplaceAll(tt.expected, "\r\n", "\n"))
			if result != expected {
				t.Errorf("uncommentContent() got:\n%v\nwant:\n%v", result, expected)
			}
		})
	}
}

func TestAddMarkersIfNeeded(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Add markers to text without markers",
			text:     "Test text",
			expected: markerStart + "Test text" + markerEnd,
		},
		{
			name:     "Text already has markers",
			text:     markerStart + "Test text" + markerEnd,
			expected: markerStart + "Test text" + markerEnd,
		},
		{
			name:     "Empty text",
			text:     "",
			expected: markerStart + markerEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addMarkersIfNeeded(tt.text)
			if result != tt.expected {
				t.Errorf("addMarkersIfNeeded() got:\n%v\nwant:\n%v", result, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		style    CommentStyle
		wantOk   bool
		expected *Comment
	}{
		{
			name: "Parse Go multi-line comment",
			content: "/*\n * " + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n *\n * " +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n * ") +
				"\n *\n * " + markerStart + "End of MIT License" + markerEnd + "\n */",
			style:  getCommentStyle("test.go"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.go"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse Go single-line comment",
			content: "// " + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n//\n// " +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n// ") +
				"\n//\n// " + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.go"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.go"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse Go single-line comment without space",
			content: "//" + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n//\n//" +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n//") +
				"\n//\n//" + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.go"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.go"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse Python style comment",
			content: "# " + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n#\n# " +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n# ") +
				"\n#\n# " + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.py"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.py"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse C++ style comment",
			content: "// " + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n//\n// " +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n// ") +
				"\n//\n// " + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.cpp"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.cpp"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse Python style comment with no space after #",
			content: "#" + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n#\n#" +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n#") +
				"\n#\n#" + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.py"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.py"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse C++ style comment with no space after //",
			content: "//" + markerStart + "Copyright (c) 2025 Test User" + markerEnd + "\n//\n//" +
				strings.ReplaceAll(readMITLicense(t), "\n", "\n//") +
				"\n//\n//" + markerStart + "End of MIT License" + markerEnd,
			style:  getCommentStyle("test.cpp"),
			wantOk: true,
			expected: &Comment{
				Style:  getCommentStyle("test.cpp"),
				Header: markerStart + "Copyright (c) 2025 Test User" + markerEnd,
				Body:   readMITLicense(t),
				Footer: markerStart + "End of MIT License" + markerEnd,
			},
		},
		{
			name: "Parse invalid comment",
			content: "This is not a valid comment\n" +
				readMITLicense(t) +
				"\nNo proper structure",
			style:  getCommentStyle("test.go"),
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, ok := Parse(tt.content, tt.style)
			if ok != tt.wantOk {
				t.Errorf("Parse() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if !tt.wantOk {
				return
			}

			// Normalize newlines and whitespace for comparison
			if strings.TrimSpace(comment.Header) != strings.TrimSpace(tt.expected.Header) {
				t.Errorf("Parse() header got:\n%v\nwant:\n%v", comment.Header, tt.expected.Header)
			}
			if strings.TrimSpace(comment.Body) != strings.TrimSpace(tt.expected.Body) {
				t.Errorf("Parse() body got:\n%v\nwant:\n%v", comment.Body, tt.expected.Body)
			}
			if strings.TrimSpace(comment.Footer) != strings.TrimSpace(tt.expected.Footer) {
				t.Errorf("Parse() footer got:\n%v\nwant:\n%v", comment.Footer, tt.expected.Footer)
			}
		})
	}
}
