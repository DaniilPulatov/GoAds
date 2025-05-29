package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/internal/repository/ad"
	"ads-service/internal/repository/user"
	customLogger "ads-service/pkg/logger"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMockAdminService_GetAllAds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		defer mockUserRepo.AssertExpectations(t)

		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})
		expectedAds := []entities.Ad{
			{ID: 1, AuthorID: "1", Title: "ad1"},
			{ID: 2, AuthorID: "1", Title: "ad2"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedAds, nil)

		ads, err := service.GetAllAds(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedAds, ads)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		defer mockUserRepo.AssertExpectations(t)

		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetAll", mock.Anything).
			Return([]entities.Ad{}, repoerr.ErrSelection)

		ads, err := service.GetAllAds(context.Background())
		assert.Error(t, err)
		assert.Nil(t, ads)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		defer mockUserRepo.AssertExpectations(t)

		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetAll", mock.Anything).Return([]entities.Ad{}, nil)

		ads, err := service.GetAllAds(context.Background())
		assert.Error(t, err)
		assert.Nil(t, ads)
	})
}

func TestMockAdminService_DeleteAd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		defer mockUserRepo.AssertExpectations(t)

		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})
		mockRepo.On("Delete", mock.Anything, 1).Return(nil)

		err := service.DeleteAd(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("invalid ad id", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		err := service.DeleteAd(context.Background(), 0)
		assert.Error(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("Delete", mock.Anything, 2).Return(assert.AnError)

		err := service.DeleteAd(context.Background(), 2)
		assert.Error(t, err)
	})
}

func TestMockAdminService_Approve(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 1}
		mockRepo.On("GetByID", mock.Anything, 1).Return(adEntity, nil)
		mockRepo.On("Approve", mock.Anything, 1, mock.AnythingOfType("*entities.Ad")).Return(nil)

		err := service.Approve(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 2).Return(nil, nil)

		err := service.Approve(context.Background(), 2)
		assert.Error(t, err)
	})

	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 3).Return(nil, assert.AnError)

		err := service.Approve(context.Background(), 3)
		assert.Error(t, err)
	})

	t.Run("approve error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 4}
		mockRepo.On("GetByID", mock.Anything, 4).Return(adEntity, nil)
		mockRepo.On("Approve", mock.Anything, 4, mock.AnythingOfType("*entities.Ad")).Return(assert.AnError)

		err := service.Approve(context.Background(), 4)
		assert.Error(t, err)
	})
}

func TestMockAdminService_Reject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 1}
		mockRepo.On("GetByID", mock.Anything, 1).Return(adEntity, nil)
		mockRepo.On("Reject", mock.Anything, 1, mock.AnythingOfType("*entities.Ad")).Return(nil)

		err := service.Reject(context.Background(), 1, "bad")
		assert.NoError(t, err)
	})

	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 2).Return(nil, nil)

		err := service.Reject(context.Background(), 2, "bad")
		assert.Error(t, err)
	})

	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 3).Return(nil, assert.AnError)

		err := service.Reject(context.Background(), 3, "bad")
		assert.Error(t, err)
	})

	t.Run("reject error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 4}
		mockRepo.On("GetByID", mock.Anything, 4).Return(adEntity, nil)
		mockRepo.On("Reject", mock.Anything, 4, mock.AnythingOfType("*entities.Ad")).Return(assert.AnError)

		err := service.Reject(context.Background(), 4, "bad")
		assert.Error(t, err)
	})
}

func TestMockAdminService_GetStatistics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		expectedStats := entities.AdStatistics{Total: 10}
		mockRepo.On("GetStatistics", mock.Anything).Return(expectedStats, nil)

		stats, err := service.GetStatistics(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockUserRepo := user.MockUserRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewAdminService(&mockRepo, &mockUserRepo, customLogger.Logger{})

		mockRepo.On("GetStatistics", mock.Anything).Return(entities.AdStatistics{}, assert.AnError)

		stats, err := service.GetStatistics(context.Background())
		assert.Error(t, err)
		assert.Equal(t, entities.AdStatistics{}, stats)
	})
}
