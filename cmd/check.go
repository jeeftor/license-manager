package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/license-manager/internal/processor"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for license headers in files",
	Long:  `Check if files have the specified license headers`,
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
		return p.Check()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
