package entities

import (
	"time"
)

// Status - represent allowed statuses to be used for ad.
type Status string

// The only allowed statuses.
const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

// Ad - represent ad, contains reference to the user (AuthorID) and reference to the category of the ad.
type Ad struct {
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          Status
	Title           string
	Location        string
	Description     string
	RejectionReason string
	AuthorID        string
	CategoryID      int
	ID              int
	IsActive        bool
}

// AdFile - represents file that user will attach to the ad, contains reference to the ad (AdID)
// and path to the file (URL).
type AdFile struct {
	CreatedAt time.Time
	FileName  string
	URL       string
	AdID      int
	ID        int
}

type AdFilter struct {
	DateFrom   time.Time
	DateTo     time.Time
	Status     string
	UserID     string
	CategoryID int
	Limit      int
	Page       int
}

type AdStatistics struct {
	Total     int
	Published int
	Pending   int
	Rejected  int
}
