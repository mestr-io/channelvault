package app_helpers

import (
	"fmt"
	"runtime"
)

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
