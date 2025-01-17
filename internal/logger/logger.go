package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Logger handles all logging operations
type Logger struct {
	verbose bool
	colors  map[string]*color.Color
}

// NewLogger creates a new Logger instance
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		colors: map[string]*color.Color{
			"error":   color.New(color.FgRed),
			"warning": color.New(color.FgYellow),
			"success": color.New(color.FgGreen),
			"info":    color.New(color.FgCyan),
		},
	}
}

// LogError logs an error message
func (l *Logger) LogError(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", l.colors["error"].Sprint("ERROR:"), fmt.Sprintf(format, args...))
}

// LogWarning logs a warning message
func (l *Logger) LogWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", l.colors["warning"].Sprint("WARNING:"), fmt.Sprintf(format, args...))
}

// LogSuccess logs a success message
func (l *Logger) LogSuccess(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", l.colors["success"].Sprint("✓"), fmt.Sprintf(format, args...))
}

// LogInfo logs an info message
func (l *Logger) LogInfo(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", l.colors["info"].Sprint("INFO:"), fmt.Sprintf(format, args...))
}

// LogVerbose logs a message only in verbose mode
func (l *Logger) LogVerbose(format string, args ...interface{}) {
	if l.verbose {
		l.LogInfo(format, args...)
	}
}

// Prompt asks the user for confirmation
func (l *Logger) Prompt(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", message)
	
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// PrintStats prints operation statistics
func (l *Logger) PrintStats(stats map[string]int) {
	if len(stats) == 0 {
		return
	}

	fmt.Println("\nSummary:")
	if stats["added"] > 0 {
		l.LogSuccess("Added license to %d files", stats["added"])
	}
	if stats["existing"] > 0 {
		l.LogWarning("License already exists in %d files (use 'update' command to modify)", stats["existing"])
	}
	if stats["skipped"] > 0 {
		l.LogInfo("Skipped %d files", stats["skipped"])
	}
	if stats["failed"] > 0 {
		l.LogError("Failed to process %d files", stats["failed"])
	}
}
