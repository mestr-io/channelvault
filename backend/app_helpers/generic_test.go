package app_helpers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLineAndNo(t *testing.T) {
	assert := assert.New(t)

	thePath := GetLineAndNo(1) //keep at line 13 and in file generic_test.go or UPDATE assert

	_, result := filepath.Split(thePath)

	assert.Equal(result, "generic_test.go:13")
}
