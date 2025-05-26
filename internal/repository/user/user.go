package user

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) error
	GetByPhone(ctx context.Context, phone string) (*entities.User, error)
	IsExists(ctx context.Context, phone string) (bool, error)
	/*
		GetByID(ctx context.Context, userID string) (*entities.User, error)
		Update(ctx context.Context, user entities.User) error
		Delete(ctx context.Context, userID string) error
	*/
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user entities.User) error {
	insertQuery := `INSERT INTO users (first_name, last_name, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, insertQuery, user.ID, user.FName, user.LName, user.Phone, user.Role, user.PasswordHash)
	if err != nil {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*entities.User, error) {
	selectQuery := `SELECT id, first_name, last_name, phone, role, password_hash, created_at, updated_at FROM users WHERE phone = $1`
	row := r.db.QueryRow(ctx, selectQuery, phone)
	var user entities.User
	err := row.Scan(&user.ID, &user.FName, &user.LName, &user.Phone, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) IsExists(ctx context.Context, phone string) (bool, error) {
	selectQuery := `SELECT COUNT(*) FROM users WHERE phone = $1`
	row := r.db.QueryRow(ctx, selectQuery, phone)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error checking user existence:", err)
		return false, repoerr.ErrScan
	}
	return count > 0, nil
}

/*
func (r *userRepo) GetByID(ctx context.Context, userID string) (*entities.User, error) {
	selectQuery := `SELECT id, first_name, last_name, phone, role, password_hash, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, selectQuery, userID)
	var user entities.User
	err := row.Scan(&user.ID, &user.FName, &user.LName, &user.Phone, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user entities.User) error {
	updateQuery := `UPDATE users SET first_name = $1, last_name = $2, phone = $3, role = $4, password_hash = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.Exec(ctx, updateQuery, user.FName, user.LName, user.Phone, user.Role, user.PasswordHash, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, userID string) error {
	deleteQuery := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}
	return nil
}
*/
