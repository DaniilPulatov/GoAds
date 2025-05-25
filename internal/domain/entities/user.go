package entities

import uuid "github.com/jackc/pgtype/ext/gofrs-uuid"

// Role - for role-based auth.
type Role string

// The only allowed roles.
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	Role         Role
	FName        string
	LName        string
	PasswordHash string
	Phone        string
	ID           uuid.UUID
}
