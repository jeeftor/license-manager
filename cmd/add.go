package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/processor"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add license headers to files",
	Long:  `Add license headers to files that don't already have them`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("license file (--license) is required for add command")
		}

		appCfg := config.AppConfig{
			// File paths
			LicenseFile: cfgLicense,
			Input:       cfgInput,
			Skip:        cfgSkip,

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
		err = p.Add()

		cmd.SilenceUsage = true
		return err
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
