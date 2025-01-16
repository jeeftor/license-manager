package styles

// CommentStyle represents how comments should be formatted for a specific language
type CommentStyle struct {
	Language    string
	Single      string
	MultiStart  string
	MultiEnd    string
	MultiPrefix string
	LinePrefix  string
	PreferMulti bool
}

// Common comment styles for different file extensions
var extensionStyles = map[string]CommentStyle{
	".rb":    {Language: "ruby", Single: "#", MultiStart: "=begin", MultiEnd: "=end", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".js":    {Language: "javascript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".jsx":   {Language: "javascript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".ts":    {Language: "typescript", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".tsx":   {Language: "typescript", Single: "//", MultiStart: "{/*", MultiEnd: "*/}", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".py":    {Language: "python", Single: "#", MultiStart: "'''", MultiEnd: "'''", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
	".go":    {Language: "go", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".java":  {Language: "java", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".cpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".c":     {Language: "c", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".h":     {Language: "c", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".hpp":   {Language: "cpp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".cs":    {Language: "csharp", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".php":   {Language: "php", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".swift": {Language: "swift", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".rs":    {Language: "rust", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".kt":    {Language: "kotlin", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".scala": {Language: "scala", Single: "//", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".html":  {Language: "html", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
	".css":   {Language: "css", Single: "", MultiStart: "/*", MultiEnd: "*/", MultiPrefix: " *", LinePrefix: " ", PreferMulti: true},
	".sh":    {Language: "shell", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".bash":  {Language: "shell", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".zsh":   {Language: "shell", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".fish":  {Language: "shell", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yaml":  {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".yml":   {Language: "yaml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".toml":  {Language: "toml", Single: "#", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".ini":   {Language: "ini", Single: ";", MultiStart: "", MultiEnd: "", MultiPrefix: "", LinePrefix: " ", PreferMulti: false},
	".xml":   {Language: "xml", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
	".md":    {Language: "markdown", Single: "", MultiStart: "<!--", MultiEnd: "-->", MultiPrefix: "", LinePrefix: " ", PreferMulti: true},
}

// GetCommentStyle returns the appropriate comment style for a given file extension
func GetCommentStyle(extension string) CommentStyle {
	if style, ok := extensionStyles[extension]; ok {
		return style
	}

	// Default to no comments for unknown file types
	return CommentStyle{
		Language: "text",
	}
}
