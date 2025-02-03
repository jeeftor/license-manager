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
		want    string
	}{
		{
			name:    "single line license",
			license: "Copyright 2025 Example Corp",
			style: styles.HeaderFooterStyle{
				Header: "License Header",
				Footer: "License Footer",
			},
			want: `'''
 * License Header
 * Copyright 2025 Example Corp
 * License Footer
 '''`,
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
			want: `'''
 * License Header
 * Copyright 2025 Example Corp
 * All rights reserved.
 *
 * Licensed under MIT License
 * License Footer
 '''`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := logger.NewLogger(logger.InfoLevel)
			h := NewPythonHandler(testLogger, tt.style)
			commentStyle := styles.LanguageExtensions[".py"]
			got := h.FormatLicense(tt.license, commentStyle, tt.style)

			// Normalize line endings
			got = strings.ReplaceAll(got, "\r\n", "\n")
			tt.want = strings.ReplaceAll(tt.want, "\r\n", "\n")

			if got != tt.want {
				t.Errorf("PythonHandler.FormatLicense() = %q, want %q", got, tt.want)
			}
		})
	}
}
