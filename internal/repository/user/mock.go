//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package user

import (
	"ads-service/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAdRepo struct {
	mock.Mock
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
