//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package auth
import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// Token — копия нужной структуры
type Token struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

// AuthRepositoryMock — мок интерфейса AuthRepository
type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) Create(ctx context.Context, rtoken Token) error {
	args := m.Called(ctx, rtoken)
	return args.Error(0)
}

func (m *AuthRepositoryMock) Get(ctx context.Context, userID string) (*Token, error) {
	args := m.Called(ctx, userID)
	if token, ok := args.Get(0).(*Token); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AuthRepositoryMock) Update(ctx context.Context, rtoken Token) error {
	args := m.Called(ctx, rtoken)
	return args.Error(0)
}

func (m *AuthRepositoryMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}