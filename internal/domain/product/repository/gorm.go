package repository

import (
	"context"
	"errors"
	"orders/internal/product"

	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{
		db: db,
	}
}

func (r *GormProductRepository) CreateProduct(ctx context.Context, product *product.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *GormProductRepository) GetProducts(ctx context.Context) ([]product.Product, error) {
	var prods []product.Product

	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&prods).Error
	if err != nil {
		return nil, err
	}

	return prods, nil
}

func (r *GormProductRepository) GetProduct(ctx context.Context, productId string) (*product.Product, error) {
	var prod product.Product

	err := r.db.WithContext(ctx).Find(&prod, "id = ?", productId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	return &prod, nil
}

func (r *GormProductRepository) UpdateProduct(ctx context.Context, prod *product.Product) error {
	return r.db.WithContext(ctx).Save(prod).Error
}

func (r *GormProductRepository) DeleteProduct(ctx context.Context, productId string) error {
	res := r.db.WithContext(ctx).Delete(&product.Product{}, "id = ?", productId)

	if res.RowsAffected == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

func (r *GormProductRepository) GetProductsByIds(ctx context.Context, productIds []string) ([]product.Product, error) {
	var prods []product.Product
	err := r.db.WithContext(ctx).
		Where("id IN ?", productIds).
		Find(&prods).Error

	if err != nil {
		return nil, err
	}

	return prods, nil
}
