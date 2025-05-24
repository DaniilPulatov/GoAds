package entities

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)
// Role - for role based auth.
type Role string

// The only allowed roles.
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	CretedAt     time.Time
	UpdatedAt    time.Time
	ID           uuid.UUID
	Role         Role
	Username     string
	PasswordHash string
	Phone        string
}
