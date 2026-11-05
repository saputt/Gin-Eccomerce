package order

import (
	"context"
)

type OrderRepository interface {
	CreateOrderTX(ctx context.Context, order *Order) error
	CancelOrder(ctx context.Context, suborderIds []string) error
	CancelSuborder(ctx context.Context, suborderId string) error
	FindMySuborderForBuyer(ctx context.Context, userId string) ([]Order, error)
	FindMySuborderForSeller(ctx context.Context, userId string) ([]Order, error)
	FindOrder(ctx context.Context, orderId string) (*Order, error)
	FindSuborder(ctx context.Context, suborderId string) (*Suborder, error)
}
