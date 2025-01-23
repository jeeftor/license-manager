// internal/processor/errors.go
package processor

import "license-manager/internal/license"

// CheckError represents a failure in license checking
type CheckError struct {
	msg    string
	Status license.Status
}

func NewCheckError(msg string, status license.Status) *CheckError {
	return &CheckError{msg: msg, Status: status}
}

func (e *CheckError) Error() string {
	return e.msg
}
