package errors

import "fmt"

// LicenseError represents a license-related error
type LicenseError struct {
	Message string
	File    string
}

func (e *LicenseError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("license error in %s: %s", e.File, e.Message)
	}
	return fmt.Sprintf("license error: %s", e.Message)
}

// ValidationError represents a validation-related error
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// FileError represents a file operation error
type FileError struct {
	Message string
	Path    string
	Op      string
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s failed for %s: %s", e.Op, e.Path, e.Message)
}

// NewLicenseError creates a new LicenseError
func NewLicenseError(message string, file string) *LicenseError {
	return &LicenseError{Message: message, File: file}
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, field string) *ValidationError {
	return &ValidationError{Message: message, Field: field}
}

// NewFileError creates a new FileError
func NewFileError(message string, path string, op string) *FileError {
	return &FileError{Message: message, Path: path, Op: op}
}
