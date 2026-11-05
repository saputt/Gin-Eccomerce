package product

import (
	"time"
)

type ProductStatus string

const (
	StatusActive  ProductStatus = "ACTIVE"
	StatusArchive ProductStatus = "ARCHIVE"
)

type Product struct {
	Id          string        `json:"id" gorm:"primaryKey;"`
	Name        string        `json:"name" gorm:"not null;"`
	Description string        `json:"description" gorm:"omitempty;"`
	Status      ProductStatus `json:"status" gorm:"not null;default:ACTIVE"`
	Price       float32       `json:"price" gorm:"not null;"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`

	SellerId string `json:"seller_id"`
}
