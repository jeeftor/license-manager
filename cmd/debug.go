package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"license-manager/internal/processor"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug license markers in files",
	Long:  `Show license markers in files by making invisible markers visible`,
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(cfgInput)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		// Get the current style to check its header/footer
		style := processor.GetPresetStyle(cfgPresetStyle)
		
		fmt.Println("File contents with markers made visible:")
		fmt.Println("---------------------------------------")
		
		// Show the raw header/footer with markers visible
		fmt.Printf("Header should be: %s\n", strings.ReplaceAll(strings.ReplaceAll(
			style.Header, "\u200B", color.New(color.FgRed).Sprint("[START]")), 
			"\u200C", color.New(color.FgRed).Sprint("[END]")))
		fmt.Printf("Footer should be: %s\n", strings.ReplaceAll(strings.ReplaceAll(
			style.Footer, "\u200B", color.New(color.FgRed).Sprint("[START]")), 
			"\u200C", color.New(color.FgRed).Sprint("[END]")))
		fmt.Println("\nActual file contents:")
		fmt.Println("--------------------")
		
		// Show the file contents with markers visible
		debugContent := strings.ReplaceAll(string(content), "\u200B", color.New(color.FgRed).Sprint("[START]"))
		debugContent = strings.ReplaceAll(debugContent, "\u200C", color.New(color.FgRed).Sprint("[END]"))
		fmt.Println(debugContent)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
