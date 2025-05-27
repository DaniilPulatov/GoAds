package auth

import (
	"ads-service/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Create(ctx context.Context, rtoken entities.Token) error
	Get(ctx context.Context, token string) (*entities.Token, error)
	Update(ctx context.Context, rtoken entities.Token) error
	Delete(ctx context.Context, token string) error
	// CleanUp(ctx context.Context) error
}

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) AuthRepository {
	return &authRepo{db: db}
}
