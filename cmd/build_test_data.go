// cmd/build_test_data.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildTestDataCmd = &cobra.Command{
	Use:   "build-test-data",
	Short: "Generate test files for all supported languages",
	Long:  `Creates a test_data directory with hello world programs in all supported languages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return BuildTestData()
	},
}

func init() {
	rootCmd.AddCommand(buildTestDataCmd)
}

// TestFile represents a test file to be created
type TestFile struct {
	SourcePath string // Path to template file
	TargetPath string // Path where the file should be created in test_data
}

// getTestFiles returns a list of test files to create
func getTestFiles() []TestFile {
	return []TestFile{
		{"templates/python/hello.py", "test_data/python/hello.py"},
		{"templates/python/hello_utf8.py", "test_data/python/hello_utf8.py"},
		{"templates/ruby/hello.rb", "test_data/ruby/hello.rb"},
		{"templates/javascript/hello.js", "test_data/javascript/hello.js"},
		{"templates/javascript/component.jsx", "test_data/javascript/component.jsx"},
		{"templates/typescript/hello.ts", "test_data/typescript/hello.ts"},
		{"templates/typescript/component.tsx", "test_data/typescript/component.tsx"},
		{"templates/java/HelloWorld.java", "test_data/java/HelloWorld.java"},
		{"templates/go/hello.go", "test_data/go/hello.go"},
		{"templates/go/hello_with_directive.go", "test_data/go/hello_with_directive.go"},
		{"templates/c/hello.c", "test_data/c/hello.c"},
		{"templates/c/hello.h", "test_data/c/hello.h"},
		{"templates/cpp/hello.cpp", "test_data/cpp/hello.cpp"},
		{"templates/cpp/hello.hpp", "test_data/cpp/hello.hpp"},
		{"templates/csharp/Hello.cs", "test_data/csharp/Hello.cs"},
		{"templates/php/hello.php", "test_data/php/hello.php"},
		{"templates/swift/hello.swift", "test_data/swift/hello.swift"},
		{"templates/rust/hello.rs", "test_data/rust/hello.rs"},
		{"templates/shell/hello.sh", "test_data/shell/hello.sh"},
		{"templates/shell/hello.bash", "test_data/shell/hello.bash"},
		{"templates/yaml/config.yml", "test_data/yaml/config.yml"},
		{"templates/yaml/config.yaml", "test_data/yaml/config.yaml"},
		{"templates/perl/hello.pl", "test_data/perl/hello.pl"},
		{"templates/perl/Hello.pm", "test_data/perl/Hello.pm"},
		{"templates/r/hello.r", "test_data/r/hello.r"},
		{"templates/html/index.html", "test_data/html/index.html"},
		{"templates/xml/hello.xml", "test_data/xml/hello.xml"},
		{"templates/css/style.css", "test_data/css/style.css"},
		{"templates/scss/style.scss", "test_data/scss/style.scss"},
		{"templates/sass/style.sass", "test_data/sass/style.sass"},
		{"templates/lua/hello.lua", "test_data/lua/hello.lua"},
		// Add more files here as templates are created
	}
}

// BuildTestData creates test files from templates
func BuildTestData() error {
	files := getTestFiles()

	// Create test_data directory if it doesn't exist
	if err := os.MkdirAll("test_data", 0755); err != nil {
		return fmt.Errorf("failed to create test_data directory: %w", err)
	}

	// Create test files from templates
	for _, file := range files {
		// Create target directory if it doesn't exist
		targetDir := filepath.Dir(file.TargetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Read template content
		content, err := os.ReadFile(file.SourcePath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", file.SourcePath, err)
		}

		// Write content to target file
		if err := os.WriteFile(file.TargetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.TargetPath, err)
		}

		fmt.Printf("Created %s\n", file.TargetPath)
	}

	fmt.Printf("\nSuccessfully created test files in test_data directory\n")
	return nil
}
