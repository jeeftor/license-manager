package processor

import (
	"path/filepath"
)

type CommentStyle struct {
	Single      string // Single line comment prefix
	MultiStart  string // Multi-line comment start
	MultiEnd    string // Multi-line comment end
	PreferMulti bool   // Whether to prefer multi-line comments
	FileType    string // Type of file for special handling (e.g., "go", "python", "shell")
}

var extensionStyles = map[string]CommentStyle{
	".py":    {Single: "#", MultiStart: `"""`, MultiEnd: `"""`, PreferMulti: true, FileType: "python"},
	".rb":    {Single: "#", MultiStart: "=begin", MultiEnd: "=end", PreferMulti: false, FileType: "ruby"},
	".js":    {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
	".jsx":   {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "javascript"},
	".ts":    {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "typescript"},
	".tsx":   {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "typescript"},
	".java":  {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "java"},
	".go":    {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "go"},
	".c":     {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "c"},
	".cpp":   {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
	".hpp":   {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "cpp"},
	".cs":    {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "csharp"},
	".php":   {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "php"},
	".swift": {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "swift"},
	".rs":    {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "rust"},
	".sh":    {Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", PreferMulti: false, FileType: "shell"},
	".bash":  {Single: "#", MultiStart: ": <<'END'", MultiEnd: "END", PreferMulti: false, FileType: "shell"},
	".yml":   {Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
	".yaml":  {Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "yaml"},
	".pl":    {Single: "#", MultiStart: "=pod", MultiEnd: "=cut", PreferMulti: false, FileType: "perl"},
	".pm":    {Single: "#", MultiStart: "=pod", MultiEnd: "=cut", PreferMulti: false, FileType: "perl"},
	".r":     {Single: "#", MultiStart: "", MultiEnd: "", PreferMulti: false, FileType: "r"},
	".html":  {Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "html"},
	".xml":   {Single: "", MultiStart: "<!--", MultiEnd: "-->", PreferMulti: true, FileType: "xml"},
	".css":   {Single: "", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "css"},
	".scss":  {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "scss"},
	".sass":  {Single: "//", MultiStart: "/*", MultiEnd: "*/", PreferMulti: true, FileType: "sass"},
	".lua":   {Single: "--", MultiStart: "--[[", MultiEnd: "--]]", PreferMulti: true, FileType: "lua"},
}

func getCommentStyle(filename string) CommentStyle {
	ext := filepath.Ext(filename)
	if style, ok := extensionStyles[ext]; ok {
		return style
	}
	// Default to C-style comments if unknown
	return CommentStyle{
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		PreferMulti: true,
		FileType:    "unknown",
	}
}
