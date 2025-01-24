package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var shortOutput bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if shortOutput {
			fmt.Println(version)
			return
		}
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Built:      %s\n", date)
		fmt.Printf("Git commit: %s\n", commit)
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Go version: %s\n", runtime.Version())
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "Print only version number")
	rootCmd.AddCommand(versionCmd)
}
