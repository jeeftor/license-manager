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

func ParseLogLevel(level string) LogLevel {

	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "notice":
		return NoticeLevel
	case "warn":
		return WarningLevel
	case "warning":
		return WarningLevel
	default:
		return ErrorLevel
	}
}

type LogLevel int

// Logger handles all logging operations
type Logger struct {
	colors map[string]*color.Color
	level  LogLevel
}

// NewLogger creates a new Logger instance
func NewLogger(LogLevel LogLevel) *Logger {
	return &Logger{
		level: LogLevel,
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
		var outputText string
		if len(args) > 0 {
			// Always use fmt.Sprintf for consistency with direct Printf
			outputText = fmt.Sprintf(format, args...)
		} else {
			// When no args are provided and format contains %%, we need to handle it specially
			// to match fmt.Printf behavior
			if strings.Contains(format, "%%") {
				outputText = fmt.Sprintf(format)
			} else {
				outputText = format
			}
		}

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
			fmt.Printf("%s %s\n", prefix, outputText)
		} else {
			fmt.Printf("%s\n", outputText)
		}
	}
}

// LogError logs an error message
func (l *Logger) LogError(format string, args ...interface{}) {
	l.Log(ErrorLevel, true, format, args...)
}

func (l *Logger) LogDebug(format string, args ...interface{}) {
	l.Log(DebugLevel, true, format, args...)
}

// LogWarning logs a warning message
func (l *Logger) LogWarning(format string, args ...interface{}) {
	l.Log(WarningLevel, true, format, args...)
}

// LogSuccess logs a success message
func (l *Logger) LogSuccess(format string, args ...interface{}) {
	successPrefix := l.colors["success"].Sprint("✓ ")
	l.Log(NoticeLevel, false, successPrefix+format, args...)
}

// LogInfo logs an info message
func (l *Logger) LogInfo(format string, args ...interface{}) {
	l.Log(InfoLevel, true, format, args...)
}

// LogInfo logs an notice message
func (l *Logger) LogNotice(format string, args ...interface{}) {
	l.Log(NoticeLevel, true, format, args...)
}

// LogQuestion formats a question message and returns it
func (l *Logger) LogQuestion(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", l.colors["question"].Sprint("?"), fmt.Sprintf(format, args...))
}

// LogVerbose logs a message only in verbose mode
func (l *Logger) LogVerbose(format string, args ...interface{}) {
	l.Log(DebugLevel, true, format, args...)
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
		fmt.Printf("%s license to %d files\n", operation, stats["added"])
	}
	if stats["existing"] > 0 {
		fmt.Printf(
			"License already exists in %d files (use 'update' command to modify)\n",
			stats["existing"],
		)
	}
	if stats["skipped"] > 0 {
		fmt.Printf("Skipped %d files\n", stats["skipped"])
	}
	if stats["failed"] > 0 {
		fmt.Printf("Failed to process %d files\n", stats["failed"])
	}
}
