package entities

import "time"

type RefreshToken struct {
	ExpiredAt time.Time
	Token     string
	UserID    int
}
