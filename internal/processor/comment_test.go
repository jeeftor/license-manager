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
			comment: func() *Comment {
				c := &Comment{
					Style: getCommentStyle("test.go"),
					Body:  mitLicense,
				}
				c.SetHeaderAndFooterStyle("minimal")
				return c
			}(),
			expected: "/*\n * ─────────────────────────────────────\n *\n * " +
				strings.ReplaceAll(mitLicense, "\n", "\n * ") +
				"\n *\n * ─────────────────────────────────────\n */",
		},
		{
			name: "Python style single-line comment",
			comment: func() *Comment {
				c := &Comment{
					Style: getCommentStyle("test.py"),
					Body:  mitLicense,
				}
				c.SetHeaderAndFooterStyle("minimal")
				return c
			}(),
			expected: "# ─────────────────────────────────────\n#\n# " +
				strings.ReplaceAll(mitLicense, "\n", "\n# ") +
				"\n#\n# ─────────────────────────────────────",
		},
		{
			name: "HTML style comment",
			comment: func() *Comment {
				c := &Comment{
					Style: getCommentStyle("test.html"),
					Body:  mitLicense,
				}
				c.SetHeaderAndFooterStyle("minimal")
				return c
			}(),
			expected: "<!--\n─────────────────────────────────────\n\n" +
				mitLicense +
				"\n\n─────────────────────────────────────\n-->",
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

func TestParseLicenseBlock(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantHeader  string
		wantBody    string
		wantFooter  string
		wantSuccess bool
	}{
		{
			name: "Parse Go multi-line comment",
			content: `/*
 * ─────────────────────────────────────
 *
 * MIT License
 *
 * Copyright (c) [year] [fullname]
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 * ─────────────────────────────────────
 */`,
			wantHeader:  "─────────────────────────────────────",
			wantBody:    "MIT License\n\nCopyright (c) [year] [fullname]\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.",
			wantFooter:  "─────────────────────────────────────",
			wantSuccess: true,
		},
		{
			name:        "Parse invalid comment",
			content:     "Invalid comment",
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the license block
			header, body, footer, ok := ParseLicenseComponents(tt.content)
			if ok != tt.wantSuccess {
				t.Errorf("ParseLicenseComponents() success = %v, want %v", ok, tt.wantSuccess)
				return
			}

			if ok {
				if header != tt.wantHeader {
					t.Errorf("ParseLicenseComponents() header got:\n%v\nwant:\n%v", header, tt.wantHeader)
				}
				if body != tt.wantBody {
					t.Errorf("ParseLicenseComponents() body got:\n%v\nwant:\n%v", body, tt.wantBody)
				}
				if footer != tt.wantFooter {
					t.Errorf("ParseLicenseComponents() footer got:\n%v\nwant:\n%v", footer, tt.wantFooter)
				}
			}
		})
	}
}

func TestMarkerOperations(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantHas  bool
		stripped string
		marked   string
	}{
		{
			name:     "Text without markers",
			text:     "Test text",
			wantHas:  false,
			stripped: "Test text",
			marked:   markerStart + "Test text" + markerEnd,
		},
		{
			name:     "Text with markers",
			text:     markerStart + "Test text" + markerEnd,
			wantHas:  true,
			stripped: "Test text",
			marked:   markerStart + "Test text" + markerEnd,
		},
		{
			name:     "Empty text",
			text:     "",
			wantHas:  false,
			stripped: "",
			marked:   markerStart + markerEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasMarkers(tt.text); got != tt.wantHas {
				t.Errorf("hasMarkers() = %v, want %v", got, tt.wantHas)
			}

			if got := stripMarkers(tt.text); got != tt.stripped {
				t.Errorf("stripMarkers() = %v, want %v", got, tt.stripped)
			}

			if got := addMarkers(tt.text); got != tt.marked {
				t.Errorf("addMarkers() = %v, want %v", got, tt.marked)
			}
		})
	}
}
