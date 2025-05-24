package entities

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

// Status - represent alowed statuses to be used for ad.
type Status string

// The only allowed statuses.
const (
	StatusDraft    Status = "draft"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

// Ad - represent ad, contains refernce to the user (AuthorID) and reference to the category oo the ad.
type Ad struct {
	CreatedAt       time.Time
	UpdatedAt       time.Time
	AuthorID        uuid.UUID
	Files           []AdFile // files that are attached to the ad
	Title           string
	Description     string
	Status          Status
	RejectionReason string
	Location        string
	ID              int
	CategoryID      int
	IsActive        bool
}

// AdFile - represents file that user will attach to the ad, contains reference to the ad (AdID) and path to the file (URL).
type AdFile struct {
	CreatedAt time.Time
	FileName  string
	URL       string
	ID        int
	AdID      int
}
