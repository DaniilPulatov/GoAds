package ad

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdRepo struct {
	db *pgxpool.Pool
}

func NewAdRepo(db *pgxpool.Pool) *AdRepo {
	return &AdRepo{db: db}
}

var (
	adSelect = `SELECT 
		ads.id AS ad_id, ads.title, ads.description,
		ads.author_id, ads.category_id, ads.status,
		ads.is_active, ads.created_at, ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
	FROM ads
	LEFT JOIN ad_files ON ad_files.ad_id = ads.id
	GROUP BY ads.id
	ORDER BY ads.id;
	`
	adSelectLimit = `SELECT
		ads.id AS ad_id, ads.title, ads.description,
		ads.author_id, ads.category_id, ads.status, 
		ads.is_active, ads.created_at, ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
		FROM ads
		LEFT JOIN ad_files ON ad_files.ad_id = ads.id
		GROUP BY ads.id
		ORDER BY ads.id
		Limit $1;`
	adSelectByID = `SELECT
		ads.id AS ad_id, ads.title, ads.description,
		ads.author_id, ads.category_id, ads.status, 
		ads.is_active, ads.created_at, ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
		FROM ads
		LEFT JOIN ad_files ON ad_files.ad_id = ads.id
		WHERE ads.id = $1
		GROUP BY ads.id
		ORDER BY ads.id;`
	adSelectByUserID = `SELECT
		ads.id AS ad_id, ads.title, ads.description,
		ads.author_id, ads.category_id, ads.status, 
		ads.is_active, ads.created_at, ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
		FROM ads
		LEFT JOIN ad_files ON ad_files.ad_id = ads.id
		WHERE ads.author_id = $1
		GROUP BY ads.id
		ORDER BY ads.id;`
	adSelectByUserIDLimit = `SELECT
		ads.id AS ad_id, ads.title, ads.description,
		ads.author_id, ads.category_id, ads.status, 
		ads.is_active, ads.created_at, ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
		FROM ads
		LEFT JOIN ad_files ON ad_files.ad_id = ads.id
		WHERE ads.author_id = $1
		GROUP BY ads.id
		ORDER BY ads.id
		Limit $2;`
	adSelectByCat = `SELECT ads.id AS ad_id, ads.title,
		ads.description, ads.author_id, ads.category_id,
		ads.status, ads.is_active, ads.created_at,
		ads.updated_at,
		COALESCE(json_agg(
			json_build_object(
				'id', ad_files.id,
				'url', ad_files.url,
				'file_name', ad_files.file_name
			) ORDER BY ad_files.id
		) FILTER (WHERE ad_files.id IS NOT NULL), '[]') AS files
		FROM ads
		LEFT JOIN ad_files ON ad_files.ad_id = ads.id
		WHERE ads.category_id = $1
		GROUP BY ads.id
		ORDER BY ads.id;
	`
	adUpdate = ` UPDATE ads
		SET 
		title = $1,
		description = $2,
		location = $3,
		category_id = $4,
		status = $5,
		rejection_reason = $6,
		is_active = $7,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $8;`
)

func (r AdRepo) GetAll(ctx context.Context) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, adSelect)
	if err != nil {
		log.Println("Query error:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ads found")
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}

	return r.queryRows(rows)
}
func (r AdRepo) GetSome(ctx context.Context, limit int) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, adSelectLimit, limit)
	if err != nil {
		log.Println("Query error:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ads found")
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}
	return r.queryRows(rows)
}

func (r AdRepo) GetByID(ctx context.Context, id int) (*entities.Ad, error) {
	row := r.db.QueryRow(ctx, adSelectByID, id)
	ad, err := r.queryRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ad found with ID:", id)
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}
	return ad, nil

}
func (r AdRepo) GetByUserID(ctx context.Context, userUUid string) (*entities.Ad, error) {
	row := r.db.QueryRow(ctx, adSelectByUserID, userUUid)
	ad, err := r.queryRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ad found with ID:", userUUid)
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}
	return ad, nil
}
func (r AdRepo) GetSomeByUserID(ctx context.Context, userID string, limit int) ([]entities.Ad, error) {
	rows, err := r.db.Query(ctx, adSelectLimit, userID, limit)
	if err != nil {
		log.Println("Query error:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ads found")
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}
	return r.queryRows(rows)
}

func (r AdRepo) GetByCat(ctx context.Context, catID int) (*entities.Ad, error) {
	row := r.db.QueryRow(ctx, adSelectByCat, catID)
	ad, err := r.queryRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ads found with category ID:", catID)
			return nil, repoerr.ErrNoRows
		}
		log.Println("Database selection error:", err)
		return nil, repoerr.ErrSelection
	}
	return ad, nil
}

func (r AdRepo) Update(ctx context.Context, ad *entities.Ad) error {
	if _, err := r.db.Exec(ctx, adUpdate,
		ad.Title,
		ad.Description,
		ad.Location,
		ad.CategoryID,
		entities.StatusDraft, // status is set to draft on update
		ad.RejectionReason,
		true,
		ad.ID,
	); err != nil {
		log.Println("Update error:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No ad found with ID:", ad.ID)
			return repoerr.ErrNoRows
		}
		log.Println("Database update error:", err)
		return repoerr.ErrUpdate
	}
	return nil
}
func (r AdRepo) Delete(ctx context.Context, adID int) error {
	// Implementation will go here
	return nil
}

// queryRows - helper function that processes the rows returned from the database query and unmarshals the JSON data for files.
func (r AdRepo) queryRows(rows pgx.Rows) ([]entities.Ad, error) {
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
		if err := json.Unmarshal(filesJSON, &ad.Files); err != nil {
			log.Println("Unmarshal error:", err)
			return nil, repoerr.ErrJsonUnmarshal
		}

		if ad.Files == nil {
			ad.Files = []entities.AdFile{} // Ensure Files is initialized to an empty slice if nil
		}
		ads = append(ads, ad)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
		return nil, repoerr.ErrSelection
	}
	return ads, nil
}

func (r AdRepo) queryRow(row pgx.Row) (*entities.Ad, error) {
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
	if err := json.Unmarshal(filesJSON, &ad.Files); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, repoerr.ErrJsonUnmarshal
	}

	if ad.Files == nil {
		ad.Files = []entities.AdFile{}
	}

	return &ad, nil
}
