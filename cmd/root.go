// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jeeftor/license-manager/internal/force"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cc "github.com/ivanpirog/coloredcobra"
)

const (
	envPrefix = "LM"
)

type commentStyleFlag struct {
	value *force.ForceCommentStyle
}

func (f *commentStyleFlag) String() string {
	if f.value == nil {
		return string(force.No) // default value
	}
	return string(*f.value)
}

func (f *commentStyleFlag) Set(s string) error {
	switch force.ForceCommentStyle(s) {
	case force.No, force.Single, force.Multi:
		*f.value = force.ForceCommentStyle(s)
		return nil
	default:
		return fmt.Errorf("must be one of no, single, or multi")
	}
}

func (f *commentStyleFlag) Type() string {
	return "commentStyle"
}

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

	// buildVersion color
	versionColor = color.RGB(85, 85, 85) // 242 gray
)

var versionString = GetVersionString()

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

var rootCmd = &cobra.Command{
	Use:   "license-manager",
	Short: color.CyanString("A tool to manage license headers in source files"),
	Long: logo + "\n\n" + color.BlueString(
		"license-manager",
	) + color.WhiteString(
		`is a CLI tool that helps manage license headers in source files.
It can add, remove, update, and check license headers in multiple files using patterns.

`,
	) + color.YellowString(
		"Environment variables:",
	) + `
  ` + color.CyanString(
		"LM_LICENSE",
	) + `   Path to license text file
`,
}

func Execute() error {
	cobra.MousetrapHelpText = ""

	// Configure help template
	helpTemplate := `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}{{end}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

	rootCmd.SetHelpTemplate(helpTemplate)

	// Configure help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long:  `Help provides help for any command in the application.`,
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				c.Root().Usage()
			} else {
				cmd.Help()
			}
		},
	})

	// Configure help flags
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for this command")

	// Silence cobra error output since we handle it ourselves
	rootCmd.SilenceErrors = true

	// Configure output colors
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiYellow + cc.Bold + cc.Underline,
		Commands:      cc.HiBlue + cc.Bold,
		Example:       cc.Italic,
		ExecName:      cc.Bold + cc.Red,
		CmdShortDescr: cc.Green,
		Flags:         cc.Bold + cc.Green,
	})

	// Configure error handling
	err := rootCmd.Execute()
	if err != nil {
		if exitErr, ok := err.(*ExitError); ok {
			os.Exit(exitErr.Code)
		}
		log := logger.NewLogger(logger.ParseLogLevel(cfgLogLevel))
		log.LogError("Error: %v", err)
		os.Exit(1)
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgPresetStyle, "style", "hash", "Preset style for header/footer (run styles command for list)")
	rootCmd.PersistentFlags().Var(&commentStyleFlag{&cfgForceCommentStyle}, "comments",
		"Force comment style (no|single|multi)")
	// set default value
	cfgForceCommentStyle = force.No

	rootCmd.PersistentFlags().StringVar(&cfgLicense, "license", "", "Path to license text file")

	rootCmd.PersistentFlags().
		StringSliceVar(&cfgInputs, "input", []string{}, "Inputs file patterns")
	rootCmd.PersistentFlags().StringSliceVar(&cfgSkips, "skip", []string{}, "Patterns to skip")

	rootCmd.PersistentFlags().
		StringVar(&cfgLogLevel, "log-level", "notice", "Log level (debug, info, notice, warn, error)")
}

func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Update variables with viper values
	if viper.IsSet("license") {
		cfgLicense = viper.GetString("license")
	}

	if viper.IsSet("log-level") {
		cfgLogLevel = viper.GetString("log-level")
	}
	if viper.IsSet("style") {
		cfgPresetStyle = viper.GetString("style")
	}
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
