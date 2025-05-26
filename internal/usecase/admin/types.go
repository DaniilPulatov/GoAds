package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/ad"
	"ads-service/internal/repository/user"
	"context"
)

type AdminAdvertisementService interface {
	GetAllAds(ctx context.Context) ([]entities.Ad, error)
	GetStatistics(ctx context.Context) (entities.AdStatistics, error)
	DeleteAd(ctx context.Context, adID int) error
	DeleteImage(ctx context.Context, adID int, imageID int, adminID string) error
	Approve(ctx context.Context, adID int, adminID string) error
	Reject(ctx context.Context, adID int, adminID string) error
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

type service struct {
	adRepo   ad.AdRepository
	userRepo user.UserRepository
}

func NewAdminService(adRepo ad.AdRepository, userRepo user.UserRepository) AdminAdvertisementService {
	return &service{
		adRepo:   adRepo,
		userRepo: userRepo,
	}
}
