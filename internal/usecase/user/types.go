package user

import (
	"ads-service/internal/domain/entities"
	adRepo "ads-service/internal/repository/ad"
	adfile "ads-service/internal/repository/adFile"
	customLogger "ads-service/pkg/logger"
	"context"
)

type UserAdvertisementService interface {
	CreateDraft(ctx context.Context, userID string, ad *entities.Ad) error
	GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error)
	UpdateMyAd(ctx context.Context, userID string, ad *entities.Ad) error
	DeleteMyAd(ctx context.Context, userID string, adID int) error
	SubmitForModeration(ctx context.Context, userID string, adID int) error
	AddImageToMyAd(ctx context.Context, userID string, file *entities.AdFile) error
	GetImagesToMyAd(ctx context.Context, userID string, adID int) ([]entities.AdFile, error)
	DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error
}

type service struct {
	repo     adRepo.AdRepository
	fileRepo adfile.AdFileRepository
	logger   customLogger.Logger
}

func NewUserService(repo adRepo.AdRepository, fileRepo adfile.AdFileRepository, logTool customLogger.Logger) UserAdvertisementService {
	return &service{
		repo:     repo,
		fileRepo: fileRepo,
		logger:   logTool,
	}
}
