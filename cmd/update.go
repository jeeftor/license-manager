package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
	"license-manager/internal/styles"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update license headers in files",
	Long:  `Update existing license headers in files with new content`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("%s", "license file (--license) is required for update command")
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
		return p.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
