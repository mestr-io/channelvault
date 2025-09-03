package app

import (
	"context"
	"testing"
	"time"

	"github.com/mestr-io/slack-export-viewer/backend/models"

	"github.com/stretchr/testify/mock"

	"github.com/wailsapp/wails/v2/pkg/runtime" //para runtime.Quit
)

// MockRepository is a mock implementation of the repository.Repository interface.
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

func (m *MockRepository) GetChannels() ([]models.Channel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Channel), args.Error(1)
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
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SlackImport), args.Error(1)
}

func (m *MockRepository) MarkFileAsDownloaded(ctx context.Context, importID int64, fileID string) error {
	return nil
}

func (m *MockRepository) FinishedImportUpdate(ctx context.Context, imp *models.SlackImport) error {
	args := m.Called(ctx, imp)
	return args.Error(0)
}

func (m *MockRepository) FinishedImportUpdateWithError(ctx context.Context, imp *models.SlackImport) error {
	args := m.Called(ctx, imp)
	return args.Error(0)
}

func (m *MockRepository) DeleteImport(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) RenameImport(ctx context.Context, importID int64, name string) error {
	args := m.Called(ctx, importID, name)
	return args.Error(0)
}

func (m *MockRepository) CreateImport(ctx context.Context, path string, name string) (*models.SlackImport, error) {
	args := m.Called(ctx, path, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SlackImport), args.Error(1)
}

func (m *MockRepository) NoFinishedImports(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) ListImports(ctx context.Context) ([]models.SlackImport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SlackImport), args.Error(1)
}

func (m *MockRepository) ListPendingFiles(ctx context.Context, importID int64) ([]string, error) {
	args := m.Called(ctx, importID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) ListPendingBinaries(context.Context, int64, int) ([]models.File, error) {
	return []models.File{}, nil
}

func (m *MockRepository) BulkInsertMessages(ctx context.Context, ChannelID string, messages []models.Message) error {
	args := m.Called(ctx, ChannelID, messages)
	return args.Error(0)
}

func (m *MockRepository) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
	args := m.Called(ctx, channelID, cursor, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockRepository) ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error) {
	args := m.Called(ctx, channelID, pageSize, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.File), args.Error(1)
}

func (m *MockRepository) GetThreadByTs(ctx context.Context, threadTs string, channelID string) ([]models.Message, error) {
	args := m.Called(ctx, threadTs, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockRepository) GetFileById(ctx context.Context, fileId string) (*models.File, error) {
	args := m.Called(ctx, fileId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockRepository) SetupSchemas() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) BulkInsertUsers(ctx context.Context, users []models.User) error {
	args := m.Called(ctx, users)
	return args.Error(0)
}

func (m *MockRepository) GetUserAvatar(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) StoreImport(ctx context.Context, path string) (*models.SlackImport, error) {
	args := m.Called(ctx, path)
	// Handle the case where the first return value might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SlackImport), args.Error(1)
}

func (m *MockRepository) FinishedImport(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TearDownFakeApp(ctx context.Context, ch chan string) {
	//fmt.Printf("After %d second(s) ending test", afterSeconds)
	runtime.Quit(ctx)
	ch <- "test ends"
}

func WaitTestToEnd(t *testing.T, ch chan string) {
	select {
	case /*msg := */ <-ch:
		/*fmt.Println(msg)*/
	case <-time.After(2 * time.Second):
		t.Fatalf("Timeout")
	}
}
