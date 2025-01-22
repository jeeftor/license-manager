// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

var (
	cfgLicense      string
	cfgInputs       []string
	cfgSkips        []string
	cfgPrompt       bool
	cfgDryRun       bool
	cfgVerbose      bool   // Add verbose flag
	cfgPresetStyle  string // header/footer style
	cfgPreferMulti  bool   // prefer multi-line comments where supported
	checkIgnoreFail bool   // Added for check command

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
	rootCmd.PersistentFlags().BoolVar(&cfgPreferMulti, "multi", true, "Prefer multi-line comments where supported")

	rootCmd.PersistentFlags().StringVar(&cfgLicense, "license", "", "Path to license text file")

	rootCmd.PersistentFlags().StringSliceVar(&cfgInputs, "input", []string{}, "Inputs file patterns")
	rootCmd.PersistentFlags().StringSliceVar(&cfgSkips, "skip", []string{}, "Patterns to skip")

	rootCmd.PersistentFlags().BoolVar(&cfgPrompt, "prompt", false, "Prompt before processing each file")
	rootCmd.PersistentFlags().BoolVar(&cfgDryRun, "dry-run", false, "Show which files would be processed without making changes")
	rootCmd.PersistentFlags().BoolVar(&cfgVerbose, "verbose", false, "Enable verbose output")
}

func ProcessPatterns(patterns []string) string {
	var result []string
	for _, p := range patterns {
		// Split on commas if present
		parts := strings.Split(p, ",")
		result = append(result, parts...)
	}
	return strings.Join(result, ",")
}
