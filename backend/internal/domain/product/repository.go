package product

import (
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product Product) (Product, error)
	GetByID(ctx context.Context, id int64) (Product, error)
	GetBySKU(ctx context.Context, sku string) (Product, error)
	GetAll(ctx context.Context, filter ListProductFilter) ([]Product, int64, error)
	Update(ctx context.Context, product UpdateProductRequest) error
	Delete(ctx context.Context, id int64) error
}
