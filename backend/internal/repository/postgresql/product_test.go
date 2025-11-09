package postgresql

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/naxumi/bnsp-jwd/internal/pkg/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProductRepo(t *testing.T) (productDomain.ProductRepository, *database.DB, func()) {
	// Get DSN from environment or use default test database
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:root@localhost:5432/product_management?sslmode=disable"
	}

	db, err := database.NewPostgreSQLDB(dsn)
	require.NoError(t, err, "Failed to connect to test database")

	repo := NewProductRepository(db)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		_, _ = db.Exec(context.Background(), "DELETE FROM products WHERE sku LIKE 'TEST-%'")
		db.Close()
	}

	return repo, db, cleanup
}

func TestProductRepository_Create_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	desc := "Test Description"
	newProduct := productDomain.Product{
		SKU:         "TEST-SKU-001",
		Name:        "Test Product",
		Description: &desc,
		Price:       decimal.NewFromInt(10000),
		Stock:       100,
		Category:    "Electronics",
		Status:      productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	assert.NotZero(t, createdProduct.ID)
	assert.Equal(t, newProduct.SKU, createdProduct.SKU)
	assert.Equal(t, newProduct.Name, createdProduct.Name)
	assert.Equal(t, newProduct.Description, createdProduct.Description)
	assert.True(t, newProduct.Price.Equal(createdProduct.Price))
	assert.Equal(t, newProduct.Stock, createdProduct.Stock)
	assert.Equal(t, newProduct.Category, createdProduct.Category)
	assert.Equal(t, newProduct.Status, createdProduct.Status)
	assert.NotZero(t, createdProduct.CreatedAt)
	assert.NotZero(t, createdProduct.UpdatedAt)
}

