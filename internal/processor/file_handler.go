package processor

import (
	"io/fs"
	"license-manager/internal/errors"
	"license-manager/internal/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

// FileHandler handles file operations
type FileHandler struct {
	logger *logger.Logger
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(logger *logger.Logger) *FileHandler {
	return &FileHandler{
		logger: logger,
	}
}

// FindFiles finds all files matching the input pattern
func (fh *FileHandler) FindFiles(pattern string) ([]string, error) {
	var files []string

	// Check if pattern is a direct file path
	if info, err := os.Stat(pattern); err == nil && !info.IsDir() {
		if isProcessableFile(pattern) {
			return []string{pattern}, nil
		}
		return nil, errors.NewFileError("file is not a processable type", pattern, "validate")
	}

	// Handle glob patterns
	matches, err := doublestar.Glob(pattern)
	if err != nil {
		return nil, errors.NewFileError("invalid glob pattern", pattern, "glob")
	}

	// Process each match
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			fh.logger.LogError("Error accessing %s: %v", match, err)
			continue
		}

		if info.IsDir() {
			// Walk directory
			err := filepath.WalkDir(match, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fh.logger.LogError("Error walking directory %s: %v", path, err)
					return nil
				}

				if !d.IsDir() && isProcessableFile(path) {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				fh.logger.LogError("Error walking directory %s: %v", match, err)
			}
		} else if isProcessableFile(match) {
			files = append(files, match)
		}
	}

	if len(files) == 0 {
		return nil, errors.NewFileError("no matching files found", pattern, "find")
	}

	return files, nil
}

// ReadFile reads a file and returns its content
func (fh *FileHandler) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", errors.NewFileError("failed to read file", path, "read")
	}
	return string(content), nil
}

// WriteFile writes content to a file
func (fh *FileHandler) WriteFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return errors.NewFileError("failed to write file", path, "write")
	}
	return nil
}

// BackupFile creates a backup of a file
func (fh *FileHandler) BackupFile(path string) error {
	content, err := fh.ReadFile(path)
	if err != nil {
		return err
	}

	backupPath := path + ".bak"
	err = fh.WriteFile(backupPath, content)
	if err != nil {
		return errors.NewFileError("failed to create backup", path, "backup")
	}

	return nil
}

// RestoreFile restores a file from its backup
func (fh *FileHandler) RestoreFile(path string) error {
	backupPath := path + ".bak"
	content, err := fh.ReadFile(backupPath)
	if err != nil {
		return err
	}

	err = fh.WriteFile(path, content)
	if err != nil {
		return err
	}

	return os.Remove(backupPath)
}

// isProcessableFile checks if a file should be processed based on its extension
func isProcessableFile(path string) bool {
	// Skip hidden files and directories
	if strings.HasPrefix(filepath.Base(path), ".") {
		return false
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}

	// List of supported file extensions
	supportedExts := map[string]bool{
		".go":   true,
		".py":   true,
		".js":   true,
		".jsx":  true,
		".ts":   true,
		".tsx":  true,
		".java": true,
		".c":    true,
		".cpp":  true,
		".h":    true,
		".hpp":  true,
		".rs":   true,
		".rb":   true,
		".php":  true,
		".cs":   true,
	}

	return supportedExts[ext]
}
