package cmd

import (
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
	"license-manager/internal/styles"
)

var checkIgnoreFail bool

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for license headers in files",
	Long:  `Check if files have the specified license headers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("%s", "license file (--license) is required for check command")
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
		err := p.Check()

		// If it's a CheckError, we don't want to show usage
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
