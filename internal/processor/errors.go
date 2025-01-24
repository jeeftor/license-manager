// internal/processor/errors.go
package processor

import (
	"license-manager/internal/license"
)

// CheckError represents an error during license checking
type CheckError struct {
	Status license.Status
	Msg    string
}

func (e *CheckError) Error() string {
	return e.Msg
}

// NewCheckError creates a new CheckError
func NewCheckError(status license.Status, msg string) *CheckError {
	return &CheckError{
		Status: status,
		Msg:    msg,
	}
}
