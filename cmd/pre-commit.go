package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeeftor/license-manager/internal/config"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/processor"
	"github.com/spf13/cobra"
)

var (
	licensePath string
	logLevel    string
	addFlag     bool
	updateFlag  bool
)
var preCommitCmd = &cobra.Command{
	Use:   "pre-commit [files...]",
	Short: "Run license checks on specified files",
	Long:  `Automatically checks license headers on files passed by pre-commit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// Check if license file exists
		if _, err := os.Stat(licensePath); os.IsNotExist(err) {
			return fmt.Errorf(`License file not found at: %s

To specify a different license file, update your .pre-commit-config.yaml:

repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: vX.X.X
    hooks:
      - id: license-manager
        args: [--license, path/to/your/LICENSE]`, licensePath)
		}

		// Use files passed directly as arguments
		if len(args) == 0 {
			fmt.Println("No files to check")
			return nil
		}

		//TODO: Print out all the files we received

		// Rest of your existing code...
		appCfg := config.AppConfig{
			LicenseFile:  licensePath,
			Inputs:       strings.Join(args, ","),
			Skips:        ProcessPatterns(cfgSkips),
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go", // default
			LogLevel:     logger.ParseLogLevel(logLevel),
			Interactive:  false,
			Force:        false,
			IgnoreFail:   false,
			IsPreCommit:  true,
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		if addFlag {
			err := p.Add()
			if err != nil {
				return err
			}
		}

		if updateFlag {
			err := p.Update()
			if err != nil {
				return err
			}
		}

		return p.Check()
	},
}

func init() {
	rootCmd.AddCommand(preCommitCmd)
	preCommitCmd.Flags().StringVar(&licensePath, "license", "./LICENSE", "Path to license file")
	preCommitCmd.Flags().StringVar(&logLevel, "log-level", "info", "Logging level (debug, info, warn, error)")

	preCommitCmd.Flags().BoolVar(&addFlag, "add", false, "Add license headers to files")
	preCommitCmd.Flags().BoolVar(&updateFlag, "update", false, "Update existing license headers")
}
