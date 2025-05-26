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
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Role         Role      `json:"role"`
	FName        string    `json:"first_name"`
	LName        string    `json:"last_name"`
	PasswordHash string    `json:"password_hash" `
	Password     string    `json:"password" binding:"required"` // temp field for registration
	// will be removed after rehashing the password
	Phone string `json:"phone" binding:"required"`
	ID    string `json:"id"`
}
