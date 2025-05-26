package usecaseerr

var (
	ErrAccessDenied      = Error("access denied")
	ErrGettingAllAds     = Error("error getting all ads from database")
	ErrGettingAdByID     = Error("error getting ad by ID from database")
	ErrNoAds             = Error("no ads found")
	ErrDeletingAd        = Error("error deleting ad from database")
	ErrApprovingAd       = Error("error approving ad")
	ErrRejectingAd       = Error("error rejecting ad")
	ErrGettingStatistics = Error("error getting ad statistics")
)
