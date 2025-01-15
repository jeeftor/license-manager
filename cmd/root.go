// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgLicense     string
	cfgInput       string
	cfgSkip        string
	cfgPrompt      bool
	cfgDryRun      bool
	cfgVerbose     bool   // Add verbose flag
	cfgPresetStyle string // header/footer style
)

var rootCmd = &cobra.Command{
	Use:   "license-manager",
	Short: "A tool to manage license headers in source files",
	Long: `license-manager is a CLI tool that helps manage license headers in source files.
It can add, remove, update, and check license headers in multiple files using patterns.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.PersistentFlags().StringVar(&cfgPresetStyle, "style", "simple", "Preset style for header/footer (e.g., simple, modern, elegant)")

	rootCmd.PersistentFlags().StringVar(&cfgLicense, "license", "", "Path to license text file")
	rootCmd.PersistentFlags().StringVar(&cfgInput, "input", "", "Input file patterns (comma-separated)")
	rootCmd.PersistentFlags().StringVar(&cfgSkip, "skip", "", "Patterns to skip (comma-separated)")
	rootCmd.PersistentFlags().BoolVar(&cfgPrompt, "prompt", false, "Prompt before processing each file")
	rootCmd.PersistentFlags().BoolVar(&cfgDryRun, "dry-run", false, "Show which files would be processed without making changes")
	rootCmd.PersistentFlags().BoolVar(&cfgVerbose, "verbose", false, "Enable verbose output")
}
