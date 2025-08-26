package app

import (
	"context"

	"github.com/mestr-io/slack-export-viewer/backend/models"

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
