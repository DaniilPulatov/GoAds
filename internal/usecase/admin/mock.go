//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package admin

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) GetAllAds(ctx context.Context) ([]entities.Ad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Ad), args.Error(1)
}

func (m *MockAdminService) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(entities.AdStatistics), args.Error(1)
}

func (m *MockAdminService) DeleteAd(ctx context.Context, adID int) error {
	args := m.Called(ctx, adID)
	return args.Error(0)
}

/*
func (m *MockAdminService) DeleteFile(ctx context.Context, adID int, imageID int, adminID string) error {
	args := m.Called(ctx, adID, imageID, adminID)
	return args.Error(0)
}
*/

func (m *MockAdminService) Approve(ctx context.Context, adID int) error {
	args := m.Called(ctx, adID)
	return args.Error(0)
}

func (m *MockAdminService) Reject(ctx context.Context, adID int, reason string) error {
	args := m.Called(ctx, adID, reason)
	return args.Error(0)
}

var _ AdminAdvertisementService = (*MockAdminService)(nil)
