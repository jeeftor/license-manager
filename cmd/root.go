// cmd/root.go
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

import cc "github.com/ivanpirog/coloredcobra"

const (
	envPrefix = "LM"
)

// Version information
var (
	version = "dev" // Will be overwritten at build time
	date    = func() string {
		if version == "dev" {
			return time.Now().Format("2006-01-02")
		}
		return "unknown"
	}()
	commit = ""
)

var (
	// Blue gradient (51, 45, 39, 33)
	blue1 = color.RGB(87, 207, 255) // 51
	blue2 = color.RGB(77, 183, 255) // 45
	blue3 = color.RGB(62, 158, 255) // 39
	blue4 = color.RGB(43, 134, 255) // 33

	// Purple gradient (183, 177, 171, 165)
	purple1 = color.RGB(215, 183, 255) // 183
	purple2 = color.RGB(207, 171, 255) // 177
	purple3 = color.RGB(199, 159, 255) // 171
	purple4 = color.RGB(191, 147, 255) // 165

	// Version color
	versionColor = color.RGB(85, 85, 85) // 242 gray
)

// Build the version string
var versionString = version + " (" + date + " " + commit + ")"

//var logo = "\x1b[38;5;51m" + `
//▗▖   ▗▄▄▄▖ ▗▄▄▖▗▄▄▄▖▗▖  ▗▖ ▗▄▄▖▗▄▄▄▖ ` + "\x1b[38;5;45m" + `
//▐▌     █  ▐▌   ▐▌   ▐▛▚▖▐▌▐▌   ▐▌    ` + "\x1b[38;5;39m" + `
//▐▌     █  ▐▌   ▐▛▀▀▘▐▌ ▝▜▌ ▝▀▚▖▐▛▀▀▘ ` + "\x1b[38;5;33m" + `
//▐▙▄▄▖▗▄█▄▖▝▚▄▄▖▐▙▄▄▖▐▌  ▐▌▗▄▄▞▘▐▙▄▄▖` + "\x1b[0m" + "\n" +
//	"\x1b[38;5;183m" + `   ▗▖  ▗▖ ▗▄▖ ▗▖  ▗▖ ▗▄▖  ▗▄▄▖▗▄▄▄▖▗▄▄▖ ` + "\x1b[38;5;177m" + `
//   ▐▛▚▞▜▌▐▌ ▐▌▐▛▚▖▐▌▐▌ ▐▌▐▌   ▐▌   ▐▌ ▐▌` + "\x1b[38;5;171m" + `
//   ▐▌  ▐▌▐▛▀▜▌▐▌ ▝▜▌▐▛▀▜▌▐▌▝▜▌▐▛▀▀▘▐▛▀▚▖` + "\x1b[38;5;165m" + `
//   ▐▌  ▐▌▐▌ ▐▌▐▌  ▐▌▐▌ ▐▌▝▚▄▞▘▐▙▄▄▖▐▌ ▐▌` + "\x1b[0m" + "\n" +
//	"\x1b[38;5;242m" + versionString + "\x1b[0m"

var logo = "" +
	blue1.Sprint(`▗▖   ▗▄▄▄▖ ▗▄▄▖▗▄▄▄▖▗▖  ▗▖ ▗▄▄▖▗▄▄▄▖`) + "\n" +
	blue2.Sprint(`▐▌     █  ▐▌   ▐▌   ▐▛▚▖▐▌▐▌   ▐▌`) + "\n" +
	blue3.Sprint(`▐▌     █  ▐▌   ▐▛▀▀▘▐▌ ▝▜▌ ▝▀▚▖▐▛▀▀▘`) + "\n" +
	blue4.Sprint(`▐▙▄▄▖▗▄█▄▖▝▚▄▄▖▐▙▄▄▖▐▌  ▐▌▗▄▄▞▘▐▙▄▄▖`) + "\n" +
	purple1.Sprint(`   ▗▖  ▗▖ ▗▄▖ ▗▖  ▗▖ ▗▄▖  ▗▄▄▖▗▄▄▄▖▗▄▄▖`) + "\n" +
	purple2.Sprint(`   ▐▛▚▞▜▌▐▌ ▐▌▐▛▚▖▐▌▐▌ ▐▌▐▌   ▐▌   ▐▌ ▐▌`) + "\n" +
	purple3.Sprint(`   ▐▌  ▐▌▐▛▀▜▌▐▌ ▝▜▌▐▛▀▜▌▐▌▝▜▌▐▛▀▀▘▐▛▀▚▖`) + "\n" +
	purple4.Sprint(`   ▐▌  ▐▌▐▌ ▐▌▐▌  ▐▌▐▌ ▐▌▝▚▄▞▘▐▙▄▄▖▐▌ ▐▌`) + "\n" +
	versionColor.Sprint(versionString)

