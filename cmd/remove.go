package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/processor"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove license headers from files",
	Long:  `Remove license headers from files that have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appCfg := config.AppConfig{
			// File paths
			LicenseFile: cfgLicense, // Optional for remove command
			Inputs:      ProcessPatterns(cfgInputs),
			Skips:       ProcessPatterns(cfgSkips),

			// Style settings
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go", // default
			PreferMulti:  cfgPreferMulti,

			// Behavior flags
			Verbose: cfgVerbose,

			Force:      false,
			IgnoreFail: false,
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		err = p.Remove()

		cmd.SilenceUsage = true
		return err
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
