# License Manager

A CLI tool for managing license headers in source code files. This tool helps you add, remove, update, and check license headers across your codebase with ease and precision. It utilizes non-printing Unicode character to mark the start & end of the License comment such that you can easier update or remove the license information.

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/license-manager)](https://goreportcard.com/report/github.com/yourusername/license-manager)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- **Add License Headers**: Automatically add license headers to files that don't have them
- **Remove License Headers**: Clean up license headers from files
- **Update License Headers**: Replace existing license headers with new content
- **Check License Headers**: Verify if files have the correct license headers
- **Multiple File Support**: Process multiple files using glob patterns
- **Customizable Comment Styles**: Supports various programming language comment styles
- **Dry Run Mode**: Preview changes before applying them
- **Interactive Mode**: Confirm changes for each file
- **Verbose Output**: Detailed logging for better visibility
- **Skip Patterns**: Exclude specific files or directories
- **Preset Header Styles**: Choose from predefined header/footer styles


### Fully Tested Languages
This is a list of languages that have been fully tested.

- [ ] C
- [ ] C++
- [ ] C#
- [ ] CSS (Cascading Style Sheets)
- [x] Go
- [ ] HTML (HyperText Markup Language)
- [ ] Java
- [ ] JavaScript
- [ ] Lua
- [ ] Perl
- [ ] PHP (Hypertext Preprocessor)
- [ ] Python
- [ ] R
- [ ] Ruby
- [ ] Rust
- [ ] SASS (Syntactically Awesome Style Sheets)
- [ ] SCSS (Sassy CSS)
- [ ] Shell Scripts (.sh/.bash/.zsh)
- [ ] Swift
- [ ] TypeScript
- [ ] XML (Extensible Markup Language)
- [ ] YAML (YAML Ain't Markup Language)

## Installation

### Using Go

```bash
go install github.com/yourusername/license-manager@latest
```

### From Source

```bash
git clone https://github.com/yourusername/license-manager.git
cd license-manager
go build
```

## Usage

### Basic Commands

```bash
# Add license headers
license-manager add --license LICENSE.txt --input "**/*.go"

# Remove license headers
license-manager remove --input "**/*.go"

# Update existing license headers
license-manager update --license NEW_LICENSE.txt --input "**/*.go"

# Check license headers
license-manager check --license LICENSE.txt --input "**/*.go"
```

### Command Options

- `--license`: Path to the license template file (required for add/update/check)
- `--input`: Glob pattern for input files (e.g., "**/*.go" for all Go files)
- `--skip`: Glob pattern for files to skip
- `--prompt`: Enable interactive mode to confirm each change
- `--dry-run`: Preview changes without applying them
- `--verbose`: Enable detailed logging
- `--preset-style`: Choose a predefined header/footer style

### Examples

```bash
# Add license headers to all Go files, excluding tests
license-manager add --license LICENSE.txt --input "**/*.go" --skip "**/*_test.go"

# Update license headers in Python files with confirmation
license-manager update --license NEW_LICENSE.txt --input "**/*.py" --prompt

# Check license headers in JavaScript files with detailed output
license-manager check --license LICENSE.txt --input "**/*.js" --verbose

# Remove license headers from C++ files in dry-run mode
license-manager remove --input "**/*.cpp" --dry-run
```

## Configuration

### Comment Styles

The tool automatically detects appropriate comment styles based on file extensions:
- Go: `// comment`
- Python: `# comment`
- JavaScript/TypeScript: `// comment`
- C/C++: `/* comment */`
- And many more...

### Header Styles

Choose from various preset header/footer styles using the `--preset-style` flag:
- `standard`: Simple comment block
- `box`: Boxed comment style
- `line`: Line-separated comments
- And more...

## Building from Source

Requirements:
- Go 1.19 or higher

```bash
# Clone the repository
git clone https://github.com/yourusername/license-manager.git

# Change to project directory
cd license-manager

# Build the project
go build

# Run tests
go test ./...

# Install locally
go install
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI interface
- Inspired by various license management tools in the open-source community

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
