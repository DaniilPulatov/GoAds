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
	// DeleteFile(ctx context.Context, adID int, imageID int, adminID string) error
	Approve(ctx context.Context, adID int) error
	Reject(ctx context.Context, adID int, reason string) error
}

/*
	type FileDeleter interface {
		Delete(ctx context.Context, file *entities.AdFile) (url string, err error)
	}
*/
type service struct {
	// fileDel  FileDeleter
	adRepo   ad.AdRepository
	userRepo user.UserRepository
}

func NewAdminService(adRepo ad.AdRepository, userRepo user.UserRepository) AdminAdvertisementService {
	return &service{
		// fileDel: fileDel,
		adRepo:   adRepo,
		userRepo: userRepo,
	}
}
