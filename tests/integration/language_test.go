package integration

import (
	"os"
	"testing"
)

type languageTest struct {
	language string
	patterns []string
	t        *testing.T
	helper   *testHelper
}

func TestLanguages(t *testing.T) {
	tests := map[string]func(*languageTest, *testing.T){
		"Add":         (*languageTest).testAdd,
		"CheckFail":   (*languageTest).testCheckFail,
		"UpdateCheck": (*languageTest).testUpdateCheck,
		"Remove":      (*languageTest).testRemove,
	}

	for testName, testFn := range tests {
		t.Run(testName, func(t *testing.T) {
			resetTestData(t) // Reset before each test type

			for _, tc := range languageDefinitions[5:6] {
				tc := tc
				t.Run(tc.language, func(t *testing.T) {
					t.Parallel()
					lt := &languageTest{
						language: tc.language,
						patterns: tc.patterns,
						t:        t,
						helper:   newTestHelper(t),
					}
					testFn(lt, t)
				})
			}
		})
	}
}

func (lt *languageTest) testAdd(t *testing.T) {

	patterns := lt.helper.getPattern(languageData{patterns: lt.patterns})

	if err := lt.helper.verifyLicenseMissing(patterns); err != nil {
		t.Fatal(err)
	}

	if _, _, err := lt.helper.runLicenseCommand("add", patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyLicensePresent(patterns); err != nil {
		t.Fatal(err)
	}
}

func (lt *languageTest) testCheckFail(t *testing.T) {

	patterns := lt.helper.getPattern(languageData{patterns: lt.patterns})

	licenseFile := createLicenseFile(t, "incorrect license")
	defer os.Remove(licenseFile)

	if _, _, err := lt.helper.runLicenseCommand("add", patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyLicensePresent(patterns); err != nil {
		t.Fatal(err)
	}

	newLicenseFile := createLicenseFile(t, "new license")
	defer os.Remove(newLicenseFile)

	if _, _, err := lt.helper.runLicenseCommand("check", patterns); err != nil {
		t.Fatal(err)
	}
}

func (lt *languageTest) testUpdateCheck(t *testing.T) {

	patterns := lt.helper.getPattern(languageData{patterns: lt.patterns})

	licenseFile := createLicenseFile(t, "original license")
	defer os.Remove(licenseFile)

	if _, _, err := lt.helper.runLicenseCommand("add", patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyLicensePresent(patterns); err != nil {
		t.Fatal(err)
	}

	newLicenseFile := createLicenseFile(t, "new license")
	defer os.Remove(newLicenseFile)

	if _, _, err := lt.helper.runLicenseCommand("update", patterns); err != nil {
		t.Fatal(err)
	}

	if _, _, err := lt.helper.runLicenseCommand("check", patterns); err != nil {
		t.Fatal(err)
	}
}

func (lt *languageTest) testRemove(t *testing.T) {

	patterns := lt.helper.getPattern(languageData{patterns: lt.patterns})

	if _, _, err := lt.helper.runLicenseCommand("add", patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyLicensePresent(patterns); err != nil {
		t.Fatal(err)
	}

	if _, _, err := lt.helper.runLicenseCommand("remove", patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyContentMatchesTemplate(patterns); err != nil {
		t.Fatal(err)
	}

	if err := lt.helper.verifyLicenseMissing(patterns); err != nil {
		t.Fatal(err)
	}
}
