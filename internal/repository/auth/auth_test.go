package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/pkg/db"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthRepo_Create(t *testing.T) {
	t.Run("error deleting old token", func(t *testing.T) {
		mockDB := new(db.MockPool)
		defer mockDB.AssertExpectations(t)

		repo := &authRepo{pool: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("delete error")).Once()

		err := repo.Create(context.Background(), entities.Token{UserID: "user123"})
		assert.Equal(t, repoerr.ErrTokenDeleteFailed, err)
	})

	t.Run("error inserting new token", func(t *testing.T) {
		mockDB := new(db.MockPool)
		defer mockDB.AssertExpectations(t)

		repo := &authRepo{pool: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything, []interface{}{"user123"}).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("insert error")).Once()

		err := repo.Create(context.Background(), entities.Token{UserID: "user123"})
		assert.Equal(t, repoerr.ErrCreatingToken, err)
	})

	t.Run("successfully created", func(t *testing.T) {
		mockDB := new(db.MockPool)
		defer mockDB.AssertExpectations(t)

		repo := &authRepo{pool: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything, []interface{}{"user123"}).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()

		err := repo.Create(context.Background(), entities.Token{UserID: "user123"})
		assert.Nil(t, err)
	})
}

func TestGetToken_Success(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	expectedToken := &Token{
		UserID:    "user-123",
		Token:     "refresh-token-xyz",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("Get", ctx, "user-123").Return(expectedToken, nil)

	actualToken, err := mockRepo.Get(ctx, "user-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, actualToken)
	mockRepo.AssertExpectations(t)
}

func TestGetToken_Error(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	mockRepo.On("Get", ctx, "non-existent-user").Return(nil, errors.New("user not found"))

	actualToken, err := mockRepo.Get(ctx, "non-existent-user")

	assert.Error(t, err)
	assert.Nil(t, actualToken)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateToken_Success(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	rtoken := Token{
		UserID:    "user-123",
		Token:     "updated-token",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	mockRepo.On("Update", ctx, rtoken).Return(nil)

	err := mockRepo.Update(ctx, rtoken)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateToken_Error(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	rtoken := Token{
		UserID:    "user-123",
		Token:     "",
		ExpiresAt: time.Now().Add(48 * time.Hour), // Пустой токен
	}

	mockRepo.On("Update", ctx, rtoken).Return(errors.New("invalid token"))

	err := mockRepo.Update(ctx, rtoken)

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteToken_Success(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	userID := "user-123"

	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := mockRepo.Delete(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteToken_Error(t *testing.T) {
	mockRepo := new(AuthRepositoryMock)
	ctx := context.Background()

	userID := "non-existent-user"

	mockRepo.On("Delete", ctx, userID).Return(errors.New("user not found"))

	err := mockRepo.Delete(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
