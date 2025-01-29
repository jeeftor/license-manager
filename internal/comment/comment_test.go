package comment

import (
	"strings"
	"testing"

	"license-manager/internal/language"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
)

// MockLanguageHandler implements language.LanguageHandler for testing
type MockLanguageHandler struct {
	*language.GenericHandler
}

func NewMockLanguageHandler() *MockLanguageHandler {
	return &MockLanguageHandler{
		GenericHandler: language.NewGenericHandler(nil, styles.HeaderFooterStyle{}),
	}
}

func (h *MockLanguageHandler) ScanBuildDirectives(content string) ([]string, int) {
	return nil, 0
}

func TestExtractComponents(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		stripMarkers bool
		wantHeader   string
		wantBody     string
		wantFooter   string
		wantSuccess  bool
	}{
		{
			name: "C-style comment with markers",
			content: `/* Copyright (c) 2025 */
This is the body text
/* All rights reserved */`,
			stripMarkers: true,
			wantHeader:   "Copyright (c) 2025",
			wantBody:     "This is the body text",
			wantFooter:   "All rights reserved",
			wantSuccess:  true,
		},
		{
			name: "C-style comment without stripping",
			content: `/* Copyright (c) 2025 */
This is the body text
/* All rights reserved */`,
			stripMarkers: false,
			wantHeader:   "/* Copyright (c) 2025 */",
			wantBody:     "This is the body text",
			wantFooter:   "/* All rights reserved */",
			wantSuccess:  true,
		},
		{
			name: "Python triple single quote comment",
			content: `'''
 * License Header
 * Copyright 2025 Example Corp
 * License Footer
 '''

def main():
    pass`,
			stripMarkers: true,
			wantHeader:   "License Header",
			wantBody:     "Copyright 2025 Example Corp",
			wantFooter:   "License Footer",
			wantSuccess:  true,
		},
		{
			name: "Python triple double quote comment",
			content: `"""
 * License Header
 * Copyright 2025 Example Corp
 * License Footer
 """

def main():
    pass`,
			stripMarkers: true,
			wantHeader:   "License Header",
			wantBody:     "Copyright 2025 Example Corp",
			wantFooter:   "License Footer",
			wantSuccess:  true,
		},
		{
			name: "Python hash comment",
			content: `# License Header
# Copyright 2025 Example Corp
# License Footer

def main():
    pass`,
			stripMarkers: true,
			wantHeader:   "License Header",
			wantBody:     "Copyright 2025 Example Corp",
			wantFooter:   "License Footer",
			wantSuccess:  true,
		},
		{
			name: "HTML-style comment",
			content: `<!-- MIT License -->
Permission is hereby granted
<!-- End License -->`,
			stripMarkers: true,
			wantHeader:   "MIT License",
			wantBody:     "Permission is hereby granted",
			wantFooter:   "End License",
			wantSuccess:  true,
		},
		{
			name: "Known style with stars",
			content: `****************************************
Permission is hereby granted, free of charge
****************************************`,
			stripMarkers: false,
			wantHeader:   "****************************************",
			wantBody:     "Permission is hereby granted, free of charge",
			wantFooter:   "****************************************",
			wantSuccess:  true,
		},
		{
			name: "Known style with box",
			content: `+------------------------------------+
Licensed under Apache 2.0
+------------------------------------+`,
			stripMarkers: false,
			wantHeader:   "+------------------------------------+",
			wantBody:     "Licensed under Apache 2.0",
			wantFooter:   "+------------------------------------+",
			wantSuccess:  true,
		},
		{
			name: "Unknown style should fail",
			content: `Unknown Header Style
Some body text
Unknown Footer Style`,
			stripMarkers: false,
			wantSuccess:  false,
		},
		{
			name:        "Empty content",
			content:     "",
			wantSuccess: false,
		},
		{
			name:        "Single line",
			content:     "Just one line",
			wantSuccess: false,
		},
		{
			name: "No clear header/footer",
			content: `Line 1
Line 2
Line 3`,
			wantSuccess: false,
		},
		{
			name: "Multiple empty lines with stars style",
			content: `****************************************


Body content


****************************************`,
			stripMarkers: false,
			wantHeader:   "****************************************",
			wantBody:     "Body content",
			wantFooter:   "****************************************",
			wantSuccess:  true,
		},
		{
			name: "Mixed styles should fail",
			content: `****************************************
Body content
+------------------------------------+`,
			stripMarkers: false,
			wantSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var style styles.CommentLanguage
			if strings.Contains(tt.content, `'''`) || strings.Contains(tt.content, `"""`) {
				style = styles.CommentLanguage{
					Single:     "#",
					MultiStart: "'''",
					MultiEnd:   "'''",
				}
				if strings.Contains(tt.content, `"""`) {
					style.MultiStart = `"""`
					style.MultiEnd = `"""`
				}
			} else {
				style = styles.CommentLanguage{
					Single:     "//",
					MultiStart: "/*",
					MultiEnd:   "*/",
				}
			}
			testLogger := logger.NewLogger(logger.InfoLevel)
			gotHeader, gotBody, gotFooter, gotSuccess := language.ExtractComponents(testLogger, tt.content, tt.stripMarkers, style)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("ExtractComponents() success = %v, want %v", gotSuccess, tt.wantSuccess)
				return
			}

			if !tt.wantSuccess {
				return
			}

			if gotHeader != tt.wantHeader {
				t.Errorf("ExtractComponents() header = %q, want %q", gotHeader, tt.wantHeader)
			}
			if gotBody != tt.wantBody {
				t.Errorf("ExtractComponents() body = %q, want %q", gotBody, tt.wantBody)
			}
			if gotFooter != tt.wantFooter {
				t.Errorf("ExtractComponents() footer = %q, want %q", gotFooter, tt.wantFooter)
			}
		})
	}
}

