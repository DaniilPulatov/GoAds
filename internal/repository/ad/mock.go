//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package ad

import (
	"ads-service/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAdRepo struct {
	mock.Mock
}

type MockAdFileRepo struct {
	mock.Mock
}

func (m *MockAdFileRepo) Create(ctx context.Context, file *entities.AdFile) (int, error) {
	args := m.Called(ctx, file)
	if id, ok := args.Get(0).(int); ok {
		return id, args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockAdFileRepo) GetAll(ctx context.Context, adID int) ([]entities.AdFile, error) {
	args := m.Called(ctx, adID)
	if files, ok := args.Get(0).([]entities.AdFile); ok {
		return files, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAdFileRepo) Delete(ctx context.Context, file *entities.AdFile) (string, error) {
	args := m.Called(ctx, file)
	if url, ok := args.Get(0).(string); ok {
		return url, args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockAdRepo) Create(ctx context.Context, ad *entities.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

func (m *MockAdRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Ad), args.Error(1)
}

func (m *MockAdRepo) GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockAdRepo) GetAll(ctx context.Context) ([]entities.Ad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockAdRepo) Update(ctx context.Context, ad *entities.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

func (m *MockAdRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAdRepo) Approve(ctx context.Context, id int, ad *entities.Ad) error {
	args := m.Called(ctx, id, ad)
	return args.Error(0)
}

func (m *MockAdRepo) Reject(ctx context.Context, id int, ad *entities.Ad) error {
	args := m.Called(ctx, id, ad)
	return args.Error(0)
}

func (m *MockAdRepo) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(entities.AdStatistics), args.Error(1)
}
