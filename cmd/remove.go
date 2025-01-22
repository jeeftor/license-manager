package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/processor"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove license headers from files",
	Long:  `Remove license headers from files that have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("license file (--license) is required for remove command")
		}

		appCfg := config.AppConfig{
			// File paths
			LicenseFile: cfgLicense,
			Inputs:      ProcessPatterns(cfgInputs),
			Skips:       ProcessPatterns(cfgSkips),

			// Style settings
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go", // default
			PreferMulti:  cfgPreferMulti,

			// Behavior flags
			Verbose:     cfgVerbose,
			Interactive: cfgPrompt,
			DryRun:      cfgDryRun,
			Force:       false,
			IgnoreFail:  false,
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
