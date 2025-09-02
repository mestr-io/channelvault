package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestGetSqliteDatabasePath(t *testing.T) {
	// To make this test hermetic, we need to temporarily override the XDG base directories.
	// We'll set them to a temporary directory and restore them after the test.
	oldXDGDataHome := os.Getenv("XDG_DATA_HOME")
	tempDir, err := os.MkdirTemp("", "test-xdg-")
	assert.NoError(t, err)
	defer func(){
		err = os.RemoveAll(tempDir)
		if err != nil {
			panic(err)
		}
	}()

	err = os.Setenv("XDG_DATA_HOME", tempDir)
	if err != nil {
		panic(err)
	}
	// Reset XDG base directories to pick up the new environment variable
	xdg.Reload()

	defer func() {
		err = os.Setenv("XDG_DATA_HOME", oldXDGDataHome)
		if err != nil {
			panic(err)
		}
		xdg.Reload()
	}()

	appName := "TestApp"
	expectedPath := filepath.Join(tempDir, appName, "data.db")

	path, err := getSqliteDatabasePath(appName)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)

	// Check that the directory was created
	_, err = os.Stat(filepath.Dir(path))
	assert.NoError(t, err, "The directory should have been created by xdg.DataFile")
}
