package order

import (
	"context"
	"orders/internal/domain/inventory"
	"orders/internal/domain/payment"
	"orders/internal/domain/product"
)

type ProductServiceContract interface {
	GetProductsByIds(ctx context.Context, productIds []string) ([]product.Product, error)
}

type InventoryServiceContract interface {
	BulkUpdateStock(ctx context.Context, dto inventory.BulkUpdateStockRequest) error
}

type PaymentServiceContract interface {
	CreateInvoice(ctx context.Context, dto payment.CreateInvoiceInput) (*payment.CreateInvoiceResponse, error)
}
