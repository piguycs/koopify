package checkout

import "errors"

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrInsufficientStock      = errors.New("insufficient stock for product")
	ErrProductNotFound        = errors.New("product not found")
	ErrInvalidQuantity        = errors.New("invalid quantity")
	ErrCheckoutCreationFailed = errors.New("failed to create checkout session")
	ErrOrderUpdateFailed      = errors.New("failed to update order")
)
