package adfile

import (
	"ads-service/internal/domain/entities"
	repoerr "ads-service/internal/errs/repoErr"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx"
)

func (r adFileRepo) Create(ctx context.Context, file *entities.AdFile) (int, error) {
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

func (r adFileRepo) Delete(ctx context.Context, file *entities.AdFile) (string, error) {
	var (
		selectQuery = `SELECT url FROM ad_files WHERE id = $1 RETURNING url`
		delteQuery  = `DELETE FROM ad_files WHERE id = $1 AND ad_id = $2`
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

	if _, err := r.db.Exec(ctx, delteQuery, file.ID, file.AdID); err != nil {
		log.Println("Error deleting ad file:", err)
		return "", repoerr.ErrFileDeletion
	}
	return url, nil
}

func (r adFileRepo) GetAll(ctx context.Context, adID int) ([]entities.AdFile, error) {
	var (
		selectQuery = `SELECT id, ad_id, file_name, url, created_at FROM ad_files WHERE ad_id = $1`
		files       []entities.AdFile
	)

	rows, err := r.db.Query(ctx, selectQuery, adID)
	if err != nil {
		log.Println("Error selecting ad files:", err)
		return nil, repoerr.ErrFileSelection
	}
	defer rows.Close()

	for rows.Next() {
		var file entities.AdFile
		if err := rows.Scan(&file.ID, &file.AdID, &file.FileName, &file.URL, &file.CreatedAt); err != nil {
			log.Println("Error scanning ad file:", err)
			return nil, repoerr.ErrJSONUnmarshal
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over ad files:", err)
		return nil, repoerr.ErrSelection
	}

	return files, nil
}
