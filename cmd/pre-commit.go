package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jeeftor/license-manager/internal/logger"

	"github.com/jeeftor/license-manager/internal/config"
	"github.com/jeeftor/license-manager/internal/processor"

	"github.com/spf13/cobra"
)

var preCommitCmd = &cobra.Command{
	Use:   "pre-commit",
	Short: "Run license checks on staged files for pre-commit",
	Long:  `Automatically checks license headers on staged Go files before committing`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Silence usage output if command succeeds
		cmd.SilenceUsage = true

		// Get staged Go files
		stagedFiles, err := getStagedFiles()
		if err != nil {
			return fmt.Errorf("failed to get staged files: %v", err)
		}

		// If no staged Go files, exit successfully
		if len(stagedFiles) == 0 {
			fmt.Println("No staged files to check")
			return nil
		}

		// Prepare configuration
		appCfg := config.AppConfig{
			LicenseFile: cfgLicense,
			Inputs:      strings.Join(stagedFiles, ","),
			Skips:       ProcessPatterns(cfgSkips),

			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go", // default

			LogLevel: logger.DebugLevel,

			Interactive: false, // Typically want non-interactive in pre-commit
			Force:       false,
			IgnoreFail:  false, // We want to fail if license checks fail
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		p := processor.NewFileProcessor(procCfg)
		return p.Check()
	},
}

// getStagedFiles retrieves the list of all staged files
func getStagedFiles() ([]string, error) {
	// Use git command to get staged files
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split into files, removing empty lines
	var stagedFiles []string
	for _, file := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(file) != "" {
			stagedFiles = append(stagedFiles, file)
		}
	}

	return stagedFiles, nil
}

func init() {
	rootCmd.AddCommand(preCommitCmd)

	// Reuse existing flags from check command, but with pre-commit specific defaults
	preCommitCmd.Flags().StringVar(&cfgLicense, "license", "./LICENSE", "Path to license file")
}
