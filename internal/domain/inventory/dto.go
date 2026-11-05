package inventory

type UpdateStockRequest struct {
	ProductId string `json:"product_id"`
	Quantity  uint16 `json:"quantity" binding:"min=1"`
}

type BulkUpdateStockRequest struct {
	Products []UpdateStockRequest `json:"products" binding:"required,dive"`
}
