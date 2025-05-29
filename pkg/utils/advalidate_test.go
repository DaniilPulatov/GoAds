package utils

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/pkgerr/utilserr"
	"errors"
	"testing"
)

func TestValidateAd(t *testing.T) {
	t.Run("valid ad with empty title", func(t *testing.T) {
		ad := &entities.Ad{
			Title:      "",
			CategoryID: 1,
		}
		err := ValidateAd(ad)
		if !errors.Is(err, utilserr.ErrTitleRequired) {
			t.Errorf("expected %v, got %v", utilserr.ErrTitleRequired, err)
		}

	})

	t.Run("valid ad with invalid category ID", func(t *testing.T) {
		ad := &entities.Ad{
			Title:      "Valid Ad",
			CategoryID: 0,
		}
		err := ValidateAd(ad)
		if !errors.Is(err, utilserr.ErrCategoryRequired) {
			t.Errorf("expected %v, got %v", utilserr.ErrCategoryRequired, err)
		}
	})

	t.Run("valid ad with all fields", func(t *testing.T) {
		ad := &entities.Ad{
			Title:      "Valid Ad",
			CategoryID: 1,
		}
		err := ValidateAd(ad)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("valid ad with empty description", func(t *testing.T) {
		ad := &entities.Ad{
			Title:       "Valid Ad",
			CategoryID:  1,
			Description: "",
		}
		err := ValidateAd(ad)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("valid ad with all fields including description", func(t *testing.T) {
		ad := &entities.Ad{
			Title:       "Valid Ad",
			CategoryID:  1,
			Description: "This is a valid ad description.",
		}
		err := ValidateAd(ad)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("valid ad with negative category ID", func(t *testing.T) {
		ad := &entities.Ad{
			Title:      "Valid Ad",
			CategoryID: -1,
		}
		err := ValidateAd(ad)
		if !errors.Is(err, utilserr.ErrCategoryRequired) {
			t.Errorf("expected %v, got %v", utilserr.ErrCategoryRequired, err)
		}
	})
}
