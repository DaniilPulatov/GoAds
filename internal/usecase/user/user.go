package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/internal/errs/usecaseerr"
	"ads-service/pkg/utils"
	"time"

	"context"
	"fmt"
	"os"
	"strings"
)

const fileDirPerm = 0o750 // Permissions for the directory where ad images are stored

func (s *service) CreateDraft(ctx context.Context, userID string, adEntity *entities.Ad) error {
	err := utils.ValidateAd(adEntity)
	if err != nil {
		s.logger.ERROR(err)
		return usecaseerr.ErrInvalidParams
	}

	adEntity.AuthorID = userID
	adEntity.Status = entities.StatusPending
	adEntity.IsActive = false
	now := time.Now().UTC()
	adEntity.CreatedAt = now
	adEntity.UpdatedAt = now

	err = s.repo.Create(ctx, adEntity)
	if err != nil {
		s.logger.ERROR("error creating ad:", err)
		return repoerr.ErrInsert
	}

	return nil
}

func (s *service) GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error) {
	ads, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.ERROR("error getting my ads: ", err)
		return nil, repoerr.ErrGettingAdsByUserID
	}
	if len(ads) == 0 {
		s.logger.ERROR("user not found")
		return nil, usecaseerr.ErrUserNotHaveAds
	}
	s.logger.INFO("ads retrieved successfully: ")
	return ads, nil
}

func (s *service) UpdateMyAd(ctx context.Context, userID string, adEntity *entities.Ad) error {
	ad, err := s.repo.GetByID(ctx, adEntity.ID)
	if err != nil {
		s.logger.ERROR("error getting my ad by ID: ", err)
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil {
		s.logger.ERROR("ad is nil")
		return usecaseerr.ErrGettingAdByID
	}
	if ad.AuthorID != userID {
		s.logger.WARN("userID denied: ", userID)
		s.logger.WARN("ad author denied: ", ad.AuthorID)
		s.logger.ERROR("access denied")
		return usecaseerr.ErrAccessDenied
	}

	if err = utils.ValidateAd(adEntity); err != nil {
		s.logger.ERROR(err)
		return usecaseerr.ErrInvalidParams
	}

	ad.Title = adEntity.Title
	ad.Description = adEntity.Description
	ad.CategoryID = adEntity.CategoryID
	ad.UpdatedAt = time.Now().UTC()

	if err = s.repo.Update(ctx, ad); err != nil {
		s.logger.ERROR("error updating my ad: ", err)
		return repoerr.ErrUpdate
	}
	s.logger.INFO("my ad successfully updated")
	return nil
}

func (s *service) DeleteMyAd(ctx context.Context, userID string, adID int) error {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		s.logger.ERROR("error getting my ad by ID: ", err)
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil || ad.AuthorID != userID {
		s.logger.ERROR("user not found")
		return usecaseerr.ErrAccessDenied
	}

	if err = s.repo.Delete(ctx, adID); err != nil {
		s.logger.ERROR("error deleting my ad: ", err)
		return repoerr.ErrDelete
	}
	s.logger.INFO("my ad successfully deleted")
	return nil
}

func (s *service) SubmitForModeration(ctx context.Context, userID string, adID int) error {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		s.logger.ERROR("error submitting ad:", err)
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil || ad.AuthorID != userID {
		s.logger.ERROR("error submitting ad: invalid user")
		return usecaseerr.ErrAccessDenied
	}

	ad.Status = entities.StatusPending
	ad.UpdatedAt = time.Now().UTC()

	err = s.repo.Update(ctx, ad)
	if err != nil {
		s.logger.ERROR("error submitting ad:", err)
		return usecaseerr.ErrApprovingAd
	}
	s.logger.INFO("my ad successfully submitted")
	return nil
}

func (s *service) AddImageToMyAd(ctx context.Context, userID string, file *entities.AdFile) error {
	ad, err := s.repo.GetByID(ctx, file.AdID)
	if err != nil {
		s.logger.ERROR("error getting ads by user ID: ", userID, "\n", err)
		return repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		s.logger.ERROR("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	if !checkIfFileAllowed(file.FileName) {
		s.logger.ERROR("invalid format of the file:", file.FileName)
		return usecaseerr.ErrFileNotAllowed
	}

	dirPath := fmt.Sprintf("storage/uploadings/ad_%d", file.AdID)

	if err := os.MkdirAll(dirPath, fileDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	_, err = s.fileRepo.Create(ctx, file)
	if err != nil {
		s.logger.ERROR("error adding image to ad ", file.AdID, "\n", err)
		return repoerr.ErrFileInsertion
	}

	s.logger.INFO("ad successfully added image to ad ", file.AdID)
	return nil
}

func (s *service) DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error {
	ad, err := s.repo.GetByID(ctx, file.AdID)
	if err != nil {
		s.logger.ERROR("error getting ad by ID: ", file.AdID, "\n", err)
		return repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		s.logger.ERROR("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	url, err := s.fileRepo.Delete(ctx, file)
	if err != nil {
		s.logger.ERROR("error deleting image from ad: ", file.AdID, "\n", err)
		return repoerr.ErrFileDeletion
	}
	s.logger.ERROR("image deleted from file db successfully")

	err = os.Remove(url) // Remove the file from the filesystem
	if err != nil {
		s.logger.ERROR("error removing file from filesystem: ", err)
	}
	s.logger.INFO("ad image successfully deleted")
	return nil
}

func (s *service) GetImagesToMyAd(ctx context.Context, userID string, adID int) ([]entities.AdFile, error) {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		s.logger.ERROR("error getting ad by ID ", adID, "\n", err)
		return nil, repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		s.logger.ERROR("error: user does not own the ad")
		return nil, usecaseerr.ErrAccessDenied
	}
	files, err := s.fileRepo.GetAll(ctx, adID)
	if err != nil {
		s.logger.ERROR("error getting images for ad with ID", adID, "\n", err)
		return nil, repoerr.ErrSelection
	}
	s.logger.INFO("found ", len(files), " images for ad with id: ", adID)
	return files, nil
}

func (s *service) GetMyAdsByFilter(ctx context.Context, userID string,
	filter *entities.AdFilter) ([]entities.Ad, error) {
	filter.UserID = userID
	ads, err := s.repo.Filter(ctx, filter)
	if err != nil {
		s.logger.ERROR("error getting my ads: ", err)
		return nil, repoerr.ErrGettingAdsByUserID
	}
	if len(ads) == 0 {
		s.logger.ERROR("user not found")
		return nil, usecaseerr.ErrUserNotHaveAds
	}

	s.logger.INFO("ads retrieved successfully: ")
	return ads, nil
}

func checkIfFileAllowed(fileName string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".svg"}
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
