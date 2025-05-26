package entities

import "time"

type RefreshToken struct {
	ExpiresAt time.Time
	Token     string
	UserID    string
}
