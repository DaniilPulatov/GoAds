package entities

import "time"

type Token struct {
	ExpiresAt time.Time
	Token     string
	UserID    string
}
