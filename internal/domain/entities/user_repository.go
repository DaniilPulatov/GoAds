package entities

import (
	"context"
)

// BasicUser - interface that contains common methods for user and admin.
type BasicUser interface {
	GetAllAds(ctx context.Context) ([]Ad, error)
	GetSomeAds(ctx context.Context, limit int) ([]Ad, error)
	GetAdByID(ctx context.Context, id int) (*Ad, error)
	GetAdByUserID(ctx context.Context, userUUid string) (*Ad, error)
	GetAdsByCat(ctx context.Context, cat string) (*Ad, error)

	//CheckStatus(ctx context.Context, adID int)(string, error)

	UpdateAd(ctx context.Context, ad *Ad) error
	DeleteAd(ctx context.Context, adID int) error
}

type UserRepository interface {
	BasicUser
	AttachFile(ctx context.Context, adID int, fileURL string, userUUid string) error
	CreateAndSendAd(ctx context.Context, ad *Ad) error // automatically send for moderation
	CreateDraft(ctx context.Context, ad *Ad) error
	SendForModeration(ctx context.Context, adID int) error

	/*
		GetDrafts(ctx context.Context) ([]Ad, error)
		GetSomeDrafts(ctx context.Context, limit int) ([]Ad, error)
		GetDraftByID(ctx context.Context, id int) (Ad, error)
		CheckStatus(ctx context.Context, adID int)(string, error)
	


		UpdateAd(ctx context.Context, ad *Ad, userUUid string) error
		DeleteAd(ctx context.Context, adID int, userUUid string) error
	*/
	
	
}

type AdminRepository interface {
	BasicUser
	PublishAd(ctx context.Context, id int) error
	RejectAd(ctx context.Context, id int) error

	GetAllRejected(ctx context.Context) ([]Ad, error)
	GetAllPublished(ctx context.Context) ([]Ad, error)
	GetAllDrafts(ctx context.Context) ([]Ad, error)

	GetSomeRejected(ctx context.Context, limit int) ([]Ad, error)
	GetSomePublished(ctx context.Context, limit int) ([]Ad, error)
	GetSomeDrafts(ctx context.Context, limit int) ([]Ad, error)

	GetStatistic(ctx context.Context) (int, error) // show number of ads that are published, rejected, an in pending
	/*
	   GetNoDrafts(ctx context.Context) (int, error)
	   GetNoPublished(ctx context.Context) (int, error)
	   GetNoRejected(ctx context.Context) (int, error)
	*/
}
