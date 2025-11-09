package product

import (
	"github.com/naxumi/bnsp-jwd/internal/pkg/validator"
	"github.com/shopspring/decimal"
)

// ========================================
// PRODUCT DTOs
// ========================================

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	SKU         string          `json:"sku"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
	Category    string          `json:"category"`
	Status      ProductStatus   `json:"status"`
}

func (r *CreateProductRequest) Validate() error {
	var errs validator.ValidationErrors

	// SKU
	if validator.IsEmpty(r.SKU) {
		errs = append(errs, validator.ValidationError{
			Field:   "sku",
			Message: "sku is required",
		})
	}
	if len(r.SKU) > 100 {
		errs = append(errs, validator.ValidationError{
			Field:   "sku",
			Message: "sku must not exceed 100 characters",
		})
	}

	// Name
	if validator.IsEmpty(r.Name) {
		errs = append(errs, validator.ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	}

	// Price
	if r.Price.LessThan(decimal.Zero) {
		errs = append(errs, validator.ValidationError{
			Field:   "price",
			Message: "price must be greater than or equal to 0",
		})
	}
	// NUMERIC(10,2) allows max 99,999,999.99
	maxPrice := decimal.NewFromFloat(99999999.99)
	if r.Price.GreaterThan(maxPrice) {
		errs = append(errs, validator.ValidationError{
			Field:   "price",
			Message: "price must not exceed 99,999,999.99",
		})
	}

	// Stock
	if r.Stock < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "stock",
			Message: "stock must be greater than or equal to 0",
		})
	}

	// Category
	if validator.IsEmpty(r.Category) {
		errs = append(errs, validator.ValidationError{
			Field:   "category",
			Message: "category is required",
		})
	}
	if len(r.Category) > 100 {
		errs = append(errs, validator.ValidationError{
			Field:   "category",
			Message: "category must not exceed 100 characters",
		})
	}

	// Status
	if r.Status != ProductStatusActive && r.Status != ProductStatusInactive {
		errs = append(errs, validator.ValidationError{
			Field:   "status",
			Message: "status must be either 'Active' or 'Inactive'",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// UpdateProductRequest represents the request to update an existing product
type UpdateProductRequest struct {
	ID          int64            `json:"id"`
	SKU         *string          `json:"sku,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Price       *decimal.Decimal `json:"price,omitempty"`
	Stock       *int             `json:"stock,omitempty"`
	Category    *string          `json:"category,omitempty"`
	Status      *ProductStatus   `json:"status,omitempty"`
	ImageURL    *string          `json:"image_url,omitempty"`
}

