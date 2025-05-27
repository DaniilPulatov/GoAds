//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package adfile

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAdFileRepository struct {
	mock.Mock
}

func (m *MockAdFileRepository) Create(ctx context.Context, file *entities.AdFile) (int, error) {
	args := m.Called(ctx, file)
	return args.Int(0), args.Error(1)
}

func (m *MockAdFileRepository) Delete(ctx context.Context, file *entities.AdFile) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

func (m *MockAdFileRepository) GetAll(ctx context.Context, adID int) ([]entities.AdFile, error) {
	args := m.Called(ctx, adID)
	return args.Get(0).([]entities.AdFile), args.Error(1)
}
