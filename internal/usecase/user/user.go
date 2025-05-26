package user

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	usecaseerr "ads-service/internal/errs/usecaseErr"

	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

const fileDirPerm = 0o600 // Permissions for the directory where ad images are stored

func (s *service) CreateDraft(ctx context.Context, userID string, adEntity *entities.Ad) error {
	return nil //TODO: implement
}

func (s *service) GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error) {
	return nil, nil //TODO: implement
}

func (s *service) UpdateMyAd(ctx context.Context, userID string, adEntity *entities.Ad) error {
	return nil //TODO: implement
}
func (s *service) DeleteMyAd(ctx context.Context, userID string, adID int) error {
	return nil //TODO: implement
}
func (s *service) SubmitForModeration(ctx context.Context, userID string, adID int) error {
	return nil //TODO: implement
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
	_, err = s.fileRepo.AddImage(ctx, file)
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
	url, err := s.fileRepo.DeleteImage(ctx, file)
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
	files, err := s.fileRepo.GetAllAdImages(ctx, adID)
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
