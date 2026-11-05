package order

import (
	"orders/internal/domain/payment"
	"time"
)

type SuborderStatus string

const (
	StatusPending   SuborderStatus = "PENDING"
	StatusProcess   SuborderStatus = "PROCESS"
	StatusShip      SuborderStatus = "SHIP"
	StatusCompleted SuborderStatus = "COMPLETED"
	StatusCancel    SuborderStatus = "CANCEL"
)

type Order struct {
	Id       string  `json:"id" gorm:"primaryKey"`
	Subtotal float64 `json:"subtotal" gorm:"not null"`
	BuyerId  string  `json:"buyer_id" gorm:"not null"`

	CreatedAt time.Time         `json:"created_at"`
	Suborders []Suborder        `gorm:"foreignKey:OrderId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payment   *payment.Payments `json:"payment,omitempty" gorm:"foreignKey:OrderId;references:Id"`
}

type Suborder struct {
	Id           string         `json:"id" gorm:"primaryKey"`
	OrderId      string         `json:"order_id" gorm:"not null"`
	SellerId     string         `json:"seller_id" gorm:"not null"`
	Status       SuborderStatus `json:"status" gorm:"default:PENDING"`
	Subtotal     float32        `json:"subtotal"`
	ShippingCost float32        `json:"shipping_cost"`
	CreatedAt    time.Time      `json:"created_at"`

	Order      *Order      `gorm:"foreignKey:OrderId"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type OrderItem struct {
	Id         string  `json:"id" gorm:"primaryKey"`
	ProductId  string  `json:"product_id"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	Total      float64 `json:"total" gorm:"not null"`
	Price      float64 `json:"price" gorm:"not null"`
	SuborderId string  `json:"suborder_id"`
}
