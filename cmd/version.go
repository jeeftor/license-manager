package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
)

var shortOutput bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if shortOutput {
			fmt.Println(buildVersion)
			return
		}
		versionColor := color.New(color.FgCyan, color.Bold)
		buildColor := color.New(color.FgYellow)
		commitColor := color.New(color.FgGreen)
		osArchColor := color.New(color.FgMagenta)
		goVersionColor := color.New(color.FgRed)
		whiteColor := color.New(color.FgWhite)
		pathColor := color.New(color.FgBlue)

		whiteColor.Printf("Version: ")
		versionColor.Printf("%s\n", buildVersion)

		whiteColor.Printf("Built:   ")
		buildColor.Printf("%s\n", GetFormattedBuildTime())

		whiteColor.Printf("Commit:  ")
		commitColor.Printf("%s\n", buildCommit)

		whiteColor.Printf("OS/Arch: ")
		osArchColor.Printf("%s/%s\n", runtime.GOOS, runtime.GOARCH)

		whiteColor.Printf("Go:      ")
		goVersionColor.Printf("%s\n", runtime.Version())

		exe, err := os.Executable()
		exePath := "Unknown"
		if err == nil {
			exePath, _ = filepath.Abs(exe)
		}

		whiteColor.Printf("Binary:  ")
		pathColor.Printf("%s\n", exePath)

	},
}

func init() {
	versionCmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "Print only version number")
	rootCmd.AddCommand(versionCmd)
}
