package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/internal/errs/usecaseerr"
	"ads-service/internal/repository/ad"
	adfile "ads-service/internal/repository/adFile"
	customLogger "ads-service/pkg/logger"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateDraft(t *testing.T) {
	t.Run("title is empty", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		err := service.CreateDraft(context.Background(), "1", &entities.Ad{
			Title: "",
		})

		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrInvalidParams, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("Create", mock.Anything, mock.Anything).
			Return(repoerr.ErrInsert)
		err := service.CreateDraft(context.Background(), "1",
			&entities.Ad{Title: "ok", Description: "desc", CategoryID: 1})

		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrInsert, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		err := service.CreateDraft(context.Background(), "1",
			&entities.Ad{Title: "ok", Description: "desc", CategoryID: 1})
		assert.NoError(t, err)
	})
}

func TestService_UpdateMyAd(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("GetByID", mock.Anything, 1).
			Return((*entities.Ad)(nil), errors.New("err"))
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1})

		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1})

		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("invalid params", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1, Title: ""})
		assert.Equal(t, usecaseerr.ErrInvalidParams, err)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 1, AuthorID: "1", Title: "ok", Description: "desc", CategoryID: 1}
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(repoerr.ErrUpdate)
		err := service.UpdateMyAd(context.Background(), "1", adEntity)

		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrUpdate, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		adEntity := &entities.Ad{ID: 1, AuthorID: "1", Title: "ok", Description: "desc", CategoryID: 1}
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(nil)
		err := service.UpdateMyAd(context.Background(), "1", adEntity)

		assert.NoError(t, err)
	})
}

func TestService_DeleteMyAd(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 1).
			Return((*entities.Ad)(nil), errors.New("err"))

		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)

		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Delete", mock.Anything, 1).
			Return(errors.New("err"))

		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrDelete, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Delete", mock.Anything, 1).
			Return(nil)

		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.NoError(t, err)
	})
}

func TestService_GetMyAds(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByUserID", mock.Anything, "1").
			Return([]entities.Ad{}, errors.New("db error"))

		ads, err := service.GetMyAds(context.Background(), "1")
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAdsByUserID, err)
	})

	t.Run("user has no ads", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByUserID", mock.Anything, "1").
			Return([]entities.Ad{}, nil)

		ads, err := service.GetMyAds(context.Background(), "1")
		assert.Nil(t, ads)
		assert.Equal(t, usecaseerr.ErrUserNotHaveAds, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		expectedAds := []entities.Ad{
			{ID: 1, AuthorID: "1", Title: "ad1"},
			{ID: 2, AuthorID: "1", Title: "ad2"},
		}
		mockRepo.On("GetByUserID", mock.Anything, "1").
			Return(expectedAds, nil)

		ads, err := service.GetMyAds(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, expectedAds, ads)
	})
}

func TestService_SubmitForModeration(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		mockRepo.On("GetByID", mock.Anything, 1).
			Return((*entities.Ad)(nil), errors.New("err"))

		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)

		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 1, AuthorID: "1"}
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(adEntity, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(errors.New("err"))

		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrApprovingAd, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})

		adEntity := &entities.Ad{ID: 1, AuthorID: "1"}
		mockRepo.On("GetByID", mock.Anything, 1).
			Return(adEntity, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).
			Return(nil)

		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.NoError(t, err)
	})
}

func TestService_AddImageToMyAd(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return((*entities.Ad)(nil), errors.New("err"))

		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "img.jpg"})
		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)

		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "img.jpg"})
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file not allowed", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)

		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.exe"})
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrFileNotAllowed, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFileRepo.On("Create", mock.Anything, mock.Anything).
			Return(-1, repoerr.ErrFileInsertion)

		err := service.AddImageToMyAd(context.Background(), "1",
			&entities.AdFile{AdID: 1, FileName: "img.jpg"})
		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrFileInsertion, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFileRepo.On("Create", mock.Anything, mock.Anything).
			Return(1, nil)

		err := service.AddImageToMyAd(context.Background(), "1",
			&entities.AdFile{AdID: 1, FileName: "img.jpg"})
		assert.NoError(t, err)
	})
}

func TestService_GetImagesToMyAd(t *testing.T) {
	t.Run("get by id error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(nil, errors.New("db error"))

		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)

		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFileRepo.On("GetAll", mock.Anything, 1).
			Return([]entities.AdFile{}, errors.New("file error"))

		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFileRepo}
		expectedFiles := []entities.AdFile{
			{ID: 1, AdID: 1, FileName: "img1.jpg"},
			{ID: 2, AdID: 1, FileName: "img2.jpg"},
		}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFileRepo.On("GetAll", mock.Anything, 1).
			Return(expectedFiles, nil)

		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedFiles, files)
	})
}

func TestService_DeleteMyAdImage(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFile.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFile}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(nil, errors.New("err"))

		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFile.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFile}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "2"}, nil)

		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFile.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFile}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Delete", mock.Anything, mock.Anything).
			Return("", repoerr.ErrFileDeletion)

		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, repoerr.ErrFileDeletion, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFile.AssertExpectations(t)

		service := &service{repo: &mockRepo, fileRepo: &mockFile}

		mockRepo.On("GetByID", mock.Anything, 1).
			Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Delete", mock.Anything, mock.Anything).
			Return("file.jpg", nil)

		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.NoError(t, err)
	})
}

func TestService_GetMyAdsByFilter(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		filter := entities.AdFilter{}
		mockRepo.On("Filter", mock.Anything, mock.Anything).
			Return([]entities.Ad{}, errors.New("db error"))

		ads, err := service.GetMyAdsByFilter(context.Background(), "1", &filter)
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAdsByUserID, err)
	})

	t.Run("user has no ads", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		filter := entities.AdFilter{}
		mockRepo.On("Filter", mock.Anything, mock.Anything).
			Return([]entities.Ad{}, nil)

		ads, err := service.GetMyAdsByFilter(context.Background(), "1", &filter)
		assert.Nil(t, ads)
		assert.Equal(t, usecaseerr.ErrUserNotHaveAds, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFileRepo := adfile.MockAdFileRepository{}
		defer mockRepo.AssertExpectations(t)
		defer mockFileRepo.AssertExpectations(t)

		service := NewUserService(&mockRepo, &mockFileRepo, customLogger.Logger{})
		expectedAds := []entities.Ad{
			{ID: 1, AuthorID: "1", Title: "ad1"},
			{ID: 2, AuthorID: "1", Title: "ad2"},
		}
		filter := entities.AdFilter{}
		mockRepo.On("Filter", mock.Anything, mock.Anything).
			Return(expectedAds, nil)

		ads, err := service.GetMyAdsByFilter(context.Background(), "1", &filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedAds, ads)
	})
}
