package repository

import (
	"context"
	"errors"
	"orders/internal/common/helper"
	"orders/internal/domain/inventory"

	"gorm.io/gorm"
)

type GormInventoryRepository struct {
	db *gorm.DB
}

func NewGormInventoryRepository(db *gorm.DB) *GormInventoryRepository {
	return &GormInventoryRepository{
		db: db,
	}
}

func (r *GormInventoryRepository) CreateInventory(ctx context.Context, inv *inventory.Inventories) error {
	return r.db.WithContext(ctx).Save(inv).Error
}

func (r *GormInventoryRepository) UpdateStock(ctx context.Context, productId string, quantity int) (*inventory.Inventories, error) {
	var inv inventory.Inventories

	err := r.db.WithContext(ctx).
		Model(inventory.Inventories{}).
		Where("product_id = ?", productId).
		Update("stock", gorm.Expr("quantity + ?", quantity)).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, inventory.ErrInventoryNotFound
		}
		return nil, err
	}

	err = r.db.WithContext(ctx).Find(&inv).Error
	if err != nil {
		return nil, err
	}

	return &inv, err
}

func (r *GormInventoryRepository) FindByProduct(ctx context.Context, productId string) (*inventory.Inventories, error) {
	var inv inventory.Inventories
	err := r.db.WithContext(ctx).Find(&inv, "id = ?", productId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, inventory.ErrInventoryNotFound
		}
	}
	return &inv, nil
}

func (r *GormInventoryRepository) BulkUpdateStock(ctx context.Context, dto inventory.BulkUpdateStockRequest) error {
	db := helper.ExtractTX(ctx)
	if db == nil {
		db = r.db.WithContext(ctx)
	}

	for _, item := range dto.Products {
		res := db.Model(inventory.Inventories{}).
			Where("product_id = ? AND stock >= ?", item.ProductId, item.Quantity).
			Update("stock", gorm.Expr("stock - ?", item.Quantity))

		if res.RowsAffected == 0 {
			return inventory.ErrOutOfStock
		}

		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}
