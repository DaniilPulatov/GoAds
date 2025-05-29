//nolint:all // testpackage
package ad

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
)

func TestAdRepo_Create(t *testing.T) {
	t.Run("error at create ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrInsert)

		err := pool.Create(context.Background(), &entities.Ad{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrInsert, err)
	})

	t.Run("success at create ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)

		err := pool.Create(context.Background(), &entities.Ad{})
		assert.Nil(t, err)
	})

	t.Run("not found at create ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrInsert)

		err := pool.Create(context.Background(), &entities.Ad{ID: 1})
		assert.Equal(t, repoerr.ErrInsert, err)
	})

	t.Run("error at create ad with invalid data", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("invalid data"))

		err := pool.Create(context.Background(), &entities.Ad{Title: ""})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrInsert, err)
	})
}

func TestAdRepo_GetByID(t *testing.T) {
	t.Run("error at get ad by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

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

	t.Run("success at get ad by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan",
			mock.Anything, // id
			mock.Anything, // author_id
			mock.Anything, // title
			mock.Anything, // description
			mock.Anything, // category_id
			mock.Anything, // status
			mock.Anything, // is_active
			mock.Anything, // created_at
			mock.Anything, // updated_at
			mock.Anything, // location
		).Return(nil)

		ad, err := pool.GetByID(context.Background(), 1)
		assert.NotNil(t, ad)
		assert.Nil(t, err)
	})

	t.Run("error at get ad by id with invalid id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgx.ErrNoRows)

		ad, err := pool.GetByID(context.Background(), -1)
		assert.Nil(t, ad)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}

func TestAdRepo_GetByUserID(t *testing.T) {
	t.Run("error at get by user id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow, errors.New("db error"))
		ads, err := pool.GetByUserID(context.Background(), "user1")

		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAdsByUserID, err)
	})

	t.Run("success at get by user id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		// Первый вызов Next() — true, второй — false (одна строка)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		ads, err := pool.GetByUserID(context.Background(), "user1")

		assert.NotNil(t, ads)
		assert.Nil(t, err)
	})
}

func TestAdRepo_GetAll(t *testing.T) {

	t.Run("error at get all", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(new(db.MockRows), errors.New("db error"))

		ads, err := pool.GetAll(context.Background())
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAllAds, err)
	})

	t.Run("success at get all", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		ads, err := pool.GetAll(context.Background())
		assert.NotNil(t, ads)
		assert.Nil(t, err)
	})
}

func TestAdRepo_Update(t *testing.T) {
	t.Run("error at update ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrUpdate)

		err := pool.Update(context.Background(), &entities.Ad{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrUpdate, err)
	})

	t.Run("success at update ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil)

		err := pool.Update(context.Background(), &entities.Ad{})
		assert.Nil(t, err)
	})

	t.Run("not found at update ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Update(context.Background(), &entities.Ad{ID: 1})
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}

func TestAdRepo_Delete(t *testing.T) {
	t.Run("error at delete ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrDelete)

		err := pool.Delete(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrDelete, err)
	})

	t.Run("success at delete ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("DELETE 1")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Delete(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("not found at delete ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("DELETE 0")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Delete(context.Background(), 1)
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}

func TestAdRepo_Approve(t *testing.T) {
	t.Run("error at approve ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrApproval)

		err := pool.Approve(context.Background(), 1, &entities.Ad{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrApproval, err)
	})

	t.Run("success at approve ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("UPDATE 1")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Approve(context.Background(), 1, &entities.Ad{})
		assert.Nil(t, err)
	})

	t.Run("not found at approve ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Approve(context.Background(), 1, &entities.Ad{ID: 1})
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}

func TestAdRepo_Reject(t *testing.T) {
	t.Run("error at reject ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, repoerr.ErrRejection)

		err := pool.Reject(context.Background(), 1, &entities.Ad{})
		assert.NotNil(t, err)
		assert.Equal(t, repoerr.ErrRejection, err)
	})

	t.Run("success at reject ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("UPDATE 1")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Reject(context.Background(), 1, &entities.Ad{})
		assert.Nil(t, err)
	})

	t.Run("not found at reject ad", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(tag, nil)

		err := pool.Reject(context.Background(), 1, &entities.Ad{ID: 1})
		assert.Equal(t, repoerr.ErrAdNotFound, err)
	})
}

func TestAdRepo_GetStatistics(t *testing.T) {
	t.Run("error at get statistics", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("scan error"))

		stats, err := pool.GetStatistics(context.Background())
		assert.Equal(t, entities.AdStatistics{}, stats)
		assert.Equal(t, repoerr.ErrGettingStatistics, err)
	})

	t.Run("success at get statistics", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		stats, err := pool.GetStatistics(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, stats)
	})
}

func TestAdRepo_Filter(t *testing.T) {
	t.Run("error at filter ads", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(new(db.MockRows), errors.New("db error"))

		ads, err := pool.Filter(context.Background(), &entities.AdFilter{})
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrGettingAllAds, err)
	})

	t.Run("success at filter ads", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		ads, err := pool.Filter(context.Background(), &entities.AdFilter{})
		assert.NotNil(t, ads)
		assert.Nil(t, err)
	})

	t.Run("scan error at filter ads", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).
			Return(errors.New("scan error")).Once()
		mockRows.On("Close").Return()

		ads, err := pool.Filter(context.Background(), &entities.AdFilter{})
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrScan, err)
	})

	t.Run("rows error at filter ads", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRows := new(db.MockRows)
		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		pool := &adRepo{db: mockPool}
		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRows, nil)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows error"))
		mockRows.On("Close").Return()

		ads, err := pool.Filter(context.Background(), &entities.AdFilter{})
		assert.Nil(t, ads)
		assert.Equal(t, repoerr.ErrScan, err)
	})
}
