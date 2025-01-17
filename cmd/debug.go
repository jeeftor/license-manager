package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"license-manager/internal/config"
	"license-manager/internal/styles"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug license markers in files",
	Long:  `Show license markers in files by making invisible markers visible`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgInput == "" {
			return fmt.Errorf("input file (--input) is required for debug command")
		}

		appCfg := config.AppConfig{
			// File paths
			Input: cfgInput,

			// Style settings
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go",
			PreferMulti:  cfgPreferMulti,

			// Debug needs verbose output
			Verbose: true,
		}

		// Read the input file
		content, err := os.ReadFile(appCfg.Input)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		// Get the header/footer style for debugging
		style := styles.Get(appCfg.HeaderStyle)

		fmt.Println("File contents with markers made visible:")
		fmt.Println("---------------------------------------")

		// Show the configured style's header/footer with markers visible
		if style.Name != "" {
			fmt.Printf("Style: %s\n", style.Name)
			fmt.Printf("Header should be: %s\n", strings.ReplaceAll(strings.ReplaceAll(
				style.Header, "\u200B", color.New(color.FgRed).Sprint("[START]")),
				"\u200C", color.New(color.FgRed).Sprint("[END]")))
			fmt.Printf("Footer should be: %s\n", strings.ReplaceAll(strings.ReplaceAll(
				style.Footer, "\u200B", color.New(color.FgRed).Sprint("[START]")),
				"\u200C", color.New(color.FgRed).Sprint("[END]")))
			fmt.Println()
		}

		// Show the actual file contents with markers visible
		fmt.Println("Actual file contents:")
		fmt.Println("--------------------")
		debugContent := strings.ReplaceAll(string(content), "\u200B", color.New(color.FgRed).Sprint("[START]"))
		debugContent = strings.ReplaceAll(debugContent, "\u200C", color.New(color.FgRed).Sprint("[END]"))
		fmt.Println(debugContent)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
