package importer

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the repository.Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SetupSchemas() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) NoFinishedImports(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) StoreImport(ctx context.Context, path string) (*models.Import, error) {
	args := m.Called(ctx, path)
	// Handle the case where the first return value might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Import), args.Error(1)
}

func (m *MockRepository) FinishedImport(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) BulkInsertUsers(users []models.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockRepository) BulkInsertChannels(channels []models.Channel) error {
	args := m.Called(channels)
	return args.Error(0)
}

func (m *MockRepository) BulkInsertMessages(ctx context.Context, ChannelID string, messages []models.Message) error {
	args := m.Called(ctx, ChannelID, messages)
	return args.Error(0)
}

func (m *MockRepository) GetChannels() ([]models.Channel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Channel), args.Error(1)
}

func (m *MockRepository) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
	args := m.Called(ctx, channelID, cursor, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockRepository) GetThreadByTs(ctx context.Context, threadTs string, channelID string) ([]models.Message, error) {
	args := m.Called(ctx, threadTs, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockRepository) GetUserAvatar(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetFileById(ctx context.Context, fileId string) (*models.File, error) {
	args := m.Called(ctx, fileId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockRepository) ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error) {
	args := m.Called(ctx, channelID, pageSize, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.File), args.Error(1)
}

// createTestZip creates a temporary zip file for testing purposes.
func createTestZip(t *testing.T, files map[string]string) string {
	t.Helper()
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for name, content := range files {
		f, err := zipWriter.Create(name)
		assert.NoError(t, err)
		_, err = f.Write([]byte(content))
		assert.NoError(t, err)
	}
	assert.NoError(t, zipWriter.Close())

	tmpFile, err := os.CreateTemp("", "test-*.zip")
	assert.NoError(t, err)
	_, err = tmpFile.Write(buf.Bytes())
	assert.NoError(t, err)
	assert.NoError(t, tmpFile.Close())

	return tmpFile.Name()
}

func TestNewImporter(t *testing.T) {
	mockRepo := new(MockRepository)
	importer := NewImporter(context.Background(), "/path/to/file", mockRepo)
	assert.NotNil(t, importer)
	assert.Equal(t, "/path/to/file", importer.path)
	assert.Equal(t, mockRepo, importer.repository)
}

func TestImport(t *testing.T) {
	mockRepo := new(MockRepository)

	users := []models.User{{ID: "U1", Name: "user1"}}
	channels := []models.Channel{{ID: "C1", Name: "channel1"}}
	messages := []models.Message{{Text: "Hello"}}

	usersJSON, _ := json.Marshal(users)
	channelsJSON, _ := json.Marshal(channels)
	messagesJSON, _ := json.Marshal(messages)

	files := map[string]string{
		"users.json":         string(usersJSON),
		"channels.json":      string(channelsJSON),
		"C1/2025-08-26.json": string(messagesJSON),
	}
	zipPath := createTestZip(t, files)
	defer os.Remove(zipPath)

	importer := NewImporter(context.Background(), zipPath, mockRepo)
	progressChan := make(chan Progress, 10)

	mockRepo.On("BulkInsertUsers", mock.AnythingOfType("[]models.User")).Return(nil).Once()
	mockRepo.On("BulkInsertChannels", mock.AnythingOfType("[]models.Channel")).Return(nil).Once()
	mockRepo.On("BulkInsertMessages", mock.Anything, "C1", mock.AnythingOfType("[]models.Message")).Return(nil).Once()

	err := importer.Import(progressChan)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)

	// Check progress channel
	progress, ok := <-progressChan
	assert.True(t, ok)
	assert.Equal(t, "C1 / 2025-08-26.json", progress.FileName)
	_, ok = <-progressChan
	assert.False(t, ok, "Channel should be closed")
}

func TestImport_MissingFiles(t *testing.T) {
	mockRepo := new(MockRepository)

	t.Run("Missing users.json", func(t *testing.T) {
		zipPath := createTestZip(t, map[string]string{"channels.json": "[]"})
		defer os.Remove(zipPath)
		importer := NewImporter(context.Background(), zipPath, mockRepo)
		err := importer.Import(make(chan Progress, 1))
		assert.ErrorContains(t, err, "couldnt find file users.json")
	})

	t.Run("Missing channels.json", func(t *testing.T) {
		zipPath := createTestZip(t, map[string]string{"users.json": "[]"})
		defer os.Remove(zipPath)
		importer := NewImporter(context.Background(), zipPath, mockRepo)
		mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
		err := importer.Import(make(chan Progress, 1))
		assert.ErrorContains(t, err, "couldnt find file channels.json")
		mockRepo.AssertExpectations(t)
	})
}

func TestImport_ContextCancelled(t *testing.T) {
	mockRepo := new(MockRepository)
	zipPath := createTestZip(t, map[string]string{
		"users.json":         "[]",
		"channels.json":      "[]",
		"C1/2025-08-26.json": "[]",
	})
	defer os.Remove(zipPath)

	ctx, cancel := context.WithCancel(context.Background())
	importer := NewImporter(ctx, zipPath, mockRepo)

	mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil)
	mockRepo.On("BulkInsertChannels", mock.Anything).Return(nil)

	// Cancel context right away
	cancel()

	err := importer.Import(make(chan Progress, 1))
	assert.ErrorIs(t, err, context.Canceled)
}

func TestProcessFunctions(t *testing.T) {
	mockRepo := new(MockRepository)
	importer := NewImporter(context.Background(), "", mockRepo)

	t.Run("processUsers error", func(t *testing.T) {
		mockRepo.On("BulkInsertUsers", mock.Anything).Return(errors.New("db error")).Once()
		// Create a dummy zip file to test processUsers
		zipPath := createTestZip(t, map[string]string{"users.json": "[{}]"})
		defer os.Remove(zipPath)
		r, _ := zip.OpenReader(zipPath)
		defer r.Close()
		f, _ := findFile("users.json", r.File)
		err := importer.processUsers(f)
		assert.ErrorContains(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestFindFile(t *testing.T) {
	zipPath := createTestZip(t, map[string]string{"test.txt": "hello"})
	defer os.Remove(zipPath)

	r, err := zip.OpenReader(zipPath)
	assert.NoError(t, err)
	defer r.Close()

	t.Run("File found", func(t *testing.T) {
		file, err := findFile("test.txt", r.File)
		assert.NoError(t, err)
		assert.Equal(t, "test.txt", filepath.Base(file.Name))
	})

	t.Run("File not found", func(t *testing.T) {
		_, err := findFile("notfound.txt", r.File)
		assert.Error(t, err)
	})
}
