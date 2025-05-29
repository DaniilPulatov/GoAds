package adfile

import (
	"ads-service/internal/domain/entities"
	"ads-service/pkg/db"
	customLogger "ads-service/pkg/logger"
	"context"
)

type AdFileRepository interface {
	Create(ctx context.Context, file *entities.AdFile) (int, error)
	GetAll(ctx context.Context, adID int) ([]entities.AdFile, error)
	Delete(ctx context.Context, file *entities.AdFile) (string, error)
}
type adFileRepo struct {
	pool   db.Pool
	logger customLogger.Logger
}

func NewAdFileRepo(db db.Pool, logTool customLogger.Logger) AdFileRepository {
	return &adFileRepo{pool: db, logger: logTool}
}
