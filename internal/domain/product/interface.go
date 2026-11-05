package product

import (
	"context"
	"orders/internal/domain/inventory"
)

type InventoryServiceContract interface {
	CreateInventory(ctx context.Context, inv *inventory.Inventories) (*inventory.Inventories, error)
}
