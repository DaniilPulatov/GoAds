package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/repository/auth"
	"ads-service/internal/repository/user"
	"context"
)

type UserAuthService interface {
	Register(ctx context.Context, user entities.User) error
	Login(ctx context.Context, phone, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type authService struct {
	userRepo user.UserRepository
	authRepo auth.AuthRepository
}

func NewAuthService(userRepo user.UserRepository, authRepo auth.AuthRepository) UserAuthService {
	return &authService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
