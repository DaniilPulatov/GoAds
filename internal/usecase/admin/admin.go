package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/usecaseErr"
	"ads-service/internal/repository/ad"
	"ads-service/internal/repository/user"
	"context"
	"log"
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

func (s *service) GetAllAds(ctx context.Context) ([]entities.Ad, error) {
	ads, err := s.adRepo.GetAll(ctx)
	if err != nil {
		log.Println("error getting all ads:", err)
		return nil, usecaseerr.ErrGettingAllAds
	}
	if len(ads) == 0 {
		return nil, usecaseerr.ErrNoAds
	}

	return ads, nil
}

func (s *service) DeleteAd(ctx context.Context, adID int) error {
	err := s.adRepo.Delete(ctx, adID)
	if err != nil {
		return usecaseerr.ErrDeletingAd
	}

	return nil
}

func (s *service) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	return entities.AdStatistics{}, nil // TODO: implement
}

func (s *service) DeleteImage(ctx context.Context, adID, imageID int, adminID string) error {
	return nil // TODO: implement
}

func (s *service) Approve(ctx context.Context, adID int, adminID string) error {
	return nil // TODO: implement
}

func (s *service) Reject(ctx context.Context, adID int, adminID string) error {
	return nil // TODO: implement
}

// Additional methods for the admin service can be added here as needed.

// This service will handle administrative tasks related to advertisements,
// such as managing ad statuses and retrieving statistics.

// The methods will interact with the ad repository to perform the necessary operations.

// IsAdmin should be called in middleware to check access
func (s *service) IsAdmin(ctx context.Context, userID string) (bool, error) {
	userByID, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, usecaseerr.ErrGettingUser
	}
	if userByID == nil {
		return false, usecaseerr.ErrUserNotFound
	}

	return userByID.Role == entities.RoleAdmin, nil
}
