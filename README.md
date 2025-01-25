# License Manager

![GitHub commit activity](https://img.shields.io/github/commit-activity/w/jeeftor/license-manager)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/jeeftor/license-manager/latest)
![GitHub Release Date](https://img.shields.io/github/release-date/jeeftor/license-manager)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/jeeftor/license-manager/total)
![GitHub Repo stars](https://img.shields.io/github/stars/jeeftor/license-manager)

![docs/logo.png](docs/logo.png)
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
This is a list of languages that have been ~fully~ mostly tested.

![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=status.c.status&style=for-the-badge&label=c&color=status.c.color)

![Dynamic C Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.c.status&label=C&labelColor=%24.c.color)

![Dynamic C++ Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.cpp.status&label=C%2B%2B&labelColor=%24.cpp.color)

![Dynamic C# Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.csharp.status&label=C%23&labelColor=%24.csharp.color)

![Dynamic CSS Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.css.status&label=CSS&labelColor=%24.css.color)

![Dynamic Go Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.go.status&label=Go&labelColor=%24.go.color)

![Dynamic HTML Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.html.status&label=HTML&labelColor=%24.html.color)

![Dynamic INI Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.ini.status&label=INI&labelColor=%24.ini.color)

![Dynamic Java Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.java.status&label=Java&labelColor=%24.java.color)

![Dynamic JavaScript Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.javascript.status&label=JavaScript&labelColor=%24.javascript.color)

![Dynamic Kotlin Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.kotlin.status&label=Kotlin&labelColor=%24.kotlin.color)

![Dynamic Markdown Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.markdown.status&label=Markdown&labelColor=%24.markdown.color)

![Dynamic PHP Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=%24.php.status&label=PHP&labelColor=%24.php.color)

![Dynamic Python Badge](https://img.shields


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

# License Manager

A CLI tool for managing license headers in source code files. This tool helps you add, remove, update, and check license headers across your codebase with ease and precision. It utilizes non-printing Unicode character to mark the start & end of the License comment such that you can easier update or remove the license information.

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/license-manager)](https://goreportcard.com/report/github.com/yourusername/license-manager)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)


## Installation

You can install the License Manager through several methods:

### Using Homebrew (macOS and Linux)

The easiest way to install on macOS and Linux is through Homebrew:

```bash
# Add the tap repository
brew tap jeeftor/tap

# Install the package
brew install license-manager
```

To upgrade to the latest version:
```bash
brew upgrade license-manager
```

### Using Debian Package (Ubuntu/Debian)

For Debian-based Linux distributions, you can download and install the .deb package:

```bash
# Download the latest release
curl -LO "https://github.com/jeeftor/license-manager/releases/latest/download/license-manager_$(curl -s https://api.github.com/repos/jeeftor/license-manager/releases/latest | grep tag_name | cut -d '"' -f 4 | cut -c 2-)_linux_amd64.deb"

# Install the package
sudo dpkg -i license-manager_*_linux_amd64.deb

# Install dependencies if needed
sudo apt-get install -f
```

You can also download the .deb package directly from the [releases page](https://github.com/jeeftor/license-manager/releases).

### Using Go

If you prefer to install using Go:

```bash
go install github.com/jeeftor/license-manager@latest
```

### From Source

For the latest development version:

```bash
git clone https://github.com/jeeftor/license-manager.git
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
