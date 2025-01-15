package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove license headers from files",
	Long:  `Remove license headers from files that have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := processor.Config{

			LicenseText: cfgLicense,
			Input:       cfgInput,
			Skip:        cfgSkip,
			Prompt:      cfgPrompt,
			DryRun:      cfgDryRun,
			Verbose:     cfgVerbose,
		}

		p := processor.NewFileProcessor(config)
		return p.Remove()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
