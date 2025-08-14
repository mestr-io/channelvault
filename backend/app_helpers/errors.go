package app_helpers

import (
	"errors"
	"fmt"
	"runtime"
)

// Finally executes a function and captures any error, returning it along with the previous error.
func Finally(f func() error, previous error) (err error) {
	e := f()

	err = errors.Join(previous, AddFileAndLineNoToErr(e))

	return
}

// AddFileAndLineNoToErr adds the file and line number of the caller to the error message.
// It retrieves the file and line number from the call stack and prepends it into the error message
// as "file:line".
func AddFileAndLineNoToErr(e error) (err error) {
	where := GetLineAndNo(3)
	err = fmt.Errorf("%s %w", where, e)

	return
}

// GetLineAndNo retrieves the file name and line number of the caller at a specified depth.
// returns a string formatted as "file:line".
func GetLineAndNo(caller int) string {
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		file = "unknown"
		line = 0
	}

	return fmt.Sprintf("%s:%d", file, line)
}
