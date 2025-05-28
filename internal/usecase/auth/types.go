package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/auth"
	"ads-service/internal/repository/user"
	customLogger "ads-service/pkg/logger"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, phone, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

type userAuthService struct {
	userRepo user.UserRepository
	authRepo auth.AuthRepository
	logger   customLogger.Logger
}

func NewAuthService(userRepo user.UserRepository, authRepo auth.AuthRepository, logger customLogger.Logger) AuthService {
	return &userAuthService{
		logger:   logger,
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
