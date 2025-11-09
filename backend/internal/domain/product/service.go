package product

import (
	"context"
	"mime/multipart"
)

type ProductService interface {
	// Create a new product
	CreateProduct(ctx context.Context, req CreateProductRequest) (ProductResponse, error)

	// Retrieve product by ID or SKU
	GetProduct(ctx context.Context, id int64) (ProductResponse, error)
	GetProductBySKU(ctx context.Context, sku string) (ProductResponse, error)

	// Update and delete
	UpdateProduct(ctx context.Context, req UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int64) error

	// List products with pagination/filtering
	ListProducts(ctx context.Context, filter ListProductFilter) (ListProductResponse, error)

	// Upload and delete product image
	UploadImage(ctx context.Context, id int64, file multipart.File, fileHeader *multipart.FileHeader) error
	DeleteImage(ctx context.Context, id int64) error
}
