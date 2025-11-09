package product

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock File implements multipart.File interface
type MockFile struct {
	*strings.Reader
}

func (m *MockFile) Close() error {
	return nil
}

func (m *MockFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func NewMockFile(content string) *MockFile {
	return &MockFile{Reader: strings.NewReader(content)}
}

// Mock Repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product productDomain.Product) (productDomain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(productDomain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (productDomain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(productDomain.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (productDomain.Product, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(productDomain.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter productDomain.ListProductFilter) ([]productDomain.Product, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]productDomain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Update(ctx context.Context, product productDomain.UpdateProductRequest) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock File Service
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadProductImage(ctx context.Context, productID string, file io.Reader, filename string) (string, error) {
	args := m.Called(ctx, productID, file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) UploadAvatar(ctx context.Context, employeeID string, file io.Reader, filename string) (string, error) {
	args := m.Called(ctx, employeeID, file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) UploadDocument(ctx context.Context, employeeID string, file io.Reader, filename string, documentType string) (string, error) {
	args := m.Called(ctx, employeeID, file, filename, documentType)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) UploadAttendanceProof(ctx context.Context, employeeID string, date time.Time, file io.Reader, filename string, clockType string) (string, error) {
	args := m.Called(ctx, employeeID, date, file, filename, clockType)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) UploadLeaveAttachment(ctx context.Context, employeeID string, file io.Reader, filename string) (string, error) {
	args := m.Called(ctx, employeeID, file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) UploadCompanyLogo(ctx context.Context, companyUsername string, file io.Reader, filename string) (string, error) {
	args := m.Called(ctx, companyUsername, file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) DeleteFile(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockFileService) GetFileURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, path, expiry)
	return args.String(0), args.Error(1)
}

// Tests for CreateProduct
func TestProductService_CreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	req := productDomain.CreateProductRequest{
		SKU:      "TEST-SKU-001",
		Name:     "Test Product",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	now := time.Now()
	expectedProduct := productDomain.Product{
		ID:        1,
		SKU:       req.SKU,
		Name:      req.Name,
		Price:     req.Price,
		Stock:     req.Stock,
		Category:  req.Category,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p productDomain.Product) bool {
		return p.SKU == req.SKU && p.Name == req.Name && p.Price.Equal(req.Price)
	})).Return(expectedProduct, nil)

	result, err := service.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, result.ID)
	assert.Equal(t, expectedProduct.SKU, result.SKU)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.True(t, expectedProduct.Price.Equal(result.Price))
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	req := productDomain.CreateProductRequest{
		SKU:      "TEST-SKU-001",
		Name:     "Test Product",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	// Simulate repository error
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(productDomain.Product{}, errors.New("database error"))

	_, err := service.CreateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_DuplicateSKU(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	req := productDomain.CreateProductRequest{
		SKU:      "DUPLICATE-SKU",
		Name:     "Test Product",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	pgErr := &pgconn.PgError{Code: "23505"}
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(productDomain.Product{}, pgErr)

	_, err := service.CreateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductSKUExists, err)
	mockRepo.AssertExpectations(t)
}

// Tests for GetProduct
func TestProductService_GetProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	now := time.Now()
	expectedProduct := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).
		Return(expectedProduct, nil)

	result, err := service.GetProduct(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, result.ID)
	assert.Equal(t, expectedProduct.SKU, result.SKU)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	mockRepo.On("GetByID", mock.Anything, int64(999)).
		Return(productDomain.Product{}, errors.New("product not found: no rows in result set"))

	_, err := service.GetProduct(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get product")
	mockRepo.AssertExpectations(t)
}

// Tests for GetProductBySKU
func TestProductService_GetProductBySKU_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	now := time.Now()
	expectedProduct := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetBySKU", mock.Anything, "TEST-SKU-001").
		Return(expectedProduct, nil)

	result, err := service.GetProductBySKU(context.Background(), "TEST-SKU-001")

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.SKU, result.SKU)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductBySKU_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	mockRepo.On("GetBySKU", mock.Anything, "NONEXISTENT").
		Return(productDomain.Product{}, pgx.ErrNoRows)

	_, err := service.GetProductBySKU(context.Background(), "NONEXISTENT")

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

// Tests for UpdateProduct
func TestProductService_UpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	name := "Updated Product"
	price := decimal.NewFromInt(15000)

	req := productDomain.UpdateProductRequest{
		ID:    1,
		Name:  &name,
		Price: &price,
	}

	mockRepo.On("Update", mock.Anything, req).
		Return(nil)

	err := service.UpdateProduct(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	name := "Updated Product"
	req := productDomain.UpdateProductRequest{
		ID:   1,
		Name: &name,
	}

	// Simulate repository error
	mockRepo.On("Update", mock.Anything, req).
		Return(errors.New("database error"))

	err := service.UpdateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	name := "Updated Product"
	req := productDomain.UpdateProductRequest{
		ID:   999,
		Name: &name,
	}

	mockRepo.On("Update", mock.Anything, req).
		Return(pgx.ErrNoRows)

	err := service.UpdateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_DuplicateSKU(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	sku := "DUPLICATE-SKU"
	req := productDomain.UpdateProductRequest{
		ID:  1,
		SKU: &sku,
	}

	pgErr := &pgconn.PgError{Code: "23505"}
	mockRepo.On("Update", mock.Anything, req).
		Return(pgErr)

	err := service.UpdateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductSKUExists, err)
	mockRepo.AssertExpectations(t)
}

