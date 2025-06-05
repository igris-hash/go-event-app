package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ResponseWithPagination is the standard API response structure with pagination
type ResponseWithPagination struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination information
type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	PageSize     int `json:"page_size"`
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendSuccessWithPagination sends a success response with pagination
func SendSuccessWithPagination(c *gin.Context, message string, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, ResponseWithPagination{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// SendCreated sends a success response with 201 status code
func SendCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendBadRequest sends a 400 bad request response
func SendBadRequest(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

// SendUnauthorized sends a 401 unauthorized response
func SendUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

// SendForbidden sends a 403 forbidden response
func SendForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

// SendNotFound sends a 404 not found response
func SendNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

// SendInternalError sends a 500 internal server error response
func SendInternalError(c *gin.Context, message string, err error) {
	errorMessage := "Internal server error"
	if err != nil {
		errorMessage = err.Error()
	}
	
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error:   errorMessage,
	})
}

// SendValidationError sends a 422 unprocessable entity response
func SendValidationError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
} 