// internal/processor/errors.go
package processor

// CheckError represents a failure in license checking
type CheckError struct {
	msg string
}

func NewCheckError(msg string) *CheckError {
	return &CheckError{msg: msg}
}

func (e *CheckError) Error() string {
	return e.msg
}
