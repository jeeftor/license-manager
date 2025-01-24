package dev

import (
	"bufio"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"license-manager/internal/styles"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var testFiles []singleTestFile

type TestStageDefinition struct {
	Name   string
	Helper func(test *singleTestFile, t *testing.T)
}

type singleTestFile struct {
	filePath         string
	templateFilePath string
	fileName         string
	language         styles.CommentLanguage
}

//var projectRoot = getProjectRoot()

func gatherTestFiles(searchPath string, outputPath string) []singleTestFile {
	var files []singleTestFile

	for ext, style := range styles.LanguageExtensions {
		pattern := filepath.Join(projectRoot, searchPath, "**", "*"+ext)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			panic(err)
		}

		for _, match := range matches {
			tf := singleTestFile{
				fileName:         filepath.Base(match),
				language:         style,
				filePath:         filepath.Join(projectRoot, outputPath, style.Language, filepath.Base(match)),
				templateFilePath: filepath.Join(projectRoot, searchPath, style.Language, filepath.Base(match)),
			}

			if err := resetFile(tf); err != nil {
				panic(err)
			}

			fmt.Printf("Added file: %s with language: %s\n", tf.fileName, style.Language)
			files = append(files, tf)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].language.Language != files[j].language.Language {
			return files[i].language.Language < files[j].language.Language
		}
		return files[i].fileName < files[j].fileName
	})

	return files
}

func init() {
	testFiles = gatherTestFiles(templateDir, testDir)
	fmt.Printf("Total files found: %d\n", len(testFiles))

}

func resetFile(tf singleTestFile) error {

	srcPath := filepath.Join(projectRoot, templateDir, tf.language.Language, tf.fileName)
	dstPath := filepath.Join(projectRoot, testDir, tf.language.Language, tf.fileName)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed creating directory: %v", err)
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed reading template: %v", err)
	}

	return os.WriteFile(dstPath, data, 0644)
}

func getProjectRoot() string {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(filepath.Dir(wd))
}

func TestMatrix(t *testing.T) {
	testStage := []TestStageDefinition{
		//{Name: "Add", Helper: testAddAndCheck},
		{Name: "AddModifyUpdateRemove", Helper: testWorkflow},
		//{Name: "Update", Helper: testUpdate},
		//{Name: "Remove", Helper: testRemove},
	}

	for _, file := range testFiles {
		file := file // capture varibale in loop
		for _, stage := range testStage {
			test := stage // capture loop var

			pathComponents := []string{file.language.Language, file.fileName}
			if len(testStage) > 1 {
				pathComponents = append(pathComponents, test.Name)
			}
			t.Run(fmt.Sprintf("%s", strings.Join(pathComponents, "/")), func(t *testing.T) {
				t.Parallel()

				if err := resetFile(file); err != nil {
					t.Fatal(err)
				}
				test.Helper(&file, t)
			})
		}
	}
}

func testAddAndCheck(file *singleTestFile, t *testing.T) {
	// Add license test implementation

	err := verifyLicenseMissing(file, devLicense)
	if err != nil {
		t.Fatalf("License Error: %v\n", err)
	}

	_, _, err = AddLicense(file.filePath, devLicense)
	if err != nil {
		t.Fatalf("Error adding license: %v\n", err)
	}

	err = verifyLicenseExists(file, devLicense)
	if err != nil {
		t.Fatalf("License Error: %v\n", err)

	}
}

func verifyLicenseExists(file *singleTestFile, licensePath string) error {
	return checkLicenseWithErrorValue(file, licensePath, 0)
}

func verifyLicenseMissing(file *singleTestFile, licensePath string) error {
	return checkLicenseWithErrorValue(file, licensePath, 1)
}

func verifyLicenseMismatch(file *singleTestFile, licensePath string) error {
	return checkLicenseWithErrorValue(file, licensePath, 2)
}

func extractErrorText(input string) string {
	var errors []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "ERROR:") {
			errors = append(errors, scanner.Text())
		}
	}
	return strings.Join(errors, "\n")
}

func checkLicenseWithErrorValue(file *singleTestFile, licensePath string, wanted_exit_code int) error {

	stdout, stderr, err := CheckLicense(file.filePath, licensePath)
	if wanted_exit_code == 0 {
		if stderr != "" {
			return fmt.Errorf("license check failed: %v\n %s\nStderr: %s", err, extractErrorText(stdout), stderr)
		} else {
			return nil
		}
	}
	if wanted_exit_code == 1 { // License Mismatch
		if strings.Contains(stderr, "exit status 1") {
			return nil
		} else {
			return fmt.Errorf("license check failed: %v\n %s\nStderr: %s", err, extractErrorText(stdout), stderr)
		}
	}

	if wanted_exit_code == 2 { // License Missing
		if strings.Contains(stderr, "exit status 2") {
			return nil
		} else {
			return fmt.Errorf("license check failed: %v\n %s\nStderr: %s", err, extractErrorText(stdout), stderr)
		}
	}

	return nil

}

func testCheck(file *singleTestFile, t *testing.T) {
	// Check license test implementation
}

func testWorkflow(file *singleTestFile, t *testing.T) {
	err := verifyLicenseMissing(file, devLicense)
	if err != nil {
		t.Fatalf("Step 1 - License Error: %v\n", err)
	}

	_, _, err = AddLicense(file.filePath, devLicense)
	if err != nil {
		t.Fatalf("Step 2 - Error adding license: %v\n", err)
	}

	err = verifyLicenseExists(file, devLicense)
	if err != nil {
		t.Fatalf("Step 3 - License Error: %v\n", err)

	}
	err = verifyLicenseMismatch(file, mitLicense)
	if err != nil {
		t.Fatalf("Step 4 - License Error: %v\n", err)
	}

	// Update the license now

	_, _, err = UpdateLicense(file.filePath, mitLicense)
	if err != nil {
		t.Fatalf("Step 5 - Error updating license: %v\n", err)
	}

	err = verifyLicenseExists(file, mitLicense)
	if err != nil {
		t.Fatalf("Step 6 - License Error: %v\n", err)
	}

	_, _, err = RemoveLicense(file.filePath)
	if err != nil {
		t.Fatalf("Step 7 - Error removing license: %v\n", err)
	}

	err = verifyLicenseMissing(file, mitLicense)
	if err != nil {
		t.Fatalf("Step 8 - License Error: %v\n", err)
	}

	err = verifyLicenseMissing(file, devLicense)
	if err != nil {
		t.Fatalf("Step 9 - License Error: %v\n", err)
	}

	// Lastly verify the template file matches the rest file?

	diffText, err := diffFiles(file.filePath, file.templateFilePath)
	if err != nil {
		t.Fatalf("Step 10 - Diff Error: %v\n%s", err, diffText)
	}

}

func testRemove(file *singleTestFile, t *testing.T) {
	// Remove license test implementation
}

func diffFiles(file1, file2 string) (string, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return "", err
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return "", err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), true)

	return dmp.DiffPrettyText(diffs), nil
}
