package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	DebugLevel = iota
	InfoLevel
	NoticeLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

type LogLevel int

// Logger handles all logging operations
type Logger struct {
	verbose bool
	colors  map[string]*color.Color
	level   LogLevel
}

// NewLogger creates a new Logger instance
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		level:   DebugLevel,
		colors: map[string]*color.Color{
			"error":   color.New(color.FgRed),
			"warning": color.New(color.FgYellow),
			"success": color.New(color.FgGreen),
			"notice":  color.New(color.FgBlue), // Example color for notices
			"info":    color.New(color.FgCyan),
			"debug":   color.New(color.FgMagenta),
		},
	}
}

func (l *Logger) Log(level LogLevel, showPrefix bool, format string, args ...interface{}) {
	if l.level <= level {

		var prefix string
		switch level {
		case DebugLevel:
			prefix = l.colors["debug"].Sprint("DEBUG:")
		case InfoLevel:
			prefix = l.colors["info"].Sprint("INFO:")
		case NoticeLevel:
			prefix = l.colors["notice"].Sprint("NOTICE:")
		case WarningLevel:
			prefix = l.colors["warning"].Sprint("WARNING:")
		case ErrorLevel, FatalLevel:
			prefix = l.colors["error"].Sprint("ERROR:")
		}
		if showPrefix {
			fmt.Printf("%s %s\n", prefix, fmt.Sprintf(format, args...))
		} else {
			fmt.Printf("%s\n", fmt.Sprintf(format, args...))
		}
		//fmt.Printf("%s %s\n", prefix, fmt.Sprintf(format, args...))
	}
}

// LogError logs an error message
func (l *Logger) LogError(format string, args ...interface{}) {
	l.Log(ErrorLevel, true, format, args)
	//fmt.Printf("%s %s\n", l.colors["error"].Sprint("ERROR:"), fmt.Sprintf(format, args...))
}

func (l *Logger) LogDebug(format string, args ...interface{}) {
	l.Log(DebugLevel, true, format, args)
	//fmt.Printf("%s %s\n", l.colors["error"].Sprint("ERROR:"), fmt.Sprintf(format, args...))
}

// LogWarning logs a warning message
func (l *Logger) LogWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", l.colors["warning"].Sprint("WARNING:"), fmt.Sprintf(format, args...))
}

// LogSuccess logs a success message
func (l *Logger) LogSuccess(format string, args ...interface{}) {

	successPrefix := l.colors["success"].Sprint("✓ ")
	l.Log(NoticeLevel, false, successPrefix+format, args...)
}

// LogInfo logs an info message
func (l *Logger) LogInfo(format string, args ...interface{}) {
	//fmt.Printf("%s %s\n", l.colors["info"].Sprint("INFO:"), fmt.Sprintf(format, args...))
	l.Log(InfoLevel, true, format, args)
}

// LogQuestion formats a question message and returns it
func (l *Logger) LogQuestion(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", l.colors["question"].Sprint("?"), fmt.Sprintf(format, args...))
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
func (l *Logger) PrintStats(stats map[string]int, operation string) {
	if len(stats) == 0 {
		return
	}

	fmt.Println("\nSummary:")
	if stats["added"] > 0 {
		l.LogSuccess("%s license to %d files", operation, stats["added"])
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
