clean:
	rm -rf dist
build:
	goreleaser build --snapshot --clean

clean-test-dir:
	rm -rf test_data
test-dir:
	go run main.go build-test-data

check_py:
	go run main.go check --input "./test_data/**/*.py" --verbose --license ./LICENSE
add_py:
	go run main.go add --input "./test_data/**/*.py" --verbose --license ./LICENSE
modify_py:
	sed -i '' '8,15d' ./test_data/python/hello.py
update_py:
	 go run main.go update --input "./test_data/**/*.py" --verbose --license ./LICENSE
debug_py:
	go run main.go debug --input "./test_data/python/hello.py"



check_go:
	go run main.go check --input "./test_data/**/*.go" --verbose --license ./LICENSE
add_go:
	go run main.go add --input "./test_data/**/*.go" --verbose --license ./LICENSE
modify_go:
	sed -i '' '8,15d' ./test_data/go/hello.go
update_go:
	 go run main.go update --input "./test_data/**/*.go" --verbose --license ./LICENSE
debug_go:
	go run main.go debug --input "./test_data/go/hello.go" --verbose





check:
	go run main.go check --input "./test_data/**/*.*"  --license ./LICENSE
add:
	go run main.go add --input "./test_data/**/*.*" --verbose --license ./LICENSE
update:
	 go run main.go update --input "./test_data/**/*.*" --verbose --license ./LICENSE
