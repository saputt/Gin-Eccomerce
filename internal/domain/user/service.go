package user

import (
	"context"
	"fmt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
	_, err := s.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, ErrEmailAlreadyExist
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*User, error) {
	fmt.Println("masuk ke service")

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
