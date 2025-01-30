# Pre-commit Hooks

This package provides a variety of different license related pre-commit hooks you can use with the code.


## Example


Check for licenses

```yaml

repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.2.4
    hooks:
      - id: license-manager
```

 Example setup for a Go project.
  - Use `mit.txt` as the license file
  - Process only `Python` and `Go` files
  - Ignore the `vendor` directory

```yaml
repos:
  - repo: https://github.com/jeeftor/license-manager
    rev: v0.2.5
    hooks:
      - id: check
        exclude: ^vendor/|/vendor/
        types_or: [go, python]
        args: [--license, ./mit.txt]
```
