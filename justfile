# Core variables
go_cmd := "go run main.go"
license_path := "./LICENSE"
test_data_dir := "./test_data"

# Define language patterns
languages := "c cpp csharp css go html java javascript lua perl php python r ruby rust sass scss shell swift typescript xml yaml"

# Define patterns for each language
c_patterns := "*.c *.h"
cpp_patterns := "*.cpp *.hpp *.cc *.hh"
csharp_patterns := "*.cs"
css_patterns := "*.css"
go_patterns := "*.go"
html_patterns := "*.html *.htm"
java_patterns := "*.java"
javascript_patterns := "*.js *.jsx"
lua_patterns := "*.lua"
perl_patterns := "*.pl *.pm"
php_patterns := "*.php"
python_patterns := "*.py"
r_patterns := "*.r *.R"
ruby_patterns := "*.rb"
rust_patterns := "*.rs"
sass_patterns := "*.sass"
scss_patterns := "*.scss"
shell_patterns := "*.sh *.bash"
swift_patterns := "*.swift"
typescript_patterns := "*.ts *.tsx"
xml_patterns := "*.xml"
yaml_patterns := "*.yml *.yaml"

# Helper function to get patterns for a language
[private]
get-patterns lang:
    #!/usr/bin/env bash
    case "{{lang}}" in
        "c") echo "{{c_patterns}}" ;;
        "cpp") echo "{{cpp_patterns}}" ;;
        "csharp") echo "{{csharp_patterns}}" ;;
        "css") echo "{{css_patterns}}" ;;
        "go") echo "{{go_patterns}}" ;;
        "html") echo "{{html_patterns}}" ;;
        "java") echo "{{java_patterns}}" ;;
        "javascript") echo "{{javascript_patterns}}" ;;
        "lua") echo "{{lua_patterns}}" ;;
        "perl") echo "{{perl_patterns}}" ;;
        "php") echo "{{php_patterns}}" ;;
        "python") echo "{{python_patterns}}" ;;
        "r") echo "{{r_patterns}}" ;;
        "ruby") echo "{{ruby_patterns}}" ;;
        "rust") echo "{{rust_patterns}}" ;;
        "sass") echo "{{sass_patterns}}" ;;
        "scss") echo "{{scss_patterns}}" ;;
        "shell") echo "{{shell_patterns}}" ;;
        "swift") echo "{{swift_patterns}}" ;;
        "typescript") echo "{{typescript_patterns}}" ;;
        "xml") echo "{{xml_patterns}}" ;;
        "yaml") echo "{{yaml_patterns}}" ;;
    esac

# Ensure test directory exists before running commands
[private]
ensure-test-dir:
    #!/usr/bin/env bash
    if [ ! -d "{{test_data_dir}}" ]; then
        echo "Test directory does not exist. Creating test files..."
        {{go_cmd}} build-test-data
    fi

# Build and clean commands
clean:
    rm -rf dist junit.xml integration-status.json test-output.json

build:
    goreleaser build --snapshot --clean

clean-test-dir:
    rm -rf {{test_data_dir}}

test-dir:
    {{go_cmd}} build-test-data

# Generic command for any language
[private]
run-command lang cmd *FLAGS: ensure-test-dir
    #!/usr/bin/env bash
    patterns=$(just get-patterns {{lang}})

    echo "Running {{cmd}} for {{lang}} files... [$patterns]"

    # Split patterns and create find arguments
    find_args=""
    for pattern in $patterns; do
        if [ -n "$find_args" ]; then
            find_args="$find_args -o"
        fi
        find_args="$find_args -name \"$pattern\""
    done

    echo "Looking for files in: {{test_data_dir}}/{{lang}}"
    echo "Command: find {{test_data_dir}}/{{lang}} -type f \( $find_args \) 2>/dev/null"

    # Use eval to properly handle the complex find command
    files=$(eval "find {{test_data_dir}}/{{lang}} -type f \( $find_args \) 2>/dev/null")

    if [ -n "$files" ]; then
        echo "Found files: $files"
        echo "$files" | xargs -I {} {{go_cmd}} {{cmd}} --input {} {{FLAGS}} --license {{license_path}}
    else
        echo "Warning: No {{lang}} files found matching patterns: $patterns"
        echo "If this is unexpected, try running 'just test-dir' to recreate test files"
    fi

# Language-specific commands
add lang: (run-command lang "add" "--log-level" "notice") # "--verbose")
check lang: (run-command lang "check" "--log-level" "notice")
update lang: (run-command lang "update" "--log-level" "notice")
debug lang: (run-command lang "debug")
remove lang: (run-command lang "remove")
modify lang: ensure-test-dir
    #!/usr/bin/env bash
    find {{test_data_dir}}/{{lang}} -type f -name "hello.*" 2>/dev/null | \
    xargs -I {} sed -i '' '8,15d' {}

# Combined commands for all languages
add-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just add $lang
    done

check-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just check $lang
    done

update-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just update $lang
    done

debug-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just debug $lang
    done

remove-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just remove $lang
    done

modify-all: ensure-test-dir
    #!/usr/bin/env bash
    for lang in {{languages}}; do
        just modify $lang
    done

# Pre-commit hook
pre-commit:
    {{go_cmd}} pre-commit --license {{license_path}} || true

# Show available languages
list-languages:
    @echo "Available languages:"
    @echo {{languages}} | tr ' ' '\n' | sed 's/^/- /'
