package inventory

import "context"

type InventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) CreateInventory(ctx context.Context, inv *Inventories) (*Inventories, error) {
	err := s.repo.CreateInventory(ctx, inv)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (s *InventoryService) FindByProduct(ctx context.Context, productId string) (*Inventories, error) {
	return s.repo.FindByProduct(ctx, productId)
}

func (s *InventoryService) UpdateStock(ctx context.Context, productId string, quantity uint16) (*Inventories, error) {
	return s.repo.UpdateStock(ctx, productId, int(quantity))
}

func (s *InventoryService) BulkUpdateStock(ctx context.Context, dto BulkUpdateStockRequest) error {
	return s.repo.BulkUpdateStock(ctx, dto)
}
