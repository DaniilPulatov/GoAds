package auth

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshRepo(db *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) SavaRefreshToken(ctx context.Context, refreshToken entities.RefreshToken) error {
	if _, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens(expires_at, token, user_id)
		values($1, $2, $3);`,
		refreshToken.ExpiredAt, refreshToken.Token, refreshToken.UserID); err != nil {
		log.Println("while inserting into refresh_tokens:", err)
		return repoerr.ErrSavingToken
	}
	log.Println("refresh token repo cerated successfully")
	return nil
}
