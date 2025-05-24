package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

/*
	func (r UserRepo) GetSomeAds(ctx context.Context, limit int) ([]*entities.Ad, error)
	func (r UserRepo) GetAdByID(ctx context.Context, id int) (*entities.Ad, error)
	func (r UserRepo) GetAdByUserID(ctx context.Context, userUUid string) (*entities.Ad, error)
	func (r UserRepo) GetAdsByCat(ctx context.Context, cat string) (*entities.Ad, error)
	func (r UserRepo) UpdateAd(ctx context.Context, ad *entities.Ad) error
	func (r UserRepo) DeleteAd(ctx context.Context, adID int) error

	//CheckStatus(ctx context.Context, adID int)(string, error)
*/
