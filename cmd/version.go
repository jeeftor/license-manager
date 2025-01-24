package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Built:      %s\n", date)
		fmt.Printf("Git commit: %s\n", commit)
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Go version: %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
