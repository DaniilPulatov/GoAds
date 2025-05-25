package ad

import (
	"ads-service/internal/domain/entities"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdRepository interface {
	Create(ctx context.Context, ad *entities.Ad) error
	GetByID(ctx context.Context, id int) (*entities.Ad, error)
	GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error)
	GetAll(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error)
	Update(ctx context.Context, ad *entities.Ad) error
	Delete(ctx context.Context, id int) error
}

type adRepo struct {
	db *pgxpool.Pool
}

func NewAdRepo(db *pgxpool.Pool) AdRepository {
	return &adRepo{db: db}
}

func (r adRepo) GetAll(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error) {
	return nil, errors.New("not implemented")
}

func (r adRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	return nil, nil
}

func (r adRepo) GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error) {
	return nil, nil
}
func (r adRepo) Create(ctx context.Context, ad *entities.Ad) error {
	return nil
}
func (r adRepo) Update(ctx context.Context, ad *entities.Ad) error {
	return nil
}
func (r adRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (r adRepo) ChangeStatus(ctx context.Context, id int, status, adminID string) error {
	return nil
}
