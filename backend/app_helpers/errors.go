package app_helpers

import (
	"errors"
	"fmt"
	"runtime"
)

// Finally executes a function and captures any error, returning it along with the previous error.
func Finally(f func() error, previous error) (err error) {
	e := f()

	err = errors.Join(previous, createError("%s", e.Error()))
	return
}

func createError(format string, v ...any) (err error) {
	where := getLineAndNo(3)
	s := fmt.Sprintf(format, v...)
	err = fmt.Errorf("%s %s", where, s)

	return
}

func getLineAndNo(caller int) string {
	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		file = "unknown"
		line = 0
	}

	return fmt.Sprintf("%s:%d", file, line)
}
