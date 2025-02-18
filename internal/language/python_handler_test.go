package language

import (
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

func TestPythonHandler_PreservePreamble(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantPreamble  string
		wantRest      string
		wantSeparator string
	}{
		{
			name: "with shebang and encoding",
			content: `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
def main():
    pass`,
			wantPreamble:  "#!/usr/bin/env python3\n# -*- coding: utf-8 -*-",
			wantRest:      "def main():\n    pass",
			wantSeparator: "\n",
		},
		{
			name: "with shebang only",
			content: `#!/usr/bin/env python3
def main():
    pass`,
			wantPreamble:  "#!/usr/bin/env python3",
			wantRest:      "def main():\n    pass",
			wantSeparator: "\n",
		},
		{
			name: "with encoding only",
			content: `# -*- coding: utf-8 -*-
def main():
    pass`,
			wantPreamble:  "# -*- coding: utf-8 -*-",
			wantRest:      "def main():\n    pass",
			wantSeparator: "\n",
		},
		{
			name: "no preamble",
			content: `def main():
    pass`,
			wantPreamble:  "",
			wantRest:      "def main():\n    pass",
			wantSeparator: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := logger.NewLogger(logger.InfoLevel)
			h := NewPythonHandler(testLogger, styles.HeaderFooterStyle{})
			gotPreamble, gotRest := h.PreservePreamble(tt.content)

			if gotPreamble != tt.wantPreamble {
				t.Errorf(
					"PythonHandler.PreservePreamble() preamble = %q, want %q",
					gotPreamble,
					tt.wantPreamble,
				)
			}
			if gotRest != tt.wantRest {
				t.Errorf(
					"PythonHandler.PreservePreamble() rest = %q, want %q",
					gotRest,
					tt.wantRest,
				)
			}
		})
	}
}

func TestPythonHandler_FormatLicense(t *testing.T) {
	tests := []struct {
		name    string
		license string
		style   styles.HeaderFooterStyle
		want    FullLicenseBlock
	}{
		{
			name:    "single line license",
			license: "Copyright 2025 Example Corp",
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: FullLicenseBlock{
				String: `'''
 * License Header
 * Copyright 2025 Example Corp
 * License Footer
 '''`,
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "multi line license",
			license: `Copyright 2025 Example Corp
All rights reserved.

Licensed under MIT License`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: FullLicenseBlock{
				String: `'''
 * License Header
 * Copyright 2025 Example Corp
 * All rights reserved.
 *
 * Licensed under MIT License
 * License Footer
 '''`,
				Header: "License Header",
				Body: `Copyright 2025 Example Corp
All rights reserved.

Licensed under MIT License`,
				Footer: "License Footer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := logger.NewLogger(logger.InfoLevel)
			h := NewPythonHandler(testLogger, tt.style)
			commentStyle := styles.LanguageExtensions[".py"]
			got := h.FormatLicense(tt.license, commentStyle, tt.style)

			// Normalize line endings
			gotString := strings.ReplaceAll(got.String, "\r\n", "\n")
			wantString := strings.ReplaceAll(tt.want.String, "\r\n", "\n")

			if gotString != wantString {
				t.Errorf(
					"PythonHandler.FormatLicense() String = %q, want %q",
					gotString,
					wantString,
				)
			}

			// Compare normalized components
			gotHeader := normalizeText(got.Header)
			wantHeader := normalizeText(tt.want.Header)
			if gotHeader != wantHeader {
				t.Errorf("Header = %q, want %q", gotHeader, wantHeader)
			}

			gotBody := normalizeText(got.Body)
			wantBody := normalizeText(tt.want.Body)
			if gotBody != wantBody {
				t.Errorf("Body = %q, want %q", gotBody, wantBody)
			}

			gotFooter := normalizeText(got.Footer)
			wantFooter := normalizeText(tt.want.Footer)
			if gotFooter != wantFooter {
				t.Errorf("Footer = %q, want %q", gotFooter, wantFooter)
			}
		})
	}
}

func TestPythonHandler_ExtractComponents(t *testing.T) {
	tests := []struct {
		name    string
		content string
		style   styles.HeaderFooterStyle
		want    ExtractedComponents
	}{
		{
			name: "triple double quotes with markers",
			content: `"""
​License Header‌
Copyright 2025 Example Corp
​License Footer‌
"""`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "triple single quotes with markers",
			content: `'''
​License Header‌
Copyright 2025 Example Corp
​License Footer‌
'''`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "hash comments with markers",
			content: `# ​License Header‌
# Copyright 2025 Example Corp
# ​License Footer‌`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "unicode escape sequences",
			content: `"""
\u200bLicense Header\u200c
Copyright 2025 Example Corp
\u200bLicense Footer\u200c
"""`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "mixed style - hash in triple quotes",
			content: `"""
# ​License Header‌
# Copyright 2025 Example Corp
# ​License Footer‌
"""`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
		{
			name: "with extra whitespace",
			content: `"""
   ​License Header‌
   Copyright 2025 Example Corp
   ​License Footer‌
"""`,
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: ExtractedComponents{
				Header: "License Header",
				Body:   "Copyright 2025 Example Corp",
				Footer: "License Footer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := logger.NewLogger(logger.InfoLevel)
			h := NewPythonHandler(testLogger, tt.style)
			got, success := h.ExtractComponents(tt.content)

			if !success {
				t.Errorf("PythonHandler.ExtractComponents() success = false, want true")
			}

			// Compare normalized versions
			gotHeader := normalizeText(got.Header)
			wantHeader := normalizeText(tt.want.Header)
			if gotHeader != wantHeader {
				t.Errorf("Header = %q, want %q", gotHeader, wantHeader)
			}

			gotBody := normalizeText(got.Body)
			wantBody := normalizeText(tt.want.Body)
			if gotBody != wantBody {
				t.Errorf("Body = %q, want %q", gotBody, wantBody)
			}

			gotFooter := normalizeText(got.Footer)
			wantFooter := normalizeText(tt.want.Footer)
			if gotFooter != wantFooter {
				t.Errorf("Footer = %q, want %q", gotFooter, wantFooter)
			}
		})
	}
}