// Tests for DeleteProduct
func TestProductService_DeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	now := time.Now()
	product := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		ImageURL:  nil, // No image
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).
		Return(product, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).
		Return(nil)

	err := service.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_WithImage(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	imageURL := "http://localhost:8080/uploads/products/1/image.jpg"
	now := time.Now()
	product := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		ImageURL:  &imageURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).
		Return(product, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).
		Return(nil)
	mockFileService.On("DeleteFile", mock.Anything, mock.Anything).
		Return(nil)

	err := service.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockFileService.AssertExpectations(t)
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	mockRepo.On("GetByID", mock.Anything, int64(999)).
		Return(productDomain.Product{}, pgx.ErrNoRows)

	err := service.DeleteProduct(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

// Tests for ListProducts
func TestProductService_ListProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	filter := productDomain.ListProductFilter{
		Page:      1,
		Limit:     10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	now := time.Now()
	products := []productDomain.Product{
		{
			ID:        1,
			SKU:       "SKU-001",
			Name:      "Product 1",
			Price:     decimal.NewFromInt(10000),
			Stock:     100,
			Category:  "Electronics",
			Status:    productDomain.ProductStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			SKU:       "SKU-002",
			Name:      "Product 2",
			Price:     decimal.NewFromInt(20000),
			Stock:     50,
			Category:  "Electronics",
			Status:    productDomain.ProductStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(f productDomain.ListProductFilter) bool {
		return f.Page == 1 && f.Limit == 10
	})).Return(products, int64(2), nil)

	result, err := service.ListProducts(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.TotalCount)
	assert.Len(t, result.Products, 2)
	assert.Equal(t, "SKU-001", result.Products[0].SKU)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ListProducts_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	filter := productDomain.ListProductFilter{
		Page:      1,
		Limit:     10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	// Simulate repository error
	mockRepo.On("GetAll", mock.Anything, mock.Anything).
		Return([]productDomain.Product{}, int64(0), errors.New("database error"))

	_, err := service.ListProducts(context.Background(), filter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list products")
	mockRepo.AssertExpectations(t)
}

// Tests for UploadImage
func TestProductService_UploadImage_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	// Create a mock file
	file := NewMockFile("fake image content")
	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     100,
	}

	now := time.Now()
	existingProduct := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		ImageURL:  nil, // No existing image
		CreatedAt: now,
		UpdatedAt: now,
	}

	uploadedPath := "products/1/product-1-image.jpg"
	fullURL := "http://localhost:8080/uploads/products/1/product-1-image.jpg"

	mockRepo.On("GetByID", mock.Anything, int64(1)).
		Return(existingProduct, nil)

	mockFileService.On("UploadProductImage", mock.Anything, "1", mock.Anything, "product-1-image.jpg").
		Return(uploadedPath, nil)

	mockFileService.On("GetFileURL", mock.Anything, uploadedPath, time.Duration(0)).
		Return(fullURL, nil)

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(req productDomain.UpdateProductRequest) bool {
		return req.ID == 1 && req.ImageURL != nil && *req.ImageURL == fullURL
	})).Return(nil)

	err := service.UploadImage(context.Background(), 1, file, fileHeader)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockFileService.AssertExpectations(t)
}

func TestProductService_UploadImage_ReplaceExisting(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	file := NewMockFile("fake image content")
	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     100,
	}

	oldImageURL := "http://localhost:8080/uploads/products/1/old-image.jpg"
	now := time.Now()
	existingProduct := productDomain.Product{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		ImageURL:  &oldImageURL, // Has existing image
		CreatedAt: now,
		UpdatedAt: now,
	}

	uploadedPath := "products/1/product-1-image.jpg"
	fullURL := "http://localhost:8080/uploads/products/1/product-1-image.jpg"

	mockRepo.On("GetByID", mock.Anything, int64(1)).
		Return(existingProduct, nil)

	mockFileService.On("UploadProductImage", mock.Anything, "1", mock.Anything, "product-1-image.jpg").
		Return(uploadedPath, nil)

	mockFileService.On("GetFileURL", mock.Anything, uploadedPath, time.Duration(0)).
		Return(fullURL, nil)

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(req productDomain.UpdateProductRequest) bool {
		return req.ID == 1 && req.ImageURL != nil && *req.ImageURL == fullURL
	})).Return(nil)

	mockFileService.On("DeleteFile", mock.Anything, "products/1/old-image.jpg").
		Return(nil)

	err := service.UploadImage(context.Background(), 1, file, fileHeader)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockFileService.AssertExpectations(t)
}

func TestProductService_UploadImage_InvalidFileType(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	file := NewMockFile("fake file content")
	fileHeader := &multipart.FileHeader{
		Filename: "test.pdf", // Invalid file type
		Size:     100,
	}

	err := service.UploadImage(context.Background(), 1, file, fileHeader)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrInvalidImageFormat, err)
	mockFileService.AssertNotCalled(t, "UploadProductImage")
	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestProductService_UploadImage_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockFileService := new(MockFileService)

	service := &ProductServiceImpl{
		repository:  mockRepo,
		fileService: mockFileService,
	}

	file := NewMockFile("fake image content")
	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     100,
	}

	mockRepo.On("GetByID", mock.Anything, int64(999)).
		Return(productDomain.Product{}, pgx.ErrNoRows)

	err := service.UploadImage(context.Background(), 999, file, fileHeader)

	assert.Error(t, err)
	assert.Equal(t, productDomain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
	mockFileService.AssertNotCalled(t, "UploadProductImage")
}
