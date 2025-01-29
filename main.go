// main.go
package main

import (
	"fmt"
	"github.com/jeeftor/license-manager/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if exitErr, ok := err.(*cmd.ExitError); ok {
			os.Exit(exitErr.Code)
		}
		os.Exit(1)
	}
}
