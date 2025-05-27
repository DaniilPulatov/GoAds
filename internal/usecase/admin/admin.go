package admin

import (
	"ads-service/internal/domain/entities"
	usecaseerr "ads-service/internal/errs/usecaseerr"
	"context"
	"log"
	"time"
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

func (s *service) DeleteImage(ctx context.Context, adID, imageID int, adminID string) error {
	return nil // TODO: implement
}

func (s *service) Approve(ctx context.Context, adID int) error {
	repoAd, err := s.adRepo.GetByID(ctx, adID)
	if err != nil {
		log.Printf("failed to get ad by id %d: %v", adID, err)
		return usecaseerr.ErrGettingAdByID
	}
	if repoAd == nil {
		log.Printf("ad with id %d not found", adID)
		return usecaseerr.ErrGettingAdByID
	}

	repoAd.Status = entities.StatusApproved
	repoAd.IsActive = true
	repoAd.UpdatedAt = time.Now().UTC()

	if err = s.adRepo.Approve(ctx, adID, repoAd); err != nil {
		log.Printf("failed to approve ad id %d: %v", adID, err)
		return usecaseerr.ErrApprovingAd
	}

	return nil
}

func (s *service) Reject(ctx context.Context, adID int, reason string) error {
	repoAd, err := s.adRepo.GetByID(ctx, adID)
	if err != nil {
		log.Printf("failed to get ad by id %d: %v", adID, err)
		return usecaseerr.ErrGettingAdByID
	}
	if repoAd == nil {
		log.Printf("ad with id %d not found", adID)
		return usecaseerr.ErrGettingAdByID
	}

	repoAd.Status = entities.StatusRejected
	repoAd.IsActive = false
	repoAd.UpdatedAt = time.Now().UTC()
	repoAd.RejectionReason = reason

	if err = s.adRepo.Reject(ctx, adID, repoAd); err != nil {
		log.Printf("failed to reject ad id %d: %v", adID, err)
		return usecaseerr.ErrRejectingAd
	}

	return nil
}

func (s *service) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	statistics, err := s.adRepo.GetStatistics(ctx)
	if err != nil {
		log.Println("error getting statistics:", err)
		return entities.AdStatistics{}, usecaseerr.ErrGettingStatistics
	}
	return statistics, nil
}
