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
	GetByID(ctx context.Context, id int) (*entities.Ad, error)
	GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error)
	GetAll(ctx context.Context, filter *entities.AdFilter) ([]entities.Ad, error)
	Update(ctx context.Context, ad *entities.Ad) error
	Delete(ctx context.Context, id int, userID string) error
	ChangeStatus(ctx context.Context, id int, status string, adminID string) error
	AddImage(ctx context.Context, file *entities.AdFile) (int, error)
	DeleteImage(ctx context.Context, file *entities.AdFile) (string, error)
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
func (r adRepo) Delete(ctx context.Context, id int, userID string) error {
	return nil
}

func (r adRepo) ChangeStatus(ctx context.Context, id int, status, adminID string) error {
	return nil
}
func (r adRepo) AddImage(ctx context.Context, file *entities.AdFile) (int, error) {
	var (
		insertQuery = `INSERT INTO ad_files (ad_id, file_name, url) VALUES ($1, $2, $3) RETURNING id`
		fileID int
	)
	row := r.db.QueryRow(ctx, insertQuery, file.AdID, file.FileName, file.URL)
	err := row.Scan(&fileID)
	if  err != nil {
		log.Println("Error scanning fileID:", err)
		return -1, repoerr.ErrFileInsertion
	}
	return fileID, nil
}

func (r adRepo) DeleteImage(ctx context.Context, file *entities.AdFile) (string, error) {
	var (
		selectQuery = `SELECT url FROM ad_files WHERE id = $1 RETURNING url`
		delteQuery = `DELETE FROM ad_files WHERE id = $1 AND ad_id = $2`
		url string
	)
	
	row := r.db.QueryRow(ctx, selectQuery, file.ID)
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ad file found with ID:", file.ID)
			return "", repoerr.ErrFileNotFound
		}
		log.Println("Error selecting ad file:", err)
		return "", repoerr.ErrSelection
	}
	
	if _, err := r.db.Exec(ctx, delteQuery, file.ID, file.AdID); err != nil {
		log.Println("Error deleting ad file:", err)
		return "", repoerr.ErrFileDeletion
	}
	return url, nil
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
