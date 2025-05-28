package adfile

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/pkg/db"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAdFileRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)
		file := &entities.AdFile{AdID: 1, FileName: "file.jpg", URL: "http://example.com/file.jpg"}

		mockRow := new(db.MockRow)
		mockRow.On("Scan", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
			*(args[0].(*int)) = 10
		}).Return(nil)
		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 3
		})).Return(mockRow)

		id, err := repo.Create(context.Background(), file)
		assert.NoError(t, err)
		assert.Equal(t, 10, id)
	})

	t.Run("error scan", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)
		file := &entities.AdFile{}

		mockRow := new(db.MockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("scan error"))
		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 3
		})).Return(mockRow)

		id, err := repo.Create(context.Background(), file)
		assert.Error(t, err)
		assert.Equal(t, -1, id)
		assert.Equal(t, repoerr.ErrFileInsertion, err)
	})
}

func TestAdFileRepo_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)

		adID := 1
		mockRows := new(db.MockRows)
		mockPool.On("Query", mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 1 && args[0] == adID
		})).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*(args[0].(*int)) = 1
			*(args[1].(*int)) = adID
			*(args[2].(*string)) = "file.jpg"
			*(args[3].(*string)) = "http://example.com/file.jpg"
			*(args[4].(*time.Time)) = time.Now()
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Once()

		files, err := repo.GetAll(context.Background(), adID)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, "file.jpg", files[0].FileName)
	})

	t.Run("query error", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)

		mockRows := new(db.MockRows)
		mockPool.On("Query", mock.Anything, mock.Anything,
			mock.MatchedBy(func(args []interface{}) bool {
				return len(args) == 1 && args[0] == 1
			})).Return(mockRows, errors.New("query error"))

		files, err := repo.GetAll(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrFileSelection, err)
	})
}

func TestAdFileRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)
		file := &entities.AdFile{ID: 1, AdID: 1}

		mockRow := new(db.MockRow)
		mockRow.On("Scan", mock.AnythingOfType("*string")).Run(func(args mock.Arguments) {
			*(args[0].(*string)) = "http://example.com/file.jpg"
		}).Return(nil)
		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			mock.MatchedBy(func(arg interface{}) bool {
				if args, ok := arg.([]interface{}); ok && len(args) == 1 {
					if id, ok := args[0].(int); ok {
						return id == file.ID
					}
				}
				return false
			})).Return(mockRow)
		mockPool.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{file.ID, file.AdID}).
			Return(pgconn.CommandTag{}, nil)

		url, err := repo.Delete(context.Background(), file)
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com/file.jpg", url)
	})

	t.Run("not found", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)
		file := &entities.AdFile{ID: 1, AdID: 1}

		mockRow := new(db.MockRow)
		mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)
		mockPool.On("QueryRow", mock.Anything, mock.Anything,
			mock.MatchedBy(func(args []interface{}) bool {
				return len(args) == 1 && args[0] == file.ID
			})).Return(mockRow)

		url, err := repo.Delete(context.Background(), file)
		assert.Equal(t, "", url)
		assert.Equal(t, repoerr.ErrFileNotFound, err)
	})
}
