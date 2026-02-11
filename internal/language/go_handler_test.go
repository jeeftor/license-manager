package language

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
	"github.com/stretchr/testify/assert"
)

func TestGoHandler_isDirective(t *testing.T) {
	handler := NewGoHandler(nil, styles.HeaderFooterStyle{})
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Build directive", "//go:build darwin", true},
		{"Generate directive", "//go:generate mockgen", true},
		{"Plus build directive", "// +build darwin linux", true},
		{"Plus build no space", "//+build darwin linux", true},
		{"Regular comment", "// This is a comment", false},
		{"Code line", "package main", false},
		{"Empty line", "", false},
		{"Indented directive", "    //go:build darwin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.isDirective(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGoHandler_isGenerateDirective(t *testing.T) {
	handler := NewGoHandler(nil, styles.HeaderFooterStyle{})
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Generate directive", "//go:generate mockgen", true},
		{"Build directive", "//go:build darwin", false},
		{"Regular comment", "// This is a comment", false},
		{"Indented generate", "    //go:generate mockgen", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.isGenerateDirective(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGoHandler_ScanBuildDirectives(t *testing.T) {
	testLogger := logger.NewLogger(logger.InfoLevel)
	handler := NewGoHandler(testLogger, styles.HeaderFooterStyle{})

	tests := []struct {
		name           string
		content        string
		wantDirectives []string
		wantEndIndex   int
	}{
		{
			name: "single build directive",
			content: `//go:build linux
package main`,
			wantDirectives: []string{
				"//go:build linux",
			},
			wantEndIndex: 1,
		},
		{
			name: "multiple build directives with one blank line",
			content: `//go:build linux
// +build linux

//go:build !windows
// +build !windows

package main`,
			wantDirectives: []string{
				"//go:build linux",
				"// +build linux",
				"",
				"//go:build !windows",
				"// +build !windows",
				"",
			},
			wantEndIndex: 6,
		},
		{
			name: "real-world example with build and generate directives",
			content: `//go:build (darwin || linux) && !arm
//go:build !windows
//go:build cgo
//go:build !no_protobuf
//go:build go1.18
//go:build darwin || linux
// +build darwin linux

//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
//go:generate stringer -type=MyEnumType
//go:generate command
//go:linkname
package main2`,
			wantDirectives: []string{
				"//go:build (darwin || linux) && !arm",
				"//go:build !windows",
				"//go:build cgo",
				"//go:build !no_protobuf",
				"//go:build go1.18",
				"//go:build darwin || linux",
				"// +build darwin linux",
				"",
				"//go:generate mockgen -source=myfile.go",
				"//go:generate protoc --go_out=. myproto.proto",
				"//go:generate stringer -type=MyEnumType",
				"//go:generate command",
				"//go:linkname",
			},
			wantEndIndex: 13,
		},
		{
			name: "directives separated by multiple blank lines",
			content: `//go:build linux
// +build linux


//go:build !windows
// +build !windows

package main`,
			wantDirectives: []string{
				"//go:build linux",
				"// +build linux",
				"",
				"",
				"//go:build !windows",
				"// +build !windows",
			},
			wantEndIndex: 7,
		},
		{
			name: "directives after package declaration",
			content: `package main

//go:build linux
// +build linux`,
			wantDirectives: nil,
			wantEndIndex:   0,
		},
		{
			name: "mixed directive styles",
			content: `//go:build linux
//+build linux
// +build !windows

package main`,
			wantDirectives: []string{
				"//go:build linux",
				"//+build linux",
				"// +build !windows",
				"",
			},
			wantEndIndex: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirectives, gotEndIndex := handler.ScanBuildDirectives(tt.content)

			// Detailed error messages for better debugging
			if len(gotDirectives) != len(tt.wantDirectives) {
				t.Errorf(
					"ScanBuildDirectives() got %d directives, want %d\nGot directives:\n%s\nWant directives:\n%s",
					len(gotDirectives),
					len(tt.wantDirectives),
					strings.Join(gotDirectives, "\n"),
					strings.Join(tt.wantDirectives, "\n"),
				)
				return
			}

			for i := range gotDirectives {
				if strings.TrimSpace(gotDirectives[i]) != strings.TrimSpace(tt.wantDirectives[i]) {
					t.Errorf(
						"ScanBuildDirectives() directive[%d] = %q, want %q",
						i,
						gotDirectives[i],
						tt.wantDirectives[i],
					)
				}
			}

			assert.Equal(t, tt.wantEndIndex, gotEndIndex, "EndIndex mismatch")
		})
	}
}

func TestGoHandler_PreservePreamble(t *testing.T) {
	handler := NewGoHandler(nil, styles.HeaderFooterStyle{})

	tests := []struct {
		name         string
		content      string
		wantPreamble string
		wantRest     string
	}{
		{
			name: "build and generate directives",
			content: `//go:build (darwin || linux) && !arm
//go:build !windows
//go:build cgo
//go:build !no_protobuf
//go:build go1.18
//go:build darwin || linux
// +build darwin linux

//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
//go:generate stringer -type=MyEnumType
//go:generate command
//go:linkname

package main2

import "fmt"

func main() {}`,
			wantPreamble: `//go:build (darwin || linux) && !arm
//go:build !windows
//go:build cgo
//go:build !no_protobuf
//go:build go1.18
//go:build darwin || linux
// +build darwin linux

//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
//go:generate stringer -type=MyEnumType
//go:generate command
//go:linkname
`,
			wantRest: `package main2

import "fmt"

func main() {}`,
		},
		{
			name: "only build directives",
			content: `//go:build linux
// +build linux

package main

func main() {}`,
			wantPreamble: `//go:build linux
// +build linux
`,
			wantRest: `package main

func main() {}`,
		},
		{
			name: "directives with comments",
			content: `//go:build linux
// +build linux

// This is a regular comment
package main

func main() {}`,
			wantPreamble: `//go:build linux
// +build linux
`,
			wantRest: `// This is a regular comment
package main

func main() {}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPreamble, gotRest := handler.PreservePreamble(tt.content)
			assert.Equal(t, tt.wantPreamble, gotPreamble, "Preamble mismatch")
			assert.Equal(t, tt.wantRest, gotRest, "Rest content mismatch")
		})
	}
}

// TestGoHandler_ScanBuildDirectivesFromTemplates tests the handler against actual template files
func TestGoHandler_ScanBuildDirectivesFromTemplates(t *testing.T) {
	templateDir := "../../templates/go"
	files, err := os.ReadDir(templateDir)
	if err != nil {
		t.Skip("Skipping template tests: template directory not found")
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".go.tmpl") {
			continue
		}

		t.Run(strings.TrimSuffix(file.Name(), ".tmpl"), func(t *testing.T) {
			testLogger := logger.NewLogger(logger.InfoLevel)
			content, err := os.ReadFile(filepath.Join(templateDir, file.Name()))
			if err != nil {
				t.Fatalf("Failed to read template file: %v", err)
			}

			handler := NewGoHandler(testLogger, styles.HeaderFooterStyle{})
			directives, endIndex := handler.ScanBuildDirectives(string(content))

			// Verify directives come before package declaration
			packageLine := -1
			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "package ") {
					packageLine = i
					break
				}
			}

			if packageLine >= 0 && endIndex > packageLine {
				t.Errorf(
					"Directives end at line %d, but package declaration is at line %d",
					endIndex,
					packageLine,
				)
			}

			// Verify files with _with_directive suffix contain directives
			if strings.HasSuffix(file.Name(), "_with_directive.go.tmpl") {
				assert.Greater(
					t,
					len(directives),
					0,
					"Expected directives in file with _with_directive suffix",
				)
			}
		})
	}
}
