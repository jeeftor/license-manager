package styles

import "strings"

// CommentLanguage represents how comments should be formatted for a specific language
type CommentLanguage struct {
	Language    string
	Single      string
	MultiStart  string
	MultiEnd    string
	MultiPrefix string
	LinePrefix  string
	PreferMulti bool
}

// LanguageExtensions Common comment styles for different file extensions
var LanguageExtensions = map[string]CommentLanguage{
	".rb": {
		Language:    "ruby",
		Single:      "#",
		MultiStart:  "=begin",
		MultiEnd:    "=end",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".js": {
		Language:    "javascript",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	//".jsx": {Language: "javascript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".ts": {
		Language:    "typescript",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	//".tsx":   {Language: "typescript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".py": {
		Language:    "python",
		Single:      "#",
		MultiStart:  "'''",
		MultiEnd:    "'''",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".go": {
		Language:    "go",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".java": {
		Language:    "java",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".cpp": {
		Language:    "cpp",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".c": {
		Language:    "c",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".h": {
		Language:    "c",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".hpp": {
		Language:    "cpp",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".cs": {
		Language:    "csharp",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".php": {
		Language:    "php",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".swift": {
		Language:    "swift",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".rs": {
		Language:    "rust",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".kt": {
		Language:    "kotlin",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".scala": {
		Language:    "scala",
		Single:      "//",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".html": {
		Language:    "html",
		Single:      "",
		MultiStart:  "<!--",
		MultiEnd:    "-->",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".css": {
		Language:    "css",
		Single:      "",
		MultiStart:  "/*",
		MultiEnd:    "*/",
		MultiPrefix: " *",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	// Shell script extensions
	".sh": {
		Language:    "shell",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".bash": {
		Language:    "shell",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".zsh": {
		Language:    "shell",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".fish": {
		Language:    "shell",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".yaml": {
		Language:    "yaml",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".yml": {
		Language:    "yaml",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".toml": {
		Language:    "toml",
		Single:      "#",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".ini": {
		Language:    "ini",
		Single:      ";",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: false,
	},
	".xml": {
		Language:    "xml",
		Single:      "",
		MultiStart:  "<!--",
		MultiEnd:    "-->",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	".md": {
		Language:    "markdown",
		Single:      "",
		MultiStart:  "<!--",
		MultiEnd:    "-->",
		MultiPrefix: "",
		LinePrefix:  " ",
		PreferMulti: true,
	},
	"GENERIC": {
		Language:    "unset",
		Single:      "",
		MultiStart:  "",
		MultiEnd:    "",
		MultiPrefix: "",
		LinePrefix:  "",
		PreferMulti: false,
	},
}

// StripCommentMarkers removes comment markers from a line of text based on the language style
func (c *CommentLanguage) StripCommentMarkers(line string) string {
	if line == "" {
		return line
	}

	// Handle multi-line style markers
	if c.MultiStart != "" {
		line = strings.TrimPrefix(line, c.MultiStart)
		line = strings.TrimSuffix(line, c.MultiEnd)
	}

	// Handle single-line style markers
	if c.Single != "" {
		line = strings.TrimPrefix(line, c.Single)
	}

	// Handle line prefix
	if c.MultiPrefix != "" {
		line = strings.TrimPrefix(line, c.MultiPrefix)
	}
	if c.LinePrefix != "" {
		line = strings.TrimPrefix(line, c.LinePrefix)
	}

	return strings.TrimSpace(line)
}

// GetLanguageCommentStyle returns the appropriate comment style for a given file extension
func GetLanguageCommentStyle(extension string) CommentLanguage {
	// Add dot prefix if missing
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	if style, ok := LanguageExtensions[extension]; ok {
		return style
	}

	// Default to no comments for unknown file types
	return CommentLanguage{
		Language: "text",
	}
}
