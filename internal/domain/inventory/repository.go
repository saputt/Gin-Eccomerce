package inventory

import "context"

type InventoryRepository interface {
	CreateInventory(ctx context.Context, inventory *Inventories) error
	FindByProduct(ctx context.Context, productId string) (*Inventories, error)
	UpdateStock(ctx context.Context, productId string, quantity int) (*Inventories, error)
	BulkUpdateStock(ctx context.Context, dto BulkUpdateStockRequest) error
}
