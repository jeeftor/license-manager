package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestLoggerOutputComparison compares the two different logging methods
func TestLoggerOutputComparison(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
	}{
		// Original test cases
		{
			name:   "simple string no args",
			format: "test message",
			args:   nil,
		},
		{
			name:   "string with integer arg",
			format: "Scanning %d Directories:",
			args:   []interface{}{1},
		},
		{
			name:   "string with string arg",
			format: "Processing file: %s",
			args:   []interface{}{"./test_data/file.txt"},
		},
		{
			name:   "string with brackets",
			format: "Processing file: [%s]",
			args:   []interface{}{"./test_data/file.txt"},
		},
		{
			name:   "multiple args",
			format: "File: %s, Lines: %d",
			args:   []interface{}{"./test_data/file.txt", 100},
		},
		{
			name:   "empty args slice",
			format: "Just a message",
			args:   []interface{}{},
		},
		{
			name:   "empty string arg",
			format: "Message: %s",
			args:   []interface{}{""},
		},
		// New edge cases
		{
			name:   "format string with percent sign",
			format: "Progress: 100%%",
			args:   nil,
		},
		{
			name:   "format string with multiple percent signs",
			format: "Progress: %d%% of %d%%",
			args:   []interface{}{50, 100},
		},
		{
			name:   "nil argument",
			format: "Value: %v",
			args:   []interface{}{nil},
		},
		{
			name:   "format string with newlines",
			format: "Line1\nLine2: %s\nLine3",
			args:   []interface{}{"middle"},
		},
		{
			name:   "empty format string",
			format: "",
			args:   nil,
		},
		{
			name:   "format string with just spaces",
			format: "   ",
			args:   nil,
		},
		{
			name:   "mixed type arguments",
			format: "%v %d %s %t %f",
			args:   []interface{}{"str", 42, "text", true, 3.14},
		},
		{
			name:   "unicode characters",
			format: "Unicode: %s %s",
			args:   []interface{}{"🚀", "⭐"},
		},
		{
			name:   "format string with quotes",
			format: "Quoted: \"%s\" and '%s'",
			args:   []interface{}{"double", "single"},
		},
		{
			name:   "format string with tabs",
			format: "Col1\t%s\tCol2\t%s",
			args:   []interface{}{"val1", "val2"},
		},
		{
			name:   "special characters in arguments",
			format: "Special: %s",
			args:   []interface{}{"\n\t\r"},
		},
		{
			name:   "wrong number of arguments (too few)",
			format: "%s %s %s",
			args:   []interface{}{"one", "two"},
		},
		{
			name:   "wrong number of arguments (too many)",
			format: "%s",
			args:   []interface{}{"one", "two", "three"},
		},
		{
			name:   "wrong format specifier",
			format: "%d",
			args:   []interface{}{"not a number"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(false, 0)

			// Capture output from direct Printf
			directOutput := captureOutput(func() {
				fmt.Printf("%s %s\n",
					logger.colors["info"].Sprint("INFO:"),
					fmt.Sprintf(tt.format, tt.args...))
			})

			// Capture output from Log method
			logOutput := captureOutput(func() {
				logger.Log(InfoLevel, true, tt.format, tt.args...)
			})

			// Clean both outputs
			directOutput = stripANSI(directOutput)
			logOutput = stripANSI(logOutput)

			// Compare the outputs
			if directOutput != logOutput {
				t.Errorf("\nFormat: %q\nArgs: %v\nDirect Printf output: %q\nLog method output: %q",
					tt.format, tt.args, directOutput, logOutput)
			}
		})
	}
}

// stripANSI removes ANSI color codes from a string
func stripANSI(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	return strings.ReplaceAll(str, "\x1b[0m", "")
}
