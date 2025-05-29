package ad

import (
	"ads-service/internal/domain/entities"
	"ads-service/pkg/db"
	customLogger "ads-service/pkg/logger"
	"context"
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
	Filter(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error)
}

type adRepo struct {
	db     db.Pool
	logger customLogger.Logger
}

func NewAdRepo(pool db.Pool, logger customLogger.Logger) AdRepository {

	return &adRepo{db: pool, logger: logger}
}
