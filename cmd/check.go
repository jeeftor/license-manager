package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/processor"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for license headers in files",
	Long:  `Check if files have the specified license headers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("license file (--license) is required for check command")
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
			IgnoreFail:  checkIgnoreFail, // Special flag for check command
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		err = p.Check()

		// If it's a CheckError, don't show usage
		if _, isCheckError := err.(*processor.CheckError); isCheckError {
			cmd.SilenceUsage = true
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&checkIgnoreFail, "ignore-fail", false, "Return exit code 0 even if checks fail")
}
