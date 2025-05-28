package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/internal/errs/usecaseerr"
	"ads-service/pkg/utils"
	"time"

	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

const fileDirPerm = 0o600 // Permissions for the directory where ad images are stored

func (s *service) CreateDraft(ctx context.Context, userID string, adEntity *entities.Ad) error {
	err := utils.ValidateAd(adEntity)
	if err != nil {
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
		log.Println("error creating ad:", err)
		return repoerr.ErrInsert
	}

	return nil
}

func (s *service) GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error) {
	ads, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, repoerr.ErrGettingAdsByUserID
	}
	if len(ads) == 0 {
		return nil, usecaseerr.ErrUserNotHaveAds
	}

	return ads, nil
}

func (s *service) UpdateMyAd(ctx context.Context, userID string, adEntity *entities.Ad) error {
	log.Println(adEntity.Title)
	log.Println(adEntity.ID)
	log.Println(adEntity.Description)

	ad, err := s.repo.GetByID(ctx, adEntity.ID)
	if err != nil {
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil || ad.AuthorID != userID {
		return usecaseerr.ErrAccessDenied
	}

	if err = utils.ValidateAd(adEntity); err != nil {
		return usecaseerr.ErrInvalidParams
	}

	ad.Title = adEntity.Title
	ad.Description = adEntity.Description
	ad.CategoryID = adEntity.CategoryID
	ad.UpdatedAt = time.Now().UTC()

	if err = s.repo.Update(ctx, ad); err != nil {
		return repoerr.ErrUpdate
	}

	return nil
}

func (s *service) DeleteMyAd(ctx context.Context, userID string, adID int) error {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil || ad.AuthorID != userID {
		return usecaseerr.ErrAccessDenied
	}

	if err = s.repo.Delete(ctx, adID); err != nil {
		return repoerr.ErrDelete
	}

	return nil
}

func (s *service) SubmitForModeration(ctx context.Context, userID string, adID int) error {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		return usecaseerr.ErrGettingAdByID
	}
	if ad == nil || ad.AuthorID != userID {
		return usecaseerr.ErrAccessDenied
	}

	ad.Status = entities.StatusPending
	ad.UpdatedAt = time.Now().UTC()

	err = s.repo.Update(ctx, ad)
	if err != nil {
		return usecaseerr.ErrApprovingAd
	}

	return nil
}

func (s *service) AddImageToMyAd(ctx context.Context, userID string, file *entities.AdFile) error {
	ad, err := s.repo.GetByID(ctx, file.AdID)
	if err != nil {
		log.Printf("error getting ads by user ID %s: %v", userID, err)
		return repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		log.Println("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	if !checkIfFileAllowed(file.FileName) {
		log.Printf("error: file %s is not allowed", file.FileName)
		return usecaseerr.ErrFileNotAllowed
	}

	dirPath := fmt.Sprintf("storage/ad_%d", file.AdID)

	if err := os.MkdirAll(dirPath, fileDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// TODO: delete logs after debugging
	log.Printf("dirPath: %s", dirPath)
	log.Printf("file.FileName: %s", file.FileName)

	file.URL = dirPath + "/" + file.FileName

	log.Printf("file.URL: %s", file.URL)
	_, err = s.fileRepo.Create(ctx, file)
	if err != nil {
		log.Printf("error adding image to ad %d: %v", file.AdID, err)
		return repoerr.ErrFileInsertion
	}
	return nil
}

func (s *service) DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error {
	ad, err := s.repo.GetByID(ctx, file.AdID)
	if err != nil {
		log.Printf("error getting ad by ID %d: %v", file.AdID, err)
		return repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		log.Println("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	url, err := s.fileRepo.Delete(ctx, file)
	if err != nil {
		log.Printf("error deleting image from ad %d: %v", file.AdID, err)
		return repoerr.ErrFileDeletion
	}
	log.Println("image deleted from file db successfully")

	err = os.Remove(url) // Remove the file from the filesystem
	if err != nil {
		log.Printf("error removing file from filesystem: %v", err)
	}

	return nil
}

func (s *service) GetImagesToMyAd(ctx context.Context, userID string, adID int) ([]entities.AdFile, error) {
	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		log.Printf("error getting ad by ID %d: %v", adID, err)
		return nil, repoerr.ErrSelection
	}
	if ad != nil && ad.AuthorID != userID {
		log.Println("error: user does not own the ad")
		return nil, usecaseerr.ErrAccessDenied
	}
	files, err := s.fileRepo.GetAll(ctx, adID)
	if err != nil {
		log.Printf("error getting images for ad %d: %v", adID, err)
		return nil, repoerr.ErrSelection
	}
	log.Printf("found %d images for ad %d", len(files), adID)
	return files, nil
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
