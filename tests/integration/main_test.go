package integration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"license-manager/internal/license"
	"license-manager/internal/styles"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var testFiles []singleTestFile

var testStatusByLanguage = make(map[string]string) // "Pass" or "Fail"
var statusMutex = &sync.Mutex{}

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

		// Initialize the status for the language
		testStatusByLanguage[style.Language] = "Pass"

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

func writeIntegrationStatus(t *testing.T) {
	// After all tests have completed, write the results to JSON
	//var statusMap []map[string]string
	//for lang, status := range testStatusByLanguage {
	//	statusMap = append(statusMap, map[string]string{
	//		"Language": lang,
	//		"Status":   status,
	//	})
	//}

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(testStatusByLanguage, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write to file
	outputPath := filepath.Join(getProjectRoot(), "integration-status.json")
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write integration status to file: %v", err)
	}
}

func TestMatrix(t *testing.T) {
	defer writeIntegrationStatus(t)

	testStage := []TestStageDefinition{
		//{Name: "Add", Helper: testAddAndCheck},
		{Name: "AddModifyUpdateRemove", Helper: testWorkflow},
		//{Name: "Update", Helper: testUpdate},
		//{Name: "Remove", Helper: testRemove},
	}

	// When testing stuff you want to use this
	// for _, file := range testFiles[0:1] {
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

				// Check if the test failed and update status
				if t.Failed() {
					testStatusByLanguage[file.language.Language] = "Fail"
				}
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
	return checkLicenseWithErrorValue(file, licensePath, []int{int(license.FullMatch)})
}

func verifyLicenseMissing(file *singleTestFile, licensePath string) error {
	return checkLicenseWithErrorValue(file, licensePath, []int{int(license.NoLicense)})
}

func verifyLicenseMismatch(file *singleTestFile, licensePath string) error {
	return checkLicenseWithErrorValue(file, licensePath, []int{
		int(license.ContentMismatch),
		int(license.StyleMismatch),
		int(license.ContentAndStyleMismatch),
	})
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

func checkLicenseWithErrorValue(file *singleTestFile, licensePath string, wantedExitCodes []int) error {
	stdout, stderr, err := CheckLicense(file.filePath, licensePath)

	// Extract exit code from stderr
	exitCode := 0
	if strings.Contains(stderr, "exit status") {
		fmt.Sscanf(stderr, "exit status %d", &exitCode)
	}

	// Check if the exit code is one of the wanted codes
	for _, code := range wantedExitCodes {
		if code == 0 && exitCode == 0 {
			if stderr != "" {
				return fmt.Errorf("license check failed: %v\n %s\nStderr: %s", err, extractErrorText(stdout), stderr)
			}
			return nil
		}
		if code == exitCode {
			return nil
		}
	}

	return fmt.Errorf("license check failed: got exit code %d, wanted one of %v\n %s\nStderr: %s",
		exitCode, wantedExitCodes, extractErrorText(stdout), stderr)
}

func testCheck(file *singleTestFile, t *testing.T) {
	// Check license test implementation
}

func testWorkflow(file *singleTestFile, t *testing.T) {
	err := verifyLicenseMissing(file, devLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 1 - License Error: %v\n", err)
	}

	_, _, err = AddLicense(file.filePath, devLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 2 - Error adding license: %v\n", err)
	}

	err = verifyLicenseExists(file, devLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 3 - Verify Add - License Error: %v\n", err)

	}
	err = verifyLicenseMismatch(file, mitLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 4 - License Error: %v\n", err)
	}

	// Update the license now

	_, _, err = UpdateLicense(file.filePath, mitLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 5 - Error updating license: %v\n", err)
	}

	err = verifyLicenseExists(file, mitLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 6 - License Error: %v\n", err)
	}

	_, _, err = RemoveLicense(file.filePath)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 7 - Error removing license: %v\n", err)
	}

	err = verifyLicenseMissing(file, mitLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 8 - Verify Removal - License Error: %v\n", err)
	}

	err = verifyLicenseMissing(file, devLicense)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 9 - License Error: %v\n", err)
	}

	// Lastly verify the template file matches the rest file?

	diffText, err := diffFiles(file.filePath, file.templateFilePath)
	if err != nil {
		testStatusByLanguage[file.language.Language] = "Fail"
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
