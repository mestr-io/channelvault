package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetChannels(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "")

	expectedChannels := []models.Channel{{ID: "C1", Name: "channel1"}}
	mockRepo.On("GetChannels").Return(expectedChannels, nil).Once()

	channels, err := app.GetChannels()

	assert.NoError(t, err)
	assert.Equal(t, expectedChannels, channels)
	mockRepo.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "")
	app.ctx = context.Background()

	expectedMessages := []models.Message{{Text: "Hello"}}
	mockRepo.On("GetMessages", mock.Anything /*ctx*/, "C1", "", 10).Return(expectedMessages, nil).Once()

	messages, err := app.GetMessages("C1", "", 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetUserAvatar(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "")
	app.ctx = context.Background()

	expectedAvatar := "http://example.com/avatar.png"
	mockRepo.On("GetUserAvatar", mock.Anything /*ctx*/, "U1").Return(expectedAvatar, nil).Once()

	avatar, err := app.GetUserAvatar("U1")

	assert.NoError(t, err)
	assert.Equal(t, expectedAvatar, avatar)
	mockRepo.AssertExpectations(t)
}

func TestGetThreadByTs(t *testing.T) {
	mockRepo := new(MockRepository)
	app := New(mockRepo, "")
	app.ctx = context.Background()

	expectedMessages := []models.Message{{Text: "Reply"}}
	mockRepo.On("GetThreadByTs", mock.Anything /*ctx*/, "123.456", "C1", int64(0)).Return(expectedMessages, nil).Once()

	messages, err := app.GetThreadByTs("123.456", "C1", 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetBase64Image(t *testing.T) {
	// Create a temporary directory for files
	tempDir, err := os.MkdirTemp("", "test-files")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	mockRepo := new(MockRepository)
	app := New(mockRepo, tempDir)
	app.ctx = context.Background()

	file := &models.File{
		ID:       "F1",
		Name:     "test.jpg",
		Mimetype: "image/jpeg",
	}

	t.Run("Image already exists locally", func(t *testing.T) {
		// Create a dummy file
		dummyFilePath := filepath.Join(tempDir, "F1_test.jpg")
		// NOTE: "test data" in base64 is dGVzdCBkYXRh
		err := os.WriteFile(dummyFilePath, []byte("test data"), 0644)
		assert.NoError(t, err)

		mockRepo.On("GetFileById", mock.Anything, "F1").Return(file, nil).Once()

		b64, err := app.GetBase64Image("F1")

		assert.NoError(t, err)
		/* "test data" in base64 is dGVzdCBkYXRh */
		assert.Equal(t, "data:image/jpeg;base64,dGVzdCBkYXRh", b64)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Image needs to be downloaded", func(t *testing.T) {
		// Setup a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// NOTE: "downloaded data" in base64 is ZG93bmxvYWRlZCBkYXRh
			_, err := fmt.Fprint(w, "downloaded data")
			if err != nil {
				t.Fatalf("Failed to write response: %v", err)
			}
		}))
		defer server.Close()

		fileWithURL := &models.File{
			ID:       "F2",
			Name:     "test2.jpg",
			Mimetype: "image/jpeg",
			FileURL:  server.URL,
		}

		mockRepo.On("GetFileById", mock.Anything, "F2").Return(fileWithURL, nil).Once()

		b64, err := app.GetBase64Image("F2")

		assert.NoError(t, err)
		assert.Equal(t, "data:image/jpeg;base64,ZG93bmxvYWRlZCBkYXRh", b64)
		mockRepo.AssertExpectations(t)

		// Check if the file was saved
		_, err = os.Stat(filepath.Join(tempDir, "F2_test2.jpg"))
		assert.NoError(t, err)
	})
}
