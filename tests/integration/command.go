package integration

import (
	"bytes"
	"os/exec"
	"testing"
)

func runCommand(t *testing.T, args ...string) (string, string, error) {
	cmdArgs := append([]string{"run", "main.go"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
