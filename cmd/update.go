package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update license headers in files",
	Long:  `Update existing license headers in files with new content`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := processor.Config{
			Header:      cfgHeader,
			Footer:      cfgFooter,
			LicenseText: cfgLicense,
			Input:       cfgInput,
			Skip:        cfgSkip,
			Prompt:      cfgPrompt,
			DryRun:      cfgDryRun,
			Verbose:     cfgVerbose,
		}

		p := processor.NewFileProcessor(config)
		return p.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
