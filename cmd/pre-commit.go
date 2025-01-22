package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"license-manager/internal/config"
	"license-manager/internal/processor"

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
		stagedFiles, err := getStagedGoFiles()
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
			PreferMulti:  cfgPreferMulti,

			Verbose:     cfgVerbose,
			Interactive: false, // Typically want non-interactive in pre-commit
			DryRun:      false,
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

// getStagedGoFiles retrieves the list of staged Go files
func getStagedGoFiles() ([]string, error) {
	// Use git command to get staged Go files
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split and filter for Go files
	var stagedGoFiles []string
	for _, file := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(file) != "" && strings.HasSuffix(file, ".go") {
			stagedGoFiles = append(stagedGoFiles, file)
		}
	}

	return stagedGoFiles, nil
}

func init() {
	rootCmd.AddCommand(preCommitCmd)

	// Reuse existing flags from check command, but with pre-commit specific defaults
	preCommitCmd.Flags().StringVar(&cfgLicense, "license", "./LICENSE", "Path to license file")
	preCommitCmd.Flags().BoolVar(&cfgVerbose, "verbose", false, "Enable verbose output")
}
