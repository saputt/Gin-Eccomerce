package order

import "errors"

var (
	ErrOrderNotFound    = errors.New("Order not found")
	ErrSuborderNotFound = errors.New("Suborder not found")
	ErrUnauthorized     = errors.New("You're not authorized to access this order")
	ErrProductNotFound  = errors.New("Product not found")
	ErrCancelDeclined   = errors.New("Cancel declined, out of time")
)
