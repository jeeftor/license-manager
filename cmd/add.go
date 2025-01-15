package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add license headers to files",
	Long:  `Add license headers to files that don't already have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("%s", "license file (--license) is required for add command")
		}

		config := processor.Config{
			LicenseText: cfgLicense,
			Input:       cfgInput,
			Skip:        cfgSkip,
			Prompt:      cfgPrompt,
			DryRun:      cfgDryRun,
			Verbose:     cfgVerbose,
		}

		p := processor.NewFileProcessor(config)
		err := p.Add()
		
		// Don't show usage for any error from Add()
		cmd.SilenceUsage = true
		
		return err
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
