package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/license"
	"license-manager/internal/processor"
	"os"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for license headers in files",
	Long:  `Check if files have the specified license headers\n Exit Code 0 - License Found\nExit Code 1 - License Missing \nExit Code 2 - License Mismatch`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// CLI validation errors should show usage
		if cfgLicense == "" {
			return fmt.Errorf("license file (--license) is required for check command")
		}
		if cfgInputs == nil {
			return fmt.Errorf("input pattern (--input) is required for check command")
		}

		// After validation passes, silence usage since any further errors are execution errors
		cmd.SilenceUsage = true

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
			IgnoreFail:  checkIgnoreFail,
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		err = p.Check()

		if licErr, ok := err.(*processor.CheckError); ok {
			switch licErr.Status {
			case license.NoLicense:
				os.Exit(1)
			case license.DifferentLicense:
				os.Exit(2)
			}
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&checkIgnoreFail, "ignore-fail", false, "Return exit code 0 even if checks fail")
}