var (
	cfgLicense      string
	cfgInputs       []string
	cfgSkips        []string
	cfgPrompt       bool
	cfgDryRun       bool
	cfgVerbose      bool   // Add verbose flag
	cfgPresetStyle  string // header/footer style
	cfgPreferMulti  bool   // prefer multi-line comments where supported
	checkIgnoreFail bool   // Added for check command

)

var rootCmd = &cobra.Command{
	Use:   "license-manager",
	Short: color.CyanString("A tool to manage license headers in source files"),
	Long: logo + "\n\n" + color.BlueString("license-manager") + color.WhiteString(`is a CLI tool that helps manage license headers in source files.
It can add, remove, update, and check license headers in multiple files using patterns.

`) + color.YellowString("Environment variables:") + `
  ` + color.CyanString("LM_LICENSE") + `   Path to license text file
`,
}

func Execute() error {
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiYellow + cc.Bold + cc.Underline,
		Commands:      cc.HiBlue + cc.Bold,
		Example:       cc.Italic,
		ExecName:      cc.Bold + cc.Red,
		CmdShortDescr: cc.Green,
		//FlagsDescr:    cc.Green,
		Flags: cc.Bold + cc.Green,
	})
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPresetStyle, "style", "simple", "Preset style for header/footer (e.g., simple, modern, elegant)")
	rootCmd.PersistentFlags().BoolVar(&cfgPreferMulti, "multi", true, "Prefer multi-line comments where supported")

	rootCmd.PersistentFlags().StringVar(&cfgLicense, "license", "", "Path to license text file")

	rootCmd.PersistentFlags().StringSliceVar(&cfgInputs, "input", []string{}, "Inputs file patterns")
	rootCmd.PersistentFlags().StringSliceVar(&cfgSkips, "skip", []string{}, "Patterns to skip")

	rootCmd.PersistentFlags().BoolVar(&cfgPrompt, "prompt", false, "Prompt before processing each file")
	rootCmd.PersistentFlags().BoolVar(&cfgDryRun, "dry-run", false, "Show which files would be processed without making changes")
	rootCmd.PersistentFlags().BoolVar(&cfgVerbose, "verbose", false, "Enable verbose output")
}
func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Update variables with viper values
	if viper.IsSet("license") {
		cfgLicense = viper.GetString("license")
	}
	//if viper.IsSet("style") {
	//	cfgPresetStyle = viper.GetString("style")
	//}
	//if viper.IsSet("multi") {
	//	cfgPreferMulti = viper.GetBool("multi")
	//}
	//if viper.IsSet("prompt") {
	//	cfgPrompt = viper.GetBool("prompt")
	//}
	//if viper.IsSet("dry-run") {
	//	cfgDryRun = viper.GetBool("dry-run")
	//}
	//if viper.IsSet("verbose") {
	//	cfgVerbose = viper.GetBool("verbose")
	//}
}

func ProcessPatterns(patterns []string) string {
	var result []string
	for _, p := range patterns {
		// Split on commas if present
		parts := strings.Split(p, ",")
		result = append(result, parts...)
	}
	return strings.Join(result, ",")
}
