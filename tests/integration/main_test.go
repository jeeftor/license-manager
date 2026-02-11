package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jeeftor/license-manager/internal/license"
	"github.com/jeeftor/license-manager/internal/styles"

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
				fileName: filepath.Base(match),
				language: style,
				filePath: filepath.Join(
					projectRoot,
					outputPath,
					style.Language,
					filepath.Base(match),
				),
				templateFilePath: filepath.Join(
					projectRoot,
					searchPath,
					style.Language,
					filepath.Base(match),
				),
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
func writeIntegrationStatus() {
	mtn, _ := time.LoadLocation("America/Denver")
	now := time.Now().In(mtn)

	type langStatus struct {
		Color    string `json:"color"`
		ColorHex string `json:"colorHex"`
		ColorRgb string `json:"colorRgb"`
		Lang     string `json:"lang"`
		Status   string `json:"status"`
		Emoji    string `json:"emoji"`
		Text     string `json:"text"`
	}

	output := make(map[string]interface{})
	output["date"] = now.Format("Jan 2 2006")
	output["time"] = now.Format("3:04pm") + " Mountain Time"

	for lang, status := range testStatusByLanguage {
		colorName := "red"
		colorHex := "ff0000"
		colorRgb := "rgb(255,0,0)"
		emoji := "❌"
		if status == "Pass" {
			colorName = "green"
			colorHex = "00FF00"
			colorRgb = "rgb(0,255,0)"
			emoji = "✅"
		}
		output[lang] = langStatus{
			Color:    colorName,
			ColorHex: colorHex,
			ColorRgb: colorRgb,
			Lang:     lang,
			Status:   status,
			Emoji:    emoji,
			Text:     fmt.Sprintf("%s %s", emoji, status),
		}
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(filepath.Join(getProjectRoot(), "integration-status.json"), jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write integration status to file: %v", err)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()

	writeIntegrationStatus()
	os.Exit(code)
}

func TestMatrix(t *testing.T) {

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

func checkLicenseWithErrorValue(
	file *singleTestFile,
	licensePath string,
	wantedExitCodes []int,
) error {
	stdout, stderr, err := CheckLicense(file.filePath, licensePath)

	// Extract exit code from stderr. go run wraps the program exit code —
	// it always returns exit code 1 for any non-zero program exit.
	// The actual program exit code appears as "exit status N" at the end of stderr.
	exitCode := 0
	if err != nil {
		// Scan from the end of stderr for "exit status N"
		lines := strings.Split(stderr, "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "exit status ") {
				fmt.Sscanf(line, "exit status %d", &exitCode)
				break
			}
		}
	}

	// Check if the exit code is one of the wanted codes
	for _, code := range wantedExitCodes {
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
	switch err := err.(type) {
	case *LineEndingDiffError:
		// Line ending differences are acceptable, don't fail the test
		t.Logf("Note: Files differ only in line endings")
	case *ContentDiffError:
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 10 - Content differs:\n%s", diffText)
	case nil:
		// Files are identical, this is good
	default:
		testStatusByLanguage[file.language.Language] = "Fail"
		t.Fatalf("Step 10 - Unexpected error comparing files: %v", err)
	}

}

func testRemove(file *singleTestFile, t *testing.T) {
	// Remove license test implementation
}

// Custom error types
type LineEndingDiffError struct{}

func (e *LineEndingDiffError) Error() string {
	return "files differ only in line endings"
}

type ContentDiffError struct {
	diff string
}

func (e *ContentDiffError) Error() string {
	return fmt.Sprintf("files have different content:\n%s", e.diff)
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

	// If lengths differ by 1 and the longer one ends with a line ending,
	// check if that's the only difference
	lenDiff := len(content1) - len(content2)
	if abs(lenDiff) == 1 {
		var longer, shorter []byte
		if len(content1) > len(content2) {
			longer = content1
			shorter = content2
		} else {
			longer = content2
			shorter = content1
		}

		lastByte := longer[len(longer)-1]
		if lastByte == '\n' || lastByte == '\r' {
			// Check if everything else matches
			if bytes.Equal(longer[:len(longer)-1], shorter) {
				return "", &LineEndingDiffError{}
			}
		}
	}

	// If we get here, check for actual content differences
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), true)

	if len(diffs) > 1 || (len(diffs) == 1 && diffs[0].Type != diffmatchpatch.DiffEqual) {
		diffText := dmp.DiffPrettyText(diffs)
		return diffText, &ContentDiffError{diff: diffText}
	}

	return "", nil
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//func diffFiles(file1, file2 string) (string, error) {
//	content1, err := os.ReadFile(file1)
//	if err != nil {
//		return "", err
//	}
//
//	content2, err := os.ReadFile(file2)
//	if err != nil {
//		return "", err
//	}
//
//	dmp := diffmatchpatch.New()
//	diffs := dmp.DiffMain(string(content1), string(content2), true)
//
//	// Check if there are any differences
//	if len(diffs) > 1 || (len(diffs) == 1 && diffs[0].Type != diffmatchpatch.DiffEqual) {
//		diffText := dmp.DiffPrettyText(diffs)
//		return diffText, fmt.Errorf("files are different")
//	}
//
//	return "", nil
//}
