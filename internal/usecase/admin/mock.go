//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package admin

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllAds(ctx context.Context) ([]entities.Ad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockService) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(entities.AdStatistics), args.Error(1)
}

func (m *MockService) DeleteAd(ctx context.Context, adID int) error {
	args := m.Called(ctx, adID)
	return args.Error(0)
}
/*
func (m *MockService) DeleteFile(ctx context.Context, adID int, imageID int, adminID string) error {
	args := m.Called(ctx, adID, imageID, adminID)
	return args.Error(0)
}
*/
func (m *MockService) Approve(ctx context.Context, adID int) error {
	args := m.Called(ctx, adID)
	return args.Error(0)
}

func (m *MockService) Reject(ctx context.Context, adID int, reason string) error {
	args := m.Called(ctx, adID, reason)
	return args.Error(0)
}

var _ AdminAdvertisementService = (*MockService)(nil)
