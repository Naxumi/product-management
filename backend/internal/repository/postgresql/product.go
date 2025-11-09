package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/naxumi/bnsp-jwd/internal/pkg/database"
	"github.com/shopspring/decimal"
)

type productRepositoryImpl struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) productDomain.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(ctx context.Context, newProduct productDomain.Product) (productDomain.Product, error) {
	q := GetQuerier(ctx, r.db)

	query := `
		INSERT INTO products (sku, name, description, price, stock, category, status, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := q.QueryRow(ctx, query,
		newProduct.SKU,
		newProduct.Name,
		newProduct.Description,
		newProduct.Price,
		newProduct.Stock,
		newProduct.Category,
		newProduct.Status,
		newProduct.ImageURL,
	).Scan(&newProduct.ID, &newProduct.CreatedAt, &newProduct.UpdatedAt)
	if err != nil {
		return productDomain.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	return newProduct, nil
}

func (r *productRepositoryImpl) GetByID(ctx context.Context, id int64) (productDomain.Product, error) {
	q := GetQuerier(ctx, r.db)

	query := `
		SELECT id, sku, name, description, price, stock, category, status, image_url, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product productDomain.Product
	err := q.QueryRow(ctx, query, id).
		Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.Status,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return productDomain.Product{}, fmt.Errorf("product not found: %w", err)
		}
		return productDomain.Product{}, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return product, nil
}

func (r *productRepositoryImpl) GetBySKU(ctx context.Context, sku string) (productDomain.Product, error) {
	q := GetQuerier(ctx, r.db)

	query := `
		SELECT id, sku, name, description, price, stock, category, status, image_url, created_at, updated_at
		FROM products
		WHERE sku = $1
	`

	var product productDomain.Product
	err := q.QueryRow(ctx, query, sku).
		Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.Status,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return productDomain.Product{}, fmt.Errorf("product not found: %w", err)
		}
		return productDomain.Product{}, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

func (r *productRepositoryImpl) GetAll(ctx context.Context, filter productDomain.ListProductFilter) ([]productDomain.Product, int64, error) {
	q := GetQuerier(ctx, r.db)

	whereClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if filter.Name != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Name+"%")
		argIdx++
	}

	if filter.SKU != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("sku ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.SKU+"%")
		argIdx++
	}

	if filter.Category != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("category ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Category+"%")
		argIdx++
	}

	if filter.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argIdx))
		args = append(args, decimal.NewFromFloat(*filter.MinPrice))
		argIdx++
	}

	if filter.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, decimal.NewFromFloat(*filter.MaxPrice))
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereSQL)
	var total int64
	err := q.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Default sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}

	query := fmt.Sprintf(`
		SELECT id, sku, name, description, price, stock, category, status, image_url, created_at, updated_at
		FROM products
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, sortBy, sortOrder, argIdx, argIdx+1)

	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []productDomain.Product
	for rows.Next() {
		var product productDomain.Product
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.Status,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, product productDomain.UpdateProductRequest) error {
	q := GetQuerier(ctx, r.db)

	updates := []string{}
	args := []interface{}{}
	argIdx := 1

	if product.SKU != nil {
		updates = append(updates, fmt.Sprintf("sku = $%d", argIdx))
		args = append(args, *product.SKU)
		argIdx++
	}

	if product.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *product.Name)
		argIdx++
	}

	if product.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *product.Description)
		argIdx++
	}

	if product.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, *product.Price)
		argIdx++
	}

	if product.Stock != nil {
		updates = append(updates, fmt.Sprintf("stock = $%d", argIdx))
		args = append(args, *product.Stock)
		argIdx++
	}

	if product.Category != nil {
		updates = append(updates, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, *product.Category)
		argIdx++
	}

	if product.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *product.Status)
		argIdx++
	}

	if product.ImageURL != nil {
		updates = append(updates, fmt.Sprintf("image_url = $%d", argIdx))
		args = append(args, *product.ImageURL)
		argIdx++
	}

	if len(updates) == 0 {
		// No fields to update, just return success
		return nil
	}

	args = append(args, product.ID)
	query := fmt.Sprintf(`
		UPDATE products
		SET %s, updated_at = NOW()
		WHERE id = $%d
	`, strings.Join(updates, ", "), argIdx)

	commandTag, err := q.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id int64) error {
	q := GetQuerier(ctx, r.db)

	query := `
		DELETE FROM products
		WHERE id = $1
	`

	commandTag, err := q.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
