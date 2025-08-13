package app_helpers

import (
	"errors"
	"fmt"
	"runtime"
)

// Error formats an error message with the file and line number where it was called.
func Error(format string, v ...any) (err error) {
	where := getLineAndNo(3)
	s := fmt.Sprintf(format, v...)
	err = fmt.Errorf("%s %s", where, s)

	return
}

// Finally executes a function and captures any error, returning it along with the previous error.
func Finally(f func() error, previous error) (err error) {
	e := f()

	err = errors.Join(previous, Error("%s", e.Error()))
	return
}

func getLineAndNo(caller int) string {
	if caller == 0 {
		caller = 2 // Default to the caller of the caller of this function
	}
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		file = "unknown"
		line = 0
	}

	return fmt.Sprintf("%s:%d", file, line)
}
