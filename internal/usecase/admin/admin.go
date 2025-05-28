package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/usecaseerr"
	"context"
	"time"
)

func (s *service) GetAllAds(ctx context.Context) ([]entities.Ad, error) {
	ads, err := s.adRepo.GetAll(ctx)
	if err != nil {
		s.logger.ERROR("error getting all ads:", err)
		return nil, usecaseerr.ErrGettingAllAds
	}
	if len(ads) == 0 {
		s.logger.ERROR("no ads found ", usecaseerr.ErrNoAds)
		return nil, usecaseerr.ErrNoAds
	}
	s.logger.INFO("all ads retrivied successfully")
	return ads, nil
}

func (s *service) DeleteAd(ctx context.Context, adID int) error {
	if adID <= 0 {
		s.logger.ERROR("invalid ad ID")
		return usecaseerr.ErrInvalidParams
	}

	err := s.adRepo.Delete(ctx, adID)
	if err != nil {
		s.logger.ERROR("error deleting ad:", err)
		return usecaseerr.ErrDeletingAd
	}
	s.logger.INFO("ad deleted successfully")
	return nil
}

/*
	func (s *service) DeleteFile(ctx context.Context, file *entities.AdFile) (error) {
		url, err := s.fileDel.Delete(ctx, file)
		if err != nil{
			s.logger.ERROR("error while deleteing image: ", err)
			return err // TODO: wrap error
		}
		if err := os.Remove(url); err != nil{
			log.Printf("error while deleteing file at url:%v\n%v\n", url, err)
			return err // TODO: wrap error
		}
		s.logger.ERROR("image deleted successfully")
		return nil
	}
*/
func (s *service) Approve(ctx context.Context, adID int) error {
	repoAd, err := s.adRepo.GetByID(ctx, adID)
	if err != nil {
		s.logger.ERROR("error getting ad:", err)
		return usecaseerr.ErrGettingAdByID
	}
	if repoAd == nil {
		s.logger.ERROR("no ad found ", usecaseerr.ErrNoAds)
		return usecaseerr.ErrGettingAdByID
	}

	repoAd.Status = entities.StatusApproved
	repoAd.IsActive = true
	repoAd.UpdatedAt = time.Now().UTC()

	if err = s.adRepo.Approve(ctx, adID, repoAd); err != nil {
		s.logger.ERROR("error approving ad:", err)
		return usecaseerr.ErrApprovingAd
	}
	s.logger.INFO("ad approved successfully")
	return nil
}

func (s *service) Reject(ctx context.Context, adID int, reason string) error {
	repoAd, err := s.adRepo.GetByID(ctx, adID)
	if err != nil {
		s.logger.ERROR("error getting ad:", err)
		return usecaseerr.ErrGettingAdByID
	}
	if repoAd == nil {
		s.logger.ERROR("no ad found ", usecaseerr.ErrGettingAdByID)
		return usecaseerr.ErrGettingAdByID
	}

	repoAd.Status = entities.StatusRejected
	repoAd.IsActive = false
	repoAd.UpdatedAt = time.Now().UTC()
	repoAd.RejectionReason = reason

	if err = s.adRepo.Reject(ctx, adID, repoAd); err != nil {
		s.logger.ERROR("error rejecting ad:", err)
		return usecaseerr.ErrRejectingAd
	}
	s.logger.INFO("ad rejected successfully")
	return nil
}

func (s *service) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	statistics, err := s.adRepo.GetStatistics(ctx)
	if err != nil {
		s.logger.ERROR("error getting statistics:", err)
		return entities.AdStatistics{}, usecaseerr.ErrGettingStatistics
	}
	return statistics, nil
}
