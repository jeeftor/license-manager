# Pre-commit Hooks

This package provides a variety of license-related pre-commit hooks to manage license headers in your code. Always run `pre-commit autoupdate` to ensure you have the latest version of the hooks.

## Available Hooks

 Hook ID | Name | Description | Example Usage |
|---------|------|-------------|---------------|
| verify | License Manager - Check | Checks license headers in staged files | `- id: verify` |
| verify-debug | License Manager - Check (Debug) | Checks licenses with debug logging enabled | `- id: verify-debug` |
| add-missing | License Manager - Add Missing | Adds missing license headers to files | `- id: add-missing` |
| update | License Manager - Update | Updates existing license blocks | `- id: update`|
| version | License Manager - Version | Prints version information | `- id: version` |


## Example Configurations

### Basic License Check


```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.3.3 # This tag is very old -> run pre-commit autoupdate
    hooks:
      - id: verify
      # This configuraiton isn't recommended
```





### Language-Specific Configuration

Process only specific file types and exclude certain directories:

Using a custom `mit.txt` as a license file

```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.3.3
    hooks:
      - id: verify
        exclude: ^vendor/|/vendor/|node_modules/
        types_or: [go, python, javascript]
        args: [--license, ./mit.txt]
```

Using the standard `LICENSE` but with a `swords` style

```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.3.3
    hooks:
      - id: verify
        exclude: ^vendor/|/vendor/|node_modules/
        types_or: [go, python, javascript]
        args: [--style, swords]
```



### License Check - Debug Configuration

Use this configuration when troubleshooting license issues:

```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.3.3
    hooks:
      - id: verify-debug
        args: [--license, ./apache.txt]
```

### Automatic License Addition
Automatically add licenses to new files:

```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.3.3
    hooks:
      - id: add-missing
        args: [--license, ./license.txt]
        types_or: [go, python, java, javascript]
```

## Important Notes

1. Always run `pre-commit autoupdate` to ensure you have the latest version of the hooks.
2. Each hook supports the following common arguments:
    - `--license`: Path to your license template file
    - `--log-level`: Logging level (notice/debug)
3. The `types_or` field can be used to specify which file types to process
4. Use `exclude` to skip certain directories or files
5. All hooks support file passing and require pre-commit version 3.0.0 or higher
