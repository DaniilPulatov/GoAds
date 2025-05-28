package user

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/pkg/db"
	customLogger "ads-service/pkg/logger"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepo_CreateUser(t *testing.T) {
	t.Run("error creating user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))

		_, err := pool.CreateUser(context.Background(), &entities.User{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrCreationUser, err)
	})

	t.Run("successful creating user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		repo := &userRepo{db: mockPool}
		expectedID := "test-id"

		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.AnythingOfType("*string")).Run(func(args mock.Arguments) {
			ptr := args.Get(0).(*string)
			*ptr = expectedID
		}).Return(nil)

		user := &entities.User{}
		id, err := repo.CreateUser(context.Background(), user)
		assert.Nil(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, expectedID, id)
	})
}

func TestUserRepo_GetByPhone(t *testing.T) {
	t.Run("error getting user by phone", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("db error"))

		_, err := pool.GetByPhone(context.Background(), "1234567890")
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrScan, err)
	})

	t.Run("user not found by phone", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgx.ErrNoRows)

		user, err := pool.GetByPhone(context.Background(), "1234567890")
		assert.Nil(t, user)
		assert.Equal(t, pgx.ErrNoRows, err)
	})

	t.Run("successful getting user by phone", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		expectedUser := &entities.User{
			ID:           "test-id",
			FName:        "John",
			LName:        "Doe",
			Phone:        "1234567890",
			Role:         entities.RoleUser,
			PasswordHash: "hashed-password",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			[]interface{}{"1234567890"}).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*(args[0].(*string)) = expectedUser.ID
				*(args[1].(*string)) = expectedUser.FName
				*(args[2].(*string)) = expectedUser.LName
				*(args[3].(*string)) = expectedUser.Phone
				*(args[4].(*entities.Role)) = expectedUser.Role
				*(args[5].(*string)) = expectedUser.PasswordHash
				*(args[6].(*time.Time)) = expectedUser.CreatedAt
				*(args[7].(*time.Time)) = expectedUser.UpdatedAt
			}).Return(nil)

		user, err := pool.GetByPhone(context.Background(), expectedUser.Phone)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.FName, user.FName)
		assert.Equal(t, expectedUser.LName, user.LName)
		assert.Equal(t, expectedUser.Phone, user.Phone)
		assert.Equal(t, expectedUser.Role, user.Role)
		assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	})
}

func TestUserRepo_IsExists(t *testing.T) {
	t.Run("error checking user existence", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))

		exists, err := pool.IsExists(context.Background(), "1234567890")
		assert.False(t, exists)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrScan, err)
	})

	t.Run("user does not exist", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			ptr := args.Get(0).(*int)
			*ptr = 0
		}).Return(nil)

		exists, err := pool.IsExists(context.Background(), "1234567890")
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("user exists", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			ptr := args.Get(0).(*int)
			*ptr = 1
		}).Return(nil)

		exists, err := pool.IsExists(context.Background(), "1234567890")
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}

func TestUserRepo_GetAllUser(t *testing.T) {
	t.Run("error getting all users", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(new(db.MockRows), errors.New("db error"))

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, users)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("successful getting all users", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		expectedUsers := []entities.User{
			{
				ID:    "1",
				FName: "Alice",
				LName: "Smith",
				Phone: "1234567890",
				Role:  entities.RoleUser,
			},
			{
				ID:    "2",
				FName: "Bob",
				LName: "Johnson",
				Phone: "0987654321",
				Role:  entities.RoleAdmin,
			},
		}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		// Первый пользователь
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*(args[0].(*string)) = expectedUsers[0].ID
				*(args[1].(*string)) = expectedUsers[0].FName
				*(args[2].(*string)) = expectedUsers[0].LName
				*(args[3].(*string)) = expectedUsers[0].Phone
				*(args[4].(*entities.Role)) = expectedUsers[0].Role
			}).Return(nil).Once()

		// Второй пользователь
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*(args[0].(*string)) = expectedUsers[1].ID
				*(args[1].(*string)) = expectedUsers[1].FName
				*(args[2].(*string)) = expectedUsers[1].LName
				*(args[3].(*string)) = expectedUsers[1].Phone
				*(args[4].(*entities.Role)) = expectedUsers[1].Role
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("no users found", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, err)
		assert.Empty(t, users)
	})

	t.Run("scan error", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("scan error")).Once()
		mockRows.On("Close").Return()

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, users)
		assert.Equal(t, repoerr.ErrScan, err)
	})

	t.Run("rows error", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows error")).Once()
		mockRows.On("Close").Return()

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, users)
		assert.Equal(t, repoerr.ErrScan, err)
	})
}

