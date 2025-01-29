package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/jeeftor/license-manager/internal/styles"
)

var stylesCmd = &cobra.Command{
	Use:   "styles",
	Short: "List available license header styles",
	Long:  `Display all available preset styles for license headers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.CyanString("Available header styles:"))
		fmt.Println()

		for i, name := range styles.List() {
			style := styles.Get(name)
			fmt.Printf("Style [%0d]: %s\n", i+1, color.BlueString(style.Name))
			fmt.Printf("Description: %s\n", color.WhiteString(style.Description))
			fmt.Printf("Header: %s\n", color.GreenString(style.Header))
			fmt.Printf("Footer: %s\n", color.GreenString(style.Footer))
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(stylesCmd)
}
