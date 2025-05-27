package ad

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdRepository interface {
	Create(ctx context.Context, ad *entities.Ad) error
	GetByID(ctx context.Context, id int) (*entities.Ad, error)
	GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error)
	GetAll(ctx context.Context) ([]entities.Ad, error)
	Update(ctx context.Context, ad *entities.Ad) error
	Delete(ctx context.Context, id int) error
	Approve(ctx context.Context, id int, ad *entities.Ad) error
	Reject(ctx context.Context, id int, ad *entities.Ad) error
	GetStatistics(ctx context.Context) (entities.AdStatistics, error)
}

type adRepo struct {
	db *pgxpool.Pool
}

func NewAdRepo(db *pgxpool.Pool) AdRepository {
	return &adRepo{db: db}
}
