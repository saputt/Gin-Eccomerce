package payment

type CreateInvoiceInput struct {
	OrderID string
	Amount  float64
}

type CreateInvoiceResponse struct {
	InvoiceID  string
	InvoiceURL string
	Status     string
}
