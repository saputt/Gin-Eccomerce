package repository

import (
	"context"
	"orders/internal/common/helper"
	"orders/internal/payment"

	"gorm.io/gorm"
)

type PaymentGormRepository struct {
	db *gorm.DB
}

func NewPaymentGormRepository(db *gorm.DB) *PaymentGormRepository {
	return &PaymentGormRepository{
		db: db,
	}
}

func (r *PaymentGormRepository) CreatePayment(ctx context.Context, payment *payment.Payments) error {
	db := helper.ExtractTX(ctx)
	if db == nil {
		db = r.db
	}

	err := db.Save(&payment).Error
	if err != nil {
		return err
	}

	return nil
}
