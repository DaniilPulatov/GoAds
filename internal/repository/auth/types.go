package auth

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthRepository interface {
	CreateToken(ctx context.Context, rtoken entities.RefreshToken) error
	GetToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	UpdateToken(ctx context.Context, rtoken entities.RefreshToken) error
	DeleteToken(ctx context.Context, token string) error
	CleanUp(ctx context.Context) error
}

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) AuthRepository {
	return &authRepo{db: db}
}
