package errors

import "fmt"

// New ...
func New(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
