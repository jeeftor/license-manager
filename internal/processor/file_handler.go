package processor

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeeftor/license-manager/internal/errors"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"

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

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fh.logger.LogError("Error getting working directory: %v", err)
		return false
	}

	// Convert absolute path to relative path for matching
	normalizedPath := filepath.ToSlash(path)
	if strings.HasPrefix(normalizedPath, cwd) {
		normalizedPath = normalizedPath[len(cwd)+1:]
	}
	if strings.HasPrefix(normalizedPath, "./") {
		normalizedPath = normalizedPath[2:]
	}

	for _, pattern := range strings.Split(fh.skip, ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Normalize pattern but keep the ** globstar
		pattern = filepath.ToSlash(pattern)
		if strings.HasPrefix(pattern, "./") {
			pattern = pattern[2:]
		}

		// Make sure the pattern ends with /** if it's a directory pattern
		if !strings.Contains(pattern, "*") {
			pattern = strings.TrimSuffix(pattern, "/") + "/**"
		}

		matched, err := doublestar.Match(pattern, normalizedPath)
		if err != nil {
			fh.logger.LogError("Error matching skip pattern %s: %v", pattern, err)
			continue
		}
		if matched {
			fh.logger.LogDebug("Path %s matched skip pattern %s", normalizedPath, pattern)
			return true
		}
	}
	return false
}

// FindFiles finds all files matching the input pattern
func (fh *FileHandler) FindFiles(pattern string) ([]string, error) {
	var allFiles []string

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.NewFileError("failed to get working directory", "", "find")
	}
	fh.logger.LogDebug("Current working directory: %s", cwd)

	// Split input patterns
	patterns := strings.Split(pattern, ",")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Normalize pattern
		p = filepath.ToSlash(p)

		// Check if pattern is a direct file path
		if !strings.Contains(p, "*") && !strings.Contains(p, "?") && !strings.Contains(p, "[") {
			absPath := p
			if !filepath.IsAbs(p) {
				absPath = filepath.Join(cwd, p)
			}

			// Check skip pattern BEFORE adding to allFiles
			if fh.shouldSkip(absPath) {
				fh.logger.LogDebug("Skipping file: %s", absPath)
				continue
			}

			if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
				if isProcessableFile(absPath) {
					allFiles = append(allFiles, absPath)
				}
			}
			continue
		}

		// Handle glob patterns
		matches, err := doublestar.Glob(p)
		if err != nil {
			fh.logger.LogError("Invalid glob pattern %s: %v", p, err)
			continue
		}

		for _, match := range matches {
			absMatch := match
			if !filepath.IsAbs(match) {
				absMatch = filepath.Join(cwd, match)
			}

			// Check skip pattern BEFORE processing
			if fh.shouldSkip(absMatch) {
				fh.logger.LogDebug("Skipping matched path: %s", absMatch)
				continue
			}

			info, err := os.Stat(absMatch)
			if err != nil {
				fh.logger.LogError("Error accessing %s: %v", absMatch, err)
				continue
			}

			if info.IsDir() {
				err := filepath.WalkDir(
					absMatch,
					func(path string, d fs.DirEntry, err error) error {
						if err != nil {
							return nil
						}

						// Check skip pattern for each walked file
						if fh.shouldSkip(path) {
							return filepath.SkipDir
						}

						if !d.IsDir() && isProcessableFile(path) {
							allFiles = append(allFiles, path)
						}
						return nil
					},
				)
				if err != nil {
					fh.logger.LogError("Error walking directory %s: %v", absMatch, err)
				}
			} else if isProcessableFile(absMatch) {
				allFiles = append(allFiles, absMatch)
			}
		}
	}

	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	var uniqueFiles []string
	for _, f := range allFiles {
		if !seen[f] {
			seen[f] = true
			uniqueFiles = append(uniqueFiles, f)
		}
	}

	if len(uniqueFiles) == 0 {
		return nil, errors.NewFileError("no matching files found", pattern, "find")
	}

	return uniqueFiles, nil
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
