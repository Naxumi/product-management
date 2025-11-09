package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Product Service
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req productDomain.CreateProductRequest) (productDomain.ProductResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(productDomain.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProduct(ctx context.Context, id int64) (productDomain.ProductResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(productDomain.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProductBySKU(ctx context.Context, sku string) (productDomain.ProductResponse, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(productDomain.ProductResponse), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, req productDomain.UpdateProductRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) ListProducts(ctx context.Context, filter productDomain.ListProductFilter) (productDomain.ListProductResponse, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(productDomain.ListProductResponse), args.Error(1)
}

func (m *MockProductService) UploadImage(ctx context.Context, id int64, file multipart.File, fileHeader *multipart.FileHeader) error {
	args := m.Called(ctx, id, file, fileHeader)
	return args.Error(0)
}

func (m *MockProductService) DeleteImage(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Tests for CreateProduct Handler
func TestProductHandler_CreateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	reqBody := productDomain.CreateProductRequest{
		SKU:      "TEST-SKU-001",
		Name:     "Test Product",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	expectedResp := productDomain.ProductResponse{
		ID:        1,
		SKU:       reqBody.SKU,
		Name:      reqBody.Name,
		Price:     reqBody.Price,
		Stock:     reqBody.Stock,
		Category:  reqBody.Category,
		Status:    reqBody.Status,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	mockService.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r productDomain.CreateProductRequest) bool {
		return r.SKU == reqBody.SKU && r.Name == reqBody.Name
	})).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/product", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, false, response["success"])
}

func TestProductHandler_CreateProduct_ValidationError(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	// Missing required fields
	reqBody := productDomain.CreateProductRequest{
		Name: "Test Product",
		// Missing SKU, Price, Stock, Category, Status
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, false, response["success"])
	mockService.AssertNotCalled(t, "CreateProduct")
}

func TestProductHandler_CreateProduct_DuplicateSKU(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	reqBody := productDomain.CreateProductRequest{
		SKU:      "DUPLICATE-SKU",
		Name:     "Test Product",
		Price:    decimal.NewFromInt(10000),
		Stock:    100,
		Category: "Electronics",
		Status:   productDomain.ProductStatusActive,
	}

	mockService.On("CreateProduct", mock.Anything, mock.Anything).
		Return(productDomain.ProductResponse{}, productDomain.ErrProductSKUExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, false, response["success"])
	mockService.AssertExpectations(t)
}

// Tests for GetProduct Handler
func TestProductHandler_GetProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	expectedResp := productDomain.ProductResponse{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	mockService.On("GetProduct", mock.Anything, int64(1)).
		Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_InvalidID(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, false, response["success"])
	mockService.AssertNotCalled(t, "GetProduct")
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("GetProduct", mock.Anything, int64(999)).
		Return(productDomain.ProductResponse{}, productDomain.ErrProductNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, false, response["success"])
	mockService.AssertExpectations(t)
}

// Tests for GetProductBySKU Handler
func TestProductHandler_GetProductBySKU_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	expectedResp := productDomain.ProductResponse{
		ID:        1,
		SKU:       "TEST-SKU-001",
		Name:      "Test Product",
		Price:     decimal.NewFromInt(10000),
		Stock:     100,
		Category:  "Electronics",
		Status:    productDomain.ProductStatusActive,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	mockService.On("GetProductBySKU", mock.Anything, "TEST-SKU-001").
		Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/sku/TEST-SKU-001", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("sku", "TEST-SKU-001")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetProductBySKU(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProductBySKU_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("GetProductBySKU", mock.Anything, "NONEXISTENT").
		Return(productDomain.ProductResponse{}, productDomain.ErrProductNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/sku/NONEXISTENT", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("sku", "NONEXISTENT")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetProductBySKU(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for UpdateProduct Handler
func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	name := "Updated Product"
	price := decimal.NewFromInt(15000)

	reqBody := productDomain.UpdateProductRequest{
		ID:    1,
		Name:  &name,
		Price: &price,
	}

	mockService.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(r productDomain.UpdateProductRequest) bool {
		return r.ID == 1 && r.Name != nil && *r.Name == name
	})).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_UpdateProduct_ValidationError(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	// Invalid ID
	reqBody := productDomain.UpdateProductRequest{
		ID: -1,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	mockService.AssertNotCalled(t, "UpdateProduct")
}

func TestProductHandler_UpdateProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	name := "Updated Product"
	reqBody := productDomain.UpdateProductRequest{
		ID:   999,
		Name: &name,
	}

	mockService.On("UpdateProduct", mock.Anything, mock.Anything).
		Return(productDomain.ErrProductNotFound)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for DeleteProduct Handler
func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("DeleteProduct", mock.Anything, int64(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_InvalidID(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "DeleteProduct")
}

func TestProductHandler_DeleteProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("DeleteProduct", mock.Anything, int64(999)).
		Return(productDomain.ErrProductNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for ListProducts Handler
func TestProductHandler_ListProducts_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	expectedResp := productDomain.ListProductResponse{
		TotalCount: 2,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Products: []productDomain.ProductResponse{
			{
				ID:       1,
				SKU:      "SKU-001",
				Name:     "Product 1",
				Price:    decimal.NewFromInt(10000),
				Stock:    100,
				Category: "Electronics",
				Status:   productDomain.ProductStatusActive,
			},
			{
				ID:       2,
				SKU:      "SKU-002",
				Name:     "Product 2",
				Price:    decimal.NewFromInt(20000),
				Stock:    50,
				Category: "Electronics",
				Status:   productDomain.ProductStatusActive,
			},
		},
	}

	mockService.On("ListProducts", mock.Anything, mock.Anything).
		Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_ListProducts_WithFilters(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	expectedResp := productDomain.ListProductResponse{
		TotalCount: 1,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Products: []productDomain.ProductResponse{
			{
				ID:       1,
				SKU:      "SKU-001",
				Name:     "Test Product",
				Price:    decimal.NewFromInt(10000),
				Stock:    100,
				Category: "Electronics",
				Status:   productDomain.ProductStatusActive,
			},
		},
	}

	mockService.On("ListProducts", mock.Anything, mock.Anything).
		Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/product?name=Test&category=Electronics&status=Active", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	mockService.AssertExpectations(t)
}

// Tests for UploadImage Handler
func TestProductHandler_UploadImage_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("UploadImage", mock.Anything, int64(1), mock.Anything, mock.Anything).
		Return(nil)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	io.WriteString(part, "fake image content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/product/1/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UploadImage(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response["success"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_UploadImage_InvalidID(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/product/invalid/image", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UploadImage(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "UploadImage")
}

func TestProductHandler_UploadImage_NoFile(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/product/1/image", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UploadImage(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "UploadImage")
}

func TestProductHandler_UploadImage_ProductNotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := &ProductHandlerImpl{productService: mockService}

	mockService.On("UploadImage", mock.Anything, int64(999), mock.Anything, mock.Anything).
		Return(productDomain.ErrProductNotFound)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	io.WriteString(part, "fake image content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/product/999/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UploadImage(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
