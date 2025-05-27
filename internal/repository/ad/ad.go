package ad

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx"
)

func (r adRepo) GetAll(ctx context.Context) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at, location
		FROM ads
	`)
	if err != nil {
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

	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, repoerr.ErrScan
	}

	return ads, nil
}

func (r adRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	var ad entities.Ad
	err := r.db.QueryRow(ctx, `
		SELECT 
		    id, author_id, title, description, category_id, 
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
		SELECT 
		    id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at, location
		FROM ads
		WHERE author_id = $1`, userID)
	if err != nil {
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

	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, repoerr.ErrScan
	}

	return ads, nil
}

func (r adRepo) Create(ctx context.Context, ad *entities.Ad) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO ads(
            author_id, title, description, location, category_id, 
            status, is_active, created_at, updated_at
        ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		ad.AuthorID, ad.Title, ad.Description, ad.Location, ad.CategoryID,
		ad.Status, ad.IsActive, ad.CreatedAt, ad.UpdatedAt)
	if err != nil {
		log.Println("while inserting into ads:", err)
		return repoerr.ErrInsert
	}

	return nil
}

func (r adRepo) Update(ctx context.Context, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET title = $1, description = $2, location = $3, category_id = $4,
			status = $5, is_active = $6, updated_at = $7
		WHERE id = $8;`, ad.Title, ad.Description, ad.Location, ad.CategoryID,
		ad.Status, ad.IsActive, ad.UpdatedAt, ad.ID)
	if err != nil {
		log.Println("Error updating ad:", err)
		return repoerr.ErrUpdate
	}

	if row.RowsAffected() == 0 {
		log.Println("No ad found with ID:", ad.ID)
		return repoerr.ErrAdNotFound
	}

	return nil
}

func (r adRepo) Delete(ctx context.Context, id int) error {
	row, err := r.db.Exec(ctx, `
		DELETE FROM ads
		WHERE id = $1;`, id)
	if err != nil {
		log.Println("Error deleting ad:", err)
		return repoerr.ErrDelete
	}
	if row.RowsAffected() == 0 {
		log.Println("No ad found with ID:", id)
		return repoerr.ErrAdNotFound
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

func (r adRepo) GetStatistics(ctx context.Context) (entities.AdStatistics, error) {
	var statistics entities.AdStatistics
	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'approved') AS published,
			COUNT(*) FILTER (WHERE status = 'pending') AS pending,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected
		FROM ads;
	`
	err := r.db.QueryRow(ctx, query).Scan(
		&statistics.Total,
		&statistics.Published,
		&statistics.Pending,
		&statistics.Rejected,
	)
	if err != nil {
		log.Println("Error getting ad statistics:", err)
		return statistics, repoerr.ErrGettingStatistics
	}

	return statistics, nil
}
