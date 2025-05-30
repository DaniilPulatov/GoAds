package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/pkg/db"
	customLogger "ads-service/pkg/logger"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (string, error) // return id and error
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	GetAllUser(ctx context.Context) ([]entities.User, error)
	GetByPhone(ctx context.Context, phone string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
	IsExists(ctx context.Context, phone string) (bool, error)
}

type userRepo struct {
	db     db.Pool
	logger customLogger.Logger
}

func NewUserRepo(pool db.Pool, logTool customLogger.Logger) UserRepository {
	return &userRepo{db: pool, logger: logTool}
}
