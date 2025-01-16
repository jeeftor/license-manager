package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LoadStandardLicense loads a standard license from the templates/licenses directory
// name should be one of: mit, apache, gpl, bsd, mpl
// year and fullname will be substituted in the license text where appropriate
func LoadStandardLicense(name, fullname string) (string, error) {
	name = strings.ToLower(name)
	filename := ""
	switch name {
	case "mit":
		filename = "mit.txt"
	case "apache":
		filename = "apache.txt"
	case "gpl":
		filename = "gpl.txt"
	case "bsd":
		filename = "bsd.txt"
	case "mpl":
		filename = "mpl.txt"
	default:
		return "", fmt.Errorf("unknown standard license: %s", name)
	}

	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// The templates directory should be relative to the executable
	licensePath := filepath.Join(filepath.Dir(exePath), "..", "templates", "licenses", filename)

	// Read the license file
	content, err := os.ReadFile(licensePath)
	if err != nil {
		return "", fmt.Errorf("failed to read license file %s: %v", filename, err)
	}

	// Replace placeholders
	licenseText := string(content)
	currentYear := time.Now().Year()
	licenseText = strings.ReplaceAll(licenseText, "[year]", fmt.Sprintf("%d", currentYear))
	licenseText = strings.ReplaceAll(licenseText, "[fullname]", fullname)

	return licenseText, nil
}
