package repoerr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrSelection = Error("error selecting ads from database")
	ErrInsert    = Error("error inserting ad into database")
	ErrUpdate    = Error("error updating ad in database")
	ErrDelete    = Error("error deleting ad from database")
	ErrNoRows    = Error("no rows found in database for the query")

	ErrFileSelection = Error("error selecting ad files from database")
	ErrFileInsertion = Error("error inserting ad file into database")
	ErrFileDeletion  = Error("error deleting ad file from database")
	ErrFileNotFound  = Error("ad file not found in database")

	ErrJSONUnmarshal = Error("error unmarshalling JSON data from database")
)
