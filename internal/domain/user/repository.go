package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, userId string) (*User, error)
}
