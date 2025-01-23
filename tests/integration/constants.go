package integration

import (
	"license-manager/internal/styles"
	"os"
	"path/filepath"
)

type languageData struct {
	language string
	patterns []string
}

type singleFileTestInput struct {
	language styles.CommentLanguage
	filePath string
}

var languageDefinitions = []languageData{
	{"batch", []string{"test_data/batch/*.bat"}},
	{"c", []string{"test_data/c/*.c", "test_data/c/*.h"}},
	{"cpp", []string{"test_data/cpp/*.cpp", "test_data/cpp/*.hpp"}},
	{"csharp", []string{"test_data/csharp/*.cs"}},
	{"css", []string{"test_data/css/*.css"}},
	{"go", []string{"test_data/go/*.go"}},
	{"html", []string{"test_data/html/*.html"}},
	{"java", []string{"test_data/java/*.java"}},
	{"javascript", []string{"test_data/javascript/*.js", "test_data/javascript/*.jsx"}},
	{"lua", []string{"test_data/lua/*.lua"}},
	{"perl", []string{"test_data/perl/*.pl", "test_data/perl/*.pm"}},
	{"php", []string{"test_data/php/*.php"}},
	{"python", []string{"test_data/python/*.py"}},
	{"r", []string{"test_data/r/*.r"}},
	{"ruby", []string{"test_data/ruby/*.rb"}},
	{"rust", []string{"test_data/rust/*.rs"}},
	{"sass", []string{"test_data/sass/*.sass"}},
	{"scss", []string{"test_data/scss/*.scss"}},
	{"shell", []string{"test_data/shell/*.sh", "test_data/shell/*.bash"}},
	{"swift", []string{"test_data/swift/*.swift"}},
	{"typescript", []string{"test_data/typescript/*.ts", "test_data/typescript/*.tsx"}},
	{"xml", []string{"test_data/xml/*.xml"}},
	{"yaml", []string{"test_data/yaml/*.yaml", "test_data/yaml/*.yml"}},
}

var projectRoot string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectRoot = filepath.Dir(filepath.Dir(wd))
}
