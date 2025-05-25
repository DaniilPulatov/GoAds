package admin

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/ad"
)

type AdminAdvertisementService interface {
	GetAllAds(filter *entities.AdFilter) ([]entities.Ad, error)
	GetStatistics() (entities.AdStatistics, error)
	ChangeAdStatus(adID string, status string, adminID string) error
	DeleteAd(adID string, adminID string) error
	DeleteImage(adID string, imageID string, adminID string) error
}

type Service struct {
	adRepo ad.AdRepository
}

func NewAdminService(adRepo ad.AdRepository) AdminAdvertisementService {
	return &Service{
		adRepo: adRepo,
	}
}
func (s *Service) GetAllAds(filter *entities.AdFilter) ([]entities.Ad, error) {
	return nil, nil // TODO: implement
}
func (s *Service) ChangeAdStatus(adID string, status string, adminID string) error {
	return nil // TODO: implement
}
func (s *Service) DeleteAd(adID string, adminID string) error {
	return nil // TODO: implement
}
func (s *Service) GetStatistics() (entities.AdStatistics, error) {
	return entities.AdStatistics{}, nil // TODO: implement
}
func (s *Service) DeleteImage(adID string, imageID string, adminID string) error {
	return nil // TODO: implement
}

// Additional methods for the admin service can be added here as needed.

// This service will handle administrative tasks related to advertisements,
// such as managing ad statuses and retrieving statistics.

// The methods will interact with the ad repository to perform the necessary operations.