func TestProductRepository_Create_DuplicateSKU(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	desc := "Test Description"
	newProduct := productDomain.Product{
		SKU:         "TEST-SKU-DUP",
		Name:        "Test Product",
		Description: &desc,
		Price:       decimal.NewFromInt(10000),
		Stock:       100,
		Category:    "Electronics",
		Status:      productDomain.ProductStatusActive,
	}

	// Create first product
	_, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Try to create duplicate
	_, err = repo.Create(context.Background(), newProduct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value")
}

func TestProductRepository_GetByID_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a product first
	desc := "Test Description 2"
	newProduct := productDomain.Product{
		SKU:         "TEST-SKU-002",
		Name:        "Test Product 2",
		Description: &desc,
		Price:       decimal.NewFromInt(20000),
		Stock:       50,
		Category:    "Books",
		Status:      productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Get the product by ID
	foundProduct, err := repo.GetByID(context.Background(), createdProduct.ID)
	require.NoError(t, err)

	assert.Equal(t, createdProduct.ID, foundProduct.ID)
	assert.Equal(t, createdProduct.SKU, foundProduct.SKU)
	assert.Equal(t, createdProduct.Name, foundProduct.Name)
	assert.True(t, createdProduct.Price.Equal(foundProduct.Price))
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	_, err := repo.GetByID(context.Background(), 999999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

func TestProductRepository_GetBySKU_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a product first
	desc := "Test Description 3"
	newProduct := productDomain.Product{
		SKU:         "TEST-SKU-003",
		Name:        "Test Product 3",
		Description: &desc,
		Price:       decimal.NewFromInt(30000),
		Stock:       75,
		Category:    "Clothing",
		Status:      productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Get the product by SKU
	foundProduct, err := repo.GetBySKU(context.Background(), "TEST-SKU-003")
	require.NoError(t, err)

	assert.Equal(t, createdProduct.ID, foundProduct.ID)
	assert.Equal(t, createdProduct.SKU, foundProduct.SKU)
	assert.Equal(t, createdProduct.Name, foundProduct.Name)
}

func TestProductRepository_GetBySKU_NotFound(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	_, err := repo.GetBySKU(context.Background(), "NONEXISTENT-SKU")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

func TestProductRepository_GetAll_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create multiple products
	products := []productDomain.Product{
		{
			SKU:      "TEST-SKU-LIST-1",
			Name:     "Product 1",
			Price:    decimal.NewFromInt(10000),
			Stock:    100,
			Category: "Electronics",
			Status:   productDomain.ProductStatusActive,
		},
		{
			SKU:      "TEST-SKU-LIST-2",
			Name:     "Product 2",
			Price:    decimal.NewFromInt(20000),
			Stock:    50,
			Category: "Books",
			Status:   productDomain.ProductStatusActive,
		},
	}

	for _, p := range products {
		_, err := repo.Create(context.Background(), p)
		require.NoError(t, err)
	}

	// Get all products
	filter := productDomain.ListProductFilter{
		Page:  1,
		Limit: 10,
	}

	foundProducts, total, err := repo.GetAll(context.Background(), filter)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(foundProducts), 2)
	assert.GreaterOrEqual(t, total, int64(2))
}

func TestProductRepository_GetAll_WithFilters(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a specific product
	newProduct := productDomain.Product{
		SKU:      "TEST-SKU-FILTER",
		Name:     "Filtered Product",
		Price:    decimal.NewFromInt(15000),
		Stock:    25,
		Category: "Toys",
		Status:   productDomain.ProductStatusActive,
	}

	_, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Filter by category
	category := "Toys"
	filter := productDomain.ListProductFilter{
		Page:     1,
		Limit:    10,
		Category: &category,
	}

	foundProducts, total, err := repo.GetAll(context.Background(), filter)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, total, int64(1))
	for _, p := range foundProducts {
		assert.Equal(t, "Toys", p.Category)
	}
}

func TestProductRepository_Update_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a product first
	newProduct := productDomain.Product{
		SKU:      "TEST-SKU-UPDATE",
		Name:     "Original Name",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Update the product
	time.Sleep(time.Millisecond * 100) // Ensure updated_at is different

	updatedName := "Updated Name"
	updatedPrice := decimal.NewFromInt(15000)
	updatedStock := 150
	updatedCategory := "Gadgets"
	updatedStatus := productDomain.ProductStatusInactive

	updateReq := productDomain.UpdateProductRequest{
		ID:       createdProduct.ID,
		Name:     &updatedName,
		Price:    &updatedPrice,
		Stock:    &updatedStock,
		Category: &updatedCategory,
		Status:   &updatedStatus,
	}

	err = repo.Update(context.Background(), updateReq)
	require.NoError(t, err)

	// Verify the update
	foundProduct, err := repo.GetByID(context.Background(), createdProduct.ID)
	require.NoError(t, err)

	assert.Equal(t, updatedName, foundProduct.Name)
	assert.True(t, updatedPrice.Equal(foundProduct.Price))
	assert.Equal(t, updatedStock, foundProduct.Stock)
	assert.Equal(t, updatedCategory, foundProduct.Category)
	assert.Equal(t, updatedStatus, foundProduct.Status)
	assert.NotEqual(t, createdProduct.UpdatedAt, foundProduct.UpdatedAt)
}

func TestProductRepository_Update_PartialUpdate(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a product first
	newProduct := productDomain.Product{
		SKU:      "TEST-SKU-PARTIAL",
		Name:     "Original Name",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Update only the name
	updatedName := "Partially Updated Name"
	updateReq := productDomain.UpdateProductRequest{
		ID:   createdProduct.ID,
		Name: &updatedName,
	}

	err = repo.Update(context.Background(), updateReq)
	require.NoError(t, err)

	// Verify only name changed
	foundProduct, err := repo.GetByID(context.Background(), createdProduct.ID)
	require.NoError(t, err)

	assert.Equal(t, updatedName, foundProduct.Name)
	assert.True(t, createdProduct.Price.Equal(foundProduct.Price))
	assert.Equal(t, createdProduct.Stock, foundProduct.Stock)
}

func TestProductRepository_Update_NotFound(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	updatedName := "Updated Name"
	updateReq := productDomain.UpdateProductRequest{
		ID:   999999,
		Name: &updatedName,
	}

	err := repo.Update(context.Background(), updateReq)
	assert.Error(t, err)
	// Could be either "no rows affected" or "no rows in result set"
	assert.True(t,
		strings.Contains(err.Error(), "no rows affected") ||
			strings.Contains(err.Error(), "no rows in result set"),
		"error should mention no rows: %v", err)
}

func TestProductRepository_Delete_Success(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	// Create a product first
	newProduct := productDomain.Product{
		SKU:      "TEST-SKU-DELETE",
		Name:     "Product to Delete",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	createdProduct, err := repo.Create(context.Background(), newProduct)
	require.NoError(t, err)

	// Delete the product
	err = repo.Delete(context.Background(), createdProduct.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.GetByID(context.Background(), createdProduct.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no rows")
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	repo, _, cleanup := setupProductRepo(t)
	defer cleanup()

	err := repo.Delete(context.Background(), 999999)
	assert.Error(t, err)
	// Could be either "no rows affected" or "no rows in result set"
	assert.True(t,
		strings.Contains(err.Error(), "no rows affected") ||
			strings.Contains(err.Error(), "no rows in result set"),
		"error should mention no rows: %v", err)
}
