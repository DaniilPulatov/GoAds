package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/ad"
)

type UserAdvertisementService interface {
	CreateDraft(userID string, ad *entities.Ad) error
	GetMyAds(userID string) ([]entities.Ad, error)
	UpdateMyAd(userID string, ad *entities.Ad) error
	DeleteMyAd(userID string, adID string) error
	SubmitForModeration(userID string, adID string) error
	AddImageToMyAd(userID string, adID string, imageURL string) error
	DeleteMyAdImage(userID string, adID string, imageID string) error
}

type Service struct {
	adRepo ad.AdRepository
}

func NewUserService(repo ad.AdRepository) UserAdvertisementService {
	return &Service{
		adRepo: repo,
	}
}

func (s *Service) CreateDraft(userID string, adEntity *entities.Ad) error {
	return nil //TODO: implement
}

func (s *Service) GetMyAds(userID string) ([]entities.Ad, error) {
	return nil, nil //TODO: implement
}

func (s *Service) UpdateMyAd(userID string, adEntity *entities.Ad) error {
	return nil //TODO: implement
}
func (s *Service) DeleteMyAd(userID, adID string) error {
	return nil //TODO: implement
}
func (s *Service) SubmitForModeration(userID, adID string) error {
	return nil //TODO: implement
}
func (s *Service) AddImageToMyAd(userID, adID, imageURL string) error {
	return nil //TODO: implement
}

func (s *Service) DeleteMyAdImage(userID, adID, imageID string) error {
	return nil //TODO: implement
}
