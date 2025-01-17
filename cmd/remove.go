package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
	"license-manager/internal/styles"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove license headers from files",
	Long:  `Remove license headers from files that have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("%s", "license file (--license) is required for remove command")
		}

		config := &processor.Config{
			Input:       cfgInput,
			Skip:        cfgSkip,
			Prompt:      cfgPrompt,
			DryRun:      cfgDryRun,
			Verbose:     cfgVerbose,
			PreferMulti: cfgPreferMulti,
		}

		style := styles.GetPresetStyle(cfgPresetStyle)
		p := processor.NewFileProcessor(config, cfgLicense, style)
		return p.Remove()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
