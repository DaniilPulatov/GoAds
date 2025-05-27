package adfile

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdFileRepository interface {
	Create(ctx context.Context, file *entities.AdFile) (int, error)
	GetAll(ctx context.Context, adID int) ([]entities.AdFile, error)
	Delete(ctx context.Context, file *entities.AdFile) (string, error)
}
type adFileRepo struct {
	db *pgxpool.Pool
}

func NewAdFileRepo(db *pgxpool.Pool) AdFileRepository {
	return &adFileRepo{db: db}
}
