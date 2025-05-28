package entities

import "time"

// Role - for role-based auth.
type Role string

// The only allowed roles.
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         Role
	FName        string
	LName        string
	PasswordHash string
	Password     string
	Phone        string
	ID           string
}
