package product

import "context"

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProducts(ctx context.Context) ([]Product, error)
	GetProductsByIds(ctx context.Context, productIds []string) ([]Product, error)
	GetProduct(ctx context.Context, productId string) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, productId string) error
}
