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

func (m *MockRepository) CountPendingBinaries(ctx context.Context, importID int64) (int64, error) {
	return 10, nil
}

func (m *MockRepository) StartDownload(ctx context.Context, importID int64) error {
	return nil
}

func (m *MockRepository) FinishDownload(ctx context.Context, importID int64) error {
	return nil
}

// channels
func (m *MockRepository) GetChannels() (channels []models.Channel, err error) {
	panic("GetChannels not defined in mock")
}

func (m *MockRepository) BulkInsertChannels(ctx context.Context, channels []models.Channel) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRepository) MarkCheckpointAsCompleted(ctx context.Context, importID int64, path string, errMsg string) error {
	args := m.Called(ctx, importID, path, errMsg)
	return args.Error(0)
}

func (m *MockRepository) CreateCheckpointPlaceHolders(ctx context.Context, importID int64, paths []string) error {
	args := m.Called(ctx, importID, paths)
	return args.Error(0)
}

func (m *MockRepository) GetImport(ctx context.Context, id int64) (*models.SlackImport, error) {
	panic("GetImport not defined in mock")
}

func (m *MockRepository) FinishedImportUpdate(ctx context.Context, imp *models.SlackImport) error {
	panic("FinishedImportUpdate not defined in mock")
}

func (m *MockRepository) FinishedImportUpdateWithError(ctx context.Context, imp *models.SlackImport) error {
	panic("FinishedImportUpdateWithError not defined in mock")
}

func (m *MockRepository) DeleteImport(ctx context.Context, id int64) error {
	panic("DeleteImport not defined in mock")
}

func (m *MockRepository) RenameImport(ctx context.Context, importID int64, name string) error {
	panic("RenameImport not defined in mock")
}

func (m *MockRepository) CreateImport(ctx context.Context, path string, name string) (*models.SlackImport, error) {
	panic("CreateImport not defined in mock")
}

func (m *MockRepository) NoFinishedImports(ctx context.Context) (bool, error) {
	panic("NoFinishedImports not defined in mock")
}

func (m *MockRepository) ListImports(ctx context.Context) ([]models.SlackImport, error) {
	panic("ListImports not defined in mock")
}

func (m *MockRepository) ListPendingBinaries(context.Context, int64, int) ([]models.File, error) {
	return []models.File{}, nil
}

func (m *MockRepository) MarkFileAsDownloaded(ctx context.Context, importID int64, fileID string) error {
	return nil
}

func (m *MockRepository) ListPendingFiles(ctx context.Context, importID int64) ([]string, error) {
	args := m.Called(ctx, importID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) BulkInsertMessages(ctx context.Context, channelID string, messages []models.Message) error {
	args := m.Called(ctx, channelID, messages)
	return args.Error(0)
}

func (m *MockRepository) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
	panic("GetMessages not defined in mock")
}

func (m *MockRepository) ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error) {
	panic("ListChannelFiles not defined in mock")
}
func (m *MockRepository) GetThreadByTs(ctx context.Context, threadTs string, channelID string, importID int64) ([]models.Message, error) {
	panic("GetThreadByTs not defined in mock")
}
func (m *MockRepository) GetFileById(ctx context.Context, fileId string) (*models.File, error) {
	panic("GetFileById not defined in mock")
}

func (m *MockRepository) SetupSchemas() error {
	panic("SetupSchemas not defined in mock")
}

func (m *MockRepository) BulkInsertUsers(ctx context.Context, users []models.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockRepository) GetUserAvatar(ctx context.Context, userId string) (string, error) {
	panic("GetUserAvatar nor defined in mock")
}

func (m *MockRepository) StoreImport(ctx context.Context, path string) (*models.SlackImport, error) {
	panic("StoreImport not defined in mock")
}

func (m *MockRepository) FinishedImport(ctx context.Context, id int64) error {
	panic("FinishedImport not defined in mock")
}

func (m *MockRepository) GetUserHandle(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) SearchMessages(ctx context.Context, query string) ([]models.SearchResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SearchResult), args.Error(1)
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
	defer func() {
		err = tmpFile.Close()
		if err != nil {
			panic(err)
		}
	}()

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
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertMessages", mock.Anything, "C1", mock.Anything).Return(nil).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)

				zipPath := createTestZip(t, map[string]string{
					"users.json":                    validUsersJSON,
					"channels.json":                 validChannelsJSON,
					"C1/2022-01-01.json":            validMessagesJSON,
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
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
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
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
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
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
				zipPath := createTestZip(t, map[string]string{
					"users.json": validUsersJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to process users: db error",
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error processing channels",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
				zipPath := createTestZip(t, map[string]string{
					"users.json":    validUsersJSON,
					"channels.json": validChannelsJSON,
				})
				return zipPath, context.Background()
			},
			wantErr:       true,
			expectedError: "fail to process channels: db error",
			checkMock: func(t *testing.T, mockRepo *MockRepository) {
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything, mock.Anything)
				mockRepo.AssertCalled(t, "BulkInsertChannels", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error processing messages",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertMessages", mock.Anything, "C1", mock.Anything).Return(errors.New("db error")).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
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
				mockRepo.AssertCalled(t, "BulkInsertUsers", mock.Anything, mock.Anything)
				mockRepo.AssertCalled(t, "BulkInsertChannels", mock.Anything, mock.Anything)
				mockRepo.AssertCalled(t, "BulkInsertMessages", mock.Anything, "C1", mock.Anything)
			},
		},
		{
			name: "Context cancelled",
			setup: func(t *testing.T, mockRepo *MockRepository) (string, context.Context) {
				mockRepo.On("CreateCheckpointPlaceHolders", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("BulkInsertUsers", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("BulkInsertChannels", mock.Anything, mock.Anything).Return(nil).Once()
				mockRepo.On("MarkCheckpointAsCompleted", mock.Anything, mock.Anything, mock.Anything, "").Return(nil)
				mockRepo.On("ListPendingFiles", mock.Anything, mock.Anything).Return(nil, nil)
				// No call to BulkInsertMessages is expected as it should be cancelled before.

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
				defer func() {
					err := os.Remove(zipPath)
					if err != nil {
						panic(err)
					}
				}()
			}
			importIt := models.SlackImport{ID: 66, Path: zipPath}
			importer := NewImporter(ctx, mockRepo, &importIt)
			progressChan := make(chan models.Progress, 10)

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
