package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgHeader  string
	cfgFooter  string
	cfgLicense string
	cfgInput   string
	cfgSkip    string
	cfgPrompt  bool
	cfgDryRun  bool
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
	rootCmd.PersistentFlags().StringVar(&cfgHeader, "header", "/* LICENSE HEADER */", "License header text")
	rootCmd.PersistentFlags().StringVar(&cfgFooter, "footer", "/* LICENSE FOOTER */", "License footer text")
	rootCmd.PersistentFlags().StringVar(&cfgLicense, "license", "", "Path to license text file")
	rootCmd.PersistentFlags().StringVar(&cfgInput, "input", "", "Input file patterns (comma-separated)")
	rootCmd.PersistentFlags().StringVar(&cfgSkip, "skip", "", "Patterns to skip (comma-separated)")
	rootCmd.PersistentFlags().BoolVar(&cfgPrompt, "prompt", false, "Prompt before processing each file")
	rootCmd.PersistentFlags().BoolVar(&cfgDryRun, "dry-run", false, "Show which files would be processed without making changes")
}
