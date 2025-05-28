package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/usecaseerr"
	"context"
	"log"
	"os"
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
	if adID <= 0 {
		log.Println("invalid ad ID")
		return usecaseerr.ErrInvalidParams
	}

	err := s.adRepo.Delete(ctx, adID)
	if err != nil {
		log.Println("error deleting ad:", err)
		return usecaseerr.ErrDeletingAd
	}

	return nil
}

func (s *service) DeleteFile(ctx context.Context, file *entities.AdFile) (error) {
	url, err := s.fileDel.Delete(ctx, file)
	if err != nil{
		log.Println("error while deleteing image: ", err)
		return err // TODO: wrap error
	}
	if err := os.Remove(url); err != nil{
		log.Printf("error while deleteing file at url:%v\n%v\n", url, err)
		return err // TODO: wrap error
	}
	log.Println("image deleted successfully")
	return nil
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
