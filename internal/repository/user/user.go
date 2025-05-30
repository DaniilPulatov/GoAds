package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *userRepo) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO users(
				first_name, last_name, phone, 
				password_hash
		)
		VALUES($1, $2, $3, $4)
		RETURNING id;`,
		user.FName, user.LName, user.Phone, user.PasswordHash).
		Scan(&user.ID)
	if err != nil {
		r.logger.ERROR("Failed to create user")
		return "", repoerr.ErrCreationUser
	}
	r.logger.INFO("Create user successfully")

	return user.ID, nil
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*entities.User, error) {
	selectQuery := `
		SELECT id, first_name, last_name, phone, 
		       role, password_hash, created_at, updated_at
		FROM users
		WHERE phone = $1`

	row := r.db.QueryRow(ctx, selectQuery, phone)
	var user entities.User
	err := row.Scan(&user.ID, &user.FName, &user.LName, &user.Phone, &user.Role,
		&user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.logger.ERROR("Error selecting user by phone:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ERROR("User not found")
			return nil, pgx.ErrNoRows
		}
		r.logger.ERROR("Error selecting user by phone: ", err)
		return nil, repoerr.ErrScan
	}
	r.logger.INFO("User successfully retrieved")
	return &user, nil
}

func (r *userRepo) IsExists(ctx context.Context, phone string) (bool, error) {
	selectQuery := `SELECT COUNT(*) FROM users WHERE phone = $1`
	row := r.db.QueryRow(ctx, selectQuery, phone)
	var count int
	err := row.Scan(&count)
	if err != nil {
		r.logger.ERROR("Error checking user existence:", err)
		return false, repoerr.ErrScan
	}
	r.logger.INFO("Number of users found:", count)
	return count > 0, nil
}

func (r *userRepo) GetAllUser(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, first_name, last_name, phone, role
		FROM users`)
	if err != nil {
		r.logger.ERROR("Error getting users:", err)
		return nil, repoerr.ErrSelection
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err = rows.Scan(&user.ID, &user.FName, &user.LName, &user.Phone,
			&user.Role); err != nil {
			r.logger.ERROR("Error scanning users:", err)
			return nil, repoerr.ErrScan
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		r.logger.ERROR("Error iterating rows:", err)
		return nil, repoerr.ErrScan
	}
	r.logger.INFO("Users successfully retrieved")
	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(ctx, `
		SELECT id, first_name, last_name, phone, role
		FROM users
		WHERE id = $1`, userID).
		Scan(&user.ID, &user.FName, &user.LName, &user.Phone, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ERROR("No user found with ID:", userID)
			return nil, repoerr.ErrUserNotFound
		}
		r.logger.ERROR("Error selecting user:", err)
		return nil, repoerr.ErrSelection
	}
	r.logger.INFO("User successfully retrieved by ID: ", userID)

	return &user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, role = $4
		WHERE id = $5;`,
		user.FName, user.LName, user.Phone, user.Role, user.ID)
	if err != nil {
		r.logger.ERROR("Error updating user:", err)
		return repoerr.ErrUpdate
	}
	r.logger.INFO("User successfully updated")
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, userID string) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1;`, userID)
	if err != nil {
		r.logger.ERROR("Error deleting user:", err)
		return repoerr.ErrDelete
	}
	if result.RowsAffected() == 0 {
		r.logger.ERROR("No user found with ID:", userID)
		return repoerr.ErrUserNotFound
	}
	r.logger.INFO("User successfully deleted")
	return nil
}
