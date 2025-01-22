package main

import (
	"fmt"
	"io"
)

// TODO: Update to a strucutred logger like:"github.com/charmbracelet/log" or "log/slog"
// Logger is an interface for logging.
type Logger interface {
	Printf(format string, v ...interface{})
	Write(p []byte) (n int, err error)
	SetOutput(w io.Writer)
}

func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	fmt.Println(greet("World"))
}
