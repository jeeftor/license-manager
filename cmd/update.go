package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jeeftor/license-manager/internal/config"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/processor"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update license headers in files",
	Long:  `Update existing license headers in files with new content`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgLicense == "" {
			return fmt.Errorf("license file (--license) is required for update command")
		}

		appCfg := config.AppConfig{
			// File paths
			LicenseFile: cfgLicense,
			Inputs:      ProcessPatterns(cfgInputs),
			Skips:       ProcessPatterns(cfgSkips),

			// Style settings
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go", // default

			// Behavior flags
			LogLevel: logger.ParseLogLevel(cfgLogLevel),

			Force:      false,
			IgnoreFail: false,
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		err = p.Update()

		cmd.SilenceUsage = true
		return err
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
