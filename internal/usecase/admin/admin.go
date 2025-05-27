package admin

import (
	"ads-service/internal/domain/entities"
	usecaseerr "ads-service/internal/errs/usecaseErr"
	"context"
	"log"
)

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
