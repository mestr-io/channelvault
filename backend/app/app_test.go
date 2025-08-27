package app

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "/test/folder")
	assert.Equal(t, mockRepo, app.repository)
	assert.Equal(t, "/test/folder", app.filesFolder)
}

func TestGetViewState(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "")
	app.ctx = context.Background()

	t.Run("Dashboard State", func(t *testing.T) {
		mockRepo.On("NoFinishedImports", mock.Anything).Return(true, nil).Once()
		state := app.GetViewState()
		assert.Equal(t, string(ViewStateDashboard), state)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Setup State", func(t *testing.T) {
		mockRepo.On("NoFinishedImports", mock.Anything).Return(false, nil).Once()
		state := app.GetViewState()
		assert.Equal(t, string(ViewStateSetup), state)
		mockRepo.AssertExpectations(t)
	})
}

