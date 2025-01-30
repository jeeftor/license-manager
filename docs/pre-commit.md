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
