package entities

import "context"

type UserRepository interface {
	Regsister(ctx context.Context, user *User) error
	Login(ctx context.Context, phone, password string) (*User, error)
}
