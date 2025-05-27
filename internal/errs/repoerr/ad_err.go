package repoerr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrSelection         = Error("error selecting ads from database")
	ErrInsert            = Error("error inserting ad into database")
	ErrUpdate            = Error("error updating ad in database")
	ErrDelete            = Error("error deleting ad from database")
	ErrNoRows            = Error("no rows found in database for the query")
	ErrAdNotFound        = Error("no such ad in database")
	ErrUserNotHaveAds    = Error("user does not have any ads")
	ErrApproval          = Error("error approving ad")
	ErrRejection         = Error("error rejecting ad")
	ErrGettingStatistics = Error("error getting ad statistics")

	ErrGettingAllAds          = Error("error getting all ads from database")
	ErrGettingAdByID          = Error("error getting ad by ID from database")
	ErrGettingAdsByUserID     = Error("error getting ads by user ID from database")
	ErrGettingAdsByCategoryID = Error("error getting ads by category ID from database")
	ErrGettingAdsByStatus     = Error("error getting ads by status from database")
	ErrGettingAdsByDate       = Error("error getting ads by date from database")
	ErrGettingAdsByTitle      = Error("error getting ads by title from database")
	ErrGettingAdsByLocation   = Error("error getting ads by location from database")
	ErrScan                   = Error("error scanning ad from database")

	ErrFileSelection = Error("error selecting ad files from database")
	ErrFileInsertion = Error("error inserting ad file into database")
	ErrFileDeletion  = Error("error deleting ad file from database")
	ErrFileNotFound  = Error("ad file not found in database")

	ErrJSONUnmarshal = Error("error unmarshalling JSON data from database")
)
