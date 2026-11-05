package payment

import "context"

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payments) error
}
