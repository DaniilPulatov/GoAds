package utils

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/pkgerr/utilserr"
)

func ValidateAd(a *entities.Ad) error {
	if a.Title == "" {
		return utilserr.ErrTitleRequired
	}

	if a.CategoryID <= 0 {
		return utilserr.ErrCategoryRequired
	}

	return nil
}
