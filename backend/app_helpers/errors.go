package app_helpers

import (
	"errors"
	"fmt"	
	"runtime"
)

// Error formats an error message with the file and line number where it was called.
func Error(format string, v ...any) (err error) {
	where := getLineAndNo(0)
	s := fmt.Sprintf(format, v...)
	err = fmt.Errorf("%s %s", where, s)

	return
}

// Finally executes a function and captures any error, returning it along with the previous error.
func Finally(f func() error, previous error) (err error) {
	if previous != nil {
		err = previous
	}

	e := f()

	if e != nil {
		if previous != nil {
			err = errors.Join(previous, Error("%s", e.Error()))
		} else {
			err = e
		}
	}

	return
}

func getLineAndNo(caller int) string {
	if caller == 0 {
		caller = 2 // Default to the caller of this function
	}	
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		file = "unknown"
		line = 0
	}

	return fmt.Sprintf("%s:%d", file, line)
}
