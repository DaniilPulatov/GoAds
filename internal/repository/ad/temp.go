package ad

/*
package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/internal/errs/usecaseerr"
	"ads-service/internal/repository/ad"
	"ads-service/internal/repository/adFile"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
	"time"
)

type mockFileRepo struct {
	adFile.AdFileRepository
	mock.Mock
}

func (m *mockFileRepo) Create(ctx context.Context, file *entities.AdFile) (int, error) {
	args := m.Called(ctx, file)
	return args.Int(0), args.Error(1)
}
func (m *mockFileRepo) Delete(ctx context.Context, file *entities.AdFile) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}
func (m *mockFileRepo) GetAll(ctx context.Context, adID int) ([]entities.AdFile, error) {
	args := m.Called(ctx, adID)
	return args.Get(0).([]entities.AdFile), args.Error(1)
}

func TestService_CreateDraft(t *testing.T) {
	t.Run("title is empty", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		err := service.CreateDraft(context.Background(), "1", &entities.Ad{Title: ""})
		assert.Error(t, err)
		assert.Equal(t, usecaseerr.ErrInvalidParams, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(repoerr.ErrInsert)
		err := service.CreateDraft(context.Background(), "1", &entities.Ad{Title: "ok", Description: "desc"})
		assert.Error(t, err)
		assert.Equal(t, repoerr.ErrInsert, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		err := service.CreateDraft(context.Background(), "1", &entities.Ad{Title: "ok", Description: "desc"})
		assert.NoError(t, err)
	})
}

func TestService_GetMyAds(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByUserID", mock.Anything, "1").Return(nil, errors.New("err"))
		ads, err := service.GetMyAds(context.Background(), "1")
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAdsByUserID, err)
	})

	t.Run("no ads", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByUserID", mock.Anything, "1").Return([]entities.Ad{}, nil)
		ads, err := service.GetMyAds(context.Background(), "1")
		assert.Nil(t, ads)
		assert.Equal(t, usecaseerr.ErrUserNotHaveAds, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByUserID", mock.Anything, "1").Return([]entities.Ad{{ID: 1}}, nil)
		ads, err := service.GetMyAds(context.Background(), "1")
		assert.NoError(t, err)
		assert.Len(t, ads, 1)
	})
}

func TestService_UpdateMyAd(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1})
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1})
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("invalid params", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1, Title: ""})
		assert.Equal(t, usecaseerr.ErrInvalidParams, err)
	})

	t.Run("repo update error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(repoerr.ErrUpdate)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1, Title: "ok", Description: "desc"})
		assert.Equal(t, repoerr.ErrUpdate, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		err := service.UpdateMyAd(context.Background(), "1", &entities.Ad{ID: 1, Title: "ok", Description: "desc"})
		assert.NoError(t, err)
	})
}

func TestService_DeleteMyAd(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("repo delete error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Delete", mock.Anything, 1).Return(repoerr.ErrDelete)
		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.Equal(t, repoerr.ErrDelete, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Delete", mock.Anything, 1).Return(nil)
		err := service.DeleteMyAd(context.Background(), "1", 1)
		assert.NoError(t, err)
	})
}

func TestService_SubmitForModeration(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Equal(t, usecaseerr.ErrGettingAdByID, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("repo update error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(usecaseerr.ErrApprovingAd)
		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.Equal(t, usecaseerr.ErrApprovingAd, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		defer mockRepo.AssertExpectations(t)
		service := NewUserService(&mockRepo)
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		err := service.SubmitForModeration(context.Background(), "1", 1)
		assert.NoError(t, err)
	})
}

func TestService_AddImageToMyAd(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.jpg"})
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.jpg"})
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file not allowed", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.exe"})
		assert.Equal(t, usecaseerr.ErrFileNotAllowed, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Create", mock.Anything, mock.Anything).Return(-1, repoerr.ErrFileInsertion)
		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.jpg"})
		assert.Equal(t, repoerr.ErrFileInsertion, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Create", mock.Anything, mock.Anything).Return(1, nil)
		err := service.AddImageToMyAd(context.Background(), "1", &entities.AdFile{AdID: 1, FileName: "file.jpg"})
		assert.NoError(t, err)
	})
}

func TestService_DeleteMyAdImage(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Delete", mock.Anything, mock.Anything).Return("", repoerr.ErrFileDeletion)
		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.Equal(t, repoerr.ErrFileDeletion, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("Delete", mock.Anything, mock.Anything).Return("file.jpg", nil)
		err := service.DeleteMyAdImage(context.Background(), "1", &entities.AdFile{AdID: 1})
		assert.NoError(t, err)
	})
}

func TestService_GetImagesToMyAd(t *testing.T) {
	t.Run("ad not found", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("err"))
		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("access denied", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "2"}, nil)
		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, usecaseerr.ErrAccessDenied, err)
	})

	t.Run("file repo error", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("GetAll", mock.Anything, 1).Return(nil, repoerr.ErrSelection)
		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := ad.MockAdRepo{}
		mockFile := new(mockFileRepo)
		defer mockRepo.AssertExpectations(t)
		service := &service{repo: &mockRepo, fileRepo: mockFile}
		mockRepo.On("GetByID", mock.Anything, 1).Return(&entities.Ad{AuthorID: "1"}, nil)
		mockFile.On("GetAll", mock.Anything, 1).Return([]entities.AdFile{{ID: 1}}, nil)
		files, err := service.GetImagesToMyAd(context.Background(), "1", 1)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
	})
}
*/
