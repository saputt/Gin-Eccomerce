package repository

import (
	"context"
	"errors"
	"orders/internal/user"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

func (r *GormUserRepository) CreateUser(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var users user.User

	err := r.db.WithContext(ctx).Find(&users, "email = ?", email).Error

	if err != nil {
		return nil, user.ErrEmailAlreadyExist
	}

	return &users, nil
}

func (r *GormUserRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	var users user.User

	err := r.db.WithContext(ctx).Find(&users, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}

	return &users, nil
}
