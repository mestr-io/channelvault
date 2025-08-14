package app_helpers

import (
	"errors"
	"fmt"
)

// JoinAndAddFileNoToErr joins the provided error with the previous error and adds the file and line number of the caller to the error message.
func JoinAndAddFileNoToErr(e error, previous error) (err error) {

	err = errors.Join(previous, AddFileAndLineNoToErr(e))

	return
}

// AddFileAndLineNoToErr adds the file and line number of the caller to the error message.
// It retrieves the file and line number from the call stack and prepends it into the error message
// as "file:line".
func AddFileAndLineNoToErr(e error) (err error) {

	if e != nil {
		where := GetLineAndNo(3)
		err = fmt.Errorf("%s %w", where, e)

	}

	return
}
