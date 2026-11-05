package order

import (
	"context"
	"fmt"
	"orders/internal/common/helper"
	"orders/internal/domain/inventory"
	"orders/internal/domain/payment"
	"orders/internal/domain/product"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repo             OrderRepository
	productService   ProductServiceContract
	inventoryService InventoryServiceContract
	paymentService   PaymentServiceContract
	db               *gorm.DB
}

func NewOrderService(repo OrderRepository, productService ProductServiceContract) *OrderService {
	return &OrderService{
		repo:           repo,
		productService: productService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, dto CreateOrderRequest, buyerId string) (*CreateOrderResponse, error) {
	var productsIds []string

	for _, val := range dto.Products {
		productsIds = append(productsIds, val.ProductId)
	}

	products, err := s.productService.GetProductsByIds(ctx, productsIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]product.Product)
	for _, p := range products {
		productMap[p.Id] = p
	}

	var subtotal float32
	orderId := uuid.New()

	suborderMap := make(map[string]Suborder)

	for _, item := range dto.Products {
		prod, exists := productMap[item.ProductId]
		if !exists {
			return nil, fmt.Errorf("product with id %s: %w", item.ProductId, ErrProductNotFound)
		}

		total := prod.Price * float32(item.Quantity)

		subord, exists := suborderMap[prod.SellerId]
		if !exists {
			subord = Suborder{
				Id:           uuid.New().String(),
				SellerId:     prod.SellerId,
				OrderId:      orderId.String(),
				Subtotal:     0,
				ShippingCost: 0,
				OrderItems:   []OrderItem{},
				Status:       StatusPending,
			}
			suborderMap[prod.SellerId] = subord
		}

		subord.Subtotal += total
		subord.OrderItems = append(subord.OrderItems, OrderItem{
			Id:         uuid.New().String(),
			ProductId:  prod.Id,
			Total:      float64(total),
			Price:      float64(prod.Price),
			Quantity:   int(item.Quantity),
			SuborderId: subord.Id,
		})
	}

	var suborders []Suborder
	for _, val := range suborderMap {
		suborders = append(suborders, val)
	}

	ord := Order{
		Id:        orderId.String(),
		Subtotal:  float64(subtotal),
		BuyerId:   buyerId,
		Suborders: suborders,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := helper.InjectTX(ctx, tx)

		if err := s.repo.CreateOrderTX(txCtx, &ord); err != nil {
			return err
		}

		var invUpdate inventory.BulkUpdateStockRequest

		for _, val := range dto.Products {
			invUpdate.Products = append(invUpdate.Products, inventory.UpdateStockRequest{
				ProductId: val.ProductId,
				Quantity:  val.Quantity,
			})
		}

		return s.inventoryService.BulkUpdateStock(txCtx, invUpdate)
	})

	if err != nil {
		return nil, err
	}

	payment, err := s.paymentService.CreateInvoice(ctx, payment.CreateInvoiceInput{
		OrderID: ord.Id,
		Amount:  ord.Subtotal,
	})

	if err != nil {
		return nil, err
	}

	return &CreateOrderResponse{
		PaymentUrl: payment.InvoiceURL,
		Subtotal:   ord.Subtotal,
		Suborders:  ord.Suborders,
	}, nil
}

func (s *OrderService) FindOrder(ctx context.Context, userId string, orderId string) (*Order, error) {
	ord, err := s.repo.FindOrder(ctx, orderId)
	if err != nil {
		return nil, err
	}

	if ord.BuyerId != userId {
		return nil, ErrUnauthorized
	}

	return ord, nil
}

func (s *OrderService) FindSuborder(ctx context.Context, userId string, suborderId string) (*Suborder, error) {
	suborder, err := s.repo.FindSuborder(ctx, suborderId)
	if err != nil {
		return nil, err
	}

	if suborder.Order.BuyerId != userId {
		return nil, ErrUnauthorized
	}

	return suborder, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, userId string, orderId string) error {
	ord, err := s.FindOrder(ctx, userId, orderId)
	if err != nil {
		return err
	}

	maxCancelTime := 30 * time.Minute

	if time.Since(ord.CreatedAt) > maxCancelTime {
		return ErrCancelDeclined
	}

	var subordersIds []string
	for _, val := range ord.Suborders {
		subordersIds = append(subordersIds, val.Id)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := helper.InjectTX(ctx, tx)

		err = s.repo.CancelOrder(txCtx, subordersIds)
		if err != nil {
			return err
		}

		stockMap := make(map[string]uint16)
		for _, val := range ord.Suborders {
			for _, item := range val.OrderItems {
				stockMap[item.ProductId] += uint16(item.Quantity)
			}
		}

		var invUpdate inventory.BulkUpdateStockRequest

		for productId, quantity := range stockMap {
			invUpdate.Products = append(invUpdate.Products, inventory.UpdateStockRequest{
				ProductId: productId,
				Quantity:  quantity,
			})
		}

		err = s.inventoryService.BulkUpdateStock(txCtx, invUpdate)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelSuborder(ctx context.Context, userId string, suborderId string) error {
	suborder, err := s.FindSuborder(ctx, userId, suborderId)
	if err != nil {
		return err
	}

	maxCancelTime := 6 * time.Hour

	if (suborder.Status != StatusPending || suborder.Status != StatusProcess) && time.Since(suborder.Order.Payment.PaidAt) > maxCancelTime {
		return ErrCancelDeclined
	}

	err = s.repo.CancelSuborder(ctx, suborderId)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) FindMySuborderForBuyer(ctx context.Context, userId string) ([]Order, error) {
	return s.repo.FindMySuborderForBuyer(ctx, userId)
}

func (s *OrderService) FindMySuborderForSeller(ctx context.Context, userId string) ([]Order, error) {
	return s.repo.FindMySuborderForSeller(ctx, userId)
}
