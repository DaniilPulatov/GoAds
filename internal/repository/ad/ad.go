package ad

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdRepository interface {
	Create(ctx context.Context, ad *entities.Ad) error
	GetByID(ctx context.Context, id string) (*entities.Ad, error)
	GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error)
	GetAll(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error)
	Update(ctx context.Context, ad *entities.Ad) error
	Delete(ctx context.Context, id string, userID string) error
	ChangeStatus(ctx context.Context, id string, status string, adminID string) error
	AddImage(ctx context.Context, file *entities.AdFile) error
	DeleteImage(ctx context.Context, file *entities.AdFile) error
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

func (r adRepo) GetByID(ctx context.Context, id string) (*entities.Ad, error) {
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
func (r adRepo) Delete(ctx context.Context, id, userID string) error {
	return nil
}

func (r adRepo) ChangeStatus(ctx context.Context, id, status, adminID string) error {
	return nil
}
func (r adRepo) AddImage(ctx context.Context, file *entities.AdFile) error {
	return nil
}

func (r adRepo) DeleteImage(ctx context.Context, file *entities.AdFile) error {
	return nil
}

// queryRows - helper function that processes the rows returned from the database query and
// unmarshal the JSON data for files.
func (r adRepo) queryRows(rows pgx.Rows) ([]entities.Ad, error) {
	defer rows.Close()
	var ads []entities.Ad
	for rows.Next() {
		var ad entities.Ad
		var filesJSON []byte
		err := rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Description,
			&ad.AuthorID,
			&ad.CategoryID,
			&ad.Status,
			&ad.IsActive,
			&ad.CreatedAt,
			&ad.UpdatedAt,
			&filesJSON,
		)
		if err != nil {
			log.Println("Scan error:", err)
			return nil, repoerr.ErrSelection
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
		return nil, repoerr.ErrSelection
	}
	return ads, nil
}

func (r adRepo) queryRow(row pgx.Row) (*entities.Ad, error) {
	var ad entities.Ad
	var filesJSON []byte
	err := row.Scan(
		&ad.ID,
		&ad.Title,
		&ad.Description,
		&ad.AuthorID,
		&ad.CategoryID,
		&ad.Status,
		&ad.IsActive,
		&ad.CreatedAt,
		&ad.UpdatedAt,
		&filesJSON,
	)
	if err != nil {
		log.Println("Scan error:", err)
		return nil, repoerr.ErrSelection
	}

	return &ad, nil
}
