package entities

import "context"

type AdRepository interface {
	// will be used for admin only
	GetAll(ctx context.Context, userID string) ([]Ad, error)
	GetSome(ctx context.Context, userID string, limit int) ([]Ad, error)
	GetByID(ctx context.Context, userID string, id int) (*Ad, error)

	// Will be used for user and admin
	GetByUserID(ctx context.Context, userUUid string) (*Ad, error)
	GetSomeByUserID(ctx context.Context, userID string, limit int) ([]Ad, error)
	GetByCat(ctx context.Context, userID string, cat string) (*Ad, error)

	//CheckStatus(ctx context.Context, adID int)(string, error)

	Update(ctx context.Context, ad *Ad) error
	Delete(ctx context.Context, adID int) error

	AttachFile(ctx context.Context, adID int, fileURL string, userUUid string) error

	CreateAndSendAd(ctx context.Context, ad *Ad) error // automatically send for moderation
	CreateDraft(ctx context.Context, ad *Ad) error
	SendForModeration(ctx context.Context, adID int) error

	PublishAd(ctx context.Context, id int) error
	RejectAd(ctx context.Context, id int) error

	GetAllRejected(ctx context.Context) ([]Ad, error)
	GetAllPublished(ctx context.Context) ([]Ad, error)
	GetAllDrafts(ctx context.Context) ([]Ad, error)

	GetSomeRejected(ctx context.Context, limit int) ([]Ad, error)
	GetSomePublished(ctx context.Context, limit int) ([]Ad, error)
	GetSomeDrafts(ctx context.Context, limit int) ([]Ad, error)

	GetStatistic(ctx context.Context) (int, error) // show number of ads that are published, rejected, an in pending
}
