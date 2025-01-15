package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Header      string
	Footer      string
	LicenseText string
	Input       string
	Skip        string
	Prompt      bool
	DryRun      bool
}

type FileProcessor struct {
	config Config
}

func NewFileProcessor(config Config) *FileProcessor {
	// Read license text from file if provided
	if config.LicenseText != "" {
		content, err := os.ReadFile(config.LicenseText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading license file: %v\n", err)
			os.Exit(1)
		}
		config.LicenseText = string(content)
	}

	return &FileProcessor{
		config: config,
	}
}

func (fp *FileProcessor) processFiles(action func(string, string, *LicenseManager) error) error {
	var files []string
	for _, pattern := range strings.Split(fp.config.Input, ",") {
		if pattern != "" {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return fmt.Errorf("error with pattern %s: %v", pattern, err)
			}
			files = append(files, matches...)
		}
	}

	for _, file := range files {
		if err := fp.processFile(file, action); err != nil {
			return err
		}
	}
	return nil
}

func (fp *FileProcessor) processFile(filename string, action func(string, string, *LicenseManager) error) error {
	// Skip files that match skip patterns
	for _, pattern := range strings.Split(fp.config.Skip, ",") {
		if pattern != "" {
			matched, err := filepath.Match(pattern, filepath.Base(filename))
			if err != nil {
				return fmt.Errorf("invalid skip pattern %s: %v", pattern, err)
			}
			if matched {
				return nil
			}
		}
	}

	if fp.config.DryRun {
		fmt.Printf("Would process file: %s\n", filename)
		return nil
	}

	if fp.config.Prompt {
		if !promptUser(fmt.Sprintf("Process file %s?", filename)) {
			return nil
		}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	// Create a new LicenseManager with the appropriate comment style for this file
	style := getCommentStyle(filename)
	license := NewLicenseManager(fp.config.Header, fp.config.Footer, fp.config.LicenseText, style)

	return action(filename, string(content), license)
}

func (fp *FileProcessor) Add() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		if license.CheckLicense(content) {
			return nil
		}
		newContent := license.AddLicense(content)
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Remove() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		newContent := license.RemoveLicense(content)
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Update() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		newContent := license.UpdateLicense(content)
		return os.WriteFile(filename, []byte(newContent), 0644)
	})
}

func (fp *FileProcessor) Check() error {
	return fp.processFiles(func(filename, content string, license *LicenseManager) error {
		if !license.CheckLicense(content) {
			fmt.Printf("License missing or invalid in file: %s\n", filename)
		}
		return nil
	})
}

func promptUser(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", message)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
