package product

import "errors"

var (
	ErrProductNotFound = errors.New("Product not found")
	ErrUnauthorized    = errors.New("You're not authorized")
)
