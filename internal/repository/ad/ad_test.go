package ad

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"ads-service/pkg/db"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAdRepo_Create(t *testing.T) {
	t.Run("error at create ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{pool: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrInsert)

		err := pool.Create(context.Background(), &entities.Ad{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrInsert, err)
	})

	t.Run("success at create ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{pool: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)

		err := pool.Create(context.Background(), &entities.Ad{})
		assert.Nil(t, err)
	})
}

func TestAdRepo_GetByID(t *testing.T) {
	t.Run("error at get ad by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{pool: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan",
			mock.Anything, // id
			mock.Anything, // title
			mock.Anything, // description
			mock.Anything, // image_url
			mock.Anything, // user_id
			mock.Anything, // status
			mock.Anything, // is_published
			mock.Anything, // created_at
			mock.Anything, // updated_at
			mock.Anything, // rejection_reason
		).Return(pgx.ErrNoRows)

		ad, err := pool.GetByID(context.Background(), 1)
		assert.Nil(t, ad)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}
