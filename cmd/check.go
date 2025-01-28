package cmd

import (
	"fmt"
	"license-manager/internal/logger"
	"strings"

	"github.com/spf13/cobra"

	"license-manager/internal/config"
	"license-manager/internal/license"
	"license-manager/internal/processor"

	cc "github.com/ivanpirog/coloredcobra"
)

var (
	checkIgnoreFail bool
	cfgLicense      string
	cfgInputs       []string
	cfgSkips        []string
	cfgPresetStyle  string
	cfgLogLevel     string
	cfgPreferMulti  bool
)

// ExitError represents an error with an exit code
type ExitError struct {
	msg  string
	Code int
}

func (e *ExitError) Error() string {
	return e.msg
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check license headers in files",
	Long: `Check license headers in files

Exit Codes:
  0: All files match
  1: At least 1 file is missing a license
  2: Files have both content and header mismatch
  3: Files have content mismatch
  4: Files have style mismatch
`,
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

		// Create app config
		appCfg := config.AppConfig{
			LicenseFile: cfgLicense,
			Inputs:      strings.Join(cfgInputs, ","),
			Skips:       strings.Join(cfgSkips, ","),
			HeaderStyle: cfgPresetStyle,
			PreferMulti: cfgPreferMulti,
			LogLevel:    logger.ParseLogLevel(cfgLogLevel),
			IgnoreFail:  checkIgnoreFail,
		}

		cc.Init(&cc.Config{
			RootCmd:       rootCmd,
			Headings:      cc.HiYellow + cc.Bold + cc.Underline,
			Commands:      cc.HiBlue + cc.Bold,
			Example:       cc.Italic,
			ExecName:      cc.Bold + cc.Red,
			CmdShortDescr: cc.Green,
			Flags:         cc.Bold + cc.Green,
		})

		// Convert to processor config
		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return fmt.Errorf("failed to create processor config: %w", err)
		}

		// Create processor and run check
		p := processor.NewFileProcessor(procCfg)
		err = p.Check()

		if err != nil {
			if checkErr, ok := err.(*processor.CheckError); ok {
				if checkIgnoreFail {
					return nil
				}
				switch checkErr.Status {
				case license.NoLicense:
					return &ExitError{
						msg:  "license check failed: some files have missing licenses",
						Code: int(license.NoLicense),
					}
				case license.ContentMismatch:
					return &ExitError{
						msg:  "license check failed: some files have incorrect license content",
						Code: int(license.ContentMismatch),
					}
				case license.StyleMismatch:
					return &ExitError{
						msg:  "license check failed: some files have incorrect license style",
						Code: int(license.StyleMismatch),
					}
				case license.ContentAndStyleMismatch:
					return &ExitError{
						msg:  "license check failed: some files have incorrect license content and style",
						Code: int(license.ContentAndStyleMismatch),
					}
				default:
					return &ExitError{
						msg:  "license check failed: unknown error",
						Code: 5, // Keep this as a constant since it's not part of the Status enum
					}
				}
			}

			return &ExitError{
				msg:  fmt.Sprintf("license check failed: %v", err),
				Code: 5, // Keep this as a constant since it's not part of the Status enum
			}
		}
		//if cfgVerbose {
		//fmt.Printf("Exit code: 0\n")
		//}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&checkIgnoreFail, "ignore-fail", false, "Return exit code 0 even if checks fail")
}
