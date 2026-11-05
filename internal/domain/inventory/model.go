package inventory

import (
	"time"
)

type Inventories struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Stock     uint      `json:"stock" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"`
	ProductId string    `json:"product_id" gorm:"nut null"`
}
