package auth

import (
	"context"
	"orders/internal/domain/user"
)

type UserServiceContract interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) (*user.User, error)
}
