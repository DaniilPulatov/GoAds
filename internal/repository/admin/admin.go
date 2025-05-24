package admin

import (
	"ads-service/internal/entities"
	"context"
)

type AdminRepo interface {
	GetAllAds(ctx context.Context) ([]entities.Ad, error)
	GetSomeAds(ctx context.Context, limit int) ([]entities.Ad, error)
	GetAdByID(ctx context.Context, id int) (*entities.Ad, error)

	PublishAd(ctx context.Context, id int) error
	RejectAd(ctx context.Context, id int) error

	DeleteAd(ctx context.Context, adID int) error

	GetAllRejected(ctx context.Context) ([]entities.Ad, error)
	GetAllPublished(ctx context.Context) ([]entities.Ad, error)
	GetAllDrafts(ctx context.Context) ([]entities.Ad, error)

	GetSomeRejected(ctx context.Context, limit int) ([]entities.Ad, error)
	GetSomePublished(ctx context.Context, limit int) ([]entities.Ad, error)
	GetSomeDrafts(ctx context.Context, limit int) ([]entities.Ad, error)

	GetStatistic(ctx context.Context) (int, error) // show number of ads that are published, rejected, an in pending
	/*
	   GetNoDrafts(ctx context.Context) (int, error)
	   GetNoPublished(ctx context.Context) (int, error)
	   GetNoRejected(ctx context.Context) (int, error)
	*/
}
