package payment

import "context"

type PaymentService interface {
	CreateInvoice(ctx context.Context, dto CreateInvoiceInput) (*CreateInvoiceResponse, error)
}
