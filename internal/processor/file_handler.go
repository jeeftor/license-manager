package processor

import (
	"io/fs"
	"license-manager/internal/errors"
	"license-manager/internal/logger"
	"license-manager/internal/styles"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

// FileHandler handles file operations
type FileHandler struct {
	logger *logger.Logger
	skip   string
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(logger *logger.Logger) *FileHandler {
	return &FileHandler{
		logger: logger,
	}
}

// SetSkipPattern sets the skip pattern for file filtering
func (fh *FileHandler) SetSkipPattern(pattern string) {
	fh.skip = pattern
}

// shouldSkip checks if a file should be skipped based on skip patterns
func (fh *FileHandler) shouldSkip(path string) bool {
	if fh.skip == "" {
		return false
	}

	for _, pattern := range strings.Split(fh.skip, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		matched, err := doublestar.Match(pattern, path)
		if err != nil {
			fh.logger.LogError("Error matching skip pattern %s: %v", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

// FindFiles finds all files matching the input pattern
func (fh *FileHandler) FindFiles(pattern string) ([]string, error) {
	var allFiles []string

	// Split input patterns
	patterns := strings.Split(pattern, ",")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		var files []string

		// Check if pattern is a direct file path
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			if isProcessableFile(p) && !fh.shouldSkip(p) {
				files = append(files, p)
			}
			continue
		}

		// Handle glob patterns
		matches, err := doublestar.Glob(p)
		if err != nil {
			fh.logger.LogError("Invalid glob pattern %s: %v", p, err)
			continue
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

					if !d.IsDir() && isProcessableFile(path) && !fh.shouldSkip(path) {
						files = append(files, path)
					}
					return nil
				})
				if err != nil {
					fh.logger.LogError("Error walking directory %s: %v", match, err)
				}
			} else if isProcessableFile(match) && !fh.shouldSkip(match) {
				files = append(files, match)
			}
		}

		allFiles = append(allFiles, files...)
	}

	if len(allFiles) == 0 {
		return nil, errors.NewFileError("no matching files found", pattern, "find")
	}

	return allFiles, nil
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
	// Skips hidden files and directories
	if strings.HasPrefix(filepath.Base(path), ".") {
		return false
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}

	// Get the comment style for this extension
	style := styles.GetLanguageCommentStyle(ext)
	return style.Language != ""
}