func TestComment_String(t *testing.T) {
	tests := []struct {
		name     string
		style    styles.CommentLanguage
		hfStyle  styles.HeaderFooterStyle
		body     string
		expected string
	}{
		{
			name: "Single-line comment style",
			style: styles.CommentLanguage{
				Single:     "//",
				LinePrefix: " ",
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "Copyright Header",
				Footer: "License Footer",
			},
			body: "Body text\nSecond line",
			expected: `// Copyright Header

// Body text
// Second line

// License Footer`,
		},
		{
			name: "Multi-line comment style",
			style: styles.CommentLanguage{
				MultiStart:  "/*",
				MultiEnd:    "*/",
				MultiPrefix: " *",
				LinePrefix:  " ",
				PreferMulti: true,
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "MIT License",
				Footer: "End License",
			},
			body: "Permission is hereby granted",
			expected: `/*
 * MIT License
 *
 * Permission is hereby granted
 *
 * End License
*/`,
		},
		{
			name: "HTML comment style",
			style: styles.CommentLanguage{
				MultiStart:  "<!--",
				MultiEnd:    "-->",
				PreferMulti: true,
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "License Notice",
				Footer: "End Notice",
			},
			body: "Content here",
			expected: `<!--
License Notice

Content here

End Notice
-->`,
		},
		{
			name: "Python triple single quote comment style",
			style: styles.CommentLanguage{
				Single:     "#",
				MultiStart: "'''",
				MultiEnd:   "'''",
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			body: "Body text\nSecond line",
			expected: `'''
License Header

Body text
Second line

License Footer
'''`,
		},
		{
			name: "Python triple double quote comment style",
			style: styles.CommentLanguage{
				Single:     "#",
				MultiStart: `"""`,
				MultiEnd:   `"""`,
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			body: "Body text\nSecond line",
			expected: `"""
License Header

Body text
Second line

License Footer
"""`,
		},
		{
			name: "Python hash comment style",
			style: styles.CommentLanguage{
				Single: "#",
			},
			hfStyle: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			body: "Body text\nSecond line",
			expected: `# License Header
# Body text
# Second line
# License Footer`,
		},
	}

	mockHandler := NewMockLanguageHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := language.NewComment(tt.style, tt.hfStyle, tt.body, mockHandler)
			got := c.String()
			// Normalize line endings for comparison
			got = strings.ReplaceAll(got, "\r\n", "\n")
			tt.expected = strings.ReplaceAll(tt.expected, "\r\n", "\n")

			if got != tt.expected {
				t.Errorf("Comment.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestComment_Clone(t *testing.T) {
	mockHandler := NewMockLanguageHandler()
	original := language.NewComment(
		styles.CommentLanguage{
			Single: "//",
		},
		styles.HeaderFooterStyle{
			Header: "Header",
			Footer: "Footer",
		},
		"Body",
		mockHandler,
	)

	clone := original.Clone()

	// Verify clone has same values
	if clone.header != original.header {
		t.Errorf("Clone header = %q, want %q", clone.header, original.header)
	}
	if clone.body != original.body {
		t.Errorf("Clone body = %q, want %q", clone.body, original.body)
	}
	if clone.footer != original.footer {
		t.Errorf("Clone footer = %q, want %q", clone.footer, original.footer)
	}

	// Verify modifying clone doesn't affect original
	clone.SetBody("Modified Body")
	if original.body == clone.body {
		t.Error("Modifying clone body affected original")
	}

	clone.SetHeaderFooterStyle(styles.HeaderFooterStyle{
		Header: "Modified Header",
		Footer: "Modified Footer",
	})
	if original.header == clone.header || original.footer == clone.footer {
		t.Error("Modifying clone header/footer affected original")
	}
}
