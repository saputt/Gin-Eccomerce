package payment

import (
	"time"

	"gorm.io/datatypes"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "PENDING"
	StatusPaid    PaymentStatus = "PAID"
	StatusExpired PaymentStatus = "EXPIRED"
	StatusFailed  PaymentStatus = "FAILED"
)

type Payments struct {
	Id      string `json:"id" gorm:"primaryKey"`
	OrderId string `json:"order_id" gorm:"uniqueIndex"`

	XenditInvoiceId      string `json:"xendit_invoice_id" gorm:"uniqueIndex,not null"`
	XenditPaymentMethod  string `json:"xendit_payment_method"`
	XenditPaymentChannel string `json:"xendit_payment_channel"`

	Amount     float32 `json:"amount" gorm:"not null"`
	PaymentUrl string  `json:"payment_url"`

	Status      PaymentStatus  `json:"status" gorm:"not null,default:PENDING"`
	RawResponse datatypes.JSON `json:"raw_response" gorm:"type:json"`
	ExpiredAt   time.Time      `json:"expired_at"`
	PaidAt      time.Time      `json:"paid_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
