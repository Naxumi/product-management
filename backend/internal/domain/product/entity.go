package product

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "Active"
	ProductStatusInactive ProductStatus = "Inactive"
)

type Product struct {
	ID          int64
	SKU         string
	Name        string
	Description *string
	Price       decimal.Decimal
	Stock       int
	Category    string
	Status      ProductStatus
	ImageURL    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
