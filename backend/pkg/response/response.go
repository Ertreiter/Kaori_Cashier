package response

import (
	"github.com/gin-gonic/gin"
)

// APIResponse is a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains pagination info
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with pagination meta
func SuccessWithMeta(c *gin.Context, statusCode int, data interface{}, meta *Meta) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// Common error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodePaymentFailed  = "PAYMENT_FAILED"
	ErrCodeVoucherInvalid = "VOUCHER_INVALID"
)

// BadRequest sends a 400 error
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, ErrCodeBadRequest, message)
}

// Unauthorized sends a 401 error
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, ErrCodeUnauthorized, message)
}

// Forbidden sends a 403 error
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, ErrCodeForbidden, message)
}

// NotFound sends a 404 error
func NotFound(c *gin.Context, message string) {
	Error(c, 404, ErrCodeNotFound, message)
}

// Conflict sends a 409 error
func Conflict(c *gin.Context, message string) {
	Error(c, 409, ErrCodeConflict, message)
}

// InternalError sends a 500 error
func InternalError(c *gin.Context, message string) {
	Error(c, 500, ErrCodeInternal, message)
}

// ValidationError sends a 422 error
func ValidationError(c *gin.Context, message string) {
	Error(c, 422, ErrCodeValidation, message)
}

// CalculateTotalPages calculates total pages for pagination
func CalculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		pages++
	}
	return pages
}
