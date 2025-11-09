package product

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/naxumi/bnsp-jwd/internal/pkg/database"
	"github.com/naxumi/bnsp-jwd/internal/service/file"
)

type ProductServiceImpl struct {
	db          *database.DB
	repository  productDomain.ProductRepository
	fileService file.FileService
}

func NewProductService(db *database.DB, repository productDomain.ProductRepository, fileService file.FileService) productDomain.ProductService {
	return &ProductServiceImpl{
		db:          db,
		repository:  repository,
		fileService: fileService,
	}
}

// UploadImage implements productDomain.ProductService.
func (s *ProductServiceImpl) UploadImage(ctx context.Context, id int64, file multipart.File, fileHeader *multipart.FileHeader) error {
	// Validate file is provided
	if fileHeader == nil {
		return productDomain.ErrImageRequired
	}

	// Validate file size (5MB max)
	const maxFileSize = 5 * 1024 * 1024 // 5MB in bytes
	if fileHeader.Size > maxFileSize {
		return productDomain.ErrImageTooLarge
	}

	// Validate file type
	ext := filepath.Ext(fileHeader.Filename)
	ext = strings.ToLower(ext)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif"}
	isValidExt := false
	for _, allowed := range allowedExts {
		if ext == allowed {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return productDomain.ErrInvalidImageFormat
	}

	// Get existing product to check for old image
	existingProduct, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Generate unique filename
	uniqueFilename := fmt.Sprintf("product-%d-image%s", id, ext)

	// Upload file using fileService
	uploadedPath, err := s.fileService.UploadProductImage(ctx, fmt.Sprintf("%d", id), file, uniqueFilename)
	if err != nil {
		return fmt.Errorf("failed to upload product image: %w", err)
	}

	// Get full URL for the uploaded file
	imageURL, err := s.fileService.GetFileURL(ctx, uploadedPath, 0)
	if err != nil {
		return fmt.Errorf("failed to get image URL: %w", err)
	}

	// Update product's image URL in the repository
	updateReq := productDomain.UpdateProductRequest{
		ID:       id,
		ImageURL: &imageURL,
	}
	if err := s.repository.Update(ctx, updateReq); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to update product image URL: %w", err)
	}

	// Delete old image if exists
	if existingProduct.ImageURL != nil && *existingProduct.ImageURL != "" {
		// Extract relative path from URL if it's a full URL
		oldImagePath := *existingProduct.ImageURL
		if len(oldImagePath) > 0 && (strings.HasPrefix(oldImagePath, "http://") || strings.HasPrefix(oldImagePath, "https://")) {
			// It's a full URL, extract the path after /uploads/
			parts := strings.Split(oldImagePath, "/uploads/")
			if len(parts) > 1 {
				oldImagePath = parts[1]
			}
		}
		if err := s.fileService.DeleteFile(ctx, oldImagePath); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to delete old image %s: %v\n", oldImagePath, err)
		}
	}

	return nil
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req productDomain.CreateProductRequest) (productDomain.ProductResponse, error) {

	newProduct := productDomain.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Status:      req.Status,
	}

	createdProduct, err := s.repository.Create(ctx, newProduct)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique violation
			return productDomain.ProductResponse{}, productDomain.ErrProductSKUExists
		}
		return productDomain.ProductResponse{}, fmt.Errorf("failed to create product: %w", err)
	}

	return productDomain.ProductResponse{
		ID:          createdProduct.ID,
		SKU:         createdProduct.SKU,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		Stock:       createdProduct.Stock,
		Category:    createdProduct.Category,
		Status:      createdProduct.Status,
		ImageURL: func() *string {
			if createdProduct.ImageURL == nil || *createdProduct.ImageURL == "" {
				return nil
			}
			return createdProduct.ImageURL
		}(),
		CreatedAt: createdProduct.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: createdProduct.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, id int64) (productDomain.ProductResponse, error) {
	p, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ProductResponse{}, productDomain.ErrProductNotFound
		}
		return productDomain.ProductResponse{}, fmt.Errorf("failed to get product: %w", err)
	}

	return productDomain.ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Status:      p.Status,
		ImageURL: func() *string {
			if p.ImageURL == nil || *p.ImageURL == "" {
				return nil
			}
			return p.ImageURL
		}(),
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *ProductServiceImpl) GetProductBySKU(ctx context.Context, sku string) (productDomain.ProductResponse, error) {
	p, err := s.repository.GetBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ProductResponse{}, productDomain.ErrProductNotFound
		}
		return productDomain.ProductResponse{}, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return productDomain.ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Status:      p.Status,
		ImageURL: func() *string {
			if p.ImageURL == nil || *p.ImageURL == "" {
				return nil
			}
			return p.ImageURL
		}(),
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req productDomain.UpdateProductRequest) error {

	if err := s.repository.Update(ctx, req); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique violation
			return productDomain.ErrProductSKUExists
		}
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	// Get product to check for image
	product, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Delete product from database
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Delete product image if exists
	if product.ImageURL != nil && *product.ImageURL != "" {
		imagePath := *product.ImageURL
		// Extract relative path from URL if it's a full URL
		if len(imagePath) > 0 && (strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://")) {
			// It's a full URL, extract the path after /uploads/
			parts := strings.Split(imagePath, "/uploads/")
			if len(parts) > 1 {
				imagePath = parts[1]
			}
		}
		if err := s.fileService.DeleteFile(ctx, imagePath); err != nil {
			// Log error but don't fail the operation since product is already deleted
			fmt.Printf("Warning: failed to delete product image %s: %v\n", imagePath, err)
		}
	}

	return nil
}

// DeleteImage implements productDomain.ProductService.
func (s *ProductServiceImpl) DeleteImage(ctx context.Context, id int64) error {
	// Get product to check if image exists
	product, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product.ImageURL == nil || *product.ImageURL == "" {
		return fmt.Errorf("product has no image to delete")
	}

	// Extract relative path from URL if it's a full URL
	imagePath := *product.ImageURL
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		// It's a full URL, extract the path after /uploads/
		parts := strings.Split(imagePath, "/uploads/")
		if len(parts) > 1 {
			imagePath = parts[1]
		}
	}

	// Delete the physical file
	if err := s.fileService.DeleteFile(ctx, imagePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	// Update product to remove image URL
	emptyString := ""
	updateReq := productDomain.UpdateProductRequest{
		ID:       id,
		ImageURL: &emptyString,
	}

	if err := s.repository.Update(ctx, updateReq); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ErrProductNotFound
		}
		return fmt.Errorf("failed to update product image URL: %w", err)
	}

	return nil
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, filter productDomain.ListProductFilter) (productDomain.ListProductResponse, error) {
	products, total, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productDomain.ListProductResponse{}, productDomain.ErrProductNotFound
		}
		return productDomain.ListProductResponse{}, fmt.Errorf("failed to list products: %w", err)
	}

	var productResponses []productDomain.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, productDomain.ProductResponse{
			ID:          p.ID,
			SKU:         p.SKU,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			Status:      p.Status,
			ImageURL: func() *string {
				if p.ImageURL == nil || *p.ImageURL == "" {
					return nil
				}
				return p.ImageURL
			}(),
			CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	startIdx := (filter.Page-1)*filter.Limit + 1
	endIdx := startIdx + len(productResponses) - 1
	showing := fmt.Sprintf("Showing %d to %d of %d products", startIdx, endIdx, total)

	return productDomain.ListProductResponse{
		TotalCount: total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(totalPages),
		Showing:    showing,
		Products:   productResponses,
	}, nil
}
