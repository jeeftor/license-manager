# License Manager
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeeftor/license-manager)](https://goreportcard.com/report/github.com/jeeftor/license-manager)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

![GitHub commit activity](https://img.shields.io/github/commit-activity/w/jeeftor/license-manager)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/jeeftor/license-manager/latest)
![GitHub Release Date](https://img.shields.io/github/release-date/jeeftor/license-manager)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/jeeftor/license-manager/total)
![GitHub Repo stars](https://img.shields.io/github/stars/jeeftor/license-manager)

![docs/logo.png](docs/logo.png)
A CLI tool and a [Pre-Commit Hook](docs/pre-commit.md) for managing license headers in source code files. This tool helps you add, remove, update, and check license headers across your codebase with ease and precision. It utilizes non-printing Unicode character to mark the start & end of the License comment such that you can easier update or remove the license information.

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


| ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=cpp.text&label=cpp) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=csharp.text&label=csharp) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=css.text&label=css) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=go.text&label=go) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=html.text&label=html) |
|----------|---------|---------|---------|---------|
| ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=ini.text&label=ini) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=java.text&label=java) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=javascript.text&label=javascript) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=kotlin.text&label=kotlin) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=markdown.text&label=markdown) |
| ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=php.text&label=php) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=python.text&label=python) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=ruby.text&label=ruby) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=rust.text&label=rust) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=scala.text&label=scala) |
| ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=shell.text&label=shell) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=swift.text&label=swift) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=toml.text&label=toml) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=typescript.text&label=typescript) | ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=xml.text&label=xml) |
| ![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgist.githubusercontent.com%2Fjeeftor%2Ff639b71257cceeb283a30cba77ee17c9%2Fraw%2Fintegration-status.json&query=yaml.text&label=yaml) |         |         |         |         |


## Installation

The following details CLI based installation. License manager can also be used as a [docs/pre-commit.md](docs/pre-commit.md) hook.

# License Manager

A CLI tool for managing license headers in source code files. This tool helps you add, remove, update, and check license headers across your codebase with ease and precision. It utilizes non-printing Unicode character to mark the start & end of the License comment such that you can easier update or remove the license information.


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
go install github.com/jeeftor/license-manager/cmd/license-manager@latest
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
license-manager remove --input "**/*.go" --input "**/*.py" --skip "./vendor/**" --skip "./venv/**"

# Update existing license headers
license-manager update --license NEW_LICENSE.txt --input "**/*.go"

# Check license headers
license-manager check --license LICENSE.txt --input "**/*.go"
```

### Command Options

- `--license` _string_      Path to license text file (required for add/update/check)
- `--input` _strings_      Input file patterns (can be comma-separated or multiple flags)
- `--skip` _strings_       Patterns to skip (can be comma-separated or multiple flags)
- `--style` _string_       Preset style for header/footer (default "hash")
- `--comments` _string_    Force comment style (no|single|multi)
- `--log-level` _string_   Log level (debug|info|notice|warn|error) (default "notice")

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
