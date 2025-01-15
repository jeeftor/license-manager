package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/license-manager/internal/processor"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add license headers to files",
	Long:  `Add license headers to files that don't already have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := processor.Config{
			Header:      cfgHeader,
			Footer:      cfgFooter,
			LicenseText: cfgLicense,
			Input:       cfgInput,
			Skip:        cfgSkip,
			Prompt:      cfgPrompt,
			DryRun:      cfgDryRun,
		}

		p := processor.NewFileProcessor(config)
		return p.Add()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
