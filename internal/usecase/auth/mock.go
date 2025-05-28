//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package auth

import "github.com/stretchr/testify/mock"

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) IsAdmin(userID string) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}
