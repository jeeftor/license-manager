package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCommand(args ...string) (string, string, error) {
	cmdArgs := append([]string{"run", "cmd/license-manager/main.go"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	stdoutStr := strings.TrimSpace(stdout.String())
	stderrStr := strings.TrimSpace(stderr.String())

	return stdoutStr, stderrStr, err
}

func fileExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file: %s - %v", path, err)
	}
	return nil
}

func AddLicense(inputFile string, license string) (string, string, error) {
	if err := fileExists(inputFile); err != nil {
		return "", "", err
	}
	if err := fileExists(license); err != nil {
		return "", "", err
	}
	args := []string{"add", "--input", inputFile, "--license", license}
	return runCommand(args...)
}

func UpdateLicense(inputFile string, license string) (string, string, error) {
	if err := fileExists(inputFile); err != nil {
		return "", "", err
	}
	if err := fileExists(license); err != nil {
		return "", "", err
	}
	args := []string{"update", "--input", inputFile, "--license", license}
	return runCommand(args...)
}

func CheckLicense(inputFile string, license string) (string, string, error) {
	if err := fileExists(inputFile); err != nil {
		return "", "", err
	}
	if err := fileExists(license); err != nil {
		return "", "", err
	}
	args := []string{"check", "--input", inputFile, "--license", license}
	return runCommand(args...)
}

func RemoveLicense(inputFile string) (string, string, error) {
	if err := fileExists(inputFile); err != nil {
		return "", "", err
	}
	args := []string{"remove", "--input", inputFile}
	return runCommand(args...)
}
