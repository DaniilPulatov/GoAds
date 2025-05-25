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
	Delete(ctx context.Context, id int) error
	Approve(ctx context.Context, id int, ad *entities.Ad) error
	Reject(ctx context.Context, id int, ad *entities.Ad) error
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
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, title, description, category_id, 
		status, is_active, created_at, updated_at, location
		FROM ads
	`)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerr.ErrNoRows
		}
		log.Println("Error getting ads:", err)
		return nil, repoerr.ErrGettingAllAds
	}

	defer rows.Close()

	var ads []entities.Ad
	for rows.Next() {
		var ad entities.Ad
		if err = rows.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt, &ad.Location); err != nil {
			log.Println("Scan error:", err)
			return nil, repoerr.ErrScan
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

func (r adRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	var ad entities.Ad
	err := r.db.QueryRow(ctx, `
		SELECT id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at, location
		FROM ads
		WHERE id = $1`, id).
		Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt, &ad.Location)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ad found with ID:", id)
			return nil, repoerr.ErrAdNotFound
		}
		log.Println("Error selecting ad:", err)
		return nil, repoerr.ErrGettingAdByID
	}

	return &ad, nil
}

func (r adRepo) GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, title, description, category_id, 
		status, is_active, created_at, updated_at, location
		FROM ads
		WHERE author_id = $1`, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ads found for user with ID:", userID)
			return nil, repoerr.ErrUserNotHaveAds
		}
		log.Println("Error getting user ads:", err)
		return nil, repoerr.ErrGettingAdsByUserID
	}

	defer rows.Close()

	var ads []entities.Ad
	for rows.Next() {
		var ad entities.Ad
		if err = rows.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt, &ad.Location); err != nil {
			log.Println("Scan error:", err)
			return nil, repoerr.ErrScan
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

func (r adRepo) Create(ctx context.Context, ad *entities.Ad) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO ads(author_id, title, description, location, category_id)
		VALUES($1, $2, $3, $4, $5);`,
		ad.AuthorID, ad.Title, ad.Description, ad.Location, ad.CategoryID)
	if err != nil {
		log.Println("while inserting into ads:", err)
		return repoerr.ErrInsert
	}

	return nil
}

func (r adRepo) Update(ctx context.Context, ad *entities.Ad) error {
	res, err := r.db.Exec(ctx, `
		UPDATE ads
		SET title = $1, description = $2, location = $3, category_id = $4,
			status = $5, is_active = $6, updated_at = $7
		WHERE id = $8;`, ad.Title, ad.Description, ad.Location, ad.CategoryID,
		ad.Status, ad.IsActive, ad.UpdatedAt, ad.ID)
	if err != nil {
		log.Println("Error updating ad:", err)
		return repoerr.ErrUpdate
	}

	if res.RowsAffected() == 0 {
		log.Println("No ad found with ID:", ad.ID)
		return repoerr.ErrAdNotFound
	}

	return nil
}

func (r adRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM ads
		WHERE id = $1;`, id)
	if err != nil {
		log.Println("Error deleting ad:", err)
		return repoerr.ErrDelete
	}

	return nil
}

func (r adRepo) Approve(ctx context.Context, id int, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET status = $1, is_active = $2, updated_at = $3
		WHERE id = $4;`,
		ad.Status, ad.IsActive, ad.UpdatedAt, id)
	if err != nil {
		log.Println("Error approving ad:", err)
		return repoerr.ErrApproval
	}
	if row.RowsAffected() == 0 {
		log.Println("No ad found with ID:", id)
		return repoerr.ErrApproval
	}

	return nil
}

func (r adRepo) Reject(ctx context.Context, id int, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET status = $1, rejection_reason = $2, is_active = $3, updated_at = $4
		WHERE id = $5;`, ad.Status, ad.RejectionReason, ad.IsActive, ad.UpdatedAt, id)
	if err != nil {
		log.Println("Error rejecting ad:", err)
		return repoerr.ErrRejection
	}
	if row.RowsAffected() == 0 {
		log.Println("No ad found with ID:", id)
		return repoerr.ErrRejection
	}

	return nil
}

func (r adRepo) AddImage(ctx context.Context, file *entities.AdFile) (int, error) {
	var (
		insertQuery = `INSERT INTO ad_files (ad_id, file_name, url) VALUES ($1, $2, $3) RETURNING id`
		fileID      int
	)
	row := r.db.QueryRow(ctx, insertQuery, file.AdID, file.FileName, file.URL)
	err := row.Scan(&fileID)
	if err != nil {
		log.Println("Error scanning fileID:", err)
		return -1, repoerr.ErrFileInsertion
	}
	return fileID, nil
}

func (r adRepo) DeleteImage(ctx context.Context, file *entities.AdFile) (string, error) {
	var (
		selectQuery = `SELECT url FROM ad_files WHERE id = $1 RETURNING url`
		deleteQuery = `DELETE FROM ad_files WHERE id = $1 AND ad_id = $2`
		url         string
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

	if _, err := r.db.Exec(ctx, deleteQuery, file.ID, file.AdID); err != nil {
		log.Println("Error deleting ad file:", err)
		return "", repoerr.ErrFileDeletion
	}
	return url, nil
}
