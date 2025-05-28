//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package user

import (
	"ads-service/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	args := m.Called(ctx, user)
	if userID, ok := args.Get(0).(string); ok {
		return userID, args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	args := m.Called(ctx, userID)
	if user, ok := args.Get(0).(*entities.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) GetAllUser(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]entities.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) GetByPhone(ctx context.Context, phone string) (*entities.User, error) {
	args := m.Called(ctx, phone)
	if user, ok := args.Get(0).(*entities.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) IsExists(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	if exists, ok := args.Get(0).(bool); ok {
		return exists, args.Error(1)
	}
	return false, args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, ad *entities.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Ad), args.Error(1)
}

func (m *MockUserRepo) GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockUserRepo) GetAll(ctx context.Context) ([]entities.Ad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, ad *entities.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) Approve(ctx context.Context, id int, ad *entities.Ad) error {
	args := m.Called(ctx, id, ad)
	return args.Error(0)
}

func (m *MockUserRepo) Reject(ctx context.Context, id int, ad *entities.Ad) error {
	args := m.Called(ctx, id, ad)
	return args.Error(0)
}

func (m *MockUserRepo) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(entities.AdStatistics), args.Error(1)
}
