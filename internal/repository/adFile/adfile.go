package adfile

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/errs/repoerr"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r adFileRepo) Create(ctx context.Context, file *entities.AdFile) (int, error) {
	var (
		insertQuery = `INSERT INTO ad_files (ad_id, file_name, url) VALUES ($1, $2, $3) RETURNING id;`
		fileID      int
	)

	row := r.pool.QueryRow(ctx, insertQuery, file.AdID, file.FileName, file.URL)
	err := row.Scan(&fileID)
	if err != nil {
		r.logger.ERROR("Error scanning fileID:", err)
		return -1, repoerr.ErrFileInsertion
	}
	r.logger.INFO("Successfully created ad file with ID:", file.AdID)
	return fileID, nil
}

func (r adFileRepo) Delete(ctx context.Context, file *entities.AdFile) (string, error) {
	var (
		selectQuery = `SELECT url FROM ad_files WHERE id = $1`
		deleteQuery = `DELETE FROM ad_files WHERE url=$1 AND ad_id = $2;`
		url         string
	)

	row := r.pool.QueryRow(ctx, selectQuery, file.ID)
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ERROR("No ad file found with ID: ", file.ID)
			return "", repoerr.ErrFileNotFound
		}
		r.logger.ERROR("Error selecting ad file: ", err)
		return "", repoerr.ErrSelection
	}

	if _, err := r.pool.Exec(ctx, deleteQuery, url, file.AdID); err != nil {
		r.logger.ERROR("Error deleting ad file :", err)
		return "", repoerr.ErrFileDeletion
	}
	r.logger.INFO("Deleted ad file successfully", file)
	return url, nil
}

func (r adFileRepo) GetAll(ctx context.Context, adID int) ([]entities.AdFile, error) {
	var (
		selectQuery = `SELECT id, ad_id, file_name, url, created_at FROM ad_files WHERE ad_id = $1`
		files       []entities.AdFile
	)

	rows, err := r.pool.Query(ctx, selectQuery, adID)
	if err != nil {
		r.logger.ERROR("Error selecting ad files:", err)
		return nil, repoerr.ErrFileSelection
	}
	defer rows.Close()

	for rows.Next() {
		var file entities.AdFile
		if err := rows.Scan(&file.ID, &file.AdID, &file.FileName, &file.URL, &file.CreatedAt); err != nil {
			r.logger.ERROR("Error scanning ad file:", err)
			return nil, repoerr.ErrJSONUnmarshal
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		r.logger.ERROR("Error iterating over ad files:", err)
		return nil, repoerr.ErrSelection
	}
	r.logger.INFO("Retrieved all ad files", len(files))

	return files, nil
}
