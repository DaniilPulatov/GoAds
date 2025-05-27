package adfile

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/pkg/db"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

/*
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
			mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.MatchedBy(func(arg interface{}) bool {
				if args, ok := arg.([]interface{}); ok && len(args) == 1 {
					if id, ok := args[0].(int); ok {
						return id == file.ID
					}
				}
				return false
			})).Return(mockRow)
			mockPool.On("Exec", mock.Anything, mock.Anything, file.ID, file.AdID).Return(nil, nil)

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
			mockRow.On("Scan", mock.Anything).Return(errors.New("no rows in result set"))
			mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
				return len(args) == 1 && args[0] == file.ID
			})).Return(mockRow)

			url, err := repo.Delete(context.Background(), file)
			assert.Equal(t, "", url)
			assert.Equal(t, repoerr.ErrFileNotFound, err)
		})
	}

	func TestAdFileRepo_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows) // Using db.MockRows as per your provided mock
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		repo := NewAdFileRepo(mockPool)

		adID := 1
		// FIX: Pass adID as a slice of interface{} for variadic argument matching
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Once()
		// FIX: Use mock.AnythingOfType("[]interface {}") for Scan's variadic dest arguments
		// And correctly cast args[N] to ([]interface{}) to access the pointers
		mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			dest := args.Get(0).([]interface{})
			*(dest[0].(*int)) = 1
			*(dest[1].(*int)) = adID
			*(dest[2].(*string)) = "file.jpg"
			*(dest[3].(*string)) = "http://example.com/file.jpg"
			*(dest[4].(*time.Time)) = time.Now()
		}).Return(nil).Once() // Add .Once() for clarity and strictness

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once() // Add .Once()
		mockRows.On("Close").Return().Once()  // Add .Once()

		files, err := repo.GetAll(context.Background(), adID)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, "file.jpg", files[0].FileName)
	})

	t.Run("query error", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)
		repo := NewAdFileRepo(mockPool)

		// FIX: Pass adID as a slice of interface{} for variadic argument matching
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("query error")).Once()

		files, err := repo.GetAll(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, files)
		assert.Equal(t, repoerr.ErrFileSelection, err)
	})
}
*/
