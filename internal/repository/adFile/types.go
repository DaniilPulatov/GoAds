package adfile

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdFileRepository interface {
	AddImage(ctx context.Context, file *entities.AdFile) (int, error)
	GetAllAdImages(ctx context.Context, adID int) ([]entities.AdFile, error)
	DeleteImage(ctx context.Context, file *entities.AdFile) (string, error)
}
type adFileRepo struct {
	db *pgxpool.Pool
}

func NewAdFileRepo(db *pgxpool.Pool) AdFileRepository {
	return &adFileRepo{db: db}
}
