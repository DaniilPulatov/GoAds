//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package auth

import (
	"ads-service/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *MockAuthService) Login(ctx context.Context, phone, password string) (string, string, error) {
	args := m.Called(ctx, phone, password)
	if accessToken, ok := args.Get(0).(string); ok {
		if refreshToken, ok := args.Get(1).(string); ok {
			return accessToken, refreshToken, args.Error(2)
		}
	}
	return "", "", args.Error(2)
}

func (m *MockAuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	args := m.Called(ctx, refreshToken)
	if accessToken, ok := args.Get(0).(string); ok {
		if newRefreshToken, ok := args.Get(1).(string); ok {
			return accessToken, newRefreshToken, args.Error(2)
		}
	}
	return "", "", args.Error(2)
}

func (m *MockAuthService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	if isAdmin, ok := args.Get(0).(bool); ok {
		return isAdmin, args.Error(1)
	}
	return false, args.Error(1)
}
