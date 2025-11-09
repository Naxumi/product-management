package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/naxumi/bnsp-jwd/internal/handler/http/response"
)

type ProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetProductBySKU(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
	DeleteImage(w http.ResponseWriter, r *http.Request)
}

type ProductHandlerImpl struct {
	productService productDomain.ProductService
}

func NewProductHandler(productService productDomain.ProductService) ProductHandler {
	return &ProductHandlerImpl{
		productService: productService,
	}
}

// UploadImage implements ProductHandler.
func (h *ProductHandlerImpl) UploadImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	// Parse multipart form (max 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		response.BadRequest(w, "Failed to parse form", nil)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
		response.BadRequest(w, "Image file is required", nil)
		return
	}
	defer file.Close()

	err = h.productService.UploadImage(r.Context(), id, file, fileHeader)
	if err != nil {
		log.Printf("Error uploading image for product ID %d: %v", id, err)
		response.HandleError(w, err)
		return
	}

	response.SuccessWithMessage(w, "Image uploaded successfully", nil)
}

// DeleteImage implements ProductHandler.
func (h *ProductHandlerImpl) DeleteImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	err = h.productService.DeleteImage(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting image for product ID %d: %v", id, err)
		response.HandleError(w, err)
		return
	}

	response.SuccessWithMessage(w, "Image deleted successfully", nil)
}

func (h *ProductHandlerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productDomain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request format: %v", err)
		response.BadRequest(w, "Invalid request format", nil)
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		response.HandleError(w, err)
		return
	}

	createdProduct, err := h.productService.CreateProduct(r.Context(), req)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		response.HandleError(w, err)
		return
	}

	response.Created(w, "Product created successfully", createdProduct)
}

func (h *ProductHandlerImpl) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		log.Printf("Error getting product with ID %d: %v", id, err)
		response.HandleError(w, err)
		return
	}

	response.Success(w, product)
}

func (h *ProductHandlerImpl) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")
	if sku == "" {
		response.BadRequest(w, "SKU is required", nil)
		return
	}

	product, err := h.productService.GetProductBySKU(r.Context(), sku)
	if err != nil {
		log.Printf("Error getting product with SKU %s: %v", sku, err)
		response.HandleError(w, err)
		return
	}

	response.Success(w, product)
}

func (h *ProductHandlerImpl) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req productDomain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request format: %v", err)
		response.BadRequest(w, "Invalid request format", nil)
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		response.HandleError(w, err)
		return
	}

	err := h.productService.UpdateProduct(r.Context(), req)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		response.HandleError(w, err)
		return
	}

	response.SuccessWithMessage(w, "Product updated successfully", nil)
}

func (h *ProductHandlerImpl) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	err = h.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting product with ID %d: %v", id, err)
		response.HandleError(w, err)
		return
	}

	response.SuccessWithMessage(w, "Product deleted successfully", nil)
}

func (h *ProductHandlerImpl) ListProducts(w http.ResponseWriter, r *http.Request) {
	var filter productDomain.ListProductFilter

	// Parse query parameters
	queryParams := r.URL.Query()

	// Search & Filter
	if name := queryParams.Get("name"); name != "" {
		filter.Name = &name
	}
	if sku := queryParams.Get("sku"); sku != "" {
		filter.SKU = &sku
	}
	if category := queryParams.Get("category"); category != "" {
		filter.Category = &category
	}
	if status := queryParams.Get("status"); status != "" {
		productStatus := productDomain.ProductStatus(status)
		filter.Status = &productStatus
	}

	// Price range
	if minPrice := queryParams.Get("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &price
		}
	}
	if maxPrice := queryParams.Get("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &price
		}
	}

	// Pagination
	if page := queryParams.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}
	if limit := queryParams.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	// Sorting
	if sortBy := queryParams.Get("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder := queryParams.Get("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	// Validate filter
	if err := filter.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		response.HandleError(w, err)
		return
	}

	products, err := h.productService.ListProducts(r.Context(), filter)
	if err != nil {
		log.Printf("Error listing products: %v", err)
		response.HandleError(w, err)
		return
	}

	response.Success(w, products)
}
