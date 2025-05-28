package ad

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
)

func (r adRepo) Create(ctx context.Context, ad *entities.Ad) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO ads(
            author_id, title, description, category_id, 
            status, is_active, created_at, updated_at
        ) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		ad.AuthorID, ad.Title, ad.Description, ad.CategoryID,
		ad.Status, ad.IsActive, ad.CreatedAt, ad.UpdatedAt)
	if err != nil {
		r.logger.ERROR("while inserting into ads:", err)
		return repoerr.ErrInsert
	}
	r.logger.INFO("Ad created successfully", ad.ID)
	return nil
}

func (r adRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	var ad entities.Ad
	err := r.db.QueryRow(ctx, `
		SELECT 
		    id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at
		FROM ads
		WHERE id = $1`, id).
		Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ERROR("No ad found with ID: ", id)
			return nil, repoerr.ErrAdNotFound
		}
		r.logger.ERROR("Error selecting ad: ", err)
		return nil, repoerr.ErrGettingAdByID
	}
	r.logger.INFO("AD retrived by ID successfully, ID: ", ad.ID)
	return &ad, nil
}

func (r adRepo) GetByUserID(ctx context.Context, userID string) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
		    id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at
		FROM ads
		WHERE author_id = $1`, userID)
	if err != nil {
		r.logger.ERROR("Error getting user ads: ", err)
		return nil, repoerr.ErrGettingAdsByUserID
	}
	defer rows.Close()

	var ads []entities.Ad
	for rows.Next() {
		var ad entities.Ad
		if err = rows.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
			log.Println("Scan error:", err)
			return nil, repoerr.ErrScan
		}
		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		r.logger.ERROR("Error iterating rows:", err)
		return nil, repoerr.ErrScan
	}
	r.logger.INFO("ADs retrived by user ID successfully, user ID: ", userID)
	return ads, nil
}

func (r adRepo) GetAll(ctx context.Context) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    id, author_id, title, description, category_id, 
			status, is_active, created_at, updated_at
		FROM ads
	`)
	if err != nil {
		r.logger.ERROR("Error getting ads:", err)
		return nil, repoerr.ErrGettingAllAds
	}

	defer rows.Close()

	var ads []entities.Ad
	for rows.Next() {
		var ad entities.Ad
		if err = rows.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.CategoryID, &ad.Status,
			&ad.IsActive, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
			r.logger.ERROR("Scan error: ", err)
			return nil, repoerr.ErrScan
		}
		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		r.logger.ERROR("Error iterating rows: ", err)
		return nil, repoerr.ErrScan
	}
	r.logger.INFO("ADs retrived successfully ")
	return ads, nil
}

func (r adRepo) Update(ctx context.Context, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET title = $1, description = $2, category_id = $3,
			status = $4, is_active = $5, updated_at = $6
		WHERE id = $7;`, ad.Title, ad.Description, ad.CategoryID,
		ad.Status, ad.IsActive, ad.UpdatedAt, ad.ID)
	if err != nil {
		r.logger.ERROR("Error updating ad: ", err)
		return repoerr.ErrUpdate
	}

	if row.RowsAffected() == 0 {
		r.logger.ERROR("No ad found with ID: ", ad.ID)
		return repoerr.ErrAdNotFound
	}
	r.logger.INFO("AD updated successfully, ID: ", ad.ID)
	return nil
}

func (r adRepo) Delete(ctx context.Context, id int) error {
	row, err := r.db.Exec(ctx, `
		DELETE FROM ads
		WHERE id = $1;`, id)
	if err != nil {
		r.logger.ERROR("Error deleting ad: ", err)
		return repoerr.ErrDelete
	}
	if row.RowsAffected() == 0 {
		r.logger.ERROR("No ad found with ID: ", id)
		return repoerr.ErrAdNotFound
	}
	r.logger.INFO("Ad deleted successfully: ", id)
	return nil
}

func (r adRepo) Approve(ctx context.Context, id int, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET status = $1, is_active = $2, updated_at = $3
		WHERE id = $4;`,
		ad.Status, ad.IsActive, ad.UpdatedAt, id)
	if err != nil {
		r.logger.ERROR("Error approving ad: ", err)
		return repoerr.ErrApproval
	}
	if row.RowsAffected() == 0 {
		r.logger.ERROR("No ad found with ID: ", id)
		return repoerr.ErrAdNotFound
	}
	r.logger.INFO("Ad approved successfully, ID: ", id)
	return nil
}

func (r adRepo) Reject(ctx context.Context, id int, ad *entities.Ad) error {
	row, err := r.db.Exec(ctx, `
		UPDATE ads
		SET status = $1, rejection_reason = $2, is_active = $3, updated_at = $4
		WHERE id = $5;`, ad.Status, ad.RejectionReason, ad.IsActive, ad.UpdatedAt, id)
	if err != nil {
		r.logger.ERROR("Error rejecting ad: ", err)
		return repoerr.ErrRejection
	}
	if row.RowsAffected() == 0 {
		r.logger.ERROR("No ad found with ID: ", id)
		return repoerr.ErrAdNotFound
	}
	r.logger.INFO("Ad rejected successfully, ID: ", id)
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
		r.logger.ERROR("Error getting ad statistics: ", err)
		return statistics, repoerr.ErrGettingStatistics
	}
	r.logger.INFO("Statistics retrieved successfully ")
	return statistics, nil
}