func (r *UpdateProductRequest) Validate() error {
	var errs validator.ValidationErrors

	// ID
	if r.ID <= 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "id",
			Message: "id must be a positive integer",
		})
	}

	// SKU
	if r.SKU != nil {
		if validator.IsEmpty(*r.SKU) {
			errs = append(errs, validator.ValidationError{
				Field:   "sku",
				Message: "sku must not be empty",
			})
		}
		if len(*r.SKU) > 100 {
			errs = append(errs, validator.ValidationError{
				Field:   "sku",
				Message: "sku must not exceed 100 characters",
			})
		}
	}

	// Name
	if r.Name != nil {
		if validator.IsEmpty(*r.Name) {
			errs = append(errs, validator.ValidationError{
				Field:   "name",
				Message: "name must not be empty",
			})
		}
	}

	// Price
	if r.Price != nil && r.Price.LessThan(decimal.Zero) {
		errs = append(errs, validator.ValidationError{
			Field:   "price",
			Message: "price must be greater than or equal to 0",
		})
	}
	// NUMERIC(10,2) allows max 99,999,999.99
	if r.Price != nil {
		maxPrice := decimal.NewFromFloat(99999999.99)
		if r.Price.GreaterThan(maxPrice) {
			errs = append(errs, validator.ValidationError{
				Field:   "price",
				Message: "price must not exceed 99,999,999.99",
			})
		}
	}

	// Stock
	if r.Stock != nil && *r.Stock < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "stock",
			Message: "stock must be greater than or equal to 0",
		})
	}

	// Category
	if r.Category != nil {
		if validator.IsEmpty(*r.Category) {
			errs = append(errs, validator.ValidationError{
				Field:   "category",
				Message: "category must not be empty",
			})
		}
		if len(*r.Category) > 100 {
			errs = append(errs, validator.ValidationError{
				Field:   "category",
				Message: "category must not exceed 100 characters",
			})
		}
	}

	// Status
	if r.Status != nil && *r.Status != ProductStatusActive && *r.Status != ProductStatusInactive {
		errs = append(errs, validator.ValidationError{
			Field:   "status",
			Message: "status must be either 'Active' or 'Inactive'",
		})
	}

	// ImageURL
	if r.ImageURL != nil && len(*r.ImageURL) > 2048 {
		errs = append(errs, validator.ValidationError{
			Field:   "image_url",
			Message: "image_url must not exceed 2048 characters",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ProductResponse represents the response containing product data
type ProductResponse struct {
	ID          int64           `json:"id"`
	SKU         string          `json:"sku"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
	Category    string          `json:"category"`
	Status      ProductStatus   `json:"status"`
	ImageURL    *string         `json:"image_url,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// ListProductFilter represents the filter for listing products
type ListProductFilter struct {
	// Search & Filter
	Name     *string        `json:"name,omitempty"`
	SKU      *string        `json:"sku,omitempty"`
	Category *string        `json:"category,omitempty"`
	Status   *ProductStatus `json:"status,omitempty"`
	MinPrice *float64       `json:"min_price,omitempty"`
	MaxPrice *float64       `json:"max_price,omitempty"`

	// Pagination
	Page  int `json:"page"`
	Limit int `json:"limit"`

	// Sorting
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

func (f *ListProductFilter) Validate() error {
	var errs validator.ValidationErrors

	// Page validation
	if f.Page < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "page",
			Message: "page must be a positive number",
		})
	}
	if f.Page == 0 {
		f.Page = 1 // Default page
	}

	// Limit validation
	if f.Limit < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "limit",
			Message: "limit must be a positive number",
		})
	}
	if f.Limit == 0 {
		f.Limit = 20 // Default limit
	}
	if f.Limit > 100 {
		errs = append(errs, validator.ValidationError{
			Field:   "limit",
			Message: "limit must not exceed 100",
		})
	}

	// Price validation
	if f.MinPrice != nil && *f.MinPrice < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "min_price",
			Message: "min_price must be greater than or equal to 0",
		})
	}
	if f.MaxPrice != nil && *f.MaxPrice < 0 {
		errs = append(errs, validator.ValidationError{
			Field:   "max_price",
			Message: "max_price must be greater than or equal to 0",
		})
	}
	if f.MinPrice != nil && f.MaxPrice != nil && *f.MinPrice > *f.MaxPrice {
		errs = append(errs, validator.ValidationError{
			Field:   "price",
			Message: "min_price must be less than or equal to max_price",
		})
	}

	// Sort validation
	if f.SortBy != "" {
		validSortFields := []string{"id", "sku", "name", "price", "stock", "category", "status", "created_at", "updated_at"}
		if !validator.IsInSlice(f.SortBy, validSortFields) {
			errs = append(errs, validator.ValidationError{
				Field:   "sort_by",
				Message: "sort_by must be one of: id, sku, name, price, stock, category, status, created_at, updated_at",
			})
		}
	} else {
		f.SortBy = "created_at" // Default sort
	}

	if f.SortOrder != "" {
		validSortOrders := []string{"asc", "desc"}
		if !validator.IsInSlice(f.SortOrder, validSortOrders) {
			errs = append(errs, validator.ValidationError{
				Field:   "sort_order",
				Message: "sort_order must be one of: asc, desc",
			})
		}
	} else {
		f.SortOrder = "desc" // Default descending
	}

	// Status validation
	if f.Status != nil && *f.Status != ProductStatusActive && *f.Status != ProductStatusInactive {
		errs = append(errs, validator.ValidationError{
			Field:   "status",
			Message: "status must be either 'Active' or 'Inactive'",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ListProductResponse represents the paginated response for listing products
type ListProductResponse struct {
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
	Showing    string            `json:"showing"`
	Products   []ProductResponse `json:"products"`
}
