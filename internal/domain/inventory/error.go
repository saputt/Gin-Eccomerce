package inventory

import "errors"

var (
	ErrInventoryNotFound = errors.New("Inventory doesnt exist")
	ErrOutOfStock        = errors.New("Out Of Stock")
)
