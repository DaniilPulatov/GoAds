package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/ad"
	"context"
)

type AdminAdvertisementService interface {
	GetAllAds(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error)
	GetStatistics(ctx context.Context) (entities.AdStatistics, error)
	DeleteAd(ctx context.Context, adID int, adminID string) error
	DeleteImage(ctx context.Context, adID int, imageID int, adminID string) error
	Approve(ctx context.Context, adID int, adminID string) error
	Reject(ctx context.Context, adID int, adminID string) error
}

type service struct {
	adRepo ad.AdRepository
}

func NewAdminService(adRepo ad.AdRepository) AdminAdvertisementService {
	return &service{
		adRepo: adRepo,
	}
}
