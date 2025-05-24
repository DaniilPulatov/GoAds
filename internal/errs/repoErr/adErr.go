package repoerr

import "errors"

var (
	ErrSelection = errors.New("error selecting ads from database")
	ErrInsert    = errors.New("error inserting ad into database")
	ErrUpdate    = errors.New("error updating ad in database")
	ErrDelete    = errors.New("error deleting ad from database")
	ErrNoRows   = errors.New("no rows found in database for the query")

	ErrJsonUnmarshal = errors.New("error unmarshalling JSON data from database")
)
