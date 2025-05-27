package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/pkg/db"
	"context"
)

type AuthRepository interface {
	Create(ctx context.Context, rtoken entities.Token) error
	Get(ctx context.Context, token string) (*entities.Token, error)
	Update(ctx context.Context, rtoken entities.Token) error
	Delete(ctx context.Context, token string) error
	// CleanUp(ctx context.Context) error
}

type authRepo struct {
	pool db.Pool
}

func NewAuthRepo(pool db.Pool) AuthRepository {
	return &authRepo{pool: pool}
}
