package user

import (
	"ads-service/internal/entities"
	"context"
)

type UserRepository interface{
	// View only user's ads
	GetAllAds(ctx context.Context, userUUid string) ([]entities.Ad, error)
	GetSomeAds(ctx context.Context, limit int, userUUid string) ([]entities.Ad, error)
	GetAdByID(ctx context.Context, id int, userUUid string) (*entities.Ad, error)
	GetAdsByCat(ctx context.Context, cat string, userUUid string) (*entities.Ad, error)
/*
	GetDrafts(ctx context.Context) ([]entities.Ad, error)
	GetSomeDrafts(ctx context.Context, limit int) ([]entities.Ad, error)
	GetDraftByID(ctx context.Context, id int) (entities.Ad, error)
	CheckStatus(ctx context.Context, adID int)(string, error)
*/
	CreateAndSendAd(ctx context.Context, ad *entities.Ad) error // automatically send for moderation
	CreateDraft(ctx context.Context, ad *entities.Ad) error
	SendForModeration(ctx context.Context, adID int) error

	UpdateAd(ctx context.Context, ad *entities.Ad, userUUid string) error
	DeleteAd(ctx context.Context, adID int, userUUid string) error

	AttachFile(ctx context.Context, adID int, fileURL string, userUUid string) error
}
