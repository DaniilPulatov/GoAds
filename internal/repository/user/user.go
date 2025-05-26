package user

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (string, error) // return id and error
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	GetAllUser(ctx context.Context) ([]entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewAdRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r userRepo) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO users(
				first_name, last_name, phone, 
				password_hash, role
		)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id;`,
		user.FName, user.LName, user.Phone, user.PasswordHash, user.Role).
		Scan(&user.ID)
	if err != nil {
		return "", repoerr.ErrCreationUser
	}

	return user.ID, nil
}

func (r userRepo) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(ctx, `
		SELECT id, first_name, last_name, phone, role
		FROM users
		WHERE id = $1`, userID).
		Scan(&user.ID, &user.FName, &user.LName, &user.Phone, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No user found with ID:", userID)
			return nil, repoerr.ErrUserNotFound
		}
		log.Println("Error selecting user:", err)
		return nil, repoerr.ErrSelection
	}

	return &user, nil
}

func (r userRepo) GetAllUser(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, first_name, last_name, phone, role
		FROM users`)
	if err != nil {
		log.Println("Error getting users:", err)
		return nil, repoerr.ErrSelection
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err = rows.Scan(&user.ID, &user.FName, &user.LName, &user.Phone,
			&user.Role); err != nil {
			log.Println("Error scanning users:", err)
			return nil, repoerr.ErrScan
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, repoerr.ErrScan
	}

	return users, nil
}

func (r userRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, role = $4
		WHERE id = $5;`,
		user.FName, user.LName, user.Phone, user.Role, user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return repoerr.ErrUpdate
	}

	return nil
}

func (r userRepo) DeleteUser(ctx context.Context, userID string) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1;`, userID)
	if err != nil {
		log.Println("Error deleting user:", err)
		return repoerr.ErrDelete
	}
	if result.RowsAffected() == 0 {
		log.Println("No user found with ID:", userID)
		return repoerr.ErrUserNotFound
	}

	return nil
}
