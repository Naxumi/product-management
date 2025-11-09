package product

import "errors"

var (
	// Product Errors
	ErrProductNotFound      = errors.New("product not found")
	ErrProductSKUExists     = errors.New("product with this SKU already exists")
	ErrInvalidProductStatus = errors.New("invalid product status")
	ErrProductIDRequired    = errors.New("product ID is required")
	ErrInvalidPrice         = errors.New("invalid price")
	ErrInvalidStock         = errors.New("invalid stock")
	ErrInvalidImageFormat   = errors.New("invalid image format, only JPG, JPEG, PNG, GIF are allowed")
	ErrImageTooLarge        = errors.New("image file size exceeds maximum limit of 5MB")
	ErrImageRequired        = errors.New("image file is required")
)