func TestUserRepo_GetUserByID(t *testing.T) {
	t.Run("error getting user by ID", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(errors.New("db error"))

		user, err := pool.GetUserByID(context.Background(), "test-id")
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrSelection, err)
	})

	t.Run("user not found by ID", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(pgx.ErrNoRows)

		user, err := pool.GetUserByID(context.Background(), "test-id")
		assert.Nil(t, user)
		assert.Equal(t, repoerr.ErrUserNotFound, err)
	})

	t.Run("successful getting user by ID", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		expectedUser := &entities.User{
			ID:        "test-id",
			FName:     "John",
			LName:     "Doe",
			Phone:     "1234567890",
			Role:      entities.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			[]interface{}{"test-id"}).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*(args[0].(*string)) = expectedUser.ID
				*(args[1].(*string)) = expectedUser.FName
				*(args[2].(*string)) = expectedUser.LName
				*(args[3].(*string)) = expectedUser.Phone
				*(args[4].(*entities.Role)) = expectedUser.Role
			}).Return(nil)

		user, err := pool.GetUserByID(context.Background(), expectedUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.FName, user.FName)
		assert.Equal(t, expectedUser.LName, user.LName)
		assert.Equal(t, expectedUser.Phone, user.Phone)
		assert.Equal(t, expectedUser.Role, user.Role)
	})
}

func TestUserRepo_UpdateUser(t *testing.T) {
	t.Run("error updating user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("db error"))

		err := pool.UpdateUser(context.Background(), &entities.User{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrUpdate, err)
	})

	t.Run("successful_updating_user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil)

		err := pool.UpdateUser(context.Background(), &entities.User{})
		assert.Nil(t, err)
	})
}

func TestUserRepo_DeleteUser(t *testing.T) {
	t.Run("error deleting user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("db error"))

		err := pool.DeleteUser(context.Background(), "test-id")
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrDelete, err)
	})

	t.Run("successful deleting user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything,
			mock.Anything).Return(pgconn.NewCommandTag("DELETE 1"), nil)

		err := pool.DeleteUser(context.Background(), "test-id")
		assert.Nil(t, err)
	})

	t.Run("delete user not found", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("DELETE 0"), nil)

		err := pool.DeleteUser(context.Background(), "not-found-id")
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrUserNotFound, err)
	})
}

func TestUserRepo_GetUserByID_Args(t *testing.T) {
	t.Run("QueryRow called с правильными аргументами", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		id := "id"
		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			mock.MatchedBy(func(args []interface{}) bool {
				return len(args) == 1 && args[0] == "id"
			}),
		).Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgx.ErrNoRows)

		user, err := pool.GetUserByID(context.Background(), id)
		assert.Nil(t, user)
		assert.Equal(t, repoerr.ErrUserNotFound, err)
	})
}

func TestUserRepo_UpdateUser_Args(t *testing.T) {
	t.Run("Exec_called_with_правильными_аргументами", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewUserRepo(mockPool, customLogger.Logger{})

		user := &entities.User{
			ID:    "id",
			FName: "Имя",
			LName: "Фамилия",
			Phone: "123",
			Role:  "user",
		}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil)

		err := repo.UpdateUser(context.Background(), user)
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetAllUser_ScanAndErr(t *testing.T) {
	t.Run("rows.Err возвращает ошибку после Next", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows error")).Once()
		mockRows.On("Close").Return()

		users, err := pool.GetAllUser(context.Background())
		assert.Nil(t, users)
		assert.Equal(t, repoerr.ErrScan, err)
	})
}

func TestUserRepo_GetByPhone_Args(t *testing.T) {
	t.Run("QueryRow called with right argument", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		phone := "9876543210"
		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			[]interface{}{"9876543210"}).Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgx.ErrNoRows)

		user, err := pool.GetByPhone(context.Background(), phone)
		assert.Nil(t, user)
		assert.Equal(t, pgx.ErrNoRows, err)
	})
}

func TestUserRepo_IsExists_Args(t *testing.T) {
	t.Run("QueryRow called с правильными аргументами", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &userRepo{db: mockPool}
		phone := "555"
		mockPool.On("QueryRow", mock.Anything, mock.Anything, []interface{}{phone}).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			ptr := args.Get(0).(*int)
			*ptr = 1
		}).Return(nil)

		exists, err := pool.IsExists(context.Background(), phone)
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
