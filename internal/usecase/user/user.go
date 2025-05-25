package user

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	usecaseerr "ads-service/internal/errs/usecaseErr"
	"ads-service/internal/repository/ad"
	"context"
	"log"
	"os"
)

type UserAdvertisementService interface {
	CreateDraft(ctx context.Context, userID string, ad *entities.Ad) error
	GetMyAds(ctx context.Context, userID string) ([]entities.Ad, error)
	UpdateMyAd(ctx context.Context, userID string, ad *entities.Ad) error
	DeleteMyAd(ctx context.Context, userID string, adID int) error
	SubmitForModeration(ctx context.Context, userID string, adID int) error
	AddImageToMyAd(ctx context.Context, userID string, file *entities.AdFile) error
	DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error
}

type service struct {
	repo ad.AdRepository
}

func NewUserService(repo ad.AdRepository) UserAdvertisementService {
	return &service{
		repo: repo,
	}
}

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
	if ad.AuthorID != userID {
		log.Println("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	fileID, err := s.repo.AddImage(ctx, file)
	if err != nil {
		log.Printf("error adding image to ad %d: %v", file.AdID, err)
		return repoerr.ErrFileInsertion
	}
	_ = fileID // TODO: save file into storage/Post{id}/fileID

	return nil
}

func (s *service) DeleteMyAdImage(ctx context.Context, userID string, file *entities.AdFile) error {
	ad, err := s.repo.GetByID(ctx, file.AdID)
	if err != nil {
		log.Printf("error getting ad by ID %d: %v", file.AdID, err)
		return repoerr.ErrSelection
	}
	if ad.AuthorID != userID {
		log.Println("error: user does not own the ad")
		return usecaseerr.ErrAccessDenied
	}
	url, err := s.repo.DeleteImage(ctx, file)
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
