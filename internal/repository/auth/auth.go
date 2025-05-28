package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *authRepo) Create(ctx context.Context, rtoken entities.Token) error {
	if err := r.Delete(ctx, rtoken.UserID); err != nil {
		r.logger.ERROR("Error deleting existing token for user:", rtoken.UserID, "Error:", err)
		return repoerr.ErrTokenDeleteFailed
	}

	insertQuery := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, insertQuery, rtoken.UserID, rtoken.Token, rtoken.ExpiresAt)
	if err != nil {
		r.logger.ERROR("Error creating token:", err)
		return repoerr.ErrCreatingToken
	}
	r.logger.INFO("Add token successfully")
	return nil
}
func (r *authRepo) Get(ctx context.Context, userID string) (*entities.Token, error) {
	selectQuery := `SELECT token, expires_at FROM refresh_tokens WHERE user_id = $1`
	row := r.pool.QueryRow(ctx, selectQuery, userID)
	var token entities.Token
	err := row.Scan(&token.Token, &token.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ERROR("Token not found for user: ", userID)
			return nil, repoerr.ErrTokenNotFound
		}
		r.logger.ERROR("Error selecting token:", err)
		return nil, repoerr.ErrTokenSelectFailed
	}
	r.logger.INFO("Get token successfully, user ID: ", userID)
	return &token, nil
}
func (r *authRepo) Update(ctx context.Context, rtoken entities.Token) error {
	updateQuery := `UPDATE refresh_tokens SET token = $1, expires_at = $2 WHERE user_id = $3`
	_, err := r.pool.Exec(ctx, updateQuery, rtoken.Token, rtoken.ExpiresAt, rtoken.UserID)
	if err != nil {
		r.logger.ERROR("Error updating token:", err)
		return repoerr.ErrTokenUpdateFailed
	}
	r.logger.INFO("Update token successfully")
	return nil
}
func (r *authRepo) Delete(ctx context.Context, userID string) error {
	deleteQuery := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, deleteQuery, userID)
	if err != nil {
		r.logger.ERROR("Error deleting token for user:", userID, "Error:", err)
		return repoerr.ErrTokenDeleteFailed
	}
	r.logger.INFO("Delete token successfully")
	return nil
}

/*
// TODO: Implement a cleanup function to remove expired tokens with pg_cron or similar
func (r *authRepo) CleanUp(ctx context.Context) error {
	deleteQuery := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.pool.Exec(ctx, deleteQuery)
	if err != nil {
		log.Println("Error cleaning up expired tokens:", err)
		return repoerr.ErrTokenDeleteFailed
	}
	return nil
}
*/
