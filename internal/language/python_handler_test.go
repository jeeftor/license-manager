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
	// Use real preset styles so the output matches what FormatComment produces.
	// Python's comment style has MultiPrefix="" and LinePrefix="", so content
	// goes directly between """ markers without any * prefixes.
	hashStyle := styles.Get("hash")
	commentStyle := styles.GetLanguageCommentStyle(".py")

	t.Run("single line license", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)
		got := h.FormatLicense("Copyright 2025 Example Corp", commentStyle, hashStyle)

		// Should use """ delimiters (Python multi-line)
		if !strings.Contains(got.String, `"""`) {
			t.Errorf("Expected \"\"\" delimiters in output, got:\n%s", got.String)
		}

		// Should NOT have * prefixes (that's Go style, not Python)
		if strings.Contains(got.String, " * ") {
			t.Errorf("Should not have Go-style ' * ' prefixes in Python output:\n%s", got.String)
		}

		// Should contain the license text
		if !strings.Contains(got.String, "Copyright 2025 Example Corp") {
			t.Errorf("Output should contain license text, got:\n%s", got.String)
		}

		// Should contain hash header/footer pattern
		if !strings.Contains(got.String, "######") {
			t.Errorf("Output should contain hash pattern, got:\n%s", got.String)
		}

		// Body should contain the license text
		gotBody := normalizeText(got.Body)
		if !strings.Contains(gotBody, "Copyright 2025 Example Corp") {
			t.Errorf("Body = %q, should contain license text", gotBody)
		}
	})

	t.Run("multi line license", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)
		license := "Copyright 2025 Example Corp\nAll rights reserved.\n\nLicensed under MIT License"
		got := h.FormatLicense(license, commentStyle, hashStyle)

		// Should contain all license lines
		if !strings.Contains(got.String, "Copyright 2025 Example Corp") {
			t.Error("Missing first license line")
		}
		if !strings.Contains(got.String, "All rights reserved.") {
			t.Error("Missing second license line")
		}
		if !strings.Contains(got.String, "Licensed under MIT License") {
			t.Error("Missing third license line")
		}

		// Body should contain all lines
		gotBody := normalizeText(got.Body)
		if !strings.Contains(gotBody, "Copyright 2025 Example Corp") {
			t.Errorf("Body should contain license text, got %q", gotBody)
		}
	})
}

func TestPythonHandler_ExtractComponents(t *testing.T) {
	// Use a real preset style so styles.Infer() can match the header/footer.
	// The old tests used fake style names like "License Header" which Infer() can never match.
	hashStyle := styles.Get("hash")
	licenseText := "Copyright 2025 Example Corp"

	t.Run("system-generated multi-line license", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)

		// Generate a license block using FormatComment (what production uses)
		commentStyle := styles.GetLanguageCommentStyle(".py")
		formatted := FormatComment(licenseText, commentStyle, hashStyle)

		// Now extract components from the formatted license
		got, success := h.ExtractComponents(formatted)
		if !success {
			t.Fatalf(
				"ExtractComponents() returned false for system-generated license:\n%s",
				formatted,
			)
		}

		// Header should contain the hash style pattern
		if !strings.Contains(got.Header, "#") {
			t.Errorf("Header should contain hash pattern, got %q", got.Header)
		}

		// Body should contain the license text
		gotBody := normalizeText(got.Body)
		if !strings.Contains(gotBody, "Copyright 2025 Example Corp") {
			t.Errorf("Body should contain license text, got %q", gotBody)
		}

		// Footer should contain the hash style pattern
		if !strings.Contains(got.Footer, "#") {
			t.Errorf("Footer should contain hash pattern, got %q", got.Footer)
		}
	})

	t.Run("system-generated license with preamble", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)

		commentStyle := styles.GetLanguageCommentStyle(".py")
		formatted := FormatComment(licenseText, commentStyle, hashStyle)

		// Add a shebang preamble before the license
		contentWithPreamble := "#!/usr/bin/env python3\n" + formatted + "\n\nimport os\n"

		got, success := h.ExtractComponents(contentWithPreamble)
		if !success {
			t.Fatalf(
				"ExtractComponents() returned false for license with preamble:\n%s",
				contentWithPreamble,
			)
		}

		// Preamble should contain the shebang
		if !strings.Contains(got.Preamble, "#!/usr/bin/env python3") {
			t.Errorf("Preamble should contain shebang, got %q", got.Preamble)
		}

		// Body should contain the license text
		gotBody := normalizeText(got.Body)
		if !strings.Contains(gotBody, "Copyright 2025 Example Corp") {
			t.Errorf("Body should contain license text, got %q", gotBody)
		}

		// Rest should contain the import
		if !strings.Contains(got.Rest, "import os") {
			t.Errorf("Rest should contain remaining code, got %q", got.Rest)
		}
	})

	t.Run("no license present", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)

		content := "#!/usr/bin/env python3\n\nimport os\n\ndef main():\n    pass\n"
		_, success := h.ExtractComponents(content)
		if success {
			t.Error("ExtractComponents() should return false when no license is present")
		}
	})

	t.Run("empty content", func(t *testing.T) {
		testLogger := logger.NewLogger(logger.ErrorLevel)
		h := NewPythonHandler(testLogger, hashStyle)

		_, success := h.ExtractComponents("")
		if success {
			t.Error("ExtractComponents() should return false for empty content")
		}
	})

	t.Run("roundtrip with different styles", func(t *testing.T) {
		// Test that extraction works with multiple preset styles
		for _, styleName := range []string{"hash", "equals", "stars", "simple"} {
			t.Run(styleName, func(t *testing.T) {
				style := styles.Get(styleName)
				testLogger := logger.NewLogger(logger.ErrorLevel)
				h := NewPythonHandler(testLogger, style)

				commentStyle := styles.GetLanguageCommentStyle(".py")
				formatted := FormatComment(licenseText, commentStyle, style)

				got, success := h.ExtractComponents(formatted)
				if !success {
					t.Fatalf("ExtractComponents() failed for style %q:\n%s", styleName, formatted)
				}

				gotBody := normalizeText(got.Body)
				if !strings.Contains(gotBody, "Copyright 2025 Example Corp") {
					t.Errorf(
						"Body should contain license text for style %q, got %q",
						styleName,
						gotBody,
					)
				}
			})
		}
	})
}
