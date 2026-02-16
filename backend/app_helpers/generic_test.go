package app_helpers

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLineAndNo(t *testing.T) {
	assert := assert.New(t)
	thePath := GetLineAndNo(1) //keep at line 13 and in file generic_test.go or UPDATE assert

	_, result := filepath.Split(thePath)

	assert.Equal(result, "generic_test.go:13")
}

func TestGetLineAndNoBis(t *testing.T) {
	result := GetLineAndNo(0)
	if !strings.Contains(result, "generic_test.go") && !strings.Contains(result, "generic.go") {
		t.Errorf("Expected file name in result, got: %s", result)
	}
	if !strings.Contains(result, ":") {
		t.Errorf("Expected ':' in result, got: %s", result)
	}
}
