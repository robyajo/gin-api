package handler

import (
	"github.com/gin-gonic/gin"
)

// APIResponse struktur standar respons JSON
type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta metadata informasi pagination
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

// PaginatedData payload data dengan pagination
type PaginatedData struct {
	Items      interface{}    `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

// SuccessResponse mengembalikan respon sukses HTTP
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// PaginatedResponse mengembalikan respon daftar berhalaman HTTP
func PaginatedResponse(c *gin.Context, code int, message string, items interface{}, page, limit, total int) {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	c.JSON(code, APIResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data: PaginatedData{
			Items: items,
			Pagination: PaginationMeta{
				CurrentPage: page,
				PerPage:     limit,
				Total:       total,
				TotalPages:  totalPages,
			},
		},
	})
}

// ErrorResponse mengembalikan respon gagal HTTP
func ErrorResponse(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, APIResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}
