package product

import (
	"context"
	"orders/internal/domain/inventory"

	"github.com/google/uuid"
)

type ProductService struct {
	repo             ProductRepository
	inventoryService InventoryServiceContract
}

func NewProductService(repo ProductRepository, inventoryService InventoryServiceContract) *ProductService {
	return &ProductService{
		repo:             repo,
		inventoryService: inventoryService,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, dto CreateProductRequest, userId string) (*ProductResponse, error) {
	product := Product{
		Id:          uuid.NewString(),
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		Price:       dto.Price,
		SellerId:    userId,
	}

	err := s.repo.CreateProduct(ctx, &product)

	if err != nil {
		return nil, err
	}

	inventory := inventory.Inventories{
		ProductId: product.Id,
		Stock:     dto.Stock | 0,
	}

	_, err = s.inventoryService.CreateInventory(ctx, &inventory)
	if err != nil {
		return nil, err
	}

	return ToProductResponse(&product), nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, dto *UpdateProductRequest, productId string, userId string) (*ProductResponse, error) {
	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		return nil, err
	}

	if product.SellerId != userId {
		return nil, ErrUnauthorized
	}

	if dto.Description != "" {
		product.Description = dto.Description
	}

	if dto.Name != "" {
		product.Name = dto.Name
	}

	if dto.Price != 0 {
		product.Price = dto.Price
	}

	if dto.Status != "" {
		product.Status = dto.Status
	}

	err = s.repo.UpdateProduct(ctx, product)
	return ToProductResponse(product), err
}

func (s *ProductService) GetProduct(ctx context.Context, productId string) (*ProductResponse, error) {
	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		return nil, err
	}
	return ToProductResponse(product), nil
}

func (s *ProductService) GetProducts(ctx context.Context) ([]ProductResponse, error) {
	products, err := s.repo.GetProducts(ctx)
	return ToProductResponseList(products), err
}

func (s *ProductService) DeleteProduct(ctx context.Context, productId string, userId string) error {
	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		return ErrProductNotFound
	}

	if product.SellerId != userId {
		return ErrUnauthorized
	}

	err = s.repo.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}
