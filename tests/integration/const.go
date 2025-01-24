package integration

import (
	"os"
	"path/filepath"
)

var projectRoot string
var mitLicense string
var devLicense string

const testDir = "test_data"
const templateDir = "templates"

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectRoot = filepath.Dir(filepath.Dir(wd))

	mitLicense = filepath.Join(projectRoot, templateDir, "licenses", "mit.txt")
	devLicense = filepath.Join(projectRoot, templateDir, "licenses", "dev.txt")
}
