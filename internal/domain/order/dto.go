package order

type OrderItemRequest struct {
	ProductId string `json:"product_id"`
	Quantity  uint16 `json:"quantity" binding:"min=1"`
}

type CreateOrderRequest struct {
	Products []OrderItemRequest `json:"products" binding:"required,dive"`
}

type CreateOrderResponse struct {
	Id         string     `json:"id" gorm:"primaryKey"`
	Subtotal   float64    `json:"subtotal" gorm:"not null"`
	Suborders  []Suborder `json:"suborders"`
	PaymentUrl string     `json:"payment_url"`
}
