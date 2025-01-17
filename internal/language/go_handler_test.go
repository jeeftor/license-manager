package language

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"license-manager/internal/styles"
)

func TestGoHandler_ScanBuildDirectives(t *testing.T) {
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
				"",
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
				"",
			},
			wantEndIndex: 14,
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
			},
			wantEndIndex: 3,
		},
		{
			name: "no directives",
			content: `package main

func main() {}`,
			wantDirectives: nil,
			wantEndIndex:   0,
		},
		{
			name: "directives at end of file",
			content: `//go:build linux
// +build linux
`,
			wantDirectives: []string{
				"//go:build linux",
				"// +build linux",
				"",
			},
			wantEndIndex: 3,
		},
		{
			name: "directives with comments in between",
			content: `//go:build linux
// +build linux

// This is a regular comment
//go:build !windows
// +build !windows

package main`,
			wantDirectives: []string{
				"//go:build linux",
				"// +build linux",
				"",
			},
			wantEndIndex: 3,
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
		{
			name: "directives after package declaration",
			content: `package main

//go:build linux
// +build linux`,
			wantDirectives: nil,
			wantEndIndex:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGoHandler(styles.HeaderFooterStyle{})
			gotDirectives, gotEndIndex := handler.ScanBuildDirectives(tt.content)

			// Compare directives
			if len(gotDirectives) != len(tt.wantDirectives) {
				t.Errorf("ScanBuildDirectives() got %d directives, want %d", len(gotDirectives), len(tt.wantDirectives))
				t.Errorf("Got directives:\n%s", strings.Join(gotDirectives, "\n"))
				t.Errorf("Want directives:\n%s", strings.Join(tt.wantDirectives, "\n"))
				return
			}
			for i := range gotDirectives {
				if strings.TrimSpace(gotDirectives[i]) != strings.TrimSpace(tt.wantDirectives[i]) {
					t.Errorf("ScanBuildDirectives() directive[%d] = %q, want %q", i, gotDirectives[i], tt.wantDirectives[i])
				}
			}

			// Compare end index
			if gotEndIndex != tt.wantEndIndex {
				t.Errorf("ScanBuildDirectives() gotEndIndex = %v, want %v", gotEndIndex, tt.wantEndIndex)
			}
		})
	}
}

func TestGoHandler_PreservePreamble(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantPreamble  string
		wantRest      string
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
			name: "only generate directives",
			content: `//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto

package main

func main() {}`,
			wantPreamble: `//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
`,
			wantRest: `package main

func main() {}`,
		},
		{
			name: "no directives",
			content: `package main

func main() {}`,
			wantPreamble: "",
			wantRest: `package main

func main() {}`,
		},
		{
			name: "directives at end of file",
			content: `//go:build linux
// +build linux`,
			wantPreamble: `//go:build linux
// +build linux
`,
			wantRest: "",
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
			handler := NewGoHandler(styles.HeaderFooterStyle{})
			gotPreamble, gotRest := handler.PreservePreamble(tt.content)

			// Compare preamble
			if gotPreamble != tt.wantPreamble {
				t.Errorf("PreservePreamble() preamble = %q, want %q", gotPreamble, tt.wantPreamble)
			}

			// Compare rest
			if gotRest != tt.wantRest {
				t.Errorf("PreservePreamble() rest = %q, want %q", gotRest, tt.wantRest)
			}
		})
	}
}

func TestGoHandler_ScanBuildDirectivesFromTemplates(t *testing.T) {
	templateDir := "../../templates/go"
	files, err := os.ReadDir(templateDir)
	if err != nil {
		t.Fatalf("Failed to read template directory: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(templateDir, file.Name()))
			if err != nil {
				t.Fatalf("Failed to read template file: %v", err)
			}

			handler := NewGoHandler(styles.HeaderFooterStyle{})
			directives, endIndex := handler.ScanBuildDirectives(string(content))

			// Log the results for inspection
			t.Logf("File: %s", file.Name())
			t.Logf("Found %d directives", len(directives))
			t.Logf("Directives end at line %d", endIndex)
			for i, d := range directives {
				t.Logf("Directive[%d]: %q", i, d)
			}

			// Basic validation
			if strings.HasSuffix(file.Name(), "_with_directive.go") {
				if len(directives) == 0 {
					t.Error("Expected directives in file with _with_directive suffix")
				}
			}

			// Verify that directives come before package declaration
			packageLine := -1
			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "package ") {
					packageLine = i
					break
				}
			}

			if packageLine >= 0 && endIndex > packageLine {
				t.Errorf("Directives end at line %d, but package declaration is at line %d", endIndex, packageLine)
			}
		})
	}
}
