package response

import (
	"errors"
	"net/http"

	productDomain "github.com/naxumi/bnsp-jwd/internal/domain/product"
	"github.com/naxumi/bnsp-jwd/internal/pkg/validator"
)

// HandleError maps domain errors to HTTP responses
func HandleError(w http.ResponseWriter, err error) {
	// Check if it's a validation error
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		ValidationError(w, validationErrs.ToMap())
		return
	}

	switch {
	// Product domain errors
	case errors.Is(err, productDomain.ErrProductNotFound):
		NotFound(w, "Product not found")
	case errors.Is(err, productDomain.ErrProductSKUExists):
		Conflict(w, "Product with this SKU already exists")
	case errors.Is(err, productDomain.ErrInvalidProductStatus):
		BadRequest(w, "Invalid product status", nil)
	case errors.Is(err, productDomain.ErrProductIDRequired):
		BadRequest(w, "Product ID is required", nil)
	case errors.Is(err, productDomain.ErrInvalidPrice):
		BadRequest(w, "Invalid price", nil)
	case errors.Is(err, productDomain.ErrInvalidStock):
		BadRequest(w, "Invalid stock", nil)
	case errors.Is(err, productDomain.ErrInvalidImageFormat):
		BadRequest(w, "Invalid image format, only JPG, JPEG, PNG, GIF are allowed", nil)
	case errors.Is(err, productDomain.ErrImageTooLarge):
		BadRequest(w, "Image file size exceeds maximum limit of 5MB", nil)
	case errors.Is(err, productDomain.ErrImageRequired):
		BadRequest(w, "Image file is required", nil)

	// Default
	default:
		InternalServerError(w, "An unexpected error occurred")
	}
}
