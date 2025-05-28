//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package user

import (
	"ads-service/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateDraft(ctx context.Context, userID string, ad *entities.Ad) error {
	args := m.Called(ctx, userID, ad)
	return args.Error(0)
}

func (m *MockUserService) GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error) {
	args := m.Called(ctx, userID)
	if ads, ok := args.Get(0).([]entities.Ad); ok {
		return ads, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) UpdateMyAd(ctx context.Context, userID string, ad *entities.Ad) error {
	args := m.Called(ctx, userID, ad)
	return args.Error(0)
}

func (m *MockUserService) DeleteMyAd(ctx context.Context, userID string, adID int) error {
	args := m.Called(ctx, userID, adID)
	return args.Error(0)
}

func (m *MockUserService) SubmitForModeration(ctx context.Context, userID string, adID int) error {
	args := m.Called(ctx, userID, adID)
	return args.Error(0)
}

func (m *MockUserService) AddImageToMyAd(ctx context.Context, userID string, file *entities.AdFile) error {
	args := m.Called(ctx, userID, file)
	return args.Error(0)
}

func (m *MockUserService) GetImagesToMyAd(ctx context.Context, userID string, adID int) ([]entities.AdFile, error) {
	args := m.Called(ctx, userID, adID)
	if files, ok := args.Get(0).([]entities.AdFile); ok {
		return files, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error {
	args := m.Called(ctx, userID, file)
	return args.Error(0)
}

var _ UserAdvertisementService = (*MockUserService)(nil)
