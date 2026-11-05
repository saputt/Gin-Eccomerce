package product

import "time"

type CreateProductRequest struct {
	Name        string        `json:"name" binding:"required,min=5,max=50"`
	Price       float32       `json:"price" binding:"required,min=1"`
	Stock       uint          `json:"stock" binding:"omitempty,min=1"`
	Description string        `json:"description"`
	Status      ProductStatus `json:"status" binding:"oneof=ACTIVE ARCHIVE"`
}

type UpdateProductRequest struct {
	Name        string        `json:"name" binding:"omitempty,min=5,max=50"`
	Price       float32       `json:"price" binding:"omitempty,min=1"`
	Description string        `json:"description" binding:"omitempty"`
	Status      ProductStatus `json:"status" binding:"omitempty,oneof=ACTIVE ARCHIVE"`
}

type ProductResponse struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProductStatus `json:"status"`
	Price       float32       `json:"price"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	SellerId    string        `json:"seller_id"`
}

func ToProductResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		SellerId:    p.SellerId,
	}
}

func ToProductResponseList(products []Product) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		if res := ToProductResponse(&p); res != nil {
			responses = append(responses, *res)
		}
	}
	return responses
}
