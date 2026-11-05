package repository

import (
	"context"
	"errors"
	"orders/internal/common/helper"
	"orders/internal/domain/order"

	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{
		db: db,
	}
}

func (r *GormOrderRepository) CreateOrderTX(ctx context.Context, order *order.Order) error {
	tx := helper.ExtractTX(ctx)
	return tx.Create(order).Error
}

func (r *GormOrderRepository) CancelOrder(ctx context.Context, subordersIds []string) error {
	db := helper.ExtractTX(ctx)
	if db == nil {
		db = r.db.WithContext(ctx)
	}

	res := db.
		Model(&order.Suborder{}).
		Where("id IN ?", subordersIds).
		Update("status", order.StatusCancel)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return order.ErrOrderNotFound
	}

	return nil
}

func (r *GormOrderRepository) FindMySuborderForBuyer(ctx context.Context, buyerId string) ([]order.Suborder, error) {
	var ord []order.Suborder
	err := r.db.WithContext(ctx).Find(&ord, "buyer_id = ?", buyerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	return ord, nil
}

func (r *GormOrderRepository) FindMySuborderForSeller(ctx context.Context, sellerId string) ([]order.Suborder, error) {
	var ord []order.Suborder
	err := r.db.WithContext(ctx).Find(&ord, "seller_id = ?", sellerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	return ord, nil
}

func (r *GormOrderRepository) FindOrder(ctx context.Context, orderId string) (*order.Order, error) {
	var ord order.Order
	err := r.db.WithContext(ctx).
		Joins("Payments").
		Where("id = ?", orderId).
		First(&ord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	return &ord, nil
}

func (r *GormOrderRepository) FindSuborder(ctx context.Context, suborderId string) (*order.Suborder, error) {
	var suborder order.Suborder

	err := r.db.WithContext(ctx).
		Joins("Order").
		Where("id = ?", suborderId).
		First(&suborder).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrSuborderNotFound
		}
		return nil, err
	}

	return &suborder, nil
}

func (r *GormOrderRepository) CancelSuborder(ctx context.Context, suborderId string) error {
	tx := helper.ExtractTX(ctx)
	if tx == nil {
		tx = r.db
	}

	res := tx.
		Model(&order.Suborder{}).
		Where("id ? = ", suborderId).
		Update("status", order.StatusCancel)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return order.ErrOrderNotFound
	}

	return nil
}
