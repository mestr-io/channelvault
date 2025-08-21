package importer

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	//"github.com/mestr-io/slack-export-viewer/backend/repository"
)

// MockRepository is a mock implementation of the repository.Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BulkInsertUsers(users []models.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockRepository) BulkInsertChannels(channels []models.Channel) error {
	args := m.Called(channels)
	return args.Error(0)
}

func (m *MockRepository) BulkInserMessages(ctx context.Context, channelID string, messages []models.Message) error {
	args := m.Called(ctx, channelID, messages)
	return args.Error(0)
}

//channels	
func (m *MockRepository) GetChannels() (channels []models.Channel, err error) {
	panic("GetChannels not defined in mock")
}

//imports
func (m *MockRepository) StoreImport(ctx context.Context, path string) (*models.Import, error) {
	panic("StoreImport not defined in mock")	
}
func (m *MockRepository) FinishedImport(ctx context.Context, ID int64) error {
	panic("FinishedImport not defined in mock")
}
func (m *MockRepository) NoFinishedImports(ctx context.Context) (bool, error) {
	panic("NotFinishedImports not defined in mock")
}

//messages
func (m *MockRepository) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
	panic("GetMessages not defined in mock")
}

func (m *MockRepository) ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error) {
	panic("ListChannelFiles not defined in mock")
}
func (m *MockRepository) GetThreadByTs(ctx context.Context, threadTs string, channelID string) ([]models.Message, error) {
	panic("GetThreadByTs not defined in mock")
}
func (m *MockRepository) GetFileById(ctx context.Context, fileId string) (*models.File, error) {
	panic("GetFileById not defined in mock")
}

//setup
func (m *MockRepository) SetupSchemas() error {
	panic("SetupSchemas not defined in mock")
}

//users	
func (m *MockRepository) GetUserAvatar(ctx context.Context, userId string) (string, error) {
	panic("GetUserAvatar nor defined in mock")
}

// Helper function to create a test zip file in memory.
func createTestZip(t *testing.T, files map[string]string) string {
	t.Helper()
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for name, content := range files {
		f, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("Failed to create file in zip: %v", err)
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			t.Fatalf("Failed to write to file in zip: %v", err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}

	// Create a temporary file to write the zip content to.
	tmpFile, err := os.CreateTemp("", "test-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Return the path to the temporary file.
	return tmpFile.Name()
}


func TestImporter_Imports(t *testing.T) {
	validUsersJSON := `[{"id": "U1", "name": "user1"}]`
	validChannelsJSON := `[{"id": "C1", "name": "channel1"}]`
	validMessagesJSON := `[{"user": "U1", "text": "hello", "ts": "1620000000.000000"}]`

	tests := []struct {
		name          string
		setup         func(t *testing.T, mockRepo *MockRepository) (path string, ctx context.Context)
		wantErr       bool
		expectedError string
		checkMock     func(t *testing.T, mockRepo *MockRepository)
	}{
		{
			name: "Successful import",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInserMessages", mock.Anything, "C1", mock.Anything).Return(nil).Once()

				zipPath := createTestZip(t, map[string]string{
					"users.json":             validUsersJSON,
					"channels.json":          validChannelsJSON,
					"C1/2022-01-01.json":     validMessagesJSON,
					"ignore_this_dir/some_file.txt": "data",
				})
				return zipPath, context.Background()
			},
			wantErr: false,
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertExpectations(t)
			},
		},
		{
			name: "Invalid zip path",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				return "nonexistent.zip", context.Background()
			},
			wantErr:       true,
			expectedError: "fail to open the zip export file",
		},
		{
			name: "Missing users.json",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				zipPath := createTestZip(t, map[string]string{
					"channels.json": validChannelsJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to open unzip users.json file: FATAL: couldnt find file users.json",
		},
		{
			name: "Missing channels.json",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
				zipPath := createTestZip(t, map[string]string{
					"users.json": validUsersJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to open unzip channels.json file: FATAL: couldnt find file channels.json",
		},
		{
			name: "Error processing users",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(errors.New("db error")).Once()
				zipPath := createTestZip(t, map[string]string{
					"users.json": validUsersJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to process users: db error",
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything)
			},
		},
		{
			name: "Error processing channels",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything).Return(errors.New("db error")).Once()
				zipPath := createTestZip(t, map[string]string{
					"users.json":    validUsersJSON,
					"channels.json": validChannelsJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to process channels: db error",
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything)
				mockRepo.AssertCalled(t, "BulkInsertChannels", mock.Anything)
			},
		},
		{
			name: "Error processing messages",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInserMessages", mock.Anything, "C1", mock.Anything).Return(errors.New("db error")).Once()
				zipPath := createTestZip(t, map[string]string{
					"users.json":         validUsersJSON,
					"channels.json":      validChannelsJSON,
					"C1/2022-01-01.json": validMessagesJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to process messages for channel C1: db error",
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything)
				mockRepo.AssertCalled(t, "BulkInsertChannels", mock.Anything)
				mockRepo.AssertCalled(t, "BulkInserMessages", mock.Anything, "C1", mock.Anything)
			},
		},
		{
			name: "Context cancelled",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("BulkInsertUsers", mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything).Return(nil).Once()
				// No call to BulkInserMessages is expected as it should be cancelled before.

				zipPath := createTestZip(t, map[string]string{
					"users.json":         validUsersJSON,
					"channels.json":      validChannelsJSON,
					"C1/2022-01-01.json": validMessagesJSON,
				})
				ctx, cancel := context.WithCancel(context.Background())
				// Cancel the context immediately
				cancel()
				return zipPath, ctx
			},
			wantErr:       true,
			expectedError: context.Canceled.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			zipPath, ctx := tt.setup(t, mockRepo)
			if zipPath != "nonexistent.zip" {
				defer os.Remove(zipPath) // Clean up the temp file
			}

			importer := NewImporter(ctx, zipPath, mockRepo)
			progressChan := make(chan Progress, 10)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Drain the channel to prevent the sender from blocking
				for range progressChan {
				}
			}()

			err := importer.Import(progressChan)

			wg.Wait() // Wait for the channel drainer to finish

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.checkMock != nil {
				tt.checkMock(t, mockRepo)
			}
		})
	}
}

func TestFindFile(t *testing.T) {
	files := []*zip.File{
		{FileHeader: zip.FileHeader{Name: "users.json"}},
		{FileHeader: zip.FileHeader{Name: "channel1/messages.json"}},
	}

	t.Run("finds existing file", func(t *testing.T) {
		f, err := findFile("users.json", files)
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, "users.json", filepath.Base(f.Name))
	})

	t.Run("finds existing file in subdir", func(t *testing.T) {
		f, err := findFile("messages.json", files)
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, "messages.json", filepath.Base(f.Name))
	})


	t.Run("returns error for non-existing file", func(t *testing.T) {
		f, err := findFile("nonexistent.json", files)
		assert.Error(t, err)
		assert.Nil(t, f)
		assert.Equal(t, "FATAL: couldnt find file nonexistent.json", err.Error())
	})
}
