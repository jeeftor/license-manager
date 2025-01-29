// cmd/build_test_data.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeeftor/license-manager/internal/logger"

	"github.com/spf13/cobra"
)

var log = logger.NewLogger(logger.DebugLevel)

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

// getTestFiles returns a list of test files to create by scanning the templates directory
func getTestFiles() []TestFile {
	var files []TestFile
	templatesDir := "templates"

	// Walk through the templates directory
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Get the relative path from templates directory
		relPath, err := filepath.Rel(templatesDir, path)
		if err != nil {
			log.LogError("Error getting relative path for %s: %v", path, err)
			return nil
		}

		// Create the corresponding test_data path
		testPath := filepath.Join("test_data", relPath)

		files = append(files, TestFile{
			SourcePath: path,
			TargetPath: testPath,
		})

		return nil
	})

	if err != nil {
		log.LogError("Error walking templates directory: %v", err)
		return nil
	}

	if len(files) == 0 {
		log.LogError("No template files found in %s", templatesDir)
	} else {
		log.LogInfo("Found %d template files", len(files))
	}

	return files
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
